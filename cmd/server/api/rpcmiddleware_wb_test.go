package api

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
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
		storageMock := mock.NewMockStorage(ctrl)
		storageMock.EXPECT().ReadAll(context.Background()).Return([]scan.ERC721UniversalContract{
			{
				Address: common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
				Block:   uint64(0),
				BaseURI: "evochain1/collectionId/",
			},
		}, nil).Times(2)

		b, err := checkContractInList("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", storageMock)
		if err != nil {
			t.Fatalf("got %T, expected nil error", err)
		}
		if !b {
			t.Fatalf("got %v, expected true", b)
		}

		b, err = checkContractInList("0x36CB70039FE1bd36b4659858d4c4D0cBcafd743A", storageMock)
		if err != nil {
			t.Fatalf("got %T, expected nil error", err)
		}
		if b {
			t.Fatalf("got %v, expected false", b)
		}
	})
}
