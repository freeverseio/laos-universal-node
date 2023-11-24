package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain/contract"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	mockStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	"github.com/freeverseio/laos-universal-node/internal/repository"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/state/mock"
	v1 "github.com/freeverseio/laos-universal-node/internal/state/v1"
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
		storedContracts             [][]byte
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
			storedContracts: [][]byte{
				[]byte("contract_0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
			expectedContracts: []string{"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},
		},
		{
			c: config.Config{
				StartingBlock: 1,
				BlocksMargin:  0,
				BlocksRange:   50,
				WaitingTime:   1 * time.Second,
				Contracts:     []string{"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},
			},
			l1LatestBlock:               101,
			name:                        "scan events one time with block number in db",
			blockNumberDB:               "100",
			expectedStartingBlock:       100,
			blockNumberTimes:            2,
			scanEventsTimes:             1,
			scanNewUniversalEventsTimes: 0,
			txCommitTimes:               0,
			txDiscardTimes:              0,
			newLatestBlock:              "102",
			expectedContracts:           []string{"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},
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
			storedContracts: [][]byte{
				[]byte("contract_0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
			expectedContracts: []string{"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},
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
			txCommitTimes:               0,
			txDiscardTimes:              0,
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
				Return(nil, big.NewInt(int64(tt.l1LatestBlock)), nil).
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
					Return(tt.storedContracts, nil).
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
		Contracts:     []string{"0x0"},
	}
	ctx, cancel := getContext()
	defer cancel()

	client, scanner, storage, _ := getMocks(t)

	client.EXPECT().BlockNumber(ctx).
		Return(uint64(101), nil).
		AnyTimes()

	scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(51), c.Contracts).
		Return(nil, big.NewInt(51), nil).
		Times(1)
	scanner.EXPECT().ScanEvents(ctx, big.NewInt(52), big.NewInt(101), c.Contracts).
		Return(nil, big.NewInt(101), nil).
		Times(1)

	storage.EXPECT().Get([]byte("contract_0x0")).
		Return([]byte(""), nil).
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
			storage := mockStorage.NewMockService(ctrl)
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

func TestScanEvoChainOnce(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		c                      config.Config
		l1LatestBlock          uint64
		blockNumberDB          string
		blockNumberTimes       int
		scanEventsTimes        int
		expectedFromBlock      uint64
		expectedToBlock        uint64
		expectedNewLatestBlock string
		errorScanEvents        error
		errorSaveBlockNumber   error
		errorGetBlockNumber    error
		errorGetL1LatestBlock  error
		expectedError          error
	}{
		{
			name: "scan evo chain OK",
			c: config.Config{
				StartingBlock:   0,
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:          250,
			blockNumberTimes:       1,
			blockNumberDB:          "100",
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: "151",
		},
		{
			name: "scan evo chain OK with no starting block in DB and 0 in config",
			c: config.Config{
				StartingBlock:   0,
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:          250,
			blockNumberTimes:       2,
			blockNumberDB:          "",
			expectedFromBlock:      250,
			expectedToBlock:        250,
			expectedNewLatestBlock: "251",
		},
		{
			name: "scan evo chain OK with no starting block in DB and 100 in config",
			c: config.Config{
				EvoStartingBlock: 100,
				EvoBlocksMargin:  0,
				EvoBlocksRange:   50,
				WaitingTime:      1 * time.Second,
				Contracts:        []string{},
			},
			l1LatestBlock:          250,
			blockNumberTimes:       1,
			blockNumberDB:          "",
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: "151",
		},
		{
			name: "scan evo chain with an error getting the block number from L1",
			c: config.Config{
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:          250,
			blockNumberTimes:       1,
			blockNumberDB:          "",
			expectedNewLatestBlock: "151",
			errorGetL1LatestBlock:  errors.New("error getting block number from L1"),
			expectedError:          errors.New("error retrieving the latest block from chain: error getting block number from L1"),
		},
		{
			name: "scan evo chain with an error getting the block number from DB",
			c: config.Config{
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:          250,
			blockNumberTimes:       0,
			blockNumberDB:          "",
			expectedNewLatestBlock: "151",
			errorGetBlockNumber:    errors.New("error getting block number from DB"),
			expectedError:          errors.New("error retrieving the current block from storage: error getting block number from DB"),
		},
		{
			name: "scan evo chain and getting an error saving the block number",
			c: config.Config{
				EvoStartingBlock: 100,
				EvoBlocksMargin:  0,
				EvoBlocksRange:   50,
				WaitingTime:      1 * time.Second,
				Contracts:        []string{},
			},
			l1LatestBlock:          250,
			blockNumberTimes:       1,
			blockNumberDB:          "",
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: "151",
			errorSaveBlockNumber:   errors.New("error saving block number"),
			expectedError:          nil, // in this case we break the loop and don't return an error
		},
		{
			name: "scan evo chain and getting an error from scan events",
			c: config.Config{
				EvoStartingBlock: 100,
				EvoBlocksMargin:  0,
				EvoBlocksRange:   50,
				WaitingTime:      1 * time.Second,
				Contracts:        []string{},
			},
			l1LatestBlock:          250,
			blockNumberTimes:       1,
			blockNumberDB:          "100",
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: "151",
			errorScanEvents:        errors.New("error scanning events"),
			expectedError:          nil, // in this case we break the loop and don't return an error
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
				Return(tt.l1LatestBlock, tt.errorGetL1LatestBlock).
				Times(tt.blockNumberTimes)

			storage.EXPECT().Get([]byte("evo_current_block")).
				Return([]byte(tt.blockNumberDB), tt.errorGetBlockNumber).
				Times(1)

			if tt.errorGetL1LatestBlock == nil && tt.errorGetBlockNumber == nil {
				scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedFromBlock)), big.NewInt(int64(tt.expectedToBlock)), nil).
					Return(nil, big.NewInt(int64(tt.expectedToBlock)), tt.errorScanEvents).
					Do(func(_ context.Context, _ *big.Int, _ *big.Int, _ []string) {
						if tt.errorScanEvents != nil {
							cancel() // we cancel the loop since we only want one iteration
						}
					},
					).Times(1)

				if tt.errorScanEvents == nil {
					storage.EXPECT().Set([]byte("evo_current_block"), []byte(tt.expectedNewLatestBlock)).
						Return(tt.errorSaveBlockNumber).Do(
						func(_ []byte, _ []byte) {
							cancel() // we cancel the loop since we only want one iteration
						},
					).Times(1)

					storage.EXPECT().NewTransaction().Return(tx)
					tx.EXPECT().Discard().Return()
				}
			}

			err := scanEvoChain(ctx, &tt.c, client, scanner, repository.New(storage), v1.NewStateService(storage))
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Fatalf(`got error "%v", expected error: "%v"`, err, tt.expectedError)
			}
		})
	}
}

func TestScanEvoChainWithEvents(t *testing.T) {
	t.Parallel()
	mintedWithExternalURIEventHash := "0xa7135052b348b0b4e9943bae82d8ef1c5ac225e594ef4271d12f0744cfc98348"
	eventLog := &types.Log{
		Topics: []common.Hash{
			common.HexToHash(mintedWithExternalURIEventHash),
			common.HexToHash("0x000000000000000000000000c112bde959080c5b46e73749e3e170f47123e85a"),
		},
		Data:        common.Hex2Bytes("00000000000000000000000000000000000000003d5b1313de887a00000000003d5b1313de887a0000000000c112bde959080c5b46e73749e3e170f47123e85a0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000002e516d4e5247426d7272724862754b4558375354544d326f68325077324d757438674863537048706a367a7a637375000000000000000000000000000000000000"),
		BlockNumber: 100,
	}
	evoAbi, err := abi.JSON(strings.NewReader(contract.EvolutionMetaData.ABI))
	if err != nil {
		t.Fatalf(`got error "%v", expected error: "%v"`, err, nil)
	}
	event, err := parseMintedWithExternalURI(eventLog, &evoAbi)
	if err != nil {
		t.Fatalf(`got error "%v", expected error: "%v"`, err, nil)
	}
	tests := []struct {
		name                   string
		c                      config.Config
		l1LatestBlock          uint64
		blockNumberDB          string
		blockNumberTimes       int
		scanEventsTimes        int
		expectedFromBlock      uint64
		expectedToBlock        uint64
		expectedNewLatestBlock string
		errorScanEvents        error
		errorSaveBlockNumber   error
		errorGetBlockNumber    error
		errorGetL1LatestBlock  error
		expectedError          error
	}{
		{
			name: "scan evo chain OK",
			c: config.Config{
				StartingBlock:   0,
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:          250,
			blockNumberTimes:       1,
			blockNumberDB:          "100",
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: "151",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()
			client, scanner, storage, _ := getMocks(t)
			storage2, tx := getMocksFromState(t)

			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, tt.errorGetL1LatestBlock).
				Times(tt.blockNumberTimes)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedFromBlock)), big.NewInt(int64(tt.expectedToBlock)), nil).
				Return([]scan.Event{event}, big.NewInt(int64(tt.expectedToBlock)), tt.errorScanEvents).Times(1)

			storage.EXPECT().Get([]byte("evo_current_block")).
				Return([]byte(tt.blockNumberDB), tt.errorGetBlockNumber).
				Times(1)

			storage.EXPECT().Set([]byte("evo_current_block"), []byte(tt.expectedNewLatestBlock)).
				Return(tt.errorSaveBlockNumber).Do(
				func(_ []byte, _ []byte) {
					cancel() // we cancel the loop since we only want one iteration
				},
			).Times(1)

			tx.EXPECT().GetEvoChainEvents(gomock.Any()).Return(nil, nil).Times(1)
			tx.EXPECT().StoreEvoChainMintEvents(common.HexToAddress("0x0000000000000000000000000000000000000000"), gomock.Any()).Return(nil).Times(1)

			storage2.EXPECT().NewTransaction().Return(tx)
			tx.EXPECT().Discard().Return()

			err := scanEvoChain(ctx, &tt.c, client, scanner, repository.New(storage), storage2)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Fatalf(`got error "%v", expected error: "%v"`, err, tt.expectedError)
			}
		})
	}
}

