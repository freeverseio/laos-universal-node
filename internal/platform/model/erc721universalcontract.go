package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ERC721UniversalContract struct {
	Address common.Address `json:"address"`
	BaseURI string         `json:"base_uri"`
}

type MintedWithExternalURI struct {
	Slot        *big.Int
	To          common.Address
	TokenURI    string
	TokenId     *big.Int
	BlockNumber uint64
	Timestamp   uint64
}
