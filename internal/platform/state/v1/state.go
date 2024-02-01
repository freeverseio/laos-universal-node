package v1

import (
	"fmt"
	"log/slog"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	evolutionContractState "github.com/freeverseio/laos-universal-node/internal/platform/state/contract/evolution"
	ownershipContractState "github.com/freeverseio/laos-universal-node/internal/platform/state/contract/ownership"
	evolutionSyncState "github.com/freeverseio/laos-universal-node/internal/platform/state/sync/evolution"
	ownershipSyncState "github.com/freeverseio/laos-universal-node/internal/platform/state/sync/ownership"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/account"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumerated"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/enumeratedtotal"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/ownership"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

type service struct {
	storageService storage.Service
}

// NewStateService creates a new state service
func NewStateService(storageService storage.Service) state.Service {
	return &service{
		storageService: storageService,
	}
}

// Creates a new state transaction
func (s *service) NewTransaction() (state.Tx, error) {
	storageTx := s.storageService.NewTransaction()
	accountTree, err := account.NewTree(storageTx)
	if err != nil {
		return nil, err
	}

	return &tx{
		ownershipTrees:         make(map[common.Address]ownership.Tree),
		enumeratedTrees:        make(map[common.Address]enumerated.Tree),
		enumeratedTotalTrees:   make(map[common.Address]enumeratedtotal.Tree),
		accountTree:            accountTree,
		tx:                     storageTx,
		OwnershipContractState: ownershipContractState.NewService(storageTx),
		EvolutionContractState: evolutionContractState.NewService(storageTx),
		OwnershipSyncState:     ownershipSyncState.NewService(storageTx),
		EvolutionSyncState:     evolutionSyncState.NewService(storageTx),
	}, nil
}

type tx struct {
	tx                   storage.Tx
	ownershipTrees       map[common.Address]ownership.Tree
	enumeratedTrees      map[common.Address]enumerated.Tree
	enumeratedTotalTrees map[common.Address]enumeratedtotal.Tree
	accountTree          account.Tree
	state.OwnershipContractState
	state.EvolutionContractState
	state.OwnershipSyncState
	state.EvolutionSyncState
}

// isTreeSetForContact returns true if the tree is set
func (t *tx) isTreeSetForContract(contract common.Address) bool {
	_, ok := t.ownershipTrees[contract]
	return ok
}

// createTreesForContract creates new trees for contract (ownership, enumerated, and enumeratedtotal)
func (t *tx) createTreesForContract(contract common.Address) (
	ownershipTree ownership.Tree,
	enumeratedTree enumerated.Tree,
	enumeratedTotalTree enumeratedtotal.Tree,
	err error,
) {
	slog.Debug("creating trees for contract", "contract", contract.String())
	accountData, err := t.accountTree.AccountData(contract)
	if err != nil {
		return nil, nil, nil, err
	}

	ownershipTree, err = ownership.NewTree(contract, accountData.OwnershipRoot, t.tx)
	if err != nil {
		return nil, nil, nil, err
	}

	enumeratedTree, err = enumerated.NewTree(contract, accountData.EnumeratedRoot, t.tx)
	if err != nil {
		return nil, nil, nil, err
	}

	enumeratedTotalTree, err = enumeratedtotal.NewTree(contract, accountData.EnumeratedTotalRoot, accountData.TotalSupply, t.tx)
	if err != nil {
		return nil, nil, nil, err
	}

	return ownershipTree, enumeratedTree, enumeratedTotalTree, nil
}

// setTreesForContract sets trees for contract in memory
func (t *tx) setTreesForContract(
	contract common.Address,
	ownershipTree ownership.Tree,
	enumeratedTree enumerated.Tree,
	enumeratedTotalTree enumeratedtotal.Tree,
) {
	slog.Debug("setting trees for contract", "contract", contract.String())

	t.ownershipTrees[contract] = ownershipTree
	t.enumeratedTrees[contract] = enumeratedTree
	t.enumeratedTotalTrees[contract] = enumeratedTotalTree
}

func (t *tx) loadContractStateFromAccountTree(contract common.Address) error {
	slog.Debug("LoadContractState", "contract", contract.String())

	accountData, err := t.accountTree.AccountData(contract)
	if err != nil {
		return err
	}

	enumeratedTree, ok := t.enumeratedTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	enumeratedTotalTree, ok := t.enumeratedTotalTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	ownershipTree, ok := t.ownershipTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	enumeratedTotalTree.SetRoot(accountData.EnumeratedTotalRoot)
	enumeratedTotalTree.SetTotalSupply(accountData.TotalSupply)
	enumeratedTree.SetRoot(accountData.EnumeratedRoot)
	ownershipTree.SetRoot(accountData.OwnershipRoot)

	return nil
}

// TODO check if it can be merged with LoadContractState
// LoadContractTrees loads the merkle trees in memory for contractAddress
func (t *tx) LoadContractTrees(contractAddress common.Address) error {
	slog.Debug("LoadContractTrees", "contract", contractAddress.String())
	if !t.isTreeSetForContract(contractAddress) {
		ownTree, enumTree, enumTotTree, err := t.createTreesForContract(contractAddress)
		if err != nil {
			return err
		}
		t.setTreesForContract(contractAddress, ownTree, enumTree, enumTotTree)
		return nil // when creating new trees we load contract state directly
	}
	return t.loadContractStateFromAccountTree(contractAddress)
}

// OwnerOf returns the owner of the token
func (t *tx) OwnerOf(contract common.Address, tokenId *big.Int) (common.Address, error) {
	slog.Debug("OwnerOf", "contract", contract.String(), "tokenId", tokenId.String())
	ownershipTree, ok := t.ownershipTrees[contract]
	if !ok {
		return common.Address{}, fmt.Errorf("contract %s does not exist", contract.String())
	}
	return ownershipTree.OwnerOf(tokenId)
}

// BalanceOf returns the balance of the owner
func (t *tx) BalanceOf(contract, owner common.Address) (*big.Int, error) {
	slog.Debug("BalanceOf", "contract", contract.String(), "owner", owner.String())
	enumeratedTree, ok := t.enumeratedTrees[contract]
	if !ok {
		return big.NewInt(0), fmt.Errorf("contract %s does not exist", contract.String())
	}

	balance, err := enumeratedTree.BalanceOfOwner(owner)
	if err != nil {
		return big.NewInt(0), err
	}

	return big.NewInt(int64(balance)), nil
}

// TokenOfOwnerByIndex returns the token of the owner by index
func (t *tx) TokenOfOwnerByIndex(contract, owner common.Address, idx int) (*big.Int, error) {
	slog.Debug("TokenOfOwnerByIndex", "contract", contract.String(), "owner", owner.String(), "idx", idx)
	enumeratedTree, ok := t.enumeratedTrees[contract]
	if !ok {
		return big.NewInt(0), fmt.Errorf("contract %s does not exist", contract.String())
	}

	return enumeratedTree.TokenOfOwnerByIndex(owner, uint64(idx))
}

// Transfer transfers ownership of the token. From, To, and TokenID are set in event
func (t *tx) Transfer(contract common.Address, eventTransfer *model.ERC721Transfer) error {
	slog.Debug("Transfer", "contract",
		contract.String(),
		"From", eventTransfer.From.String(),
		"To", eventTransfer.To.String(), "tokenId",
		eventTransfer.TokenId.String())
	ownershipTree, ok := t.ownershipTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	err := ownershipTree.Transfer(eventTransfer)
	if err != nil {
		return err
	}

	tokenData, err := ownershipTree.TokenData(eventTransfer.TokenId)
	if err != nil {
		return err
	}

	if !tokenData.Minted {
		return nil
	}

	enumeratedTree, ok := t.enumeratedTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	err = enumeratedTree.Transfer(true, eventTransfer)
	if err != nil {
		return err
	}

	// if transfer is to zero address (burn) we have to modify the enumeratedTotal tree
	if eventTransfer.To.Cmp(common.Address{}) == 0 {
		enumeratedTotalTree, ok := t.enumeratedTotalTrees[contract]
		if !ok {
			return fmt.Errorf("contract %s does not exist", contract.String())
		}

		tokenIdLast, err := enumeratedTotalTree.TokenByIndex(int(enumeratedTotalTree.TotalSupply()) - 1)
		if err != nil {
			return err
		}

		err = enumeratedTotalTree.Burn(tokenData.Idx)
		if err != nil {
			return err
		}

		tokenDataLast, err := ownershipTree.TokenData(tokenIdLast)
		if err != nil {
			return err
		}

		tokenDataLast.Idx = tokenData.Idx
		return ownershipTree.SetTokenData(tokenDataLast, tokenIdLast)
	}

	return nil
}

// Mint creates a new token
func (t *tx) Mint(contract common.Address, mintEvent *model.MintedWithExternalURI) error {
	slog.Debug("Mint", "contract", contract.String(), "tokenId", mintEvent.TokenId.String())
	enumeratedTotalTree, ok := t.enumeratedTotalTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	err := enumeratedTotalTree.Mint(mintEvent.TokenId)
	if err != nil {
		return err
	}

	ownershipTree, ok := t.ownershipTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	err = ownershipTree.Mint(mintEvent, int(enumeratedTotalTree.TotalSupply())-1)
	if err != nil {
		return err
	}

	tokenData, err := ownershipTree.TokenData(mintEvent.TokenId)
	if err != nil {
		return err
	}

	enumeratedTree, ok := t.enumeratedTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}
	return enumeratedTree.Mint(mintEvent.TokenId, tokenData.SlotOwner)
}

