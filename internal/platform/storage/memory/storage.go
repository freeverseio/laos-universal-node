package memory

import (
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

type data struct {
	data map[string][]byte
}

type service struct {
	committed data
}

// New creates a new in-memory storage service.
func New() *service {
	s := service{}
	s.committed.data = make(map[string][]byte)
	return &s
}

// NewTransaction creates a new storage transaction.
func (b *service) NewTransaction() storage.Tx {
	temp := b.committed
	temp.data = make(map[string][]byte)
	for k, v := range b.committed.data {
		temp.data[k] = v
	}
	return &tx{temp, b}
}

func (b *service) Get(key []byte) ([]byte, error) {
	value, has := b.committed.data[string(key)]
	if !has {
		return nil, nil
	}

	return value, nil
}

// GetKeysWithPrefix returns nil. it is just put here for interface compatibility
func (b *service) GetKeysWithPrefix(prefix []byte) ([][]byte, error) {
	return nil, nil
}

// Set updates a key/value pair in the storage service.
func (b *service) Set(key, value []byte) error {
	b.committed.data[string(key)] = value
	return nil
}

// Delete deletes a key
func (b *service) Delete(key []byte) error {
	delete(b.committed.data, string(key))
	return nil
}

// Close closes the storage service.
func (b service) Close() {}

type tx struct {
	temp data
	s    *service
}

// Set updates a key/value pair in the storage service.
func (b tx) Set(key, value []byte) error {
	b.temp.data[string(key)] = value
	return nil
}

// Get stores a key/value pair in the storage service.
func (b tx) Get(key []byte) ([]byte, error) {
	value, has := b.temp.data[string(key)]
	if !has {
		return nil, nil
	}

	return value, nil
}

// GetKeysWithPrefix added for interface compatibility
func (b tx) GetKeysWithPrefix(prefix []byte) [][]byte {
	// TODO implement this if we have to use it for testing purposes
	return nil
}

// Delete deletes a key.
func (b tx) Delete(key []byte) error {
	delete(b.temp.data, string(key))
	return nil
}

// Discard implemented just to make it compatible with storage interface
func (b *tx) Discard() {
}

// Commit commits the storage transaction.
func (b *tx) Commit() error {
	b.s.committed = b.temp
	return nil
}
