package jellyfish

import (
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
	"github.com/lazyledger/smt"
	"golang.org/x/crypto/blake2b"
)

const Null = "0x0000000000000000000000000000000000000000000000000000000000000000"

type smtStore struct {
	prefix string
	s      storage.Tx
}

// Gets value for a given key
func (b smtStore) Get(key []byte) ([]byte, error) {
	return b.s.Get(append([]byte(b.prefix), key...))
}

// Deletes key value pair
func (b smtStore) Delete(key []byte) error {
	return b.s.Delete(append([]byte(b.prefix), key...))
}

// Sets value for a given key
func (b smtStore) Set(key, value []byte) error {
	return b.s.Set(append([]byte(b.prefix), key...), value)
}

type jellyfish struct {
	tree *smt.SparseMerkleTree
}

func New(store storage.Tx, prefix string) (merkletree.MerkleTree, error) {
	slog.Debug("create jelly fish merkle tree with prefix", "prefix", prefix)
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}

	hasher, err := blake2b.New256(nil)
	if err != nil {
		return nil, err
	}

	tree := smt.NewSparseMerkleTree(smtStore{prefix: prefix, s: store}, hasher)
	return &jellyfish{tree}, nil
}

func (j *jellyfish) Leaf(idx *big.Int) (common.Hash, error) {
	val, err := j.tree.Get([]byte(idx.String()))
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(val), nil
}

func (j *jellyfish) SetLeaf(idx *big.Int, hash common.Hash) error {
	_, err := j.tree.Update([]byte(idx.String()), hash.Bytes())
	return err
}

func (j *jellyfish) Root() common.Hash {
	return common.BytesToHash(j.tree.Root())
}

func (j *jellyfish) SetRoot(hash common.Hash) {
	j.tree.SetRoot(hash.Bytes())
}
