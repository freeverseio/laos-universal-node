package account_test

import (
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
		_, err := account.NewTree(nil)
		assert.Error(t, err, "store is nil")
	})

	t.Run(`initial root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := account.NewTree(tx)
		assert.NilError(t, err)
		assert.Equal(t, tr.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	})

	t.Run(`set  contract data changes the root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr1, err := account.NewTree(tx)
		assert.NilError(t, err)

		testData := account.AccountData{
			EnumeratedRoot:      common.HexToHash("0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55"),
			EnumeratedTotalRoot: common.HexToHash("0xcf46e158742177f61d06cb049d82c7d4aeb7420205d0e1c1bacc45406acde8f3"),
			OwnershipRoot:       common.HexToHash("0x59d7de0b77f377095267336d574c03d9c444c5cbbdfe03997e16aa1ff0df6798"),
			TotalSupply:         100,
		}

		err = tr1.SetAccountData(&testData, common.HexToAddress("0x500"))
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0x76d53bd50ccd68af8b623d6d3eccc311c5a7d92db6a06993bde49da9b1c82e9a")

		data, err := tr1.AccountData(common.HexToAddress("0x500"))
		assert.NilError(t, err)
		assert.Equal(t, data.EnumeratedRoot.Cmp(testData.EnumeratedRoot), 0)
		assert.Equal(t, data.EnumeratedTotalRoot.Cmp(testData.EnumeratedTotalRoot), 0)
		assert.Equal(t, data.OwnershipRoot.Cmp(testData.OwnershipRoot), 0)
		assert.Equal(t, data.TotalSupply, testData.TotalSupply)
	})
}

func TestTag(t *testing.T) {
	t.Parallel()
	t.Run(`tag root before setting merkle roots. checkout at that root returns state before setting merkle roots`, func(t *testing.T) {
		t.Parallel()
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr1, err := account.NewTree(tx)
		assert.NilError(t, err)

		testData := account.AccountData{
			EnumeratedRoot:      common.HexToHash("0x390efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55"),
			EnumeratedTotalRoot: common.HexToHash("0xcf46e158742177f61d06cb049d82c7d4aeb7420205d0e1c1bacc45406acde8f3"),
			OwnershipRoot:       common.HexToHash("0x59d7de0b77f377095267336d574c03d9c444c5cbbdfe03997e16aa1ff0df6798"),
			TotalSupply:         100,
		}

		err = tr1.SetAccountData(&testData, common.HexToAddress("0x500"))
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0xd2dd7011bbd026ffc2340f1e6eb5e7d6a8e7756e5985eb63d860bbbf93762a92")

		data, err := tr1.AccountData(common.HexToAddress("0x500"))
		assert.NilError(t, err)
		assert.Equal(t, data.EnumeratedRoot.Cmp(testData.EnumeratedRoot), 0)
		assert.Equal(t, data.EnumeratedTotalRoot.Cmp(testData.EnumeratedTotalRoot), 0)
		assert.Equal(t, data.OwnershipRoot.Cmp(testData.OwnershipRoot), 0)
		assert.Equal(t, data.TotalSupply, testData.TotalSupply)

		err = tr1.TagRoot(1)
		assert.NilError(t, err)

		testData2 := account.AccountData{
			EnumeratedRoot:      common.HexToHash("0x090efb1b494cf9fec34922b9e6c80adfaeb1a488e7abc52d40d034adb6527c55"),
			EnumeratedTotalRoot: common.HexToHash("0x0f46e158742177f61d06cb049d82c7d4aeb7420205d0e1c1bacc45406acde8f3"),
			OwnershipRoot:       common.HexToHash("0x09d7de0b77f377095267336d574c03d9c444c5cbbdfe03997e16aa1ff0df6798"),
			TotalSupply:         100,
		}

		err = tr1.SetAccountData(&testData, common.HexToAddress("0x500"))
		assert.NilError(t, err)
		assert.Equal(t, tr1.Root().String(), "0xd2dd7011bbd026ffc2340f1e6eb5e7d6a8e7756e5985eb63d860bbbf93762a92")

		data2, err := tr1.AccountData(common.HexToAddress("0x500"))
		assert.NilError(t, err)
		assert.Equal(t, data2.EnumeratedRoot.Cmp(testData2.EnumeratedRoot), 0)
		assert.Equal(t, data2.EnumeratedTotalRoot.Cmp(testData2.EnumeratedTotalRoot), 0)
		assert.Equal(t, data2.OwnershipRoot.Cmp(testData2.OwnershipRoot), 0)
		assert.Equal(t, data2.TotalSupply, testData2.TotalSupply)

		err = tr1.Checkout(1)
		assert.NilError(t, err)

		data3, err := tr1.AccountData(common.HexToAddress("0x500"))
		assert.NilError(t, err)
		assert.Equal(t, data3.EnumeratedRoot.Cmp(testData.EnumeratedRoot), 0)
		assert.Equal(t, data3.EnumeratedTotalRoot.Cmp(testData.EnumeratedTotalRoot), 0)
		assert.Equal(t, data3.OwnershipRoot.Cmp(testData.OwnershipRoot), 0)
		assert.Equal(t, data3.TotalSupply, testData.TotalSupply)
	})

	t.Run(`checkout at block which tag does not exist returns error`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tr, err := account.NewTree(tx)
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

		tr, err := account.NewTree(tx)
		assert.NilError(t, err)

		err = tr.TagRoot(1)
		assert.NilError(t, err)

		err = tr.TagRoot(2)
		assert.NilError(t, err)

		err = tr.DeleteRootTag(1)
		assert.NilError(t, err)
		err = tx.Commit()
		assert.NilError(t, err)

		tx = service.NewTransaction()
		tr, err = account.NewTree(tx)
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

		tr, err := account.NewTree(tx)
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
