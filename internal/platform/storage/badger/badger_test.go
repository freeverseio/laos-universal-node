package badger_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
)

var db *badger.DB // badger.DB is thread-safe

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

	service := badgerStorage.NewService(db)

	expectedKey, expectedVal := []byte("expectedKey"), []byte("expectedValue")
	performTransaction(t, expectedKey, expectedVal, service)

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

	service := badgerStorage.NewService(db)

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

func TestStorageGetKeysWithPrefix(t *testing.T) {
	t.Parallel()

	service := badgerStorage.NewService(db)
	expectedKeys := [][]byte{
		[]byte("entry_1"),
		[]byte("entry_2"),
	}
	performTransaction(t, expectedKeys[0], []byte(""), service)
	performTransaction(t, expectedKeys[1], []byte(""), service)
	performTransaction(t, []byte("third"), []byte(""), service)

	got, err := service.GetKeysWithPrefix([]byte("entry_"))
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}
	if len(got) == 0 {
		t.Fatal("got 0 keys when 3 keys were expected")
	}
	if !bytes.Equal(got[0], expectedKeys[0]) {
		t.Fatalf("got key %v, expected %v", string(got[0]), string(expectedKeys[0]))
	}
	if !bytes.Equal(got[1], expectedKeys[1]) {
		t.Fatalf("got key %v, expected %v", string(got[1]), string(expectedKeys[1]))
	}
}

func TestSetGetDeleteOnTransaction(t *testing.T) {
	t.Parallel()
	t.Helper()

	service := badgerStorage.NewService(db)
	t.Helper()
	tx := service.NewTransaction()
	defer tx.Discard()

	err := tx.Set([]byte("key"), []byte("value"))
	if err != nil {
		t.Fatal(err)
	}

	buf, err := tx.Get([]byte("key"))
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}
	if !bytes.Equal(buf, []byte("value")) {
		t.Fatalf("got %v, expected %v", string(buf), "value")
	}

	err = tx.Delete([]byte("key"))
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}

	buf, err = tx.Get([]byte("key"))
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}
	if !bytes.Equal(buf, []byte("")) {
		t.Fatalf("got %v, expected %v", string(buf), "value")
	}
}

func TestStorageGetNoKeysWithPrefix(t *testing.T) {
	t.Parallel()

	service := badgerStorage.NewService(db)

	got, err := service.GetKeysWithPrefix([]byte("idonotexisteither_"))
	if err != nil {
		t.Fatalf("got error %s, expecting no error", err.Error())
	}
	if len(got) > 0 {
		t.Fatalf("got %d keys when 0 keys were expected", len(got))
	}
}

func performTransaction(t *testing.T, key, val []byte, service storage.Service) {
	t.Helper()
	tx := service.NewTransaction()
	defer tx.Discard()
	err := tx.Set(key, val)
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
