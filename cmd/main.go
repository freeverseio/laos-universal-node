package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
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

const historyLength = 256

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
		return fmt.Errorf("error initializing storage: %w", err)
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
		return scanUniversalChain(ctx, c, client, s, stateService)
	})

	// Evolution chain scanner
	group.Go(func() error {
		client, err := ethclient.Dial(c.EvoRpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}
		// TODO check if chain ID match with the one in DB (call "compareChainIDs")
		s := scan.NewScanner(client)
		return scanEvoChain(ctx, c, client, s, stateService)
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

func shouldDiscover(tx state.Tx, contracts []string) (bool, error) {
	if len(contracts) == 0 {
		return true, nil
	}
	for i := 0; i < len(contracts); i++ {
		hasContract, err := tx.HasERC721UniversalContract(contracts[i])
		if err != nil {
			return false, err
		}
		if !hasContract {
			return true, nil
		}
	}
	return false, nil
}

func scanUniversalChain(ctx context.Context, c *config.Config, client scan.EthClient,
	s scan.Scanner, stateService state.Service,
) error {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	startingBlockDB, err := tx.GetCurrentOwnershipBlock()
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

			tx = stateService.NewTransaction()

			// discovering new contracts deployed on the ownership chain
			shouldDiscover, err := shouldDiscover(tx, c.Contracts)
			if err != nil {
				slog.Error("error occurred reading contracts from storage", "err", err.Error())
				break
			}
			if shouldDiscover {
				if err = discoverContracts(ctx, s, startingBlock, lastBlock, tx); err != nil {
					break
				}
			}

			var contractsAddress []string
			// choosing which contracts to scan
			if len(c.Contracts) > 0 {
				// if contracts come from flag, consider only those that have been discovered (whether in this iteration or previously)
				var existingContracts []string
				existingContracts, err = tx.GetExistingERC721UniversalContracts(c.Contracts)
				if err != nil {
					slog.Error("error occurred checking if user-provided contracts exist in storage", "err", err.Error())
					break
				}
				contractsAddress = append(contractsAddress, existingContracts...)
			} else {
				dbContracts := tx.GetAllERC721UniversalContracts()
				contractsAddress = append(contractsAddress, dbContracts...)
			}

			// load merkle trees for all contracts whose events have to be scanned for
			if err = loadMerkleTrees(tx, contractsAddress); err != nil {
				slog.Error("error creating merkle trees", "err", err)
				break
			}

			// scanning contracts for events on the ownership chain
			var lastScannedBlock *big.Int
			if len(contractsAddress) > 0 {
				var scanEvents []scan.Event
				scanEvents, lastScannedBlock, err = s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), contractsAddress)
				if err != nil {
					slog.Error("error occurred while scanning events", "err", err.Error())
					break
				}

				// getting transfer events from scan events
				var modelTransferEvents map[string][]model.ERC721Transfer
				modelTransferEvents, err = getModelTransferEvents(ctx, client, scanEvents)
				if err != nil {
					slog.Error("error parsing transfer events", "err", err.Error())
					break
				}

				// retrieving minted events and update the state accordingly
				if err = readEventsAndUpdateState(client, contractsAddress, modelTransferEvents, tx); err != nil {
					slog.Error("error occurred", "err", err.Error())
					break
				}
			} else {
				lastScannedBlock = big.NewInt(int64(lastBlock))
			}

			nextStartingBlock := lastScannedBlock.Uint64() + 1

			if err = tagRootsUntilBlock(tx, contractsAddress, nextStartingBlock); err != nil {
				slog.Error("error occurred while tagging roots", "err", err.Error())
				break
			}

			if err = tx.SetCurrentOwnershipBlock(nextStartingBlock); err != nil {
				slog.Error("error occurred while storing current block", "err", err.Error())
				break
			}
			if err = tx.Commit(); err != nil {
				slog.Error("error occurred while committing transaction", "err", err.Error())
				break
			}
			startingBlock = nextStartingBlock
		}
	}
}

func scanEvoChain(ctx context.Context, c *config.Config, client scan.EthClient, s scan.Scanner, stateService state.Service) error {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	startingBlockDB, err := tx.GetCurrentEvoBlock()
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
			tx = stateService.NewTransaction()
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

			_, lastScannedBlock, err := s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), nil)
			if err != nil {
				slog.Error("error occurred while scanning LaosEvolution events", "err", err.Error())
				break
			}

			nextStartingBlock := lastScannedBlock.Uint64() + 1
			if err = tx.SetCurrentEvoBlock(nextStartingBlock); err != nil {
				slog.Error("error occurred while storing current block", "err", err.Error())
				break
			}
			if err = tx.Commit(); err != nil {
				slog.Error("error occurred while committing transaction", "err", err.Error())
				break
			}
			startingBlock = nextStartingBlock
		}
	}
}

