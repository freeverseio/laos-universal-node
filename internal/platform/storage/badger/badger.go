package badger

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

type Badger struct {
	db *badger.DB
}

func NewService(db *badger.DB) storage.Service {
	return Badger{
		db: db,
	}
}

func (b Badger) Set(key, value []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		if err != nil {
			return err
		}
		return nil
	})
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

type Tx struct {
	tx *badger.Txn
}

func (b Badger) NewTransaction() storage.Tx {
	return Tx{
		b.db.NewTransaction(true),
	}
}

func (t Tx) Commit() error {
	return t.tx.Commit()
}

func (t Tx) Discard() {
	t.tx.Discard()
}

func (t Tx) Set(key, value []byte) error {
	return t.tx.Set(key, value)
}

func (t Tx) Get(key []byte) ([]byte, error) {
	item, err := t.tx.Get(key)
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}
