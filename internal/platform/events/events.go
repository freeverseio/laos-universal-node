package events

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
)

type Events interface {
	FilterEventLogs(ctx context.Context, firstBlock, lastBlock *big.Int, topics [][]common.Hash, contracts ...common.Address) ([]types.Log, error)
}

func NewEvents(ownershipChainClient blockchain.EthClient, evoChainClient blockchain.EthClient, contract common.Address) Events {
	return events{
		ownershipChainClient: ownershipChainClient,
		evoChainClient:       evoChainClient,
		contract:             contract,
	}
}

type events struct {
	ownershipChainClient blockchain.EthClient
	evoChainClient       blockchain.EthClient
	contract             common.Address
}

func (s events) FilterEventLogs(ctx context.Context, firstBlock, lastBlock *big.Int, topics [][]common.Hash, contracts ...common.Address) ([]types.Log, error) {
	ownershipLogs, err := filterEventLogs(s.ownershipChainClient, ctx, firstBlock, lastBlock, topics, contracts...)
	if err != nil {
		return nil, err
	}
	slog.Info("ownershipLogs", "ownershipLogs", ownershipLogs)

	evoChainLogs, err := filterEventLogs(s.evoChainClient, ctx, firstBlock, lastBlock, topics, contracts...)
	if err != nil {
		return nil, err
	}
	slog.Info("evoChainLogs", "evoChainLogs", evoChainLogs)

	return mergeEventLogs(ownershipLogs, evoChainLogs), nil
}

func mergeEventLogs(ownershipLogs, evoChainLogs []types.Log) []types.Log {
	return append(ownershipLogs, evoChainLogs...)
}

func filterEventLogs(client blockchain.EthClient, ctx context.Context, firstBlock, lastBlock *big.Int, topics [][]common.Hash, contracts ...common.Address) ([]types.Log, error) {
	return client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: firstBlock,
		ToBlock:   lastBlock,
		Addresses: contracts,
		Topics:    topics,
	})
}
