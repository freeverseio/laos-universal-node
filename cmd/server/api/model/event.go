package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type EventLog struct {
	Address          common.Address `json:"address,omitempty"`
	Topics           []common.Hash  `json:"topics,omitempty"`
	Data             []byte         `json:"data,omitempty"`
	BlockNumber      *big.Int       `json:"blockNumber,omitempty"`
	TransactionHash  common.Hash    `json:"transactionHash,omitempty"`
	TransactionIndex *big.Int       `json:"transactionIndex,omitempty"`
	BlockHash        common.Hash    `json:"blockHash,omitempty"`
	LogIndex         *big.Int       `json:"logIndex,omitempty"`
	Removed          bool           `json:"removed,omitempty"`
}
