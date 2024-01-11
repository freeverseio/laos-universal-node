package ownership_test

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/sync/ownership"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	"go.uber.org/mock/gomock"
)

func TestSetGetLastOwnershipBlock(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockService(mockCtrl)
	mockStorageTransaction := mock.NewMockTx(mockCtrl)
	mockStorage.EXPECT().NewTransaction().Return(mockStorageTransaction)

	stateService := v1.NewStateService(mockStorage)
	tx := stateService.NewTransaction()

	block := model.Block{
		Number:    1,
		Timestamp: 1,
		Hash:      common.HexToHash("0x123"),
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	_ = encoder.Encode(block) // omit error since block is constant

	mockStorageTransaction.EXPECT().Set([]byte("ownership_last_block"), buf.Bytes()).Return(nil)
	mockStorageTransaction.EXPECT().Set([]byte("ownership_block_1"), buf.Bytes())

	err := tx.SetLastOwnershipBlock(block)
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}
	mockStorageTransaction.EXPECT().Get([]byte("ownership_last_block")).Return(buf.Bytes(), nil)

	newBlock, err := tx.GetLastOwnershipBlock()
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}

	if newBlock.Number != block.Number {
		t.Fatalf("got block number %d, expecting %d", newBlock.Number, block.Number)
	}

	if newBlock.Timestamp != block.Timestamp {
		t.Fatalf("got block timestamp %d, expecting %d", newBlock.Timestamp, block.Timestamp)
	}

	if newBlock.Hash != block.Hash {
		t.Fatalf("got block hash %s, expecting %s", newBlock.Hash.String(), block.Hash.String())
	}
}

func TestSetGetCurrentEvoEventsIndexForOwnershipContract(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockService(mockCtrl)
	mockStorageTransaction := mock.NewMockTx(mockCtrl)
	mockStorage.EXPECT().NewTransaction().Return(mockStorageTransaction)

	stateService := v1.NewStateService(mockStorage)
	tx := stateService.NewTransaction()

	contract := "0x123"
	mockStorageTransaction.EXPECT().Set([]byte("ownership_contract_evo_current_index_"+contract), []byte("50")).Return(nil)

	err := tx.SetCurrentEvoEventsIndexForOwnershipContract(contract, uint64(50))
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}
	mockStorageTransaction.EXPECT().Get([]byte("ownership_contract_evo_current_index_"+contract)).Return([]byte("50"), nil)

	result, err := tx.GetCurrentEvoEventsIndexForOwnershipContract(contract)
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}

	if result != 50 {
		t.Fatalf("got %d, expecting %d", result, 50)
	}
}

func TestGetAllStoredBlockNumbers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		mockBlockNumbers []string
		expectedNumbers  []uint64
		expectError      bool
	}{
		{
			name:             "SingleBlockNumber",
			mockBlockNumbers: []string{"ownership_block_1"},
			expectedNumbers:  []uint64{1},
		},
		{
			name:             "MultipleBlockNumbers",
			mockBlockNumbers: []string{"ownership_block_01", "ownership_block_02", "ownership_block_03"},
			expectedNumbers:  []uint64{1, 2, 3},
		},
		{
			name:             "NoBlockNumbers",
			mockBlockNumbers: []string{},
			expectedNumbers:  []uint64{},
		},

		{
			name:             "WithErrorInvalidNumber",
			mockBlockNumbers: []string{"ownership_block_1", "ownership_block_a"},
			expectedNumbers:  nil,
			expectError:      true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			mockStorageTransaction := mock.NewMockTx(mockCtrl)
			defer mockCtrl.Finish()
			mockStorageTransaction.EXPECT().
				GetKeysWithPrefix([]byte("ownership_block_")).
				Return(convertToByteSliceArray(tc.mockBlockNumbers))

			service := ownership.NewService(mockStorageTransaction)
			numbers, err := service.GetAllStoredBlockNumbers()
			if tc.expectError {
				if err == nil {
					t.Errorf("expected an error in test case %s, but got none", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("got error %v, expected no error in test case %s", err, tc.name)
				}

				if !compareSlices(numbers, tc.expectedNumbers) {
					t.Errorf("got %v, expected %v in test case %s", numbers, tc.expectedNumbers, tc.name)
				}
			}
		})
	}
}

func TestSetLastOwnershipBlock(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                      string
		blockNumber               uint64
		expectedOwnershipBlockTag string
	}{
		{
			name:                      "Block 1",
			blockNumber:               1,
			expectedOwnershipBlockTag: "ownership_block_000000000000000001",
		},
		{
			name:                      "Block 2",
			blockNumber:               2,
			expectedOwnershipBlockTag: "ownership_block_000000000000000002",
		},
		{
			name:                      "Block 3",
			blockNumber:               3,
			expectedOwnershipBlockTag: "ownership_block_000000000000000003",
		},
		{
			name:                      "Block 1254",
			blockNumber:               1254,
			expectedOwnershipBlockTag: "ownership_block_000000000000001254",
		},
		{
			name:                      "Blocknumer with 18 digits",
			blockNumber:               123654258965487545,
			expectedOwnershipBlockTag: "ownership_block_123654258965487545",
		},
		{
			name:                      "Blocknumer with more than 18 digits",
			blockNumber:               8888888754587958787,
			expectedOwnershipBlockTag: "ownership_block_8888888754587958787",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			mockStorageTransaction := mock.NewMockTx(mockCtrl)
			defer mockCtrl.Finish()

			block := model.Block{
				Number:    tc.blockNumber,
				Timestamp: 1,
				Hash:      common.HexToHash("0x123"),
			}

			var buf bytes.Buffer
			encoder := gob.NewEncoder(&buf)
			_ = encoder.Encode(block) // omit error since block is constant

			mockStorageTransaction.EXPECT().Set([]byte("ownership_last_block"), buf.Bytes()).Return(nil)
			mockStorageTransaction.EXPECT().Set([]byte(tc.expectedOwnershipBlockTag), buf.Bytes()).Return(nil)

			service := ownership.NewService(mockStorageTransaction)
			err := service.SetLastOwnershipBlock(block)
			if err != nil {
				t.Errorf("got error %v, expected no error in test case %s", err, tc.name)
			}
		})
	}
}

func convertToByteSliceArray(strs []string) [][]byte {
	var result [][]byte
	for _, s := range strs {
		result = append(result, []byte(s))
	}
	return result
}

func compareSlices(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
