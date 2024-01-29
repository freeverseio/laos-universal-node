package enumerated_test

import (
	"fmt"
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
	t.Run(`init contract address is common.Hash{} returns error `, func(t *testing.T) {
		t.Parallel()
		_, err := enumerated.NewTree(common.Address{}, common.Hash{}, nil)
		assert.Error(t, err, "contract address is 0x0000000000000000000000000000000000000000")
	})

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := enumerated.NewTree(common.HexToAddress("0x500"), common.Hash{}, nil)
		assert.Error(t, err, "store is nil")
	})

	t.Run(`initial root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), common.HexToHash("0x1"), tx)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000001")
	})

	t.Run(`transfer of the token that is not minted does not change state`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)

		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
		err = tr.Transfer(false, &model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		})
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")

		balance1, err := tr.BalanceOfOwner(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, balance1, uint64(0))
		balance2, err := tr.BalanceOfOwner(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, balance2, uint64(0))

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xdbe504d75c24341e90d36073d631b7cc9ffdca8df2071e14a3c71c5b8c2ffd5b")

		balance1, err = tr.BalanceOfOwner(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, balance1, uint64(0))
		balance2, err = tr.BalanceOfOwner(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, balance2, uint64(1))

		err = tr.Transfer(false, &model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(2),
		})

		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xdbe504d75c24341e90d36073d631b7cc9ffdca8df2071e14a3c71c5b8c2ffd5b")
		balance1, err = tr.BalanceOfOwner(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, balance1, uint64(0))
		balance2, err = tr.BalanceOfOwner(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, balance2, uint64(1))
	})

	t.Run(`mint tokens to address, setting root returns correct data`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x5ce00d2afc3d832a1cd6383355aeb85283a0b0004fad0efc599324ed9057737b")

		balance, err := tr.BalanceOfOwner(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, balance, uint64(1))

		token, err := tr.TokenOfOwnerByIndex(common.HexToAddress("0x1"), 0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		err = tr.Mint(big.NewInt(2), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xdbaa67bf186eb370fafb60f542854bec45b56b2f6a29d83d6d206fddf5a7f8bd")

		balance, err = tr.BalanceOfOwner(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, balance, uint64(2))

		token, err = tr.TokenOfOwnerByIndex(common.HexToAddress("0x1"), 0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		token, err = tr.TokenOfOwnerByIndex(common.HexToAddress("0x1"), 1)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(2)), 0)

		tr.SetRoot(common.HexToHash("0x5ce00d2afc3d832a1cd6383355aeb85283a0b0004fad0efc599324ed9057737b"))
		balance, err = tr.BalanceOfOwner(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, balance, uint64(1))

		token, err = tr.TokenOfOwnerByIndex(common.HexToAddress("0x1"), 0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		_, err = tr.TokenOfOwnerByIndex(common.HexToAddress("0x1"), 1)
		assert.Error(t, err, "index 1 out of range")
	})

	t.Run(`tokens minted in different contracts`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x5ce00d2afc3d832a1cd6383355aeb85283a0b0004fad0efc599324ed9057737b")

		tr1, err := enumerated.NewTree(common.HexToAddress("0x501"), common.Hash{}, tx)
		assert.NilError(t, err)

		err = tr1.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0x5ce00d2afc3d832a1cd6383355aeb85283a0b0004fad0efc599324ed9057737b")
	})

	t.Run(`transfer token  works correctly`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), common.Hash{}, tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x5ce00d2afc3d832a1cd6383355aeb85283a0b0004fad0efc599324ed9057737b")

		balance, err := tr.BalanceOfOwner(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, balance, uint64(1))

		token, err := tr.TokenOfOwnerByIndex(common.HexToAddress("0x1"), 0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)

		balance, err = tr.BalanceOfOwner(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, balance, uint64(0))

		err = tr.Transfer(true, &model.ERC721Transfer{
			From:    common.HexToAddress("0x1"),
			To:      common.HexToAddress("0x2"),
			TokenId: big.NewInt(1),
		})
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x5a8a15f0f6e2e2c551e0dadfb31c1f1d2f08377328eb20812e6e3df86b979218")

		balance, err = tr.BalanceOfOwner(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, balance, uint64(0))

		_, err = tr.TokenOfOwnerByIndex(common.HexToAddress("0x1"), 0)
		assert.Error(t, fmt.Errorf("index 0 out of range"), err.Error())

		balance, err = tr.BalanceOfOwner(common.HexToAddress("0x2"))
		assert.NilError(t, err)
		assert.Equal(t, balance, uint64(1))

		token, err = tr.TokenOfOwnerByIndex(common.HexToAddress("0x2"), 0)
		assert.NilError(t, err)
		assert.Equal(t, token.Cmp(big.NewInt(1)), 0)
	})
}
