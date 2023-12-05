package model

import (
	"github.com/ethereum/go-ethereum/common"
)

// TODO add timestamp
type ERC721UniversalContract struct {
	Address           common.Address
	CollectionAddress common.Address
	BlockNumber       uint64
}
