package universal

import (
	"context"
	"fmt"
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

const (
	safeBlockMargin = 250
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
	RecoverFromReorg(ctx context.Context, startingBlock uint64) (*model.Block, error)
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

// RecoverFromReorg is called when a reorg is detected. It will recursively check for reorgs until it finds a block without reorg.
// It will set the last ownership block to the block without reorg and delete all block hashes and block numbers after the block without reorg.
// It will checkout the merkle tree at the block without reorg, commit the transaction to flush the data to disk,
// and return the block without reorg.
func (p *processor) RecoverFromReorg(ctx context.Context, currentBlock uint64) (*model.Block, error) {
	// Start a transaction
	tx, err := p.stateService.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Discard()

	storedBlockNumbers, err := tx.GetAllStoredBlockNumbers()
	if err != nil {
		return nil, err
	}
	// Check for reorg recursively
	blockWithoutReorg, err := p.findBlockWithoutReorg(ctx, tx, currentBlock, storedBlockNumbers)
	if err != nil {
		return nil, err
	}
	// set last ownership block to block without reorg
	if errLastOwnership := tx.SetLastOwnershipBlock(*blockWithoutReorg); errLastOwnership != nil {
		return nil, errLastOwnership
	}
	// deleting all block hashes after the block without reorg
	if errDeleteOrphanBlockData := tx.DeleteOrphanBlockData(blockWithoutReorg.Number); errDeleteOrphanBlockData != nil {
		return nil, errDeleteOrphanBlockData
	}
	// deleting all root tags after the block without reorg
	if errDeleteOrphanRootTags := tx.DeleteOrphanRootTags(int64(blockWithoutReorg.Number)+1, int64(currentBlock)); errDeleteOrphanRootTags != nil {
		return nil, errDeleteOrphanRootTags
	}

	err = tx.Checkout(int64(blockWithoutReorg.Number))
	if err != nil {
		return nil, err
	}

	if errCommit := tx.Commit(); errCommit != nil {
		return nil, errCommit
	}

	return blockWithoutReorg, nil
}

func (p *processor) findBlockWithoutReorg(ctx context.Context, tx state.Tx, currentBlock uint64, storedBlockNumbers []uint64) (*model.Block, error) {
	blockNumberToCheck, found := getNextLowerBlockNumber(currentBlock, storedBlockNumbers)
	if !found { // no lower block number found
		// we get a safe block number to start from
		return getSafeBlock(currentBlock), nil
	}

	blockToCheck, err := tx.GetOwnershipBlock(blockNumberToCheck)
	if err != nil {
		slog.Error("error retrieving block data", "blockNumber", blockNumberToCheck, "err", err.Error())
		return nil, err
	}

	err = p.checkBlockForReorg(ctx, blockToCheck)
	switch e := err.(type) {
	case nil:
		// no Reorg detected
		return &blockToCheck, e
	case ReorgError:
		// reorg, continue checking the previous blocks
		return p.findBlockWithoutReorg(ctx, tx, blockNumberToCheck, storedBlockNumbers)
	default:
		// Other error occurred
		return nil, err
	}
}

func getNextLowerBlockNumber(currentBlock uint64, storedBlockNumbers []uint64) (blockNumber uint64, exist bool) {
	var maxLowerBlock uint64
	found := false

	// Finding the maximum number lower than the current block in the copy
	for _, blockNumber := range storedBlockNumbers {
		if blockNumber < currentBlock {
			if !found || blockNumber > maxLowerBlock {
				maxLowerBlock = blockNumber
				found = true
			}
		}
	}

	if !found {
		return 0, false
	}
	return maxLowerBlock, found
}

func (p *processor) checkBlockForReorg(ctx context.Context, lastBlockToCheck model.Block) error {
	if (lastBlockToCheck.Hash == common.Hash{}) {
		return fmt.Errorf("no hash stored in the database for block %d", lastBlockToCheck.Number)
	}
	// Verify whether the hash of the last block from the previous iteration remains unchanged;
	// if it differs, it indicates a reorganization has taken place.
	previousStoredBlock := lastBlockToCheck.Number
	slog.Debug("verifying chain consistency on block number", "previousLastBlock", previousStoredBlock)
	previousLastBlockData, err := p.client.HeaderByNumber(ctx, big.NewInt(int64(previousStoredBlock)))
	if err != nil {
		slog.Error("error occurred while retrieving new start range block", "err", err.Error())
		return err
	}
	h := previousLastBlockData.Hash()
	fmt.Println("previousLastBlockData.Hash()", h)
	// If the hash is the same, it means there was no reorganization
	if previousLastBlockData.Hash().Cmp(lastBlockToCheck.Hash) != 0 {
		return ReorgError{Block: previousStoredBlock, ChainHash: previousLastBlockData.Hash(), StorageHash: lastBlockToCheck.Hash}
	}

	return nil
}

func (p *processor) IsEvoSyncedWithOwnership(ctx context.Context, lastOwnershipBlock uint64) (bool, error) {
	lastBlockHeader, err := p.client.HeaderByNumber(ctx, big.NewInt(int64(lastOwnershipBlock)))
	if err != nil {
		return false, err
	}

	tx, err := p.stateService.NewTransaction()
	if err != nil {
		return false, err
	}
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
	tx, err := p.stateService.NewTransaction()
	if err != nil {
		slog.Error("error occurred while creating transaction", "err", err.Error())
		return err
	}
	defer tx.Discard()

	// Get the last stored block number from the database
	previousLastBlockDB, err := tx.GetLastOwnershipBlock()
	if err != nil {
		slog.Error("error occurred while reading LaosEvolution end range block hash", "err", err.Error())
		return err
	}

	lastBlockData, err := getLastBlockData(ctx, p.client, lastBlock)
	if err != nil {
		return err
	}

	shouldDiscover, err := p.discoverer.ShouldDiscover(tx, startingBlock, lastBlock)
	if err != nil {
		slog.Error("error occurred reading contracts from storage", "err", err.Error())
		return err
	}

	newContracts := make(map[common.Address]uint64)
	if shouldDiscover {
		newContracts, err = p.discoverer.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
		if err != nil {
			return err
		}
	}

	contracts, err := p.discoverer.GetContracts(tx)
	if err != nil {
		return err
	}

	transferEvents := make(map[uint64]map[string][]model.ERC721Transfer)
	if len(contracts) > 0 {
		transferEvents, err = p.updater.GetModelTransferEvents(ctx, startingBlock, lastBlock, contracts)
		if err != nil {
			return err
		}
	}

	err = p.updater.UpdateState(ctx, tx, contracts, newContracts, transferEvents, startingBlock, lastBlockData)
	if err != nil {
		return err
	}

	slog.Debug("setting ownership end range block hash for block number",
		"blockNumber", lastBlockData.Number, "blockHash", lastBlockData.Hash, "timestamp", lastBlockData.Timestamp)

	if err = tx.SetLastOwnershipBlock(lastBlockData); err != nil {
		slog.Error("error occurred while storing end range block hash", "err", err.Error())
		return err
	}

	// check for reorgs
	// During the initial iteration, no hash is stored in the database, so this code block is bypassed.
	if (previousLastBlockDB.Hash != common.Hash{}) {
		// we check the previously stored last block and check if it is still on the same branch as the current last block
		// otherwise we return a reorg error
		err = p.checkBlockForReorg(ctx, previousLastBlockDB)
		if err != nil {
			return err
		}
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

func getSafeBlock(currentBlockNumber uint64) *model.Block {
	if currentBlockNumber < safeBlockMargin {
		return &model.Block{Number: 0}
	}
	return &model.Block{Number: currentBlockNumber - safeBlockMargin}
}
