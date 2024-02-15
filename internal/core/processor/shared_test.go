package shared_test

import (
	"context"
	"fmt"
	"testing"

	shared "github.com/freeverseio/laos-universal-node/internal/core/processor"
	blockchainMock "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
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
			{"should return the left side of the minimum", 20, 11},
			{"should return the right side of the minimum", 10, 10},
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
	t.Run("should return the error from BlockNumber", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient, mockStateService := getMocks(ctrl)

		blockHelper := shared.NewBlockHelper(mockClient, mockStateService, 10, 0, 1)
		expectedErr := fmt.Errorf("an error occurred")
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

func TestGetOwnershipInitStartingBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient, mockStateService := getMocks(ctrl)

	tests := []struct {
		name                  string
		startingBlockData     model.Block
		userStartingBlock     uint64
		chainLatestBlock      uint64
		expectedStartingBlock uint64
		blockNumberTimes      int
	}{
		{
			name:                  "should use starting block from storage",
			startingBlockData:     model.Block{Number: 10},
			userStartingBlock:     0,
			chainLatestBlock:      0,
			expectedStartingBlock: 11,
			blockNumberTimes:      0,
		},
		{
			name:                  "should use user provided starting block",
			startingBlockData:     model.Block{Number: 0},
			userStartingBlock:     20,
			chainLatestBlock:      0,
			expectedStartingBlock: 20,
			blockNumberTimes:      0,
		},
		{
			name:                  "should use latest block from chain when no starting block provided",
			startingBlockData:     model.Block{Number: 0},
			userStartingBlock:     0,
			chainLatestBlock:      30,
			expectedStartingBlock: 30,
			blockNumberTimes:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := stateMock.NewMockTx(ctrl)
			mockStateService.EXPECT().NewTransaction().Return(tx, nil)
			tx.EXPECT().Discard().Times(1)
			tx.EXPECT().GetLastOwnershipBlock().Return(tt.startingBlockData, nil).Times(1)
			mockClient.EXPECT().BlockNumber(context.Background()).Return(tt.chainLatestBlock, nil).Times(tt.blockNumberTimes)

			helper := shared.NewBlockHelper(mockClient, mockStateService, 100, 10, tt.userStartingBlock)
			actualStartingBlock, err := helper.GetOwnershipInitStartingBlock(context.Background())
			if err != nil {
				t.Errorf("got error '%v' while no error was expected", err)
			}
			if actualStartingBlock != tt.expectedStartingBlock {
				t.Errorf("got %d, expected starting block %d", tt.expectedStartingBlock, actualStartingBlock)
			}
		})
	}
}

func getMocks(ctrl *gomock.Controller) (*blockchainMock.MockEthClient, *stateMock.MockService) {
	mockClient := blockchainMock.NewMockEthClient(ctrl)
	mockStateService := stateMock.NewMockService(ctrl)
	return mockClient, mockStateService
}
