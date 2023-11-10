package model

import "github.com/ethereum/go-ethereum/common"

type ERC721UniversalContract struct {
	Address common.Address `json:"address"`
	BaseURI string         `json:"base_uri"`
}
