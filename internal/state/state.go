package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/state/enumerated"
	"github.com/freeverseio/laos-universal-node/internal/state/enumeratedtotal"
	"github.com/freeverseio/laos-universal-node/internal/state/ownership"
)

const (
	ContractPrefix = "contract_"
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
	OwnershipContractState
	EvolutionContractState
	OwnershipBlockState
	EvolutionBlockState
}

// State interface defines functions to interact with state of the blockchain
type State interface {
	CreateTreesForContract(contract common.Address) (ownership.Tree, enumerated.Tree, enumeratedtotal.Tree, error)
	SetTreesForContract(contract common.Address,
		ownershipTree ownership.Tree,
		enumeratedTree enumerated.Tree,
		enumeratedTotalTree enumeratedtotal.Tree)

	OwnerOf(contract common.Address, tokenId *big.Int) (common.Address, error)
	BalanceOf(contract, owner common.Address) (*big.Int, error)
	TokenOfOwnerByIndex(contract, owner common.Address, idx int) (*big.Int, error)
	TotalSupply(contract common.Address) (int64, error)
	TokenByIndex(contract common.Address, idx int) (*big.Int, error)
	Transfer(contract common.Address, eventTransfer *model.ERC721Transfer) error
	Mint(contract common.Address, tokenId *big.Int) error
	IsTreeSetForContract(contract common.Address) bool
	Get(key string) ([]byte, error)
	TagRoot(contract common.Address, blockNumber int64) error
	DeleteRootTag(contract common.Address, blockNumber int64) error
	GetLastTaggedBlock(contract common.Address) (int64, error)
	Checkout(contract common.Address, blockNumber int64) error
}

type OwnershipContractState interface {
	StoreERC721UniversalContracts(universalContracts []model.ERC721UniversalContract) error
	GetExistingERC721UniversalContracts(contracts []string) ([]string, error)
	GetCollectionAddress(contract string) (common.Address, error)
	GetAllERC721UniversalContracts() []string
}

type EvolutionContractState interface {
	GetMintedWithExternalURIEvents(contract string) ([]model.MintedWithExternalURI, error)
}

type OwnershipBlockState interface {
	SetCurrentEvoBlockForOwnershipContract(contract string, blockNumber uint64) error
	GetCurrentEvoBlockForOwnershipContract(contract string) (uint64, error)
	SetCurrentOwnershipBlock(number uint64) error
	GetCurrentOwnershipBlock() (uint64, error)
}

type EvolutionBlockState interface {
	SetCurrentEvoBlock(number uint64) error
	GetCurrentEvoBlock() (uint64, error)
}
