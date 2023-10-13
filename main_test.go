package main

import (
	"bytes"
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/freeverseio/laos-universal-node/config"
	"github.com/freeverseio/laos-universal-node/scan/mock"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/slog"
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

	climock.EXPECT().BlockNumber(ctx).Return(c.StartingBlock, nil)
	toBlock := int64(c.StartingBlock) + int64(c.BlocksRange)
	scanmock.EXPECT().ScanEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(toBlock)).Return(nil, nil)

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
			if logOutput.String() == "<nil>" {
				time.Sleep(time.Millisecond * 500)
			}
		}
	}

	// compare log output
}
