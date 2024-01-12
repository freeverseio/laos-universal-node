package storage

type Tx interface {
	Commit() error
	Discard()
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	GetKeysWithPrefix(prefix []byte, reverse ...bool) [][]byte
	Len() int
	ClearAll() error
}

type Service interface {
	NewTransaction() Tx
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
	GetKeysWithPrefix(prefix []byte, reverse ...bool) ([][]byte, error)
}
