package main

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/freeverseio/laos-universal-node/internal/config"
	mockStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	"github.com/freeverseio/laos-universal-node/internal/repository"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
	"go.uber.org/mock/gomock"
)

// TODO check test coverage
func TestRunScanOk(t *testing.T) {
	t.Parallel()
	tests := []struct {
		c                     config.Config
		l1LatestBlock         uint64
		name                  string
		blockNumberDB         string
		blockNumberTimes      int
		scanEventsTimes       int
		txCommitTimes         int
		txDiscardTimes        int
		expectedStartingBlock uint64
		newLatestBlock        string
	}{
		{
			c: config.Config{
				StartingBlock: 1,
				BlocksMargin:  0,
				BlocksRange:   100,
				WaitingTime:   1 * time.Second,
			},
			l1LatestBlock:         101,
			expectedStartingBlock: 1,
			name:                  "scan events one time",
			blockNumberTimes:      2,
			scanEventsTimes:       1,
			txCommitTimes:         1,
			txDiscardTimes:        1,
			newLatestBlock:        "102",
		},
		{
			c: config.Config{
				StartingBlock: 1,
				BlocksMargin:  0,
				BlocksRange:   50,
				WaitingTime:   1 * time.Second,
			},
			l1LatestBlock:         101,
			name:                  "scan events one time with block number in db",
			blockNumberDB:         "100",
			expectedStartingBlock: 100,
			blockNumberTimes:      2,
			scanEventsTimes:       1,
			txCommitTimes:         1,
			txDiscardTimes:        1,
			newLatestBlock:        "102",
		},
		{
			c: config.Config{
				BlocksMargin: 0,
				BlocksRange:  50,
				WaitingTime:  1 * time.Second,
			},
			l1LatestBlock:         100,
			name:                  "scan events with last block from blockchain",
			expectedStartingBlock: 100,
			blockNumberTimes:      3,
			scanEventsTimes:       1,
			txCommitTimes:         1,
			txDiscardTimes:        1,
			newLatestBlock:        "101",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner, storage, tx := getMocks(t)
			var expecetedContracts []string
			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, nil).
				Times(tt.blockNumberTimes)
			scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock))).
				Return(nil, nil).
				Times(tt.scanEventsTimes)
			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock)), expecetedContracts).
				Return(nil, nil).
				Times(tt.scanEventsTimes)
			tx.EXPECT().Commit().
				Return(nil).
				Times(tt.txCommitTimes)
			tx.EXPECT().Discard().
				Times(tt.txDiscardTimes)
			storage.EXPECT().NewTransaction().
				Return(tx).
				Times(tt.txCommitTimes)

			storage.EXPECT().GetKeysWithPrefix([]byte("contract_")).
				Return([][]byte{}, nil).
				Times(1)
			storage.EXPECT().Get([]byte("current_block")).
				Return([]byte(tt.blockNumberDB), nil).
				Times(1)
			storage.EXPECT().Set([]byte("current_block"), []byte(tt.newLatestBlock)).
				Return(nil).
				Times(1)

			err := runScan(ctx, &tt.c, client, scanner, repository.New(storage))
			if err != nil {
				t.Fatalf(`got error "%v" when no error was expeceted`, err)
			}
		})
	}
}

func TestRunScanTwice(t *testing.T) {
	t.Parallel()
	c := config.Config{
		StartingBlock: 1,
		BlocksMargin:  0,
		BlocksRange:   50,
		WaitingTime:   1 * time.Second,
	}
	ctx, cancel := getContext()
	defer cancel()

	client, scanner, storage, tx := getMocks(t)
	var expecetedContracts []string
	client.EXPECT().BlockNumber(ctx).
		Return(uint64(101), nil).
		Times(3)
	scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(int64(51))).
		Return(nil, nil).
		Times(1)
	scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(52)), big.NewInt(int64(101))).
		Return(nil, nil).
		Times(1)
	scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(51), expecetedContracts).
		Return(nil, nil).
		Times(1)
	scanner.EXPECT().ScanEvents(ctx, big.NewInt(52), big.NewInt(101), expecetedContracts).
		Return(nil, nil).
		Times(1)
	tx.EXPECT().Commit().
		Return(nil).
		Times(2)
	tx.EXPECT().Discard().
		Times(2)
	storage.EXPECT().NewTransaction().
		Return(tx).
		Times(2)
	storage.EXPECT().GetKeysWithPrefix([]byte("contract_")).
		Return([][]byte{}, nil).
		Times(2)
	storage.EXPECT().Get([]byte("current_block")).
		Return([]byte(""), nil).
		Times(1)
	storage.EXPECT().Set([]byte("current_block"), []byte("52")).
		Return(nil).
		Times(1)
	storage.EXPECT().Set([]byte("current_block"), []byte("102")).
		Return(nil).
		Times(1)

	err := runScan(ctx, &c, client, scanner, repository.New(storage))
	if err != nil {
		t.Fatalf(`got error "%v" when no error was expeceted`, err)
	}
}

func TestRunScanError(t *testing.T) {
	t.Parallel()
	c := config.Config{
		StartingBlock: 0,
	}
	ctx, cancel := getContext()
	defer cancel()

	client, scanner, storage, _ := getMocks(t)

	expectedErr := errors.New("block number error")
	client.EXPECT().BlockNumber(ctx).
		Return(uint64(0), expectedErr).
		Times(1)
	storage.EXPECT().Get([]byte("current_block")).
		Return([]byte(""), nil).
		Times(1)

	err := runScan(ctx, &c, client, scanner, repository.New(storage))
	if err == nil {
		t.Fatalf(`got no error when error "%v" was expected`, expectedErr)
	}
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), 100*time.Millisecond)
}

func getMocks(t *testing.T) (*mock.MockEthClient, *mock.MockScanner, *mockStorage.MockStorage, *mockStorage.MockTx) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockEthClient(ctrl), mock.NewMockScanner(ctrl),
		mockStorage.NewMockStorage(ctrl), mockStorage.NewMockTx(ctrl)
}
