package account

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree/jellyfish"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	prefix            = "account"
	treePrefix        = prefix + "tree/"
	headRootKeyPrefix = prefix + "head/"
	leafDataPrefix    = prefix + "data/"
	tagPrefix         = prefix + "tags/"
	lastTagPrefix     = prefix + "lasttag/"
)

// MerkleTreeRoots defines the roots from enumerated, enumerated total and ownership merkle trees
// placed in data of the leaf of the tree
type MerkleTreeRoots struct {
	Enumerated      common.Hash
	EnumeratedTotal common.Hash
	Ownership       common.Hash
}

// Tree defines interface for the account tree
type Tree interface {
	Root() common.Hash
	MerkleTreeRoots(merkleTreeID *big.Int) (*MerkleTreeRoots, error)
	SetMerkleTreeRoots(roots *MerkleTreeRoots, merkleTreeID *big.Int) error
	TagRoot(blockNumber int64) error
	GetLastTaggedBlock() (int64, error)
	Checkout(blockNumber int64) error
}

type tree struct {
	contract common.Address
	mt       merkletree.MerkleTree
	store    storage.Tx
	tx       bool
}

// TODO NewTree should be GetTree (same for enumerated and enumeratedtotal packages)

// NewTree creates a new merkleTree with a custom storage
func NewTree(contract common.Address, store storage.Tx) (Tree, error) {
	if contract.Cmp(common.Address{}) == 0 {
		return nil, errors.New("contract address is " + common.Address{}.String())
	}

	t, err := jellyfish.New(store, treePrefix+contract.String())
	if err != nil {
		return nil, err
	}

	root, err := headRoot(contract, store)
	if err != nil {
		return nil, err
	}

	t.SetRoot(root)
	slog.Debug("accountTree", "HEAD", root.String())

	return &tree{contract, t, store, false}, err
}

// SetMerkleTreeRoots updates the MerkleTreeRoots
func (b *tree) SetMerkleTreeRoots(data *MerkleTreeRoots, accountID *big.Int) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	hash := crypto.Keccak256Hash(buf)
	if err := b.store.Set([]byte(leafDataPrefix+b.contract.String()+"/"+hash.String()), buf); err != nil {
		return err
	}

	return b.mt.SetLeaf(accountID, hash)
}

// MerkleTreeRoots returns the merkle trees roots
func (b *tree) MerkleTreeRoots(accountID *big.Int) (*MerkleTreeRoots, error) {
	leafHash, err := b.mt.Leaf(accountID)
	if err != nil {
		return &MerkleTreeRoots{}, err
	}
	if leafHash.String() == jellyfish.Null {
		return &MerkleTreeRoots{}, err
	}

	buf, err := b.store.Get([]byte(leafDataPrefix + b.contract.String() + "/" + leafHash.String()))
	if err != nil {
		return &MerkleTreeRoots{}, err
	}

	var roots MerkleTreeRoots
	if err := json.Unmarshal(buf, &roots); err != nil {
		return &MerkleTreeRoots{}, err
	}

	return &roots, nil
}

func (b *tree) Root() common.Hash {
	return b.mt.Root()
}

func headRoot(contract common.Address, store storage.Tx) (common.Hash, error) {
	buf, err := store.Get([]byte(headRootKeyPrefix + contract.String()))
	if err != nil {
		return common.Hash{}, err
	}

	if len(buf) == 0 {
		return common.Hash{}, nil
	}

	return common.BytesToHash(buf), nil
}

func setHeadRoot(contract common.Address, store storage.Tx, root common.Hash) error {
	return store.Set([]byte(headRootKeyPrefix+contract.String()), root.Bytes())
}

// TagRoot stores a root value for the block so that it can be checked later
func (b *tree) TagRoot(blockNumber int64) error {
	tagKey := tagPrefix + b.contract.String() + "/" + strconv.FormatInt(blockNumber, 10)
	root := b.Root()
	err := b.store.Set([]byte(tagKey), root.Bytes())
	if err != nil {
		return err
	}

	lastTagKey := lastTagPrefix + b.contract.String()
	return b.store.Set([]byte(lastTagKey), []byte(strconv.FormatInt(blockNumber, 10)))
}

func (b *tree) GetLastTaggedBlock() (int64, error) {
	lastTagKey := lastTagPrefix + b.contract.String()
	buf, err := b.store.Get([]byte(lastTagKey))
	if err != nil {
		return 0, err
	}
	if len(buf) == 0 {
		return 0, nil
	}

	return strconv.ParseInt(string(buf), 10, 64)
}

// Checkout sets the current root to the one that is tagged for a blockNumber.
func (b *tree) Checkout(blockNumber int64) error {
	tagKey := tagPrefix + b.contract.String() + "/" + strconv.FormatInt(blockNumber, 10)
	buf, err := b.store.Get([]byte(tagKey))
	if err != nil {
		return err
	}

	if len(buf) == 0 {
		return errors.New("no tag found for this block number " + strconv.FormatInt(blockNumber, 10))
	}

	newRoot := common.BytesToHash(buf)
	b.mt.SetRoot(newRoot)
	return setHeadRoot(b.contract, b.store, newRoot)
}

// DeleteRootTag deletes root tag without loading the tree
func DeleteRootTag(tx storage.Tx, contract string, blockNumber int64) error {
	tagKey := tagPrefix + contract + "/" + strconv.FormatInt(blockNumber, 10)
	return tx.Delete([]byte(tagKey))
}
