package sparsemt

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

// Tree implements sparse merkle tree
type Tree struct {
	store   storage.Tx
	prefix  string
	depth   uint
	root    common.Hash
	nLeaves *big.Int
}

// New creates a new merkle tree with storage capability, depth and prefix
func New(depth uint, store storage.Tx, prefix string) (merkletree.MerkleTree, error) {
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	return &Tree{
		store:   store,
		depth:   depth,
		nLeaves: new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(depth)), nil),
		prefix:  prefix,
	}, nil
}

// SetRoot sets the root of the tree to "hash"
func (b *Tree) SetRoot(hash common.Hash) {
	b.root = hash
}

// Root returns the root of the tree
func (b *Tree) Root() common.Hash {
	return b.root
}

// Proof calculates the proofs from leaf at index "idx" going up from the root. It returns a slice of hashes of siblings
// of each node that is found in the path from the root to the leaf
func (b *Tree) Proof(idx *big.Int) ([]common.Hash, error) {
	if idx.Cmp(b.nLeaves) != -1 {
		return nil, fmt.Errorf("Tree:Proof:out of bounds")
	}

	node, err := GetNode(b.store, b.root, b.prefix)
	if err != nil {
		return nil, err
	}
	var proof []common.Hash

	for i := int(b.depth - 1); i >= 0; i-- {
		if node == nil {
			return nil, fmt.Errorf("trunk branch in level %v, tree depth %v", i, b.depth)
		}
		if idx.Bit(i) == 1 {
			proof = append(proof, node.L)
			if node, err = GetNode(b.store, node.R, b.prefix); err != nil {
				return nil, err
			}
		} else {
			proof = append(proof, node.R)
			if node, err = GetNode(b.store, node.L, b.prefix); err != nil {
				return nil, err
			}
		}
	}

	return proof, nil
}

// Leaf returns the leaf at index "idx"
func (b *Tree) Leaf(idx *big.Int) (common.Hash, error) {
	path, err := b.Path(idx)
	if err != nil {
		return common.Hash{}, err
	}

	if len(path) == 0 {
		return b.root, nil
	}

	return path[b.depth-1], nil
}

// SetLeaf sets a new leaf at index "idx"
func (b *Tree) SetLeaf(idx *big.Int, hash common.Hash) error {
	proof, err := b.Proof(idx)
	if err != nil {
		return err
	}

	var node Node
	for i := 0; i < int(b.depth); i++ {
		if idx.Bit(i) == 1 {
			// Examine the leaf bit to determine whether to place the hash as the left (L) or right (R) sibling.
			// If idx.Big(i) == 1, the new hash is placed in the right (R) sibling, while the left (L) sibling
			// receives the corresponding proof hash.
			node.L = proof[int(b.depth)-i-1]
			node.R = hash
		} else {
			node.L = hash
			node.R = proof[int(b.depth)-i-1]
		}
		hash, err = PutNode(b.store, node, b.prefix)
		if err != nil {
			return err
		}
	}

	b.root = hash

	return nil
}

// Path returns the path (i.e. the hashes) leading to the leaf at index "idx"
func (b *Tree) Path(idx *big.Int) ([]common.Hash, error) {
	if idx.Cmp(b.nLeaves) != -1 {
		return nil, fmt.Errorf("Tree:Path:out of bounds")
	}

	var path []common.Hash
	if b.depth == 0 {
		return path, nil
	}

	hash := b.root
	for i := int(b.depth - 1); i >= 0; i-- {
		node, err := GetNode(b.store, hash, b.prefix)
		if err != nil {
			return nil, err
		}
		if node == nil {
			return nil, fmt.Errorf("unexistent node. tree is broken")
		}
		if idx.Bit(i) == 1 {
			hash = node.R
		} else {
			hash = node.L
		}
		path = append(path, hash)
	}

	return path, nil
}

// CountLeaves returns the total number of leaves
func (b *Tree) CountLeaves() *big.Int {
	return new(big.Int).Set(b.nLeaves)
}

// Depth returns the depth of the tree
func (b *Tree) Depth() uint {
	return b.depth
}
