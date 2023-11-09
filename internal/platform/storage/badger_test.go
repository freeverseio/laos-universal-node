package storage_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

var (
	expectedKey, expectedVal = []byte("expectedKey"), []byte("expectedValue")
	db                       *badger.DB // badger.DB is thread-safe
)

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		log.Fatalf("error occurred setting up storage tests: %s", err.Error())
	}
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestStorageGetData(t *testing.T) {
	t.Parallel()

	service := storage.New(db)

	performTransaction(t, service)

	got, err := service.Get(expectedKey)
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}
	if !bytes.Equal(got, expectedVal) {
		t.Fatalf("value retrieved from DB is %v, expected %v", string(got), string(expectedVal))
	}
}

func TestStorageGetNoData(t *testing.T) {
	t.Parallel()

	service := storage.New(db)

	got, err := service.Get([]byte("idonotexist"))
	if err == nil {
		t.Fatal("got no error, expecting badger.ErrKeyNotFound")
	}
	if err != badger.ErrKeyNotFound {
		t.Fatalf("got error %s, expecting badger.ErrKeyNotFound", err.Error())
	}
	if got != nil {
		t.Fatalf("value retrieved from DB is %v, expected nil", string(got))
	}
}

func performTransaction(t *testing.T, service storage.Storage) {
	t.Helper()
	tx := service.NewTransaction()
	defer tx.Discard()
	err := tx.Set(expectedKey, expectedVal)
	if err != nil {
		t.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}
}

func setup() error {
	var err error
	db, err = badger.Open(
		badger.DefaultOptions("").
			WithInMemory(true).
			WithLoggingLevel(badger.ERROR))
	if err != nil {
		return err
	}
	return nil
}

func teardown() {
	if err := db.Close(); err != nil {
		log.Printf("error occurred closing the DB: %s", err.Error())
	}
}
