package blockmapper

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/block/helper"
	"github.com/freeverseio/laos-universal-node/internal/core/block/search"
	shared "github.com/freeverseio/laos-universal-node/internal/core/worker"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type Worker interface {
	Run(ctx context.Context)
}

type worker struct {
	waitingTime          time.Duration
	ownershipBlockHelper *helper.Helper
	evoBlockHelper       *helper.Helper
	ownershipClient      blockchain.EthClient
	evoClient            blockchain.EthClient
	blockSearch          search.Search
	stateService         state.Service
}

func New(c *config.Config, ownershipClient, evoClient blockchain.EthClient, stateService state.Service) Worker {
	return &worker{
		waitingTime:     c.WaitingTime,
		ownershipClient: ownershipClient,
		evoClient:       evoClient,
		ownershipBlockHelper: helper.New(
			ownershipClient,
			stateService,
			uint64(c.BlocksRange),
			uint64(c.BlocksMargin),
			c.StartingBlock,
		),
		evoBlockHelper: helper.New(
			evoClient,
			stateService,
			uint64(c.EvoBlocksRange),
			uint64(c.EvoBlocksMargin),
			c.EvoStartingBlock,
		),
		blockSearch:  search.New(ownershipClient, evoClient),
		stateService: stateService,
	}
}

func (w *worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return
		default:
			if err := w.ExecuteMapping(ctx); err != nil {
				slog.Error("error occurred while performing block mapping", "err", err)
			}
		}
	}
}

func (w *worker) ExecuteMapping(ctx context.Context) error {
	tx, err := w.stateService.NewTransaction()
	if err != nil {
		err = fmt.Errorf("error occurred creating transaction: %w", err)
		return err
	}
	defer tx.Discard()

	// check last mapped ownership block from storage
	lastMappedOwnershipBlock, err := tx.GetLastMappedOwnershipBlockNumber()
	if err != nil {
		return fmt.Errorf("error occurred retrieving the latest mapped ownership block from storage: %w", err)
	}

	// get evo block starting point to resume mapping procedure
	evoBlock, err := w.getInitialEvoBlock(ctx, lastMappedOwnershipBlock, tx)
	if err != nil {
		return err
	}

	// compare the last mapped ownership block with the last processed ownership block
	lastProcessedOwnershipBlock, err := tx.GetLastOwnershipBlock()
	if err != nil {
		return fmt.Errorf("error occurred retrieving the last processed ownership block from storage: %w", err)
	}
	if lastMappedOwnershipBlock >= lastProcessedOwnershipBlock.Number {
		slog.Debug("mapped block has reached the last processed ownership block, waiting to process more blocks before mapping again...")
		shared.Wait(ctx, w.waitingTime)
		return nil
	}

	// given the evo block timestamp, find the corresponding ownership block number
	evoHeader, err := w.evoClient.HeaderByNumber(ctx, big.NewInt(int64(evoBlock)))
	if err != nil {
		return fmt.Errorf("error occurred retrieving block number %d from evolution chain %w:", evoBlock, err)
	}
	toMapOwnershipBlock, err := w.blockSearch.GetOwnershipBlockByTimestamp(ctx, evoHeader.Time, lastMappedOwnershipBlock)
	if err != nil {
		return fmt.Errorf("error occurred searching for ownership block number by target timestamp %d (evolution block number %d): %w",
			evoHeader.Time, evoBlock, err)
	}

	// set ownership block -> evo block mapping
	err = tx.SetOwnershipEvoBlockMapping(toMapOwnershipBlock, evoBlock)
	if err != nil {
		return fmt.Errorf("error setting ownership block number %d (key) to evo block number %d (value) in storage: %w",
			toMapOwnershipBlock, evoBlock, err)
	}
	err = tx.SetLastMappedOwnershipBlockNumber(toMapOwnershipBlock)
	if err != nil {
		return fmt.Errorf("error setting the last mapped ownership block number %d in storage: %w", toMapOwnershipBlock, err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

func (w *worker) getInitialEvoBlock(ctx context.Context, lastMappedOwnershipBlock uint64, tx state.Tx) (uint64, error) {
	var evoBlock uint64
	var err error
	// if a mapped block is found in storage, resume mapping from the next one
	if lastMappedOwnershipBlock > 0 {
		evoBlock, err = tx.GetMappedEvoBlockNumber(lastMappedOwnershipBlock)
		if err != nil {
			return 0, fmt.Errorf("error occurred retrieving the mapped evolution block number by ownership block %d from storage: %w",
				lastMappedOwnershipBlock, err)
		}
		evoBlock++
	} else {
		// if no block has ever been mapped, start mapping from the oldest user-defined block
		evoBlock, err = w.getOldestUserDefinedBlock(ctx)
		if err != nil {
			return 0, err
		}
	}

	return evoBlock, nil
}

func (w *worker) getOldestUserDefinedBlock(ctx context.Context) (uint64, error) {
	ownershipStartingBlock, err := w.ownershipBlockHelper.GetOwnershipInitStartingBlock(ctx)
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving the ownership init starting block: %w", err)
	}
	evoStartingBlock, err := w.evoBlockHelper.GetEvoInitStartingBlock(ctx)
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving the evolution init starting block: %w", err)
	}
	ownershipHeader, err := w.ownershipClient.HeaderByNumber(ctx, big.NewInt(int64(ownershipStartingBlock)))
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving block number %d from ownership chain: %w",
			ownershipStartingBlock, err)
	}
	evoHeader, err := w.evoClient.HeaderByNumber(ctx, big.NewInt(int64(evoStartingBlock)))
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving block number %d from evolution chain: %w",
			evoStartingBlock, err)
	}
	oldestBlock := evoStartingBlock
	// if the user-defined ownership block was produced before the user-defined evolution block,
	// look for the evolution block corresponding to that ownership block in time
	if ownershipHeader.Time < evoHeader.Time {
		oldestBlock, err = w.blockSearch.GetEvolutionBlockByTimestamp(ctx, ownershipHeader.Time)
		if err != nil {
			return 0, fmt.Errorf("error occurred searching for evolution block number by target timestamp %d (ownership block number %d): %w",
				ownershipHeader.Time, ownershipStartingBlock, err)
		}
	}
	return oldestBlock, nil
}
