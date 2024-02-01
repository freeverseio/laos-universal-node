package evolution_test

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
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
	mockStorageTransaction.EXPECT().Get([]byte("accounthead/")).Return(nil, nil)
	stateService := v1.NewStateService(mockStorage)
	tx, err := stateService.NewTransaction()
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}

	block := model.Block{
		Number:    1,
		Timestamp: 1,
		Hash:      common.HexToHash("0x123"),
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	_ = encoder.Encode(block) // omit error since block is constant

	mockStorageTransaction.EXPECT().Set([]byte("evo_last_block"), buf.Bytes()).Return(nil)

	err = tx.SetLastEvoBlock(block)
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

func TestSetGetNextNextEvoEventBlock(t *testing.T) {
	t.Parallel()
	t.Run("stores minted events", func(t *testing.T) {
		t.Parallel()
		db := createBadger(t)
		tx, err := createBadgerTransaction(t, db)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		err = tx.SetNextEvoEventBlock(common.HexToAddress("0x500").Hex(), 100)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		block, err := tx.GetNextEvoEventBlock(common.HexToAddress("0x500").Hex(), 0)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		if block != 100 {
			t.Errorf(`got block %d when 100 was expected`, block)
		}
	})
}

func createBadgerTransaction(t *testing.T, db *badger.DB) (state.Tx, error) {
	t.Helper()
	badgerService := badgerStorage.NewService(db)
	stateService := v1.NewStateService(badgerService)
	return stateService.NewTransaction()
}

func createBadger(t *testing.T) *badger.DB {
	t.Helper()
	db, err := badger.Open(
		badger.DefaultOptions("").
			WithInMemory(true).
			WithLoggingLevel(badger.ERROR).WithMemTableSize(1 << 30))
	if err != nil {
		t.Fatalf("error initializing storage: %v", err)
	}

	return db
}
