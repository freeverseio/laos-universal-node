package v1

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	evolutionContractState "github.com/freeverseio/laos-universal-node/internal/platform/state/contract/evolution"
	ownershipContractState "github.com/freeverseio/laos-universal-node/internal/platform/state/contract/ownership"
	evolutionSyncState "github.com/freeverseio/laos-universal-node/internal/platform/state/sync/evolution"
	ownershipSyncState "github.com/freeverseio/laos-universal-node/internal/platform/state/sync/ownership"
	accountTreeMock "github.com/freeverseio/laos-universal-node/internal/platform/state/tree/account/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumerated"
	enumeratedTreeMock "github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumerated/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumeratedtotal"
	enumeratedTotalTreeMock "github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumeratedtotal/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/ownership"
	ownershipTreeMock "github.com/freeverseio/laos-universal-node/internal/platform/state/tree/ownership/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"
)

func TestTransfer(t *testing.T) {
	t.Parallel()

	t.Run(`contracts not set`, func(t *testing.T) {
		t.Parallel()
		memoryService := memory.New()
		stateService := NewStateService(memoryService)

		tx, err := stateService.NewTransaction()
		assert.NilError(t, err)

		eventTransfer := model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		}

		err = tx.Transfer(common.HexToAddress("0x500"), &eventTransfer)
		assert.Error(t, err, "contract 0x0000000000000000000000000000000000000500 does not exist")
	})

	t.Run(`transfer token that is not minted`, func(t *testing.T) {
		t.Parallel()
		ctrl, _, _, ownershipTree, _, transaction := getMocksAndTransaction(t)
		defer ctrl.Finish()

		eventTransfer := model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		}

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x2"), Minted: false, Idx: 0}

		ownershipTree.EXPECT().Transfer(&eventTransfer).Return(nil)
		ownershipTree.EXPECT().TokenData(eventTransfer.TokenId).Return(&tokenData, nil)

		err := transaction.Transfer(common.HexToAddress("0x500"), &eventTransfer)
		assert.NilError(t, err)
	})

	t.Run(`transfer token that is minted`, func(t *testing.T) {
		t.Parallel()
		ctrl, enumeratedTree, _, ownershipTree, _, transaction := getMocksAndTransaction(t)
		defer ctrl.Finish()

		eventTransfer := model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		}

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x2"), Minted: true, Idx: 0}

		ownershipTree.EXPECT().Transfer(&eventTransfer).Return(nil)
		ownershipTree.EXPECT().TokenData(eventTransfer.TokenId).Return(&tokenData, nil)
		enumeratedTree.EXPECT().Transfer(true, &eventTransfer).Return(nil)

		err := transaction.Transfer(common.HexToAddress("0x500"), &eventTransfer)
		assert.NilError(t, err)
	})

	t.Run(`burn token that is minted`, func(t *testing.T) {
		t.Parallel()
		ctrl, enumeratedTree, enumeratedTotalTree, ownershipTree, _, transaction := getMocksAndTransaction(t)
		defer ctrl.Finish()

		eventTransfer := model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.Address{},
			TokenId: big.NewInt(1),
		}

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x2"), Minted: true, Idx: 0}

		ownershipTree.EXPECT().Transfer(&eventTransfer).Return(nil)
		ownershipTree.EXPECT().TokenData(eventTransfer.TokenId).Return(&tokenData, nil)
		enumeratedTree.EXPECT().Transfer(true, &eventTransfer).Return(nil)
		enumeratedTotalTree.EXPECT().TotalSupply().Return(int64(15))
		enumeratedTotalTree.EXPECT().TokenByIndex(14).Return(big.NewInt(10), nil)
		enumeratedTotalTree.EXPECT().Burn(int(0)).Return(nil)

		tokenData2 := ownership.TokenData{SlotOwner: common.HexToAddress("0x3"), Minted: true, Idx: 0}
		ownershipTree.EXPECT().TokenData(big.NewInt(10)).Return(&tokenData2, nil)
		tokenData2.Idx = 0
		ownershipTree.EXPECT().SetTokenData(&tokenData2, big.NewInt(10)).Return(nil)

		err := transaction.Transfer(common.HexToAddress("0x500"), &eventTransfer)
		assert.NilError(t, err)
	})
}

