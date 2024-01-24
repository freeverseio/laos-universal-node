package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ERC721Event struct {
	Address   common.Address
	Data      []byte
	Topics    []common.Hash
	BlockHash common.Hash
	TxHash    common.Hash
	TxIndex   uint
	Index     uint
	Removed   bool
}

type MintedWithExternalURI struct {
	Slot        *big.Int
	To          common.Address
	TokenURI    string
	TokenId     *big.Int
	BlockNumber uint64
	Timestamp   uint64
	TxIndex     uint64
}
