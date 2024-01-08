package ownership_test

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
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
