package evolution

import (
	"bytes"
	"encoding/gob"
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
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(events); err != nil {
		return err
	}

	return s.tx.Set([]byte(eventsPrefix+strings.ToLower(contract)), buf.Bytes())
}

func (s *service) GetMintedWithExternalURIEvents(contract string) ([]model.MintedWithExternalURI, error) {
	value, err := s.tx.Get([]byte(eventsPrefix + strings.ToLower(contract)))
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	var mintedEvents []model.MintedWithExternalURI
	decoder := gob.NewDecoder(bytes.NewBuffer(value))
	if err := decoder.Decode(&mintedEvents); err != nil {
		return nil, err
	}
	return mintedEvents, nil
}
