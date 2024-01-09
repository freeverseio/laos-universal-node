package discoverer

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	uValidator "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer/validator"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type Discoverer interface {
	ShouldDiscover(tx state.Tx, startingBlock, lastBlock uint64) (bool, error)
	GetContracts(tx state.Tx) ([]string, error)
	DiscoverContracts(ctx context.Context, tx state.Tx, startingBlock, lastBlock uint64) error
}

type discoverer struct {
	client    blockchain.EthClient
	contracts []string
	scanner   scan.Scanner
	validator uValidator.Validator
}

func New(
	client blockchain.EthClient,
	contracts []string,
	scanner scan.Scanner,
	validator uValidator.Validator,
) Discoverer {
	return &discoverer{
		client:    client,
		contracts: contracts,
		scanner:   scanner,
		validator: validator,
	}
}

func (d *discoverer) ShouldDiscover(tx state.Tx, startingBlock, lastBlock uint64) (bool, error) {
	if len(d.contracts) == 0 {
		return true, nil
	}
	for i := 0; i < len(d.contracts); i++ {
		hasContract, err := tx.HasERC721UniversalContract(d.contracts[i])
		if err != nil {
			return false, err
		}
		if !hasContract {
			return true, nil
		}
	}
	return false, nil
}

func (d *discoverer) DiscoverContracts(ctx context.Context, tx state.Tx, startingBlock, lastBlock uint64) error {
	scannedContracts, err := d.scanner.ScanNewUniversalEvents(ctx,
		big.NewInt(int64(startingBlock)),
		big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred while discovering new universal events", "err", err.Error())
		return err
	}

	for i := range scannedContracts {
		contract, err := d.validator.Validate(scannedContracts[i])
		if err != nil {
			continue
		}
		if err = tx.StoreERC721UniversalContracts([]model.ERC721UniversalContract{contract}); err != nil {
			slog.Error("error occurred while storing universal contract(s)", "err", err.Error())
			return err
		}

		if err = loadMerkleTree(tx, contract.Address); err != nil {
			slog.Error("error creating merkle trees for newly discovered universal contract(s)", "err", err)
			return err
		}

		// check if there are mint events for this contract
		mintEvents, err := tx.GetMintedWithExternalURIEvents(contract.CollectionAddress.String())
		if err != nil {
			slog.Error("error occurred retrieving evochain minted events for ownership contract: %w", err)
			return err
		}

		header, err := d.client.HeaderByNumber(ctx, big.NewInt(int64(contract.BlockNumber)))
		if err != nil {
			return err
		}

		ownershipContractEvoEventIndex, err := updateStateWithMintEvents(contract.Address, tx, mintEvents, header.Time)
		if err != nil {
			slog.Error("error occurred updating state with mint events", "err", err)
			return err
		}

		if err = tx.SetCurrentEvoEventsIndexForOwnershipContract(contract.Address.String(), ownershipContractEvoEventIndex); err != nil {
			return fmt.Errorf("error updating current evochain event index %d for ownership contract %s: %w",
				ownershipContractEvoEventIndex, strings.ToLower(contract.Address.String()), err)
		}

		if err = tx.TagRoot(contract.Address, int64(contract.BlockNumber)); err != nil {
			slog.Error("error occurred tagging roots for newly discovered universal contract(s)", "err", err.Error())
			return err
		}
	}

	return nil
}

func loadMerkleTree(tx state.Tx, contractAddress common.Address) error {
	if !tx.IsTreeSetForContract(contractAddress) {
		ownership, enumerated, enumeratedTotal, err := tx.CreateTreesForContract(contractAddress)
		if err != nil {
			return err
		}
		tx.SetTreesForContract(contractAddress, ownership, enumerated, enumeratedTotal)
	}
	return nil
}

func updateStateWithMintEvents(
	contract common.Address,
	tx state.Tx,
	mintedEvents []model.MintedWithExternalURI,
	timestampContract uint64,
) (uint64, error) {
	for i := range mintedEvents {
		if mintedEvents[i].Timestamp > timestampContract {
			return uint64(i), nil
		}
		if err := tx.Mint(contract, &mintedEvents[i]); err != nil {
			return 0, fmt.Errorf("error updating mint state for contract %s and token id %d: %w",
				contract, mintedEvents[i].TokenId, err)
		}
	}
	return uint64(len(mintedEvents)), nil
}

func (d *discoverer) GetContracts(tx state.Tx) ([]string, error) {
	if len(d.contracts) > 0 {
		return tx.GetExistingERC721UniversalContracts(d.contracts)
	}
	return tx.GetAllERC721UniversalContracts(), nil
}
