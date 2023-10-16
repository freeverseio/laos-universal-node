package config

import (
	"flag"

	"log/slog"
)

type Config struct {
	BlocksMargin    uint
	BlocksRange     uint
	ContractAddress string
	Debug           bool
	Rpc             string
	StartingBlock   uint64
}

func Load() *Config {
	// 0xBC4... is the Bored Ape ERC721 Ethereum contract
	blocksRange := flag.Uint("blocks_range", 100, "Amount of blocks the scanner processes")
	blocksMargin := flag.Uint("blocks_margin", 70, "Number of blocks to assume finality")
	contractAddress := flag.String("contract", "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D", "Web3 address of the smart contract")
	debug := flag.Bool("debug", false, "Set logs to debug level")
	rpc := flag.String("rpc", "https://eth.llamarpc.com", "URL of the RPC node of an evm-compatible blockchain")
	startingBlock := flag.Uint64("starting_block", 18288287, "Initial block where the scanning process should start from")

	flag.Parse()

	c := &Config{
		BlocksMargin:    *blocksMargin,
		BlocksRange:     *blocksRange,
		ContractAddress: *contractAddress,
		Debug:           *debug,
		Rpc:             *rpc,
		StartingBlock:   *startingBlock,
	}

	return c
}

func (c *Config) LogFields() {
	slog.Debug("config loaded", slog.Group("config", "rpc", c.Rpc, "contract", c.ContractAddress,
		"starting_block", c.StartingBlock, "blocks_margin", c.BlocksMargin, "blocks_range", c.BlocksRange, "debug", c.Debug))
}
