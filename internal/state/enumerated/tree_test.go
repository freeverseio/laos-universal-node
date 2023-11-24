package enumerated_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
	"github.com/freeverseio/laos-universal-node/internal/state/enumerated"
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
		assert.Equal(t, tr.Root().String(), "0xd2c0580ba1c042026990043e8c9c4191e0e41b3893d6d1fd2fe8192c8e5b8dbb")

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
		assert.Equal(t, tr.Root().String(), "0xd2c0580ba1c042026990043e8c9c4191e0e41b3893d6d1fd2fe8192c8e5b8dbb")
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
		assert.Equal(t, tr.Root().String(), "0x0390a835b13dc7d39fac829657704b2b93fad32bdd80643debfcd1726a8dd166")

		tokens, err := tr.TokensOf(common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, len(tokens), 1)
		assert.Equal(t, tokens[0].Cmp(big.NewInt(1)), 0)

		err = tr.Mint(big.NewInt(2), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xd16d265b662c771be40c2bb76db2f3e7520d42c203594369f7dcd516a2b18743")

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
		assert.Equal(t, tr.Root().String(), "0x0390a835b13dc7d39fac829657704b2b93fad32bdd80643debfcd1726a8dd166")

		tr1, err := enumerated.NewTree(common.HexToAddress("0x501"), tx)
		assert.NilError(t, err)

		err = tr1.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0x0390a835b13dc7d39fac829657704b2b93fad32bdd80643debfcd1726a8dd166")
	})

	t.Run(`transfer token  works correctly`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Mint(big.NewInt(1), common.HexToAddress("0x1"))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0390a835b13dc7d39fac829657704b2b93fad32bdd80643debfcd1726a8dd166")

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
		assert.Equal(t, tr.Root().String(), "0x7b0a894e5d132e9c6d0f09fc9ff8bf67d280335be2d298dd918294e3ae76f213")

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
		assert.Equal(t, tr.Root().String(), "0x0390a835b13dc7d39fac829657704b2b93fad32bdd80643debfcd1726a8dd166")

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
		assert.Equal(t, tr.Root().String(), "0x7b0a894e5d132e9c6d0f09fc9ff8bf67d280335be2d298dd918294e3ae76f213")
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
	t.Run(`Find the first tag that has the same state as current block number`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
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

		tr, err := enumerated.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		blockNumber, err := tr.FindBlockWithTag(4)
		assert.NilError(t, err)
		assert.Equal(t, blockNumber, int64(0))
	})
}
