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
	prefix       = "enumerated/"
	treePrefix   = prefix + "tree/"
	tokensPrefix = prefix + "tokens/"
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
	SetRoot(root common.Hash)
}

type tree struct {
	contract common.Address
	mt       merkletree.MerkleTree
	store    storage.Tx
	tx       bool
}

// NewTree creates a new merkleTree with a custom storage
func NewTree(contract common.Address, root common.Hash, store storage.Tx) (Tree, error) {
	if contract.Cmp(common.Address{}) == 0 {
		return nil, errors.New("contract address is " + common.Address{}.String())
	}

	t, err := jellyfish.New(store, treePrefix+contract.String())
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

	return b.mt.SetLeaf(position, hash)
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

	return b.mt.SetLeaf(position, hash)
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

// SetRoot sets the current root to the one that is tagged for a blockNumber.
func (b *tree) SetRoot(root common.Hash) {
	b.mt.SetRoot(root)
}
