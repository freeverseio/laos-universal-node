package blockmapper

import (
	"context"
	"math/big"
	"time"

	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
)

type BlockCorrectionFactor int

const (
	OwershipBlockFactor BlockCorrectionFactor = 1
	EvoChainBlockFactor BlockCorrectionFactor = -1
)

func (b BlockCorrectionFactor) UiInt64() uint64 {
	return uint64(int(b))
}

type correctionFuncType func(uint64, bool) uint64

var OwershipBlockCorrectionFunc correctionFuncType = func(blockNumber uint64, sameTimestamp bool) uint64 {
	return blockNumber + OwershipBlockFactor.UiInt64()
}

var EvoChainBlockCorrectionFunc correctionFuncType = func(blockNumber uint64, sameTimestamp bool) uint64 {
	if sameTimestamp {
		return blockNumber
	}
	return blockNumber + EvoChainBlockFactor.UiInt64()
}

type Worker interface {
	Run(ctx context.Context) error
	SearchBlockByTimestamp(targetTimestamp int64, client blockchain.EthClient, correctionFunction func(uint64, bool) uint64) (uint64, error)
}

type worker struct {
	waitingTime     time.Duration
	clientOwnership blockchain.EthClient
	clientEvoChain  blockchain.EthClient
}

func New(c *config.Config, client blockchain.EthClient) Worker {
	return &worker{
		waitingTime:     c.WaitingTime,
		clientOwnership: client,
	}
}

func (w *worker) Run(ctx context.Context) error {

	// check last mapped own block from db

	// if no block is mapped, start from the beginning
	// compare own starting block with evo starting block
	// we take the oldest and start mapping from there

	// find the first evo block (from the oldest block timestamp)
	// -> find own block number for the evo block timestamp
	// -> map the key own block to value evo block
	// evo block ++
	// ...

	return nil
}

// SearchBlockByTimestamp performs a binary search to find the block number for a given timestamp.
// It assumes block timestamps are strictly increasing.
func (bs *worker) SearchBlockByTimestamp(targetTimestamp int64, client blockchain.EthClient, correctionFunc func(uint64, bool) uint64) (uint64, error) {
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
		case midTimestamp < uint64(targetTimestamp):
			left = mid + 1
		case midTimestamp > uint64(targetTimestamp):
			right = mid - 1
		default:
			return correctionFunc(mid, true), nil // Apply the correction function here
		}
	}

	return correctionFunc(right, false), nil // Apply the correction function here
}
