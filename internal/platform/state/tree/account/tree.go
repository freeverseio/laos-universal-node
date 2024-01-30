package account

import (
	"encoding/json"
	"errors"
	"log/slog"
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

// AccountData defines the roots from enumerated, enumerated total and ownership merkle trees
// placed in data of the leaf of the tree
type AccountData struct {
	EnumeratedRoot        common.Hash
	EnumeratedTotalRoot   common.Hash
	OwnershipRoot         common.Hash
	TotalSupply           int64
	LastProcessedEvoBlock uint64
}

// Tree defines interface for the account tree
type Tree interface {
	Root() common.Hash
	AccountData(contract common.Address) (*AccountData, error)
	SetAccountData(data *AccountData, accountAddress common.Address) error
	TagRoot(blockNumber int64) error
	GetLastTaggedBlock() (int64, error)
	Checkout(blockNumber int64) error
	DeleteRootTag(blockNumber int64) error
}

type tree struct {
	mt    merkletree.MerkleTree
	store storage.Tx
	tx    bool
}

// TODO NewTree should be GetTree (same for enumerated and enumeratedtotal packages)

// NewTree creates a new merkleTree with a custom storage
func NewTree(store storage.Tx) (Tree, error) {
	t, err := jellyfish.New(store, treePrefix)
	if err != nil {
		return nil, err
	}

	root, err := headRoot(store)
	if err != nil {
		return nil, err
	}

	t.SetRoot(root)
	slog.Debug("accountTree", "HEAD", root.String())

	return &tree{t, store, false}, err
}

// SetAccountData updates the MerkleTreeRoots
func (b *tree) SetAccountData(data *AccountData, address common.Address) error {
	slog.Debug("SetAccountData", "data", data, "address", address.String())
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	hash := crypto.Keccak256Hash(buf)
	if err := b.store.Set([]byte(leafDataPrefix+hash.String()), buf); err != nil {
		return err
	}

	err = b.mt.SetLeaf(address.Big(), hash)
	if err != nil {
		return err
	}

	slog.Debug("accountTree", "HEAD", b.Root().String())
	return setHeadRoot(b.store, b.Root())
}

// AccountData returns the merkle trees roots
func (b *tree) AccountData(address common.Address) (*AccountData, error) {
	leafHash, err := b.mt.Leaf(address.Big())
	if err != nil {
		return &AccountData{}, err
	}
	if leafHash.String() == jellyfish.Null {
		return &AccountData{}, err
	}

	buf, err := b.store.Get([]byte(leafDataPrefix + leafHash.String()))
	if err != nil {
		return &AccountData{}, err
	}

	var roots AccountData
	if err := json.Unmarshal(buf, &roots); err != nil {
		return &AccountData{}, err
	}

	return &roots, nil
}

func (b *tree) Root() common.Hash {
	return b.mt.Root()
}

func headRoot(store storage.Tx) (common.Hash, error) {
	buf, err := store.Get([]byte(headRootKeyPrefix))
	if err != nil {
		return common.Hash{}, err
	}

	if len(buf) == 0 {
		return common.Hash{}, nil
	}

	return common.BytesToHash(buf), nil
}

func setHeadRoot(store storage.Tx, root common.Hash) error {
	return store.Set([]byte(headRootKeyPrefix), root.Bytes())
}

// TagRoot stores a root value for the block so that it can be checked later
func (b *tree) TagRoot(blockNumber int64) error {
	tagKey := tagPrefix + strconv.FormatInt(blockNumber, 10)
	root := b.Root()
	err := b.store.Set([]byte(tagKey), root.Bytes())
	if err != nil {
		return err
	}

	return b.store.Set([]byte(lastTagPrefix), []byte(strconv.FormatInt(blockNumber, 10)))
}

func (b *tree) GetLastTaggedBlock() (int64, error) {
	buf, err := b.store.Get([]byte(lastTagPrefix))
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
	tagKey := tagPrefix + strconv.FormatInt(blockNumber, 10)
	buf, err := b.store.Get([]byte(tagKey))
	if err != nil {
		return err
	}

	if len(buf) == 0 {
		return errors.New("no tag found for this block number " + strconv.FormatInt(blockNumber, 10))
	}

	newRoot := common.BytesToHash(buf)
	b.mt.SetRoot(newRoot)
	return setHeadRoot(b.store, newRoot)
}

// DeleteRootTag deletes root tag without loading the tree
func (b *tree) DeleteRootTag(blockNumber int64) error {
	tagKey := tagPrefix + strconv.FormatInt(blockNumber, 10)
	return b.store.Delete([]byte(tagKey))
}

