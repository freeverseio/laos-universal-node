package universal

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/freeverseio/laos-universal-node/internal/config"
	contractDiscoverer "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer"
	contractUpdater "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/updater"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type ReorgError struct {
	Block       uint64
	ChainHash   common.Hash
	StorageHash common.Hash
}

func (e ReorgError) Error() string {
	return "reorg error"
}

type Processor interface {
	GetInitStartingBlock(ctx context.Context) (uint64, error)
	GetLastBlock(ctx context.Context, startingBlock uint64) (uint64, error)
	VerifyChainConsistency(ctx context.Context, startingBlock uint64) error
	IsEvoSyncedWithOwnership(ctx context.Context, lastOwnershipBlock uint64) (bool, error)

	ProcessUniversalBlockRange(ctx context.Context, startingBlock, lastBlock uint64) error
}

type processor struct {
	client              blockchain.EthClient
	stateService        state.Service
	scanner             scan.Scanner
	configStartingBlock uint64
	configBlocksMargin  uint64
	configBlocksRange   uint64
	discoverer          contractDiscoverer.Discoverer
	updater             contractUpdater.Updater
}

func NewProcessor(client blockchain.EthClient,
	stateService state.Service,
	scanner scan.Scanner,
	c *config.Config,
	discoverer contractDiscoverer.Discoverer,
	updater contractUpdater.Updater,
) *processor {
	return &processor{
		client:              client,
		stateService:        stateService,
		scanner:             scanner,
		configStartingBlock: c.StartingBlock,
		configBlocksMargin:  uint64(c.BlocksMargin),
		configBlocksRange:   uint64(c.BlocksRange),
		discoverer:          discoverer,
		updater:             updater,
	}
}

func (p *processor) GetInitStartingBlock(ctx context.Context) (uint64, error) {
	tx := p.stateService.NewTransaction()
	defer tx.Discard()
	startingBlockData, err := tx.GetLastOwnershipBlock()
	if err != nil {
		return 0, fmt.Errorf("error retrieving the current block from storage: %w", err)
	}

	if startingBlockData.Number != 0 {
		slog.Debug("ignoring user provided starting block, using last updated block from storage", "startingBlock", startingBlockData.Number)
		return startingBlockData.Number + 1, nil
	}

	if p.configStartingBlock != 0 {
		slog.Debug("using user provided starting block", "startingBlock", p.configStartingBlock)
		return p.configStartingBlock, nil
	}

	startingBlock, err := p.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("error retrieving the latest block from chain: %w", err)
	}

	slog.Debug("using latestBlock from blockchain as startingBlock", "startingBlock", startingBlock)
	return startingBlock, nil
}

func (p *processor) GetLastBlock(ctx context.Context, startingBlock uint64) (uint64, error) {
	l1LatestBlock, err := p.client.BlockNumber(ctx)
	if err != nil {
		slog.Error("error retrieving the latest block", "err", err.Error())
		return 0, err
	}

	return min(startingBlock+p.configBlocksRange, l1LatestBlock-p.configBlocksMargin), nil
}

func (p *processor) VerifyChainConsistency(ctx context.Context, startingBlock uint64) error {
	tx := p.stateService.NewTransaction()
	defer tx.Discard()

	lastBlockDB, err := tx.GetLastOwnershipBlock()
	if err != nil {
		slog.Error("error occurred while reading LaosEvolution end range block hash", "err", err.Error())
		return err
	}

	// During the initial iteration, no hash is stored in the database, so this code block is bypassed.
	if (lastBlockDB.Hash == common.Hash{}) {
		return nil
	}

	// Verify whether the hash of the last block from the previous iteration remains unchanged;
	// if it differs, it indicates a reorganization has taken place.
	previousLastBlock := startingBlock - 1
	slog.Debug("verifying chain consistency on block number", "previousLastBlock", previousLastBlock)
	previousLastBlockData, err := p.client.BlockByNumber(ctx, big.NewInt(int64(previousLastBlock)))
	if err != nil {
		slog.Error("error occurred while retrieving new start range block", "err", err.Error())
		return err
	}

	if previousLastBlockData.Hash().Cmp(lastBlockDB.Hash) != 0 {
		return ReorgError{Block: previousLastBlock, ChainHash: previousLastBlockData.Hash(), StorageHash: lastBlockDB.Hash}
	}

	return nil
}

func (p *processor) IsEvoSyncedWithOwnership(ctx context.Context, lastOwnershipBlock uint64) (bool, error) {
	lastBlockHeader, err := p.client.HeaderByNumber(ctx, big.NewInt(int64(lastOwnershipBlock)))
	if err != nil {
		return false, err
	}

	tx := p.stateService.NewTransaction()
	defer tx.Discard()

	lastEvoBLockData, err := tx.GetLastEvoBlock()
	if err != nil {
		return false, err
	}

	slog.Debug("IsEvoSyncedWithOwnership", "evo_block_number", lastEvoBLockData.Number, "evo_block_timestamp", lastEvoBLockData.Timestamp,
		"lastOwnershipBlock", lastOwnershipBlock, "lastOwnershipBlockTimestamp", lastBlockHeader.Time)

	if lastEvoBLockData.Timestamp < lastBlockHeader.Time {
		return false, nil
	}

	return true, nil
}

func (p *processor) ProcessUniversalBlockRange(ctx context.Context, startingBlock, lastBlock uint64) error {
	tx := p.stateService.NewTransaction()
	defer tx.Discard()

	lastBlockData, err := getLastBlockData(ctx, p.client, lastBlock)
	if err != nil {
		return err
	}

	shouldDiscover, err := p.discoverer.ShouldDiscover(tx, startingBlock, lastBlock)
	if err != nil {
		slog.Error("error occurred reading contracts from storage", "err", err.Error())
		return err
	}
	if shouldDiscover {
		errDiscover := p.discoverer.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
		if errDiscover != nil {
			return errDiscover
		}
	}

	contracts, err := p.discoverer.GetContracts(tx)
	if err != nil {
		return err
	}

	if len(contracts) > 0 {
		transferEvents, errTr := p.updater.GetModelTransferEvents(ctx, startingBlock, lastBlock, contracts)
		if errTr != nil {
			return errTr
		}
		err = p.updater.UpdateState(ctx, tx, contracts, transferEvents, lastBlockData)
		if err != nil {
			return err
		}
	}

	slog.Debug("setting ownership end range block hash for block number",
		"blockNumber", lastBlockData.Number, "blockHash", lastBlockData.Hash, "timestamp", lastBlockData.Timestamp)
	if err = tx.SetLastOwnershipBlock(lastBlockData); err != nil {
		slog.Error("error occurred while storing end range block hash", "err", err.Error())
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("error committing transaction", "err", err.Error())
		return err
	}

	return nil
}

func getLastBlockData(ctx context.Context, client blockchain.EthClient, lastBlock uint64) (model.Block, error) {
	block, err := client.BlockByNumber(ctx, big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred retrieving ownership end range block", "lastBlock", lastBlock, "err", err.Error())
		return model.Block{}, err
	}

	return model.Block{
		Number:    lastBlock,
		Timestamp: block.Header().Time,
		Hash:      block.Hash(),
	}, nil
}
