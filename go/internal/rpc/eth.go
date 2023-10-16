package rpc

import (
	"context"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/rpc/erc721"
)

// Block represents an Ethereum block.
type Block struct {
	Number           hexutil.Uint64 `json:"number"`
	Hash             ecommon.Hash   `json:"hash"`
	ParentHash       ecommon.Hash   `json:"parentHash"`
	Timestamp        hexutil.Uint64 `json:"timestamp"`
	Transactions     []Transaction  `json:"transactions"`
	TransactionsRoot ecommon.Hash   `json:"transactionsRoot"`
}

// Eth is an interface that defines Ethereum RPC operations.
type Eth interface {
	ChainId() *hexutil.Big
	BlockNumber(ctx context.Context) (hexutil.Uint64, error)
	Call(t Transaction, blockNumber string) (hexutil.Bytes, error)
	GetBalance(address ecommon.Address, blockNumber string) (hexutil.Uint64, error)
	GetCode(address ecommon.Address, blockNumber string) (hexutil.Bytes, error)
	GetBlockByNumber(blockNumber string, includeTransactions bool) (*Block, error)
}

type EthService struct {
	ethcli            blockchain.EthClient
	contractAddr ecommon.Address
	chainID           uint64
}

// NewEthService creates a new instance of ethService with the given parameters,
// and returns it as an Eth interface.
func NewEthService(
	ethcli blockchain.EthClient,
	contractAddr ecommon.Address,
	chainID uint64,
) Eth {
	return &EthService{
		ethcli:            ethcli,
		contractAddr: contractAddr,
		chainID:           chainID,
	}
}

// ChainId returns the chain ID of the ethService as a hexutil.Big.
// nolint:revive // needs to be named ChainId for EVM compatibility.
func (b *EthService) ChainId() *hexutil.Big {
	return (*hexutil.Big)(big.NewInt(int64(b.chainID)))
}

// BlockNumber returns a hardcoded value of 0 as the block number.
func (b *EthService) BlockNumber(_ context.Context) (hexutil.Uint64, error) {
	return hexutil.Uint64(0), nil
}

// GetBlockByNumber returns the block information for the specified block number.
// We return an empty object(this is needed for Metamask integration)
func (b *EthService) GetBlockByNumber(blockNumber string, _ bool) (*Block, error) {
	return &Block{
		// adding empty transactions otherwise it will be nil
		Transactions: []Transaction{},
	}, nil
}

// Transaction represents an Ethereum transaction.
type Transaction struct {
	Data string
	To   string
}

// Call processes an Ethereum transaction call by delegating to erc721.ProcessCall.
func (b *EthService) Call(t Transaction, blockNumber string) (hexutil.Bytes, error) {
	log.Println("Call")
	to := common.HexToAddress(t.To)
	if to != b.contractAddr {
		return nil, errors.New("to != b.contractAddr")
	}
	return erc721.ProcessCall(t.Data, common.HexToAddress(t.To),  b.ethcli, b.contractAddr, b.chainID)
}

// GetBalance returns a hardcoded value of 0 as the balance for a given Ethereum address.
func (b *EthService) GetBalance(_ ecommon.Address, blockNumber string) (hexutil.Uint64, error) {
	return 0, nil
}

// GetCode returns a hardcoded value of 0 as the code for a given Ethereum address.
func (b *EthService) GetCode(_ ecommon.Address, blockNumber string) (hexutil.Bytes, error) {
	return nil, nil
}

