package enumeratedtotal

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree/sparsemt"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	prefix               = "enumeratedtotal/"
	treePrefix           = prefix + "tree/"
	headRootKeyPrefix    = prefix + "head/"
	totalSupplyPrefix    = prefix + "totalsupply/"
	tokensPrefix         = prefix + "tokens/"
	tagPrefix            = prefix + "tags/"
	totalSupplyTagPrefix = prefix + "tags/totalsupply"
	treeDepth            = 64
)

// Tree defines interface for enumerated total tree
type Tree interface {
	Root() common.Hash
	Mint(tokenId *big.Int) error
	Burn(idx int) error
	TokenByIndex(idx int) (*big.Int, error)
	TotalSupply() (int64, error)
	TagRoot(blockNumber int64) error
	Checkout(blockNumber int64) error
	FindBlockWithTag(blockNumber int64) (int64, error)
}

// EnumeratedTokensTree is used to store enumerated tokens of each owner
type tree struct {
	contract common.Address
	mt       merkletree.MerkleTree
	store    storage.Tx
	tx       bool
}

// NewTree creates a new merkleTree with a custom storage
func NewTree(contract common.Address, store storage.Tx) (Tree, error) {
	if contract.Cmp(common.Address{}) == 0 {
		return nil, errors.New("contract address is " + common.Address{}.String())
	}

	t, err := sparsemt.New(treeDepth, store, treePrefix+contract.String())
	if err != nil {
		return nil, err
	}

	root, err := headRoot(contract, store)
	if err != nil {
		return nil, err
	}

	t.SetRoot(root)
	slog.Debug("enumeratedTotalTree", "HEAD", root.String())

	return &tree{contract, t, store, false}, err
}

// Mint creates a new token
func (b *tree) Mint(tokenId *big.Int) error {
	totalSupply, err := b.TotalSupply()
	if err != nil {
		return err
	}

	err = b.SetTokenToIndex(int(totalSupply), tokenId)
	if err != nil {
		return err
	}

	totalSupply++
	err = b.SetTotalSupply(totalSupply)
	if err != nil {
		return err
	}

	return setHeadRoot(b.contract, b.store, b.Root())
}

// Burn removes token )
func (b *tree) Burn(idx int) error {
	totalSupply, err := b.TotalSupply()
	if err != nil {
		return err
	}

	if idx >= int(totalSupply) {
		return errors.New("index out of totalSupply range")
	}

	tokenId, err := b.TokenByIndex(int(totalSupply - 1))
	if err != nil {
		return err
	}

	err = b.SetTokenToIndex(idx, tokenId)
	if err != nil {
		return err
	}

	err = b.SetTokenToIndex(int(totalSupply)-1, big.NewInt(0))
	if err != nil {
		return err
	}

	totalSupply--
	err = b.SetTotalSupply(totalSupply)
	if err != nil {
		return err
	}

	return setHeadRoot(b.contract, b.store, b.Root())
}

// SetTotalSupply sets to total number of token in the contract
func (b *tree) SetTotalSupply(totalSupply int64) error {
	return b.store.Set([]byte(totalSupplyPrefix+b.contract.String()), []byte(strconv.FormatInt(totalSupply, 10)))
}

func (b *tree) TotalSupply() (int64, error) {
	buf, err := b.store.Get([]byte(totalSupplyPrefix + b.contract.String()))
	if err != nil {
		return 0, err
	}

	if len(buf) == 0 {
		return 0, nil
	}

	return strconv.ParseInt(string(buf), 10, 64)
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
	totalSupply, err := b.TotalSupply()
	if err != nil {
		return nil, err
	}

	if idx >= int(totalSupply) {
		return big.NewInt(0), errors.New("index out of totalSupply range")
	}

	leaf, err := b.mt.Leaf(big.NewInt(int64(idx)))
	if err != nil {
		return nil, err
	}

	if leaf.String() == sparsemt.Null {
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

	totalSupply, err := b.TotalSupply()
	if err != nil {
		return err
	}
	tagTotalSupplyKey := totalSupplyTagPrefix + b.contract.String() + "/" + strconv.FormatInt(blockNumber, 10)
	return b.store.Set([]byte(tagTotalSupplyKey), []byte(strconv.FormatInt(totalSupply, 10)))
}

// Checkout sets the current root to the one that is tagged for a blockNumber.
func (b *tree) Checkout(blockNumber int64) error {
	tagKey := tagPrefix + b.contract.String() + "/" + strconv.FormatInt(blockNumber, 10)
	buf, err := b.store.Get([]byte(tagKey))
	if err != nil {
		return err
	}

	if len(buf) == 0 {
		return errors.New("no tog found for this block number " + strconv.FormatInt(blockNumber, 10))
	}

	newRoot := common.BytesToHash(buf)
	b.mt.SetRoot(newRoot)
	err = setHeadRoot(b.contract, b.store, newRoot)
	if err != nil {
		return err
	}

	tagTotalSupplyKey := totalSupplyTagPrefix + b.contract.String() + "/" + strconv.FormatInt(blockNumber, 10)
	buf, err = b.store.Get([]byte(tagTotalSupplyKey))
	if err != nil {
		return err
	}

	return b.store.Set([]byte(totalSupplyPrefix+b.contract.String()), buf)
}

// FindBlockWithTag returns the first previous blockNumber that has been tagged if the tag for the blockNumber does not
// exist
func (b *tree) FindBlockWithTag(blockNumber int64) (int64, error) {
	for {
		if blockNumber == 0 {
			return 0, nil
		}

		buf, err := b.store.Get([]byte(tagPrefix + b.contract.String() + "/" + strconv.FormatInt(blockNumber, 10)))
		if err != nil {
			return 0, err
		}

		if len(buf) != 0 {
			return blockNumber, nil
		}

		blockNumber--
	}
}
