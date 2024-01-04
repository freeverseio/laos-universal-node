package evolution

import (
	"bytes"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	lastBlock = "evo_last_block"
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (s *service) SetLastEvoBlock(block model.Block) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(block); err != nil {
		return err
	}

	return s.tx.Set([]byte(lastBlock), buf.Bytes())
}

func (s *service) GetLastEvoBlock() (model.Block, error) {
	defaultBlock := model.Block{
		Number:    0,
		Timestamp: 0,
		Hash:      common.Hash{},
	}

	value, err := s.tx.Get([]byte(lastBlock))
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
