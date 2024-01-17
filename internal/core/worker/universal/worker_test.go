package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/universal"
	mockProcessor "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/mock"
	worker "github.com/freeverseio/laos-universal-node/internal/core/worker/universal"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"go.uber.org/mock/gomock"
)

func TestRun_SuccessfulExecutionWithReorgAndRecovery(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProcessorService := mockProcessor.NewMockProcessor(mockCtrl)

	startingBlocks := []uint64{90, 80}
	verifyReorgErrors := []error{
		universal.ReorgError{Block: 90, ChainHash: common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e"), StorageHash: common.HexToHash("0x123")},
		nil,
	}

	mockProcessorService.EXPECT().GetInitStartingBlock(gomock.Any()).Return(startingBlocks[0], nil)

	for i := 0; i < len(startingBlocks); i++ {
		mockProcessorService.EXPECT().GetLastBlock(ctx, startingBlocks[i]).Return(startingBlocks[i], nil)
		mockProcessorService.EXPECT().IsEvoSyncedWithOwnership(ctx, startingBlocks[i]).Return(true, nil)
		mockProcessorService.EXPECT().ProcessUniversalBlockRange(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(verifyReorgErrors[i]).
			Do(func(ctx context.Context, startingBlock, lastBlock uint64) {
				if startingBlock == startingBlocks[len(startingBlocks)-1] {
					cancel()
				}
			})
	}
	mockProcessorService.EXPECT().RecoverFromReorg(ctx, startingBlocks[0]).Return(&model.Block{
		Number: startingBlocks[1],
		Hash:   common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e"),
	}, nil).Times(1)
	w := worker.New(&config.Config{WaitingTime: 1 * time.Second}, mockProcessorService)

	err := w.Run(ctx)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
