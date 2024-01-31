package discoverer_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/mock/gomock"

	cDiscoverer "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer"
	mockValidator "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer/validator/mock"
	mockClient "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	mockScan "github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
)

func TestShouldDiscover(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		contracts        []string
		expectedDiscover bool
		expectedError    error
	}{
		{
			name:      "no contracts in the list",
			contracts: []string{},

			expectedDiscover: true,
			expectedError:    nil,
		},
		{
			name:             "contract should be discovered",
			contracts:        []string{"contract1"},
			expectedDiscover: true,
			expectedError:    nil,
		},
		{
			name:             "contract should not be discovered",
			contracts:        []string{"contract1"},
			expectedDiscover: false,
			expectedError:    nil,
		},
		{
			name:             "error happened on execution",
			contracts:        []string{"contract1"},
			expectedDiscover: false,
			expectedError:    fmt.Errorf("error happened on checking contracts"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tx, client, _, _ := createMocks(t)

			discoverer := cDiscoverer.New(client, tt.contracts, nil, nil)
			if len(tt.contracts) == 1 {
				tx.EXPECT().HasERC721UniversalContract(tt.contracts[0]).Return(!tt.expectedDiscover, tt.expectedError)
			}

			shouldDiscover, err := discoverer.ShouldDiscover(tx, 10, 20)

			assertError(t, tt.expectedError, err)
			if shouldDiscover != tt.expectedDiscover {
				t.Fatalf("expected shouldDiscover to be %v, got %v", tt.expectedDiscover, shouldDiscover)
			}
		})
	}
}

func TestGetContracts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		contracts         []string
		expectedContracts []string
	}{
		{
			name:              "GetAllERC721UniversalContracts",
			contracts:         []string{},
			expectedContracts: []string{"contract1", "contract2"},
		},
		{
			name:              "GetExistingERC721UniversalContracts",
			contracts:         []string{"contract1"},
			expectedContracts: []string{"contract1"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tx, client, _, _ := createMocks(t)

			discoverer := cDiscoverer.New(client, tt.contracts, nil, nil)
			if len(tt.contracts) > 0 {
				tx.EXPECT().GetExistingERC721UniversalContracts(tt.contracts).Return(tt.expectedContracts, nil)
			} else {
				tx.EXPECT().GetAllERC721UniversalContracts().Return(tt.expectedContracts)
			}

			contracts, err := discoverer.GetContracts(tx)

			assertError(t, nil, err)
			if len(tt.expectedContracts) != len(contracts) {
				t.Fatalf("mismatch in expected contacts")
			}
		})
	}
}

func TestDiscoverContractsErrorOnScanning(t *testing.T) {
	t.Parallel()

	t.Run("error when scanning universal contracts", func(t *testing.T) {
		t.Parallel()

		event := scan.EventNewERC721Universal{
			BaseURI:            "https://uloc.io/GlobalConsensus(3)/Parachain(9999)/AccountKey20(0x0000000000000000000000000000000000000000)/",
			BlockNumber:        123,
			NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
		}
		errorOnScanning := fmt.Errorf("error when scanning universal contracts")
		expectedError := fmt.Errorf("error when scanning universal contracts")

		ctx := context.TODO()
		tx, client, scanner, validator := createMocks(t)

		startingBlock := uint64(100)
		lastBlock := uint64(200)

		d := cDiscoverer.New(client, []string{}, scanner, validator)

		// Mock the scanner's ScanNewUniversalEvents method
		scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock))).
			Return([]scan.EventNewERC721Universal{event}, errorOnScanning)

		_, err := d.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
		assertError(t, expectedError, err)
	})
}

func TestDiscoverContractsErrorOnValidating(t *testing.T) {
	t.Parallel()
	t.Run("error when validating events, contract is ignored nil returned", func(t *testing.T) {
		t.Parallel()
		event := scan.EventNewERC721Universal{
			BaseURI:            "https://uloc.io/GlobalConsensus(3)/Parachain(9999)/AccountKey20(0x0000000000000000000000000000000000000000)/",
			BlockNumber:        123,
			NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
		}
		expectedContract := model.ERC721UniversalContract{
			Address:           common.HexToAddress("0xc3dd09d5387fa0ab798e0adc152d15b8d1a299df"),
			CollectionAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			BlockNumber:       123,
		}
		errorOnValidation := fmt.Errorf("error when validating events")

		ctx := context.TODO()
		tx, client, scanner, validator := createMocks(t)

		startingBlock := uint64(100)
		lastBlock := uint64(200)

		d := cDiscoverer.New(client, []string{}, scanner, validator)

		// Mock the scanner's ScanNewUniversalEvents method
		scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock))).
			Return([]scan.EventNewERC721Universal{event}, nil)
		validator.EXPECT().Validate(event).Return(expectedContract, errorOnValidation)

		_, err := d.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
		assertError(t, nil, err)
	})
}

