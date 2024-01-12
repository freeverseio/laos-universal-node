package jellyfish

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
	"github.com/pokt-network/smt"
)

type jellyfish struct {
	tree   *smt.SMT
	prefix string
	store  storage.Tx
}

func New(store storage.Tx, prefix string) (merkletree.MerkleTree, error) {
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	return &jellyfish{
		prefix: prefix,
		store:  store,
		tree: smt.NewSparseMerkleTrie(
			store,
			crypto.NewKeccakState(),
			smt.WithValueHasher(nil),
		),
	}, nil
}

func (j *jellyfish) Leaf(idx *big.Int) (common.Hash, error) {
	val, err := j.tree.Get([]byte(j.prefix + "/" + idx.String()))
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(val), nil
}

func (j *jellyfish) SetLeaf(idx *big.Int, hash common.Hash) error {
	err := j.tree.Update([]byte(j.prefix+"/"+idx.String()), hash.Bytes())
	if err != nil {
		return err
	}
	return j.tree.Commit()
}

func (j *jellyfish) Root() common.Hash {
	return common.BytesToHash(j.tree.Root())
}

func (j *jellyfish) SetRoot(hash common.Hash) {
	j.tree = smt.ImportSparseMerkleTrie(j.store, crypto.NewKeccakState(), hash.Bytes(), smt.WithValueHasher(nil))
}

// added only to comply with the interface
func (j *jellyfish) Path(idx *big.Int) ([]common.Hash, error) {
	return nil, nil
}

// added only to comply with the interface
func (j *jellyfish) Proof(idx *big.Int) ([]common.Hash, error) {
	return nil, nil
}

// added only to comply with the interface
func (j *jellyfish) CountLeaves() *big.Int {
	return big.NewInt(0)
}

// added only to comply with the interface
func (j *jellyfish) Depth() uint {
	return 0
}
