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

func TestRunScanOk(t *testing.T) {
	t.Parallel()
	tests := []struct {
		c                           config.Config
		l1LatestBlock               uint64
		name                        string
		blockNumberDB               string
		blockNumberTimes            int
		scanEventsTimes             int
		scanNewUniversalEventsTimes int
		txCommitTimes               int
		txDiscardTimes              int
		expectedStartingBlock       uint64
		newLatestBlock              string
		expectedContracts           []string
	}{
		{
			c: config.Config{
				StartingBlock: 1,
				BlocksMargin:  0,
				BlocksRange:   100,
				WaitingTime:   1 * time.Second,
			},
			l1LatestBlock:               101,
			expectedStartingBlock:       1,
			name:                        "scan events one time",
			blockNumberTimes:            2,
			scanEventsTimes:             1,
			scanNewUniversalEventsTimes: 1,
			txCommitTimes:               1,
			txDiscardTimes:              1,
			newLatestBlock:              "102",
		},
		{
			c: config.Config{
				StartingBlock: 1,
				BlocksMargin:  0,
				BlocksRange:   50,
				WaitingTime:   1 * time.Second,
			},
			l1LatestBlock:               101,
			name:                        "scan events one time with block number in db",
			blockNumberDB:               "100",
			expectedStartingBlock:       100,
			blockNumberTimes:            2,
			scanEventsTimes:             1,
			scanNewUniversalEventsTimes: 1,
			txCommitTimes:               1,
			txDiscardTimes:              1,
			newLatestBlock:              "102",
		},
		{
			c: config.Config{
				BlocksMargin: 0,
				BlocksRange:  50,
				WaitingTime:  1 * time.Second,
			},
			l1LatestBlock:               100,
			name:                        "scan events with last block from blockchain",
			expectedStartingBlock:       100,
			blockNumberTimes:            3,
			scanEventsTimes:             1,
			scanNewUniversalEventsTimes: 1,
			txCommitTimes:               1,
			txDiscardTimes:              1,
			newLatestBlock:              "101",
		},
		{
			c: config.Config{
				StartingBlock: 1,
				BlocksMargin:  0,
				BlocksRange:   50,
				WaitingTime:   1 * time.Second,
				Contracts:     []string{"0x0", "0x1"},
			},
			l1LatestBlock:               101,
			name:                        "scan events with last contracts from user",
			blockNumberDB:               "100",
			expectedStartingBlock:       100,
			blockNumberTimes:            2,
			scanEventsTimes:             1,
			scanNewUniversalEventsTimes: 0,
			txCommitTimes:               1,
			txDiscardTimes:              1,
			newLatestBlock:              "102",
			expectedContracts:           []string{"0x0", "0x1"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner, storage, tx := getMocks(t)
			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, nil).
				Times(tt.blockNumberTimes)

			scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock))).
				Return(nil, nil).
				Times(tt.scanNewUniversalEventsTimes)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock)), tt.expectedContracts).
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

			if tt.c.Contracts == nil || len(tt.c.Contracts) == 0 {
				storage.EXPECT().GetKeysWithPrefix([]byte("contract_")).
					Return([][]byte{}, nil).
					Times(1)
			} else {
				for _, contract := range tt.c.Contracts {
					storage.EXPECT().Get([]byte("contract_"+contract)).
						Return([]byte("1"), nil).
						Times(1)
				}
			}
			storage.EXPECT().Get([]byte("current_block")).
				Return([]byte(tt.blockNumberDB), nil).
				Times(1)
			storage.EXPECT().Set([]byte("current_block"), []byte(tt.newLatestBlock)).
				Return(nil).
				Times(1)

			err := scanUniversalChain(ctx, &tt.c, client, scanner, repository.New(storage))
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
	var expectedContracts []string
	client.EXPECT().BlockNumber(ctx).
		Return(uint64(101), nil).
		Times(3)
	scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(int64(51))).
		Return(nil, nil).
		Times(1)
	scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(52)), big.NewInt(int64(101))).
		Return(nil, nil).
		Times(1)
	scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(51), expectedContracts).
		Return(nil, nil).
		Times(1)
	scanner.EXPECT().ScanEvents(ctx, big.NewInt(52), big.NewInt(101), expectedContracts).
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

	err := scanUniversalChain(ctx, &c, client, scanner, repository.New(storage))
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

	err := scanUniversalChain(ctx, &c, client, scanner, repository.New(storage))
	if err == nil {
		t.Fatalf(`got no error when error "%v" was expected`, expectedErr)
	}
}

