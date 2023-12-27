package model

import "github.com/ethereum/go-ethereum/common"

type Block struct {
	Number    uint64
	Timestamp uint64
	Hash      common.Hash
}
