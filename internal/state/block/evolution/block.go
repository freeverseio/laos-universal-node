package evolution

import (
	"strconv"

	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	currentBlock = "evo_current_block"
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
	defer s.tx.Discard()
	value, err := s.tx.Get([]byte(currentBlock))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(value), 10, 64)
}

func (s *service) SetCurrentEvoBlock(number uint64) error {
	return s.tx.Set([]byte(currentBlock), []byte(strconv.FormatUint(number, 10)))
}
