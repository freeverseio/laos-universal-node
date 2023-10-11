package main

import (
	"context"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/laos-universal-node/config"
	"github.com/freeverseio/laos-universal-node/scanner"
)

var (
	version = "undefined"
)

func main() {
	c := config.Load()

	setLogger(c.Debug)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	cli, err := ethclient.Dial(c.Rpc)
	if err != nil {
		slog.Error("error instantiating eth client", "err", err.Error())
		os.Exit(1)
	}

	contract := common.HexToAddress(c.ContractAddress)
	if c.StartingBlock == 0 {
		c.StartingBlock, err = getL1LatestBlock(ctx, cli)
		if err != nil {
			slog.Error("error retrieving the latest block", "err", err.Error())
			os.Exit(1)
		}
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			l1LatestBlock, err := getL1LatestBlock(ctx, cli)
			if err != nil {
				slog.Error("error retrieving the latest block", "err", err.Error())
				break
			}
			lastBlock := calculateLastBlock(c.StartingBlock, l1LatestBlock, c.BlocksRange, c.BlocksMargin)
			if lastBlock < c.StartingBlock {
				slog.Debug("last calculated block is behind starting block, continue...")
				break
			}
			_, err = scanner.ScanEvents(cli, contract, big.NewInt(int64(c.StartingBlock)), big.NewInt(int64(lastBlock)))
			if err != nil {
				slog.Error("error occurred while scanning events", "err", err.Error())
				break
			}
			c.StartingBlock = lastBlock + 1
		}
	}
}

func getL1LatestBlock(ctx context.Context, cli *ethclient.Client) (uint64, error) {
	lastBlock, err := cli.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return lastBlock, nil
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
