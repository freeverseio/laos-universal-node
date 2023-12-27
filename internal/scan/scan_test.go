package scan_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
	"go.uber.org/mock/gomock"
)

const (
	transferEventHash               = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	approveEventHash                = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
	approveForAllEventHash          = "0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31"
	newERC721UniversalEventHash     = "0x74b81bc88402765a52dad72d3d893684f472a679558f3641500e0ee14924a10a"
	newCollectionEventHash          = "0x5b84d9550adb7000df7bee717735ecd3af48ea3f66c6886d52e8227548fb228c"
	mintedWithExternalURIEventHash  = "0xa7135052b348b0b4e9943bae82d8ef1c5ac225e594ef4271d12f0744cfc98348"
	evolvedWithExternalURIEventHash = "0xdde18ad2fe10c12a694de65b920c02b851c382cf63115967ea6f7098902fa1c8"
)

func TestParseEvents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		fromBlock           *big.Int
		toBlock             *big.Int
		address             common.Address
		contracts           []model.ERC721UniversalContract
		eventLogs           []types.Log
		headerByNumberTimes int
	}{
		{
			name:      "it should parse Transfer events",
			fromBlock: big.NewInt(0),
			toBlock:   big.NewInt(100),
			address:   common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			contracts: []model.ERC721UniversalContract{
				{
					Address:           common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
					CollectionAddress: common.HexToAddress("johndoe/collection"),
				},
			},
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(transferEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x00000000000000000000000066666f58de1bcd762a5e5c5aff9cc3c906d66666"),
						common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
					},
					BlockNumber: 100,
				},
			},
		},
		{
			name:      "it should parse Approval events",
			fromBlock: big.NewInt(0),
			toBlock:   big.NewInt(90),
			address:   common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			contracts: []model.ERC721UniversalContract{
				{
					Address:           common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
					CollectionAddress: common.HexToAddress("johndoe/collection"),
				},
			},
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(approveEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
						common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
					},
					Data:        common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
					BlockNumber: 90,
				},
			},
		},
		{
			name:      "it should parse ApprovalForAll events",
			fromBlock: big.NewInt(0),
			toBlock:   big.NewInt(100),
			address:   common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			contracts: []model.ERC721UniversalContract{
				{
					Address:           common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
					CollectionAddress: common.HexToAddress("johndoe/collection"),
				},
			},
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(approveForAllEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x0000000000000000000000001e0049783f008a0085193e00003d00cd54003c71"),
					},
					Data:        common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
					BlockNumber: 100,
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
						common.HexToHash(" 0x000000000000000000000000c112bde959080c5b46e73749e3e170f47123e85a"),
					},
					Data:        common.Hex2Bytes("000000000000000000000000fffffffffffffffffffffffe00000000000000e5"),
					BlockNumber: 100,
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
						common.HexToHash("0x000000000000000000000000c112bde959080c5b46e73749e3e170f47123e85a"),
					},
					Data:        common.Hex2Bytes("00000000000000000000000000000000000000003d5b1313de887a00000000003d5b1313de887a0000000000c112bde959080c5b46e73749e3e170f47123e85a0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000002e516d4e5247426d7272724862754b4558375354544d326f68325077324d757438674863537048706a367a7a637375000000000000000000000000000000000000"),
					BlockNumber: 100,
				},
			},
			headerByNumberTimes: 1,
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
						common.HexToHash("0x000000000000000000000001684bc8f81250ad3f7f930b27586b799a4dda957b"),
					},
					Data:        common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001674657374696e67315f63616c6164616e5f31376e6f7600000000000000000000"),
					BlockNumber: 100,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli := getMockEthClient(t)

			s := scan.NewScanner(cli)

			cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: tt.fromBlock,
				ToBlock:   tt.toBlock,
				Addresses: []common.Address{tt.address},
			}).Return(tt.eventLogs, nil)

			cli.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(tt.eventLogs[0].BlockNumber))).Return(&types.Header{Time: uint64(time.Now().Unix())}, nil).Times(tt.headerByNumberTimes)

			events, err := s.ScanEvents(context.Background(), tt.fromBlock, tt.toBlock, []string{tt.address.String()})
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

