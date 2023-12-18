package evolution_test

import (
	"testing"

	"github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	v1 "github.com/freeverseio/laos-universal-node/internal/state/v1"
	"go.uber.org/mock/gomock"
)

func TestGetCurrentEvoBlock(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockService(mockCtrl)
	mockStorageTransaction := mock.NewMockTx(mockCtrl)
	mockStorage.EXPECT().NewTransaction().Return(mockStorageTransaction)

	stateService := v1.NewStateService(mockStorage)
	tx := stateService.NewTransaction()

	mockStorageTransaction.EXPECT().Get([]byte("evo_current_block")).Return([]byte("1"), nil)

	currentBlock, err := tx.GetCurrentEvoBlock()
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}

	if currentBlock != 1 {
		t.Fatalf("got currentBlock %d, expecting 1", currentBlock)
	}
}

func TestSetCurrentEvoBlock(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockService(mockCtrl)
	mockStorageTransaction := mock.NewMockTx(mockCtrl)
	mockStorage.EXPECT().NewTransaction().Return(mockStorageTransaction)

	stateService := v1.NewStateService(mockStorage)
	tx := stateService.NewTransaction()

	mockStorageTransaction.EXPECT().Set([]byte("evo_current_block"), []byte("1")).Return(nil)

	err := tx.SetCurrentEvoBlock(uint64(1))
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}
}