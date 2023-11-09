package storage_test

import (
	"bytes"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

var expectedKey, expectedVal = []byte("expectedKey"), []byte("expectedValue")

func TestStorageGet(t *testing.T) {
	t.Parallel()

	db := createDB(t)
	service := storage.New(db)

	performTransaction(t, service)

	got, err := service.Get(expectedKey)
	if err != nil {
		t.Fatalf("got error %s while expecting no error", err.Error())
	}
	if !bytes.Equal(got, expectedVal) {
		t.Fatalf("value retrieved from DB is %v, expected %v", string(got), string(expectedVal))
	}
}

func createDB(t *testing.T) *badger.DB {
	t.Helper()
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Fatal(err)
	}
	return db
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