func TestScanOnlyValidEvents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		eventLogs      []types.Log
		expectedEvents int
	}{
		{
			name: "it only returns Transfer, Approval and ApprovalForAll",
			eventLogs: []types.Log{
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
						// Event hash is not included in the list of events to be parsed
						common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x00000000000000000000000066666f58de1bcd762a5e5c5aff9cc3c906d66666"),
						common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
					},
				},
			},
			expectedEvents: 3,
		},
		{
			name: "it does not parse Transfer with unexpected topics length",
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(transferEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x00000000000000000000000066666f58de1bcd762a5e5c5aff9cc3c906d66666"),
					},
				},
			},
			expectedEvents: 0,
		},
		{
			name: "it does not parse Approval with unexpected topics length",
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(approveEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
						common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
					},
					Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
			expectedEvents: 0,
		},
		{
			name: "it does not parse ApprovalForAll with unexpected topics length",
			eventLogs: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(approveForAllEventHash),
						common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
					},
					Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
			expectedEvents: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cli := getMockEthClient(t)

			fromBlock := big.NewInt(0)
			toBlock := big.NewInt(100)
			address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
			contracts := []string{address.String()}
			s := scan.NewScanner(cli)

			cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: fromBlock,
				ToBlock:   toBlock,
				Addresses: []common.Address{address},
			}).Return(tt.eventLogs, nil)

			events, err := s.ScanEvents(context.Background(), fromBlock, toBlock, contracts)
			if err != nil {
				t.Fatalf("error occurred when scanning events %v", err.Error())
			}

			if len(events) != tt.expectedEvents {
				t.Fatalf("error scanning events: %v events exepected, got %v", 4, len(events))
			}
		})
	}
}

func TestScanEvents(t *testing.T) {
	t.Parallel()

	t.Run("it returns when there are no events", func(t *testing.T) {
		t.Parallel()

		cli := getMockEthClient(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
		contracts := []string{
			address.String(),
		}

		s := scan.NewScanner(cli)

		eventLogs := []types.Log{}

		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
		}).Return(eventLogs, nil)

		events, err := s.ScanEvents(context.Background(), fromBlock, toBlock, contracts)
		if err != nil {
			t.Fatalf("nil error expected, got %v", err)
		}

		if events != nil {
			t.Fatalf("nil events expected, got %v", events)
		}
	})
	t.Run("it does not parse events when blockchain does not return any event", func(t *testing.T) {
		t.Parallel()

		cli := getMockEthClient(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)

		s := scan.NewScanner(cli)
		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: nil,
		}).Return(nil, nil)

		events, err := s.ScanEvents(context.Background(), fromBlock, toBlock, []string{})
		if err != nil {
			t.Errorf("got error %s when scanning events while no error was expected", err.Error())
		}

		if len(events) > 0 {
			t.Fatalf("got events %v when no events where expected", events)
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
				cli := getMockEthClient(t)

				fromBlock := big.NewInt(0)
				toBlock := big.NewInt(100)
				address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
				contracts := []string{
					address.String(),
				}

				s := scan.NewScanner(cli)

				cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
					FromBlock: fromBlock,
					ToBlock:   toBlock,
					Addresses: []common.Address{address},
				}).Return(nil, tt.error)

				_, err := s.ScanEvents(context.Background(), fromBlock, toBlock, contracts)
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
	contract := model.ERC721UniversalContract{
		Address:           address,
		CollectionAddress: common.HexToAddress("evochain1/collectionId/"),
	}

	tests := []struct {
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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cli := getMockEthClient(t)
			s := scan.NewScanner(cli, contract.Address.String())

			cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: fromBlock,
				ToBlock:   toBlock,
				Addresses: []common.Address{address},
			}).Return(tt.events, tt.filterLogsError).Times(tt.filterLogsExpectedTimes)

			_, err := s.ScanNewUniversalEvents(context.Background(), fromBlock, toBlock)
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
					Data: common.Hex2Bytes("000000000000000000000000c3dd09d5387fa0ab798e0adc152d15b8d1a299df0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008c756c6f633a2f2f476c6f62616c436f6e73656e7375732867656e292f50617261636861696e2832393030292f4163636f756e744b6579323028307830303030303030303030303030303030303030303030303031306663346161303133356166376263356434386665373564613332646262353262643936333162292f47656e6572616c4b657928363636290000000000000000000000000000000000000000"),
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
			cli := getMockEthClient(t)
			s := scan.NewScanner(cli)

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

