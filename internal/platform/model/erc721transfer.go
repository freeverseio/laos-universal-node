package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ERC721Transfer struct {
	From        common.Address
	To          common.Address
	TokenId     *big.Int
	BlockNumber uint64
	Timestamp   uint64
}
