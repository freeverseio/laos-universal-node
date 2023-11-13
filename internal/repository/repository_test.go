package repository_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	"github.com/freeverseio/laos-universal-node/internal/repository"
	"go.uber.org/mock/gomock"
)

func TestStoreERC721UniversalContracts(t *testing.T) {
	t.Parallel()
	t.Run("should execute StoreERC721UniversalContracts without error", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockStorage := mock.NewMockStorage(mockCtrl)
		service := repository.New(mockStorage)

		universalContracts := []model.ERC721UniversalContract{
			{
				Address: common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
				BaseURI: "http://baseuri1.com",
			},
			{
				Address: common.HexToAddress("0x0"),
				BaseURI: "http://baseuri2.com",
			},
		}

		tx := mock.NewMockTx(mockCtrl)
		mockStorage.EXPECT().NewTransaction().Return(tx)

		for _, contract := range universalContracts {
			tx.EXPECT().Set(
				[]byte("contract_"+contract.Address.String()),
				[]byte(contract.BaseURI),
			).Return(nil)
		}

		tx.EXPECT().Commit().Return(nil)
		tx.EXPECT().Discard().Times(1)
		err := service.StoreERC721UniversalContracts(universalContracts)
		if err != nil {
			t.Fatalf("got error %s, expecting no error", err.Error())
		}
	})

	t.Run("should execute StoreERC721UniversalContracts and handle an error", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockStorage := mock.NewMockStorage(mockCtrl)
		service := repository.New(mockStorage)

		universalContracts := []model.ERC721UniversalContract{
			{
				Address: common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
				BaseURI: "http://baseuri1.com",
			},
			{
				Address: common.HexToAddress("0x0"),
				BaseURI: "http://baseuri2.com",
			},
		}

		tx := mock.NewMockTx(mockCtrl)
		mockStorage.EXPECT().NewTransaction().Return(tx)

		errExpected := fmt.Errorf("error")
		tx.EXPECT().Set(
			[]byte("contract_"+universalContracts[0].Address.String()),
			[]byte(universalContracts[0].BaseURI),
		).Return(errExpected)

		tx.EXPECT().Discard().Times(1)
		err := service.StoreERC721UniversalContracts(universalContracts)
		if err == nil {
			t.Fatalf("got no error, expecting en error %s", errExpected.Error())
		}
	})
}
