package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/state/enumerated"
	"github.com/freeverseio/laos-universal-node/internal/state/enumeratedtotal"
	"github.com/freeverseio/laos-universal-node/internal/state/ownership"
)

// Service interface is used for initializing and terminating state transaction.
type Service interface {
	NewTransaction() Tx
}

// Tx interface wraps all the available actions for state
type Tx interface {
	Discard()
	Commit() error

	State
	ContractState
}

// State interface defines functions to interact with state of the blockchain
type State interface {
	Discard()
	Commit() error

	CreateTreesForContract(contract common.Address) (ownership.Tree, enumerated.Tree, enumeratedtotal.Tree, error)
	SetTreesForContract(contract common.Address,
		ownershipTree ownership.Tree,
		enumeratedTree enumerated.Tree,
		enumeratedTotalTree enumeratedtotal.Tree) error

	OwnerOf(contract common.Address, tokenId *big.Int) (common.Address, error)
	BalanceOf(contract, owner common.Address) (*big.Int, error)
	TokenOfOwnerByIndex(contract, owner common.Address, idx int) (*big.Int, error)
	TotalSupply(contract common.Address) (int64, error)
	TokenByIndex(contract common.Address, idx int) (*big.Int, error)
	Transfer(contract common.Address, eventTransfer scan.EventTransfer) error
	Mint(contract common.Address, tokenId *big.Int) error
}

type ContractState interface {
	StoreERC721UniversalContracts(universalContracts []model.ERC721UniversalContract) error
	StoreEvoChainMintEvents(contract common.Address, events []model.EventMintedWithExternalURI) error
	EvoChainMintEvents(contract common.Address) ([]model.EventMintedWithExternalURI, error)
}
