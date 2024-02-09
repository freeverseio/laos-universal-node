package blockmapper

import (
	"context"
	"math/big"
	"time"

	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
)

type Worker interface {
	Run(ctx context.Context) error
	SearchBlockByTimestamp(targetTimestamp int64) (uint64, error)
}

type worker struct {
	waitingTime time.Duration
	client      blockchain.EthClient
}

func New(c *config.Config, client blockchain.EthClient) Worker {
	return &worker{
		waitingTime: c.WaitingTime,
		client:      client,
	}
}

func (w *worker) Run(ctx context.Context) error {
	return nil
}

// SearchBlockByTimestamp performs a binary search to find the block number for a given timestamp.
// It assumes block timestamps are strictly increasing.
func (bs *worker) SearchBlockByTimestamp(targetTimestamp int64) (uint64, error) {
	var (
		left  uint64 = 0
		right uint64
	)

	// Get the latest block to establish the upper bound
	header, err := bs.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	right = header.Number.Uint64()

	for left <= right {
		mid := left + (right-left)/2
		midHeader, err := bs.client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(mid))
		if err != nil {
			return 0, err
		}
		midTimestamp := midHeader.Time

		// Check if midTimestamp matches the target
		if midTimestamp < uint64(targetTimestamp) {
			left = mid + 1
		} else if midTimestamp > uint64(targetTimestamp) {
			right = mid - 1
		} else {
			// Exact match found, return this block number
			return mid, nil
		}
	}

	// Return the closest smaller block number if exact match is not found
	if right < left {
		return right, nil
	}

	return left, nil
}
