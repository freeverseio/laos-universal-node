package repository_test

import (
	"fmt"
	"strings"
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

		mockStorage := mock.NewMockService(mockCtrl)
		service := repository.New(mockStorage)

		universalContracts := []model.ERC721UniversalContract{
			{
				Address:           common.HexToAddress("0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"),
				CollectionAddress: common.HexToAddress("http://baseuri1.com"),
			},
			{
				Address:           common.HexToAddress("0x0"),
				CollectionAddress: common.HexToAddress("http://baseuri2.com"),
			},
		}

		tx := mock.NewMockTx(mockCtrl)
		mockStorage.EXPECT().NewTransaction().Return(tx)

		for _, contract := range universalContracts {
			addressLowerCase := strings.ToLower(contract.Address.String())
			tx.EXPECT().Set(
				[]byte("contract_"+addressLowerCase),
				contract.CollectionAddress.Bytes(),
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

		mockStorage := mock.NewMockService(mockCtrl)
		service := repository.New(mockStorage)

		universalContracts := []model.ERC721UniversalContract{
			{
				Address:           common.HexToAddress("0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"),
				CollectionAddress: common.HexToAddress("http://baseuri1.com"),
			},
			{
				Address:           common.HexToAddress("0x0"),
				CollectionAddress: common.HexToAddress("http://baseuri2.com"),
			},
		}

		tx := mock.NewMockTx(mockCtrl)
		mockStorage.EXPECT().NewTransaction().Return(tx)

		errExpected := fmt.Errorf("error")
		tx.EXPECT().Set(
			[]byte("contract_"+strings.ToLower(universalContracts[0].Address.String())),
			universalContracts[0].CollectionAddress.Bytes(),
		).Return(errExpected)

		tx.EXPECT().Discard().Times(1)
		err := service.StoreERC721UniversalContracts(universalContracts)
		if err == nil {
			t.Fatalf("got no error, expecting en error %s", errExpected.Error())
		}
	})
}

func TestGetAllERC721UniversalContracts(t *testing.T) {
	t.Parallel()
	t.Run("should execute GetAllERC721UniversalContracts without error", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockStorage := mock.NewMockService(mockCtrl)
		service := repository.New(mockStorage)

		keys := [][]byte{
			[]byte("contract_0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			[]byte("contract_0x0"),
		}

		mockStorage.EXPECT().GetKeysWithPrefix([]byte("contract_")).Return(keys, nil)

		contracts, err := service.GetAllERC721UniversalContracts()
		if err != nil {
			t.Fatalf("got error %s, expecting no error", err.Error())
		}

		if len(contracts) != len(keys) {
			t.Fatalf("got %d contracts, expecting %d contracts", len(contracts), len(keys))
		}

		for i := 0; i < len(contracts); i++ {
			contract := strings.TrimPrefix(string(keys[i]), "contract_")
			if contract != contracts[i] {
				t.Fatalf("got contract %s, expecting %s", contracts[i], string(keys[i]))
			}
		}
	})

	t.Run("should execute GetAllERC721UniversalContracts and handle an error", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockStorage := mock.NewMockService(mockCtrl)
		service := repository.New(mockStorage)

		errExpected := fmt.Errorf("error")
		mockStorage.EXPECT().GetKeysWithPrefix([]byte("contract_")).Return([][]byte{}, errExpected)

		_, err := service.GetAllERC721UniversalContracts()
		if err == nil {
			t.Fatalf("got no error, expecting en error %s", errExpected.Error())
		}
	})
}

func TestHasERC721UniversalContract(t *testing.T) {
	t.Parallel()
	t.Run("should execute HasERC721UniversalContract without error", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockStorage := mock.NewMockService(mockCtrl)
		service := repository.New(mockStorage)
		mockStorage.EXPECT().Get([]byte("contract_0x0")).Return([]byte("http://baseuri2.com"), nil)

		has, err := service.HasERC721UniversalContract("0x0")
		if err != nil {
			t.Fatalf("got error %s, expecting no error", err.Error())
		}

		if !has {
			t.Fatalf("got false, expecting true")
		}
	})

	t.Run("should execute HasERC721UniversalContract and handle an error", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockStorage := mock.NewMockService(mockCtrl)
		service := repository.New(mockStorage)
		errExpected := fmt.Errorf("error")
		mockStorage.EXPECT().Get([]byte("contract_0x0")).Return([]byte(""), errExpected)

		_, err := service.HasERC721UniversalContract("0x0")
		if err == nil {
			t.Fatalf("got no error, expecting en error %s", errExpected.Error())
		}
	})
}
