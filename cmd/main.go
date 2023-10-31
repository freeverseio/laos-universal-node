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

	group.Go(func() error {
		client, err := ethclient.Dial(c.Rpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}
		s := scan.NewScanner(client, common.HexToAddress(c.ContractAddress))
		return runScan(ctx, *c, client, s)
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

func runScan(ctx context.Context, c config.Config, client scan.EthClient, s scan.Scanner) error {
	var err error
	if c.StartingBlock == 0 {
		c.StartingBlock, err = getL1LatestBlock(ctx, client)
		if err != nil {
			return fmt.Errorf("error retrieving the latest block: %w", err)
		}
		slog.Debug("latest block found", "latest_block", c.StartingBlock)
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
			lastBlock := calculateLastBlock(c.StartingBlock, l1LatestBlock, c.BlocksRange, c.BlocksMargin)
			if lastBlock < c.StartingBlock {
				slog.Debug("last calculated block is behind starting block, waiting...",
					"last_block", lastBlock, "starting_block", c.StartingBlock)
				waitBeforeNextScan(ctx, c.WaitingTime)
				break
			}
			_, err = s.ScanEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(int64(lastBlock)))
			if err != nil {
				slog.Error("error occurred while scanning events", "err", err.Error())
				break
			}
			c.StartingBlock = lastBlock + 1
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
