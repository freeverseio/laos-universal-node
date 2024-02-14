package shared

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type BlockHelper struct {
	client        blockchain.EthClient
	stateService  state.Service
	blocksRange   uint64
	blocksMargin  uint64
	startingBlock uint64
}

func NewBlockHelper(client blockchain.EthClient, stateService state.Service, blocksRange, blocksMargin, startingBlock uint64) *BlockHelper {
	return &BlockHelper{
		client:        client,
		stateService:  stateService,
		blocksRange:   blocksRange,
		blocksMargin:  blocksMargin,
		startingBlock: startingBlock,
	}
}

func (h *BlockHelper) GetLastBlock(ctx context.Context, startingBlock uint64) (uint64, error) {
	l1LatestBlock, err := h.client.BlockNumber(ctx)
	if err != nil {
		slog.Error("error retrieving the latest block", "err", err.Error())
		return 0, err
	}

	return min(startingBlock+h.blocksRange, l1LatestBlock-h.blocksMargin), nil
}

func (h *BlockHelper) GetOwnershipInitStartingBlock(ctx context.Context) (uint64, error) {
	tx, err := h.stateService.NewTransaction()
	if err != nil {
		return 0, fmt.Errorf("error creating a new transaction: %w", err)
	}
	defer tx.Discard()
	startingBlockData, err := tx.GetLastOwnershipBlock()
	if err != nil {
		return 0, fmt.Errorf("error retrieving the current block from storage: %w", err)
	}

	return h.getInitStartingBlock(ctx, startingBlockData)
}

func (h *BlockHelper) GetEvoInitStartingBlock(ctx context.Context) (uint64, error) {
	tx, err := h.stateService.NewTransaction()
	if err != nil {
		return 0, fmt.Errorf("error creating a new transaction: %w", err)
	}
	defer tx.Discard()
	startingBlockData, err := tx.GetLastEvoBlock()
	if err != nil {
		return 0, fmt.Errorf("error retrieving the current block from storage: %w", err)
	}

	return h.getInitStartingBlock(ctx, startingBlockData)
}

func (h *BlockHelper) getInitStartingBlock(ctx context.Context, startingBlockData model.Block) (uint64, error) {
	if startingBlockData.Number != 0 {
		slog.Debug("ignoring user provided starting block, using last updated block from storage", "startingBlock", startingBlockData.Number)
		return startingBlockData.Number + 1, nil
	}

	if h.startingBlock != 0 {
		slog.Debug("using user provided starting block", "startingBlock", h.startingBlock)
		return h.startingBlock, nil
	}

	l1LatestBlock, err := h.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("error retrieving the latest block from chain: %w", err)
	}

	slog.Debug("using latestBlock from blockchain as startingBlock", "startingBlock", l1LatestBlock)
	return l1LatestBlock, nil
}
