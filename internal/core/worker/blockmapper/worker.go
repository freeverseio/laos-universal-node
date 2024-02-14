package blockmapper

import (
	"context"
	"log/slog"
	"time"

	"github.com/freeverseio/laos-universal-node/internal/core/processor/blockmapper"
	shared "github.com/freeverseio/laos-universal-node/internal/core/worker"
)

type Worker interface {
	Run(ctx context.Context)
}

type worker struct {
	processor   blockmapper.Processor
	waitingTime time.Duration
}

func New(waitingTime time.Duration, processor blockmapper.Processor) Worker {
	return &worker{
		processor:   processor,
		waitingTime: waitingTime,
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
	synced, err := w.processor.IsMappingSyncedWithProcessing()
	if err != nil {
		return err
	}
	if synced {
		slog.Debug("mapped block has reached the last processed ownership block, waiting to process more blocks before mapping again...")
		shared.Wait(ctx, w.waitingTime)
		return nil
	}
	err = w.processor.MapNextBlock(ctx)
	if err != nil {
		return err
	}
	return nil
}