func TestValidBaseURI(t *testing.T) {
	t.Parallel()
	t.Run(`extract collection address from baseURI`, func(t *testing.T) {
		t.Parallel()

		baseURI := "uloc://GlobalConsensus(gen)/Parachain(2900)/AccountKey20(0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b)/GeneralKey(666)"
		expectedCollectionAddress := common.HexToAddress("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b")

		e := scan.EventNewERC721Universal{
			BaseURI: baseURI,
		}

		address, err := e.CollectionAddress()
		if err != nil {
			t.Errorf("got %s error, nil expected", err.Error())
		}

		if address != expectedCollectionAddress {
			t.Fatalf("got %d collection address, %d expected", address, expectedCollectionAddress)
		}
	})
	t.Run(`extract global consensus from baseURI`, func(t *testing.T) {
		t.Parallel()

		baseURI := "uloc://GlobalConsensus(3)/Parachain(0)/AccountKey20(0x0)/GeneralKey(1)"
		expectedGlobalConsensus := "3"

		e := scan.EventNewERC721Universal{
			BaseURI: baseURI,
		}

		globalConsensus, err := e.GlobalConsensus()
		if err != nil {
			t.Errorf("got %s error, nil expected", err.Error())
		}

		if globalConsensus != expectedGlobalConsensus {
			t.Fatalf("got %s global consensus, %s expected", globalConsensus, expectedGlobalConsensus)
		}
	})
	t.Run(`extract parachain from baseURI`, func(t *testing.T) {
		t.Parallel()

		baseURI := "uloc://GlobalConsensus(gen)/Parachain(3336)/AccountKey20(0x0)/GeneralKey(1)"
		expectedParachain := uint64(3336)

		e := scan.EventNewERC721Universal{
			BaseURI: baseURI,
		}

		parachain, err := e.Parachain()
		if err != nil {
			t.Errorf("got %s error, nil expected", err.Error())
		}

		if parachain != expectedParachain {
			t.Fatalf("got %d parachain, %d expected", parachain, expectedParachain)
		}
	})
}

func TestInvalidBaseURI(t *testing.T) {
	t.Parallel()
	t.Run(`fails when AccountKey does not exist in base uri`, func(t *testing.T) {
		t.Parallel()

		baseURI := "uloc://GlobalConsensus(gen)/Parachain(0)/GeneralKey(666)"
		expectedError := fmt.Errorf("no collection address found in base URI: %s", baseURI)

		e := scan.EventNewERC721Universal{
			BaseURI: baseURI,
		}

		_, err := e.CollectionAddress()
		if err == nil {
			t.Fatalf("got nil error, %s expected", expectedError.Error())
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("got %s error, %s expected", err.Error(), expectedError.Error())
		}
	})
	t.Run(`fails when GlobalConsensus does not exist in base uri`, func(t *testing.T) {
		t.Parallel()

		baseURI := "uloc://Parachain(0)/AccountKey20(0x0)/GeneralKey(666)"
		expectedError := fmt.Errorf("no global consensus ID found in base URI: %s", baseURI)

		e := scan.EventNewERC721Universal{
			BaseURI: baseURI,
		}

		_, err := e.GlobalConsensus()
		if err == nil {
			t.Fatalf("got nil error, %s expected", expectedError.Error())
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("got %s error, %s expected", err.Error(), expectedError.Error())
		}
	})
	t.Run(`fails when Parachain does not exist in base uri`, func(t *testing.T) {
		t.Parallel()

		baseURI := "uloc://GlobalConsensus(3)/AccountKey20(0x0)/GeneralKey(666)"
		expectedError := fmt.Errorf("no parachain ID found in base URI: %s", baseURI)

		e := scan.EventNewERC721Universal{
			BaseURI: baseURI,
		}

		_, err := e.Parachain()
		if err == nil {
			t.Fatalf("got nil error, %s expected", expectedError.Error())
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("got %s error, %s expected", err.Error(), expectedError.Error())
		}
	})
	t.Run(`fails when Parachain is malformed in base uri`, func(t *testing.T) {
		t.Parallel()

		baseURI := "uloc://GlobalConsensus(3)/Parachain(kkk)/AccountKey20(0x0)/GeneralKey(666)"
		expectedError := fmt.Errorf(`error parsing parachain value to uint: strconv.ParseUint: parsing "kkk": invalid syntax`)

		e := scan.EventNewERC721Universal{
			BaseURI: baseURI,
		}

		_, err := e.Parachain()
		if err == nil {
			t.Fatalf("got nil error, %s expected", expectedError.Error())
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("got %s error, %s expected", err.Error(), expectedError.Error())
		}
	})
}

func getMockEthClient(t *testing.T) *mock.MockEthClient {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockEthClient(ctrl)
}
