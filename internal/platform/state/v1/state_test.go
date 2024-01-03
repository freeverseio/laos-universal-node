package v1_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	enumeratedTreeMock "github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumerated/mock"
	enumeratedTotalTreeMock "github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumeratedtotal/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/ownership"
	ownershipTreeMock "github.com/freeverseio/laos-universal-node/internal/platform/state/tree/ownership/mock"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"
)

func TestTransfer(t *testing.T) {
	t.Parallel()

	t.Run(`contracts not set`, func(t *testing.T) {
		t.Parallel()
		t.Helper()

		memoryService := memory.New()
		stateService := v1.NewStateService(memoryService)

		tx := stateService.NewTransaction()

		eventTransfer := model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		}

		err := tx.Transfer(common.HexToAddress("0x500"), &eventTransfer)
		assert.Error(t, err, "contract 0x0000000000000000000000000000000000000500 does not exist")
	})

	t.Run(`transfer token that is not minted`, func(t *testing.T) {
		t.Parallel()
		t.Helper()
		ctrl := gomock.NewController(t)

		memoryService := memory.New()
		stateService := v1.NewStateService(memoryService)

		tx := stateService.NewTransaction()

		enumeratedTree := enumeratedTreeMock.NewMockTree(ctrl)
		enumeratedTotalTree := enumeratedTotalTreeMock.NewMockTree(ctrl)
		ownershipTree := ownershipTreeMock.NewMockTree(ctrl)

		tx.SetTreesForContract(common.HexToAddress("0x500"), ownershipTree, enumeratedTree, enumeratedTotalTree)

		eventTransfer := model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		}

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x2"), Minted: false, Idx: 0}

		ownershipTree.EXPECT().Transfer(&eventTransfer).Return(nil)
		ownershipTree.EXPECT().TokenData(eventTransfer.TokenId).Return(&tokenData, nil)

		err := tx.Transfer(common.HexToAddress("0x500"), &eventTransfer)
		assert.NilError(t, err)
	})

	t.Run(`transfer token that is minted`, func(t *testing.T) {
		t.Parallel()
		t.Helper()
		ctrl := gomock.NewController(t)

		memoryService := memory.New()
		stateService := v1.NewStateService(memoryService)

		tx := stateService.NewTransaction()

		enumeratedTree := enumeratedTreeMock.NewMockTree(ctrl)
		enumeratedTotalTree := enumeratedTotalTreeMock.NewMockTree(ctrl)
		ownershipTree := ownershipTreeMock.NewMockTree(ctrl)

		tx.SetTreesForContract(common.HexToAddress("0x500"), ownershipTree, enumeratedTree, enumeratedTotalTree)

		eventTransfer := model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		}

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x2"), Minted: true, Idx: 0}

		ownershipTree.EXPECT().Transfer(&eventTransfer).Return(nil)
		ownershipTree.EXPECT().TokenData(eventTransfer.TokenId).Return(&tokenData, nil)
		enumeratedTree.EXPECT().Transfer(true, &eventTransfer).Return(nil)

		err := tx.Transfer(common.HexToAddress("0x500"), &eventTransfer)
		assert.NilError(t, err)
	})

	t.Run(`burn token that is minted`, func(t *testing.T) {
		t.Parallel()
		t.Helper()
		ctrl := gomock.NewController(t)

		memoryService := memory.New()
		stateService := v1.NewStateService(memoryService)

		tx := stateService.NewTransaction()

		enumeratedTree := enumeratedTreeMock.NewMockTree(ctrl)
		enumeratedTotalTree := enumeratedTotalTreeMock.NewMockTree(ctrl)
		ownershipTree := ownershipTreeMock.NewMockTree(ctrl)

		tx.SetTreesForContract(common.HexToAddress("0x500"), ownershipTree, enumeratedTree, enumeratedTotalTree)

		eventTransfer := model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.Address{},
			TokenId: big.NewInt(1),
		}

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x2"), Minted: true, Idx: 0}

		ownershipTree.EXPECT().Transfer(&eventTransfer).Return(nil)
		ownershipTree.EXPECT().TokenData(eventTransfer.TokenId).Return(&tokenData, nil)
		enumeratedTree.EXPECT().Transfer(true, &eventTransfer).Return(nil)
		enumeratedTotalTree.EXPECT().TotalSupply().Return(int64(15), nil)
		enumeratedTotalTree.EXPECT().TokenByIndex(14).Return(big.NewInt(10), nil)
		enumeratedTotalTree.EXPECT().Burn(int(0)).Return(nil)

		tokenData2 := ownership.TokenData{SlotOwner: common.HexToAddress("0x3"), Minted: true, Idx: 0}
		ownershipTree.EXPECT().TokenData(big.NewInt(10)).Return(&tokenData2, nil)
		tokenData2.Idx = 0
		ownershipTree.EXPECT().SetTokenData(&tokenData2, big.NewInt(10)).Return(nil)

		err := tx.Transfer(common.HexToAddress("0x500"), &eventTransfer)
		assert.NilError(t, err)
	})
}