// TotalSupply returns the total number of tokens in the contract
func (t *tx) TotalSupply(contract common.Address) (int64, error) {
	slog.Debug("TotalSupply", "contract", contract.String())
	enumeratedTotalTree, ok := t.enumeratedTotalTrees[contract]
	if !ok {
		return 0, fmt.Errorf("contract %s does not exist", contract.String())
	}

	return enumeratedTotalTree.TotalSupply(), nil
}

// TokenByIndex returns the token at the index
func (t *tx) TokenByIndex(contract common.Address, idx int) (*big.Int, error) {
	slog.Debug("TokenByIndex", "contract", contract.String(), "idx", idx)
	enumeratedTotalTree, ok := t.enumeratedTotalTrees[contract]
	if !ok {
		return big.NewInt(0), fmt.Errorf("contract %s does not exist", contract.String())
	}

	return enumeratedTotalTree.TokenByIndex(idx)
}

// TokenURI returns the token URI associated with tokenId. It returns an error if tokenId does not exist
func (t *tx) TokenURI(contract common.Address, tokenId *big.Int) (string, error) {
	slog.Debug("TokenURI", "contract", contract.String(), "tokenId", tokenId.String())
	ownershipTree, ok := t.ownershipTrees[contract]
	if !ok {
		return "", fmt.Errorf("contract %s does not exist", contract.String())
	}

	tokenData, err := ownershipTree.TokenData(tokenId)
	if err != nil {
		return "", err
	}
	if !tokenData.Minted {
		return "", fmt.Errorf("tokenId %d does not exist", tokenId)
	}
	return tokenData.TokenURI, nil
}

