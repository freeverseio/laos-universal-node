package contract

import (
	"strings"

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

func (s *service) StoreERC721UniversalContracts(universalContracts []model.ERC721UniversalContract) error {
	for i := 0; i < len(universalContracts); i++ {
		addressLowerCase := strings.ToLower(universalContracts[i].Address.String())
		err := s.tx.Set([]byte(contractPrefix+addressLowerCase), []byte(universalContracts[i].BaseURI))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetExistingERC721UniversalContracts(contracts []string) ([]string, error) {
	defer s.tx.Discard()
	var existingContracts []string
	for _, k := range contracts {
		hasContract, err := s.hasERC721UniversalContract(k)
		if err != nil {
			return nil, err
		}
		if hasContract {
			existingContracts = append(existingContracts, k)
		}
	}
	return existingContracts, nil
}

func (s *service) hasERC721UniversalContract(contract string) (bool, error) {
	defer s.tx.Discard()
	lowerCaseContractAddress := strings.ToLower(contract)
	value, err := s.tx.Get([]byte(contractPrefix + lowerCaseContractAddress))
	if err != nil {
		return false, err
	}
	return value != nil, nil
}
