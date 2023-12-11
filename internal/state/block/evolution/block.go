package evolution

import (
	"strconv"

	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	currentBlock          = "evo_current_block"
	currentBlockTimestamp = "evo_current_block_timestamp"
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (s *service) GetCurrentEvoBlock() (uint64, error) {
	value, err := s.tx.Get([]byte(currentBlock))
	if err != nil {
		return 0, err
	}
	if value == nil {
		value = []byte("0")
	}
	return strconv.ParseUint(string(value), 10, 64)
}

func (s *service) SetCurrentEvoBlock(number uint64) error {
	return s.tx.Set([]byte(currentBlock), []byte(strconv.FormatUint(number, 10)))
}

func (s *service) GetCurrentEvoBlockTimestamp() (uint64, error) {
	value, err := s.tx.Get([]byte(currentBlockTimestamp))
	if err != nil {
		return 0, err
	}
	if value == nil {
		value = []byte("0")
	}
	return strconv.ParseUint(string(value), 10, 64)
}

func (s *service) SetCurrentEvoBlockTimestamp(number uint64) error {
	return s.tx.Set([]byte(currentBlockTimestamp), []byte(strconv.FormatUint(number, 10)))
}
