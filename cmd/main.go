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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/laos-universal-node/cmd/server"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/scan"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	storage, err := scan.NewFSStorage("erc721_contracts.txt")
	if err != nil {
		return fmt.Errorf("error initializing storage: %w", err)
	}

	// ERC721 Universal scanner
	group.Go(func() error {
		client, err := ethclient.Dial(c.Rpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}

		s := scan.NewScanner(client, storage)
		return scanUniversalChain(ctx, c, client, s, storage)
	})

	// Evo scanner
	group.Go(func() error {
		client, err := ethclient.Dial(c.EvoRpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}
		s := scan.NewScanner(client, nil)
		return scanEvoChain(ctx, c, client, s)
	})

	// Universal noed RPC server
	group.Go(func() error {
		rpcServer, err := server.New()
		if err != nil {
			return fmt.Errorf("failed to create RPC server: %w", err)
		}
		addr := fmt.Sprintf("0.0.0.0:%v", c.Port)
		slog.Info("starting RPC server", "listen_address", addr)
		return rpcServer.ListenAndServe(ctx, c.Rpc, addr, storage)
	})

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func scanUniversalChain(ctx context.Context, c *config.Config, client scan.EthClient, s scan.Scanner, storage scan.Storage) error {
	var err error
	startingBlock := c.StartingBlock
	if startingBlock == 0 {
		startingBlock, err = getL1LatestBlock(ctx, client)
		if err != nil {
			return fmt.Errorf("error retrieving the latest block: %w", err)
		}
		slog.Debug("latest block found", "latest_block", startingBlock)
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

			universalContracts, err := s.ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)))
			if err != nil {
				slog.Error("error occurred while discovering new universal events", "err", err.Error())
				break
			}

			// This will be replaced by a batch write to the DB
			for i := 0; i < len(universalContracts); i++ {
				if err = storage.Store(ctx, universalContracts[i]); err != nil {
					slog.Error("error occurred while storing universal contract", "err", err.Error())
					break
				}
			}

			erc721contracts, err := storage.ReadAll(context.Background())
			if err != nil {
				slog.Error("error reading contracts from storage", "err", err.Error())
				break
			}

			if len(erc721contracts) == 0 {
				slog.Debug("no contracts found")
				break
			}

			contracts := make([]common.Address, 0)
			for _, c := range erc721contracts {
				contracts = append(contracts, c.Address)
			}

			_, err = s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), contracts...)
			if err != nil {
				slog.Error("error occurred while scanning events", "err", err.Error())
				break
			}
			startingBlock = lastBlock + 1
		}
	}
}

func scanEvoChain(ctx context.Context, c *config.Config, client scan.EthClient, s scan.Scanner) error {
	var err error
	startingBlock := c.EvoStartingBlock
	if startingBlock == 0 {
		startingBlock, err = getL1LatestBlock(ctx, client)
		if err != nil {
			return fmt.Errorf("error retrieving the latest block: %w", err)
		}
		slog.Debug("latest block found", "latest_block", startingBlock)
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

			_, err = s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), common.HexToAddress(c.EvoContract))
			if err != nil {
				slog.Error("error occurred while scanning LaosEvolution events", "err", err.Error())
				break
			}

			startingBlock = lastBlock + 1
		}
	}
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
