package main

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	mockStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	"github.com/freeverseio/laos-universal-node/internal/repository"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/state/mock"
	"go.uber.org/mock/gomock"
)

func TestRunScanWithStoredContracts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		c                            config.Config
		l1LatestBlock                uint64
		name                         string
		blockNumberDB                string
		blockNumberTimes             int
		scanEventsTimes              int
		scanNewUniversalEventsTimes  int
		txCommitTimes                int
		txDiscardTimes               int
		expectedStartingBlock        uint64
		newLatestBlock               string
		storedContracts              [][]byte
		collectionAddressForContract []string
		expectedContracts            []string
		discoveredContracts          []model.ERC721UniversalContract
		scannedEvents                []scan.Event
		blockNumberTransferEvents    uint64
		timeStampTransferEvents      uint64
		blocknumberMintedEvents      uint64
		timeStampMintedEvents        uint64
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
			name:                        "scan events one time with stored contracts and updateStateWithTransfer",
			blockNumberTimes:            2,
			scanEventsTimes:             1,
			scanNewUniversalEventsTimes: 1,
			txCommitTimes:               1,
			txDiscardTimes:              1,
			newLatestBlock:              "102",
			storedContracts: [][]byte{
				[]byte("contract_0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
			collectionAddressForContract: []string{"0x0000000000000000000000000000000000000000"},
			expectedContracts:            []string{"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},
			discoveredContracts:          getERC721UniversalContracts(),
			scannedEvents:                createERC721TransferEvents(),
			blockNumberTransferEvents:    1,
			timeStampTransferEvents:      1000,
			blocknumberMintedEvents:      0,
			timeStampMintedEvents:        0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner, storage, _ := getMocks(t)
			mockState, tx2 := getMocksFromState(t)

			mockState.EXPECT().NewTransaction().Return(tx2)
			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, nil).
				Times(tt.blockNumberTimes)

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.blockNumberTransferEvents))).Return(&types.Header{
				Time: uint64(tt.timeStampTransferEvents),
			}, nil).Times(1)

			scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock))).
				Return(tt.discoveredContracts, nil).
				Times(tt.scanNewUniversalEventsTimes)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock)), tt.expectedContracts).
				Return(tt.scannedEvents, big.NewInt(int64(tt.l1LatestBlock)), nil).
				Times(tt.scanEventsTimes)

			storage.EXPECT().GetKeysWithPrefix([]byte("contract_")).
				Return(tt.storedContracts, nil).
				Times(1)

			for i, contract := range tt.storedContracts {
				// remove the prefix
				contractAddress := string(contract[9:])
				tx2.EXPECT().GetCollectionAddress(contractAddress).Return(common.HexToAddress(tt.collectionAddressForContract[i]), nil).Times(1)
				tx2.EXPECT().GetMintedWithExternalURIEvents(tt.collectionAddressForContract[i]).
					Return(getMockMintedEvents(tt.blocknumberMintedEvents, tt.timeStampMintedEvents), nil).
					Times(1)
				tx2.EXPECT().GetCurrentEvoBlockForOwnershipContract(contractAddress).Return(uint64(1), nil).Times(1)
				tx2.EXPECT().SetCurrentEvoBlockForOwnershipContract(contractAddress, uint64(1)).Return(nil).Times(1)

			}
			if len(tt.scannedEvents) > 0 {
				// TODO remove any
				tx2.EXPECT().Transfer(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}

			for _, contract := range tt.discoveredContracts {
				tx2.EXPECT().CreateTreesForContract(contract.Address).Return(nil, nil, nil, nil).Times(1)
				tx2.EXPECT().SetTreesForContract(contract.Address, nil, nil, nil).Times(1)
			}
			if len(tt.discoveredContracts) > 0 {
				tx2.EXPECT().StoreERC721UniversalContracts(tt.discoveredContracts).Return(nil).Times(1)
			}

			newLatestBlock, err := strconv.ParseUint(tt.newLatestBlock, 10, 64)
			if err != nil {
				t.Fatalf(`got error "%v" when no error was expeceted`, err)
			}

			tx2.EXPECT().SetCurrentOwnershipBlock(newLatestBlock).Return(nil).Times(1)
			tx2.EXPECT().Commit().Return(nil).Times(tt.txCommitTimes)
			tx2.EXPECT().Discard().Times(tt.txDiscardTimes)

			storage.EXPECT().Get([]byte("ownership_current_block")).
				Return([]byte(tt.blockNumberDB), nil).
				Times(1)

			err = scanUniversalChain(ctx, &tt.c, client, scanner, repository.New(storage), mockState)
			if err != nil {
				t.Fatalf(`got error "%v" when no error was expeceted`, err)
			}
		})
	}
}

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
		// {
		// 	c: config.Config{
		// 		StartingBlock: 1,
		// 		BlocksMargin:  0,
		// 		BlocksRange:   50,
		// 		WaitingTime:   1 * time.Second,
		// 		Contracts:     []string{"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},
		// 	},
		// 	l1LatestBlock:               101,
		// 	name:                        "scan events one time with block number in db",
		// 	blockNumberDB:               "100",
		// 	expectedStartingBlock:       100,
		// 	blockNumberTimes:            2,
		// 	scanEventsTimes:             1,
		// 	scanNewUniversalEventsTimes: 0,
		// 	txCommitTimes:               0,
		// 	txDiscardTimes:              0,
		// 	newLatestBlock:              "102",
		// 	expectedContracts:           []string{"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},
		// },
		// {
		// 	c: config.Config{
		// 		BlocksMargin: 0,
		// 		BlocksRange:  50,
		// 		WaitingTime:  1 * time.Second,
		// 	},
		// 	l1LatestBlock:               100,
		// 	name:                        "scan events with last block from blockchain",
		// 	expectedStartingBlock:       100,
		// 	blockNumberTimes:            3,
		// 	scanEventsTimes:             1,
		// 	scanNewUniversalEventsTimes: 1,
		// 	txCommitTimes:               1,
		// 	txDiscardTimes:              1,
		// 	newLatestBlock:              "101",
		// 	storedContracts: [][]byte{
		// 		[]byte("contract_0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
		// 	},
		// 	expectedContracts: []string{"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},
		// },
		// {
		// 	c: config.Config{
		// 		StartingBlock: 1,
		// 		BlocksMargin:  0,
		// 		BlocksRange:   50,
		// 		WaitingTime:   1 * time.Second,
		// 		Contracts:     []string{"0x0", "0x1"},
		// 	},
		// 	l1LatestBlock:               101,
		// 	name:                        "scan events with last contracts from user",
		// 	blockNumberDB:               "100",
		// 	expectedStartingBlock:       100,
		// 	blockNumberTimes:            2,
		// 	scanEventsTimes:             1,
		// 	scanNewUniversalEventsTimes: 0,
		// 	txCommitTimes:               0,
		// 	txDiscardTimes:              0,
		// 	newLatestBlock:              "102",
		// 	expectedContracts:           []string{"0x0", "0x1"},
		// },
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner, storage, _ := getMocks(t)
			mockState, tx2 := getMocksFromState(t)

			mockState.EXPECT().NewTransaction().Return(tx2)
			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, nil).
				Times(tt.blockNumberTimes)

			scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock))).
				Return(nil, nil).
				Times(tt.scanNewUniversalEventsTimes)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock)), tt.expectedContracts).
				Return(nil, big.NewInt(int64(tt.l1LatestBlock)), nil).
				Times(tt.scanEventsTimes)

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

			newLatestBlock, err := strconv.ParseUint(tt.newLatestBlock, 10, 64)
			if err != nil {
				t.Fatalf(`got error "%v" when no error was expeceted`, err)
			}

			tx2.EXPECT().GetCollectionAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A").Return(common.HexToAddress("0x0"), nil).Times(1)
			tx2.EXPECT().GetMintedWithExternalURIEvents("0x0000000000000000000000000000000000000000").
				Return(getMockMintedEvents(uint64(0), uint64(0)), nil).Times(1)
			tx2.EXPECT().GetCurrentEvoBlockForOwnershipContract("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A").Return(uint64(1), nil).Times(1)
			tx2.EXPECT().SetCurrentEvoBlockForOwnershipContract("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", uint64(1)).Return(nil).Times(1)
			tx2.EXPECT().SetCurrentOwnershipBlock(newLatestBlock).Return(nil).Times(1)
			tx2.EXPECT().Commit().Return(nil).Times(tt.txCommitTimes)
			tx2.EXPECT().Discard().Times(tt.txDiscardTimes)

			storage.EXPECT().Get([]byte("ownership_current_block")).
				Return([]byte(tt.blockNumberDB), nil).
				Times(1)

			err = scanUniversalChain(ctx, &tt.c, client, scanner, repository.New(storage), mockState)
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
	state, tx := getMocksFromState(t)
	state.EXPECT().NewTransaction().Return(tx)

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

	err := scanUniversalChain(ctx, &c, client, scanner, repository.New(storage), state)
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
	state, tx := getMocksFromState(t)
	state.EXPECT().NewTransaction().Return(tx)

	expectedErr := errors.New("block number error")
	client.EXPECT().BlockNumber(ctx).
		Return(uint64(0), expectedErr).
		Times(1)
	storage.EXPECT().Get([]byte("current_block")).
		Return([]byte(""), nil).
		Times(1)

	err := scanUniversalChain(ctx, &c, client, scanner, repository.New(storage), state)
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

			client, scanner, storage, _ := getMocks(t)
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
				}
			}

			err := scanEvoChain(ctx, &tt.c, client, scanner, repository.New(storage))
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Fatalf(`got error "%v", expected error: "%v"`, err, tt.expectedError)
			}
		})
	}
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
	state, tx := getMocksFromState(t)
	state.EXPECT().NewTransaction().Return(tx)

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

	err := scanUniversalChain(ctx, &c, client, scanner, repository.New(storage), state)
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
	return mock.NewMockEthClient(ctrl), mock.NewMockScanner(ctrl),
		mockStorage.NewMockService(ctrl), mockStorage.NewMockTx(ctrl)
}

func getMocksFromState(t *testing.T) (*mockTx.MockService, *mockTx.MockTx) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mockTx.NewMockService(ctrl), mockTx.NewMockTx(ctrl)
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

func getERC721UniversalContracts() []model.ERC721UniversalContract {
	return []model.ERC721UniversalContract{
		{
			Address:           common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			CollectionAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		},
	}
}

func createERC721TransferEvents() []scan.Event {
	var parsedEvents []scan.Event
	parsedEvents = append(parsedEvents, scan.EventTransfer{
		From:        common.HexToAddress("0x0"),
		To:          common.HexToAddress("0x0"),
		TokenId:     big.NewInt(1),
		BlockNumber: 1,
		Contract:    common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
	})
	return parsedEvents
}
