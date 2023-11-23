package evolution

import (
	"encoding/json"
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

func (s *service) GetEvoChainEvents(contract string) ([]model.MintedWithExternalURI, error) {
	defer s.tx.Discard()
	value, err := s.tx.Get([]byte(eventsPrefix + strings.ToLower(contract)))
	if err != nil {
		return nil, err
	}
	var mintedEvents []model.MintedWithExternalURI
	err = json.Unmarshal(value, &mintedEvents) // TODO check what happens when value is nil
	if err != nil {
		return nil, err
	}
	return mintedEvents, nil
}
