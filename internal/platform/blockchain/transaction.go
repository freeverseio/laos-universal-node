package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Transaction represents an Ethereum transaction.
type Transaction struct {
	Data string
	To   string
}

// Block represents an Ethereum block.
type Block struct {
	Number           hexutil.Uint64 `json:"number"`
	Hash             common.Hash    `json:"hash"`
	ParentHash       common.Hash    `json:"parentHash"`
	Timestamp        hexutil.Uint64 `json:"timestamp"`
	Transactions     []Transaction  `json:"transactions"`
	TransactionsRoot common.Hash    `json:"transactionsRoot"`
}
