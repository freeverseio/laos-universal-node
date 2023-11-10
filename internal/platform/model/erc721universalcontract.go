package model

import "github.com/ethereum/go-ethereum/common"

type ERC721UniversalContract struct {
	Address common.Address `json:"address"`
	Block   uint64         `json:"block"`
	BaseURI string         `json:"base_uri"`
}
