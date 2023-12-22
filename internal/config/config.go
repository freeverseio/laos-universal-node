package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"
)

const (
	klaosParachain         uint64 = 3336
	caladanParachain       uint64 = 2900
	klaosGlobalConsensus   string = "3"
	caladanGlobalConsensus string = "0:0x22c48a576c33970622a2b4686a8aa5e4b58350247d69fb5d8015f12a8c8e1e4c"
)

type Config struct {
	WaitingTime      time.Duration
	StartingBlock    uint64
	EvoStartingBlock uint64
	Parachain        uint64
	Contracts        []string
	Rpc              string
	EvoRpc           string
	Path             string
	GlobalConsensus  string
	BlocksMargin     uint
	BlocksRange      uint
	EvoBlocksMargin  uint
	EvoBlocksRange   uint
	Port             uint
	Debug            bool
}

func Load() (*Config, error) {
	defaultStoragePath := getDefaultStoragePath()

	blocksRange := flag.Uint("blocks_range", 100, "Amount of blocks the scanner processes")
	blocksMargin := flag.Uint("blocks_margin", 0, "Number of blocks to assume finality")
	evoBlocksRange := flag.Uint("evo_blocks_range", 100, "Amount of blocks the scanner processes on the evolution chain")
	evoBlocksMargin := flag.Uint("evo_blocks_margin", 0, "Number of blocks to assume finality on the evolution chain")
	contracts := flag.String("contracts", "", "Comma-separated list of the web3 addresses of the smart contracts to scan")
	debug := flag.Bool("debug", false, "Set logs to debug level")
	rpc := flag.String("rpc", "https://eth.llamarpc.com", "URL of the RPC node of an evm-compatible blockchain")
	evoRpc := flag.String("evo_rpc", "", "URL of the RPC evolution chain")
	port := flag.Uint("port", 5001, "HTTP port to use for the universal node server")
	startingBlock := flag.Uint64("starting_block", 18288287, "Initial block where the scanning process should start from")
	evoStartingBlock := flag.Uint64("evo_starting_block", 0, "Initial block where the scanning process should start from on the evolution chain")
	waitingTime := flag.Duration("wait", 5*time.Second, "Waiting time between scans when scanning reaches the last block")
	storagePath := flag.String("storage_path", defaultStoragePath, "Path to the storage folder")

	flag.Parse()

	c := &Config{
		BlocksMargin:     *blocksMargin,
		BlocksRange:      *blocksRange,
		EvoBlocksMargin:  *evoBlocksMargin,
		EvoBlocksRange:   *evoBlocksRange,
		Debug:            *debug,
		Rpc:              *rpc,
		EvoRpc:           *evoRpc,
		StartingBlock:    *startingBlock,
		EvoStartingBlock: *evoStartingBlock,
		WaitingTime:      *waitingTime,
		Port:             *port,
		Path:             *storagePath,
	}

	if *contracts != "" {
		c.Contracts = strings.Split(*contracts, ",")
	}

	if err := setGlobalConsensusAndParachain(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) LogFields() {
	slog.Debug("config loaded", slog.Group("config", "rpc", c.Rpc, "evo_rpc", c.EvoRpc, "contracts", c.Contracts, "starting_block", c.StartingBlock,
		"evo_starting_block", c.EvoStartingBlock, "blocks_margin", c.BlocksMargin, "evo_blocks_margin", c.EvoBlocksMargin, "blocks_range", c.BlocksRange,
		"evo_blocks_range", c.EvoBlocksRange, "debug", c.Debug, "wait", c.WaitingTime, "port", c.Port, "storage_path", c.Path))
}

func getDefaultStoragePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Warn("user home directory not found, default storage path will be under the current directory", "err", err)
		homeDir = "./"
	}
	return path.Join(homeDir, ".universalnode")
}

func setGlobalConsensusAndParachain(c *Config) error {
	switch {
	case strings.Contains(strings.ToLower(c.EvoRpc), "caladan"):
		c.GlobalConsensus = caladanGlobalConsensus
		c.Parachain = caladanParachain
	case strings.Contains(strings.ToLower(c.EvoRpc), "klaos"):
		c.GlobalConsensus = klaosGlobalConsensus
		c.Parachain = klaosParachain
	default:
		return fmt.Errorf("unknown Evolution chain rpc provided: %s", c.EvoRpc)
	}

	return nil
}
