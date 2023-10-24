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
	transferEventHash                   = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	approveEventHash                    = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
	approveForAllEventHash              = "0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31"
	newERC721BridgelessMintingEventHash = "0x821a490a0b4f9fa6744efb226f24ce4c3917ff2fca72c1750947d75a99254610"
)

func TestScanEvents(t *testing.T) {
	t.Run("it returns when there are no events", func(t *testing.T) {
		cli, storage := getMocks(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
		contracts := []scan.ERC721BridgelessContract{
			{
				Address: address,
				Block:   fromBlock.Uint64(),
				BaseURI: "johndoe/collection",
			},
		}

		s := scan.NewScanner(cli, address, storage)

		eventLogs := []types.Log{}

		storage.EXPECT().ReadAll(context.Background()).Return(contracts, nil).Times(1)
		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
		}).Return(eventLogs, nil)

		events, err := s.ScanEvents(context.Background(), fromBlock, toBlock)
		if err != nil {
			t.Fatalf("nil error expected, got %v", err)
		}
		if events != nil {
			t.Fatalf("nil events expected, got %v", events)
		}
	})
	t.Run("it should parse Transfer events", func(t *testing.T) {
		cli, storage := getMocks(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
		contracts := []scan.ERC721BridgelessContract{
			{
				Address: address,
				Block:   fromBlock.Uint64(),
				BaseURI: "johndoe/collection",
			},
		}

		s := scan.NewScanner(cli, address, storage)

		eventLogs := []types.Log{
			{
				Topics: []common.Hash{
					common.HexToHash(transferEventHash),
					common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
					common.HexToHash("0x00000000000000000000000066666f58de1bcd762a5e5c5aff9cc3c906d66666"),
					common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
				},
			},
		}

		storage.EXPECT().ReadAll(context.Background()).Return(contracts, nil).Times(1)
		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
		}).Return(eventLogs, nil)

		events, err := s.ScanEvents(context.Background(), fromBlock, toBlock)
		if err != nil {
			t.Fatalf("error occurred when scanning events %v", err.Error())
		}

		_, ok := events[0].(scan.EventTransfer)
		if !ok {
			t.Fatal("error parsing event to EventTransfer type")
		}
	})
	t.Run("it should parse Approval events", func(t *testing.T) {
		cli, storage := getMocks(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
		contracts := []scan.ERC721BridgelessContract{
			{
				Address: address,
				Block:   fromBlock.Uint64(),
				BaseURI: "johndoe/collection",
			},
		}

		s := scan.NewScanner(cli, address, storage)

		eventLogs := []types.Log{
			{
				Topics: []common.Hash{
					common.HexToHash(approveEventHash),
					common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
					common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000009f4"),
				},
				Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
			},
		}

		storage.EXPECT().ReadAll(context.Background()).Return(contracts, nil).Times(1)
		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
		}).Return(eventLogs, nil)

		events, err := s.ScanEvents(context.Background(), fromBlock, toBlock)
		if err != nil {
			t.Fatalf("error occurred when scanning events %v", err.Error())
		}

		_, ok := events[0].(scan.EventApproval)
		if !ok {
			t.Fatal("error parsing event to EventApproval type")
		}
	})
	t.Run("it should parse ApprovalForAll events", func(t *testing.T) {
		cli, storage := getMocks(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
		contracts := []scan.ERC721BridgelessContract{
			{
				Address: address,
				Block:   fromBlock.Uint64(),
				BaseURI: "johndoe/collection",
			},
		}
		s := scan.NewScanner(cli, address, storage)

		eventLogs := []types.Log{
			{
				Topics: []common.Hash{
					common.HexToHash(approveForAllEventHash),
					common.HexToHash("0x00000000000000000000000010fc4aa0135af7bc5d48fe75da32dbb52bd9631b"),
					common.HexToHash("0x0000000000000000000000001e0049783f008a0085193e00003d00cd54003c71"),
				},
				Data: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
			},
		}

		storage.EXPECT().ReadAll(context.Background()).Return(contracts, nil).Times(1)
		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
		}).Return(eventLogs, nil)

		events, err := s.ScanEvents(context.Background(), fromBlock, toBlock)
		if err != nil {
			t.Fatalf("error occurred when scanning events %v", err.Error())
		}

		_, ok := events[0].(scan.EventApprovalForAll)
		if !ok {
			t.Fatal("error parsing event to EventApprovalForAll type")
		}
	})
	t.Run("it should only parse Transfer, Approve and ApproveForAllEvents", func(t *testing.T) {
		cli, storage := getMocks(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
		contracts := []scan.ERC721BridgelessContract{
			{
				Address: address,
				Block:   fromBlock.Uint64(),
				BaseURI: "johndoe/collection",
			},
		}

		s := scan.NewScanner(cli, address, storage)

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

		storage.EXPECT().ReadAll(context.Background()).Return(contracts, nil).Times(1)
		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
		}).Return(eventLogs, nil)

		events, err := s.ScanEvents(context.Background(), fromBlock, toBlock)
		if err != nil {
			t.Fatalf("error occurred when scanning events %v", err.Error())
		}

		if len(events) != 4 {
			t.Fatalf("error scanning events: %v events exepected, got %v", 4, len(events))
		}
	})
	t.Run("it raises an error when call to blockchain fails", func(t *testing.T) {
		cli, storage := getMocks(t)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
		contracts := []scan.ERC721BridgelessContract{
			{
				Address: address,
				Block:   fromBlock.Uint64(),
				BaseURI: "johndoe/collection",
			},
		}

		s := scan.NewScanner(cli, address, storage)

		storage.EXPECT().ReadAll(context.Background()).Return(contracts, nil).Times(1)
		cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
		}).Return(nil, fmt.Errorf("error filtering events"))

		_, err := s.ScanEvents(context.Background(), fromBlock, toBlock)
		if err == nil {
			t.Fatal("error expected, got nil")
		}
	})
}