func TestStoreMintedWithExternalURIEventsByContract(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		scannedEvents map[common.Address][]scan.Event
		storedEvents  map[common.Address][]model.MintedWithExternalURI
		expectedError error
	}{
		{
			name:          "does not stored when no new events",
			scannedEvents: map[common.Address][]scan.Event{},
			expectedError: nil,
		},
		{
			name: "Store new event to non stored contract",
			scannedEvents: map[common.Address][]scan.Event{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(2),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Store  different contract events to existing contracts",
			scannedEvents: map[common.Address][]scan.Event{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(2),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(3),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
				},
				common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(43),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"),
					},
				},
			},
			storedEvents: map[common.Address][]model.MintedWithExternalURI{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					model.MintedWithExternalURI{
						Slot: big.NewInt(1),
						To:   common.HexToAddress("0xc52293dCf9bE02C6EF7e6F840BF0F0E2FC45c646"),
					},
				},
				common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"): {
					model.MintedWithExternalURI{
						Slot: big.NewInt(40),
						To:   common.HexToAddress("0xc52293dCf9bE02C6EF7e6F840BF0F0E2FC45c646"),
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Store  different contract events to existing and non existing contracts",
			scannedEvents: map[common.Address][]scan.Event{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(2),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(3),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
				},
				common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(43),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"),
					},
				},
			},
			storedEvents: map[common.Address][]model.MintedWithExternalURI{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					model.MintedWithExternalURI{
						Slot: big.NewInt(1),
						To:   common.HexToAddress("0xc52293dCf9bE02C6EF7e6F840BF0F0E2FC45c646"),
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage, tx := getMocksFromState(t)

			storage.EXPECT().NewTransaction().Return(tx)
			tx.EXPECT().Discard()

			events := make([]scan.Event, 0)
			eventsToStore := make(map[common.Address][]model.MintedWithExternalURI)
			for contract, ev := range tt.scannedEvents {
				events = append(events, ev...)
				scanned := mintEventsToModel(ev)

				tx.EXPECT().GetEvoChainEvents(contract).Return(tt.storedEvents[contract], nil)

				eventsToStore[contract] = append(eventsToStore[contract], tt.storedEvents[contract]...)
				eventsToStore[contract] = append(eventsToStore[contract], scanned...)
				tx.EXPECT().StoreEvoChainMintEvents(contract, eventsToStore[contract]).Return(nil)
			}
			if err := storeMintedWithExternalURIEventsByContract(storage, events); err != nil {
				t.Fatalf(`got error "%v", expected error: "%v"`, err, tt.expectedError)
			}
		})
	}
}

