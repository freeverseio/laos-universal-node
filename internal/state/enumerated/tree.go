package enumerated

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
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	prefix            = "enumerated/"
	treePrefix        = prefix + "tree/"
	headRootKeyPrefix = prefix + "head/"
	tokensPrefix      = prefix + "tokens/"
	tagPrefix         = prefix + "tags/"
	lastTagPrefix     = prefix + "lasttag/"
	treeDepth         = 160
)

// Tree is used to store enumerated tokens of each owner
type Tree interface {
	Root() common.Hash
	Transfer(minted bool, eventTransfer *model.ERC721Transfer) error
	Mint(tokenId *big.Int, owner common.Address) error
	TokensOf(owner common.Address) ([]big.Int, error)
	TagRoot(blockNumber int64) error
	GetLastTaggedBlock() (int64, error)
	DeleteRootTag(blockNumber int64) error 
	Checkout(blockNumber int64) error
	FindBlockWithTag(blockNumber int64) (int64, error)
}

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
	slog.Debug("enumeratedTokenTree", "HEAD", root.String())

	return &tree{contract, t, store, false}, err
}

// Mint creates a new token
func (b *tree) Mint(tokenId *big.Int, owner common.Address) error {
	tokens, err := b.TokensOf(owner)
	if err != nil {
		return err
	}

	tokens = append(tokens, *tokenId)
	if err := b.SetTokensToOwner(owner, tokens); err != nil {
		return err
	}

	return setHeadRoot(b.contract, b.store, b.Root())
}

// Transfer adds TokenId to the new owner and removes it from the previous owner
func (b *tree) Transfer(minted bool, eventTransfer *model.ERC721Transfer) error {
	if !minted {
		return nil
	}
	if eventTransfer.From.Cmp(common.Address{}) != 0 {
		fromTokens, err := b.TokensOf(eventTransfer.From)
		if err != nil {
			return err
		}

		for i, fromToken := range fromTokens {
			if fromToken.Cmp(eventTransfer.TokenId) == 0 {
				fromTokens[i] = fromTokens[len(fromTokens)-1]
				fromTokens = fromTokens[:len(fromTokens)-1]
				break
			}
		}

		if err := b.SetTokensToOwner(eventTransfer.From, fromTokens); err != nil {
			return err
		}
	}

	if eventTransfer.To.Cmp(common.Address{}) != 0 {
		toTokens, err := b.TokensOf(eventTransfer.To)
		if err != nil {
			return err
		}
		toTokens = append(toTokens, *eventTransfer.TokenId)
		if err := b.SetTokensToOwner(eventTransfer.To, toTokens); err != nil {
			return err
		}
	}

	return setHeadRoot(b.contract, b.store, b.Root())
}

// SetTokensToOwner sets the tokens of an owner
func (b *tree) SetTokensToOwner(owner common.Address, tokens []big.Int) error {
	buf, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	hash := crypto.Keccak256Hash(buf)
	if err := b.store.Set([]byte(tokensPrefix+b.contract.String()+"/"+hash.String()), buf); err != nil {
		return err
	}

	return b.mt.SetLeaf(owner.Big(), hash)
}

// TokensOf returns the tokens of an owner
func (b *tree) TokensOf(owner common.Address) ([]big.Int, error) {
	leaf, err := b.mt.Leaf(owner.Big())
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

	var tokens []big.Int
	if err := json.Unmarshal(buf, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
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

	lastTagKey := lastTagPrefix + b.contract.String()
	return b.store.Set([]byte(lastTagKey), []byte(strconv.FormatInt(blockNumber, 10)))
}

func (b *tree) GetLastTaggedBlock() (int64, error) {
	lastTagKey := lastTagPrefix + b.contract.String()
	buf, err := b.store.Get([]byte(lastTagKey))
	if err != nil {
		return 0, err
	}
	if len(buf) == 0{
		return 0, nil
	}

	return strconv.ParseInt(string(buf), 10, 64)
}
// DeleteRootTag deletes root tag
func (b *tree) DeleteRootTag(blockNumber int64) error {
	tagKey := tagPrefix + b.contract.String() + "/" + strconv.FormatInt(blockNumber, 10)
	return b.store.Delete([]byte(tagKey))
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
