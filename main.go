package main

import (
	"context"
	"flag"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/laos-universal-node/scanner"
)

var (
	version = "undefined"
)

func main() {
	setLogger()

	// TODO move flags to config?
	rpc := flag.String("rpc", "https://eth.llamarpc.com", "URL of the RPC node of an evm-compatible blockchain")
	// 0xBC4... is the Bored Ape ERC721 Ethereum contract
	contractAddress := flag.String("contract", "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D", "Web3 address of the smart contract")
	startingBlock := flag.Uint64("starting_block", 18288287, "Initial block where the scanning process should start from")
	blocksRange := flag.Uint("block_range", 100, "Amount of blocks the scanner processes")
	blocksMargin := flag.Uint("block_margin", 70, "Number of blocks to assume finality")
	flag.Parse()

	cli, err := ethclient.Dial(*rpc)
	if err != nil {
		slog.Error("error instantiating eth client", "err", err.Error())
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	contract := common.HexToAddress(*contractAddress)
	if *startingBlock == 0 {
		*startingBlock, err = getL1LatestBlock(cli, ctx)
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
			l1LatestBlock, err := getL1LatestBlock(cli, ctx)
			if err != nil {
				slog.Error("error retrieving the latest block", "err", err.Error())
				break
			}
			lastBlock := calculateLastBlock(*startingBlock, l1LatestBlock, *blocksRange, *blocksMargin)
			if lastBlock < *startingBlock {
				slog.Debug("last calculated block is behind starting block, continue...")
				break
			}
			_, err = scanner.ScanEvents(cli, contract, big.NewInt(int64(*startingBlock)), big.NewInt(int64(lastBlock)))
			if err != nil {
				slog.Error("error occurred while scanning events", "err", err.Error())
				break
			}
			*startingBlock = lastBlock + 1
		}
	}
}

func getL1LatestBlock(cli *ethclient.Client, ctx context.Context) (uint64, error) {
	lastBlock, err := cli.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return lastBlock, nil
}

func calculateLastBlock(startingBlock, l1LatestBlock uint64, blocksRange, blocksMargin uint) uint64 {
	return min(startingBlock+uint64(blocksRange), l1LatestBlock-uint64(blocksMargin))
}

func setLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("version", version),
	}))
	slog.SetDefault(logger)
}
