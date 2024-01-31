package events

import (
	"context"
	"log/slog"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

type LogType int

const (
	OwnershipLog LogType = iota
	EvoLog       LogType = iota
)

type UnodeLog struct {
	types.Log
	OriginalBlockNumber uint64
	LogType             LogType
}

type Events interface {
	FilterEventLogs(ctx context.Context, firstBlock, lastBlock *big.Int, topics [][]common.Hash, contracts ...common.Address) ([]types.Log, error)
}

func NewEvents(ownershipChainClient, evoChainClient blockchain.EthClient, stateService state.Service, contract common.Address) Events {
	return events{
		ownershipChainClient: ownershipChainClient,
		evoChainClient:       evoChainClient,
		stateService:         stateService,
		contract:             contract,
	}
}

type events struct {
	ownershipChainClient blockchain.EthClient
	evoChainClient       blockchain.EthClient
	stateService         state.Service
	contract             common.Address
}

func (s events) FilterEventLogs(ctx context.Context, firstBlock, lastBlock *big.Int, topics [][]common.Hash, contracts ...common.Address) ([]types.Log, error) {
	ownershipLogs, err := filterEventLogs(s.ownershipChainClient, ctx, firstBlock, lastBlock, topics, contracts...)
	if err != nil {
		return nil, err
	}

	tx := s.stateService.NewTransaction()
	firstBlockEvo, err := tx.GetCorrespondingEvoBlockNumber(uint64(firstBlock.Int64()))
	if err != nil {
		slog.Error("error getting corresponding evo block number", "err", err)
		return nil, err
	}
	lastBlockEvo, err := tx.GetCorrespondingEvoBlockNumber(uint64(lastBlock.Int64()))
	if err != nil {
		slog.Error("error getting corresponding evo block number", "err", err)
		return nil, err
	}

	evoChainLogs, err := filterEventLogs(s.evoChainClient, ctx, big.NewInt(int64(firstBlockEvo)), big.NewInt(int64(lastBlockEvo)), topics, contracts...)
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

func SortLogs(logs []UnodeLog) []UnodeLog {
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].LogType != logs[j].LogType {
			return logs[i].LogType < logs[j].LogType
		}
		if logs[i].BlockNumber == logs[j].BlockNumber {
			return logs[i].TxIndex < logs[j].TxIndex
		}
		return logs[i].BlockNumber < logs[j].BlockNumber
	})
	return logs
}

func convertToUnodeLogs(logs []types.Log, logType LogType) []UnodeLog {
	unodeLogs := make([]UnodeLog, len(logs))
	for i, log := range logs {
		if logType == EvoLog {
			//
		}
		unodeLogs[i] = UnodeLog{
			Log:                 log,
			OriginalBlockNumber: log.BlockNumber,
			LogType:             logType,
		}
	}
	return unodeLogs
}
