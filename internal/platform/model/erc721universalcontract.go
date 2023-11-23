package model

import "github.com/ethereum/go-ethereum/common"

type ERC721UniversalContract struct {
	Address           common.Address `json:"address"`
	CollectionAddress common.Address `json:"collection_address"`
}
