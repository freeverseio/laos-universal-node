package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"log/slog"

	"github.com/freeverseio/laos-universal-node/config"
	"github.com/freeverseio/laos-universal-node/scan/mock"
	"go.uber.org/mock/gomock"
)

func TestRun(t *testing.T) {
	t.Run(`it should log "error retrieving the latest block" message when an error occures while retrieving L1 lastest block`, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		climock := mock.NewMockEthClient(ctrl)
		scanmock := mock.NewMockScanner(ctrl)

		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		c := config.Config{
			StartingBlock: 1,
			BlocksMargin:  0,
			BlocksRange:   100,
		}

		var logOutput bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&logOutput, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}).WithAttrs([]slog.Attr{
			slog.String("version", version),
		}))
		slog.SetDefault(logger)

		climock.EXPECT().BlockNumber(ctx).Return(uint64(0), fmt.Errorf("error retrieving L1 latest block")).AnyTimes()

		go func() {
			run(ctx, c, climock, scanmock)
		}()

		timerCh := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Millisecond)
			close(timerCh)
		}()

		for {
			select {
			case <-timerCh:
				return
			default:
				time.Sleep(time.Millisecond)
				if !strings.Contains(logOutput.String(), "error retrieving the latest block") {
					t.Fatalf("expecting log with message: %s", "error retrieving the latest block")
				}
			}
		}

	})

	t.Run(`it should log "error occurred while scanning events" message when an error occures while scanning events`, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		climock := mock.NewMockEthClient(ctrl)
		scanmock := mock.NewMockScanner(ctrl)

		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		c := config.Config{
			StartingBlock: 1,
			BlocksMargin:  0,
			BlocksRange:   100,
		}

		var logOutput bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&logOutput, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}).WithAttrs([]slog.Attr{
			slog.String("version", version),
		}))
		slog.SetDefault(logger)

		toBlock := int64(c.StartingBlock) + int64(c.BlocksRange)
		climock.EXPECT().BlockNumber(ctx).Return(uint64(toBlock), nil).AnyTimes()
		scanmock.EXPECT().ScanEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(toBlock)).Return(nil, fmt.Errorf("error scanning events occured")).AnyTimes()

		go func() {
			run(ctx, c, climock, scanmock)
		}()

		timerCh := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Millisecond)
			close(timerCh)
		}()

		for {
			select {
			case <-timerCh:
				return
			default:
				time.Sleep(time.Millisecond)
				if !strings.Contains(logOutput.String(), "error occurred while scanning events") {
					t.Fatalf("expecting log with message: %s", "error occurred while scanning events")
				}
			}
		}

	})
}