func mintEventsToModel(events []scan.Event) []model.MintedWithExternalURI {
	modelMintEvents := make([]model.MintedWithExternalURI, 0)
	for _, e := range events {
		scanMintEvent := e.(scan.EventMintedWithExternalURI)
		modelMintEvents = append(modelMintEvents, model.MintedWithExternalURI{
			Slot:        scanMintEvent.Slot,
			To:          scanMintEvent.To,
			TokenId:     scanMintEvent.TokenId,
			TokenURI:    scanMintEvent.TokenURI,
			BlockNumber: scanMintEvent.BlockNumber,
			Timestamp:   scanMintEvent.Timestamp,
		})
	}

	return modelMintEvents
}

func TestRunScanAndCancelContext(t *testing.T) {
	t.Parallel()
	c := config.Config{
		StartingBlock: 1,
		BlocksMargin:  0,
		BlocksRange:   50,
		WaitingTime:   1 * time.Second,
		Contracts:     []string{"0x0"},
	}
	ctx, cancel := getContext()
	defer cancel()

	client, scanner, storage, _ := getMocks(t)

	client.EXPECT().BlockNumber(ctx).
		Return(uint64(101), nil).
		AnyTimes()

	scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(c.StartingBlock)), big.NewInt(51), c.Contracts).
		Do(func(ctx context.Context, _ *big.Int, _ *big.Int, _ []string) {
			cancel()
		},
		).Return(nil, big.NewInt(50), nil).Times(1)

	storage.EXPECT().Get([]byte("contract_0x0")).
		Return([]byte(""), nil).
		Times(1)
	storage.EXPECT().Get([]byte("current_block")).
		Return([]byte(""), nil).
		Times(1)
	storage.EXPECT().Set([]byte("current_block"), []byte("51")).
		Return(nil).
		Times(1)

	err := scanUniversalChain(ctx, &c, client, scanner, repository.New(storage))
	if err != nil {
		t.Fatalf(`got error "%v" when no error was expeceted`, err)
	}
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), 100*time.Millisecond)
}

