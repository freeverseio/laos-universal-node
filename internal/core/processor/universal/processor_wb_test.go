package universal

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	mockDiscoverer "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer/mock"
	mockUpdater "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/updater/mock"
	mockClient "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	mockScan "github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

func TestGetNextLowerBlockNumber(t *testing.T) {
	testCases := []struct {
		name                  string
		currentBlock          uint64
		storedBlockNumbers    []uint64
		expectedBlockNumber   uint64
		expectedModifiedSlice []uint64
		expectedFound         bool
	}{
		{
			name:                  "FindLowerBlock",
			currentBlock:          5,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 7},
			expectedBlockNumber:   4,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 7},
			expectedFound:         true,
		},
		{
			name:                  "FindLowerBlock",
			currentBlock:          7,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 7},
			expectedBlockNumber:   6,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 7},
			expectedFound:         true,
		},
		{
			name:                  "FindLowerBlock",
			currentBlock:          9,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 8},
			expectedBlockNumber:   8,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 8},
			expectedFound:         true,
		},
		{
			name:                  "BlockNotFound",
			currentBlock:          3,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 8},
			expectedBlockNumber:   0,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 8},
			expectedFound:         false,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			gotBlockNumber, found := getNextLowerBlockNumber(tc.currentBlock, tc.storedBlockNumbers)

			if tc.expectedFound != found {
				t.Errorf("got %v, expected %v", found, tc.expectedFound)
			}
			if gotBlockNumber != tc.expectedBlockNumber {
				t.Errorf("got %v, expected %v", gotBlockNumber, tc.expectedBlockNumber)
			}
			if !reflect.DeepEqual(tc.storedBlockNumbers, tc.expectedModifiedSlice) {
				t.Errorf("slice was modified to %v, expected %v", tc.storedBlockNumbers, tc.expectedModifiedSlice)
			}
		})
	}
}

func TestCheckBlockForReorg(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		startingBlock          uint64
		lastBlockDB            model.Block
		previousBlockData      *types.Header
		previousBlockDataError error
		expectedError          error
	}{
		{
			name:          "No previous block hash in database",
			startingBlock: 100,
			lastBlockDB:   model.Block{},
			expectedError: fmt.Errorf("no hash stored in the database for block %d", 0),
		},
		{
			name:          "Previous block hash matches",
			startingBlock: 100,
			lastBlockDB: model.Block{
				Number: 99,
				Hash:   common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e"),
			},
			previousBlockData: &types.Header{ParentHash: common.HexToHash("0x123")},
			expectedError:     nil,
		},

		{
			name:          "error when trying to obtain previous block from chain",
			startingBlock: 100,
			lastBlockDB: model.Block{
				Number: 99,
				Hash:   common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e"),
			},
			previousBlockData:      nil,
			previousBlockDataError: errors.New("error retrieving previous block from chain"),
			expectedError:          errors.New("error retrieving previous block from chain"),
		},
		{
			name:          "Previous block hash does not match",
			startingBlock: 100,
			lastBlockDB: model.Block{
				Number: 99,
				Hash:   common.HexToHash("0x123"),
			},
			previousBlockData:      &types.Header{ParentHash: common.HexToHash("0x123")},
			previousBlockDataError: nil,
			expectedError: ReorgError{
				Block:       99,
				ChainHash:   common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e"),
				StorageHash: common.HexToHash("0x123"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.TODO()
			stateService, _, client, _, _, _ := createMocks(t)

			if tt.lastBlockDB.Hash != (common.Hash{}) {
				client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.startingBlock-1))).
					Return(tt.previousBlockData, tt.previousBlockDataError)
			}

			p := NewProcessor(client, stateService, nil, &config.Config{}, nil, nil)
			err := p.checkBlockForReorg(ctx, tt.lastBlockDB)
			assertError(t, tt.expectedError, err)
		})
	}
}

// nolint:gocritic // many return values in function => we accept this for this test helper
func createMocks(t *testing.T) (
	*mockTx.MockService,
	*mockTx.MockTx,
	*mockClient.MockEthClient,
	*mockScan.MockScanner,
	*mockDiscoverer.MockDiscoverer,
	*mockUpdater.MockUpdater,
) {
	ctrl := gomock.NewController(t)
	return mockTx.NewMockService(ctrl), mockTx.NewMockTx(ctrl), mockClient.NewMockEthClient(ctrl), mockScan.NewMockScanner(ctrl), mockDiscoverer.NewMockDiscoverer(ctrl), mockUpdater.NewMockUpdater(ctrl)
}

func assertError(t *testing.T, expectedError, err error) {
	t.Helper()
	if expectedError != nil {
		if err == nil || err.Error() != expectedError.Error() {
			t.Fatalf(`got error "%v", expected error: "%v"`, err, expectedError)
		}
	} else {
		if err != expectedError {
			t.Fatalf(`got error "%v", expected error: "%v"`, err, expectedError)
		}
	}
}
