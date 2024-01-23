package badger_test

import (
	"bytes"
	"log"
	"os"
	"strconv"
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

func TestStorageGetNoKeysWithPrefixWithValues(t *testing.T) {
	// DON'T PARALLELIZE THIS TEST
	service := badgerStorage.NewService(db)
	keysAndValues := [][]byte{
		[]byte("prefix_001"), []byte("1"),
		[]byte("prefix_002"), []byte("2"),
		[]byte("prefix_003"), []byte("3"),
		[]byte("prefix_004"), []byte("4"),
		[]byte("prefix_005"), []byte("5"),
		[]byte("prefix_010"), []byte("10"),
		[]byte("prefix_123654487745555421516551151615165151155165515165"), []byte("111"),
	}

	for i := 0; i < len(keysAndValues); i += 2 {
		if err := service.Set(keysAndValues[i], keysAndValues[i+1]); err != nil {
			t.Fatalf("got error %s, expecting no error", err.Error())
		}
	}
	tx := service.NewTransaction()
	got := tx.GetKeysWithPrefix([]byte("prefix_"), false)

	if len(got) != 7 {
		t.Fatalf("got %d keys when 7 keys were expected", len(got))
	}
	if string(got[0]) != "prefix_001" {
		t.Fatalf("got %v, expected %v", string(got[0]), "prefix_001")
	}
	if string(got[1]) != "prefix_002" {
		t.Fatalf("got %v, expected %v", string(got[1]), "prefix_002")
	}
	if err := db.DropAll(); err != nil {
		log.Printf("error occurred closing the DB: %s", err.Error())
	}
}

func TestStorageGetNoKeysWithPrefixReverse(t *testing.T) {
	// DON'T PARALLELIZE THIS TEST
	service := badgerStorage.NewService(db)

	for i := 0; i < 1000; i++ {
		err := service.Set([]byte("prefix_"+strconv.Itoa(i)), []byte(strconv.Itoa(i)))
		if err != nil {
			t.Fatalf("got error %s, expecting no error", err.Error())
		}
	}
	tx := service.NewTransaction()
	got := tx.GetKeysWithPrefix([]byte("prefix_"), true)
	if len(got) != 1000 {
		t.Fatalf("got %d keys when 1000 keys were expected", len(got))
	}
	if string(got[0]) != "prefix_999" {
		t.Fatalf("got %v, expected %v", string(got[0]), "prefix_9999")
	}
}

func TestFilterKeysWithPrefix(t *testing.T) {
	// DON'T PARALLELIZE THIS TEST
	service := badgerStorage.NewService(db)
	keysAndValues := [][]byte{
		[]byte("prefix_001"), []byte("1"),
		[]byte("prefix_002"), []byte("2"),
		[]byte("prefix_003"), []byte("3"),
		[]byte("prefix_004"), []byte("4"),
		[]byte("prefix_005"), []byte("5"),
		[]byte("prefix_010"), []byte("10"),
	}

	for i := 0; i < len(keysAndValues); i += 2 {
		if err := service.Set(keysAndValues[i], keysAndValues[i+1]); err != nil {
			t.Fatalf("got error %s, expecting no error", err.Error())
		}
	}
	tx := service.NewTransaction()
	got := tx.FilterKeysWithPrefix([]byte("prefix_"), "002", "005")

	if len(got) != 4 {
		t.Fatalf("got %d keys when 7 keys were expected", len(got))
	}

	if err := db.DropAll(); err != nil {
		log.Printf("error occurred closing the DB: %s", err.Error())
	}
}

func TestFilterKeysWithPrefixWithValues(t *testing.T) {
	// DON'T PARALLELIZE THIS TEST
	service := badgerStorage.NewService(db)
	for i := 0; i < 100; i++ {
		formatedBlockNumber := formatBlockNumber(uint64(i), 15)
		err := service.Set([]byte("prefix_"+formatedBlockNumber), []byte(strconv.Itoa(i)))
		if err != nil {
			t.Fatalf("got error %s, expecting no error", err.Error())
		}
	}

	tx := service.NewTransaction()
	got := tx.FilterKeysWithPrefix([]byte("prefix_"), formatBlockNumber(uint64(6), 15), formatBlockNumber(uint64(10), 15))

	if len(got) != 5 {
		t.Fatalf("got %d keys when 7 keys were expected", len(got))
	}

	if err := db.DropAll(); err != nil {
		log.Printf("error occurred closing the DB: %s", err.Error())
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
	if err := db.DropAll(); err != nil {
		log.Printf("error occurred closing the DB: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		log.Printf("error occurred closing the DB: %s", err.Error())
	}
}

func formatBlockNumber(blockNumber uint64, blockNumberDigits uint16) string {
	// Convert the block number to a string
	blockNumberString := strconv.FormatUint(blockNumber, 10)
	// Pad with leading zeros if shorter
	for len(blockNumberString) < int(blockNumberDigits) {
		blockNumberString = "0" + blockNumberString
	}
	return blockNumberString
}
