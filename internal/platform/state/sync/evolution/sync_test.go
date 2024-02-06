package evolution_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/sync"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	"go.uber.org/mock/gomock"
)

func TestSetGetLastEvoBlock(t *testing.T) {
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

	mockStorageTransaction.EXPECT().Set([]byte("evo_last_block"), buf.Bytes()).Return(nil)

	err := tx.SetLastEvoBlock(block)
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}
	mockStorageTransaction.EXPECT().Get([]byte("evo_last_block")).Return(buf.Bytes(), nil)

	newBlock, err := tx.GetLastEvoBlock()
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

func TestSetGetEvoEventOwnershipBlockNumber(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockService(mockCtrl)
	mockStorageTransaction := mock.NewMockTx(mockCtrl)
	mockStorage.EXPECT().NewTransaction().Return(mockStorageTransaction)

	stateService := v1.NewStateService(mockStorage)
	tx := stateService.NewTransaction()

	contract := common.HexToAddress("0x123")

	ownershipBlockNumber := uint64(1)
	event := model.MintedWithExternalURI{
		BlockNumber: 1,
		TxIndex:     1,
	}
	key := fmt.Sprintf("%s%s_%s_%d", "evo_event_", sync.FormatNumberForSorting(event.BlockNumber, 18), contract.Hex(), event.TxIndex)
	value := strconv.FormatUint(ownershipBlockNumber, 10)
	mockStorageTransaction.EXPECT().Set([]byte(key), []byte(value)).Return(nil)

	err := tx.SetEvoEventOwnershipBlockNumber(ownershipBlockNumber, &event, contract.Hex())
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}

	mockStorageTransaction.EXPECT().Get([]byte(key)).Return([]byte("1"), nil)

	blockNumber, err := tx.GetEvoEventOwnershipBlockNumber(event.BlockNumber, event.TxIndex, contract.Hex())
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}

	if blockNumber != ownershipBlockNumber {
		t.Fatalf("got block number %d, expecting %d", blockNumber, ownershipBlockNumber)
	}
}
