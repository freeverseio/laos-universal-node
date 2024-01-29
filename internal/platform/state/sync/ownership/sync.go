package ownership

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/sync"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	contractCurrentEvoEventBlockPrefix = "ownership_contract_current_evo_event_block_"
	contractNextEvoEventBlockPrefix    = "ownership_contract_next_evo_event_block_"
	contractLastEvoEventBlockPrefix    = "ownership_contract_last_evo_event_block_"
	lastBlock                          = "ownership_last_block"
	ownershipBlockTag                  = "ownership_block_"
	blockNumberDigits                  = 18
	numberOfBlocksToKeep               = 250
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (s *service) SetOwnershipBlock(blockNumber uint64, block model.Block) error {
	formatedOwnershipBlockNumber := formatBlockNumber(blockNumber, blockNumberDigits)
	// Saving the block with blocknumber as key
	return sync.SetBlock(s.tx, ownershipBlockTag+formatedOwnershipBlockNumber, block)
}

func (s *service) SetLastOwnershipBlock(block model.Block) error {
	// Saving the block with blocknumber as key
	err := s.SetOwnershipBlock(block.Number, block)
	if err != nil {
		return err
	}
	// Saving the block with lastBlock as key
	return sync.SetBlock(s.tx, lastBlock, block)
}

func (s *service) GetLastOwnershipBlock() (model.Block, error) {
	return sync.GetBlock(s.tx, lastBlock)
}

func (s *service) GetOwnershipBlock(blockNumber uint64) (model.Block, error) {
	formatedOwnershipBlockNumber := formatBlockNumber(blockNumber, blockNumberDigits)
	return sync.GetBlock(s.tx, ownershipBlockTag+formatedOwnershipBlockNumber)
}

// SetCurrentEvoBlockForOwnershipContract is used by universal processor updater to store the last block number
func (s *service) SetCurrentEvoBlockForOwnershipContract(contract string, number uint64) error {
	return s.tx.Set([]byte(contractCurrentEvoEventBlockPrefix+strings.ToLower(contract)), []byte(strconv.FormatUint(number, 10)))
}

// GetCurrentEvoBlockForOwnershipContract is used by universal processor updater to get the last block number
func (s *service) GetCurrentEvoBlockForOwnershipContract(contract string) (uint64, error) {
	value, err := s.tx.Get([]byte(contractCurrentEvoEventBlockPrefix + strings.ToLower(contract)))
	if err != nil {
		return 0, err
	}
	if value == nil {
		value = []byte("0")
	}
	return strconv.ParseUint(string(value), 10, 64)
}

// SetNextEvoEventBlockForOwnershipContract is used by evo processor for storing the next block that has events
func (s *service) SetNextEvoEventBlockForOwnershipContract(contract string, blockNumber uint64) error {
	value, err := s.tx.Get([]byte(contractLastEvoEventBlockPrefix + strings.ToLower(contract)))
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
		contractNextEvoEventBlockPrefix,
		strings.ToLower(contract),
		strconv.FormatUint(uintValue, 10))

	err = s.tx.Set([]byte(nextKey), []byte(strconv.FormatUint(blockNumber, 10)))
	if err != nil {
		return err
	}

	return s.tx.Set([]byte(contractLastEvoEventBlockPrefix+strings.ToLower(contract)), []byte(strconv.FormatUint(blockNumber, 10)))
}

// GetNextEvoEventBlockForOwnershipContract is used by universal processor for getting the next block that has events
func (s *service) GetNextEvoEventBlockForOwnershipContract(contract string, blockNumber uint64) (uint64, error) {
	key := fmt.Sprintf("%s_%s_%s",
	contractNextEvoEventBlockPrefix,
	strings.ToLower(contract),
	strconv.FormatUint(blockNumber, 10))

	value, err := s.tx.Get([]byte(key))
	if err != nil {
		return 0, err
	}

	if value == nil && len(value) > 0 {
		return 0, nil
	}

	return strconv.ParseUint(string(value), 10, 64)
}

func (s *service) GetAllStoredBlockNumbers() ([]uint64, error) {
	var blockNumbers []uint64
	keys := s.tx.GetKeysWithPrefix([]byte(ownershipBlockTag), true)
	for i := range keys {
		blockNumberStr := strings.TrimPrefix(string(keys[i]), ownershipBlockTag)
		blockNumber, err := strconv.ParseInt(blockNumberStr, 10, 64)
		if err != nil {
			return nil, err
		}

		blockNumbers = append(blockNumbers, uint64(blockNumber))
	}
	return blockNumbers, nil
}

func (s *service) DeleteOldStoredBlockNumbers() error {
	keys := s.tx.GetKeysWithPrefix([]byte(ownershipBlockTag), true)

	// Skip the first 250 keys (newest entries)
	if len(keys) > numberOfBlocksToKeep {
		keys = keys[numberOfBlocksToKeep:]
	} else {
		return nil
	}

	// Delete all keys beyond the newest 250
	for _, key := range keys {
		err := s.tx.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) DeleteOrphanBlockData(blockNumberRef uint64) error {
	keys := s.tx.GetKeysWithPrefix([]byte(ownershipBlockTag), true)
	// Delete all keys
	for i, key := range keys {
		blockNumberStr := strings.TrimPrefix(string(keys[i]), ownershipBlockTag)
		blockNumber, err := strconv.ParseInt(blockNumberStr, 10, 64)
		if err != nil {
			return err
		}
		if blockNumber > int64(blockNumberRef) {
			err := s.tx.Delete(key)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func formatBlockNumber(blockNumber uint64, blockNumberDigits uint16) string {
	// Convert the block number to a string
	blockNumberString := strconv.FormatUint(blockNumber, 10)
	// Pad with leading zeros if shorter
	for len(blockNumberString) < int(blockNumberDigits) {
		blockNumberString = "0" + blockNumberString
	}
	return blockNumberString
}
