package ownership_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
	"github.com/freeverseio/laos-universal-node/internal/state/ownership"
	"gotest.tools/assert"
)

func TestTree(t *testing.T) {
	t.Parallel()

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := ownership.NewTree(common.Address{}, nil)
		assert.Error(t, err, "contract address is 0x0000000000000000000000000000000000000000")
	})

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := ownership.NewTree(common.HexToAddress("0x500"), nil)
		assert.Error(t, err, "store is nil")
	})

	t.Run(`initial root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	})

	t.Run(`check initial owner of the token`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		tokenId, success := new(big.Int).SetString("0101FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF10", 16)
		assert.Equal(t, success, true)

		tokenData, err := tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF10")), 0)
		assert.Equal(t, tokenData.Minted, false)

		owner, err := tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.Address{}), 0)

		err = tr.Mint(tokenId, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x68d0eae8c6603616c80449c10ad99c6dc94c67398db33ce94930cb4d544eb618")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF10")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		owner, err = tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.HexToAddress("0x01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF10")), 0)
	})

	t.Run(`mint 2 tokens`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x42715022a8e7406f75319c7e031b2e55010e2bf57229acac2092bf31054aca6b")

		tokenData, err := tr.TokenData(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x1")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		err = tr.Mint(big.NewInt(2), 1)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x1619b4647bd8d33851e6d33c580a738597eae5d5dad79f256603561138bfa4aa")

		tokenData, err = tr.TokenData(big.NewInt(2))
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 1)
	})

	t.Run(`mint tokens in different contracts`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x42715022a8e7406f75319c7e031b2e55010e2bf57229acac2092bf31054aca6b")

		tokenData, err := tr.TokenData(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x1")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		tr1, err := ownership.NewTree(common.HexToAddress("0x501"), tx)
		assert.NilError(t, err)

		err = tr1.Mint(big.NewInt(1), 0)
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0x42715022a8e7406f75319c7e031b2e55010e2bf57229acac2092bf31054aca6b")

		tokenData, err = tr1.TokenData(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x1")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)
	})

	t.Run(`transfer token slot and then mint token`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		tokenId := big.NewInt(1)
		err = tr.Transfer(&model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: tokenId,
		})
		assert.NilError(t, err)

		assert.Equal(t, tr.Root().String(), "0x5644aa80f8ded20942de9ea943b74160865a55337e4df5ffe0196a4e27f619a9")
		tokenData, err := tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.Minted, false)

		owner, err := tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.Address{}), 0)

		err = tr.Mint(tokenId, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2f0b86bf3ad4921288012a2a77b3225fd671692e17b85f1a461dc33c381883fa")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		owner, err = tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.HexToAddress("0x2")), 0)
	})

	t.Run(`mint token and then transfer it`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		tokenId := big.NewInt(1)
		err = tr.Mint(tokenId, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x42715022a8e7406f75319c7e031b2e55010e2bf57229acac2092bf31054aca6b")

		tokenData, err := tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x1")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		owner, err := tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.HexToAddress("0x1")), 0)

		err = tr.Transfer(&model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: tokenId,
		})
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2f0b86bf3ad4921288012a2a77b3225fd671692e17b85f1a461dc33c381883fa")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		owner, err = tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.HexToAddress("0x2")), 0)
	})
}

func TestTag(t *testing.T) {
	t.Parallel()
	t.Run(`tag root before transfer. checkout at that root returns state before transfer`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		tokenId := big.NewInt(1)
		err = tr.Mint(tokenId, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x42715022a8e7406f75319c7e031b2e55010e2bf57229acac2092bf31054aca6b")

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		tokenData, err := tr.TokenData(tokenId)
		assert.NilError(t, err)

		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x1")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		err = tr.Transfer(&model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: tokenId,
		})

		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2f0b86bf3ad4921288012a2a77b3225fd671692e17b85f1a461dc33c381883fa")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)

		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		err = tr.TagRoot(2)
		assert.NilError(t, err)
		err = tr.Checkout(1)
		assert.NilError(t, err)

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)

		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x1")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		err = tr.Checkout(2)
		assert.NilError(t, err)

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)

		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)
	})

	t.Run(`tag root before transfer. checkout at block which tag does not exist returns error`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Checkout(1)
		assert.Error(t, err, "no tag found for this block number 1")
	})
	t.Run(`Find the first tag that has the same state as current block number`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		err = tr.TagRoot(2)
		assert.NilError(t, err)

		blockNumber, err := tr.FindBlockWithTag(4)
		assert.NilError(t, err)
		assert.Equal(t, blockNumber, int64(2))
	})

	t.Run(`Find the first tag that has the same state as current block number. no tags return 0`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		blockNumber, err := tr.FindBlockWithTag(4)
		assert.NilError(t, err)
		assert.Equal(t, blockNumber, int64(0))
	})

	t.Run(`Tag two roots and then delete the first tag. Checkout at deleted tag gives error`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		err = tr.TagRoot(2)
		assert.NilError(t, err)

		err = tr.DeleteRootTag(1)
		assert.NilError(t, err)

		err = tr.Checkout(1)
		assert.Error(t, err, "no tag found for this block number 1")
	})
}
