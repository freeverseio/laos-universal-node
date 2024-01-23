package badger

import (
	"bytes"

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

// GetKeysWithPrefix fetches keys with a given prefix.
// The reverse parameter is optional and defaults to false (non-reverse iterator).
// Pass true as the second argument to enable reverse iteration.
func (b Badger) GetKeysWithPrefix(prefix []byte, reverse ...bool) ([][]byte, error) {
	var keys [][]byte

	// Determine the reverse setting based on the optional parameter
	isReverse := false
	if len(reverse) > 0 {
		isReverse = reverse[0]
	}

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = isReverse
		iterator := txn.NewIterator(opts)
		defer iterator.Close()

		// Append 0xff to the prefix for the seek key if in reverse mode
		var seekPrefix []byte
		if isReverse {
			seekPrefix = append([]byte(nil), prefix...)
			seekPrefix = append(seekPrefix, 0xff)
		} else {
			seekPrefix = prefix
		}

		for iterator.Seek(seekPrefix); iterator.ValidForPrefix(prefix); iterator.Next() {
			item := iterator.Item()
			key := item.KeyCopy(nil)
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

// Set sets []byte value for []byte key
func (t Tx) Set(key, value []byte) error {
	return t.tx.Set(key, value)
}

// Get returns byte for the key
func (t Tx) Get(key []byte) ([]byte, error) {
	// TODO to use t.Discard here we must first give the possibility to have t as read-only (i.e. `NewTransaction(readOnly bool)`)
	// so as a first thing, `Get` checks if t is read only, and, if it is, `defer t.Discard()`
	item, err := t.tx.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	return item.ValueCopy(nil)
}

// Delete deletes a key.
func (t Tx) Delete(key []byte) error {
	return t.tx.Delete(key)
}

// GetKeysWithPrefix looks for all the keys with the specified prefix and returns them. It doesn't return values
func (t Tx) GetKeysWithPrefix(prefix []byte, reverse ...bool) [][]byte {
	var keys [][]byte

	// Determine the reverse setting based on the optional parameter
	isReverse := false
	if len(reverse) > 0 {
		isReverse = reverse[0]
	}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Reverse = isReverse
	opts.PrefetchSize = 100
	iterator := t.tx.NewIterator(opts)
	defer iterator.Close()

	// Append 0xff to the prefix for the seek key if in reverse mode
	var seekPrefix []byte
	if isReverse {
		seekPrefix = append([]byte(nil), prefix...)
		seekPrefix = append(seekPrefix, 0xff)
	} else {
		seekPrefix = prefix
	}

	for iterator.Seek(seekPrefix); iterator.ValidForPrefix(prefix); iterator.Next() {
		item := iterator.Item()
		key := item.KeyCopy(nil)
		keys = append(keys, key)
	}

	return keys
}

func (t Tx) FilterKeysWithPrefix(prefix []byte, from, to string) [][]byte {
	var keys [][]byte

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.PrefetchSize = 100
	iterator := t.tx.NewIterator(opts)
	defer iterator.Close()

	startKey := append([]byte(nil), prefix...)
	startKey = append(startKey, []byte(from)...)
	endKey := append([]byte(nil), prefix...)
	endKey = append(endKey, []byte(to)...)

	for iterator.Seek(startKey); iterator.ValidForPrefix(prefix); iterator.Next() {
		item := iterator.Item()
		key := item.KeyCopy(nil)
		if bytes.Compare(key, []byte(endKey)) > 0 {
			break
		}
		keys = append(keys, key)
	}

	return keys
}
