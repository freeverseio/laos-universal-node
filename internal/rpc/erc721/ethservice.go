package erc721

import (
	"context"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/blockchain"
)

type EthService struct {
	Ethcli       blockchain.EthRPCClient
	ContractAddr common.Address
	ChainID      uint64
}

// ChainId returns the chain ID of the ethService as a hexutil.Big.
// nolint:revive // needs to be named ChainId for EVM compatibility.
func (b *EthService) ChainId() *hexutil.Big {
	return (*hexutil.Big)(big.NewInt(int64(b.ChainID)))
}

// BlockNumber returns a hardcoded value of 0 as the block number.
func (b *EthService) BlockNumber(_ context.Context) (hexutil.Uint64, error) {
	return hexutil.Uint64(0), nil
}

// GetBlockByNumber returns the block information for the specified block number.
// We return an empty object(this is needed for Metamask integration)
func (b *EthService) GetBlockByNumber(blockNumber string, _ bool) (*blockchain.Block, error) {
	return &blockchain.Block{
		// adding empty transactions otherwise it will be nil
		Transactions: []blockchain.Transaction{},
	}, nil
}

// Call processes an Ethereum transaction call by delegating to erc721.ProcessCall.
func (b *EthService) Call(t blockchain.Transaction, blockNumber string) (hexutil.Bytes, error) {
	log.Println("Call")
	to := common.HexToAddress(t.To)
	if to != b.ContractAddr {
		return nil, errors.New("to != b.contractAddr")
	}
	return ProcessCall(t.Data, common.HexToAddress(t.To), b.Ethcli, b.ContractAddr, b.ChainID)
}

// GetBalance returns a hardcoded value of 0 as the balance for a given Ethereum address.
func (b *EthService) GetBalance(_ common.Address, blockNumber string) (hexutil.Uint64, error) {
	return 0, nil
}

// GetCode returns a hardcoded value of 0 as the code for a given Ethereum address.
func (b *EthService) GetCode(_ common.Address, blockNumber string) (hexutil.Bytes, error) {
	return nil, nil
}
