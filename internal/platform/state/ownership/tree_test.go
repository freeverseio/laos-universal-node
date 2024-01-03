package ownership_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/ownership"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
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
		assert.Equal(t, tokenData.TokenURI, "")
		assert.Equal(t, tokenData.Minted, false)

		owner, err := tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.Address{}), 0)

		mintEvent := model.MintedWithExternalURI{
			To:       tokenData.SlotOwner,
			TokenURI: "tokenURI",
			TokenId:  tokenId,
		}

		err = tr.Mint(&mintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x9f53960be4404d6d4308044e6631b2764a120a15ea7c8f2026a6afe290e907e8")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF10")), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
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

		firstMintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x1"),
			TokenURI: "tokenURI",
			TokenId:  big.NewInt(1),
		}
		err = tr.Mint(&firstMintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x6a4dec92176ba3e77f3f92bb8ad68fbb40470448900d833c6e71cf98d1479682")

		tokenData, err := tr.TokenData(firstMintEvent.TokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(firstMintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, firstMintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		secondMintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x2"),
			TokenURI: "tokenURI",
			TokenId:  big.NewInt(2),
		}

		err = tr.Mint(&secondMintEvent, 1)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x342911e2cba943dbcb5f0076d9a752fff74121d7b9be8ddf8554c782a323984e")

		tokenData, err = tr.TokenData(secondMintEvent.TokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(secondMintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, secondMintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 1)
	})

	t.Run(`mint tokens in different contracts`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		mintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x1"),
			TokenURI: "tokenURI",
			TokenId:  big.NewInt(1),
		}

		err = tr.Mint(&mintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x6a4dec92176ba3e77f3f92bb8ad68fbb40470448900d833c6e71cf98d1479682")

		tokenData, err := tr.TokenData(mintEvent.TokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(mintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		tr1, err := ownership.NewTree(common.HexToAddress("0x501"), tx)
		assert.NilError(t, err)

		err = tr1.Mint(&mintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0x6a4dec92176ba3e77f3f92bb8ad68fbb40470448900d833c6e71cf98d1479682")

		tokenData, err = tr1.TokenData(mintEvent.TokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(mintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
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

		assert.Equal(t, tr.Root().String(), "0x4138780c2f786a1b6b5c9b5c58dcc47ca3a0e338a756a96daf7f7fa69300cbb2")
		tokenData, err := tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.TokenURI, "")
		assert.Equal(t, tokenData.Minted, false)

		owner, err := tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.Address{}), 0)

		mintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x2"),
			TokenURI: "tokenURI",
			TokenId:  tokenId,
		}

		err = tr.Mint(&mintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x3c9a06f4499054d9ca5a56f415c38c78c7a896feea00b029fd5b41b4008764c9")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(mintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
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
		mintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x1"),
			TokenURI: "tokenURI",
			TokenId:  tokenId,
		}
		err = tr.Mint(&mintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x6a4dec92176ba3e77f3f92bb8ad68fbb40470448900d833c6e71cf98d1479682")

		tokenData, err := tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(mintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
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
		assert.Equal(t, tr.Root().String(), "0x3c9a06f4499054d9ca5a56f415c38c78c7a896feea00b029fd5b41b4008764c9")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
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
		mintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x1"),
			TokenURI: "tokenURI",
			TokenId:  tokenId,
		}
		err = tr.Mint(&mintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x6a4dec92176ba3e77f3f92bb8ad68fbb40470448900d833c6e71cf98d1479682")

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		tokenData, err := tr.TokenData(tokenId)
		assert.NilError(t, err)

		assert.Equal(t, tokenData.SlotOwner.Cmp(mintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		err = tr.Transfer(&model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: tokenId,
		})

		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x3c9a06f4499054d9ca5a56f415c38c78c7a896feea00b029fd5b41b4008764c9")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)

		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		err = tr.TagRoot(2)
		assert.NilError(t, err)
		err = tr.Checkout(1)
		assert.NilError(t, err)

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)

		assert.Equal(t, tokenData.SlotOwner.Cmp(mintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		err = tr.Checkout(2)
		assert.NilError(t, err)

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)

		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
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
