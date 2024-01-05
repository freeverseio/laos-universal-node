package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	processMock "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/mock"

	"go.uber.org/mock/gomock"
)

func TestExecuteUniversalBlockRange(t *testing.T) {
	t.Parallel()

	t.Run("error getting last block", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)
		lastBlock := uint64(200)
		evoSynced := true

		ctx := context.TODO()
		processor := processMock.NewMockProcessor(gomock.NewController(t))
		w := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(uint64(0), errors.New("error getting last block"))

		_, _, err := w.executeUniversalBlockRange(ctx, evoSynced, startingBlock, lastBlock)
		assertError(t, errors.New("error getting last block"), err)
	})

	t.Run("last block before starting block", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)
		lastBlock := uint64(200)
		evoSynced := true

		ctx := context.TODO()
		processor := processMock.NewMockProcessor(gomock.NewController(t))
		w := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(uint64(90), nil)

		lastBlock, evoSynced, err := w.executeUniversalBlockRange(ctx, evoSynced, startingBlock, lastBlock)
		assertError(t, nil, err)
		if lastBlock != 99 {
			t.Fatalf(`got last block %d, expected 99`, lastBlock)
		}
		if evoSynced != true {
			t.Fatalf(`got evoSynced %v, expected true`, evoSynced)
		}
	})

	t.Run("error when checking is synced with evo", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)
		lastBlock := uint64(100)
		evoSynced := true

		ctx := context.TODO()
		processor := processMock.NewMockProcessor(gomock.NewController(t))
		w := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(uint64(150), nil)
		processor.EXPECT().IsEvoSyncedWithOwnership(ctx, uint64(150)).Return(false, errors.New("error getting last block"))

		_, _, err := w.executeUniversalBlockRange(ctx, evoSynced, startingBlock, lastBlock)
		assertError(t, errors.New("error getting last block"), err)
	})

	t.Run("chains are not synced", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)
		lastBlock := uint64(100)
		calculatedLastBlock := uint64(150)
		evoSynced := true

		ctx := context.TODO()
		processor := processMock.NewMockProcessor(gomock.NewController(t))
		w := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(calculatedLastBlock, nil)
		processor.EXPECT().IsEvoSyncedWithOwnership(ctx, calculatedLastBlock).Return(false, nil)

		lastBlockObtained, evoSynced, err := w.executeUniversalBlockRange(ctx, evoSynced, startingBlock, lastBlock)
		assertError(t, nil, err)
		if calculatedLastBlock != lastBlockObtained {
			t.Fatalf(`got last block %d, expected 0`, lastBlock)
		}
		if evoSynced != false {
			t.Fatalf(`got evoSynced %v, expected false`, evoSynced)
		}
	})

	t.Run("chains are synced error on verifying chain consistency", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)
		lastBlock := uint64(100)
		calculatedLastBlock := uint64(150)
		evoSynced := true

		ctx := context.TODO()
		processor := processMock.NewMockProcessor(gomock.NewController(t))
		w := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(calculatedLastBlock, nil)
		processor.EXPECT().IsEvoSyncedWithOwnership(ctx, calculatedLastBlock).Return(true, nil)
		processor.EXPECT().VerifyChainConsistency(ctx, startingBlock).Return(errors.New("error on verifying chain consistency"))

		_, _, err := w.executeUniversalBlockRange(ctx, evoSynced, startingBlock, lastBlock)
		assertError(t, errors.New("error on verifying chain consistency"), err)
	})

	t.Run("chains are synced and chain consistency is verified", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)
		lastBlock := uint64(100)
		calculatedLastBlock := uint64(150)
		evoSynced := true

		ctx := context.TODO()
		processor := processMock.NewMockProcessor(gomock.NewController(t))
		w := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(calculatedLastBlock, nil)
		processor.EXPECT().IsEvoSyncedWithOwnership(ctx, calculatedLastBlock).Return(true, nil)
		processor.EXPECT().VerifyChainConsistency(ctx, startingBlock).Return(nil)
		processor.EXPECT().ProcessUniversalBlockRange(ctx, startingBlock, calculatedLastBlock).Return(nil)

		lastBlock, evoSynced, err := w.executeUniversalBlockRange(ctx, evoSynced, startingBlock, lastBlock)
		assertError(t, nil, err)

		if lastBlock != 150 {
			t.Fatalf(`got last block %d, expected 150`, lastBlock)
		}
		if evoSynced != true {
			t.Fatalf(`got evoSynced %v, expected true`, evoSynced)
		}
	})
}

func assertError(t *testing.T, expectedError, err error) {
	t.Helper()
	if expectedError != nil {
		if err.Error() != expectedError.Error() {
			t.Fatalf(`got error "%v", expected error: "%v"`, err, expectedError)
		}
	} else {
		if err != expectedError {
			t.Fatalf(`got error "%v", expected error: "%v"`, err, expectedError)
		}
	}
}
