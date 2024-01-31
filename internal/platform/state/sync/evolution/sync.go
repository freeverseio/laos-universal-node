package evolution

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/sync"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	lastBlock               = "evo_last_block"
	nextEvoEventBlockPrefix = "next_evo_event_block"
	lastEvoEventBlockPrefix = "last_evo_event_block"
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

// SetNextEvoEventBlockForOwnershipContract is used by evo processor for storing the next block that has events
func (s *service) SetNextEvoEventBlock(contract string, blockNumber uint64) error {
	value, err := s.tx.Get([]byte(lastEvoEventBlockPrefix + strings.ToLower(contract)))
	if err != nil {
		return err
	}

	uintValue := uint64(0)
	if value == nil && len(value) > 0 {
		uintValue, err = strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return err
		}
	}

	nextKey := fmt.Sprintf("%s_%s_%s",
		nextEvoEventBlockPrefix,
		strings.ToLower(contract),
		strconv.FormatUint(uintValue, 10))

	err = s.tx.Set([]byte(nextKey), []byte(strconv.FormatUint(blockNumber, 10)))
	if err != nil {
		return err
	}

	return s.tx.Set([]byte(lastEvoEventBlockPrefix+strings.ToLower(contract)), []byte(strconv.FormatUint(blockNumber, 10)))
}

// GetNextEvoEventBlockForOwnershipContract is used by universal processor for getting the next block that has events
func (s *service) GetNextEvoEventBlock(contract string, blockNumber uint64) (uint64, error) {
	key := fmt.Sprintf("%s_%s_%s",
		nextEvoEventBlockPrefix,
		strings.ToLower(contract),
		strconv.FormatUint(blockNumber, 10))

	value, err := s.tx.Get([]byte(key))
	if err != nil {
		return 0, err
	}

	if len(value) == 0 {
		return 0, nil
	}

	return strconv.ParseUint(string(value), 10, 64)
}
