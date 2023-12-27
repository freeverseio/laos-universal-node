package evolution

import (
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/sync"
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
	return sync.SetLastBlock(s.tx, lastBlock, block)
}

func (s *service) GetLastEvoBlock() (model.Block, error) {
	return sync.GetLastBlock(s.tx, lastBlock)
}