func TestScanNewBridgelessMintingEventsErr(t *testing.T) {
	t.Parallel()
	cli, storage := getMocks(t)
	address := common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")
	fromBlock := big.NewInt(0)
	toBlock := big.NewInt(100)
	s := scan.NewScanner(cli, address, storage)
	contract := scan.ERC721BridgelessContract{
		Address: address,
		Block:   fromBlock.Uint64(),
		BaseURI: "evochain1/collectionId/",
	}

	tests := []struct {
		name                 string
		filterLogsError      error
		events               []types.Log
		storageExpectedTimes int
	}{
		{
			name: "error storing contracts",
			events: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(newERC721BridgelessMintingEventHash),
					},
					Data: common.Hex2Bytes("00000000000000000000000026cb70039fe1bd36b4659858d4c4d0cbcafd743a0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000001765766f636861696e312f636f6c6c656374696f6e49642f000000000000000000"),
				},
			},
			storageExpectedTimes: 1,
			filterLogsError:      nil,
		},
		{
			name:                 "error filtering logs",
			events:               nil,
			storageExpectedTimes: 0,
			filterLogsError:      fmt.Errorf("error filtering logs"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage.EXPECT().Store(context.Background(), contract).
				Return(fmt.Errorf("error storing contracts")).
				Times(tt.storageExpectedTimes)
			cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: fromBlock,
				ToBlock:   toBlock,
			}).Return(tt.events, tt.filterLogsError).Times(1)

			err := s.ScanNewBridgelessMintingEvents(context.Background(), fromBlock, toBlock)
			if err == nil {
				t.Fatalf("got no error when error %v was expected", tt.name)
			}
		})
	}
}

func TestScanNewBridgelessMintingEvents(t *testing.T) {
	t.Parallel()
	cli, storage := getMocks(t)
	address := common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")
	fromBlock := big.NewInt(0)
	toBlock := big.NewInt(100)
	s := scan.NewScanner(cli, address, storage)
	contract := scan.ERC721BridgelessContract{
		Address: address,
		Block:   fromBlock.Uint64(),
		BaseURI: "evochain1/collectionId/",
	}

	tests := []struct {
		name                 string
		events               []types.Log
		storageExpectedTimes int
	}{
		{
			name: "find and store one contract",
			events: []types.Log{
				{
					Topics: []common.Hash{
						common.HexToHash(newERC721BridgelessMintingEventHash),
					},
					Data: common.Hex2Bytes("00000000000000000000000026cb70039fe1bd36b4659858d4c4d0cbcafd743a0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000001765766f636861696e312f636f6c6c656374696f6e49642f000000000000000000"),
				},
			},
			storageExpectedTimes: 1,
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
			storageExpectedTimes: 0,
		},
		{
			name: "anonymous event found",
			events: []types.Log{
				{
					Topics: []common.Hash{},
					Data:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
				},
			},
			storageExpectedTimes: 0,
		},
		{
			name:                 "no events found",
			events:               []types.Log{},
			storageExpectedTimes: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage.EXPECT().Store(context.Background(), contract).Return(nil).Times(tt.storageExpectedTimes)
			cli.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: fromBlock,
				ToBlock:   toBlock,
			}).Return(tt.events, nil).Times(1)

			err := s.ScanNewBridgelessMintingEvents(context.Background(), fromBlock, toBlock)
			if err != nil {
				t.Fatalf("got error %v when no error was expected", err)
			}
		})
	}
}

func getMocks(t *testing.T) (*mock.MockEthClient, *mock.MockStorage) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockEthClient(ctrl), mock.NewMockStorage(ctrl)
}
