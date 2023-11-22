package contract

import (
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	contractPrefix = "contract_"
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
