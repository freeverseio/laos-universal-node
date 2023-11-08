package config

import (
	"flag"
	"log/slog"
	"strings"
	"time"
)

type Config struct {
	WaitingTime   time.Duration
	StartingBlock uint64
	Contracts     []string
	Rpc           string
	Path          string
	BlocksMargin  uint
	BlocksRange   uint
	Port          uint
	Debug         bool
}

func Load() *Config {
	// 0xBC4... is the Bored Ape ERC721 Ethereum contract
	blocksRange := flag.Uint("blocks_range", 100, "Amount of blocks the scanner processes")
	blocksMargin := flag.Uint("blocks_margin", 0, "Number of blocks to assume finality")
	contracts := flag.String("contracts", "", "Comma-separated list of the web3 addresses of the smart contracts to scan")
	debug := flag.Bool("debug", false, "Set logs to debug level")
	rpc := flag.String("rpc", "https://eth.llamarpc.com", "URL of the RPC node of an evm-compatible blockchain")
	port := flag.Uint("port", 5001, "HTTP port to use for the universal node server")
	startingBlock := flag.Uint64("starting_block", 18288287, "Initial block where the scanning process should start from")
	waitingTime := flag.Duration("wait", 5*time.Second, "Waiting time between scans when scanning reaches the last block")
	path := flag.String("db_path", "~/.universalnode", "Path to the DB folder")

	flag.Parse()

	c := &Config{
		BlocksMargin:  *blocksMargin,
		BlocksRange:   *blocksRange,
		Debug:         *debug,
		Rpc:           *rpc,
		StartingBlock: *startingBlock,
		WaitingTime:   *waitingTime,
		Port:          *port,
		Path:          *path,
	}

	if *contracts != "" {
		c.Contracts = strings.Split(*contracts, ",")
	}

	return c
}

func (c *Config) LogFields() {
	slog.Debug("config loaded", slog.Group("config", "rpc", c.Rpc, "contracts", c.Contracts,
		"starting_block", c.StartingBlock, "blocks_margin", c.BlocksMargin, "blocks_range", c.BlocksRange,
		"debug", c.Debug, "wait", c.WaitingTime, "port", c.Port, "db_path", c.Path))
}
