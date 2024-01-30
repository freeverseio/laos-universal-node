package discoverer

import (
	"context"
	"log/slog"
	"math/big"

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
	DiscoverContracts(ctx context.Context, tx state.Tx, startingBlock, lastBlock uint64) (map[common.Address]uint64, error)
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

func (d *discoverer) DiscoverContracts(
	ctx context.Context,
	tx state.Tx,
	startingBlock,
	lastBlock uint64) (map[common.Address]uint64, error) {

	scannedContracts, err := d.scanner.ScanNewUniversalEvents(ctx,
		big.NewInt(int64(startingBlock)),
		big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred while discovering new universal events", "err", err.Error())
		return nil, err
	}

	newContracts := make(map[common.Address]uint64)
	for i := range scannedContracts {
		contract, err := d.validator.Validate(scannedContracts[i])
		if err != nil {
			continue
		}
		if err = tx.StoreERC721UniversalContracts([]model.ERC721UniversalContract{contract}); err != nil {
			slog.Error("error occurred while storing universal contract(s)", "err", err.Error())
			return nil, err
		}

		newContracts[contract.Address] = contract.BlockNumber
		// don't update state when contract is discovered here. it will be updated in the updater.
		// We are passing list of newly discovered contracts to the updater also
	}

	return newContracts, nil
}

func (d *discoverer) GetContracts(tx state.Tx) ([]string, error) {
	if len(d.contracts) > 0 {
		return tx.GetExistingERC721UniversalContracts(d.contracts)
	}
	return tx.GetAllERC721UniversalContracts(), nil
}
