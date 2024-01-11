package universal

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/freeverseio/laos-universal-node/internal/config"
	shared "github.com/freeverseio/laos-universal-node/internal/core/processor"
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
	client       blockchain.EthClient
	stateService state.Service
	scanner      scan.Scanner
	*shared.BlockHelper
	discoverer contractDiscoverer.Discoverer
	updater    contractUpdater.Updater
}

func NewProcessor(client blockchain.EthClient,
	stateService state.Service,
	scanner scan.Scanner,
	c *config.Config,
	discoverer contractDiscoverer.Discoverer,
	updater contractUpdater.Updater,
) *processor {
	return &processor{
		client:       client,
		stateService: stateService,
		scanner:      scanner,
		BlockHelper: shared.NewBlockHelper(
			client,
			stateService,
			uint64(c.BlocksRange),
			uint64(c.BlocksMargin),
			c.StartingBlock,
		),
		discoverer: discoverer,
		updater:    updater,
	}
}

func (p *processor) GetInitStartingBlock(ctx context.Context) (uint64, error) {
	return p.GetOwnershipInitStartingBlock(ctx)
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
	previousLastBlockData, err := p.client.HeaderByNumber(ctx, big.NewInt(int64(previousLastBlock)))
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
	header, err := client.HeaderByNumber(ctx, big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred retrieving ownership end range block from L1", "lastBlock", lastBlock, "err", err.Error())
		return model.Block{}, err
	}

	return model.Block{
		Number:    lastBlock,
		Timestamp: header.Time,
		Hash:      header.Hash(),
	}, nil
}
