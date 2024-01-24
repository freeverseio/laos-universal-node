package enumerated

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree/jellyfish"
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
)

// Tree is used to store enumerated tokens of each owner
type Tree interface {
	Root() common.Hash
	Transfer(minted bool, eventTransfer *model.ERC721Transfer) error
	Mint(tokenId *big.Int, owner common.Address) error
	TokenOfOwnerByIndex(owner common.Address, idx uint64) (*big.Int, error)
	SetTokenToOwnerToIndex(owner common.Address, idx uint64, token *big.Int) error
	SetBalanceToOwner(owner common.Address, balance uint64) error
	BalanceOfOwner(owner common.Address) (uint64, error)
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
	slog.Debug("enumeratedTokenTree", "HEAD", root.String())

	return &tree{contract, t, store, false}, err
}

// Mint creates a new token
func (b *tree) Mint(tokenId *big.Int, owner common.Address) error {
	balance, err := b.BalanceOfOwner(owner)
	if err != nil {
		return err
	}
	if err := b.SetTokenToOwnerToIndex(owner, balance, tokenId); err != nil {
		return err
	}

	return b.SetBalanceToOwner(owner, balance+1)
}

// Transfer adds TokenId to the new owner and removes it from the previous owner
func (b *tree) Transfer(minted bool, eventTransfer *model.ERC721Transfer) error {
	if !minted {
		return nil
	}

	if eventTransfer.From.Cmp(common.Address{}) != 0 {
		fromBalance, err := b.BalanceOfOwner(eventTransfer.From)
		if err != nil {
			return err
		}
		for i := uint64(0); i < fromBalance; i++ {
			fromToken, err := b.TokenOfOwnerByIndex(eventTransfer.From, i)
			if err != nil {
				return err
			}

			if fromToken.Cmp(eventTransfer.TokenId) == 0 {
				lastToken, err := b.TokenOfOwnerByIndex(eventTransfer.From, fromBalance-1)
				if err != nil {
					return err
				}

				if err := b.SetTokenToOwnerToIndex(eventTransfer.From, i, lastToken); err != nil {
					return err
				}

				if err := b.SetBalanceToOwner(eventTransfer.From, fromBalance-1); err != nil {
					return err
				}

				break
			}
		}
	}

	if eventTransfer.To.Cmp(common.Address{}) != 0 {
		toBalance, err := b.BalanceOfOwner(eventTransfer.To)
		if err != nil {
			return err
		}

		if err := b.SetTokenToOwnerToIndex(eventTransfer.To, toBalance, eventTransfer.TokenId); err != nil {
			return err
		}

		if err := b.SetBalanceToOwner(eventTransfer.To, toBalance+1); err != nil {
			return err
		}
	}

	return nil
}

// SetTokensToOwner sets the tokens of an owner
func (b *tree) SetTokenToOwnerToIndex(owner common.Address, idx uint64, token *big.Int) error {
	buf, err := json.Marshal(token)
	if err != nil {
		return err
	}

	hash := crypto.Keccak256Hash(buf)
	if err := b.store.Set([]byte(tokensPrefix+b.contract.String()+"/"+hash.String()), buf); err != nil {
		return err
	}

	position := owner.Big()
	position = position.Lsh(position, 64)
	position = position.Add(position, big.NewInt(int64(idx+1))) // +1 because balance is stored at index 0

	if err := b.mt.SetLeaf(position, hash); err != nil {
		return err
	}

	return setHeadRoot(b.contract, b.store, b.Root())
}

// TokensOf returns the tokens of an owner
func (b *tree) TokenOfOwnerByIndex(owner common.Address, idx uint64) (*big.Int, error) {
	balance, err := b.BalanceOfOwner(owner)
	if err != nil {
		return big.NewInt(0), err
	}

	if idx >= balance {
		return big.NewInt(0), fmt.Errorf("index %d out of range", idx)
	}

	position := owner.Big()
	position = position.Lsh(position, 64)
	position = position.Add(position, big.NewInt(int64(idx+1))) // +1 because balance is stored at index 0
	leaf, err := b.mt.Leaf(position)
	if err != nil {
		return big.NewInt(0), err
	}
	if leaf.String() == jellyfish.Null {
		return big.NewInt(0), nil
	}

	buf, err := b.store.Get([]byte(tokensPrefix + b.contract.String() + "/" + leaf.String()))
	if err != nil {
		return big.NewInt(0), err
	}

	var token big.Int
	if err := json.Unmarshal(buf, &token); err != nil {
		return big.NewInt(0), err
	}

	return &token, nil
}

// SetTokensToOwner sets the tokens of an owner
func (b *tree) SetBalanceToOwner(owner common.Address, balance uint64) error {
	buf := []byte(strconv.FormatUint(balance, 10))

	hash := crypto.Keccak256Hash(buf)
	if err := b.store.Set([]byte(tokensPrefix+b.contract.String()+"/"+hash.String()), buf); err != nil {
		return err
	}

	position := owner.Big()
	position = position.Lsh(position, 64)
	if err := b.mt.SetLeaf(position, hash); err != nil {
		return err
	}

	return setHeadRoot(b.contract, b.store, b.Root())
}

func (b *tree) BalanceOfOwner(owner common.Address) (uint64, error) {
	position := owner.Big()
	position = position.Lsh(position, 64)
	leaf, err := b.mt.Leaf(position)
	if err != nil {
		return 0, err
	}
	if leaf.String() == jellyfish.Null {
		return 0, nil
	}

	buf, err := b.store.Get([]byte(tokensPrefix + b.contract.String() + "/" + leaf.String()))
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(string(buf), 10, 64)
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
