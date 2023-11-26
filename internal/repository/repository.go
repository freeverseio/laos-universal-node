package repository

import (
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage"
)

const (
	contractPrefix  = "contract_"
	chainID         = "chain_id"
	currentBlock    = "ownership_current_block"
	evoCurrentBlock = "evo_current_block"
)

type Service struct {
	storageService storage.Service
}

func New(s storage.Service) Service {
	return Service{
		s,
	}
}

// TODO I'm unused, remove me and move my tests to state
func (s *Service) StoreERC721UniversalContracts(universalContracts []model.ERC721UniversalContract) error {
	tx := s.storageService.NewTransaction()
	defer tx.Discard()
	for i := 0; i < len(universalContracts); i++ {
		addressLowerCase := strings.ToLower(universalContracts[i].Address.String())
		err := tx.Set([]byte(contractPrefix+addressLowerCase), universalContracts[i].CollectionAddress.Bytes())
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
		contract := strings.TrimPrefix(string(k), contractPrefix)
		contracts = append(contracts, contract)
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

// TODO move chain-related methods to state
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

func (s *Service) HasERC721UniversalContract(contract string) (bool, error) {
	value, err := s.get(contractPrefix + contract)
	if err != nil {
		return false, err
	}
	return value != nil, nil
}
