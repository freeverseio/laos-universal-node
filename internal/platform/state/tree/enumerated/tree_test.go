package enumerated_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumerated"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
	"gotest.tools/assert"
)

func TestTree(t *testing.T) {
	t.Parallel()
	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := enumerated.NewTree(common.Address{}, nil)
		assert.Error(t, err, "contract address is 0x0000000000000000000000000000000000000000")
	})

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := enumerated.NewTree(common.HexToAddress("0x500"), nil)
		assert.Error(t, err, "store is nil")
	})

	t.Run(`initial root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	})

	t.Run(`transfer of the token that is not minted does not change state`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
		err = tr.Transfer(false, &model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		})
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")

		tokens1, err := tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens1), 0)
		tokens2, err := tr.TokensOf(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens2), 0)

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2869b4a1411aa86d60d937c481cb6b4843432fc6342efef13fbc82d1b2bd9db5")

		tokens1, err = tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens1), 0)
		tokens2, err = tr.TokensOf(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens2), 1)

		err = tr.Transfer(false, &model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(2),
		})

		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2869b4a1411aa86d60d937c481cb6b4843432fc6342efef13fbc82d1b2bd9db5")
		tokens1, err = tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens1), 0)
		tokens2, err = tr.TokensOf(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens2), 1)
	})

	t.Run(`mint tokens to address`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xd2776ed971a71a483a279ade441a20cb67374963ba95fef6874ab6f7cfa8a63a")

		tokens, err := tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens), 1)
		assert.Equal(t, tokens[0].Cmp(big.NewInt(1)), 0)

		err = tr.Mint(big.NewInt(2), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xe89b5bd33d239dde2e9d298c7ffea488eb63a0bd44b9fa39cba3deb383d470ec")

		tokens, err = tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens), 2)
		assert.Equal(t, tokens[0].Cmp(big.NewInt(1)), 0)
		assert.Equal(t, tokens[1].Cmp(big.NewInt(2)), 0)
	})

	t.Run(`tokens minted in different contracts`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xd2776ed971a71a483a279ade441a20cb67374963ba95fef6874ab6f7cfa8a63a")

		tr1, err := enumerated.NewTree(common.HexToAddress("0x501"), tx)
		assert.NilError(t, err)

		err = tr1.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0xd2776ed971a71a483a279ade441a20cb67374963ba95fef6874ab6f7cfa8a63a")
	})

	t.Run(`transfer token  works correctly`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xd2776ed971a71a483a279ade441a20cb67374963ba95fef6874ab6f7cfa8a63a")

		tokens1, err := tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens1), 1)
		assert.Equal(t, tokens1[0].Cmp(big.NewInt(1)), 0)

		tokens2, err := tr.TokensOf(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens2), 0)

		err = tr.Transfer(true, &model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		})
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xb899cb3fcae2117e1eaffa7652259f348844200f51d1b00712ff677869d8f5ca")

		tokens1, err = tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens1), 0)

		tokens2, err = tr.TokensOf(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens2), 1)
		assert.Equal(t, tokens2[0].Cmp(big.NewInt(1)), 0)
	})
}

func TestTag(t *testing.T) {
	t.Parallel()
	t.Run(`tag root before transfer. checkout at that root returns state before transfer`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xd2776ed971a71a483a279ade441a20cb67374963ba95fef6874ab6f7cfa8a63a")

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		tokens1, err := tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens1), 1)
		assert.Equal(t, tokens1[0].Cmp(big.NewInt(1)), 0)

		tokens2, err := tr.TokensOf(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens2), 0)

		err = tr.Transfer(true, &model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		})
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xb899cb3fcae2117e1eaffa7652259f348844200f51d1b00712ff677869d8f5ca")
		err = tr.TagRoot(2)
		assert.NilError(t, err)

		tokens1, err = tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens1), 0)

		tokens2, err = tr.TokensOf(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens2), 1)
		assert.Equal(t, tokens2[0].Cmp(big.NewInt(1)), 0)

		err = tr.Checkout(1)
		assert.NilError(t, err)

		tokens1, err = tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens1), 1)
		assert.Equal(t, tokens1[0].Cmp(big.NewInt(1)), 0)

		tokens2, err = tr.TokensOf(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens2), 0)

		err = tr.Checkout(2)
		assert.NilError(t, err)

		tokens1, err = tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens1), 0)

		tokens2, err = tr.TokensOf(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens2), 1)
		assert.Equal(t, tokens2[0].Cmp(big.NewInt(1)), 0)
	})

	t.Run(`tag root before transfer. checkout at block which tag does not exist returns error`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Checkout(1)
		assert.Error(t, err, "no tag found for this block number 1")
	})
}

func TestDeleteRootTag(t *testing.T) {
	t.Parallel()
	t.Run(`Tag two roots and then delete the first tag. Checkout at deleted tag gives error`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		err = tr.TagRoot(2)
		assert.NilError(t, err)
		err = tx.Commit()
		assert.NilError(t, err)
		tx = service.NewTransaction()
		err = enumerated.DeleteRootTag(tx, common.HexToAddress("0x500").Hex(), 1)
		assert.NilError(t, err)
		err = tx.Commit()
		assert.NilError(t, err)
		tx = service.NewTransaction()
		tr, err = enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)
		err = tr.Checkout(1)
		assert.Error(t, err, "no tag found for this block number 1")
	})
}
