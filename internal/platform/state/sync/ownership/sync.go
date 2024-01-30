package ownership

import (
	"strconv"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/sync"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	lastBlock            = "ownership_last_block"
	ownershipBlockTag    = "ownership_block_"
	blockNumberDigits    = 18
	numberOfBlocksToKeep = 250
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
