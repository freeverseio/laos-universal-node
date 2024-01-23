package enumeratedtotal_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumeratedtotal"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
	"gotest.tools/assert"
)

func TestTree(t *testing.T) {
	t.Parallel()

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := enumeratedtotal.NewTree(common.Address{}, nil)
		assert.Error(t, err, "contract address is 0x0000000000000000000000000000000000000000")
	})

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), nil)
		assert.Error(t, err, "store is nil")
	})

	t.Run(`initial root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	})

	t.Run(`mint token changes root, error if index is out of bound`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x914210bafc69d69da66c2ca5cfc8b95b375868ae09c07686c0e91ec981cf3c52")

		totalSupply, err := tr.TotalSupply()
		assert.NilError(t, err)

		assert.Equal(t, totalSupply, int64(1))
		token, err := tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		assert.Equal(t, totalSupply, int64(1))
		_, err = tr.TokenByIndex(1)
		assert.Error(t, err, "index out of totalSupply range")
	})

	t.Run(`mint token, burns token token where index is out of bound`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x914210bafc69d69da66c2ca5cfc8b95b375868ae09c07686c0e91ec981cf3c52")

		totalSupply, err := tr.TotalSupply()
		assert.NilError(t, err)

		assert.Equal(t, totalSupply, int64(1))
		token, err := tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		err = tr.Burn(1)
		assert.Error(t, err, "index out of totalSupply range")
	})

	tests := []struct {
		name        string
		idxOfBurned int
		root        string
		tokenLeft   *big.Int
	}{
		{
			name:        "mint two tokens, burn the first token",
			idxOfBurned: 0,
			root:        "0xe760041b06dc8f962d743ac46291f95a623dcf7a8ce3b69543fbecd7beca9be2",
			tokenLeft:   big.NewInt(2),
		},
		{
			name:        "mint two tokens, burn the second token",
			idxOfBurned: 1,
			root:        "0x637c714ed905133d791bb37254e47da8dc068466897fdc266485a3a9cc7bf784",
			tokenLeft:   big.NewInt(1),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := memory.New()
			tx := service.NewTransaction()

			tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), tx)
			assert.NilError(t, err)

			err = tr.Mint(big.NewInt(1))
			assert.NilError(t, err)
			assert.Equal(t, tr.Root().String(), "0x914210bafc69d69da66c2ca5cfc8b95b375868ae09c07686c0e91ec981cf3c52")

			totalSupply, err := tr.TotalSupply()
			assert.NilError(t, err)

			assert.Equal(t, totalSupply, int64(1))
			token, err := tr.TokenByIndex(0)
			assert.NilError(t, err)
			assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

			err = tr.Mint(big.NewInt(2))
			assert.NilError(t, err)
			assert.Equal(t, tr.Root().String(), "0x502374ef4383f50e241caf9870ccf6e832077a3b28a3cdf0f3fcb9e85072f79d")

			totalSupply, err = tr.TotalSupply()
			assert.NilError(t, err)
			assert.Equal(t, totalSupply, int64(2))

			token, err = tr.TokenByIndex(0)
			assert.NilError(t, err)
			assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

			token, err = tr.TokenByIndex(1)
			assert.NilError(t, err)
			assert.Equal(t, token.Cmp(big.NewInt(2)), 0)

			err = tr.Burn(tt.idxOfBurned)
			assert.NilError(t, err)
			assert.Equal(t, tr.Root().String(), tt.root)

			totalSupply, err = tr.TotalSupply()
			assert.NilError(t, err)
			assert.Equal(t, totalSupply, int64(1))

			token, err = tr.TokenByIndex(0)
			assert.NilError(t, err)
			assert.Equal(t, token.Cmp(tt.tokenLeft), 0)
		})
	}

	t.Run(`mint tokens in different contracts`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x914210bafc69d69da66c2ca5cfc8b95b375868ae09c07686c0e91ec981cf3c52")

		totalSupply, err := tr.TotalSupply()
		assert.NilError(t, err)

		assert.Equal(t, totalSupply, int64(1))
		token, err := tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		tr1, err := enumeratedtotal.NewTree(common.HexToAddress("0x501"), tx)
		assert.NilError(t, err)

		err = tr1.Mint(big.NewInt(2))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x914210bafc69d69da66c2ca5cfc8b95b375868ae09c07686c0e91ec981cf3c52")

		totalSupply, err = tr1.TotalSupply()
		assert.NilError(t, err)

		assert.Equal(t, totalSupply, int64(1))
		token, err = tr1.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(2)), 0)

		totalSupply, err = tr.TotalSupply()
		assert.NilError(t, err)

		assert.Equal(t, totalSupply, int64(1))
		token, err = tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)
	})
}

func TestTag(t *testing.T) {
	t.Parallel()
	t.Run(`tag root after mints. checkout at that root returns state before the second mint`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x914210bafc69d69da66c2ca5cfc8b95b375868ae09c07686c0e91ec981cf3c52")

		totalSupply, err := tr.TotalSupply()
		assert.NilError(t, err)
		assert.Equal(t, totalSupply, int64(1))

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(2))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x502374ef4383f50e241caf9870ccf6e832077a3b28a3cdf0f3fcb9e85072f79d")

		err = tr.TagRoot(2)
		assert.NilError(t, err)

		totalSupply, err = tr.TotalSupply()
		assert.NilError(t, err)
		assert.Equal(t, totalSupply, int64(2))

		err = tr.Checkout(1)
		assert.NilError(t, err)

		totalSupply, err = tr.TotalSupply()
		assert.NilError(t, err)
		assert.Equal(t, totalSupply, int64(1))

		token0, err := tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token0.Cmp(big.NewInt(1)), 0)

		err = tr.Checkout(2)
		assert.NilError(t, err)

		totalSupply, err = tr.TotalSupply()
		assert.NilError(t, err)
		assert.Equal(t, totalSupply, int64(2))

		token0, err = tr.TokenByIndex(1)
		assert.NilError(t, err)
		assert.Equal(t, token0.Cmp(big.NewInt(2)), 0)
	})

	t.Run(`tag root before transfer. checkout at block which tag does not exist returns error`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), tx)
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

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		err = tr.TagRoot(2)
		assert.NilError(t, err)
		err = tx.Commit()
		assert.NilError(t, err)
		tx = service.NewTransaction()
		err = enumeratedtotal.DeleteRootTag(tx, common.HexToAddress("0x500").Hex(), 1)
		assert.NilError(t, err)
		err = tx.Commit()
		assert.NilError(t, err)
		tx = service.NewTransaction()
		tr, err = enumeratedtotal.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)
		err = tr.Checkout(1)
		assert.Error(t, err, "no tag found for this block number 1")
	})
}
