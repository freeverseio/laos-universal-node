package shared_test

import (
	"context"
	"errors"
	"testing"

	shared "github.com/freeverseio/laos-universal-node/internal/core/processor"
	blockchainMock "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	stateMock "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

func TestGetLastBlock(t *testing.T) {
	t.Parallel()
	t.Run("GetLastBlock happy path", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name        string
			blockNumber uint64
			expected    uint64
		}{
			{"returns the left side of the minimum", 20, 11},
			{"returns the right side of the minimum", 10, 10},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockClient, mockStateService := getMocks(ctrl)

				blockHelper := shared.NewBlockHelper(mockClient, mockStateService, 10, 0, 1)
				mockClient.EXPECT().BlockNumber(context.Background()).Return(tt.blockNumber, nil)

				got, err := blockHelper.GetLastBlock(context.Background(), 1)
				if err != nil {
					t.Fatalf("got error '%v' while no error was expected", err)
				}
				if got != tt.expected {
					t.Fatalf("got %d, expected %d", got, tt.expected)
				}
			})
		}
	})
	t.Run("handles an error from BlockNumber", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient, mockStateService := getMocks(ctrl)

		blockHelper := shared.NewBlockHelper(mockClient, mockStateService, 10, 0, 1)
		expectedErr := errors.New("an error occurred")
		mockClient.EXPECT().BlockNumber(context.Background()).Return(uint64(0), expectedErr)

		_, err := blockHelper.GetLastBlock(context.Background(), 1)
		if err == nil {
			t.Errorf("got no error while an error was expected")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("got error message '%s' while '%s' was expected", err.Error(), expectedErr.Error())
		}
	})
}

func getMocks(ctrl *gomock.Controller) (*blockchainMock.MockEthClient, *stateMock.MockService) {
	mockClient := blockchainMock.NewMockEthClient(ctrl)
	mockStateService := stateMock.NewMockService(ctrl)
	return mockClient, mockStateService
}