func TestShouldDiscover(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		contracts      []string
		mockReturn     bool
		mockReturnErr  error
		expectedResult bool
		expectingError bool
	}{
		{
			name:           "no contracts should discover",
			contracts:      []string{},
			expectedResult: true,
			expectingError: false,
		},
		{
			name:           "existing contracts should not discover",
			contracts:      []string{"contract1", "contract2"},
			mockReturn:     true,
			expectedResult: false,
			expectingError: false,
		},
		{
			name:           "at least one new contract should discover",
			contracts:      []string{"contract1", "new_contract"},
			mockReturn:     false, // for example new_contract exists but contract1 does not
			expectedResult: true,
			expectingError: false,
		},
		{
			name:           "error retrieving contract should return error",
			contracts:      []string{"contract1"},
			mockReturnErr:  errors.New("storage error"),
			expectedResult: false,
			expectingError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			storage := mockStorage.NewMockStorage(ctrl)
			repositoryService := repository.New(storage)
			for _, contract := range tt.contracts {
				var returnValue []byte
				if tt.mockReturn {
					returnValue = []byte("1")
				} else {
					returnValue = nil
				}
				storage.EXPECT().Get([]byte("contract_"+contract)).Return(returnValue, tt.mockReturnErr).Times(1)
				if !tt.mockReturn {
					break // break because we only need to check one contract
				}
			}

			result, err := shouldDiscover(repositoryService, tt.contracts)

			if tt.expectingError {
				if err == nil {
					t.Errorf("got no error nut expected an error")
				}
			} else {
				if err != nil {
					t.Errorf("got unexpected error: %v, expected none", err)
				}
				if result != tt.expectedResult {
					t.Errorf("got %v, expected %v", result, tt.expectedResult)
				}
			}
		})
	}
}

func TestCompareChainIDs(t *testing.T) {
	t.Parallel()
	ethChainID := big.NewInt(3) // Example Ethereum network chain ID
	dbChainID := "3"            // Matching chain ID for the database

	tests := []struct {
		name               string
		ethClientChainID   *big.Int
		ethClientError     error
		dbChainID          string
		dbError            error
		expectSetInDB      bool
		expectedDBSetError error
		wantError          bool
	}{
		{
			name:             "chain IDs match",
			ethClientChainID: ethChainID,
			dbChainID:        dbChainID,
			wantError:        false,
		},
		{
			name:           "error from Ethereum client",
			ethClientError: errors.New("client error"),
			wantError:      true,
		},
		{
			name:      "error from database on get",
			dbError:   errors.New("db get error"),
			wantError: true,
		},
		{
			name:               "chain ID not set in database",
			ethClientChainID:   ethChainID,
			dbChainID:          "",
			expectSetInDB:      true,
			expectedDBSetError: nil,
			wantError:          false,
		},
		{
			name:             "mismatched chain IDs",
			ethClientChainID: ethChainID,
			dbChainID:        "2",
			wantError:        true,
		},
		{
			name:               "error from database on set",
			ethClientChainID:   ethChainID,
			dbChainID:          "",
			expectSetInDB:      true,
			expectedDBSetError: errors.New("db set error"),
			wantError:          true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockClient, _, storage, _ := getMocks(t)
			repositoryService := repository.New(storage)
			ctx := context.Background()

			mockClient.EXPECT().ChainID(ctx).Return(tt.ethClientChainID, tt.ethClientError).Times(1)
			if tt.ethClientError == nil {
				storage.EXPECT().Get([]byte("chain_id")).Return([]byte(tt.dbChainID), tt.dbError).Times(1)
			}
			if tt.expectSetInDB {
				storage.EXPECT().Set([]byte("chain_id"), []byte(ethChainID.String())).Return(tt.expectedDBSetError)
			}

			err := compareChainIDs(ctx, mockClient, repositoryService)

			if (err != nil) != tt.wantError {
				t.Errorf("compareChainIDs() error = %v, wantError %v", err, tt.wantError)
			}
		})
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
