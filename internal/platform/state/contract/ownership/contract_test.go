package ownership_test

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
)

func TestStoreGetERC721UniversalContracts(t *testing.T) {
	t.Parallel()
	t.Run("stores minted events", func(t *testing.T) {
		t.Parallel()
		db := createBadger(t)
		tx, err := createBadgerTransaction(t, db)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
        uEvent := model.ERC721UniversalContract{
            Address: common.HexToAddress("0x500"),
            CollectionAddress:  common.HexToAddress("0x501"),
            BlockNumber: 10,
        }

		err = tx.StoreERC721UniversalContracts([]model.ERC721UniversalContract{uEvent})
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		contracts := tx.GetAllERC721UniversalContracts()
		
		if len(contracts) != 1 {
			t.Errorf(`got %d contracts when 1 was expected`, len(contracts))
		}
	})
}

func createBadgerTransaction(t *testing.T, db *badger.DB) (state.Tx, error) {
	t.Helper()
	badgerService := badgerStorage.NewService(db)
	stateService := v1.NewStateService(badgerService)
	return stateService.NewTransaction()
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
