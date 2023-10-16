package config

import (
	"flag"
	"log/slog"
	"time"
)

type Config struct {
	WaitingTime     time.Duration
	StartingBlock   uint64
	ContractAddress string
	Rpc             string
	BlocksMargin    uint
	BlocksRange     uint
	Port            uint
	Debug           bool
}

func Load() *Config {
	// 0xBC4... is the Bored Ape ERC721 Ethereum contract
	blocksRange := flag.Uint("blocks_range", 100, "Amount of blocks the scanner processes")
	blocksMargin := flag.Uint("blocks_margin", 70, "Number of blocks to assume finality")
	contractAddress := flag.String("contract", "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D", "Web3 address of the smart contract")
	debug := flag.Bool("debug", false, "Set logs to debug level")
	rpc := flag.String("rpc", "https://eth.llamarpc.com", "URL of the RPC node of an evm-compatible blockchain")
	port := flag.Uint("port", 5001, "HTTP port to use for the universal node server")
	startingBlock := flag.Uint64("starting_block", 18288287, "Initial block where the scanning process should start from")
	waitingTime := flag.Duration("wait", 5*time.Second, "Waiting time between scans when scanning reaches the last block")

	flag.Parse()

	c := &Config{
		BlocksMargin:    *blocksMargin,
		BlocksRange:     *blocksRange,
		ContractAddress: *contractAddress,
		Debug:           *debug,
		Rpc:             *rpc,
		StartingBlock:   *startingBlock,
		WaitingTime:     *waitingTime,
		Port:            *port,
	}

	return c
}

func (c *Config) LogFields() {
	slog.Debug("config loaded", slog.Group("config", "rpc", c.Rpc, "contract", c.ContractAddress,
		"starting_block", c.StartingBlock, "blocks_margin", c.BlocksMargin, "blocks_range", c.BlocksRange,
		"debug", c.Debug, "wait", c.WaitingTime, "port", c.Port))
}
