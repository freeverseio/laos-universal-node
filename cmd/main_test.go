package main

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"go.uber.org/mock/gomock"

	mockStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	"github.com/freeverseio/laos-universal-node/internal/repository"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
)

func TestCompareChainIDs(t *testing.T) {
	t.Parallel()
	ethChainID := big.NewInt(3) // Example Ethereum network chain ID
	dbChainID := "3"            // Matching chain ID for the database

	tests := []struct {
		name               string
		ethClientChainID   *big.Int
		ethClientError     error
		dbChainID          string
		dbError            error
		expectSetInDB      bool
		expectedDBSetError error
		wantError          bool
	}{
		{
			name:             "chain IDs match",
			ethClientChainID: ethChainID,
			dbChainID:        dbChainID,
			wantError:        false,
		},
		{
			name:           "error from Ethereum client",
			ethClientError: errors.New("client error"),
			wantError:      true,
		},
		{
			name:      "error from database on get",
			dbError:   errors.New("db get error"),
			wantError: true,
		},
		{
			name:               "chain ID not set in database",
			ethClientChainID:   ethChainID,
			dbChainID:          "",
			expectSetInDB:      true,
			expectedDBSetError: nil,
			wantError:          false,
		},
		{
			name:             "mismatched chain IDs",
			ethClientChainID: ethChainID,
			dbChainID:        "2",
			wantError:        true,
		},
		{
			name:               "error from database on set",
			ethClientChainID:   ethChainID,
			dbChainID:          "",
			expectSetInDB:      true,
			expectedDBSetError: errors.New("db set error"),
			wantError:          true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockClient, _, storage := getMocks(t)
			repositoryService := repository.New(storage)
			ctx := context.Background()

			mockClient.EXPECT().ChainID(ctx).Return(tt.ethClientChainID, tt.ethClientError).Times(1)
			if tt.ethClientError == nil {
				storage.EXPECT().Get([]byte("chain_id")).Return([]byte(tt.dbChainID), tt.dbError).Times(1)
			}
			if tt.expectSetInDB {
				storage.EXPECT().Set([]byte("chain_id"), []byte(ethChainID.String())).Return(tt.expectedDBSetError)
			}

			err := compareChainIDs(ctx, mockClient, repositoryService)

			if (err != nil) != tt.wantError {
				t.Errorf("compareChainIDs() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func getMocks(t *testing.T) (*mock.MockEthClient, *mock.MockScanner, *mockStorage.MockService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockEthClient(ctrl), mock.NewMockScanner(ctrl),
		mockStorage.NewMockService(ctrl)
}
