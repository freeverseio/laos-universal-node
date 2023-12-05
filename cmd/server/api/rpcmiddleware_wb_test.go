package api

import (
	"testing"

	mockState "github.com/freeverseio/laos-universal-node/internal/state/mock"
	"go.uber.org/mock/gomock"
)

func TestCheckContractInList(t *testing.T) {
	t.Parallel() // Run tests in parallel
	t.Run("Should checkContractInList and get the expected result", func(t *testing.T) {
		t.Parallel() // Run tests in parallel
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})
		contract1 := "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"
		contract2 := "0x36CB70039FE1bd36b4659858d4c4D0cBcafd743A"
		stateService := mockState.NewMockService(ctrl)
		tx := mockState.NewMockTx(ctrl)

		stateService.EXPECT().NewTransaction().Return(tx).Times(2)
		tx.EXPECT().Discard().Times(2)
		tx.EXPECT().HasERC721UniversalContract(contract1).Return(true, nil).Times(1)
		tx.EXPECT().HasERC721UniversalContract(contract2).Return(false, nil).Times(1)

		b, err := isContractStored(contract1, stateService)
		if err != nil {
			t.Fatalf("got %T, expected nil error", err)
		}
		if !b {
			t.Fatalf("got %v, expected true", b)
		}

		b, err = isContractStored(contract2, stateService)
		if err != nil {
			t.Fatalf("got %T, expected nil error", err)
		}
		if b {
			t.Fatalf("got %v, expected false", b)
		}
	})
}
