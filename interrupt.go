package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log/slog"
)

// Interrupt is used to gracefully shutdown the laos universal node
func Interrupt(ctx context.Context, cancel context.CancelFunc) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return nil
	case <-c:
		slog.Info("exiting gracefully ...")
		cancel()
		return nil
	}
}
