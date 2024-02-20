package blockmapper

import (
	"context"
	"fmt"
	"math/big"

	"github.com/freeverseio/laos-universal-node/internal/core/block/search"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type Processor interface {
	MapNextBlock(ctx context.Context) error
	IsMappingSyncedWithProcessing() (bool, error)
}

type processor struct {
	evoClient    blockchain.EthClient
	blockSearch  search.Search
	stateService state.Service
}

type ProcessorOption func(*processor)

func WithBlockSearch(blockSearch search.Search) ProcessorOption {
	return func(p *processor) {
		p.blockSearch = blockSearch
	}
}

func New(ownershipClient, evoClient blockchain.EthClient, stateService state.Service, options ...ProcessorOption) Processor {
	p := &processor{
		evoClient:    evoClient,
		blockSearch:  search.New(ownershipClient, evoClient),
		stateService: stateService,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

// IsMappingSyncedWithProcessing tells if the last mapped ownership block has reached the last processed ownership block
func (p *processor) IsMappingSyncedWithProcessing() (bool, error) {
	tx, err := p.stateService.NewTransaction()
	if err != nil {
		err = fmt.Errorf("error occurred creating transaction: %w", err)
		return false, err
	}
	defer tx.Discard()

	// check last mapped ownership block from storage
	lastMappedOwnershipBlock, err := tx.GetLastMappedOwnershipBlockNumber()
	if err != nil {
		return false, fmt.Errorf("error occurred retrieving the latest mapped ownership block from storage: %w", err)
	}

	// compare the last mapped ownership block with the last processed ownership block
	lastProcessedOwnershipBlock, err := tx.GetLastOwnershipBlock()
	if err != nil {
		return false, fmt.Errorf("error occurred retrieving the last processed ownership block from storage: %w", err)
	}
	if lastMappedOwnershipBlock >= lastProcessedOwnershipBlock.Number {
		return true, nil
	}
	return false, nil
}

// MapNextBlock retrieves the last mapped evo block number from storage, advances to the next one,
// looks for the corresponding ownership block in time and stores the ownership-evo block pair
func (p *processor) MapNextBlock(ctx context.Context) error {
	tx, err := p.stateService.NewTransaction()
	if err != nil {
		err = fmt.Errorf("error occurred creating transaction: %w", err)
		return err
	}
	defer tx.Discard()

	lastMappedOwnershipBlock, err := tx.GetLastMappedOwnershipBlockNumber()
	if err != nil {
		return fmt.Errorf("error occurred retrieving the latest mapped ownership block from storage: %w", err)
	}

	// get evo block starting point to resume mapping procedure
	evoBlock, err := p.getNextEvoBlockToBeMapped(ctx, lastMappedOwnershipBlock, tx)
	if err != nil {
		return err
	}

	// given the evo block timestamp, find the corresponding ownership block number
	evoHeader, err := p.evoClient.HeaderByNumber(ctx, big.NewInt(int64(evoBlock)))
	if err != nil {
		return fmt.Errorf("error occurred retrieving block number %d from evolution chain %w:", evoBlock, err)
	}
	toMapOwnershipBlock, err := p.blockSearch.GetOwnershipBlockByTimestamp(ctx, evoHeader.Time, lastMappedOwnershipBlock)
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

func (p *processor) getNextEvoBlockToBeMapped(ctx context.Context, lastMappedOwnershipBlock uint64, tx state.Tx) (uint64, error) {
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
		// if no block has ever been mapped, start mapping from the oldest scanned block
		evoBlock, err = p.getOldestScannedBlock(ctx, tx)
		if err != nil {
			return 0, err
		}
	}

	return evoBlock, nil
}

func (p *processor) getOldestScannedBlock(ctx context.Context, tx state.Tx) (uint64, error) {
	ownStartingBlock, err := tx.GetFirstOwnershipBlock()
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving the first ownership block from storage: %w", err)
	}

	evoStartingBlock, err := tx.GetFirstEvoBlock()
	if err != nil {
		return 0, fmt.Errorf("error occurred retrieving the first evolution block from storage: %w", err)
	}
	oldestBlock := evoStartingBlock.Number
	// if the first scanned ownership block was produced before the first scanned evolution block,
	// look for the evolution block corresponding to that ownership block in time
	if ownStartingBlock.Timestamp < evoStartingBlock.Timestamp {
		oldestBlock, err = p.blockSearch.GetEvolutionBlockByTimestamp(ctx, ownStartingBlock.Timestamp)
		if err != nil {
			return 0, fmt.Errorf("error occurred searching for evolution block number by target timestamp %d (ownership block number %d): %w",
				ownStartingBlock.Timestamp, ownStartingBlock.Number, err)
		}
	}
	return oldestBlock, nil
}
