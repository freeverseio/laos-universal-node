package scan_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
	"go.uber.org/mock/gomock"
)

const (
	transferEventHash               = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	approveEventHash                = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
	approveForAllEventHash          = "0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31"
	newERC721UniversalEventHash     = "0x74b81bc88402765a52dad72d3d893684f472a679558f3641500e0ee14924a10a"
	newCollectionEventHash          = "0x6eb24fd767a7bcfa417f3fe25a2cb245d2ae52293d3c4a8f8c6450a09795d289"
	mintedWithExternalURIEventHash  = "0x4b3b5da28a351f8bb73b960d7c80b2cef3e3570cb03448234dee173942c74786"
	evolvedWithExternalURIEventHash = "0x95c167d04a267f10e6b3f373c7a336dc65cf459caf048854dc32a2d37ab1607c"
)

func TestParseEvents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fromBlock *big.Int
		toBlock   *big.Int
		address   common.Address
		eventLogs []types.Log
	}{
		{
			name:      "it should parse Transfer events",
			fromBlock: big.NewInt(0),
			toBlock:   big.NewInt(100),
			address:   common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(transferEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x00000000000000000000000066666f58de1bcd762a5e5c5aff9cc3c906d66666"),
						common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
					},
				},
			},
		},
		{
			name:      "it should parse Approval events",
			fromBlock: big.NewInt(0),
			toBlock:   big.NewInt(100),
			address:   common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(approveEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
						common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
					},
					Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
		},
		{
			name:      "it should parse ApprovalForAll events",
			fromBlock: big.NewInt(0),
			toBlock:   big.NewInt(100),
			address:   common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(approveForAllEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x0000000000000000000000001e0049783f008a0085193e00003d00cd54003c71"),
					},
					Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
		},
		{
			name:      "it should parse NewCollection event",
			fromBlock: big.NewInt(0),
			toBlock:   big.NewInt(100),
			address:   common.HexToAddress("0x0000000000000000000000000000000000000403"),
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(newCollectionEventHash),
						common.HexToHash("0x0000000000000000000000009c231bfb2dbbd2c92a9f7c710d51e0796c52dce5"),
					},
					Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000094"),
				},
			},
		},
		{
			name:      "it should parse MintedWithExternalURI events",
			fromBlock: big.NewInt(0),
			toBlock:   big.NewInt(100),
			address:   common.HexToAddress("0x0000000000000000000000000000000000000403"),
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(mintedWithExternalURIEventHash),
						common.HexToHash("0x000000000000000000000000684bc8f81250ad3f7f930b27586b799a4dda957b"),
					},
					Data: common.Hex2Bytes("000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000080000000000000000000000002684bc8f81250ad3f7f930b27586b799a4dda957b000000000000000000000000000000000000000000000000000000000000000c6d795f66697273745f7572690000000000000000000000000000000000000000"),
				},
			},
		},
		{
			name:      "it should parse EvolvedWithExternalURIevents",
			fromBlock: big.NewInt(0),
			toBlock:   big.NewInt(100),
			address:   common.HexToAddress("0x0000000000000000000000000000000000000403"),
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(evolvedWithExternalURIEventHash),
						common.HexToHash("0x000000000000000000000002684bc8f81250ad3f7f930b27586b799a4dda957b"),
					},
					Data: common.Hex2Bytes("000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000"),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli, storage := getMocks(t)

			s := scan.NewScanner(cli, storage)

			cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: tt.fromBlock,
				ToBlock:   tt.toBlock,
				Addresses: []common.Address{tt.address},
			}).Return(tt.eventLogs, nil)

			events, err := s.ScanEvents(context.Background(), tt.fromBlock, tt.toBlock, tt.address)
			if err != nil {
				t.Fatalf("error occurred when scanning events %v", err.Error())
			}

			switch eventType := events[0].(type) {
			case scan.EventTransfer:
				_, ok := events[0].(scan.EventTransfer)
				if !ok {
					t.Fatal("error parsing event to EventTransfer type")
				}
			case scan.EventApproval:
				_, ok := events[0].(scan.EventApproval)
				if !ok {
					t.Fatal("error parsing event to EventApproval type")
				}
			case scan.EventApprovalForAll:
				_, ok := events[0].(scan.EventApprovalForAll)
				if !ok {
					t.Fatal("error parsing event to EventApprovalForAll type")
				}
			case scan.EventNewCollecion:
				_, ok := events[0].(scan.EventNewCollecion)
				if !ok {
					t.Fatal("error parsing event to EventNewCollection type")
				}
			case scan.EventMintedWithExternalURI:
				_, ok := events[0].(scan.EventMintedWithExternalURI)
				if !ok {
					t.Fatal("error parsing event to EventMintedWithExternalURI type")
				}
			case scan.EventEvolvedWithExternalURI:
				_, ok := events[0].(scan.EventEvolvedWithExternalURI)
				if !ok {
					t.Fatal("error parsing event to EventEvolvedWithExternalURI type")
				}
			default:
				t.Fatalf("unknown event: %v", eventType)
			}
		})
	}
}

