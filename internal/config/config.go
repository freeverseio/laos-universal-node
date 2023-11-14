package config

import (
	"flag"
	"log/slog"
	"strings"
	"time"
)

type Config struct {
	WaitingTime      time.Duration
	StartingBlock    uint64
	EvoStartingBlock uint64
	Contracts        []string
	EvoContract      string
	Rpc              string
	EvoRpc           string
	BlocksMargin     uint
	BlocksRange      uint
	EvoBlocksMargin  uint
	EvoBlocksRange   uint
	Port             uint
	Debug            bool
}

func Load() *Config {
	blocksRange := flag.Uint("blocks_range", 100, "Amount of blocks the scanner processes")
	blocksMargin := flag.Uint("blocks_margin", 0, "Number of blocks to assume finality")
	evoBlocksRange := flag.Uint("evo_blocks_range", 100, "Amount of blocks the scanner processes on the evolution chain")
	evoBlocksMargin := flag.Uint("evo_blocks_margin", 0, "Number of blocks to assume finality on the evolution chain")
	contracts := flag.String("contracts", "", "Comma-separated list of the web3 addresses of the smart contracts to scan")
	evoContract := flag.String("evo_contract", "", "Web3 addresses of the LaosEvolution smart contract")
	debug := flag.Bool("debug", false, "Set logs to debug level")
	rpc := flag.String("rpc", "https://eth.llamarpc.com", "URL of the RPC node of an evm-compatible blockchain")
	evoRpc := flag.String("evo_rpc", "", "URL of the RPC evolution chain")
	port := flag.Uint("port", 5001, "HTTP port to use for the universal node server")
	startingBlock := flag.Uint64("starting_block", 18288287, "Initial block where the scanning process should start from")
	evoStartingBlock := flag.Uint64("evo_starting_block", 0, "Initial block where the scanning process should start from on the evolution chain")
	waitingTime := flag.Duration("wait", 5*time.Second, "Waiting time between scans when scanning reaches the last block")

	flag.Parse()

	c := &Config{
		BlocksMargin:     *blocksMargin,
		BlocksRange:      *blocksRange,
		EvoBlocksMargin:  *evoBlocksMargin,
		EvoBlocksRange:   *evoBlocksRange,
		Debug:            *debug,
		EvoContract:      *evoContract,
		Rpc:              *rpc,
		EvoRpc:           *evoRpc,
		StartingBlock:    *startingBlock,
		EvoStartingBlock: *evoStartingBlock,
		WaitingTime:      *waitingTime,
		Port:             *port,
	}

	if *contracts != "" {
		c.Contracts = strings.Split(*contracts, ",")
	}

	return c
}

func (c *Config) LogFields() {
	slog.Debug("config loaded", slog.Group("config", "rpc", c.Rpc, "evo_rpc", c.EvoRpc, "contracts", c.Contracts, "LaosEvolution conctract", c.EvoContract,
		"starting_block", c.StartingBlock, "evo_starting_block", c.EvoStartingBlock, "blocks_margin", c.BlocksMargin, "evo_blocks_margin", c.EvoBlocksMargin, "blocks_range", c.BlocksRange,
		"evo_blocks_range", c.EvoBlocksRange, "debug", c.Debug, "wait", c.WaitingTime, "port", c.Port))
}
