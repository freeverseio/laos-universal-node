package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/laos-universal-node/cmd/server"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
	"github.com/freeverseio/laos-universal-node/internal/repository"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/state/v1"
	"golang.org/x/sync/errgroup"
)

var version = "undefined"

func main() {
	if err := run(); err != nil {
		slog.Error("error occurred", "err", err)
	}
}

func run() error {
	c := config.Load()

	setLogger(c.Debug)
	c.LogFields()

	db, err := badger.Open(badger.DefaultOptions(c.Path).WithLoggingLevel(badger.ERROR))
	if err != nil {
		return err
	}
	defer func() {
		err = db.Close()
		if err != nil {
			slog.Error("error closing db", "err", err)
		}
	}()

	storageService := badgerStorage.NewService(db)
	// TODO merge repositoryService and stateService into a single service
	repositoryService := repository.New(storageService)
	stateService := v1.NewStateService(storageService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	// Badger DB garbage collection
	group.Go(func() error {
		numIterations := 3
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				// garbage collection cleans up at most 1 file per iteration
				// https://dgraph.io/docs/badger/get-started/#garbage-collection
				for i := 0; i < numIterations; i++ {
					err := db.RunValueLogGC(0.5)
					if err != nil {
						if err != badger.ErrNoRewrite {
							slog.Error("error occurred while running badger GC", "err", err.Error())
						}
						break
					}
				}
			}
		}
	})

	// Ownership chain scanner
	group.Go(func() error {
		client, err := ethclient.Dial(c.Rpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}
		if err := compareChainIDs(ctx, client, repositoryService); err != nil {
			return err
		}
		s := scan.NewScanner(client, c.Contracts...)
		return scanUniversalChain(ctx, c, client, s, repositoryService)
	})

	// Evolution chain scanner
	group.Go(func() error {
		client, err := ethclient.Dial(c.EvoRpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}
		// TODO check if chain ID match with the one in DB (call "compareChainIDs")
		s := scan.NewScanner(client)
		return scanEvoChain(ctx, c, client, s, repositoryService, stateService)
	})

	// Universal node RPC server
	group.Go(func() error {
		rpcServer, err := server.New()
		if err != nil {
			return fmt.Errorf("failed to create RPC server: %w", err)
		}
		addr := fmt.Sprintf("0.0.0.0:%v", c.Port)
		slog.Info("starting RPC server", "listen_address", addr)
		return rpcServer.ListenAndServe(ctx, c.Rpc, addr, stateService)
	})

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func compareChainIDs(ctx context.Context, client scan.EthClient, repositoryService repository.Service) error {
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("error getting chain ID from Ethereum client: %w", err)
	}

	chainIdDB, err := repositoryService.GetChainID()
	if err != nil {
		return fmt.Errorf("error getting chain ID from database: %w", err)
	}

	if chainIdDB == "" {
		if err = repositoryService.SetChainID(chainId.String()); err != nil {
			return fmt.Errorf("error setting chain ID in database: %w", err)
		}
	} else if chainId.String() != chainIdDB {
		return fmt.Errorf("mismatched chain IDs: database has %s, eth client reports %s", chainIdDB, chainId.String())
	}
	return nil
}

func shouldDiscover(repositoryService repository.Service, contracts []string) (bool, error) {
	if len(contracts) == 0 {
		return true, nil
	}
	for i := 0; i < len(contracts); i++ {
		hasContract, err := repositoryService.HasERC721UniversalContract(contracts[i])
		if err != nil {
			return false, err
		}
		if !hasContract {
			return true, nil
		}
	}
	return false, nil
}

func scanUniversalChain(ctx context.Context, c *config.Config, client scan.EthClient, s scan.Scanner, repositoryService repository.Service) error {
	startingBlockDB, err := repositoryService.GetCurrentBlock()
	if err != nil {
		return fmt.Errorf("error retrieving the current block from storage: %w", err)
	}
	startingBlock, err := getStartingBlock(ctx, startingBlockDB, c.StartingBlock, client)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return nil
		default:
			l1LatestBlock, err := getL1LatestBlock(ctx, client)
			if err != nil {
				slog.Error("error retrieving the latest block", "err", err.Error())
				break
			}
			lastBlock := calculateLastBlock(startingBlock, l1LatestBlock, c.BlocksRange, c.BlocksMargin)
			if lastBlock < startingBlock {
				slog.Debug("last calculated block is behind starting block, waiting...",
					"last_block", lastBlock, "starting_block", startingBlock)
				waitBeforeNextScan(ctx, c.WaitingTime)
				break
			}

			shouldDiscover, err := shouldDiscover(repositoryService, c.Contracts)
			if err != nil {
				slog.Error("error occurred reading contracts from storage", "err", err.Error())
				break
			}

			var universalContracts []model.ERC721UniversalContract

			if shouldDiscover {
				universalContracts, err = s.ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)))
				if err != nil {
					slog.Error("error occurred while discovering new universal events", "err", err.Error())
					break
				}

				if err = repositoryService.StoreERC721UniversalContracts(universalContracts); err != nil {
					slog.Error("error occurred while storing universal contract(s)", "err", err.Error())
					break
				}
			}

			var contracts []string
			if len(c.Contracts) > 0 {
				contracts = c.Contracts
			} else {
				contracts, err = repositoryService.GetAllERC721UniversalContracts()
				if err != nil {
					slog.Error("error occurred reading contracts from storage", "err", err.Error())
					break
				}
			}
			var lastScannedBlock *big.Int
			if len(contracts) > 0 {
				_, lastScannedBlock, err = s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), contracts)
				if err != nil {
					slog.Error("error occurred while scanning events", "err", err.Error())
					break
				}
			} else {
				lastScannedBlock = big.NewInt(int64(lastBlock))
			}
			nextStartingBlock := lastScannedBlock.Uint64() + 1
			if err = repositoryService.SetCurrentBlock(strconv.FormatUint(nextStartingBlock, 10)); err != nil {
				slog.Error("error occurred while storing current block", "err", err.Error())
				break
			}
			startingBlock = nextStartingBlock
		}
	}
}

