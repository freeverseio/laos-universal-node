package ownership_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/ownership"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
	"gotest.tools/assert"
)

func TestTree(t *testing.T) {
	t.Parallel()

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := ownership.NewTree(common.Address{}, common.Hash{}, nil)
		assert.Error(t, err, "contract address is 0x0000000000000000000000000000000000000000")
	})

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := ownership.NewTree(common.HexToAddress("0x500"), common.Hash{}, nil)
		assert.Error(t, err, "store is nil")
	})

	t.Run(`initial root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	})

	t.Run(`check initial owner of the token`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
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
		assert.Equal(t, tr.Root().String(), "0x27427fdf247f6ee7ffccdb61f480fe56235f48d2c8ff798b99113e4b69e8c797")

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

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)

		firstMintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x1"),
			TokenURI: "tokenURI",
			TokenId:  big.NewInt(1),
		}
		err = tr.Mint(&firstMintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55")

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
		assert.Equal(t, tr.Root().String(), "0x7070510f5faefe9a8c20a7897b579753c4a2d3a5d08332a3fce27a96a48110a0")

		tokenData, err = tr.TokenData(secondMintEvent.TokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(secondMintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, secondMintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 1)
	})

	t.Run(`mint same token 2 times returns an error`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)

		firstMintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x1"),
			TokenURI: "tokenURI",
			TokenId:  big.NewInt(1),
		}
		err = tr.Mint(&firstMintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55")

		tokenData, err := tr.TokenData(firstMintEvent.TokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(firstMintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, firstMintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		secondMintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x2"),
			TokenURI: "tokenURI",
			TokenId:  big.NewInt(1),
		}

		err = tr.Mint(&secondMintEvent, 1)
		assert.Error(t, err, "token 1 already minted")
	})

	t.Run(`mint tokens in different contracts`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)

		mintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x1"),
			TokenURI: "tokenURI",
			TokenId:  big.NewInt(1),
		}

		err = tr.Mint(&mintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55")

		tokenData, err := tr.TokenData(mintEvent.TokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(mintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		tr1, err := ownership.NewTree(common.HexToAddress("0x501"), common.Hash{}, tx)
		assert.NilError(t, err)

		err = tr1.Mint(&mintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55")

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

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)

		tokenId := big.NewInt(1)
		err = tr.Transfer(&model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: tokenId,
		})
		assert.NilError(t, err)

		assert.Equal(t, tr.Root().String(), "0x28992c39a19a0af63b5787393dfe40691fee4a023065a8901e7aced8114940c9")
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
		assert.Equal(t, tr.Root().String(), "0x02d716d8efa798c67be692351475fbeb8025fad04ccf08308191628897aaf2fe")

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

	t.Run(`mint token and then transfer it. set root to the one before transfer, owner is correct`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := ownership.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)

		tokenId := big.NewInt(1)
		mintEvent := model.MintedWithExternalURI{
			To:       common.HexToAddress("0x1"),
			TokenURI: "tokenURI",
			TokenId:  tokenId,
		}
		err = tr.Mint(&mintEvent, 0)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55")

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
		assert.Equal(t, tr.Root().String(), "0x02d716d8efa798c67be692351475fbeb8025fad04ccf08308191628897aaf2fe")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(common.HexToAddress("0x2")), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		owner, err = tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.HexToAddress("0x2")), 0)

		tr.SetRoot(common.HexToHash("0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55"))
		assert.Equal(t, tr.Root().String(), "0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55")

		tokenData, err = tr.TokenData(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, tokenData.SlotOwner.Cmp(mintEvent.To), 0)
		assert.Equal(t, tokenData.TokenURI, mintEvent.TokenURI)
		assert.Equal(t, tokenData.Minted, true)
		assert.Equal(t, tokenData.Idx, 0)

		owner, err = tr.OwnerOf(tokenId)
		assert.NilError(t, err)
		assert.Equal(t, owner.Cmp(common.HexToAddress("0x1")), 0)
	})
}
