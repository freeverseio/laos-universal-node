package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/universal"
	contractDiscoverer "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer"
	contractUpdater "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/updater"
	utils "github.com/freeverseio/laos-universal-node/internal/core/worker"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type Worker interface {
	Run(ctx context.Context) error
}

type worker struct {
	waitingTime time.Duration
	processor   universal.Processor
}

func New(c *config.Config,
	client scan.EthClient,
	scanner scan.Scanner,
	stateService state.Service,
	discoverer contractDiscoverer.Discoverer,
	updater contractUpdater.Updater,
) Worker {
	return &worker{
		waitingTime: c.WaitingTime,
		processor: universal.NewProcessor(
			client,
			stateService,
			scanner,
			c,
			discoverer,
			updater),
	}
}

func (w *worker) Run(ctx context.Context) error {
	slog.Info("starting universal worker")
	startingBlock, err := w.processor.GetInitStartingBlock(ctx)
	if err != nil {
		return err
	}

	evoSynced := true
	lastBlock := startingBlock
	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return nil
		default:
			slog.Debug("executing block range", "startingBlock", startingBlock, "lastBlock", lastBlock, "evoSynced", evoSynced)
			prevLastBlock, wasEvoSynced, err := w.executeUniversalBlockRange(ctx, evoSynced, startingBlock, lastBlock)
			if err != nil {
				slog.Error("error occurred while processing universal block range", "err", err.Error())
				var reorgErr universal.ReorgError
				if errors.As(err, &reorgErr) {
					slog.Error("ownership chain reorganization detected",
						"blockNumber", reorgErr.Block,
						"chainHash", reorgErr.ChainHash.String(),
						"storageHash", reorgErr.StorageHash.String())
					slog.Info("***************************************************************************************")
					slog.Info("Please wipe out the database before running the node again.")
					slog.Info("***************************************************************************************")
					return reorgErr
				}
				break
			}

			evoSynced = wasEvoSynced
			// if evo is not synced (evoSynced == false) the lastBlock should be the one that is previously and don't update startingBlock
			// otherwise update startingBlock to the lastBlock + 1 and the new lastBlock will be calculated in the function
			if !evoSynced {
				lastBlock = prevLastBlock
				break
			}
			startingBlock = prevLastBlock + 1
		}
	}
}

func (w *worker) executeUniversalBlockRange(ctx context.Context,
	evoSynced bool,
	startingBlock,
	lastBlock uint64,
) (previousLastBlock uint64, wasEvoSynced bool, err error) {
	if evoSynced {
		lastBlock, err = w.processor.GetLastBlock(ctx, startingBlock)
		if err != nil {
			return 0, false, err
		}

		if lastBlock < startingBlock {
			slog.Debug("last calculated block is behind starting block, waiting...",
				"lastBlock", lastBlock, "startingBlock", startingBlock)
			utils.WaitBeforeNextScan(ctx, w.waitingTime)
			return startingBlock - 1, true, nil // return lastBlock from previous range to avoid skipping a block
		}
	}

	evoSynced, err = w.processor.IsEvoSyncedWithOwnership(ctx, lastBlock)
	if err != nil {
		slog.Error("error occurred while checking if evolution chain is synced with ownership chain", "err", err.Error())
		return 0, false, err
	}

	if !evoSynced {
		slog.Debug("evolution chain is not synced with ownership chain, waiting...")
		utils.WaitBeforeNextScan(ctx, w.waitingTime)
		return lastBlock, false, nil
	}
	err = w.processor.VerifyChainConsistency(ctx, startingBlock)
	if err != nil {
		return 0, false, err
	}

	err = w.processor.ProcessUniversalBlockRange(ctx, startingBlock, lastBlock)
	if err != nil {
		return 0, false, err
	}
	return lastBlock, true, nil
}
