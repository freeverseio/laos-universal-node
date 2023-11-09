package storage

import (
	"github.com/dgraph-io/badger/v4"
)

type Badger struct {
	db *badger.DB
}

func New(db *badger.DB) Storage {
	return Badger{
		db: db,
	}
}

func (b Badger) Get(key []byte) ([]byte, error) {
	var returnValue []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		returnValue, err = item.ValueCopy(returnValue)
		if err != nil {
			return err
		}
		return nil
	})
	return returnValue, err
}

func (b Badger) GetKeysWithPrefix(prefix []byte) ([][]byte, error) {
	var keys [][]byte
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		iterator := txn.NewIterator(opts)
		defer iterator.Close()
		for iterator.Seek(prefix); iterator.ValidForPrefix(prefix); iterator.Next() {
			item := iterator.Item()
			var key []byte
			key = item.KeyCopy(key)
			keys = append(keys, key)
		}
		return nil
	})
	return keys, err
}

func (b Badger) NewTransaction() Tx {
	return b.db.NewTransaction(true)
}