func scanEvoChain(ctx context.Context, c *config.Config, client scan.EthClient, s scan.Scanner, repositoryService repository.Service, stateService state.Service) error {
	startingBlockDB, err := repositoryService.GetEvoChainCurrentBlock()
	if err != nil {
		return fmt.Errorf("error retrieving the current block from storage: %w", err)
	}
	startingBlock, err := getStartingBlock(ctx, startingBlockDB, c.EvoStartingBlock, client)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return nil
		default:
			l1LatestBlock, err := getL1LatestBlock(ctx, client)
			if err != nil {
				slog.Error("error retrieving the latest block", "err", err.Error())
				break
			}
			lastBlock := calculateLastBlock(startingBlock, l1LatestBlock, c.EvoBlocksRange, c.EvoBlocksMargin)
			if lastBlock < startingBlock {
				slog.Debug("last calculated block is behind starting block, waiting...",
					"last_block", lastBlock, "starting_block", startingBlock)
				waitBeforeNextScan(ctx, c.WaitingTime)
				break
			}

			events, lastScannedBlock, err := s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), nil)
			if err != nil {
				slog.Error("error occurred while scanning LaosEvolution events", "err", err.Error())
				break
			}

			groupMintEvents := make(map[common.Address][]model.EventMintedWithExternalURI, 0)
			for _, event := range events {
				if e, ok := event.(scan.EventMintedWithExternalURI); ok {
					groupMintEvents[e.Contract] = append(groupMintEvents[e.Contract], model.EventMintedWithExternalURI{
						Slot:        e.Slot,
						To:          e.To,
						TokenURI:    e.TokenURI,
						TokenId:     e.TokenId,
						BlockNumber: e.BlockNumber,
						Timestamp:   e.Timestamp,
					})
				}
			}

			tx := stateService.NewTransaction()
			// nolint: gocritic // TODO to address this linting suggestion (deferInLoop), probably the whole body of the "default" case must be moved to a separate function
			defer tx.Discard()

			for contract, scannedEvents := range groupMintEvents {
				// fetch current storedEvents stord for this specific contract address
				events := make([]model.EventMintedWithExternalURI, 0)
				storedEvents, err := tx.EvoChainMintEvents(contract)
				if err != nil {
					slog.Error("error occurred while reading database", "err", err.Error())
					break
				}
				if storedEvents != nil {
					events = append(events, storedEvents...)
				}
				events = append(events, scannedEvents...)
				if err := tx.StoreEvoChainMintEvents(contract, events); err != nil {
					slog.Error("error occurred while writing events to database", "err", err.Error())
					break
				}
			}

			// TODO remember to handle SetEvoChainCurrentBlock and the future SetState of the merkle tree in the same TX
			nextStartingBlock := lastScannedBlock.Uint64() + 1
			if err = repositoryService.SetEvoChainCurrentBlock(strconv.FormatUint(nextStartingBlock, 10)); err != nil {
				slog.Error("error occurred while storing current block", "err", err.Error())
				break
			}
			startingBlock = nextStartingBlock
		}
	}
}

func getStartingBlock(ctx context.Context, startingBlockDB string, configStartingBlock uint64, client scan.EthClient) (uint64, error) {
	var startingBlock uint64
	var err error
	if startingBlockDB != "" {
		startingBlock, err = strconv.ParseUint(startingBlockDB, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing the current block from storage: %w", err)
		}
		slog.Debug("ignoring user provided starting block, using last updated block from storage", "starting_block", startingBlock)
	}

	if startingBlock == 0 {
		startingBlock = configStartingBlock
		if startingBlock == 0 {
			startingBlock, err = getL1LatestBlock(ctx, client)
			if err != nil {
				return 0, fmt.Errorf("error retrieving the latest block from chain: %w", err)
			}
			slog.Debug("latest block found", "latest_block", startingBlock)
		}
	}
	return startingBlock, nil
}

func waitBeforeNextScan(ctx context.Context, waitingTime time.Duration) {
	timer := time.NewTimer(waitingTime)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}

func getL1LatestBlock(ctx context.Context, client scan.EthClient) (uint64, error) {
	l1LatestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return l1LatestBlock, nil
}

func calculateLastBlock(startingBlock, l1LatestBlock uint64, blocksRange, blocksMargin uint) uint64 {
	return min(startingBlock+uint64(blocksRange), l1LatestBlock-uint64(blocksMargin))
}

func setLogger(debug bool) {
	// Default slog.Level is Info (0)
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("version", version),
	}))
	slog.SetDefault(logger)
}
