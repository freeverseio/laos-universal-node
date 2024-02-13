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
	lastBlock                = "ownership_last_block"
	ownershipBlockTag        = "ownership_block_"
	lastMappedOwnershipBlock = "mapped_ownership_last_block"
	mappedOwnershipBlock     = "mapped_ownership_block_"
	blockNumberDigits        = 18
	numberOfBlocksToKeep     = 250
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

// TODO move the mapping methods to another package?
func (s *service) SetLastMappedOwnershipBlockNumber(blockNumber uint64) error {
	return s.tx.Set([]byte(lastMappedOwnershipBlock), []byte(strconv.FormatUint(blockNumber, 10)))
}

func (s *service) SetOwnershipEvoBlockMapping(ownershipBlock, evoBlock uint64) error {
	return s.tx.Set([]byte(mappedOwnershipBlock+fmt.Sprint(ownershipBlock)), []byte(strconv.FormatUint(evoBlock, 10)))
}

func (s *service) GetLastMappedOwnershipBlockNumber() (uint64, error) {
	return s.getBlockNumber(lastMappedOwnershipBlock)
}

func (s *service) GetMappedEvoBlockNumber(ownershipBlock uint64) (uint64, error) {
	return s.getBlockNumber(mappedOwnershipBlock + fmt.Sprint(ownershipBlock))
}

func (s *service) getBlockNumber(key string) (uint64, error) {
	value, err := s.tx.Get([]byte(key))
	if err != nil {
		return 0, err
	}
	if value == nil {
		return 0, nil
	}
	return strconv.ParseUint(string(value), 10, 64)
}

func (s *service) SetOwnershipBlock(blockNumber uint64, block model.Block) error {
	formattedOwnershipBlockNumber := formatBlockNumber(blockNumber, blockNumberDigits)
	// Saving the block with blocknumber as key
	return sync.SetBlock(s.tx, ownershipBlockTag+formattedOwnershipBlockNumber, block)
}

func (s *service) SetLastOwnershipBlock(block model.Block) error {
	// Saving the block with blocknumber as key
	if err := s.SetOwnershipBlock(block.Number, block); err != nil {
		return err
	}
	// Saving the block with lastBlock as key
	return sync.SetBlock(s.tx, lastBlock, block)
}

func (s *service) GetLastOwnershipBlock() (model.Block, error) {
	return sync.GetBlock(s.tx, lastBlock)
}

func (s *service) GetOwnershipBlock(blockNumber uint64) (model.Block, error) {
	formattedOwnershipBlockNumber := formatBlockNumber(blockNumber, blockNumberDigits)
	return sync.GetBlock(s.tx, ownershipBlockTag+formattedOwnershipBlockNumber)
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
		if err := s.tx.Delete(key); err != nil {
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
			if err := s.tx.Delete(key); err != nil {
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
