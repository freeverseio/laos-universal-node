package ownership

import (
	"strconv"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	contractEvoCurrentBlockPrefix = "ownership_contract_evo_current_block_"
	currentBlock                  = "ownership_current_block"
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (s *service) SetCurrentEvoBlockForOwnershipContract(contract string, number uint64) error {
	return s.tx.Set([]byte(contractEvoCurrentBlockPrefix+strings.ToLower(contract)), []byte(strconv.FormatUint(number, 10)))
}

func (s *service) GetCurrentEvoBlockForOwnershipContract(contract string) (uint64, error) {
	value, err := s.tx.Get([]byte(contractEvoCurrentBlockPrefix + strings.ToLower(contract)))
	if err != nil {
		return 0, err
	}
	if value == nil {
		value = []byte("0")
	}
	return strconv.ParseUint(string(value), 10, 64)
}

func (s *service) GetCurrentOwnershipBlock() (uint64, error) {
	value, err := s.tx.Get([]byte(currentBlock))
	if err != nil {
		return 0, err
	}
	if value == nil {
		value = []byte("0")
	}
	return strconv.ParseUint(string(value), 10, 64)
}

func (s *service) SetCurrentOwnershipBlock(number uint64) error {
	return s.tx.Set([]byte(currentBlock), []byte(strconv.FormatUint(number, 10)))
}
