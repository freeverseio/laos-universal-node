package scan_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/scan"
	"github.com/freeverseio/laos-universal-node/scan/mock"
	"go.uber.org/mock/gomock"
)

const (
	transferEventHash      = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	approveEventHash       = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
	approveForAllEventHash = "0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31"
)

func TestScanEvents(t *testing.T) {
	t.Run("it should parse Transfer events", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		climock := mock.NewMockEthClient(ctrl)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		contract := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")

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

		climock.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{contract},
		}).Return(eventLogs, nil)

		events, err := scan.ScanEvents(climock, contract, fromBlock, toBlock)
		if err != nil {
			t.Fatalf("error occured when scanning events %v", err.Error())
		}

		_, ok := events[0].(scan.EventTransfer)
		if !ok {
			t.Fatal("error parsing event to EventApproval type")
		}
	})
	t.Run("it should parse Approval events", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		climock := mock.NewMockEthClient(ctrl)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		contract := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")

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

		climock.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{contract},
		}).Return(eventLogs, nil)

		events, err := scan.ScanEvents(climock, contract, fromBlock, toBlock)
		if err != nil {
			t.Fatalf("error occured when scanning events %v", err.Error())
		}

		_, ok := events[0].(scan.EventApproval)
		if !ok {
			t.Fatal("error parsing event to EventApproval type")
		}
	})
	t.Run("it should parse ApprovalForAll events", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		climock := mock.NewMockEthClient(ctrl)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		contract := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")

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

		climock.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{contract},
		}).Return(eventLogs, nil)

		events, err := scan.ScanEvents(climock, contract, fromBlock, toBlock)
		if err != nil {
			t.Fatalf("error occured when scanning events %v", err.Error())
		}

		_, ok := events[0].(scan.EventApprovalForAll)
		if !ok {
			t.Fatal("error parsing event to EventApprovalForAll type")
		}
	})
	t.Run("it should only parse Transfer, Approve and ApproveForAllEvents", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		climock := mock.NewMockEthClient(ctrl)

		fromBlock := big.NewInt(0)
		toBlock := big.NewInt(100)
		contract := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")

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

		climock.EXPECT().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{contract},
		}).Return(eventLogs, nil)

		events, err := scan.ScanEvents(climock, contract, fromBlock, toBlock)
		if err != nil {
			t.Fatalf("error occured when scanning events %v", err.Error())
		}

		if len(events) != 4 {
			t.Fatalf("error scanning events: %v events exepected, got %v", 4, len(events))
		}
	})
}
