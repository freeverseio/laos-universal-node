package contract

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	contractPrefix = "contract_"
	eventsPrefix   = "evo_events_"
)

type service struct {
	tx storage.Tx
}

func NewService(tx storage.Tx) *service {
	return &service{
		tx: tx,
	}
}

func (c *service) StoreERC721UniversalContracts(universalContracts []model.ERC721UniversalContract) error {
	for i := 0; i < len(universalContracts); i++ {
		err := c.tx.Set([]byte(contractPrefix+universalContracts[i].Address.String()), []byte(universalContracts[i].BaseURI))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *service) StoreEvoChainMintEvents(contract common.Address, events []model.MintedWithExternalURI) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(events); err != nil {
		return err
	}

	return c.tx.Set([]byte(eventsPrefix+contract.String()), buf.Bytes())
}

func (s *service) GetEvoChainEvents(contract common.Address) ([]model.MintedWithExternalURI, error) {
	defer s.tx.Discard()
	value, err := s.tx.Get([]byte(eventsPrefix + strings.ToLower(contract.Hex())))
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
