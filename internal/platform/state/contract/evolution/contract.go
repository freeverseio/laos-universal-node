package evolution

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	eventsPrefix      = "evo_events_"
	blockNumberDigits = 18
	txIndexDigits     = 8
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (s *service) StoreMintedWithExternalURIEvent(contract string, event *model.MintedWithExternalURI) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(event); err != nil {
		return err
	}
	key := fmt.Sprintf("%s%s_%s_%s", eventsPrefix,
		strings.ToLower(contract),
		formatNumberForSorting(event.BlockNumber, blockNumberDigits),
		formatNumberForSorting(event.TxIndex, txIndexDigits))

	return s.tx.Set([]byte(key), buf.Bytes())
}

func (s *service) GetMintedWithExternalURIEvents(contract string, blockNumber uint64) ([]model.MintedWithExternalURI, error) {
	key := fmt.Sprintf("%s%s_%s", eventsPrefix,
		strings.ToLower(contract),
		formatNumberForSorting(blockNumber, blockNumberDigits))

	events := s.tx.GetValuesWithPrefix([]byte(key))
	var mintedEvents []model.MintedWithExternalURI
	if len(events) == 0 {
		return mintedEvents, nil
	}

	for _, event := range events {
		var mintedEvent model.MintedWithExternalURI
		decoder := gob.NewDecoder(bytes.NewBuffer(event))
		if err := decoder.Decode(&mintedEvent); err != nil {
			return nil, err
		}
		mintedEvents = append(mintedEvents, mintedEvent)
	}
	return mintedEvents, nil
}

// we add digits to the block number and tx index to make sure the keys are sorted correctly
// since badger sorts the keys lexicographically
func formatNumberForSorting(blockNumber uint64, blockNumberDigits uint16) string {
	// Convert the block number to a string
	blockNumberString := strconv.FormatUint(blockNumber, 10)
	// Pad with leading zeros if shorter
	for len(blockNumberString) < int(blockNumberDigits) {
		blockNumberString = "0" + blockNumberString
	}
	return blockNumberString
}
