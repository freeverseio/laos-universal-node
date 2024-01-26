package ownership

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree/jellyfish"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const ones160bits = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"

const (
	prefix            = "ownership/"
	treePrefix        = prefix + "tree/"
	headRootKeyPrefix = prefix + "head/"
	tokenDataPrefix   = prefix + "data/"
	tagPrefix         = prefix + "tags/"
	lastTagPrefix     = prefix + "lasttag/"
)

// TokenData defines token data placed in data of the leaf of the tree
type TokenData struct {
	SlotOwner common.Address
	TokenURI  string
	Minted    bool
	Idx       int
}

// Tree defines interface for ownership tree
type Tree interface {
	Root() common.Hash
	Transfer(eventTransfer *model.ERC721Transfer) error
	Mint(mintEvent *model.MintedWithExternalURI, idx int) error
	TokenData(tokenId *big.Int) (*TokenData, error)
	SetTokenData(tokenData *TokenData, tokenId *big.Int) error
	OwnerOf(tokenId *big.Int) (common.Address, error)
	SetRoot(root common.Hash) error
}

// EnumeratedTokensTree is used to store enumerated tokens of each SlotOwner
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
	slog.Debug("ownershipTree", "HEAD", root.String())

	return &tree{contract, t, store, false}, err
}

// Transfer updates the SlotOwner of the token
func (b *tree) Transfer(eventTransfer *model.ERC721Transfer) error {
	tokenData, err := b.TokenData(eventTransfer.TokenId)
	if err != nil {
		return err
	}

	tokenData.SlotOwner = eventTransfer.To
	if err := b.SetTokenData(tokenData, eventTransfer.TokenId); err != nil {
		return err
	}
	return setHeadRoot(b.contract, b.store, b.Root())
}

// Mint creates a new token
func (b *tree) Mint(mintEvent *model.MintedWithExternalURI, idx int) error {
	tokenData, err := b.TokenData(mintEvent.TokenId)
	if err != nil {
		return err
	}

	if tokenData.Minted {
		return errors.New("token " + mintEvent.TokenId.String() + " already minted")
	}

	tokenData.Minted = true
	tokenData.Idx = idx
	tokenData.TokenURI = mintEvent.TokenURI
	if err := b.SetTokenData(tokenData, mintEvent.TokenId); err != nil {
		return err
	}

	return setHeadRoot(b.contract, b.store, b.Root())
}

// SetTokenData updates the tokenData
func (b *tree) SetTokenData(tokenData *TokenData, tokenId *big.Int) error {
	buf, err := json.Marshal(tokenData)
	if err != nil {
		return err
	}

	hash := crypto.Keccak256Hash(buf)
	if err := b.store.Set([]byte(tokenDataPrefix+b.contract.String()+"/"+hash.String()), buf); err != nil {
		return err
	}

	return b.mt.SetLeaf(tokenId, hash)
}

// TokenData returns the tokenData
func (b *tree) TokenData(tokenId *big.Int) (*TokenData, error) {
	leaf, err := b.mt.Leaf(tokenId)
	if err != nil {
		return &TokenData{common.Address{}, "", false, 0}, err
	}
	if leaf.String() == jellyfish.Null {
		slotOwner := initSlotOwner(tokenId)
		return &TokenData{slotOwner, "", false, 0}, nil
	}

	buf, err := b.store.Get([]byte(tokenDataPrefix + b.contract.String() + "/" + leaf.String()))
	if err != nil {
		return &TokenData{common.Address{}, "", false, 0}, err
	}

	var tokenData TokenData
	if err := json.Unmarshal(buf, &tokenData); err != nil {
		return &TokenData{common.Address{}, "", false, 0}, err
	}

	return &tokenData, nil
}

// OwnerOf returns the owner of the token
func (b *tree) OwnerOf(tokenId *big.Int) (common.Address, error) {
	tokenData, err := b.TokenData(tokenId)
	if err != nil {
		return common.Address{}, err
	}

	if tokenData.Minted {
		return tokenData.SlotOwner, nil
	}

	return common.Address{}, nil
}

func initSlotOwner(tokenId *big.Int) common.Address {
	slotOwner, _ := new(big.Int).SetString(ones160bits, 16) // omit success. this is constant and it will always be true
	return common.BigToAddress(slotOwner.And(slotOwner, tokenId))
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

// SetRoot sets the current root to the one that is tagged for a blockNumber.
func (b *tree) SetRoot(root common.Hash) error {
	b.mt.SetRoot(root)
	return setHeadRoot(b.contract, b.store, root)
}
