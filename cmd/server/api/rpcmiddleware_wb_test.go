package api

import (
	"testing"

	mockStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	"github.com/freeverseio/laos-universal-node/internal/repository"
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

		storage := mockStorage.NewMockStorage(ctrl)
		keys := [][]byte{
			[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
		}

		storage.EXPECT().GetKeysWithPrefix([]byte("contract_")).Return(keys, nil).Times(2)
		repositoryService := repository.New(storage)

		b, err := isContractInList("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", repositoryService)
		if err != nil {
			t.Fatalf("got %T, expected nil error", err)
		}
		if !b {
			t.Fatalf("got %v, expected true", b)
		}

		b, err = isContractInList("0x36CB70039FE1bd36b4659858d4c4D0cBcafd743A", repositoryService)
		if err != nil {
			t.Fatalf("got %T, expected nil error", err)
		}
		if b {
			t.Fatalf("got %v, expected false", b)
		}
	})
}
