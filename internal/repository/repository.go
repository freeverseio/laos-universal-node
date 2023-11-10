package repository

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
	"github.com/freeverseio/laos-universal-node/internal/scan"
)

const (
	contractPrefix = "contract_"
	chainID        = "chain_id"
)

type Service struct {
	storageService storage.Storage
}

func New(s storage.Storage) Service {
	return Service{
		s,
	}
}

// TODO test me
func (s *Service) StoreERC721UniversalContracts(universalContracts []scan.ERC721UniversalContract) error {
	tx := s.storageService.NewTransaction()
	defer tx.Discard()
	for i := 0; i < len(universalContracts); i++ {
		err := tx.Set([]byte(contractPrefix+universalContracts[i].Address.String()), []byte(universalContracts[i].BaseURI))
		if err != nil {
			return err
		}
	}
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetAllERC721UniversalContracts() ([]string, error) {
	var contracts []string
	keys, err := s.storageService.GetKeysWithPrefix([]byte(contractPrefix))
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		contracts = append(contracts, string(k))
	}
	return contracts, nil
}

func (s *Service) get(key string) ([]byte, error) {
	value, err := s.storageService.Get([]byte(key))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	return value, nil
}

func (s *Service) GetChainID() (string, error) {
	value, err := s.get(chainID)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (s *Service) SetChainID(chainIDValue string) error {
	return s.storageService.Set([]byte(chainID), []byte(chainIDValue))
}

// TODO decide if name should change to GetERC721UniversalContractBaseURI
func (s *Service) GetERC721UniversalContract(key string) (string, error) {
	value, err := s.storageService.Get([]byte(contractPrefix + key))
	if err != nil {
		return "", err
	}
	return string(value), nil
}
