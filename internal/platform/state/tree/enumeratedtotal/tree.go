package enumeratedtotal

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree/jellyfish"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	prefix       = "enumeratedtotal/"
	treePrefix   = prefix + "tree/"
	tokensPrefix = prefix + "tokens/"
)

// Tree defines interface for enumerated total tree
type Tree interface {
	Root() common.Hash
	Mint(tokenId *big.Int) error
	Burn(idx int) error
	TokenByIndex(idx int) (*big.Int, error)
	TotalSupply() int64
	SetTotalSupply(totalSupply int64)
	SetRoot(root common.Hash)
}

// EnumeratedTokensTree is used to store enumerated tokens of each owner
type tree struct {
	contract common.Address
	mt       merkletree.MerkleTree
	store    storage.Tx
	tx       bool
	// total supply is currently not store in the tree.
	// We should measure how much storing it in the tree (in a similar way as in enumerated tree) reduces max tx size.
	// if not too much I suggest to put in the tree because code would be simpler
	totalSupply int64
}

// NewTree creates a new merkleTree with a custom storage
func NewTree(contract common.Address, root common.Hash, totalSupply int64, store storage.Tx) (Tree, error) {
	if contract.Cmp(common.Address{}) == 0 {
		return nil, errors.New("contract address is " + common.Address{}.String())
	}

	t, err := jellyfish.New(store, treePrefix+contract.String())
	if err != nil {
		return nil, err
	}

	t.SetRoot(root)
	slog.Debug("enumeratedTotalTree", "HEAD", root.String())

	return &tree{contract, t, store, false, totalSupply}, err
}

// Mint creates a new token
func (b *tree) Mint(tokenId *big.Int) error {
	err := b.SetTokenToIndex(int(b.totalSupply), tokenId)
	if err != nil {
		return err
	}

	b.totalSupply++

	return nil
}

// Burn removes token )
func (b *tree) Burn(idx int) error {
	if idx >= int(b.totalSupply) {
		return errors.New("index out of totalSupply range")
	}

	tokenId, err := b.TokenByIndex(int(b.totalSupply - 1))
	if err != nil {
		return err
	}

	err = b.SetTokenToIndex(idx, tokenId)
	if err != nil {
		return err
	}

	err = b.SetTokenToIndex(int(b.totalSupply)-1, big.NewInt(0))
	if err != nil {
		return err
	}

	b.totalSupply--

	return nil
}

// SetTotalSupply sets to total number of token in the contract
func (b *tree) SetTotalSupply(totalSupply int64) {
	b.totalSupply = totalSupply
}

func (b *tree) TotalSupply() int64 {
	return b.totalSupply
}

// SetTokenIndex sets the token index
func (b *tree) SetTokenToIndex(idx int, token *big.Int) error {
	buf, err := json.Marshal(token)
	if err != nil {
		return err
	}

	hash := crypto.Keccak256Hash(buf)
	if err := b.store.Set([]byte(tokensPrefix+b.contract.String()+"/"+hash.String()), buf); err != nil {
		return err
	}

	return b.mt.SetLeaf(big.NewInt(int64(idx)), hash)
}

// TokenByIndex returns the token by index
func (b *tree) TokenByIndex(idx int) (*big.Int, error) {
	if idx >= int(b.totalSupply) {
		return big.NewInt(0), errors.New("index out of totalSupply range")
	}

	leaf, err := b.mt.Leaf(big.NewInt(int64(idx)))
	if err != nil {
		return nil, err
	}

	if leaf.String() == jellyfish.Null {
		return nil, nil
	}

	buf, err := b.store.Get([]byte(tokensPrefix + b.contract.String() + "/" + leaf.String()))
	if err != nil {
		return nil, err
	}

	var token big.Int
	if err := json.Unmarshal(buf, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

// Root returns the root of the tree
func (b *tree) Root() common.Hash {
	return b.mt.Root()
}

// SetRoot sets the current root to the one that is tagged for a blockNumber.
func (b *tree) SetRoot(root common.Hash) {
	b.mt.SetRoot(root)
}
