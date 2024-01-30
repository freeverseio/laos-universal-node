package shared

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

func (b *BlockHelper) GetLastBlock(ctx context.Context, startingBlock uint64) (uint64, error) {
	l1LatestBlock, err := b.client.BlockNumber(ctx)
	if err != nil {
		slog.Error("error retrieving the latest block", "err", err.Error())
		return 0, err
	}

	return min(startingBlock+b.blocksRange, l1LatestBlock-b.blocksMargin), nil
}

func (b *BlockHelper) GetOwnershipInitStartingBlock(ctx context.Context) (uint64, error) {
	tx, err := b.stateService.NewTransaction()
	if err != nil {
		return 0, fmt.Errorf("error creating a new transaction: %w", err)
	}
	defer tx.Discard()
	startingBlockData, err := tx.GetLastOwnershipBlock()
	if err != nil {
		return 0, fmt.Errorf("error retrieving the current block from storage: %w", err)
	}

	return b.getInitStartingBlock(ctx, startingBlockData)
}

func (b *BlockHelper) GetEvoInitStartingBlock(ctx context.Context) (uint64, error) {
	tx, err := b.stateService.NewTransaction()
	if err != nil {
		return 0, fmt.Errorf("error creating a new transaction: %w", err)
	}
	defer tx.Discard()
	startingBlockData, err := tx.GetLastEvoBlock()
	if err != nil {
		return 0, fmt.Errorf("error retrieving the current block from storage: %w", err)
	}

	return b.getInitStartingBlock(ctx, startingBlockData)
}

func (b *BlockHelper) getInitStartingBlock(ctx context.Context, startingBlockData model.Block) (uint64, error) {
	if startingBlockData.Number != 0 {
		slog.Debug("ignoring user provided starting block, using last updated block from storage", "startingBlock", startingBlockData.Number)
		return startingBlockData.Number + 1, nil
	}

	if b.startingBlock != 0 {
		slog.Debug("using user provided starting block", "startingBlock", b.startingBlock)
		return b.startingBlock, nil
	}

	l1LatestBlock, err := b.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("error retrieving the latest block from chain: %w", err)
	}

	slog.Debug("using latestBlock from blockchain as startingBlock", "startingBlock", l1LatestBlock)
	return l1LatestBlock, nil
}

func WaitBeforeNextRequest(ctx context.Context, waitingTime time.Duration) {
	timer := time.NewTimer(waitingTime)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}
