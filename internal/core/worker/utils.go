package utils

import (
	"context"
	"time"
)

func WaitBeforeNextScan(ctx context.Context, waitingTime time.Duration) {
	timer := time.NewTimer(waitingTime)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}
