package evolution

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/freeverseio/laos-universal-node/internal/config"
	shared "github.com/freeverseio/laos-universal-node/internal/core/processor"
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
	ProcessEvoBlockRange(ctx context.Context, startingBlock, lastBlock uint64) error
}

type processor struct {
	client       blockchain.EthClient
	stateService state.Service
	scanner      scan.Scanner
	laosHTTP     LaosRPCRequests
	waitingTime  time.Duration
	shared.BlockHelper
}

func NewProcessor(client blockchain.EthClient,
	stateService state.Service,
	scanner scan.Scanner,
	laosHTTP LaosRPCRequests,
	c *config.Config,
) *processor {
	return &processor{
		client:       client,
		stateService: stateService,
		scanner:      scanner,
		laosHTTP:     laosHTTP,
		waitingTime:  c.WaitingRPCRequestTime,
		BlockHelper: shared.NewBlockHelper(
			client,
			stateService,
			uint64(c.EvoBlocksRange),
			uint64(c.EvoBlocksMargin),
			c.EvoStartingBlock,
		),
	}
}

func (p *processor) GetInitStartingBlock(ctx context.Context) (uint64, error) {
	return p.GetEvoInitStartingBlock(ctx)
}

func (p *processor) VerifyChainConsistency(ctx context.Context, startingBlock uint64) error {
	tx, err := p.stateService.NewTransaction()
	if err != nil {
		slog.Debug("error occurred while creating new transaction", "err", err.Error())
		return err
	}
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
	tx, err := p.stateService.NewTransaction()
	if err != nil {
		slog.Debug("error occurred while creating new transaction", "err", err.Error())
		return err
	}
	defer tx.Discard()

	for {
		var ok bool
		ok, errFinalize := p.hasBlockFinalize(big.NewInt(int64(lastBlock)))
		if errFinalize != nil {
			slog.Error("error occurred while checking latest finalized block", "err", errFinalize.Error())
			return errFinalize
		}
		if ok {
			break
		}
		slog.Debug("block not finalized, waiting for finality", "block", lastBlock)

		p.waitBeforeNextRequest(ctx)
	}

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

	err = p.updateFirstBlockData(ctx, tx, startingBlock)
	if err != nil {
		return err
	}

	err = p.updateLastBlockData(ctx, tx, lastBlock)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("error committing transaction", "err", err.Error())
		return err
	}

	return nil
}

func (p *processor) updateFirstBlockData(ctx context.Context, tx state.Tx, firstBlock uint64) error {
	defaultBlock := model.Block{}
	firstBlockStorage, err := tx.GetFirstEvoBlock()
	if err != nil {
		slog.Error("error occurred retrieving first evo block from storage", "err", err)
		return err
	}
	if firstBlockStorage == defaultBlock {
		firstBlockData, err := p.client.BlockByNumber(ctx, big.NewInt(int64(firstBlock)))
		if err != nil {
			slog.Error("error occurred while fetching LaosEvolution block", "firstBlock", firstBlock, "err", err.Error())
			return err
		}
		slog.Debug("setting evo first block data for block number",
			"blockNumber", firstBlockData.Number(), "blockHash", firstBlockData.Hash(), "timestamp", firstBlockData.Header().Time)
		block := model.Block{
			Number:    firstBlock,
			Timestamp: firstBlockData.Header().Time,
			Hash:      firstBlockData.Hash(),
		}
		err = tx.SetFirstEvoBlock(block)
		if err != nil {
			slog.Error("error occurred setting first evo block to storage", "firstBlock", firstBlock, "err", err.Error())
			return err
		}
	}
	return nil
}

func (p *processor) updateLastBlockData(ctx context.Context, tx state.Tx, lastBlock uint64) error {
	lastBlockData, err := p.client.BlockByNumber(ctx, big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred while fetching LaosEvolution end range block", "lastBlock", lastBlock, "err", err.Error())
		return err
	}

	slog.Debug("setting evo end range block data for block number",
		"blockNumber", lastBlockData.Number(), "blockHash", lastBlockData.Hash(), "timestamp", lastBlockData.Header().Time)

	block := model.Block{
		Number:    lastBlock,
		Timestamp: lastBlockData.Header().Time,
		Hash:      lastBlockData.Hash(),
	}

	err = tx.SetLastEvoBlock(block)
	if err != nil {
		slog.Error("error occurred while setting lastEvoBlock to storage",
			"lastBlock", lastBlock, "err", err.Error())
		return err
	}

	return nil
}

func storeMintedWithExternalURIEventsByContract(tx state.Tx, events []scan.Event) error {
	for _, event := range events {
		e, ok := event.(scan.EventMintedWithExternalURI)
		if ok {
			externalMintEvent := &model.MintedWithExternalURI{
				Slot:        e.Slot,
				To:          e.To,
				TokenURI:    e.TokenURI,
				TokenId:     e.TokenId,
				BlockNumber: e.BlockNumber,
				Timestamp:   e.Timestamp,
				TxIndex:     e.TxIndex,
			}

			if err := tx.StoreMintedWithExternalURIEvent(e.Contract.String(), externalMintEvent); err != nil {
				return err
			}

			if err := tx.SetNextEvoEventBlock(e.Contract.String(), e.BlockNumber); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *processor) hasBlockFinalize(blockNumber *big.Int) (bool, error) {
	blockHash, err := p.laosHTTP.LatestFinalizedBlockHash()
	if err != nil {
		return false, err
	}

	finalizedBlockNumber, err := p.laosHTTP.BlockNumber(blockHash)
	if err != nil {
		return false, err
	}

	// if blockNumber > finalizedBlockNumber
	if blockNumber.Cmp(finalizedBlockNumber) == 1 {
		return false, nil
	}

	return true, nil
}

func (p *processor) waitBeforeNextRequest(ctx context.Context) {
	timer := time.NewTimer(p.waitingTime)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}
