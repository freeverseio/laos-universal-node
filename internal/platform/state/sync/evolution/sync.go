package evolution

import (
	"fmt"
	"strconv"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/sync"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	lastBlock         = "evo_last_block"
	ownershipBlockTag = "evo_block_"
	evoEvent          = "evo_event_"
	blockNumberDigits = 18
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

func (s *service) GetCorrespondingEvoBlockNumber(ownershipBlockNumber uint64) (uint64, error) {
	val, err := s.tx.Get([]byte(ownershipBlockTag + strconv.Itoa(int(ownershipBlockNumber))))
	if err != nil {
		return 0, err
	}
	if val == nil {
		return 0, nil
	}
	return strconv.ParseUint(string(val), 10, 64)
}

func (s *service) SetCorrespondingEvoBlockNumber(ownershipBlockNumber, evoBlockNumber uint64) error {
	return s.tx.Set([]byte(ownershipBlockTag+strconv.Itoa(int(ownershipBlockNumber))), []byte(strconv.Itoa(int(evoBlockNumber))))
}

func (s *service) SetEvoEventOwnershipBlockNumber(ownershipBlockNumber uint64, event *model.MintedWithExternalURI, contract string) error {
	key := fmt.Sprintf("%s%s_%s_%d", evoEvent, sync.FormatNumberForSorting(event.BlockNumber, blockNumberDigits), contract, event.TxIndex)
	value := strconv.FormatUint(ownershipBlockNumber, 10)
	return s.tx.Set([]byte(key), []byte(value))
}

func (s *service) GetEvoEventOwnershipBlockNumber(blockNumber uint64, txIndex uint64, contract string) (uint64, error) {
	key := fmt.Sprintf("%s%s_%s_%d", evoEvent, sync.FormatNumberForSorting(blockNumber, blockNumberDigits), contract, txIndex)
	value, err := s.tx.Get([]byte(key))
	if err != nil {
		return 0, err
	}
	if value == nil {
		return 0, nil
	}
	return strconv.ParseUint(string(value), 10, 64)
}
