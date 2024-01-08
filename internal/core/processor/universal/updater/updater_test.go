package updater_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"

	uUpdater "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/updater"
	mockScan "github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

func TestUpdateState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                       string
		mintEvents                 []model.MintedWithExternalURI
		transferEvents             []model.ERC721Transfer
		contract                   string
		collection                 common.Address
		lastTagBlockBeforeMint     uint64
		lastTagBlockBeforeTransfer uint64
		lastBlock                  uint64
		expectedError              error
	}{
		{
			name:                   "update mint and transfer event",
			mintEvents:             getMockMintedEvents(350, 350),
			transferEvents:         getERC721TransferEvents("0x000005555", 351, 351),
			lastTagBlockBeforeMint: 348,
			lastBlock:              351,
			contract:               "0x000005555",
			collection:             common.HexToAddress("0x4444"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.TODO()
			tx, client, scanner := createMocks(t)

			updater := uUpdater.New(client, scanner)

			tx.EXPECT().LoadMerkleTrees(common.HexToAddress(tt.contract)).Return(nil).Times(1)
			tx.EXPECT().GetCollectionAddress(tt.contract).Return(tt.collection, nil)
			tx.EXPECT().GetMintedWithExternalURIEvents(tt.collection.String()).Return(tt.mintEvents, nil)
			tx.EXPECT().GetCurrentEvoEventsIndexForOwnershipContract(tt.contract).Return(uint64(0), nil)
			tx.EXPECT().GetLastTaggedBlock(common.HexToAddress(tt.contract)).Return(int64(tt.lastTagBlockBeforeMint), nil)

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.lastTagBlockBeforeMint+1))).Return(&types.Header{Time: tt.lastTagBlockBeforeMint + 1}, nil)
			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.lastTagBlockBeforeMint+2))).Return(&types.Header{Time: tt.lastTagBlockBeforeMint + 2}, nil)
			tx.EXPECT().GetLastTaggedBlock(common.HexToAddress(tt.contract)).Return(int64(tt.lastTagBlockBeforeMint), nil)
			tx.EXPECT().TagRoot(common.HexToAddress(tt.contract), int64(tt.lastTagBlockBeforeMint+1)).Return(nil)
			tx.EXPECT().DeleteRootTag(common.HexToAddress(tt.contract), int64(tt.lastTagBlockBeforeMint+1-256)).Return(nil)

			tx.EXPECT().Mint(common.HexToAddress(tt.contract), &tt.mintEvents[0]).Return(nil)

			tx.EXPECT().GetLastTaggedBlock(common.HexToAddress(tt.contract)).Return(int64(tt.lastTagBlockBeforeMint+1), nil)
			tx.EXPECT().TagRoot(common.HexToAddress(tt.contract), int64(tt.lastTagBlockBeforeMint+2)).Return(nil)
			tx.EXPECT().DeleteRootTag(common.HexToAddress(tt.contract), int64(tt.lastTagBlockBeforeMint+2-256)).Return(nil)
			tx.EXPECT().Transfer(common.HexToAddress(tt.contract), &tt.transferEvents[0]).Return(nil)

			tx.EXPECT().SetCurrentEvoEventsIndexForOwnershipContract(tt.contract, uint64(1))
			tx.EXPECT().GetLastTaggedBlock(common.HexToAddress(tt.contract)).Return(int64(tt.lastTagBlockBeforeMint+2), nil)
			tx.EXPECT().TagRoot(common.HexToAddress(tt.contract), int64(tt.lastTagBlockBeforeMint+3)).Return(nil)
			tx.EXPECT().DeleteRootTag(common.HexToAddress(tt.contract), int64(tt.lastTagBlockBeforeMint+3-256)).Return(nil)

			err := updater.UpdateState(ctx,
				tx,
				[]string{tt.contract},
				map[string][]model.ERC721Transfer{tt.contract: tt.transferEvents}, model.Block{Number: tt.lastBlock})
			assertError(t, tt.expectedError, err)
		})
	}
}

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

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.eventBlock))).Return(&types.Header{Time: tt.eventTime}, nil)

			events, err := updater.GetModelTransferEvents(ctx, tt.startingBlock, tt.lastBlock, []string{tt.contract})
			assertError(t, tt.expectedError, err)

			if len(events[common.HexToAddress(tt.contract).String()]) != 1 {
				t.Fatalf(`wrong number of events got %v expected %v"`, len(events[tt.contract]), 1)
			}
		})
	}
}

func createMocks(t *testing.T) (*mockTx.MockTx, *mockScan.MockEthClient, *mockScan.MockScanner) {
	ctrl := gomock.NewController(t)
	return mockTx.NewMockTx(ctrl), mockScan.NewMockEthClient(ctrl), mockScan.NewMockScanner(ctrl)
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
