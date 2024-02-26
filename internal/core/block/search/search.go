package search

import (
	"context"
	"math/big"

	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
)

type blockCorrectionFactor int

const (
	ownershipBlockFactor blockCorrectionFactor = 0
	evoBlockFactor       blockCorrectionFactor = -1
)

func (b blockCorrectionFactor) uint64() uint64 {
	return uint64(b)
}

type Search interface {
	GetOwnershipBlockByTimestamp(ctx context.Context, targetTimestamp uint64, startingPoint ...uint64) (uint64, error)
	GetEvolutionBlockByTimestamp(ctx context.Context, targetTimestamp uint64, startingPoint ...uint64) (uint64, error)
}

type search struct {
	ownershipClient blockchain.EthClient
	evolutionClient blockchain.EthClient
}

func New(ownershipClient, evolutionClient blockchain.EthClient) Search {
	return &search{
		ownershipClient: ownershipClient,
		evolutionClient: evolutionClient,
	}
}

// GetOwnershipBlockByTimestamp performs a binary search in the ownership chain to find the block number for a given timestamp
func (s *search) GetOwnershipBlockByTimestamp(ctx context.Context, targetTimestamp uint64, startingPoint ...uint64) (uint64, error) {
	return s.searchBlockByTimestamp(ctx, targetTimestamp, s.ownershipClient, ownershipBlockFactor, startingPoint...)
}

// GetEvolutionBlockByTimestamp performs a binary search in the evolution chain to find the block number for a given timestamp
func (s *search) GetEvolutionBlockByTimestamp(ctx context.Context, targetTimestamp uint64, startingPoint ...uint64) (uint64, error) {
	return s.searchBlockByTimestamp(ctx, targetTimestamp, s.evolutionClient, evoBlockFactor, startingPoint...)
}

func (s *search) searchBlockByTimestamp(ctx context.Context, targetTimestamp uint64, client blockchain.EthClient, correctionFactor blockCorrectionFactor, startingPoint ...uint64) (uint64, error) {
	var left, right uint64
	if len(startingPoint) > 0 {
		left = startingPoint[0]
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	right = header.Number.Uint64()

	for left <= right {
		mid := left + (right-left)/2
		midHeader, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(mid))
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
