package sync

import (
	"bytes"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

func SetBlock(tx storage.Tx, key string, block model.Block) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(block); err != nil {
		return err
	}

	return tx.Set([]byte(key), buf.Bytes())
}

func GetBlock(tx storage.Tx, key string) (model.Block, error) {
	defaultBlock := model.Block{
		Number:    0,
		Timestamp: 0,
		Hash:      common.Hash{},
	}

	value, err := tx.Get([]byte(key))
	if err != nil {
		return defaultBlock, err
	}
	if value == nil {
		return defaultBlock, nil
	}

	var block model.Block

	decoder := gob.NewDecoder(bytes.NewBuffer(value))
	if err := decoder.Decode(&block); err != nil {
		return defaultBlock, err
	}
	return block, nil
}
