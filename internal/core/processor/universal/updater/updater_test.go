package updater_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/mock/gomock"

	uUpdater "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/updater"
	mockClient "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	mockScan "github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/state/tree/account"
)

func TestGetModelTransferEvents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		transferEvents map[string][]model.ERC721Transfer
		contract       string
		startingBlock  uint64
		lastBlock      uint64
		eventBlock     uint64
		eventTime      uint64
		expectedError  error
	}{
		{
			name:           "get transfer events",
			transferEvents: map[string][]model.ERC721Transfer{"0x000005555": getERC721TransferEvents("0x000005555", 351, 351)},
			startingBlock:  300,
			lastBlock:      360,
			eventBlock:     351,
			eventTime:      351,
			contract:       "0x000005555",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.TODO()
			_, client, scanner := createMocks(t)

			updater := uUpdater.New(client, scanner)

			eventFromScanner := createScanEventTransfer(tt.contract, tt.eventBlock)
			scanner.EXPECT().
				ScanEvents(ctx, big.NewInt(int64(tt.startingBlock)), big.NewInt(int64(tt.lastBlock)), []string{tt.contract}).
				Return([]scan.Event{eventFromScanner}, nil)

			events, err := updater.GetModelTransferEvents(ctx, tt.startingBlock, tt.lastBlock, []string{tt.contract})
			assertError(t, tt.expectedError, err)

			if len(events[351][common.HexToAddress(tt.contract).String()]) != 1 {
				t.Fatalf(`wrong number of events got %v expected %v"`, len(events[351][tt.contract]), 1)
			}
		})
	}
}

func TestGetEvoEvents(t *testing.T) {
	t.Parallel()

	t.Run("get minted events", func(t *testing.T) {
		t.Parallel()

		tx, _, _ := createMocks(t)

		events := getMockMintedEvents(352, 352)
		tx.EXPECT().GetCollectionAddress("0x000005555").Return(common.HexToAddress("0x4444"), nil)
		tx.EXPECT().AccountData(common.HexToAddress("0x000005555")).Return(&account.AccountData{
			LastProcessedEvoBlock: 351,
		}, nil)
		tx.EXPECT().GetNextEvoEventBlock(common.HexToAddress("0x4444").String(), uint64(351)).Return(uint64(352), nil)
		tx.EXPECT().GetMintedWithExternalURIEvents(common.HexToAddress("0x4444").String(), uint64(352)).Return(events, nil)

		evoBlock, events, err := uUpdater.GetEvoEvents(tx, "0x000005555", uint64(352))
		assertError(t, err, nil)
		if evoBlock != uint64(352) {
			t.Fatalf(`wrong evo block got %v expected %v"`, evoBlock, uint64(352))
		}
		if len(events) != 1 {
			t.Fatal("wrong number of events")
		}
	})
}

func TestUpdateContract(t *testing.T) {
	t.Parallel()

	t.Run("get minted events", func(t *testing.T) {
		t.Parallel()

		tx, _, _ := createMocks(t)

		evoEvents := getMockMintedEvents(352, 352)
		events := getERC721TransferEvents("0x000005555", 353, 353)

		tx.EXPECT().LoadContractTrees(common.HexToAddress("0x000005555")).Return(nil)
		tx.EXPECT().Mint(common.HexToAddress("0x000005555"), &evoEvents[0]).Return(nil)
		tx.EXPECT().Transfer(common.HexToAddress("0x000005555"), &events[0]).Return(nil)
		tx.EXPECT().UpdateContractState(common.HexToAddress("0x000005555"), uint64(352)).Return(nil)

		err := uUpdater.UpdateContract(tx, "0x000005555", evoEvents, events, uint64(353), uint64(352))
		assertError(t, err, nil)
	})
}

func TestGetBlockTimestampsParallel(t *testing.T) {
	t.Parallel()

	t.Run("get block timestamps parallel", func(t *testing.T) {
		t.Parallel()

		_, client, _ := createMocks(t)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(10)).Return(&types.Header{
			Time:   10,
			Number: big.NewInt(int64(10)),
		}, nil)
		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(11)).Return(&types.Header{
			Time:   11,
			Number: big.NewInt(int64(11)),
		}, nil)
		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(12)).Return(&types.Header{
			Time:   12,
			Number: big.NewInt(int64(12)),
		}, nil)

		blockTimestamps, err := uUpdater.GetBlockTimestampsParallel(context.Background(), client, 10, 12)

		assertError(t, err, nil)
		if len(blockTimestamps) != 3 {
			t.Fatalf(`wrong number of timestamps got %v expected %v"`, len(blockTimestamps), 3)
		}

		if blockTimestamps[10] != 10 {
			t.Fatalf(`wrong timestamp got %v expected %v"`, blockTimestamps[10], 10)
		}
		if blockTimestamps[11] != 11 {
			t.Fatalf(`wrong timestamp got %v expected %v"`, blockTimestamps[11], 11)
		}
		if blockTimestamps[12] != 12 {
			t.Fatalf(`wrong timestamp got %v expected %v"`, blockTimestamps[12], 12)
		}
	})

	t.Run("get block timestamps parallel, one call returns error", func(t *testing.T) {
		t.Parallel()
		_, client, _ := createMocks(t)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(10)).Return(&types.Header{
			Time:   10,
			Number: big.NewInt(int64(10)),
		}, nil)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(11)).Return(
			nil,
			fmt.Errorf("some error"))

		_, err := uUpdater.GetBlockTimestampsParallel(context.Background(), client, 10, 11)
		assertError(t, err, fmt.Errorf("some error"))
	})
}

func createMocks(t *testing.T) (*mockTx.MockTx, *mockClient.MockEthClient, *mockScan.MockScanner) {
	ctrl := gomock.NewController(t)
	return mockTx.NewMockTx(ctrl), mockClient.NewMockEthClient(ctrl), mockScan.NewMockScanner(ctrl)
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

func getERC721TransferEvents(contract string, blockNumber, timestamp uint64) []model.ERC721Transfer {
	return []model.ERC721Transfer{
		{
			From:        common.HexToAddress("0x01"),
			To:          common.HexToAddress("0x02"),
			TokenId:     big.NewInt(1),
			BlockNumber: blockNumber,
			Contract:    common.HexToAddress(contract),
			Timestamp:   timestamp,
		},
	}
}

func createScanEventTransfer(contract string, blockNumber uint64) scan.EventTransfer {
	return scan.EventTransfer{
		From:        common.HexToAddress("0x01"),
		To:          common.HexToAddress("0x02"),
		TokenId:     big.NewInt(1),
		BlockNumber: blockNumber,
		Contract:    common.HexToAddress(contract),
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
