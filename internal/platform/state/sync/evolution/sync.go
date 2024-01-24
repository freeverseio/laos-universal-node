package evolution

import (
	"strconv"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/sync"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	lastBlock         = "evo_last_block"
	ownershipBlockTag = "evo_block_"
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
	return sync.SetBlock(s.tx, lastBlock, block)
}

func (s *service) GetLastEvoBlock() (model.Block, error) {
	return sync.GetBlock(s.tx, lastBlock)
}

func (s *service) GetCorrespondingEvoBlockNumber(ownershipBlockNumber uint64) (uint64, error) {
	val, err := s.tx.Get([]byte(ownershipBlockTag + strconv.Itoa(int(ownershipBlockNumber))))
	if err != nil {
		return 0, err
	}
	if val == nil {
		return 0, nil
	}
	return strconv.ParseUint(string(val), 10, 64)
}

func (s *service) SetCorrespondingEvoBlockNumber(ownershipBlockNumber, evoBlockNumber uint64) error {
	return s.tx.Set([]byte(ownershipBlockTag+strconv.Itoa(int(ownershipBlockNumber))), []byte(strconv.Itoa(int(evoBlockNumber))))
}
