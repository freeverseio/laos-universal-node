package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"log/slog"

	"github.com/freeverseio/laos-universal-node/config"
	"github.com/freeverseio/laos-universal-node/scan/mock"
	"go.uber.org/mock/gomock"
)

func TestRun(t *testing.T) {
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
	scanmock.EXPECT().ScanEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(toBlock)).Return(nil, fmt.Errorf("error in scan events")).AnyTimes()

	go func() {
		run(ctx, c, climock, scanmock)
	}()

	timerCh := make(chan struct{})

	go func() {
		time.Sleep(2 * time.Second)
		// timmerCh <- struct{}{}
		close(timerCh)
	}()

	for {
		select {
		case <-timerCh:
			return
		default:
			out := logOutput.String()
			if out == "<nil>" {
				time.Sleep(time.Millisecond * 500)
			} else {
				if out != "" {
					t.Fatalf("got %s, expected empty logs", out)
				}
				return
			}
		}
	}

	// compare log output
}