func TestMinting(t *testing.T) {
	t.Parallel()
	t.Run(`mint token`, func(t *testing.T) {
		t.Parallel()
		ctrl, enumeratedTree, enumeratedTotalTree, ownershipTree, _, transaction := getMocksAndTransaction(t)
		defer ctrl.Finish()

		enumeratedTotalTree.EXPECT().Mint(big.NewInt(1)).Return(nil)
		enumeratedTotalTree.EXPECT().TotalSupply().Return(int64(2))

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

		err := transaction.Mint(common.HexToAddress("0x500"), &mintEvent)
		assert.NilError(t, err)
	})
}

func TestTokenURI(t *testing.T) {
	t.Parallel()
	t.Run(`tokenURI returns valid string when asset is minted`, func(t *testing.T) {
		t.Parallel()
		ctrl, _, _, ownershipTree, _, transaction := getMocksAndTransaction(t)
		defer ctrl.Finish()

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x3"), Minted: true, Idx: 1, TokenURI: "tokenURI"}
		tokenId := big.NewInt(1)
		ownershipTree.EXPECT().TokenData(tokenId).Return(&tokenData, nil)

		tokenURI, err := transaction.TokenURI(common.HexToAddress("0x500"), tokenId)
		if err != nil {
			t.Errorf("got error %s when no error was expected", err.Error())
		}
		if tokenURI != tokenData.TokenURI {
			t.Fatalf("got token URI %s, expected %s", tokenURI, tokenData.TokenURI)
		}
	})

	t.Run(`tokenURI returns an error when asset is not minted`, func(t *testing.T) {
		t.Parallel()
		ctrl, _, _, ownershipTree, _, transaction := getMocksAndTransaction(t)
		defer ctrl.Finish()

		tokenData := ownership.TokenData{SlotOwner: common.HexToAddress("0x0"), Minted: false, Idx: 0, TokenURI: ""}
		tokenId := big.NewInt(1)
		ownershipTree.EXPECT().TokenData(tokenId).Return(&tokenData, nil)

		expectedErr := "tokenId 1 does not exist"
		tokenURI, err := transaction.TokenURI(common.HexToAddress("0x500"), tokenId)
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
		ctrl, _, _, _, accountTree, transaction := getMocksAndTransaction(t)
		defer ctrl.Finish()

		accountTree.EXPECT().Checkout(int64(1)).Return(nil)

		err := transaction.Checkout(int64(1))
		assert.NilError(t, err)
	})
}

// nolint:gocritic // it complains about more than five results in return but it is OK for the test
func getMocksAndTransaction(t *testing.T) (ctrl *gomock.Controller,
	enumeratedTree *enumeratedTreeMock.MockTree,
	enumeratedTotalTree *enumeratedTotalTreeMock.MockTree,
	ownershipTree *ownershipTreeMock.MockTree,
	accountTree *accountTreeMock.MockTree,
	transaction tx,
) {
	t.Helper()
	ctrl = gomock.NewController(t)

	memoryService := memory.New()
	storageTx := memoryService.NewTransaction()

	enumeratedTree = enumeratedTreeMock.NewMockTree(ctrl)
	enumeratedTotalTree = enumeratedTotalTreeMock.NewMockTree(ctrl)
	ownershipTree = ownershipTreeMock.NewMockTree(ctrl)
	accountTree = accountTreeMock.NewMockTree(ctrl)

	transaction = tx{
		ownershipTrees:         make(map[common.Address]ownership.Tree),
		enumeratedTrees:        make(map[common.Address]enumerated.Tree),
		enumeratedTotalTrees:   make(map[common.Address]enumeratedtotal.Tree),
		tx:                     storageTx,
		OwnershipContractState: ownershipContractState.NewService(storageTx),
		EvolutionContractState: evolutionContractState.NewService(storageTx),
		OwnershipSyncState:     ownershipSyncState.NewService(storageTx),
		EvolutionSyncState:     evolutionSyncState.NewService(storageTx),
		accountTree:            accountTree,
	}
	transaction.ownershipTrees[common.HexToAddress("0x500")] = ownershipTree
	transaction.enumeratedTrees[common.HexToAddress("0x500")] = enumeratedTree
	transaction.enumeratedTotalTrees[common.HexToAddress("0x500")] = enumeratedTotalTree

	return ctrl, enumeratedTree, enumeratedTotalTree, ownershipTree, accountTree, transaction
}
