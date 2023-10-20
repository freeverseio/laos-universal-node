package erc721

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
)

type EthService struct {
	Ethcli blockchain.EthRPCClient
}

// ChainId returns the chain ID of the ethService as a hexutil.Big.
// nolint:revive // needs to be named ChainId for EVM compatibility.
func (b *EthService) ChainId() *hexutil.Big {
	var result string
	err := b.Ethcli.Call(&result, "eth_chainId")
	if err != nil {
		return nil
	}
	chainId, err := hexutil.DecodeBig(result)
	if err != nil {
		return (*hexutil.Big)(big.NewInt(int64(0)))
	}
	return (*hexutil.Big)(chainId)
}

// BlockNumber returns a hardcoded value of 0 as the block number.
func (b *EthService) BlockNumber(_ context.Context) (hexutil.Uint64, error) {
	var result string
	err := b.Ethcli.Call(&result, "eth_blockNumber")
	if err != nil {
		return hexutil.Uint64(0), err
	}
	blockNum, err := hexutil.DecodeUint64(result)
	if err != nil {
		return hexutil.Uint64(0), err
	}
	return hexutil.Uint64(blockNum), nil
}

// GetBlockByNumber returns the block information for the specified block number.
func (b *EthService) GetBlockByNumber(blockNumber string, includeTransactions bool) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := b.Ethcli.Call(&result, "eth_getBlockByNumber", blockNumber, includeTransactions)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Call processes an Ethereum transaction call by delegating to erc721.ProcessCall.
func (b *EthService) Call(t blockchain.Transaction, blockNumber string) (hexutil.Bytes, error) {
	return ProcessCall(t.Data, common.HexToAddress(t.To), b.Ethcli, blockNumber)
}

// GetBalance returns a hardcoded value of 0 as the balance for a given Ethereum address.
func (b *EthService) GetBalance(_ common.Address, blockNumber string) (hexutil.Uint64, error) {
	return 0, nil
}

// GetCode returns a hardcoded value of 0 as the code for a given Ethereum address.
func (b *EthService) GetCode(_ common.Address, blockNumber string) (hexutil.Bytes, error) {
	return nil, nil
}
