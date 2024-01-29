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

	t.Run(`init contract address is common.Hash{} returns error `, func(t *testing.T) {
		t.Parallel()
		_, err := enumeratedtotal.NewTree(common.Address{}, common.Hash{}, 0, nil)
		assert.Error(t, err, "contract address is 0x0000000000000000000000000000000000000000")
	})

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), common.Hash{}, 0, nil)
		assert.Error(t, err, "store is nil")
	})

	t.Run(`initial root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), common.HexToHash("0x1"), 0, tx)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000001")
	})

	t.Run(`mint token changes root, error if index is out of bound`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), common.Hash{}, 0, tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2a7c66cf9e4638104cbf07da77fe051a8aa94a675bb1c539d113052cdff1b0aa")

		totalSupply := tr.TotalSupply()

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

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), common.Hash{}, 0, tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2a7c66cf9e4638104cbf07da77fe051a8aa94a675bb1c539d113052cdff1b0aa")

		totalSupply := tr.TotalSupply()
		assert.Equal(t, totalSupply, int64(1))

		token, err := tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		err = tr.Burn(1)
		assert.Error(t, err, "index out of totalSupply range")
	})

	t.Run(`mint 2 tokens, set root after the first mint, second token does not exist`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), common.Hash{}, 0, tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2a7c66cf9e4638104cbf07da77fe051a8aa94a675bb1c539d113052cdff1b0aa")

		totalSupply := tr.TotalSupply()
		assert.Equal(t, totalSupply, int64(1))

		token, err := tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		err = tr.Mint(big.NewInt(2))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x3505457d4944236492ccb69056852e218f55eba5aed2adad48e5309f9339fcef")

		totalSupply = tr.TotalSupply()
		assert.Equal(t, totalSupply, int64(2))

		token, err = tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		token, err = tr.TokenByIndex(1)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(2)), 0)

		tr.SetRoot(common.HexToHash("0x2a7c66cf9e4638104cbf07da77fe051a8aa94a675bb1c539d113052cdff1b0aa"))
		tr.SetTotalSupply(1)
		totalSupply = tr.TotalSupply()
		assert.Equal(t, totalSupply, int64(1))

		_, err = tr.TokenByIndex(1)
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
			root:        "0xa7a87cf365b9a7452c155dd084bc0f7d38b3ab1929f9c010166f97e68906b0c9",
			tokenLeft:   big.NewInt(2),
		},
		{
			name:        "mint two tokens, burn the second token",
			idxOfBurned: 1,
			root:        "0x71db8e96fee6f61f4d680e123b6b8f7ac068a478ec50bced4eb2bb37d0cf32f8",
			tokenLeft:   big.NewInt(1),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := memory.New()
			tx := service.NewTransaction()

			tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), common.Hash{}, 0, tx)
			assert.NilError(t, err)

			err = tr.Mint(big.NewInt(1))
			assert.NilError(t, err)
			assert.Equal(t, tr.Root().String(), "0x2a7c66cf9e4638104cbf07da77fe051a8aa94a675bb1c539d113052cdff1b0aa")

			totalSupply := tr.TotalSupply()
			assert.Equal(t, totalSupply, int64(1))

			token, err := tr.TokenByIndex(0)
			assert.NilError(t, err)
			assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

			err = tr.Mint(big.NewInt(2))
			assert.NilError(t, err)
			assert.Equal(t, tr.Root().String(), "0x3505457d4944236492ccb69056852e218f55eba5aed2adad48e5309f9339fcef")

			totalSupply = tr.TotalSupply()
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

			totalSupply = tr.TotalSupply()
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

		tr, err := enumeratedtotal.NewTree(common.HexToAddress("0x500"), common.Hash{}, 0, tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2a7c66cf9e4638104cbf07da77fe051a8aa94a675bb1c539d113052cdff1b0aa")

		totalSupply := tr.TotalSupply()
		assert.Equal(t, totalSupply, int64(1))

		token, err := tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		tr1, err := enumeratedtotal.NewTree(common.HexToAddress("0x501"), common.Hash{}, 0, tx)
		assert.NilError(t, err)

		err = tr1.Mint(big.NewInt(2))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x2a7c66cf9e4638104cbf07da77fe051a8aa94a675bb1c539d113052cdff1b0aa")

		totalSupply = tr1.TotalSupply()
		assert.Equal(t, totalSupply, int64(1))

		token, err = tr1.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(2)), 0)

		totalSupply = tr.TotalSupply()
		assert.Equal(t, totalSupply, int64(1))

		token, err = tr.TokenByIndex(0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)
	})
}
