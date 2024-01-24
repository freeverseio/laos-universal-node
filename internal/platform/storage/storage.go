package storage

type Tx interface {
	Commit() error
	Discard()
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	GetKeysWithPrefix(prefix []byte, reverse ...bool) [][]byte
	FilterKeysWithPrefix(prefix []byte, from, to string) [][]byte
	GetValuesWithPrefix(prefix []byte, reverse ...bool) [][]byte
}

type Service interface {
	NewTransaction() Tx
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
	GetKeysWithPrefix(prefix []byte, reverse ...bool) ([][]byte, error)
}
