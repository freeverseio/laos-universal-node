package storage

import "github.com/dgraph-io/badger/v4"

type Badger struct {
	db *badger.DB
}

type Tx interface {
	Commit() error
	Discard()
	Set(key []byte, value []byte) error
}

func NewStorage(db *badger.DB) Storage {
	return Badger{
		db: db,
	}
}

func (b *Badger) Get(key []byte) ([]byte, error) {
	var returnValue []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		_, err = item.ValueCopy(returnValue)
		if err != nil {
			return err
		}
		return nil
	})
	return returnValue, err
}

func (b *Badger) NewTransaction() Tx {
	return b.db.NewTransaction(true)
}
