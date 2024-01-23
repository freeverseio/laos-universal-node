package v1_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
)

func TestDeleteOrphanRootTags(t *testing.T) {
	t.Run("successfully deletes orphan root tags", func(t *testing.T) {
		db := createBadger(t)
		tx := createBadgerTransaction(t, db)
		contract := common.HexToAddress("0x500")
		collection := common.HexToAddress("0x501")
		c := model.ERC721UniversalContract{
			Address:           contract,
			CollectionAddress: collection,
			BlockNumber:       100,
		}

		if err := tx.StoreERC721UniversalContracts([]model.ERC721UniversalContract{c}); err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		contracts := tx.GetAllERC721UniversalContracts()
		fmt.Println("c", contracts)

		err := tx.LoadMerkleTrees(contract)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		err = tx.TagRoot(contract, 100)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		// we can checkout the contract at block 100
		err = tx.Checkout(contract, 100)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		err = tx.Commit()
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		tx = createBadgerTransaction(t, db)
		err = tx.LoadMerkleTrees(contract)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		err = tx.DeleteOrphanRootTags(int64(100), int64(105))
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		// we can not checkout the contract at block 100 anymore
		err = tx.Checkout(contract, 100)
		if err == nil {
			t.Errorf(`got no error when an error was expected`)
		}
	})
}

func TestLoadMerkleTreesWithBadger(t *testing.T) {
	t.Run("fails when contract is 0x0", func(t *testing.T) {
		db := createBadger(t)
		tx := createBadgerTransaction(t, db)
		expectedErr := fmt.Sprintf("contract address is " + common.Address{}.String())

		err := tx.LoadMerkleTrees(common.HexToAddress("0x0"))
		if err == nil {
			t.Errorf("got no error while an error was expected")
		}
		if err != nil && err.Error() != expectedErr {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr)
		}
	})
	t.Run("successfully loads merkle trees in memory and mints 1400 assets", func(t *testing.T) {
		// the test can run for 1400 assets and it fails for 1500 for this table size
		// for the table size of 1 << 30 the test can successfully run for 18000 assets
		// but that test is slower so I am not putting it for such big number
		db := createBadger(t)
		tx := createBadgerTransaction(t, db)
		contract := common.HexToAddress("0x500")

		err := tx.LoadMerkleTrees(contract)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		for i := 1; i < 1400; i++ {
			mintEvent := model.MintedWithExternalURI{
				Slot:        big.NewInt(int64(i)),
				To:          common.HexToAddress("0x3"),
				TokenURI:    "tokenURI",
				TokenId:     big.NewInt(int64(i)),
				BlockNumber: uint64(i),
				Timestamp:   uint64(time.Now().Unix()),
			}
			err = tx.Mint(contract, &mintEvent)
			if err != nil {
				t.Errorf(`got error "%v" when no error was expected`, err)
			}
		}

		err = tx.Commit()
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		errDrop := db.DropAll()
		if errDrop != nil {
			t.Errorf(`got error "%v" when no error was expected`, errDrop)
		}
		errClose := db.Close()
		if errClose != nil {
			t.Errorf(`got error "%v" when no error was expected`, errClose)
		}
	})
}

func createBadger(t *testing.T) *badger.DB {
	t.Helper()
	db, err := badger.Open(
		badger.DefaultOptions("").
			WithInMemory(true).
			WithLoggingLevel(badger.ERROR))
	if err != nil {
		t.Fatalf("error initializing storage: %v", err)
	}

	return db
}

func createBadgerTransaction(t *testing.T, db *badger.DB) state.Tx {
	t.Helper()
	badgerService := badgerStorage.NewService(db)
	stateService := v1.NewStateService(badgerService)
	return stateService.NewTransaction()
}
