package storage

type Tx interface {
	Commit() error
	Discard()
	Set(key []byte, value []byte) error
}

type Storage interface {
	NewTransaction() Tx
	Get(key []byte) ([]byte, error)
}
