package sparsemt

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

// Null represents a 0x...0 hash
const Null = "0x0000000000000000000000000000000000000000000000000000000000000000"

// Node contains the left and right hashes of a node
type Node struct {
	L common.Hash
	R common.Hash
}

// Hash returns the hash of the node
func (b Node) Hash() common.Hash {
	if b.L.String() == Null && b.R.String() == Null {
		return common.Hash{}
	}
	return crypto.Keccak256Hash(append(b.L.Bytes(), b.R.Bytes()...))
}

// PutNode adds node in the merkletree storage
func PutNode(store storage.Tx, node Node, prefix string) (common.Hash, error) {
	hash := node.Hash()
	if hash.String() == Null {
		return common.HexToHash(Null), nil
	}

	var value []byte
	value = append(value, node.L.Bytes()...)
	value = append(value, node.R.Bytes()...)
	return hash, store.Set([]byte(prefix+hash.String()), value)
}

// GetNode returns node from the storage given the hash
func GetNode(store storage.Tx, hash common.Hash, prefix string) (*Node, error) {
	if hash.String() == Null {
		return &Node{}, nil
	}

	value, err := store.Get([]byte(prefix + hash.String()))
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, nil
	}

	result := Node{
		L: common.BytesToHash(value[:32]),
		R: common.BytesToHash(value[32:]),
	}

	return &result, nil
}
