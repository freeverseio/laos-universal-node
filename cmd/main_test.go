package main

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
	"go.uber.org/mock/gomock"
)

func TestRunScanOk(t *testing.T) {
	t.Parallel()
	tests := []struct {
		c                config.Config
		l1LatestBlock    uint64
		name             string
		blockNumberTimes int
		scanEventsTimes  int
	}{
		{
			c: config.Config{
				StartingBlock: 1,
				BlocksMargin:  0,
				BlocksRange:   100,
				WaitingTime:   1 * time.Second,
			},
			l1LatestBlock:    101,
			name:             "scan events one time",
			blockNumberTimes: 2,
			scanEventsTimes:  1,
		},
		{
			c: config.Config{
				StartingBlock: 150,
				BlocksMargin:  50,
				BlocksRange:   100,
				WaitingTime:   1 * time.Second,
			},
			l1LatestBlock:    199,
			name:             "scan events zero times",
			blockNumberTimes: 1,
			scanEventsTimes:  0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner, storage := getMocks(t)

			contract := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")

			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, nil).
				Times(tt.blockNumberTimes)
			scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(tt.c.StartingBlock)), big.NewInt(int64(tt.l1LatestBlock))).
				Return(nil, nil).
				Times(tt.scanEventsTimes)
			storage.EXPECT().ReadAll(context.Background()).Return([]scan.ERC721UniversalContract{
				{
					Address: contract,
				},
			}, nil).Times(tt.scanEventsTimes)
			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.c.StartingBlock)), big.NewInt(int64(tt.l1LatestBlock)), contract).
				Return(nil, nil).
				Times(tt.scanEventsTimes)

			err := scanUniversalChain(ctx, &tt.c, client, scanner, storage)
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

	client, scanner, storage := getMocks(t)

	contract := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
	erc721UniversalContracts := []scan.ERC721UniversalContract{
		{
			Address: contract,
		},
	}

	client.EXPECT().BlockNumber(ctx).
		Return(uint64(101), nil).
		Times(3)
	scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(int64(51))).
		Return(nil, nil).
		Times(1)
	scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(52)), big.NewInt(int64(101))).
		Return(nil, nil).
		Times(1)
	storage.EXPECT().ReadAll(context.Background()).Return(erc721UniversalContracts, nil).Times(1)
	storage.EXPECT().ReadAll(context.Background()).Return(erc721UniversalContracts, nil).Times(1)
	scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(51), contract).
		Return(nil, nil).
		Times(1)
	scanner.EXPECT().ScanEvents(ctx, big.NewInt(52), big.NewInt(101), contract).
		Return(nil, nil).
		Times(1)

	err := scanUniversalChain(ctx, &c, client, scanner, storage)
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

	client, scanner, storage := getMocks(t)

	expectedErr := errors.New("block number error")
	client.EXPECT().BlockNumber(ctx).
		Return(uint64(0), expectedErr).
		Times(1)

	err := scanUniversalChain(ctx, &c, client, scanner, storage)
	if err == nil {
		t.Fatalf(`got no error when error "%v" was expeceted`, expectedErr)
	}
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), 100*time.Millisecond)
}

func getMocks(t *testing.T) (*mock.MockEthClient, *mock.MockScanner, *mock.MockStorage) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockEthClient(ctrl), mock.NewMockScanner(ctrl), mock.NewMockStorage((ctrl))
}
