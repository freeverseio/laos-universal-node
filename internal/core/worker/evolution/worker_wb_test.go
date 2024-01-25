package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	evoProcessMock "github.com/freeverseio/laos-universal-node/internal/core/processor/evolution/mock"

	"go.uber.org/mock/gomock"
)

func TestExecuteBlockRange(t *testing.T) {
	t.Parallel()

	t.Run("error getting last block", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)

		ctx := context.TODO()
		processor := evoProcessMock.NewMockProcessor(gomock.NewController(t))
		worker := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(uint64(0), errors.New("error getting last block"))

		_, err := executeEvoBlockRange(ctx, worker, startingBlock)
		assertError(t, errors.New("error getting last block"), err)
	})

	t.Run("last block less then starting block", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)

		ctx := context.TODO()
		processor := evoProcessMock.NewMockProcessor(gomock.NewController(t))
		worker := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(uint64(90), nil)

		lastBlock, err := executeEvoBlockRange(ctx, worker, startingBlock)
		assertError(t, nil, err)
		if lastBlock != 99 {
			t.Fatalf(`got last block %d, expected 99`, lastBlock)
		}
	})

	t.Run("process evo block range returns error", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)

		ctx := context.TODO()
		processor := evoProcessMock.NewMockProcessor(gomock.NewController(t))
		worker := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(uint64(110), nil)
		processor.EXPECT().
			ProcessEvoBlockRange(ctx, startingBlock, uint64(110)).
			Return(errors.New("process evo block range returns error"))
		_, err := executeEvoBlockRange(ctx, worker, startingBlock)
		assertError(t, errors.New("process evo block range returns error"), err)
	})

	t.Run("process evo block range returns finish successfully", func(t *testing.T) {
		t.Parallel()

		waitingTimeMillisecond := 100 * time.Millisecond
		startingBlock := uint64(100)

		ctx := context.TODO()
		processor := evoProcessMock.NewMockProcessor(gomock.NewController(t))
		worker := &worker{waitingTime: waitingTimeMillisecond, processor: processor}

		processor.EXPECT().GetLastBlock(ctx, startingBlock).Return(uint64(110), nil)

		processor.EXPECT().ProcessEvoBlockRange(ctx, startingBlock, uint64(110)).Return(nil)

		lastBlock, err := executeEvoBlockRange(ctx, worker, startingBlock)
		assertError(t, nil, err)
		if lastBlock != uint64(110) {
			t.Fatalf(`got last block %d, expected %d`, lastBlock, uint64(110))
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