func TestScanEvents(t *testing.T) {
	t.Parallel()
	t.Run("it should only parse Transfer, Approve and ApproveForAllEvents", func(t *testing.T) {
		t.Parallel()

		cli, storage := getMocks(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")

		s := scan.NewScanner(cli, storage)

		eventLogs := []types.Log{
			{
				Topics: []common.Hash{
					common.HexToHash(approveForAllEventHash),
					common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
					common.HexToHash("0x0000000000000000000000001e0049783f008a0085193e00003d00cd54003c71"),
				},
				Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
			},
			{
				Topics: []common.Hash{
					common.HexToHash(approveEventHash),
					common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
					common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
				},
				Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
			},
			{
				Topics: []common.Hash{
					common.HexToHash(transferEventHash),
					common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
					common.HexToHash("0x00000000000000000000000066666f58de1bcd762a5e5c5aff9cc3c906d66666"),
					common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
				},
			},
			{
				Topics: []common.Hash{
					common.HexToHash(transferEventHash),
					common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
					common.HexToHash("0x00000000000000000000000066666f58de1bcd762a5e5c5aff9cc3c906d66666"),
					common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
				},
			},
			{
				Topics: []common.Hash{
					// Event hash is not included in the list of events to be parsed
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
					common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
					common.HexToHash("0x00000000000000000000000066666f58de1bcd762a5e5c5aff9cc3c906d66666"),
					common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
				},
			},
		}

		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
		}).Return(eventLogs, nil)

		events, err := s.ScanEvents(context.Background(), fromBlock, toBlock, address)
		if err != nil {
			t.Fatalf("error occurred when scanning events %v", err.Error())
		}

		if len(events) != 4 {
			t.Fatalf("error scanning events: %v events exepected, got %v", 4, len(events))
		}
	})

	t.Run("it returns when there are no events", func(t *testing.T) {
		t.Parallel()

		cli, storage := getMocks(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")

		s := scan.NewScanner(cli, storage)

		eventLogs := []types.Log{}

		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
		}).Return(eventLogs, nil)

		events, err := s.ScanEvents(context.Background(), fromBlock, toBlock, address)
		if err != nil {
			t.Fatalf("nil error expected, got %v", err)
		}
		if events != nil {
			t.Fatalf("nil events expected, got %v", events)
		}
	})

	t.Run("scan raises error", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			error error
		}{
			{
				name:  "it raises an error when call to blockchain fails",
				error: fmt.Errorf("error filtering events"),
			},
			{
				name:  "it raises an error when storage fails reading contracts",
				error: fmt.Errorf("error reading from storage"),
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				cli, storage := getMocks(t)

				fromBlock := big.NewInt(0)
				toBlock := big.NewInt(100)
				address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")

				s := scan.NewScanner(cli, storage)

				cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
					FromBlock: fromBlock,
					ToBlock:   toBlock,
					Addresses: []common.Address{address},
				}).Return(nil, tt.error)

				_, err := s.ScanEvents(context.Background(), fromBlock, toBlock, address)
				if err == nil {
					t.Fatalf("got nil error, expected %v", tt.error.Error())
				}
			})
		}
	})
}