func getMocks(t *testing.T) (*mock.MockEthClient, *mock.MockScanner, *mockStorage.MockService, *mockStorage.MockTx) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockTx.NewMockTx(ctrl)
	return mock.NewMockEthClient(ctrl), mock.NewMockScanner(ctrl),
		mockStorage.NewMockService(ctrl), mockStorage.NewMockTx(ctrl)
}

func getMocksFromState(t *testing.T) (*mockTx.MockService, *mockTx.MockTx) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockTx.NewMockTx(ctrl)
	return mockTx.NewMockService(ctrl), mockTx.NewMockTx(ctrl)
}

func parseMintedWithExternalURI(eL *types.Log, contractAbi *abi.ABI) (scan.EventMintedWithExternalURI, error) {
	var mintWithExternalURI scan.EventMintedWithExternalURI
	eventMintedWithExternalURI := "MintedWithExternalURI"
	err := unpackIntoInterface(&mintWithExternalURI, eventMintedWithExternalURI, contractAbi, eL)
	if err != nil {
		return mintWithExternalURI, err
	}
	mintWithExternalURI.To = common.HexToAddress(eL.Topics[1].Hex())

	return mintWithExternalURI, nil
}

func unpackIntoInterface(e scan.Event, eventName string, contractAbi *abi.ABI, eL *types.Log) error {
	err := contractAbi.UnpackIntoInterface(e, eventName, eL.Data)
	if err != nil {
		return fmt.Errorf("error unpacking the event %s: %w", eventName, err)
	}

	return nil
}