func getStartingBlock(ctx context.Context, startingBlockDB, configStartingBlock uint64, client scan.EthClient) (uint64, error) {
	var startingBlock uint64
	var err error
	if startingBlockDB != 0 {
		startingBlock = startingBlockDB
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

func readEventsAndUpdateState(client scan.EthClient, contractsAddress []string, modelTransferEvents map[string][]model.ERC721Transfer, tx state.Tx) error {
	for i := range contractsAddress {
		var mintedEvents []model.MintedWithExternalURI
		collectionAddress, err := tx.GetCollectionAddress(contractsAddress[i]) // get collection address from ownership address
		if err != nil {
			return fmt.Errorf("error occurred retrieving the collection address from the ownership contract %s: %w", contractsAddress[i], err)
		}

		// now we get the minted events from the evolution chain for the collection address
		mintedEvents, err = tx.GetMintedWithExternalURIEvents(collectionAddress.String())
		if err != nil {
			return fmt.Errorf("error occurred retrieving evochain minted events for ownership contract %s and collection address %s: %w",
				contractsAddress[i], collectionAddress.String(), err)
		}

		var evoBlock uint64
		evoBlock, err = updateState(client, mintedEvents, modelTransferEvents[contractsAddress[i]], contractsAddress[i], tx)
		if err != nil {
			return fmt.Errorf("error updating state: %w", err)
		}

		if err = tx.SetCurrentEvoBlockForOwnershipContract(contractsAddress[i], evoBlock); err != nil {
			return fmt.Errorf("error updating current evochain block %d for ownership contract %s: %w", evoBlock, contractsAddress[i], err)
		}
	}
	return nil
}

func tagRootsUntilBlock(tx state.Tx, contractsAddress []string, blockNumber uint64) error {
	for i := range contractsAddress {
		lastTaggedBlock, err := tx.GetLastTaggedBlock(common.HexToAddress(contractsAddress[i]))
		if err != nil {
			return err
		}
		for block := lastTaggedBlock + 1; block < int64(blockNumber); block++ {
			if err := tx.TagRoot(common.HexToAddress(contractsAddress[i]), block); err != nil {
				return err
			}

			if err := tx.DeleteRootTag(common.HexToAddress(contractsAddress[i]), block-historyLength); err != nil {
				return err
			}
		}
	}

	return nil
}

func loadMerkleTrees(tx state.Tx, contractsAddress []string) error {
	for i := range contractsAddress {
		if err := loadMerkleTree(tx, common.HexToAddress(contractsAddress[i])); err != nil {
			return err
		}
	}
	return nil
}

func discoverContracts(ctx context.Context, s scan.Scanner, startingBlock, lastBlock uint64, tx state.Tx) error {
	contracts, err := s.ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred while discovering new universal events", "err", err.Error())
		return err
	}

	if len(contracts) > 0 {
		if err = tx.StoreERC721UniversalContracts(contracts); err != nil {
			slog.Error("error occurred while storing universal contract(s)", "err", err.Error())
			return err
		}
	}

	for i := range contracts {
		if err = loadMerkleTree(tx, contracts[i].Address); err != nil {
			slog.Error("error creating merkle trees for newly discovered universal contract(s)", "err", err)
			return err
		}

		if err = tx.TagRoot(contracts[i].Address, int64(contracts[i].BlockNumber)); err != nil {
			slog.Error("error occurred tagging roots for newly discovered universal contract(s)", "err", err.Error())
			return err
		}
	}

	return nil
}

func loadMerkleTree(tx state.Tx, contractAddress common.Address) error {
	if !tx.IsTreeSetForContract(contractAddress) {
		ownership, enumerated, enumeratedTotal, err := tx.CreateTreesForContract(contractAddress)
		if err != nil {
			return err
		}
		tx.SetTreesForContract(contractAddress, ownership, enumerated, enumeratedTotal)
	}
	return nil
}

func updateState(client scan.EthClient, mintedEvents []model.MintedWithExternalURI, modelTransferEvents []model.ERC721Transfer, contract string, tx state.Tx) (uint64, error) {
	ownershipContractEvoBlock, err := tx.GetCurrentEvoBlockForOwnershipContract(contract)
	if err != nil {
		return 0, err
	}
	var mintedIndex int
	var transferIndex int
	lastUpdatedBlock := ownershipContractEvoBlock

	for {
		switch {
		// all events have been processed => return
		case mintedIndex >= len(mintedEvents) && transferIndex >= len(modelTransferEvents):
			return lastUpdatedBlock, nil
		// all minted events have been processed => process remaining transfer events
		case mintedIndex >= len(mintedEvents):
			if err := updateStateWithTransfer(contract, tx, &modelTransferEvents[transferIndex]); err != nil {
				return 0, err
			}
			transferIndex++
		// all transfer events have been processed => process remaining minted events
		case transferIndex >= len(modelTransferEvents):
			block, err := updateStateWithMint(client, contract, tx, &mintedEvents[mintedIndex], ownershipContractEvoBlock)
			if err != nil {
				return 0, err
			}
			lastUpdatedBlock = block
			mintedIndex++
		default:
			// if minted event's timestamp is behind transfer event's timestamp => process minted event
			if mintedEvents[mintedIndex].Timestamp < modelTransferEvents[transferIndex].Timestamp {
				block, err := updateStateWithMint(client, contract, tx, &mintedEvents[mintedIndex], ownershipContractEvoBlock)
				if err != nil {
					return 0, err
				}
				lastUpdatedBlock = block
				mintedIndex++
			} else {
				if err := updateStateWithTransfer(contract, tx, &modelTransferEvents[transferIndex]); err != nil {
					return 0, err
				}
				transferIndex++
			}
		}
	}
}

func updateStateWithTransfer(contract string, tx state.Tx, modelTransferEvent *model.ERC721Transfer) error {
	// TODO if transfer event's timestamp is ahead of global evo chain current block's timestamp => wait X seconds and read again the global evo chain current block from DB
	// we must wait because there might have been mint events whose timestamp is < this transfer event
	// maybe it is worth storing the global evo chain current block's timestamp also?!

	lastTaggedBlock, err := tx.GetLastTaggedBlock(common.HexToAddress(contract))
	if err != nil {
		return err
	}
	for block := lastTaggedBlock + 1; block < int64(modelTransferEvent.BlockNumber); block++ {
		if err := tx.TagRoot(common.HexToAddress(contract), block); err != nil {
			return err
		}

		if err := tx.DeleteRootTag(common.HexToAddress(contract), block-historyLength); err != nil {
			return err
		}
	}

	if err := tx.Transfer(common.HexToAddress(contract), modelTransferEvent); err != nil {
		return fmt.Errorf("error updating transfer state for contract %s and token id %d, from %s, to %s: %w",
			contract, modelTransferEvent.TokenId,
			modelTransferEvent.From, modelTransferEvent.To, err)
	}
	return nil
}

func updateStateWithMint(client scan.EthClient, contract string, tx state.Tx, mintedEvent *model.MintedWithExternalURI, ownershipContractEvoBlock uint64) (uint64, error) {
	updatedBlock := ownershipContractEvoBlock
	// TODO check if this is correct. Could it be that on a early termination (ctrl + c), some events on ownershipContractEvoBlock are not stored in the state?
	// if so, on the next iteration, those events won't be stored in the state because of the ">" comparison
	// example:
	// 2 events on block 10, you store the first, updatedBlock == 10. a ctrl + c comes, you don't store the second event and the method returns
	// can it be that the transaction is committed?
	lastTaggedBlock, err := tx.GetLastTaggedBlock(common.HexToAddress(contract))
	if err != nil {
		return 0, err
	}

	for {
		blockToTag := lastTaggedBlock + 1
		timestamp, err := getTimestampForBlockNumber(context.Background(), client, uint64(blockToTag))
		if err != nil {
			return 0, err
		}
		if timestamp > mintedEvent.Timestamp {
			break
		}
		if err := tx.TagRoot(common.HexToAddress(contract), blockToTag); err != nil {
			return 0, err
		}
		if err := tx.DeleteRootTag(common.HexToAddress(contract), blockToTag-historyLength); err != nil {
			return 0, err
		}
		// update lastTaggedBlock
		lastTaggedBlock = blockToTag
	}

	if mintedEvent.BlockNumber > ownershipContractEvoBlock {
		if err := tx.Mint(common.HexToAddress(contract), mintedEvent.TokenId); err != nil {
			return 0, fmt.Errorf("error updating mint state for contract %s and token id %d: %w",
				contract, mintedEvent.TokenId, err)
		}
		updatedBlock = mintedEvent.BlockNumber
	}

	return updatedBlock, nil
}

func getModelTransferEvents(ctx context.Context, client scan.EthClient, scanEvents []scan.Event) (map[string][]model.ERC721Transfer, error) {
	modelTransferEvents := make(map[string][]model.ERC721Transfer)
	for i := range scanEvents {
		if scanEvent, ok := scanEvents[i].(scan.EventTransfer); ok {
			timestamp, err := getTimestampForBlockNumber(ctx, client, scanEvent.BlockNumber)
			if err != nil {
				return nil, fmt.Errorf("error retrieving timestamp for block number %d: %w", scanEvent.BlockNumber, err)
			}
			contractString := scanEvent.Contract.String()
			modelTransferEvents[contractString] = append(modelTransferEvents[contractString], model.ERC721Transfer{
				From:        scanEvent.From,
				To:          scanEvent.To,
				TokenId:     scanEvent.TokenId,
				BlockNumber: scanEvent.BlockNumber,
				Contract:    scanEvent.Contract,
				Timestamp:   timestamp,
			})
		}
	}
	return modelTransferEvents, nil
}

func getTimestampForBlockNumber(ctx context.Context, client scan.EthClient, blockNumber uint64) (uint64, error) {
	header, err := client.HeaderByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return 0, err
	}
	return header.Time, err
}
