package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/rpc/erc721"
)

// Eth is an interface that defines Ethereum RPC operations.
type Eth interface {
	ChainId() *hexutil.Big
	BlockNumber(ctx context.Context) (hexutil.Uint64, error)
	Call(t blockchain.Transaction, blockNumber string) (hexutil.Bytes, error)
	GetBalance(address common.Address, blockNumber string) (hexutil.Uint64, error)
	GetCode(address common.Address, blockNumber string) (hexutil.Bytes, error)
	GetBlockByNumber(blockNumber string, includeTransactions bool) (map[string]interface{}, error)
}

// NewEthService creates a new instance of ethService with the given parameters,
// and returns it as an Eth interface.
func NewEthService(
	ethcli blockchain.EthRPCClient,
) Eth {
	return &erc721.EthService{
		Ethcli: ethcli,
	}
}
