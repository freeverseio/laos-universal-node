package v1_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
)

func TestLoadMerkleTrees(t *testing.T) {
	t.Parallel()
	t.Run("fails when contract is 0x0", func(t *testing.T) {
		t.Parallel()
		tx := createTransaction()
		expectedErr := fmt.Sprintf("contract address is " + common.Address{}.String())

		err := tx.LoadMerkleTrees(common.HexToAddress("0x0"))
		if err == nil {
			t.Errorf("got no error while an error was expected")
		}
		if err != nil && err.Error() != expectedErr {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr)
		}
	})
	t.Run("successfully loads merkle trees in memory", func(t *testing.T) {
		t.Parallel()
		tx := createTransaction()
		contract := common.HexToAddress("0x500")

		err := tx.LoadMerkleTrees(contract)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		// if tx.Mint works, we are sure that the merkle trees are correctly loaded in memory
		// as Mint accesses all merkle trees
		mintEvent := model.MintedWithExternalURI{
			Slot:        big.NewInt(1),
			To:          common.HexToAddress("0x3"),
			TokenURI:    "tokenURI",
			TokenId:     big.NewInt(1),
			BlockNumber: 100,
			Timestamp:   1000,
		}
		err = tx.Mint(contract, &mintEvent)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
	})
}

func TestDeleteOrphanRootTags(t *testing.T) {
	// Do not run this test in parallel as it uses bager in memory
	t.Run("successfully deletes orphan root tags", func(t *testing.T) {
		tx := createBadgerTransaction(t)
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

func createTransaction() state.Tx {
	mem := memory.New()
	stateService := v1.NewStateService(mem)
	return stateService.NewTransaction()
}

func createBadgerTransaction(t *testing.T) state.Tx {
	db, err := badger.Open(
		badger.DefaultOptions("").
			WithInMemory(true).
			WithLoggingLevel(badger.ERROR))
	if err != nil {
		t.Fatalf("error initializing storage: %v", err)
	}
	badgerService := badgerStorage.NewService(db)
	stateService := v1.NewStateService(badgerService)
	return stateService.NewTransaction()
}
