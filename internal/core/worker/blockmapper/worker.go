package blockmapper

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	sharedProcessor "github.com/freeverseio/laos-universal-node/internal/core/processor"
	shared "github.com/freeverseio/laos-universal-node/internal/core/worker"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type BlockCorrectionFactor int

const (
	OwershipBlockFactor BlockCorrectionFactor = 0
	EvoBlockFactor      BlockCorrectionFactor = -1
)

func (b BlockCorrectionFactor) uint64() uint64 {
	return uint64(b)
}

type Worker interface {
	Run(ctx context.Context)
	SearchBlockByTimestamp(targetTimestamp uint64, client blockchain.EthClient, correctionFactor BlockCorrectionFactor) (uint64, error)
}

type worker struct {
	waitingTime          time.Duration
	ownershipBlockHelper *sharedProcessor.BlockHelper
	evoBlockHelper       *sharedProcessor.BlockHelper
	clientOwnership      blockchain.EthClient
	clientEvo            blockchain.EthClient
	stateService         state.Service
}

func New(c *config.Config, ownershipClient, evoClient blockchain.EthClient, stateService state.Service) Worker {
	return &worker{
		waitingTime:     c.WaitingTime,
		clientOwnership: ownershipClient,
		clientEvo:       evoClient,
		ownershipBlockHelper: sharedProcessor.NewBlockHelper(
			ownershipClient,
			stateService,
			uint64(c.BlocksRange),
			uint64(c.BlocksMargin),
			c.StartingBlock,
		),
		evoBlockHelper: sharedProcessor.NewBlockHelper(
			evoClient,
			stateService,
			uint64(c.EvoBlocksRange),
			uint64(c.EvoBlocksMargin),
			c.EvoStartingBlock,
		),
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
			tx, err := w.stateService.NewTransaction()
			if err != nil {
				slog.Error("error occurred creating transaction", "err", err)
				break
			}
			// TODO tx must be discarded outside the infinite loop
			defer tx.Discard()
			// check last mapped own block from db
			lastMappedOwnershipBlock, err := tx.GetLastMappedOwnershipBlockNumber()
			if err != nil {
				slog.Error("error occurred retrieving the latest mapped ownership block from storage", "err", err)
				break
			}
			var evoBlock uint64
			// if no block has ever been mapped, start mapping from the oldest user-defined block
			if lastMappedOwnershipBlock == 0 {
				var ownershipHeader, evoHeader *types.Header
				var ownershipStartingBlock, evoStartingBlock uint64
				ownershipStartingBlock, err = w.ownershipBlockHelper.GetOwnershipInitStartingBlock(ctx)
				if err != nil {
					slog.Error("error occurred retrieving the ownership init starting block", "err", err)
					break
				}
				evoStartingBlock, err = w.evoBlockHelper.GetEvoInitStartingBlock(ctx)
				if err != nil {
					slog.Error("error occurred retrieving the evolution init starting block", "err", err)
					break
				}
				ownershipHeader, err = w.clientOwnership.HeaderByNumber(ctx, big.NewInt(int64(ownershipStartingBlock)))
				if err != nil {
					slog.Error("error occurred retrieving block number from ownership chain",
						"blockNumber", ownershipStartingBlock, "err", err)
					break
				}
				evoHeader, err = w.clientEvo.HeaderByNumber(ctx, big.NewInt(int64(evoStartingBlock)))
				if err != nil {
					slog.Error("error occurred retrieving block number from evolution chain",
						"blockNumber", evoStartingBlock, "err", err)
					break
				}
				evoBlock = evoStartingBlock
				if ownershipHeader.Time < evoHeader.Time {
					evoBlock, err = w.SearchBlockByTimestamp(ownershipHeader.Time, w.clientEvo, EvoBlockFactor)
					if err != nil {
						slog.Error("error occurred searching block number by target timestamp in evolution chain",
							"targetTimestamp", ownershipHeader.Time, "ownershipBlockNumber", ownershipStartingBlock, "err", err)
						break
					}
				}
			} else {
				// if a block has been mapped, start mapping from the next one
				evoBlock, err = tx.GetMappedOwnershipBlockNumber(lastMappedOwnershipBlock)
				if err != nil {
					slog.Error("error occurred retrieving the mapped evo block number from the ownership block",
						"ownershipBlock", lastMappedOwnershipBlock, "err", err)
					break
				}
				evoBlock++
			}
			dbLastOwnershipBlock, err := tx.GetLastOwnershipBlock()
			if err != nil {
				slog.Error("error occurred retrieving the last processed ownership block from storage", "err", err)
				break
			}
			// compare that block with the last processed ownership block
			if lastMappedOwnershipBlock >= dbLastOwnershipBlock.Number {
				shared.Wait(ctx, w.waitingTime)
				slog.Debug("mapped block has reached the last processed ownership block, waiting to process more blocks before mapping again...")
				continue
			}

			// given the evo block timestamp, find the corresponding ownership block number
			evoHeader, err := w.clientEvo.HeaderByNumber(ctx, big.NewInt(int64(evoBlock)))
			if err != nil {
				slog.Error("error occurred retrieving block number from evolution chain",
					"blockNumber", evoBlock, "err", err)
				break
			}
			toMapOwnershipBlock, err := w.SearchBlockByTimestamp(evoHeader.Time, w.clientOwnership, OwershipBlockFactor)
			if err != nil {
				slog.Error("error occurred searching block number by target timestamp in ownership chain",
					"targetTimestamp", evoHeader.Time, "evoBlockNumber", evoBlock, "err", err)
				break
			}
			// set ownership block -> evo block mapping
			err = tx.SetMappedOwnershipBlockNumber(toMapOwnershipBlock, evoBlock)
			if err != nil {
				slog.Error("error setting ownership block number to evo block number in storage",
					"ownershipBlockNumber", toMapOwnershipBlock, "evoBlockNumber", evoBlock, "err", err)
				break
			}
			err = tx.SetLastMappedOwnershipBlockNumber(toMapOwnershipBlock)
			if err != nil {
				slog.Error("error setting the latest mapped ownership block number in storage",
					"ownershipBlockNumber", toMapOwnershipBlock, "err", err)
				break
			}
			err = tx.Commit()
			if err != nil {
				slog.Error("error committing transaction", "err", err)
			}
		}
	}
}

// SearchBlockByTimestamp performs a binary search to find the block number for a given timestamp.
// It assumes block timestamps are strictly increasing.
func (w *worker) SearchBlockByTimestamp(targetTimestamp uint64, client blockchain.EthClient, correctionFactor BlockCorrectionFactor) (uint64, error) {
	var (
		left  uint64 = 0
		right uint64
	)

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	right = header.Number.Uint64()

	for left <= right {
		mid := left + (right-left)/2
		midHeader, err := client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(mid))
		if err != nil {
			return 0, err
		}
		midTimestamp := midHeader.Time
		switch {
		case midTimestamp < targetTimestamp:
			left = mid + 1
		case midTimestamp > targetTimestamp:
			right = mid - 1
		default:
			return mid + correctionFactor.uint64(), nil
		}
	}

	return left + correctionFactor.uint64(), nil
}
