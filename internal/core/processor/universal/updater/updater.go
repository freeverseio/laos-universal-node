package updater

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

const historyLength = 256

type Updater interface {
	GetModelTransferEvents(
		ctx context.Context,
		startingBlock,
		lastBlock uint64,
		contracts []string,
	) (map[string][]model.ERC721Transfer, error)

	UpdateState(
		ctx context.Context,
		tx state.Tx,
		contracts []string,
		modelTransferEvents map[string][]model.ERC721Transfer,
		lastBlockData model.Block,
	) error
}

type updater struct {
	client  blockchain.EthClient
	scanner scan.Scanner
}

func New(client blockchain.EthClient, scanner scan.Scanner) Updater {
	return &updater{
		client:  client,
		scanner: scanner,
	}
}

func (u *updater) GetModelTransferEvents(
	ctx context.Context,
	startingBlock,
	lastBlock uint64,
	contracts []string,
) (map[string][]model.ERC721Transfer, error) {
	scanEvents, err := u.scanner.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), contracts)
	if err != nil {
		slog.Error("error occurred while scanning events", "err", err.Error())
		return nil, err
	}
	modelTransferEvents := make(map[string][]model.ERC721Transfer)
	for i := range scanEvents {
		if scanEvent, ok := scanEvents[i].(scan.EventTransfer); ok {
			header, err := u.client.HeaderByNumber(ctx, big.NewInt(int64(scanEvent.BlockNumber)))
			if err != nil {
				return nil, fmt.Errorf("error retrieving timestamp for block number %d: %w", scanEvent.BlockNumber, err)
			}
			contractString := strings.ToLower(scanEvent.Contract.String())
			modelTransferEvents[contractString] = append(modelTransferEvents[contractString], model.ERC721Transfer{
				From:        scanEvent.From,
				To:          scanEvent.To,
				TokenId:     scanEvent.TokenId,
				BlockNumber: scanEvent.BlockNumber,
				Contract:    scanEvent.Contract,
				Timestamp:   header.Time,
			})
		}
	}
	return modelTransferEvents, nil
}

