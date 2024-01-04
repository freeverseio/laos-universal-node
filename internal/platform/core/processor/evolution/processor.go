package evolution

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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
	ProcessEvoBlockRange(ctx context.Context, startingBlock, lastBlock uint64) error
}

type processor struct {
	client              scan.EthClient
	stateService        state.Service
	scanner             scan.Scanner
	configStartingBlock uint64
	configBlocksRange   uint64
	configBlocksMargin  uint64
}

func NewProcessor(client scan.EthClient,
	stateService state.Service,
	scanner scan.Scanner,
	configStartingBlock,
	configBlocksMargin,
	configBlocksRange uint64,
) *processor {
	return &processor{
		client:              client,
		stateService:        stateService,
		scanner:             scanner,
		configStartingBlock: configStartingBlock,
		configBlocksMargin:  configBlocksMargin,
		configBlocksRange:   configBlocksRange,
	}
}

func (p *processor) GetInitStartingBlock(ctx context.Context) (uint64, error) {
	tx := p.stateService.NewTransaction()
	defer tx.Discard()
	startingBlockData, err := tx.GetLastEvoBlock()
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

	lastBlockDB, err := tx.GetLastEvoBlock()
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

func (p *processor) ProcessEvoBlockRange(ctx context.Context, startingBlock, lastBlock uint64) error {
	tx := p.stateService.NewTransaction()
	defer tx.Discard()

	events, err := p.scanner.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), nil)
	if err != nil {
		slog.Error("error occurred while scanning LaosEvolution events", "err", err.Error())
		return err
	}

	err = storeMintedWithExternalURIEventsByContract(tx, events)
	if err != nil {
		slog.Error("error occurred while storing minted events", "err", err.Error())
		return err
	}

	err = updateLastBlockData(ctx, tx, p.client, lastBlock)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("error committing transaction", "err", err.Error())
		return err
	}

	return nil
}

func updateLastBlockData(ctx context.Context, tx state.Tx, client scan.EthClient, lastBlock uint64) error {
	lastBlockData, err := client.BlockByNumber(ctx, big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred while fetching LaosEvolution end range block", "lastBlock", lastBlock, "err", err.Error())
		return err
	}

	slog.Debug("setting evo end range block data for block number",
		"blockNumber", lastBlockData.Number(), "blockHash", lastBlockData.Hash(), "Timestamps", lastBlockData.Header().Time)

	block := model.Block{
		Number:    lastBlock,
		Timestamp: lastBlockData.Header().Time,
		Hash:      lastBlockData.Hash(),
	}

	return tx.SetLastEvoBlock(block)
}

func storeMintedWithExternalURIEventsByContract(tx state.Tx, events []scan.Event) error {
	groupedMintEvents := groupEventsMintedWithExternalURIByContract(events)

	for contract, scannedEvents := range groupedMintEvents {
		// fetch current storedEvents stored for this specific contract address
		storedEvents, err := tx.GetMintedWithExternalURIEvents(contract.String())
		if err != nil {
			return err
		}

		ev := make([]model.MintedWithExternalURI, 0)
		if storedEvents != nil {
			ev = append(ev, storedEvents...)
		}
		ev = append(ev, scannedEvents...)
		if err := tx.StoreMintedWithExternalURIEvents(contract.String(), ev); err != nil {
			return err
		}
	}

	return nil
}

// groups events that are of type scan.EventMintedWithExternalURI by contract address
func groupEventsMintedWithExternalURIByContract(events []scan.Event) map[common.Address][]model.MintedWithExternalURI {
	groupMintEvents := make(map[common.Address][]model.MintedWithExternalURI, 0)
	for _, event := range events {
		if e, ok := event.(scan.EventMintedWithExternalURI); ok {
			groupMintEvents[e.Contract] = append(groupMintEvents[e.Contract], model.MintedWithExternalURI{
				Slot:        e.Slot,
				To:          e.To,
				TokenURI:    e.TokenURI,
				TokenId:     e.TokenId,
				BlockNumber: e.BlockNumber,
				Timestamp:   e.Timestamp,
			})
		}
	}
	return groupMintEvents
}
