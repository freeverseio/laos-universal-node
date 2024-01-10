package ownership

import (
	"strconv"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/sync"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	contractEvoCurrentIndexPrefix = "ownership_contract_evo_current_index_"
	lastBlock                     = "ownership_last_block"
	ownershipBlockTag             = "ownership_block_"
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (s *service) SetLastOwnershipBlock(block model.Block) error {
	// Saving the block with blocknumber as key
	err := sync.SetBlock(s.tx, ownershipBlockTag+strconv.FormatUint(block.Number, 10), block)
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
	return sync.GetBlock(s.tx, ownershipBlockTag+strconv.FormatUint(blockNumber, 10))
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

func (s *service) GetAllStoredBlockNumbers() ([]uint64, error) {
	var blockNumbers []uint64
	keys := s.tx.GetKeysWithPrefix([]byte(ownershipBlockTag))
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
