package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/account"
)

// Service interface is used for initializing and terminating state transaction.
type Service interface {
	NewTransaction() (Tx, error)
}

// Tx interface wraps all the available actions for state
type Tx interface {
	Discard()
	Commit() error

	State
	OwnershipContractState
	EvolutionContractState
	OwnershipSyncState
	EvolutionSyncState
}

// State interface defines functions to interact with state of the blockchain
type State interface {
	OwnerOf(contract common.Address, tokenId *big.Int) (common.Address, error)
	BalanceOf(contract, owner common.Address) (*big.Int, error)
	TokenOfOwnerByIndex(contract, owner common.Address, idx int) (*big.Int, error)
	TotalSupply(contract common.Address) (int64, error)
	TokenByIndex(contract common.Address, idx int) (*big.Int, error)
	TokenURI(contract common.Address, tokenId *big.Int) (string, error)
	Transfer(contract common.Address, eventTransfer *model.ERC721Transfer) error
	Mint(contract common.Address, mintEvent *model.MintedWithExternalURI) error
	LoadContractTrees(contractAddress common.Address) error
	UpdateContractState(contract common.Address, lastProcessedEvoBlock uint64) error
	AccountData(contract common.Address) (*account.AccountData, error)
	TagRoot(blockNumber int64) error
	DeleteRootTag(blockNumber int64) error
	DeleteOrphanRootTags(formBlock, toBlock int64) error
	GetLastTaggedBlock() (int64, error)
	Checkout(blockNumber int64) error
}

type OwnershipContractState interface {
	StoreERC721UniversalContracts(universalContracts []model.ERC721UniversalContract) error
	GetExistingERC721UniversalContracts(contracts []string) ([]string, error)
	GetCollectionAddress(contract string) (common.Address, error)
	GetAllERC721UniversalContracts() []string
	HasERC721UniversalContract(contract string) (bool, error)
}

type EvolutionContractState interface {
	GetMintedWithExternalURIEvents(contract string, blockNumber uint64) ([]model.MintedWithExternalURI, error)
	StoreMintedWithExternalURIEvent(contract string, event *model.MintedWithExternalURI) error
}

type OwnershipSyncState interface {
	SetLastOwnershipBlock(block model.Block) error
	GetLastOwnershipBlock() (model.Block, error)
	GetOwnershipBlock(blockNumber uint64) (model.Block, error)
	SetOwnershipBlock(blockNumber uint64, block model.Block) error
	GetAllStoredBlockNumbers() ([]uint64, error)
	DeleteOldStoredBlockNumbers() error
	DeleteOrphanBlockData(blockNumberRef uint64) error
}

type EvolutionSyncState interface {
	SetNextEvoEventBlock(contract string, blockNumber uint64) error
	GetNextEvoEventBlock(contract string, blockNumber uint64) (uint64, error)

	SetLastEvoBlock(block model.Block) error
	GetLastEvoBlock() (model.Block, error)
}
