package updater

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type Updater interface {
	GetModelTransferEvents(
		ctx context.Context,
		startingBlock,
		lastBlock uint64,
		contracts []string,
	) (map[uint64]map[string][]model.ERC721Transfer, error)

	UpdateState(
		ctx context.Context,
		tx state.Tx,
		contracts []string,
		newContracts map[common.Address]uint64,
		modelTransferEvents map[uint64]map[string][]model.ERC721Transfer,
		startingBlock uint64,
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
) (map[uint64]map[string][]model.ERC721Transfer, error) {
	scanEvents, err := u.scanner.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), contracts)
	if err != nil {
		slog.Error("error occurred while scanning events", "err", err.Error())
		return nil, err
	}

	modelTransferEvents := make(map[uint64]map[string][]model.ERC721Transfer)
	
	for i := range scanEvents {
		if scanEvent, ok := scanEvents[i].(scan.EventTransfer); ok {
			eventTransfer := model.ERC721Transfer{
				From:        scanEvent.From,
				To:          scanEvent.To,
				TokenId:     scanEvent.TokenId,
				BlockNumber: scanEvent.BlockNumber,
				Contract:    scanEvent.Contract,
				Timestamp:   0,
			}
			// timestamp will be updated later to avoid calling headerByNumber for every event.
			// Instead, it will be updated only once for every block
			contractString := strings.ToLower(eventTransfer.Contract.String())
			
			if _, ok := modelTransferEvents[scanEvent.BlockNumber]; !ok {
				modelTransferEvents[scanEvent.BlockNumber] = make(map[string][]model.ERC721Transfer)
			}
			modelTransferEvents[scanEvent.BlockNumber][contractString] = append(modelTransferEvents[scanEvent.BlockNumber][contractString], eventTransfer)
		}
	}
	return modelTransferEvents, nil
}

func (u *updater) UpdateState(
	ctx context.Context,
	tx state.Tx,
	contracts []string,
	newContracts map[common.Address]uint64,
	transferEvents map[uint64]map[string][]model.ERC721Transfer,
	startingBlock uint64,
	lastBlockData model.Block,
) error {

	for block := startingBlock; block <= lastBlockData.Number; block++ {
		slog.Debug("Zoran TEST 1", "block", block)
		header, err := u.client.HeaderByNumber(ctx, big.NewInt(int64(block)))
		if err != nil {
			slog.Debug("error retrieving header for block number", "blockNumber", block, "err", err.Error())
			return err
		}
		blockTime := header.Time
		slog.Debug("Zoran TEST 2", "block", block)
		for _, contract := range contracts {
			if blockWhenDiscovered, ok := newContracts[common.HexToAddress(contract)]; !ok {
				if block<blockWhenDiscovered {
					continue
				}
			}
			collection, err := tx.GetCollectionAddress(contract)
			if err != nil {
				return fmt.Errorf("error occurred retrieving the collection address from the ownership contract %s: %w", contract, err)
			}

			evoBlock, err := tx.GetLastProcessedEvoBlockForOwnershipContract(common.HexToAddress(contract))
			if err != nil {
				return fmt.Errorf("error occurred retrieving the last processed evo block for ownership contract %s: %w", contract, err)
			}
			evoBlockTimestamp := uint64(0)
			evoEvents := make([]model.MintedWithExternalURI, 0)
			for evoBlockTimestamp < blockTime {
				newBlock, err := tx.GetNextEvoEventBlock(strings.ToLower(collection.String()), evoBlock)
				if err != nil {
					return fmt.Errorf("error occurred retrieving next evo event block for ownership contract %s and evo block %d: %w", contract, evoBlock, err)
				}

				slog.Debug("Zoran debug after GetNextEvoBlock", "contract", contract, "newBlock", newBlock, "evoBlock", evoBlock)
				if newBlock == 0 || newBlock == evoBlock {
					break
				}

				mintedEvents, err := tx.GetMintedWithExternalURIEvents(collection.String(), newBlock)
				if err != nil {
					return fmt.Errorf("error occurred retrieving evochain minted events for ownership contract %s and collection address %s: %w",
						contract, collection.String(), err)
				}
				evoBlockTimestamp = mintedEvents[0].Timestamp
				evoBlock = newBlock
				evoEvents = append(evoEvents, mintedEvents...)
			}
			slog.Debug("Zoran TEST 3", "block", block)
			slog.Debug("zoran debug before updating state", "contract", contract, "evoEvents", len(evoEvents), "transferEvents", len(transferEvents[block][contract]))
			// Now we update state if there are new events
			if len(evoEvents) > 0 || len(transferEvents[block][contract]) > 0 {
				err = tx.LoadContractTrees(common.HexToAddress(contract))
				if err != nil {
					slog.Error("error creating merkle trees", "err", err)
					return err
				}

				for _, mintEvent := range evoEvents {
					err = tx.Mint(common.HexToAddress(contract), &mintEvent)
					if err != nil {
						return fmt.Errorf("error occurred while updating state with mint event %v: %w", mintEvent, err)
					}
				}

				for _, transferEvent := range transferEvents[block][contract] {
					err = tx.Transfer(common.HexToAddress(contract), &transferEvent)
					if err != nil {
						return fmt.Errorf("error occurred while updating state with transfer event %v: %w", transferEvent, err)
					}
				}

				// update account state
				err = tx.UpdateContractState(common.HexToAddress(contract))
				if err != nil {
					return fmt.Errorf("error occurred while updating contract state for contract %s: %w", contract, err)
				}

				err = tx.SetLastProcessedEvoBlockForOwnershipContract(common.HexToAddress(contract), evoBlock)
				if err != nil {
					return fmt.Errorf("error occurred while updating current evo block for contract %s: %w", contract, err)
				}
			}
			slog.Debug("Zoran TEST 4", "block", block)
		}
		slog.Debug("Zoran debug before tagging root", "block", block)
		if err := tx.TagRoot(int64(block)); err != nil {
			slog.Error("error occurred while tagging root", "err", err.Error())
			return err
		}
		slog.Debug("Zoran TEST 5", "block", block)
	}
	return nil
}


func getBlockTimestamp(client *ethclient.Client, blockNumber uint64, wg *sync.WaitGroup, timestamps chan<- map[uint64]uint64) {
    defer wg.Done()

    header, err := client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
    if err != nil {
        return
    }

    timestamp := header.Time

    // Send the result through the channel
    timestamps <- map[uint64]uint64{blockNumber: timestamp}
}

func getBlockTimestampsParallel(client *ethclient.Client, startingBlock, lastBlock uint64) map[uint64]uint64 {
    var wg sync.WaitGroup
    timestampsChan := make(chan map[uint64]uint64, lastBlock-startingBlock+1)

    for blockNumber := startingBlock; blockNumber <= lastBlock; blockNumber++ {
        wg.Add(1)
        go getBlockTimestamp(client, blockNumber, &wg, timestampsChan)
    }

    // Close the channel when all goroutines are done
    go func() {
        wg.Wait()
        close(timestampsChan)
    }()

    timestamps := make(map[uint64]uint64)

    // Collect results from the channel
    for result := range timestampsChan {
        for blockNumber, timestamp := range result {
            timestamps[blockNumber] = timestamp
        }
    }

    return timestamps
}