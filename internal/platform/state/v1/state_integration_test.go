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
		tx, err := createBadgerTransaction(t, db)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		contract := common.HexToAddress("0x500")
		collection := common.HexToAddress("0x501")
		c := model.ERC721UniversalContract{
			Address:           contract,
			CollectionAddress: collection,
			BlockNumber:       100,
		}

		err = tx.StoreERC721UniversalContracts([]model.ERC721UniversalContract{c})
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		contracts := tx.GetAllERC721UniversalContracts()
		fmt.Println("c", contracts)

		err = tx.LoadContractTrees(contract)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		err = tx.TagRoot(100)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		// we can checkout the contract at block 100
		err = tx.Checkout(100)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		err = tx.Commit()
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		tx, err = createBadgerTransaction(t, db)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		err = tx.LoadContractTrees(contract)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		err = tx.DeleteOrphanRootTags(int64(100), int64(105))
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		// we can not checkout the contract at block 100 anymore
		err = tx.Checkout(100)
		if err == nil {
			t.Errorf(`got no error when an error was expected`)
		}
	})
}

func TestLoadMerkleTreesWithBadger(t *testing.T) {
	t.Run("fails when contract is 0x0", func(t *testing.T) {
		db := createBadger(t)
		tx, err := createBadgerTransaction(t, db)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		expectedErr := fmt.Sprintf("contract address is " + common.Address{}.String())

		err = tx.LoadContractTrees(common.HexToAddress("0x0"))
		if err == nil {
			t.Errorf("got no error while an error was expected")
		}
		if err != nil && err.Error() != expectedErr {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr)
		}
	})
	t.Run("mint many assets in different transaction and then measure tx size", func(t *testing.T) {
		// the test can run for 1400 assets and it fails for 1500 for this table size
		// for the table size of 1 << 30 the test can successfully run for 18000 assets
		// but that test is slower so I am not putting it for such big number
		contract := common.HexToAddress("0x500")
		db := createBadger(t)
		blocks := 1
		mintsInBlock := 1000
		for block := 0; block < blocks; block++ {
			tx, err := createBadgerTransaction(t, db)
			if err != nil {
				t.Errorf(`got error "%v" when no error was expected`, err)
			}
			err = tx.LoadContractTrees(contract)
			if err != nil {
				t.Errorf(`got error "%v" when no error was expected`, err)
			}

			for mintInBlock := 0; mintInBlock < mintsInBlock; mintInBlock++ {
				mintId := block*mintsInBlock + mintInBlock

				owner := common.HexToAddress("0xB200110583D9d9F5E041FcEe024886bd00996691")
				tokenId := big.NewInt(int64(mintId))
				tokenId = tokenId.Lsh(tokenId, 160)
				tokenId = tokenId.Add(tokenId, owner.Big())

				mintEvent := model.MintedWithExternalURI{
					Slot:        big.NewInt(int64(mintId)),
					To:          common.HexToAddress("0x3"),
					TokenURI:    "tokenURI",
					TokenId:     tokenId,
					BlockNumber: uint64(block),
					Timestamp:   uint64(time.Now().Unix()),
				}
				err = tx.Mint(contract, &mintEvent)
				if err != nil {
					t.Fatal("got error when no error was expected 1", err.Error())
				}
			}
			err = tx.Commit()
			if err != nil {
				t.Fatal("got error when no error was expected 2", err.Error())
			}
		}

		tx, err := createBadgerTransaction(t, db)
		if err != nil {
			t.Fatal("got error when no error was expected 3", err.Error())
		}
		err = tx.LoadContractTrees(contract)
		if err != nil {
			t.Fatal("got error when no error was expected 3", err.Error())
		}

		err = tx.Commit()
		if err != nil {
			t.Fatal("got error when no error was expected 4", err.Error())
		}
		errDrop := db.DropAll()
		if errDrop != nil {
			t.Fatal("got error when no error was expected 5", err.Error())
		}
		errClose := db.Close()
		if errClose != nil {
			t.Fatal("got error when no error was expected 6", err.Error())
		}
	})
}

func TestStoreAngGetMintedWithExternalURIEvents(t *testing.T) {
	t.Run("stores and gets mint events", func(t *testing.T) {
		db := createBadger(t)
		tx, err := createBadgerTransaction(t, db)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		err = tx.StoreMintedWithExternalURIEvents(common.HexToAddress("0x500").Hex(), &model.MintedWithExternalURI{
			Slot:        big.NewInt(1),
			To:          common.HexToAddress("0x3"),
			TokenURI:    "tokenURI",
			TokenId:     big.NewInt(1),
			BlockNumber: 100,
			Timestamp:   1000,
			TxIndex:     1,
		})
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		events, err := tx.GetMintedWithExternalURIEvents(common.HexToAddress("0x500").Hex(), 100)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		if len(events) != 1 {
			t.Errorf(`got %d events when 1 was expected`, len(events))
		}
		if events[0].Slot.Cmp(big.NewInt(1)) != 0 {
			t.Errorf(`got slot %d when 1 was expected`, events[0].Slot)
		}
		if errDrop := db.DropAll(); errDrop != nil {
			t.Errorf(`got error "%v" when no error was expected`, errDrop)
		}
		if errClose := db.Close(); errClose != nil {
			t.Errorf(`got error "%v" when no error was expected`, errClose)
		}
	})
}

func createBadger(t *testing.T) *badger.DB {
	t.Helper()
	db, err := badger.Open(
		badger.DefaultOptions("").
			WithInMemory(true).
			WithLoggingLevel(badger.ERROR).WithMemTableSize(1 << 30))
	if err != nil {
		t.Fatalf("error initializing storage: %v", err)
	}

	return db
}

func createBadgerTransaction(t *testing.T, db *badger.DB) (state.Tx, error) {
	t.Helper()
	badgerService := badgerStorage.NewService(db)
	stateService := v1.NewStateService(badgerService)
	return stateService.NewTransaction()
}
