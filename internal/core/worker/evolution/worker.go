package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/core/processor/evolution"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type Worker interface {
	Run(ctx context.Context) error
}
type worker struct {
	waitingTime time.Duration
	processor   evolution.Processor
}

func NewWorker(c *config.Config, client scan.EthClient, scanner scan.Scanner, stateService state.Service) Worker {
	return &worker{
		waitingTime: c.WaitingTime,
		processor: evolution.NewProcessor(
			client,
			stateService,
			scanner,
			c.EvoStartingBlock,
			uint64(c.EvoBlocksMargin),
			uint64(c.EvoBlocksRange)),
	}
}

func (w *worker) Run(ctx context.Context) error {
	slog.Info("starting evolution worker")
	startingBlock, err := w.processor.GetInitStartingBlock(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return nil
		default:
			lastBlock, err := executeEvoBlockRange(ctx, w, startingBlock)
			if err != nil {
				var reorgErr evolution.ReorgError
				if errors.As(err, &reorgErr) {
					slog.Error("evolution chain reorganization detected", "block number", reorgErr.Block, "chain hash", reorgErr.ChainHash.String(), "storage hash", reorgErr.StorageHash.String())
					slog.Info("***********************************************************************************************")
					slog.Info("Please wipe out the database before running the node again.")
					slog.Info("***********************************************************************************************")
					return reorgErr
				}
				break
			}

			startingBlock = lastBlock + 1
		}
	}
}

func executeEvoBlockRange(ctx context.Context, w *worker, startingBlock uint64) (uint64, error) {
	lastBlock, err := w.processor.GetLastBlock(ctx, startingBlock)
	if err != nil {
		return 0, err
	}
	if lastBlock < startingBlock {
		slog.Debug("evolution worker, last calculated block is behind starting block, waiting...",
			"lastBlock", lastBlock, "startingBlock", startingBlock)
		waitBeforeNextScan(ctx, w.waitingTime)
		return startingBlock - 1, nil // return lastBlock from previous range to avoid skipping a block
	}

	err = w.processor.VerifyChainConsistency(ctx, startingBlock)
	if err != nil {
		return 0, err
	}

	err = w.processor.ProcessEvoBlockRange(ctx, startingBlock, lastBlock)
	if err != nil {
		return 0, err
	}
	return lastBlock, nil
}

func waitBeforeNextScan(ctx context.Context, waitingTime time.Duration) {
	timer := time.NewTimer(waitingTime)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}
