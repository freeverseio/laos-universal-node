// Package merkle provides implementations for merkle tree.
package merkletree

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// MerkleTree interface defines functions to interact with merkle tree
type MerkleTree interface {
	SetLeaf(idx *big.Int, hash common.Hash) error
	Leaf(idx *big.Int) (common.Hash, error)
	Root() common.Hash
	SetRoot(hash common.Hash)
}
