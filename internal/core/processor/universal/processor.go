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
	RecoverFromReorg(ctx context.Context, startingBlock uint64) error
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
	return p.checkBlockForReorg(ctx, lastBlockDB)
}

func (p *processor) RecoverFromReorg(ctx context.Context, currentBlock uint64) error {
	// Start a transaction
	tx := p.stateService.NewTransaction()
	defer tx.Discard()

	storedBlockNumbers, err := tx.GetAllStoredBlockNumbers()
	if err != nil {
		return err
	}
	var saveBlock *model.Block
	nextBlockNumberToCheck, err := getNextLowerBlockNumber(currentBlock, storedBlockNumbers)
	if err != nil { // no lower block number found
		// we get a safe block number to start from
		saveBlock = getSafeBlock(currentBlock)
	} else {
		// Check for reorg recursively
		// TODO rename
		saveBlock, err = p.checkForReorgRecursive(ctx, tx, nextBlockNumberToCheck, storedBlockNumbers)
		if err != nil {
			return err
		}
	}

	contracts := tx.GetAllERC721UniversalContracts()
	// Process each contract
	for _, contract := range contracts {
		// Handle each contract in its own transaction
		contractTx := p.stateService.NewTransaction()
		err = checkout(contractTx, common.HexToAddress(contract), saveBlock.Number)
		if err != nil {
			contractTx.Discard()
			return err
		}
		if err := contractTx.Commit(); err != nil {
			return err // Handle commit error
		}
	}
	return nil
}

func (p *processor) checkForReorgRecursive(ctx context.Context, tx state.Tx, blockNumberToCheck uint64, storedBlockNumbers []uint64) (*model.Block, error) {
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
		nextBlockToCheck, errNextBlockNumber := getNextLowerBlockNumber(blockNumberToCheck, storedBlockNumbers)
		if errNextBlockNumber != nil { // no lower block number found
			// we return a safe block number to start from
			return getSafeBlock(blockNumberToCheck), nil
		}
		return p.checkForReorgRecursive(ctx, tx, nextBlockToCheck, storedBlockNumbers)
	default:
		// Other error occurred
		return nil, err
	}
}

func getNextLowerBlockNumber(currentBlock uint64, storedBlockNumbers []uint64) (uint64, error) {
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
		return 0, fmt.Errorf("no lower block number found")
	}
	return maxLowerBlock, nil
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

func checkout(tx state.Tx, contractAddress common.Address, blockNumber uint64) error {
	err := tx.LoadMerkleTrees(contractAddress)
	if err != nil {
		return err
	}

	err = tx.Checkout(contractAddress, int64(blockNumber))
	if err != nil {
		slog.Error("error occurred checking out merkle tree at block number", "block_number", blockNumber,
			"contract_address", contractAddress, "err", err)
		return err
	}

	return nil
}

func getSafeBlock(currentBlockNumber uint64) *model.Block {
	var safeBlockNumber uint64
	if currentBlockNumber < 250 {
		safeBlockNumber = 0
	} else {
		safeBlockNumber = currentBlockNumber - 250
	}

	return &model.Block{
		Number: safeBlockNumber,
	}
}
