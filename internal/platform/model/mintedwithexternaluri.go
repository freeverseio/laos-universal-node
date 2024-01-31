package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type MintedWithExternalURI struct {
	Slot                 *big.Int
	To                   common.Address
	TokenURI             string
	TokenId              *big.Int
	BlockNumber          uint64
	OwnershipBlockNumber uint64
	Timestamp            uint64
	TxIndex              uint64
}
