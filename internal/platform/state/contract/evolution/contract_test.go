package evolution_test

import (
	"math/big"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
)


func TestStoreMintedWithExternalURIEvents(t *testing.T) {
	t.Parallel()
	t.Run("stores minted events", func(t *testing.T) {
		t.Parallel()
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