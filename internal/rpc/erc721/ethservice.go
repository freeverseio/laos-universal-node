package erc721

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
)

type EthService struct {
	ethcli blockchain.EthRPCClient
}

func NewEthService(ethcli blockchain.EthRPCClient) *EthService {
	return &EthService{
		ethcli: ethcli,
	}
}

// ChainId returns the chain ID of the ethService as a hexutil.Big.
// nolint:revive // needs to be named ChainId for EVM compatibility.
func (b *EthService) ChainId() *hexutil.Big {
	var result string
	err := b.ethcli.Call(&result, "eth_chainId")
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
	err := b.ethcli.Call(&result, "eth_blockNumber")
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
	err := b.ethcli.Call(&result, "eth_getBlockByNumber", blockNumber, includeTransactions)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Call processes an Ethereum transaction call by delegating to erc721.ProcessCall.
func (b *EthService) Call(t blockchain.Transaction, blockNumber string) (hexutil.Bytes, error) {
	return ProcessCall(t.Data, common.HexToAddress(t.To), b.ethcli, blockNumber)
}

// GetBalance returns the balance of the specified address.
func (b *EthService) GetBalance(addr common.Address, blockNumber string) (hexutil.Uint64, error) {
	var result string
	err := b.ethcli.Call(&result, "eth_getBalance", addr, blockNumber)
	if err != nil {
		return hexutil.Uint64(0), err
	}
	balance, err := hexutil.DecodeUint64(result)
	if err != nil {
		return hexutil.Uint64(0), err
	}
	return hexutil.Uint64(balance), nil
}

func (b *EthService) GetCode(addr common.Address, blockNumber string) (hexutil.Bytes, error) {
	var result string
	err := b.ethcli.Call(&result, "eth_getCode", addr, blockNumber)
	if err != nil {
		return nil, err
	}
	return hexutil.Decode(result)
}

func (b *EthService) EstimateGas(t blockchain.Transaction) (hexutil.Bytes, error) {
	var result string
	err := b.ethcli.Call(&result, "eth_estimateGas", t)
	if err != nil {
		return nil, err
	}
	return hexutil.Decode(result)
}

func (b *EthService) GetTransactionCount(addr common.Address, blockNumber string) (hexutil.Uint64, error) {
	var result string
	err := b.ethcli.Call(&result, "eth_getTransactionCount", addr, blockNumber)
	if err != nil {
		return hexutil.Uint64(0), err
	}
	count, err := hexutil.DecodeUint64(result)
	if err != nil {
		return hexutil.Uint64(0), err
	}
	return hexutil.Uint64(count), nil
}

func (b *EthService) SendRawTransaction(data string) (hexutil.Bytes, error) {
	var result string
	err := b.ethcli.Call(&result, "eth_sendRawTransaction", data)
	if err != nil {
		return nil, err
	}
	return hexutil.Decode(result)
}