func TestScanNewUniversalEventsErr(t *testing.T) {
	t.Parallel()
	address := common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")
	fromBlock := big.NewInt(0)
	toBlock := big.NewInt(100)

	tests := []struct {
		readAllError            error
		filterLogsError         error
		name                    string
		events                  []types.Log
		filterLogsExpectedTimes int
	}{
		{
			name:                    "error filtering logs",
			events:                  nil,
			filterLogsError:         fmt.Errorf("error filtering logs"),
			filterLogsExpectedTimes: 1,
		},
		{
			name:                    "discovery fails because reading from storage fails",
			events:                  nil,
			filterLogsExpectedTimes: 0,
			filterLogsError:         nil,
			readAllError:            fmt.Errorf("error while reading contracts from storage"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cli, storage := getMocks(t)
			s := scan.NewScanner(cli, storage)

			storage.EXPECT().ReadAll(context.Background()).
				Return(nil, tt.readAllError).
				Times(1)
			cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: fromBlock,
				ToBlock:   toBlock,
				Addresses: []common.Address{address},
			}).Return(tt.events, tt.filterLogsError).Times(tt.filterLogsExpectedTimes)

			_, err := s.ScanNewUniversalEvents(context.Background(), fromBlock, toBlock, address)
			if err == nil {
				t.Fatalf("got no error, %v expected", tt.name)
			}
		})
	}
}

func TestScanNewUniversalEvents(t *testing.T) {
	t.Parallel()
	fromBlock := big.NewInt(0)
	toBlock := big.NewInt(100)

	tests := []struct {
		name                    string
		events                  []types.Log
		expectedContractsParsed int
	}{
		{
			name: "find and store one contract",
			events: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(newERC721UniversalEventHash),
					},
					Data: common.Hex2Bytes("00000000000000000000000026cb70039fe1bd36b4659858d4c4d0cbcafd743a0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000001765766f636861696e312f636f6c6c656374696f6e49642f000000000000000000"),
				},
			},
			expectedContractsParsed: 1,
		},
		{
			name: "other event types found",
			events: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(approveEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
						common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
					},
					Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
			expectedContractsParsed: 0,
		},
		{
			name: "anonymous event found",
			events: []types.Log{
				{
					Topics: []common.Hash{},
					Data:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
				},
			},
			expectedContractsParsed: 0,
		},
		{
			name:                    "no events found",
			events:                  []types.Log{},
			expectedContractsParsed: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cli, storage := getMocks(t)
			s := scan.NewScanner(cli, storage)

			cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: fromBlock,
				ToBlock:   toBlock,
			}).Return(tt.events, nil).Times(1)

			contracts, err := s.ScanNewUniversalEvents(context.Background(), fromBlock, toBlock)
			if err != nil {
				t.Fatalf("got error %v when no error was expected", err)
			}
			if len(contracts) != tt.expectedContractsParsed {
				t.Fatalf("got %d contracts, %d expected", len(contracts), tt.expectedContractsParsed)
			}
		})
	}
}

func TestScanNewUniversalEventsDiscovery(t *testing.T) {
	t.Parallel()
	fromBlock := big.NewInt(0)
	toBlock := big.NewInt(100)
	userDefinedContracts := []common.Address{
		common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
		common.HexToAddress("0x59b6c3b079a2af72ece16065c85b95ac2405459b"),
	}
	tests := []struct {
		name                    string
		storedContracts         []scan.ERC721UniversalContract
		filterLogsExpectedTimes int
	}{
		{
			name:                    "discovery returns false",
			filterLogsExpectedTimes: 0,
			storedContracts: []scan.ERC721UniversalContract{
				{
					Address: userDefinedContracts[0],
					Block:   fromBlock.Uint64(),
					BaseURI: "evochain1/collectionId/",
				},
				{
					Address: userDefinedContracts[1],
					Block:   fromBlock.Uint64(),
					BaseURI: "evochain2/collectionId/",
				},
			},
		},
		{
			name:                    "discovery returns true",
			filterLogsExpectedTimes: 1,
			storedContracts: []scan.ERC721UniversalContract{
				{
					Address: userDefinedContracts[0],
					Block:   fromBlock.Uint64(),
					BaseURI: "evochain1/collectionId/",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cli, storage := getMocks(t)
			s := scan.NewScanner(cli, storage)
			storage.EXPECT().ReadAll(context.Background()).Return(tt.storedContracts, nil).Times(1)
			cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: fromBlock,
				ToBlock:   toBlock,
				Addresses: userDefinedContracts,
			}).Return(nil, nil).Times(tt.filterLogsExpectedTimes)
			_, err := s.ScanNewUniversalEvents(context.Background(), fromBlock, toBlock, userDefinedContracts...)
			if err != nil {
				t.Fatalf("got error %v, expected nil", err)
			}
		})
	}
}

func getMocks(t *testing.T) (*mock.MockEthClient, *mock.MockStorage) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockEthClient(ctrl), mock.NewMockStorage(ctrl)
}
