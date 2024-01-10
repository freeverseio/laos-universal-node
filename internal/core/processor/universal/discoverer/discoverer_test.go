package discoverer_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	cDiscoverer "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"

	mockValidator "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer/validator/mock"
	mockScan "github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
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

		err := d.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
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

		err := d.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
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

		err := d.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
		assertError(t, expectedError, err)
	})
}

func TestDiscoverContractsErrMintEvents(t *testing.T) {
	t.Parallel()

	t.Run("error on get minted events for universal contract", func(t *testing.T) {
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
		errorOnGetMintedEvents := fmt.Errorf("error on getting minted events")
		expectedError := fmt.Errorf("error on getting minted events")

		startingBlock := uint64(100)
		lastBlock := uint64(200)

		ctx := context.TODO()
		tx, client, scanner, validator := createMocks(t)

		d := cDiscoverer.New(client, []string{}, scanner, validator)

		// Mock the scanner's ScanNewUniversalEvents method
		scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock))).
			Return([]scan.EventNewERC721Universal{event}, nil)

		validator.EXPECT().Validate(event).Return(expectedContract, nil)

		tx.EXPECT().StoreERC721UniversalContracts([]model.ERC721UniversalContract{expectedContract}).
			Return(nil)

		tx.EXPECT().LoadMerkleTrees(expectedContract.Address).Return(nil).Times(1)
		mintedEvents := getMockMintedEvents(120, 5)
		tx.EXPECT().GetMintedWithExternalURIEvents(expectedContract.CollectionAddress.String()).
			Return(mintedEvents, errorOnGetMintedEvents)

		err := d.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
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

		tx.EXPECT().LoadMerkleTrees(expectedContract.Address).Return(nil).Times(1)

		mintedEvents := getMockMintedEvents(120, 5)
		tx.EXPECT().GetMintedWithExternalURIEvents(expectedContract.CollectionAddress.String()).
			Return(mintedEvents, nil)

		client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(event.BlockNumber))).Return(&types.Header{Time: 130}, nil) // contract discovered after minted event
		tx.EXPECT().Mint(event.NewContractAddress, &mintedEvents[0]).Return(nil)
		tx.EXPECT().SetCurrentEvoEventsIndexForOwnershipContract(event.NewContractAddress.String(), uint64(1)).Return(nil)
		tx.EXPECT().TagRoot(event.NewContractAddress, int64(event.BlockNumber)).Return(nil)

		err := d.DiscoverContracts(ctx, tx, startingBlock, lastBlock)
		assertError(t, nil, err)
	})
}

func createMocks(t *testing.T) (*mockTx.MockTx, *mockScan.MockEthClient, *mockScan.MockScanner, *mockValidator.MockValidator) {
	ctrl := gomock.NewController(t)
	return mockTx.NewMockTx(ctrl), mockScan.NewMockEthClient(ctrl), mockScan.NewMockScanner(ctrl), mockValidator.NewMockValidator(ctrl)
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

func getMockMintedEvents(blockNumber, timestamp uint64) []model.MintedWithExternalURI {
	return []model.MintedWithExternalURI{
		{
			Slot:        big.NewInt(1),
			To:          common.HexToAddress("0x0"),
			TokenURI:    "",
			TokenId:     big.NewInt(1),
			BlockNumber: blockNumber,
			Timestamp:   timestamp,
		},
	}
}
