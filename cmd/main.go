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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/laos-universal-node/cmd/server"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
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

	db, err := badger.Open(badger.DefaultOptions(c.Path))
	if err != nil {
		return err
	}
	storageService := storage.New(db)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		client, err := ethclient.Dial(c.Rpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}
		scanStorage, err := scan.NewFSStorage("erc721_contracts.txt")
		if err != nil {
			return fmt.Errorf("error initializing storage: %w", err)
		}
		s := scan.NewScanner(client, scanStorage, c.Contracts...)
		return runScan(ctx, c, client, s, storageService)
	})

	group.Go(func() error {
		rpcServer, err := server.New()
		if err != nil {
			return fmt.Errorf("failed to create RPC server: %w", err)
		}
		addr := fmt.Sprintf("0.0.0.0:%v", c.Port)
		slog.Info("starting RPC server", "listen_address", addr)
		return rpcServer.ListenAndServe(ctx, c.Rpc, addr)
	})

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func runScan(ctx context.Context, c *config.Config, client scan.EthClient, s scan.Scanner, storageService storage.Storage) error {
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

			if err = func() error {
				tx := storageService.NewTransaction()
				defer tx.Discard()
				for i := 0; i < len(universalContracts); i++ {
					var value []byte
					value, err = model.MarshalERC721UniversalContract(universalContracts[i].BaseURI, universalContracts[i].Block)
					if err != nil {
						return err
					}
					err = tx.Set(universalContracts[i].Address[:], value)
					if err != nil {
						return err
					}
				}
				err = tx.Commit()
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				slog.Error("error occurred while storing universal contract", "err", err.Error())
				break
			}

			// TODO when the DB is in use, storage.ReadAll will run here and not inside ScanEvents
			_, err = s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)))
			if err != nil {
				slog.Error("error occurred while scanning events", "err", err.Error())
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
