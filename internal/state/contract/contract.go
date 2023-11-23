package contract

import (
	"bytes"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	contractPrefix = "contract_"
	evoEvents      = "evo_events_"
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

func (c *service) StoreEvoChainMintEvents(contract common.Address, events []model.EventMintedWithExternalURI) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(events); err != nil {
		return err
	}

	return c.tx.Set([]byte(evoEvents+contract.String()), buf.Bytes())
}

func (c *service) EvoChainMintEvents(contract common.Address) ([]model.EventMintedWithExternalURI, error) {
	value, err := c.tx.Get([]byte(evoEvents + contract.String()))
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(value)
	decoder := gob.NewDecoder(buffer)

	var events []model.EventMintedWithExternalURI
	if err := decoder.Decode(&events); err != nil {
		return nil, err
	}

	return events, nil
}
