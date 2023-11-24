package ownership

import (
	"strconv"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	contractEvoCurrentBlockPrefix = "ownership_contract_evo_current_block_"
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (s *service) SetCurrentEvoBlockForOwnershipContract(contract string, blockNumber uint64) error {
	return s.tx.Set([]byte(contractEvoCurrentBlockPrefix+strings.ToLower(contract)), []byte(strconv.FormatUint(blockNumber, 10)))
}

func (s *service) GetCurrentEvoBlockForOwnershipContract(contract string) (uint64, error) {
	defer s.tx.Discard()
	value, err := s.tx.Get([]byte(contractEvoCurrentBlockPrefix + strings.ToLower(contract)))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(value), 10, 64)
}
