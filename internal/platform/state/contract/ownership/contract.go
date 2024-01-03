package ownership

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
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
		err := s.tx.Set([]byte(contractPrefix+addressLowerCase), universalContracts[i].CollectionAddress.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetCollectionAddress(contract string) (common.Address, error) {
	contractLowerCase := strings.ToLower(contract)
	value, err := s.tx.Get([]byte(contractPrefix + contractLowerCase))
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(value), nil
}

func (s *service) GetExistingERC721UniversalContracts(contracts []string) ([]string, error) {
	var existingContracts []string
	for _, k := range contracts {
		hasContract, err := s.HasERC721UniversalContract(k)
		if err != nil {
			return nil, err
		}
		if hasContract {
			existingContracts = append(existingContracts, k)
		}
	}
	return existingContracts, nil
}

func (s *service) GetAllERC721UniversalContracts() []string {
	var contracts []string
	keys := s.tx.GetKeysWithPrefix([]byte(contractPrefix))
	for i := range keys {
		contract := strings.TrimPrefix(string(keys[i]), contractPrefix)
		contracts = append(contracts, contract)
	}
	return contracts
}

func (s *service) HasERC721UniversalContract(contract string) (bool, error) {
	lowerCaseContractAddress := strings.ToLower(contract)
	value, err := s.tx.Get([]byte(contractPrefix + lowerCaseContractAddress))
	if err != nil {
		return false, err
	}
	return value != nil, nil
}
