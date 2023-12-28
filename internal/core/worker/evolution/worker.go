package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/state"
)

type reorgError struct {
	block       uint64
	chainHash   common.Hash
	storageHash common.Hash
}

func (e reorgError) Error() string {
	return "reorg error"
}

type Worker interface {
	Run(ctx context.Context) error
}
type worker struct {
	startingBlock uint64
	blocksRange   uint64
	blocksMargin  uint64
	waitingTime   time.Duration
	client        scan.EthClient
	scanner       scan.Scanner
	stateService  state.Service
}

func NewWorker(c *config.Config, client scan.EthClient, scanner scan.Scanner, stateService state.Service) Worker {
	return &worker{
		startingBlock: c.EvoStartingBlock,
		blocksRange:   uint64(c.EvoBlocksRange),
		blocksMargin:  uint64(c.EvoBlocksMargin),
		waitingTime:   c.WaitingTime,
		client:        client,
		scanner:       scanner,
		stateService:  stateService,
	}
}

func (w *worker) Run(ctx context.Context) error {
	slog.Info("starting evolution worker")
	startingBlock, err := getStartingBlock(ctx, w.stateService, w.startingBlock, w.client)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return nil
		default:
			lastBlock, err := getLastBlock(ctx, w.client, startingBlock, w.blocksRange, w.blocksMargin)
			if err != nil {
				slog.Error(err.Error())
				break
			}
			if lastBlock < startingBlock {
				slog.Debug("evolution worker, last calculated block is behind starting block, waiting...",
					"lastBlock", lastBlock, "startingBlock", startingBlock)
				waitBeforeNextScan(ctx, w.waitingTime)
				break
			}

			err = processEvoBlockRange(ctx, w.client, w.stateService, w.scanner, startingBlock, lastBlock)
			if err != nil {
				var reorgErr reorgError
				if errors.As(err, &reorgErr) {
					slog.Error("evolution chain reorganization detected", "block number", reorgErr.block, "chain hash", reorgErr.chainHash.String(), "storage hash", reorgErr.storageHash.String())
					slog.Info("***********************************************************************************************")
					slog.Info("Please wipe out the database before running the node again.")
					slog.Info("***********************************************************************************************")
					return reorgErr
				}
				break
			}

			startingBlock = lastBlock + 1
		}
	}
}

func getLastBlock(ctx context.Context, client scan.EthClient, startingBlock, blocksRange, blocksMargin uint64) (uint64, error) {
	l1LatestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		slog.Error("error retrieving the latest block", "err", err.Error())
		return 0, err
	}

	return min(startingBlock+blocksRange, l1LatestBlock-blocksMargin), nil
}

func processEvoBlockRange(ctx context.Context, client scan.EthClient, stateService state.Service, s scan.Scanner, startingBlock, lastBlock uint64) error {
	tx := stateService.NewTransaction()
	defer tx.Discard()

	// retrieve the hash of the final block of the previous iteration.
	// Retrieve information about the final block in the current block range
	// Store the final block hash to verify in next iteration if a reorganization has taken place.
	err := checkChainReorg(ctx, tx, client, startingBlock, lastBlock)
	if err != nil {
		return err
	}

	events, err := s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), nil)
	if err != nil {
		slog.Error("error occurred while scanning LaosEvolution events", "err", err.Error())
		return err
	}

	err = storeMintedWithExternalURIEventsByContract(tx, events)
	if err != nil {
		slog.Error("error occurred while storing minted events", "err", err.Error())
		return err
	}

	err = updateBlock(ctx, tx, big.NewInt(int64(lastBlock)), client)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("error committing transaction", "err", err.Error())
		return err
	}

	return nil
}

func checkChainReorg(ctx context.Context, tx state.Tx, client scan.EthClient, startingBlock, lastBlock uint64) error {
	prevLastBlockHash, err := tx.GetEvoEndRangeBlockHash()
	if err != nil {
		slog.Error("error occurred while reading LaosEvolution end range block hash", "err", err.Error())
		return err
	}

	err = verifyChainConsistency(ctx, client, prevLastBlockHash, startingBlock)
	if err != nil {
		return err
	}

	endRangeBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred while fetching LaosEvolution end range block", "lastBlock", lastBlock, "err", err.Error())
		return err
	}
	slog.Debug("setting evo end range block hash for block number",
		"blockNumber", endRangeBlock.Number(), "blockHash", endRangeBlock.Hash(), "parentHash", endRangeBlock.ParentHash())

	if err = tx.SetEvoEndRangeBlockHash(endRangeBlock.Hash()); err != nil {
		slog.Error("error occurred while storing LaosEvolution end range block hash", "err", err.Error())
		return err
	}
	return nil
}