func TestDiscoverContractsErrorOnStoring(t *testing.T) {
	t.Parallel()

	t.Run("error on storing universal contract", func(t *testing.T) {
		t.Parallel()

		event := scan.EventNewERC721Universal{
			BaseURI:            "https://uloc.io/GlobalConsensus(3)/Parachain(9999)/AccountKey20(0x0000000000000000000000000000000000000000)/",
			BlockNumber:        123,
			NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
		}
		expectedContract := model.ERC721UniversalContract{
			Address:           common.HexToAddress("0xc3dd09d5387fa0ab798e0adc152d15b8d1a299df"),
			CollectionAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			BlockNumber:       123,
		}
		errorOnStoringUniversalContract := fmt.Errorf("error on storing")
		expectedError := fmt.Errorf("error on storing")

		ctx := context.TODO()
		tx, client, scanner, validator := createMocks(t)

		startingBlock := uint64(100)
		lastBlock := uint64(200)

		d := cDiscoverer.New(client, []string{}, scanner, validator)

		// Mock the scanner's ScanNewUniversalEvents method
		scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock))).
			Return([]scan.EventNewERC721Universal{event}, nil)

		validator.EXPECT().Validate(event).Return(expectedContract, nil)

		tx.EXPECT().StoreERC721UniversalContracts([]model.ERC721UniversalContract{expectedContract}).
			Return(errorOnStoringUniversalContract)

		_, err := d.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
		assertError(t, expectedError, err)
	})
}

func TestDiscoverContractsSuccess(t *testing.T) {
	t.Parallel()

	t.Run("successfully discover contract", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		tx, client, scanner, validator := createMocks(t)

		event := scan.EventNewERC721Universal{
			BaseURI:            "https://uloc.io/GlobalConsensus(3)/Parachain(9999)/AccountKey20(0x0000000000000000000000000000000000000000)/",
			BlockNumber:        123,
			NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
		}
		expectedContract := model.ERC721UniversalContract{
			Address:           common.HexToAddress("0xc3dd09d5387fa0ab798e0adc152d15b8d1a299df"),
			CollectionAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			BlockNumber:       123,
		}

		startingBlock := uint64(100)
		lastBlock := uint64(200)

		d := cDiscoverer.New(client, []string{}, scanner, validator)

		// Mock the scanner's ScanNewUniversalEvents method
		scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock))).
			Return([]scan.EventNewERC721Universal{event}, nil)

		validator.EXPECT().Validate(event).Return(expectedContract, nil)

		tx.EXPECT().StoreERC721UniversalContracts([]model.ERC721UniversalContract{expectedContract}).
			Return(nil)

		contracts, err := d.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
		assertError(t, nil, err)
		if contracts[event.NewContractAddress] != 123 {
			t.Fatalf("expected new contract %v, got %v", event.NewContractAddress, contracts[event.NewContractAddress])
		}
	})
}

func createMocks(t *testing.T) (*mockTx.MockTx, *mockClient.MockEthClient, *mockScan.MockScanner, *mockValidator.MockValidator) {
	ctrl := gomock.NewController(t)
	return mockTx.NewMockTx(ctrl), mockClient.NewMockEthClient(ctrl), mockScan.NewMockScanner(ctrl), mockValidator.NewMockValidator(ctrl)
}

func assertError(t *testing.T, expectedError, err error) {
	t.Helper()
	if expectedError != nil {
		if err.Error() != expectedError.Error() {
			t.Fatalf(`got error "%v", expected error: "%v"`, err, expectedError)
		}
	} else {
		if err != expectedError {
			t.Fatalf(`got error "%v", expected error: "%v"`, err, expectedError)
		}
	}
}
