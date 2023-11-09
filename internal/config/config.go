package config

import (
	"flag"
	"log/slog"
	"os"
	"path"
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
	defaultStoragePath := getDefaultStoragePath()

	blocksRange := flag.Uint("blocks_range", 100, "Amount of blocks the scanner processes")
	blocksMargin := flag.Uint("blocks_margin", 0, "Number of blocks to assume finality")
	contracts := flag.String("contracts", "", "Comma-separated list of the web3 addresses of the smart contracts to scan")
	debug := flag.Bool("debug", false, "Set logs to debug level")
	rpc := flag.String("rpc", "https://eth.llamarpc.com", "URL of the RPC node of an evm-compatible blockchain")
	port := flag.Uint("port", 5001, "HTTP port to use for the universal node server")
	startingBlock := flag.Uint64("starting_block", 18288287, "Initial block where the scanning process should start from")
	waitingTime := flag.Duration("wait", 5*time.Second, "Waiting time between scans when scanning reaches the last block")
	storagePath := flag.String("storage_path", defaultStoragePath, "Path to the storage folder")

	flag.Parse()

	c := &Config{
		BlocksMargin:  *blocksMargin,
		BlocksRange:   *blocksRange,
		Debug:         *debug,
		Rpc:           *rpc,
		StartingBlock: *startingBlock,
		WaitingTime:   *waitingTime,
		Port:          *port,
		Path:          *storagePath,
	}

	if *contracts != "" {
		c.Contracts = strings.Split(*contracts, ",")
	}

	return c
}

func (c *Config) LogFields() {
	slog.Debug("config loaded", slog.Group("config", "rpc", c.Rpc, "contracts", c.Contracts,
		"starting_block", c.StartingBlock, "blocks_margin", c.BlocksMargin, "blocks_range", c.BlocksRange,
		"debug", c.Debug, "wait", c.WaitingTime, "port", c.Port, "storage_path", c.Path))
}

func getDefaultStoragePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Warn("user home directory not found, default storage path will be under the current directory", "err", err)
		homeDir = "./"
	}
	return path.Join(homeDir, ".universalnode")
}