func updateBlock(ctx context.Context, tx state.Tx, lastBlock *big.Int, client scan.EthClient) error {
	nextStartingBlock := lastBlock.Uint64() + 1
	if err := tx.SetCurrentEvoBlock(nextStartingBlock); err != nil {
		slog.Error("error occurred while storing current block", "err", err.Error())
		return err
	}

	// asking for timestamp of lastBlock as nextStartingBlock does not exist yet
	timestamp, err := getTimestampForBlockNumber(ctx, client, lastBlock.Uint64())
	if err != nil {
		slog.Error("error retrieving block headers", "err", err.Error())
		return err
	}

	if err = tx.SetCurrentEvoBlockTimestamp(timestamp); err != nil {
		slog.Error("error storing block headers", "err", err.Error())
		return err
	}

	return nil
}

func storeMintedWithExternalURIEventsByContract(tx state.Tx, events []scan.Event) error {
	groupedMintEvents := groupEventsMintedWithExternalURIByContract(events)

	for contract, scannedEvents := range groupedMintEvents {
		// fetch current storedEvents stored for this specific contract address
		storedEvents, err := tx.GetMintedWithExternalURIEvents(contract.String())
		if err != nil {
			return err
		}

		ev := make([]model.MintedWithExternalURI, 0)
		if storedEvents != nil {
			ev = append(ev, storedEvents...)
		}
		ev = append(ev, scannedEvents...)
		if err := tx.StoreMintedWithExternalURIEvents(contract.String(), ev); err != nil {
			return err
		}
	}

	return nil
}

// groups events that are of type scan.EventMintedWithExternalURI by contract address
func groupEventsMintedWithExternalURIByContract(events []scan.Event) map[common.Address][]model.MintedWithExternalURI {
	groupMintEvents := make(map[common.Address][]model.MintedWithExternalURI, 0)
	for _, event := range events {
		if e, ok := event.(scan.EventMintedWithExternalURI); ok {
			groupMintEvents[e.Contract] = append(groupMintEvents[e.Contract], model.MintedWithExternalURI{
				Slot:        e.Slot,
				To:          e.To,
				TokenURI:    e.TokenURI,
				TokenId:     e.TokenId,
				BlockNumber: e.BlockNumber,
				Timestamp:   e.Timestamp,
			})
		}
	}
	return groupMintEvents
}

func getStartingBlock(ctx context.Context, stateService state.Service, configStartingBlock uint64, client scan.EthClient) (uint64, error) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	startingBlockDB, err := tx.GetCurrentEvoBlock()
	if err != nil {
		return 0, fmt.Errorf("error retrieving the current block from storage: %w", err)
	}

	var startingBlock uint64
	if startingBlockDB != 0 {
		startingBlock = startingBlockDB
		slog.Debug("ignoring user provided starting block, using last updated block from storage", "starting_block", startingBlock)
	}

	if startingBlock == 0 {
		startingBlock = configStartingBlock
		if startingBlock == 0 {
			startingBlock, err = client.BlockNumber(ctx)
			if err != nil {
				return 0, fmt.Errorf("error retrieving the latest block from chain: %w", err)
			}
			slog.Debug("latest block found", "latest_block", startingBlock)
		}
	}
	return startingBlock, nil
}

func getTimestampForBlockNumber(ctx context.Context, client scan.EthClient, blockNumber uint64) (uint64, error) {
	header, err := client.HeaderByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return 0, err
	}
	return header.Time, err
}

func waitBeforeNextScan(ctx context.Context, waitingTime time.Duration) {
	timer := time.NewTimer(waitingTime)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}

func verifyChainConsistency(ctx context.Context, client scan.EthClient, prevLastBlockHash common.Hash, startingBlock uint64) error {
	// During the initial iteration, no hash is stored in the database, so this code block is bypassed.
	// Verify whether the hash of the last block from the previous iteration remains unchanged; if it differs,
	// it indicates a reorganization has taken place.
	if prevLastBlockHash == (common.Hash{}) {
		return nil
	}

	prevIterLastBlockNumber := startingBlock - 1
	slog.Debug("verifying chain consistency on block number", "lastBlock", prevIterLastBlockNumber)
	prevIterLastBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(prevIterLastBlockNumber)))
	if err != nil {
		slog.Error("error occurred while retrieving new start range block", "err", err.Error())
		return err
	}

	if prevIterLastBlock.Hash().Cmp(prevLastBlockHash) != 0 {
		return reorgError{block: startingBlock - 1, chainHash: prevIterLastBlock.Hash(), storageHash: prevLastBlockHash}
	}

	return nil
}
