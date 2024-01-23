package evolution

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	eventsPrefix = "evo_events_"
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (s *service) StoreMintedWithExternalURIEvents(contract string, events []model.MintedWithExternalURI) error {
	for _, event := range events {
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		if err := encoder.Encode(event); err != nil {
			return err
		}
		if errSet := s.tx.Set([]byte(eventsPrefix+strings.ToLower(contract)+"_"+fmt.Sprint(event.BlockNumber)+"_"+fmt.Sprint(event.TxIndex)), buf.Bytes()); errSet != nil {
			return errSet
		}
	}
	return nil
}

func (s *service) GetMintedWithExternalURIEvents(contract string) ([]model.MintedWithExternalURI, error) {
	events := s.tx.GetValuesWithPrefix([]byte(eventsPrefix + strings.ToLower(contract) + "_"))
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
