package blockmapper

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	sharedProcessor "github.com/freeverseio/laos-universal-node/internal/core/processor"
	shared "github.com/freeverseio/laos-universal-node/internal/core/worker"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type BlockCorrectionFactor int

const (
	OwnershipBlockFactor BlockCorrectionFactor = 0
	EvoBlockFactor       BlockCorrectionFactor = -1
)

func (b BlockCorrectionFactor) uint64() uint64 {
	return uint64(b)
}

type Worker interface {
	Run(ctx context.Context)
	SearchBlockByTimestamp(targetTimestamp uint64, client blockchain.EthClient, correctionFactor BlockCorrectionFactor, startingPoint ...uint64) (uint64, error)
}

type worker struct {
	waitingTime          time.Duration
	ownershipBlockHelper *sharedProcessor.BlockHelper
	evoBlockHelper       *sharedProcessor.BlockHelper
	clientOwnership      blockchain.EthClient
	clientEvo            blockchain.EthClient
	stateService         state.Service
}

func New(c *config.Config, ownershipClient, evoClient blockchain.EthClient, stateService state.Service) Worker {
	return &worker{
		waitingTime:     c.WaitingTime,
		clientOwnership: ownershipClient,
		clientEvo:       evoClient,
		ownershipBlockHelper: sharedProcessor.NewBlockHelper(
			ownershipClient,
			stateService,
			uint64(c.BlocksRange),
			uint64(c.BlocksMargin),
			c.StartingBlock,
		),
		evoBlockHelper: sharedProcessor.NewBlockHelper(
			evoClient,
			stateService,
			uint64(c.EvoBlocksRange),
			uint64(c.EvoBlocksMargin),
			c.EvoStartingBlock,
		),
		stateService: stateService,
	}
}

func (w *worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return
		default:
			if err := w.ExecuteMapping(ctx); err != nil {
				slog.Error("error occurred while performing block mapping", "err", err)
			}
		}
	}
}

func (w *worker) ExecuteMapping(ctx context.Context) error {
	tx, err := w.stateService.NewTransaction()
	if err != nil {
		err = fmt.Errorf("error occurred creating transaction: %w", err)
		return err
	}
	defer tx.Discard()

	// check last mapped ownership block from storage
	lastMappedOwnershipBlock, err := tx.GetLastMappedOwnershipBlockNumber()
	if err != nil {
		return fmt.Errorf("error occurred retrieving the latest mapped ownership block from storage: %w", err)
	}

	var evoBlock uint64
	// if no block has ever been mapped, start mapping from the oldest user-defined block
	if lastMappedOwnershipBlock == 0 {
		evoBlock, err = w.getInitialEvoBlock(ctx)
		if err != nil {
			return err
		}
	} else {
		// if a block has been mapped, resume mapping from the next one
		evoBlock, err = tx.GetMappedEvoBlockNumber(lastMappedOwnershipBlock)
		if err != nil {
			return fmt.Errorf("error occurred retrieving the mapped evolution block number by ownership block %d from storage: %w",
				lastMappedOwnershipBlock, err)
		}
		evoBlock++
	}

	// compare the last mapped ownership block with the last processed ownership block
	lastProcessedOwnershipBlock, err := tx.GetLastOwnershipBlock()
	if err != nil {
		return fmt.Errorf("error occurred retrieving the last processed ownership block from storage: %w", err)
	}
	if lastMappedOwnershipBlock >= lastProcessedOwnershipBlock.Number {
		slog.Debug("mapped block has reached the last processed ownership block, waiting to process more blocks before mapping again...")
		shared.Wait(ctx, w.waitingTime)
		return nil
	}

	// given the evo block timestamp, find the corresponding ownership block number
	evoHeader, err := w.clientEvo.HeaderByNumber(ctx, big.NewInt(int64(evoBlock)))
	if err != nil {
		return fmt.Errorf("error occurred retrieving block number %d from evolution chain %w:", evoBlock, err)
	}
	toMapOwnershipBlock, err := w.SearchBlockByTimestamp(evoHeader.Time, w.clientOwnership, OwnershipBlockFactor, lastMappedOwnershipBlock)
	if err != nil {
		return fmt.Errorf("error occurred searching for ownership block number by target timestamp %d (evolution block number %d): %w",
			evoHeader.Time, evoBlock, err)
	}

	// set ownership block -> evo block mapping
	err = tx.SetOwnershipEvoBlockMapping(toMapOwnershipBlock, evoBlock)
	if err != nil {
		return fmt.Errorf("error setting ownership block number %d (key) to evo block number %d (value) in storage: %w",
			toMapOwnershipBlock, evoBlock, err)
	}
	err = tx.SetLastMappedOwnershipBlockNumber(toMapOwnershipBlock)
	if err != nil {
		return fmt.Errorf("error setting the last mapped ownership block number %d in storage: %w", toMapOwnershipBlock, err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

func (w *worker) getInitialEvoBlock(ctx context.Context) (uint64, error) {
	var ownershipHeader, evoHeader *types.Header
	var ownershipStartingBlock, evoStartingBlock uint64
	ownershipStartingBlock, err := w.ownershipBlockHelper.GetOwnershipInitStartingBlock(ctx)
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving the ownership init starting block: %w", err)
	}
	evoStartingBlock, err = w.evoBlockHelper.GetEvoInitStartingBlock(ctx)
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving the evolution init starting block: %w", err)
	}
	ownershipHeader, err = w.clientOwnership.HeaderByNumber(ctx, big.NewInt(int64(ownershipStartingBlock)))
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving block number %d from ownership chain: %w",
			ownershipStartingBlock, err)
	}
	evoHeader, err = w.clientEvo.HeaderByNumber(ctx, big.NewInt(int64(evoStartingBlock)))
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving block number %d from evolution chain: %w",
			evoStartingBlock, err)
	}
	evoBlock := evoStartingBlock
	// if the user-defined ownership block was produced before the user-defined evolution block,
	// look for the evolution block corresponding to that ownership block in time
	if ownershipHeader.Time < evoHeader.Time {
		evoBlock, err = w.SearchBlockByTimestamp(ownershipHeader.Time, w.clientEvo, EvoBlockFactor)
		if err != nil {
			return 0, fmt.Errorf("error occurred searching for evolution block number by target timestamp %d (ownership block number %d): %w",
				ownershipHeader.Time, ownershipStartingBlock, err)
		}
	}
	return evoBlock, nil
}

// SearchBlockByTimestamp performs a binary search to find the block number for a given timestamp.
// It assumes block timestamps are strictly increasing.
func (w *worker) SearchBlockByTimestamp(targetTimestamp uint64, client blockchain.EthClient, correctionFactor BlockCorrectionFactor, startingPoint ...uint64) (uint64, error) {
	var left, right uint64
	if len(startingPoint) > 0 {
		left = startingPoint[0]
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	right = header.Number.Uint64()

	for left <= right {
		mid := left + (right-left)/2
		midHeader, err := client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(mid))
		if err != nil {
			return 0, err
		}
		midTimestamp := midHeader.Time
		switch {
		case midTimestamp < targetTimestamp:
			left = mid + 1
		case midTimestamp > targetTimestamp:
			right = mid - 1
		default:
			return mid + correctionFactor.uint64(), nil
		}
	}

	return left + correctionFactor.uint64(), nil
}
