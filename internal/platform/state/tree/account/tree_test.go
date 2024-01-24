package account_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"gotest.tools/assert"

	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/account"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
)

func TestTree(t *testing.T) {
	t.Parallel()

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := account.NewTree(common.Address{}, nil)
		assert.Error(t, err, "contract address is 0x0000000000000000000000000000000000000000")
	})

	t.Run(`init with nil store should fail`, func(t *testing.T) {
		t.Parallel()
		_, err := account.NewTree(common.HexToAddress("0x500"), nil)
		assert.Error(t, err, "store is nil")
	})

	t.Run(`initial root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := account.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	})

	t.Run(`set  merkle roots in different contracts`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr1, err := account.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		roots1 := account.MerkleTreeRoots{
			Enumerated:      common.HexToHash("0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55"),
			EnumeratedTotal: common.HexToHash("0xcf46e158742177f61d06cb049d82c7d4aeb7420205d0e1c1bacc45406acde8f3"),
			Ownership:       common.HexToHash("0x59d7de0b77f377095267336d574c03d9c444c5cbbdfe03997e16aa1ff0df6798"),
		}

		err = tr1.SetMerkleTreeRoots(&roots1, big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0xd2dd7011bbd026ffc2340f1e6eb5e7d6a8e7756e5985eb63d860bbbf93762a92")

		leafData1, err := tr1.MerkleTreeRoots(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, leafData1.Enumerated.Cmp(roots1.Enumerated), 0)
		assert.Equal(t, leafData1.EnumeratedTotal.Cmp(roots1.EnumeratedTotal), 0)
		assert.Equal(t, leafData1.Ownership.Cmp(roots1.Ownership), 0)

		tr2, err := account.NewTree(common.HexToAddress("0x501"), tx)
		assert.NilError(t, err)

		roots2 := account.MerkleTreeRoots{
			Enumerated:      common.HexToHash("0x0c5edb6df5fef722190139079b98d5e5fe4bbaf3eed6d0f2dbed5444609a7f27"),
			EnumeratedTotal: common.HexToHash("0x253016ee3bf6160f4115c34ce91e48fa67c8930ad313a300ff531a7c650509a1"),
			Ownership:       common.HexToHash("0x3e54e3ce880dc57d7e2daedd8c7f5f902653567df67e2da5e84385ab743a5283"),
		}

		err = tr2.SetMerkleTreeRoots(&roots2, big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr2.Root().String(), "0xc87bbe8f01faab566e288a46d5358e6978decff545eb84009a23c0c198a34aaa")

		leafData2, err := tr2.MerkleTreeRoots(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, leafData2.Enumerated.Cmp(roots2.Enumerated), 0)
		assert.Equal(t, leafData2.EnumeratedTotal.Cmp(roots2.EnumeratedTotal), 0)
		assert.Equal(t, leafData2.Ownership.Cmp(roots2.Ownership), 0)
	})
}

func TestTag(t *testing.T) {
	t.Parallel()
	t.Run(`tag root before setting merkle roots. checkout at that root returns state before setting merkle roots`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := account.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		accountID := big.NewInt(1)
		roots1 := account.MerkleTreeRoots{
			Enumerated:      common.HexToHash("0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55"),
			EnumeratedTotal: common.HexToHash("0xcf46e158742177f61d06cb049d82c7d4aeb7420205d0e1c1bacc45406acde8f3"),
			Ownership:       common.HexToHash("0x59d7de0b77f377095267336d574c03d9c444c5cbbdfe03997e16aa1ff0df6798"),
		}

		err = tr.SetMerkleTreeRoots(&roots1, big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xd2dd7011bbd026ffc2340f1e6eb5e7d6a8e7756e5985eb63d860bbbf93762a92")

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		leafData, err := tr.MerkleTreeRoots(big.NewInt(1))
		assert.NilError(t, err)

		assert.Equal(t, leafData.Enumerated.Cmp(roots1.Enumerated), 0)
		assert.Equal(t, leafData.EnumeratedTotal.Cmp(roots1.EnumeratedTotal), 0)
		assert.Equal(t, leafData.Ownership.Cmp(roots1.Ownership), 0)

		roots2 := account.MerkleTreeRoots{
			Enumerated:      common.HexToHash("0x0c5edb6df5fef722190139079b98d5e5fe4bbaf3eed6d0f2dbed5444609a7f27"),
			EnumeratedTotal: common.HexToHash("0x253016ee3bf6160f4115c34ce91e48fa67c8930ad313a300ff531a7c650509a1"),
			Ownership:       common.HexToHash("0x3e54e3ce880dc57d7e2daedd8c7f5f902653567df67e2da5e84385ab743a5283"),
		}

		err = tr.SetMerkleTreeRoots(&roots2, accountID)
		assert.NilError(t, err)

		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0xc87bbe8f01faab566e288a46d5358e6978decff545eb84009a23c0c198a34aaa")

		leafData2, err := tr.MerkleTreeRoots(accountID)
		assert.NilError(t, err)

		assert.Equal(t, leafData2.Enumerated.Cmp(roots2.Enumerated), 0)
		assert.Equal(t, leafData2.EnumeratedTotal.Cmp(roots2.EnumeratedTotal), 0)
		assert.Equal(t, leafData2.Ownership.Cmp(roots2.Ownership), 0)

		err = tr.TagRoot(2)
		assert.NilError(t, err)

		err = tr.Checkout(1)
		assert.NilError(t, err)

		leafData, err = tr.MerkleTreeRoots(big.NewInt(1))
		assert.NilError(t, err)

		assert.Equal(t, leafData.Enumerated.Cmp(roots1.Enumerated), 0)
		assert.Equal(t, leafData.EnumeratedTotal.Cmp(roots1.EnumeratedTotal), 0)
		assert.Equal(t, leafData.Ownership.Cmp(roots1.Ownership), 0)

		err = tr.Checkout(2)
		assert.NilError(t, err)

		leafData, err = tr.MerkleTreeRoots(big.NewInt(1))
		assert.NilError(t, err)

		assert.Equal(t, leafData.Enumerated.Cmp(roots2.Enumerated), 0)
		assert.Equal(t, leafData.EnumeratedTotal.Cmp(roots2.EnumeratedTotal), 0)
		assert.Equal(t, leafData.Ownership.Cmp(roots2.Ownership), 0)
	})

	t.Run(`tag root before setting merkle roots. checkout at block which tag does not exist returns error`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := account.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.Checkout(1)
		assert.Error(t, err, "no tag found for this block number 1")
	})
}

func TestDeleteRootTag(t *testing.T) {
	t.Parallel()
	t.Run(`Tag two roots and then delete the first tag. Checkout at deleted tag returns error`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := account.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		err = tr.TagRoot(2)
		assert.NilError(t, err)
		err = tx.Commit()
		assert.NilError(t, err)

		tx = service.NewTransaction()
		err = account.DeleteRootTag(tx, common.HexToAddress("0x500").Hex(), 1)
		assert.NilError(t, err)
		err = tx.Commit()
		assert.NilError(t, err)

		tx = service.NewTransaction()
		tr, err = account.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)
		err = tr.Checkout(1)
		assert.Error(t, err, "no tag found for this block number 1")
	})
}

func TestGetLastTag(t *testing.T) {
	t.Parallel()
	t.Run(`Tag two roots and then get the latest tagged block`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := account.NewTree(common.HexToAddress("0x500"), tx)
		assert.NilError(t, err)

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		err = tr.TagRoot(2)
		assert.NilError(t, err)
		err = tx.Commit()
		assert.NilError(t, err)

		block, err := tr.GetLastTaggedBlock()
		assert.NilError(t, err)
		assert.Equal(t, block, int64(2))
	})
}
