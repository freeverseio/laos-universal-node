package blockmapper

import (
	"context"
	"fmt"
	"testing"
	"time"

	blockMapperMock "github.com/freeverseio/laos-universal-node/internal/core/processor/blockmapper/mock"
	"go.uber.org/mock/gomock"
)

func TestExecuteMapping(t *testing.T) {
	t.Parallel()
	t.Run("IsMappingSyncedWithProcessing returns false", func(t *testing.T) {
		t.Parallel()
		ctrl, blockMapper := getMocks(t)
		defer ctrl.Finish()

		blockMapper.EXPECT().IsMappingSyncedWithProcessing().Return(false, nil)
		blockMapper.EXPECT().MapNextBlock(context.Background()).Return(nil)

		w := worker{
			processor:   blockMapper,
			waitingTime: time.Nanosecond,
		}
		err := w.executeMapping(context.Background())
		if err != nil {
			t.Errorf("got error '%v' while no error was expected", err)
		}
	})

	t.Run("IsMappingSyncedWithProcessing returns true", func(t *testing.T) {
		t.Parallel()
		ctrl, blockMapper := getMocks(t)
		defer ctrl.Finish()

		blockMapper.EXPECT().IsMappingSyncedWithProcessing().Return(true, nil)

		w := worker{
			processor:   blockMapper,
			waitingTime: time.Nanosecond,
		}
		err := w.executeMapping(context.Background())
		if err != nil {
			t.Errorf("got error '%v' while no error was expected", err)
		}
	})

	t.Run("IsMappingSyncedWithProcessing returns an error", func(t *testing.T) {
		t.Parallel()
		expectedErr := fmt.Errorf("processor error")
		ctrl, blockMapper := getMocks(t)
		defer ctrl.Finish()

		blockMapper.EXPECT().IsMappingSyncedWithProcessing().Return(false, expectedErr)

		w := worker{
			processor:   blockMapper,
			waitingTime: time.Nanosecond,
		}
		err := w.executeMapping(context.Background())
		if err == nil || err.Error() != expectedErr.Error() {
			t.Errorf("got error '%v' while error '%v' was expected", err, expectedErr)
		}
	})

	t.Run("MapNextBlock returns an error", func(t *testing.T) {
		t.Parallel()
		expectedErr := fmt.Errorf("block mapping error")
		ctrl, blockMapper := getMocks(t)
		defer ctrl.Finish()

		blockMapper.EXPECT().IsMappingSyncedWithProcessing().Return(false, nil)
		blockMapper.EXPECT().MapNextBlock(context.Background()).Return(expectedErr)

		w := worker{
			processor:   blockMapper,
			waitingTime: time.Nanosecond,
		}
		err := w.executeMapping(context.Background())
		if err == nil || err.Error() != expectedErr.Error() {
			t.Errorf("got error '%v' while error '%v' was expected", err, expectedErr)
		}
	})
}

func getMocks(t *testing.T) (*gomock.Controller, *blockMapperMock.MockProcessor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return ctrl, blockMapperMock.NewMockProcessor(ctrl)
}
