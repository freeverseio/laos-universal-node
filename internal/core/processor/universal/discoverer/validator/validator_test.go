package validator_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer/validator"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
)

func TestValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		event            scan.EventNewERC721Universal
		expectedContract model.ERC721UniversalContract
		expectedErr      error
	}{
		{
			name: "valid event event received",
			event: scan.EventNewERC721Universal{
				BaseURI:            "https://uloc.io/GlobalConsensus(3)/Parachain(9999)/AccountKey20(0x0000000000000000000000000000000000000000)/",
				BlockNumber:        12345,
				NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
			},
			expectedContract: model.ERC721UniversalContract{
				Address:           common.HexToAddress("0xc3dd09d5387fa0ab798e0adc152d15b8d1a299df"),
				CollectionAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
				BlockNumber:       12345,
			},
			expectedErr: nil,
		},
		{
			name: "invalid parachain",
			event: scan.EventNewERC721Universal{
				BaseURI:            "https://uloc.io/GlobalConsensus(3)/Parachain(8999)/AccountKey20(0x0000000000000000000000000000000000000000)/",
				BlockNumber:        100,
				NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
			},
			expectedErr: fmt.Errorf("universal contract's base URI points to a collection in a different evochain, contract discarded"),
		},
		{
			name: "invalid consensus",
			event: scan.EventNewERC721Universal{
				BaseURI:            "https://uloc.io/GlobalConsensus(2)/Parachain(9999)/AccountKey20(0x0000000000000000000000000000000000000000)/",
				BlockNumber:        100,
				NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
			},
			expectedErr: fmt.Errorf("universal contract's base URI points to a collection in a different evochain, contract discarded"),
		},
		{
			name: "no global consensus in BaseURI",
			event: scan.EventNewERC721Universal{
				BaseURI:            "https://uloc.io/Parachain(9999)/AccountKey20(0x0000000000000000000000000000000000000000)",
				BlockNumber:        100,
				NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
			},
			expectedErr: fmt.Errorf("no global consensus ID found in base URI: https://uloc.io/Parachain(9999)/AccountKey20(0x0000000000000000000000000000000000000000)"),
		},
		{
			name: "no parachain in BaseURI",
			event: scan.EventNewERC721Universal{
				BaseURI:            "https://uloc.io/GlobalConsensus(3)/AccountKey20(0x0000000000000000000000000000000000000000)",
				BlockNumber:        100,
				NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
			},
			expectedErr: fmt.Errorf("no parachain ID found in base URI: https://uloc.io/GlobalConsensus(3)/AccountKey20(0x0000000000000000000000000000000000000000)"),
		},
		{
			name: "no collection address in BaseURI",
			event: scan.EventNewERC721Universal{
				BaseURI:            "https://uloc.io/GlobalConsensus(3)/Parachain(9999)/",
				BlockNumber:        100,
				NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
			},
			expectedErr: fmt.Errorf("no collection address found in base URI: https://uloc.io/GlobalConsensus(3)/Parachain(9999)/"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := validator.New("3", 9999)

			contract, err := v.Validate(tt.event)

			assertError(t, tt.expectedErr, err)
			if err != nil {
				if contract.Address != tt.expectedContract.Address {
					t.Fatalf("expected address to be %v, got %v", tt.expectedContract.Address, contract.Address)
				}
				if contract.CollectionAddress != tt.expectedContract.CollectionAddress {
					t.Fatalf("expected collection address to be %v, got %v", tt.expectedContract.CollectionAddress, contract.CollectionAddress)
				}
				if contract.BlockNumber != tt.expectedContract.BlockNumber {
					t.Fatalf("expected block number to be %v, got %v", tt.expectedContract.BlockNumber, contract.BlockNumber)
				}
			}
		})
	}
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
