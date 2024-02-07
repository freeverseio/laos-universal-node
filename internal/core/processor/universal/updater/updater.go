package updater

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"

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
	slog.Debug("UpdateState", "startingBlock", startingBlock, "lastBlockData.Number", lastBlockData.Number)
	blockTimestamps, err := GetBlockTimestampsParallel(ctx, u.client, startingBlock, lastBlockData.Number)
	if err != nil {
		return err
	}

	for block := startingBlock; block <= lastBlockData.Number; block++ {
		blockTime := blockTimestamps[block]
		for _, contract := range contracts {
			if blockWhenDiscovered, ok := newContracts[common.HexToAddress(contract)]; ok {
				if block < blockWhenDiscovered {
					continue
				}
			}
			evoBlock, evoEvents, err := GetEvoEvents(tx, contract, blockTime)
			if err != nil {
				return err
			}
			// Now we update contract storage if there are new events
			if len(evoEvents) > 0 || len(transferEvents[block][contract]) > 0 {
				if err := UpdateContract(tx, contract, evoEvents, transferEvents[block][contract], block, evoBlock); err != nil {
					return err
				}
			}
		}
		if err := tx.TagRoot(int64(block)); err != nil {
			return fmt.Errorf("error occurred while tagging root: %w", err)
		}
	}
	return nil
}

func UpdateContract(tx state.Tx,
	contract string,
	evoEvents []model.MintedWithExternalURI,
	transferEvents []model.ERC721Transfer,
	block uint64,
	evoBlock uint64,
) error {
	if err := tx.LoadContractTrees(common.HexToAddress(contract)); err != nil {
		return fmt.Errorf("error occurred while loading merkle trees for contract %s: %w", contract, err)
	}

	for i := range evoEvents {
		mintEvent := evoEvents[i]
		if err := tx.Mint(common.HexToAddress(contract), &mintEvent); err != nil {
			return fmt.Errorf("error occurred while updating state with mint event %v: %w", mintEvent, err)
		}
	}

	for i := range transferEvents {
		transferEvent := transferEvents[i]
		if err := tx.Transfer(common.HexToAddress(contract), &transferEvent); err != nil {
			return fmt.Errorf("error occurred while updating state with transfer event %v: %w", transferEvent, err)
		}
	}

	return tx.UpdateContractState(common.HexToAddress(contract), evoBlock)
}

func GetEvoEvents(tx state.Tx, contract string, blockTime uint64) (uint64, []model.MintedWithExternalURI, error) {
	collection, err := tx.GetCollectionAddress(contract)
	if err != nil {
		return 0, nil, fmt.Errorf("error occurred retrieving the collection address from the ownership contract %s: %w", contract, err)
	}

	accountData, err := tx.AccountData(common.HexToAddress(contract))
	if err != nil {
		return 0, nil, fmt.Errorf("error occurred retrieving the last processed evo block for ownership contract %s: %w", contract, err)
	}
	evoBlock := accountData.LastProcessedEvoBlock

	evoBlockTimestamp := uint64(0)
	evoEvents := make([]model.MintedWithExternalURI, 0)
	for evoBlockTimestamp < blockTime {
		newBlock, err := tx.GetNextEvoEventBlock(strings.ToLower(collection.String()), evoBlock)
		if err != nil {
			return 0, nil, fmt.Errorf("error occurred retrieving next evo event block for ownership contract %s and evo block %d: %w", contract, evoBlock, err)
		}

		if newBlock == 0 || newBlock == evoBlock {
			break
		}

		mintedEvents, err := tx.GetMintedWithExternalURIEvents(collection.String(), newBlock)
		if err != nil {
			return 0, nil, fmt.Errorf("error occurred retrieving evochain minted events for ownership contract %s and collection address %s: %w",
				contract, collection.String(), err)
		}

		// we get timestamp from event 0 because we don't store evo block separately in db.
		// TODO store evo block data in db when scanning evo chain. recalling block timestamp
		// should be faster then from the function that reads all events
		evoBlockTimestamp = mintedEvents[0].Timestamp
		evoBlock = newBlock
		evoEvents = append(evoEvents, mintedEvents...)
	}
	return evoBlock, evoEvents, nil
}

// GetBlockTimestampsParallel returns a map of block numbers to timestamps in parallel.
// This can increase the speed of this code 2-3x
// TODO we should see if we can get further improvements if instead of making calls in parallel we made one call in
// a batch for all block numbers in a range. possibly this could execute faster because it is only one request
func GetBlockTimestampsParallel(
	ctx context.Context,
	client blockchain.EthClient,
	startingBlock,
	lastBlock uint64,
) (map[uint64]uint64, error) {
	var wg sync.WaitGroup
	timestampsChan := make(chan map[uint64]uint64, lastBlock-startingBlock+1)
	errCh := make(chan error, lastBlock-startingBlock+1)

	for blockNumber := startingBlock; blockNumber <= lastBlock; blockNumber++ {
		wg.Add(1)
		go getBlockTimestamp(ctx, client, blockNumber, &wg, timestampsChan, errCh)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(timestampsChan)
		close(errCh)
	}()

	timestamps := make(map[uint64]uint64)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	// Collect results from the channel
	for result := range timestampsChan {
		for blockNumber, timestamp := range result {
			timestamps[blockNumber] = timestamp
		}
	}

	return timestamps, nil
}

func getBlockTimestamp(ctx context.Context,
	client blockchain.EthClient,
	blockNumber uint64,
	wg *sync.WaitGroup,
	timestamps chan<- map[uint64]uint64,
	errCh chan<- error,
) {
	defer wg.Done()

	header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		errCh <- err
		return
	}

	timestamp := header.Time

	// Send the result through the channel
	timestamps <- map[uint64]uint64{blockNumber: timestamp}
}