// TagRoot tags roots for all 3 merkle trees at the same block
func (t *tx) TagRoot(blockNumber int64) error {
	slog.Info("TagRoot", "blockNumber", strconv.FormatInt(blockNumber, 10))
	return t.accountTree.TagRoot(blockNumber)
}

func (t *tx) GetLastTaggedBlock() (int64, error) {
	slog.Debug("GetLastTaggedBlock")
	return t.accountTree.GetLastTaggedBlock()
}

func (t *tx) DeleteRootTag(blockNumber int64) error {
	slog.Debug("DeleteRootTag", "blockNumber", strconv.FormatInt(blockNumber, 10))
	return t.accountTree.DeleteRootTag(blockNumber)
}

// DeleteOrphanRootTag deletes the root tags from account merkle tree
// starting from the blockNumber until the most recent tagged block
func (t *tx) DeleteOrphanRootTags(formBlock, toBlock int64) error {
	for blockNumber := formBlock; blockNumber <= toBlock; blockNumber++ {
		if err := t.DeleteRootTag(blockNumber); err != nil {
			return fmt.Errorf("error deleting root tag at block %d: %s", blockNumber, err.Error())
		}
	}
	return nil
}

// Checkout sets the current roots to those tagged for the block
func (t *tx) Checkout(blockNumber int64) error {
	// TODO this transaction should be committed only if we want to permanently store the new root as the head
	// (when reorgs happens)
	// If we just want to read the state at current root we should not commit this transaction
	// probably the easiest and cleanest solution would be to write separate functions for creating transactions
	// NewTransactionForRead and NewTransactionForWrite instead of NewTransaction
	return t.accountTree.Checkout(blockNumber)
}

// UpdateContractState updates the contract state in the account tree
func (t *tx) UpdateContractState(contract common.Address, lastProcessedEvoBlock uint64) error {
	enumeratedTree, ok := t.enumeratedTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	enumeratedTotalTree, ok := t.enumeratedTotalTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	ownershipTree, ok := t.ownershipTrees[contract]
	if !ok {
		return fmt.Errorf("contract %s does not exist", contract.String())
	}

	accountData := account.AccountData{
		EnumeratedRoot:        enumeratedTree.Root(),
		EnumeratedTotalRoot:   enumeratedTotalTree.Root(),
		OwnershipRoot:         ownershipTree.Root(),
		TotalSupply:           enumeratedTotalTree.TotalSupply(),
		LastProcessedEvoBlock: lastProcessedEvoBlock,
	}

	slog.Debug("UpdatingContractState in the account tree", "accountData", accountData)
	return t.accountTree.SetAccountData(&accountData, contract)
}

func (t *tx) AccountData(contract common.Address) (*account.AccountData, error) {
	return t.accountTree.AccountData(contract)
}

// Discards transaction
func (t *tx) Discard() {
	t.tx.Discard()
}

// Commits transaction
func (t *tx) Commit() error {
	return t.tx.Commit()
}