func (u *updater) UpdateState(
	ctx context.Context,
	tx state.Tx,
	contracts []string,
	modelTransferEvents map[string][]model.ERC721Transfer,
	lastBlockData model.Block,
) error {
	for _, contract := range contracts {
		err := loadMerkleTree(tx, common.HexToAddress(contract))
		if err != nil {
			slog.Error("error creating merkle trees", "err", err)
			return err
		}

		collection, err := tx.GetCollectionAddress(contract)
		if err != nil {
			return fmt.Errorf("error occurred retrieving the collection address from the ownership contract %s: %w", contract, err)
		}

		// now we get the minted events from the evolution chain for the collection address
		mintedEvents, err := tx.GetMintedWithExternalURIEvents(collection.String())
		if err != nil {
			return fmt.Errorf("error occurred retrieving evochain minted events for ownership contract %s and collection address %s: %w",
				contract, collection.String(), err)
		}

		evoIndex, err := mergeAndUpdate(ctx, u.client, mintedEvents, modelTransferEvents[contract], contract, tx, lastBlockData.Timestamp)
		if err != nil {
			return fmt.Errorf("error updating state: %w", err)
		}

		if err = tx.SetCurrentEvoEventsIndexForOwnershipContract(contract, evoIndex); err != nil {
			return fmt.Errorf("error updating current evochain index %d for ownership contract %s: %w", evoIndex, contract, err)
		}

		if err = tagRootsUntilBlock(tx, contract, lastBlockData.Number+1); err != nil { // add+1 because we want to tag last block also
			slog.Error("error occurred while tagging roots", "err", err.Error())
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

func mergeAndUpdate(ctx context.Context, client blockchain.EthClient, mintedEvents []model.MintedWithExternalURI, modelTransferEvents []model.ERC721Transfer, contract string, tx state.Tx, lastBlockTimestamp uint64) (uint64, error) {
	ownershipContractEvoEventIndex, err := tx.GetCurrentEvoEventsIndexForOwnershipContract(contract)
	if err != nil {
		return 0, err
	}
	var transferIndex int
	for {
		switch {
		// all events have been processed => return
		case ownershipContractEvoEventIndex >= uint64(len(mintedEvents)) && transferIndex >= len(modelTransferEvents):
			return ownershipContractEvoEventIndex, nil
		// all minted events have been processed => process remaining transfer events
		case ownershipContractEvoEventIndex >= uint64(len(mintedEvents)):
			if err := updateStateWithTransfer(contract, tx, &modelTransferEvents[transferIndex]); err != nil {
				return 0, err
			}
			transferIndex++
		// all transfer events have been processed => process remaining minted events
		case transferIndex >= len(modelTransferEvents):
			if mintedEvents[ownershipContractEvoEventIndex].Timestamp < lastBlockTimestamp {
				err := updateStateWithMint(ctx, client, contract, tx, &mintedEvents[ownershipContractEvoEventIndex])
				if err != nil {
					return 0, err
				}
				ownershipContractEvoEventIndex++
			} else {
				return ownershipContractEvoEventIndex, nil
			}

		default:
			// if minted event's timestamp is behind transfer event's timestamp => process minted event
			if mintedEvents[ownershipContractEvoEventIndex].Timestamp < modelTransferEvents[transferIndex].Timestamp {
				err := updateStateWithMint(ctx, client, contract, tx, &mintedEvents[ownershipContractEvoEventIndex])
				if err != nil {
					return 0, err
				}
				ownershipContractEvoEventIndex++
			} else {
				if err := updateStateWithTransfer(contract, tx, &modelTransferEvents[transferIndex]); err != nil {
					return 0, err
				}
				transferIndex++
			}
		}
	}
}

func updateStateWithTransfer(contract string, tx state.Tx, transferEvent *model.ERC721Transfer) error {
	slog.Debug("updating state with transfer event", "modelTransferEvent", transferEvent, "contract", contract)
	err := tagRootsUntilBlock(tx, contract, transferEvent.BlockNumber)
	if err != nil {
		return err
	}

	return tx.Transfer(common.HexToAddress(contract), transferEvent)
}

func updateStateWithMint(ctx context.Context, client blockchain.EthClient, contract string, tx state.Tx, mintedEvent *model.MintedWithExternalURI) error {
	slog.Debug("updating state with mint event", "modelTransferEvent", mintedEvent, "contract", contract)
	block, err := getFirstOwnershipBlockAfterMintEvent(ctx, client, contract, tx, mintedEvent)
	if err != nil {
		return err
	}

	err = tagRootsUntilBlock(tx, contract, block)
	if err != nil {
		return err
	}

	return tx.Mint(common.HexToAddress(contract), mintedEvent)
}

func getFirstOwnershipBlockAfterMintEvent(ctx context.Context,
	client blockchain.EthClient,
	contract string,
	tx state.Tx,
	mintedEvent *model.MintedWithExternalURI,
) (uint64, error) {
	// At this point this function uses GetLastTaggedBlock(). I think the first ownership block after mint event
	// can be found in a way that is more clean but I will keep it now as is

	slog.Debug("finding the first ownership block after mint event", "mintedEvent", mintedEvent, "contract", contract)

	block, err := tx.GetLastTaggedBlock(common.HexToAddress(contract))
	if err != nil {
		return 0, err
	}

	for {
		block++
		header, err := client.HeaderByNumber(ctx, big.NewInt(block))
		if err != nil {
			return 0, err
		}
		if header.Time >= mintedEvent.Timestamp {
			return uint64(block), nil
		}
	}
}

func tagRootsUntilBlock(tx state.Tx, contractAddress string, blockNumber uint64) error {
	lastTaggedBlock, err := tx.GetLastTaggedBlock(common.HexToAddress(contractAddress))
	if err != nil {
		return err
	}

	for block := lastTaggedBlock + 1; block < int64(blockNumber); block++ {
		if err := tx.TagRoot(common.HexToAddress(contractAddress), block); err != nil {
			return err
		}
		if err := tx.DeleteRootTag(common.HexToAddress(contractAddress), block-historyLength); err != nil {
			return err
		}
	}
	return nil
}