func TestMinting(t *testing.T) {
	t.Parallel()
	t.Run(`mint token`, func(t *testing.T) {
		t.Parallel()
		t.Helper()
		ctrl := gomock.NewController(t)

		memoryService := memory.New()
		stateService := v1.NewStateService(memoryService)

		tx := stateService.NewTransaction()

		enumeratedTree := enumeratedTreeMock.NewMockTree(ctrl)
		enumeratedTotalTree := enumeratedTotalTreeMock.NewMockTree(ctrl)
		ownershipTree := ownershipTreeMock.NewMockTree(ctrl)

		tx.SetTreesForContract(common.HexToAddress("0x500"), ownershipTree, enumeratedTree, enumeratedTotalTree)

		enumeratedTotalTree.EXPECT().Mint(big.NewInt(1)).Return(nil)
		enumeratedTotalTree.EXPECT().TotalSupply().Return(int64(2), nil)

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x3"), Minted: true, Idx: 1, TokenURI: "tokenURI"}
		ownershipTree.EXPECT().TokenData(big.NewInt(1)).Return(&tokenData, nil)

		enumeratedTree.EXPECT().Mint(big.NewInt(1), tokenData.SlotOwner).Return(nil)

		mintEvent := model.MintedWithExternalURI{
			Slot:        big.NewInt(1),
			To:          tokenData.SlotOwner,
			TokenURI:    tokenData.TokenURI,
			TokenId:     big.NewInt(1),
			BlockNumber: 100,
			Timestamp:   1000,
		}

		ownershipTree.EXPECT().Mint(&mintEvent, 1).Return(nil)

		err := tx.Mint(common.HexToAddress("0x500"), &mintEvent)
		assert.NilError(t, err)
	})
}

func TestTokenURI(t *testing.T) {
	t.Parallel()
	t.Run(`tokenURI returns valid string when asset is minted`, func(t *testing.T) {
		t.Parallel()
		// TODO move the repeated setup code for test data to a helper method?
		ctrl := gomock.NewController(t)

		memoryService := memory.New()
		stateService := v1.NewStateService(memoryService)

		tx := stateService.NewTransaction()

		enumeratedTree := enumeratedTreeMock.NewMockTree(ctrl)
		enumeratedTotalTree := enumeratedTotalTreeMock.NewMockTree(ctrl)
		ownershipTree := ownershipTreeMock.NewMockTree(ctrl)

		tx.SetTreesForContract(common.HexToAddress("0x500"), ownershipTree, enumeratedTree, enumeratedTotalTree)

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x3"), Minted: true, Idx: 1, TokenURI: "tokenURI"}
		tokenId := big.NewInt(1)
		ownershipTree.EXPECT().TokenData(tokenId).Return(&tokenData, nil)

		tokenURI, err := tx.TokenURI(common.HexToAddress("0x500"), tokenId)
		if err != nil {
			t.Errorf("got error %s when no error was expected", err.Error())
		}
		if tokenURI != tokenData.TokenURI {
			t.Fatalf("got token URI %s, expected %s", tokenURI, tokenData.TokenURI)
		}
	})

	t.Run(`tokenURI returns an error when asset is not minted`, func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		memoryService := memory.New()
		stateService := v1.NewStateService(memoryService)

		tx := stateService.NewTransaction()

		enumeratedTree := enumeratedTreeMock.NewMockTree(ctrl)
		enumeratedTotalTree := enumeratedTotalTreeMock.NewMockTree(ctrl)
		ownershipTree := ownershipTreeMock.NewMockTree(ctrl)

		tx.SetTreesForContract(common.HexToAddress("0x500"), ownershipTree, enumeratedTree, enumeratedTotalTree)

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x0"), Minted: false, Idx: 0, TokenURI: ""}
		tokenId := big.NewInt(1)
		ownershipTree.EXPECT().TokenData(tokenId).Return(&tokenData, nil)

		expectedErr := "tokenId 1 does not exist"
		tokenURI, err := tx.TokenURI(common.HexToAddress("0x500"), tokenId)
		if err == nil {
			t.Error("got no error when error was expected")
		}
		if err.Error() != expectedErr {
			t.Errorf("got error %s, expected %s", err.Error(), expectedErr)
		}
		if tokenURI != tokenData.TokenURI {
			t.Fatalf("got token URI %s, expected %s", tokenURI, tokenData.TokenURI)
		}
	})
}

func TestCheckout(t *testing.T) {
	t.Parallel()
	t.Run(`test checkout`, func(t *testing.T) {
		t.Parallel()
		t.Helper()
		ctrl := gomock.NewController(t)

		memoryService := memory.New()
		stateService := v1.NewStateService(memoryService)

		tx := stateService.NewTransaction()

		enumeratedTree := enumeratedTreeMock.NewMockTree(ctrl)
		enumeratedTotalTree := enumeratedTotalTreeMock.NewMockTree(ctrl)
		ownershipTree := ownershipTreeMock.NewMockTree(ctrl)

		tx.SetTreesForContract(common.HexToAddress("0x500"), ownershipTree, enumeratedTree, enumeratedTotalTree)

		enumeratedTotalTree.EXPECT().Checkout(int64(1)).Return(nil)
		enumeratedTree.EXPECT().Checkout(int64(1)).Return(nil)
		ownershipTree.EXPECT().Checkout(int64(1)).Return(nil)

		err := tx.Checkout(common.HexToAddress("0x500"), int64(1))
		assert.NilError(t, err)
	})
}
