package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/universal"
	mockDiscoverer "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer/mock"
	mockProcessor "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/mock"
	mockUpdater "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/updater/mock"
	worker "github.com/freeverseio/laos-universal-node/internal/core/worker/universal"
	mockClient "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	mockScan "github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

func TestRun(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		initBlockError     error
		numberOfExecutions int
		startingBlocks     []uint64
		verifyReorgErrors  []error
		recoverFromReorg   int
	}{
		{
			name:               "successful execution with reorg and recovery",
			initBlockError:     nil,
			numberOfExecutions: 2,
			startingBlocks:     []uint64{90, 80},
			verifyReorgErrors: []error{
				universal.ReorgError{Block: 90, ChainHash: common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e"), StorageHash: common.HexToHash("0x123")},
				nil,
			},
			recoverFromReorg: 1,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockClientService := mockClient.NewMockEthClient(mockCtrl)
			mockScanner := mockScan.NewMockScanner(mockCtrl)
			mockStateService := mockTx.NewMockService(mockCtrl)
			mockDiscovererService := mockDiscoverer.NewMockDiscoverer(mockCtrl)
			mockUpdaterService := mockUpdater.NewMockUpdater(mockCtrl)
			mockProcessorService := mockProcessor.NewMockProcessor(mockCtrl)

			mockProcessorService.EXPECT().GetInitStartingBlock(gomock.Any()).Return(tc.startingBlocks[0], tc.initBlockError)
			mockProcessorService.EXPECT().ProcessUniversalBlockRange(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, startingBlock, lastBlock uint64) error {
				cancel()
				return nil
			})

			for i := 0; i < tc.numberOfExecutions; i++ {
				mockProcessorService.EXPECT().GetLastBlock(ctx, tc.startingBlocks[i]).Return(tc.startingBlocks[i], nil)
				mockProcessorService.EXPECT().IsEvoSyncedWithOwnership(ctx, tc.startingBlocks[i]).Return(true, nil)
				mockProcessorService.EXPECT().VerifyChainConsistency(ctx, tc.startingBlocks[i]).Return(tc.verifyReorgErrors[i]).Do(func(ctx context.Context, block uint64) {
					if block == tc.startingBlocks[len(tc.startingBlocks)-1] {
						cancel()
					}
				})
			}
			mockProcessorService.EXPECT().RecoverFromReorg(ctx, tc.startingBlocks[0]).Return(&model.Block{
				Number: tc.startingBlocks[1],
				Hash:   common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e"),
			}, nil).Times(tc.recoverFromReorg)

			w := worker.New(&config.Config{WaitingTime: 1 * time.Second}, mockClientService, mockScanner, mockStateService, mockDiscovererService, mockUpdaterService,
				worker.WithProcessor(mockProcessorService))

			err := w.Run(ctx)

			if tc.initBlockError != nil {
				if err == nil || err.Error() != tc.initBlockError.Error() {
					t.Errorf("expected error: %v, got: %v", tc.initBlockError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			}
		})
	}
}
