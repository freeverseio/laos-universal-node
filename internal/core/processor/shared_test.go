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
			t.Errorf("got error message '%s', expected '%s'", err.Error(), expectedErr.Error())
		}
	})
}

func TestGetInitStartingBlock(t *testing.T) {
	t.Parallel()

	t.Run("GetInitStartingBlock happy path", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name                  string
			startingBlockData     model.Block
			userStartingBlock     uint64
			chainLatestBlock      uint64
			expectedStartingBlock uint64
			blockNumberTimes      int
			getLastBlockFunc      func(*stateMock.MockTx) *gomock.Call
			targetFunc            func(*shared.BlockHelper, context.Context) (uint64, error)
		}{
			{
				name:                  "should use ownership starting block from storage",
				startingBlockData:     model.Block{Number: 10},
				userStartingBlock:     0,
				chainLatestBlock:      0,
				expectedStartingBlock: 11,
				blockNumberTimes:      0,
				getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
					return tx.EXPECT().GetLastOwnershipBlock().Return(model.Block{Number: 10}, nil)
				},
				targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
					return b.GetOwnershipInitStartingBlock(c)
				},
			},
			{
				name:                  "should use user provided ownership starting block",
				startingBlockData:     model.Block{Number: 0},
				userStartingBlock:     20,
				chainLatestBlock:      0,
				expectedStartingBlock: 20,
				blockNumberTimes:      0,
				getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
					return tx.EXPECT().GetLastOwnershipBlock().Return(model.Block{Number: 0}, nil)
				},
				targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
					return b.GetOwnershipInitStartingBlock(c)
				},
			},
			{
				name:                  "should use latest block from ownership chain when no starting block provided",
				startingBlockData:     model.Block{Number: 0},
				userStartingBlock:     0,
				chainLatestBlock:      30,
				expectedStartingBlock: 30,
				blockNumberTimes:      1,
				getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
					return tx.EXPECT().GetLastOwnershipBlock().Return(model.Block{Number: 0}, nil)
				},
				targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
					return b.GetOwnershipInitStartingBlock(c)
				},
			},
			{
				name:                  "should use evo starting block from storage",
				startingBlockData:     model.Block{Number: 10},
				userStartingBlock:     0,
				chainLatestBlock:      0,
				expectedStartingBlock: 11,
				blockNumberTimes:      0,
				getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
					return tx.EXPECT().GetLastEvoBlock().Return(model.Block{Number: 10}, nil)
				},
				targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
					return b.GetEvoInitStartingBlock(c)
				},
			},
			{
				name:                  "should use user provided evo starting block",
				startingBlockData:     model.Block{Number: 0},
				userStartingBlock:     20,
				chainLatestBlock:      0,
				expectedStartingBlock: 20,
				blockNumberTimes:      0,
				getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
					return tx.EXPECT().GetLastEvoBlock().Return(model.Block{Number: 0}, nil)
				},
				targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
					return b.GetEvoInitStartingBlock(c)
				},
			},
			{
				name:                  "should use latest block from evo chain when no starting block provided",
				startingBlockData:     model.Block{Number: 0},
				userStartingBlock:     0,
				chainLatestBlock:      30,
				expectedStartingBlock: 30,
				blockNumberTimes:      1,
				getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
					return tx.EXPECT().GetLastEvoBlock().Return(model.Block{Number: 0}, nil)
				},
				targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
					return b.GetEvoInitStartingBlock(c)
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mockClient, mockStateService := getMocks(ctrl)
				tx := stateMock.NewMockTx(ctrl)

				mockStateService.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
				tt.getLastBlockFunc(tx)
				mockClient.EXPECT().BlockNumber(context.Background()).Return(tt.chainLatestBlock, nil).Times(tt.blockNumberTimes)

				helper := shared.NewBlockHelper(mockClient, mockStateService, 100, 10, tt.userStartingBlock)
				actualStartingBlock, err := tt.targetFunc(helper, context.Background())
				if err != nil {
					t.Fatalf("got error '%v' while no error was expected", err)
				}
				if actualStartingBlock != tt.expectedStartingBlock {
					t.Fatalf("got %d, expected starting block %d", tt.expectedStartingBlock, actualStartingBlock)
				}
			})
		}
	})

	t.Run("GetInitStartingBlock errors", func(t *testing.T) {
		t.Parallel()
		t.Run("should return error when creating transaction fails", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				name       string
				targetFunc func(*shared.BlockHelper, context.Context) (uint64, error)
			}{
				{
					name: "on evo init starting block",
					targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
						return b.GetEvoInitStartingBlock(c)
					},
				},
				{
					name: "on ownership init starting block",
					targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
						return b.GetOwnershipInitStartingBlock(c)
					},
				},
			}
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					errMsg := fmt.Errorf("state service failed")
					expectedErr := fmt.Errorf("error creating a new transaction: %w", errMsg)

					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					mockClient, mockStateService := getMocks(ctrl)

					mockStateService.EXPECT().NewTransaction().Return(nil, errMsg)

					helper := shared.NewBlockHelper(mockClient, mockStateService, 100, 10, 0)
					_, err := tt.targetFunc(helper, context.Background())
					if err == nil {
						t.Fatalf("got no error when '%v' was expected", expectedErr)
					}
					if err.Error() != expectedErr.Error() {
						t.Fatalf("got error message '%s', expected '%s'", expectedErr.Error(), err.Error())
					}
				})
			}
		})
		t.Run("should return error when retrieving starting block from storage fails", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				name             string
				targetFunc       func(*shared.BlockHelper, context.Context) (uint64, error)
				getLastBlockFunc func(*stateMock.MockTx) *gomock.Call
			}{
				{
					name: "on evo init starting block",
					targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
						return b.GetEvoInitStartingBlock(c)
					},
					getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
						return tx.EXPECT().GetLastEvoBlock().Return(model.Block{}, fmt.Errorf("storage failed"))
					},
				},
				{
					name: "on ownership init starting block",
					targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
						return b.GetOwnershipInitStartingBlock(c)
					},
					getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
						return tx.EXPECT().GetLastOwnershipBlock().Return(model.Block{}, fmt.Errorf("storage failed"))
					},
				},
			}
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					errMsg := fmt.Errorf("storage failed")
					expectedErr := fmt.Errorf("error retrieving the current block from storage: %w", errMsg)

					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					mockClient, mockStateService := getMocks(ctrl)
					tx := stateMock.NewMockTx(ctrl)

					mockStateService.EXPECT().NewTransaction().Return(tx, nil)
					tx.EXPECT().Discard()
					tt.getLastBlockFunc(tx)

					helper := shared.NewBlockHelper(mockClient, mockStateService, 100, 10, 0)
					_, err := tt.targetFunc(helper, context.Background())
					if err == nil {
						t.Fatalf("got no error when '%v' was expected", expectedErr)
					}
					if err.Error() != expectedErr.Error() {
						t.Fatalf("got error message '%s', expected '%s'", expectedErr.Error(), err.Error())
					}
				})
			}
		})
		t.Run("should return error when retrieving latest block from chain fails", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				name             string
				targetFunc       func(*shared.BlockHelper, context.Context) (uint64, error)
				getLastBlockFunc func(*stateMock.MockTx) *gomock.Call
			}{
				{
					name: "on evo init starting block",
					targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
						return b.GetEvoInitStartingBlock(c)
					},
					getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
						return tx.EXPECT().GetLastEvoBlock().Return(model.Block{}, nil)
					},
				},
				{
					name: "on ownership init starting block",
					targetFunc: func(b *shared.BlockHelper, c context.Context) (uint64, error) {
						return b.GetOwnershipInitStartingBlock(c)
					},
					getLastBlockFunc: func(tx *stateMock.MockTx) *gomock.Call {
						return tx.EXPECT().GetLastOwnershipBlock().Return(model.Block{}, nil)
					},
				},
			}
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					errMsg := fmt.Errorf("node unavailable")
					expectedErr := fmt.Errorf("error retrieving the latest block from chain: %w", errMsg)

					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					mockClient, mockStateService := getMocks(ctrl)
					tx := stateMock.NewMockTx(ctrl)

					mockStateService.EXPECT().NewTransaction().Return(tx, nil)
					tx.EXPECT().Discard()
					tt.getLastBlockFunc(tx)
					mockClient.EXPECT().BlockNumber(context.Background()).Return(uint64(0), errMsg)

					helper := shared.NewBlockHelper(mockClient, mockStateService, 100, 10, 0)
					_, err := tt.targetFunc(helper, context.Background())
					if err == nil {
						t.Fatalf("got no error when '%v' was expected", expectedErr)
					}
					if err.Error() != expectedErr.Error() {
						t.Fatalf("got error message '%s', expected '%s'", expectedErr.Error(), err.Error())
					}
				})
			}
		})
	})
}

func getMocks(ctrl *gomock.Controller) (*blockchainMock.MockEthClient, *stateMock.MockService) {
	mockClient := blockchainMock.NewMockEthClient(ctrl)
	mockStateService := stateMock.NewMockService(ctrl)
	return mockClient, mockStateService
}
