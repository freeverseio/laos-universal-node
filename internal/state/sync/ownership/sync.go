package ownership

import (
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	contractEvoCurrentIndexPrefix = "ownership_contract_evo_current_index_"
	currentBlock                  = "ownership_current_block"
	endRangeBlockHash             = "ownership_end_range_block_hash"
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (s *service) SetCurrentEvoEventsIndexForOwnershipContract(contract string, number uint64) error {
	return s.tx.Set([]byte(contractEvoCurrentIndexPrefix+strings.ToLower(contract)), []byte(strconv.FormatUint(number, 10)))
}

func (s *service) GetCurrentEvoEventsIndexForOwnershipContract(contract string) (uint64, error) {
	value, err := s.tx.Get([]byte(contractEvoCurrentIndexPrefix + strings.ToLower(contract)))
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

func (s *service) SetOwnershipEndRangeBlockHash(blockHash common.Hash) error {
	return s.tx.Set([]byte(endRangeBlockHash), blockHash.Bytes())
}

func (s *service) GetOwnershipEndRangeBlockHash() (common.Hash, error) {
	value, err := s.tx.Get([]byte(endRangeBlockHash))
	if err != nil {
		return common.Hash{}, err
	}
	if value == nil {
		return common.Hash{}, nil
	}

	return common.BytesToHash(value), nil
}
