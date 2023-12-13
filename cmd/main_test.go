package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/mock/gomock"

	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain/contract"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	mockStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	"github.com/freeverseio/laos-universal-node/internal/repository"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/state/mock"
)

func TestRunScanWithStoredContracts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		c                                config.Config
		l1LatestBlock                    uint64
		name                             string
		blockNumberDB                    uint64
		blockNumberTimes                 int
		scanEventsTimes                  int
		scanNewUniversalEventsTimes      int
		txCommitTimes                    int
		txDiscardTimes                   int
		expectedStartingBlock            uint64
		newLatestBlock                   uint64
		collectionAddressForContract     []string
		expectedContracts                []string
		discoveredContracts              []model.ERC721UniversalContract
		scannedEvents                    []scan.Event
		blockNumberTransferEvents        uint64
		timeStampTransferEvents          uint64
		blocknumberMintedEvents          uint64
		ownershipContractInitialEvoBlock uint64
		timeStampMintedEvents            uint64
		timeStampLastOwnershipBlock      uint64
		timeStampLastEvoBlock            uint64
		expectedTxMintCalls              int
		endRangeBlockHash                string
		parentBlockHash                  string
	}{
		{
			c: config.Config{
				StartingBlock: 1,
				BlocksMargin:  0,
				BlocksRange:   100,
				WaitingTime:   1 * time.Second,
			},
			l1LatestBlock:                    101,
			expectedStartingBlock:            1,
			name:                             "scan events one time with stored contracts and updateStateWithTransfer",
			blockNumberTimes:                 2,
			scanEventsTimes:                  1,
			scanNewUniversalEventsTimes:      1,
			txCommitTimes:                    1,
			txDiscardTimes:                   1,
			newLatestBlock:                   102,
			collectionAddressForContract:     []string{"0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000"},
			expectedContracts:                []string{"0xc3dd09d5387fa0ab798e0adc152d15b8d1a299df", "0x26cb70039fe1bd36b4659858d4c4d0cbcafd743a"},
			discoveredContracts:              getERC721UniversalContracts(),
			scannedEvents:                    createERC721TransferEvents(),
			blockNumberTransferEvents:        1,
			timeStampTransferEvents:          1000,
			blocknumberMintedEvents:          1,
			ownershipContractInitialEvoBlock: 1,
			timeStampLastOwnershipBlock:      3000,
			timeStampLastEvoBlock:            3000,
			timeStampMintedEvents:            0,
			endRangeBlockHash:                "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:                  "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		},
		{
			c: config.Config{
				StartingBlock: 1,
				BlocksMargin:  0,
				BlocksRange:   100,
				WaitingTime:   1 * time.Second,
			},
			l1LatestBlock:                    101,
			expectedStartingBlock:            1,
			name:                             "scan events one time with stored contracts and updateStateWithTransfer",
			blockNumberTimes:                 2,
			scanEventsTimes:                  1,
			scanNewUniversalEventsTimes:      1,
			txCommitTimes:                    1,
			txDiscardTimes:                   1,
			newLatestBlock:                   102,
			collectionAddressForContract:     []string{"0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000"},
			expectedContracts:                []string{"0xc3dd09d5387fa0ab798e0adc152d15b8d1a299df", "0x26cb70039fe1bd36b4659858d4c4d0cbcafd743a"},
			discoveredContracts:              getERC721UniversalContracts(),
			scannedEvents:                    createERC721TransferEvents(),
			blockNumberTransferEvents:        1,
			timeStampTransferEvents:          1000,
			blocknumberMintedEvents:          101,
			ownershipContractInitialEvoBlock: 0,
			timeStampMintedEvents:            2000,
			timeStampLastOwnershipBlock:      3000,
			timeStampLastEvoBlock:            3000,
			expectedTxMintCalls:              1,
			endRangeBlockHash:                "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:                  "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner, _ := getMocks(t)
			mockState, tx2 := getMocksFromState(t)
			mockState.EXPECT().NewTransaction().Return(tx2).Times(3)
			tx2.EXPECT().Discard().Return().Times(2)

			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, nil).
				Times(tt.blockNumberTimes)

			tx2.EXPECT().GetCurrentEvoBlockTimestamp().Return(tt.timeStampLastEvoBlock, nil).Times(1)
			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.l1LatestBlock))).Return(&types.Header{
				Time: tt.timeStampLastOwnershipBlock,
			}, nil).Times(1)

			client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.l1LatestBlock))).
				Return(&types.Block{}, nil).
				Times(1)

			tx2.EXPECT().SetEndRangeOwnershipBlockHash(common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347")).
				Return(nil).
				Times(1)

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.blockNumberTransferEvents))).Return(&types.Header{
				Time: tt.timeStampTransferEvents,
			}, nil).Times(1)

			scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock))).
				Return(tt.discoveredContracts, nil).
				Times(tt.scanNewUniversalEventsTimes)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock)), tt.expectedContracts).
				Return(tt.scannedEvents, big.NewInt(int64(tt.l1LatestBlock)), nil).
				Times(tt.scanEventsTimes)

			tx2.EXPECT().GetAllERC721UniversalContracts().
				Return(tt.expectedContracts).
				Times(1)

			for _, contract := range tt.discoveredContracts {
				tx2.EXPECT().GetMintedWithExternalURIEvents(contract.CollectionAddress.Hex()).
					Return(getMockMintedEvents(tt.blocknumberMintedEvents, tt.timeStampMintedEvents), nil).
					Times(1)
				client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(contract.BlockNumber))).Return(&types.Header{
					Time: tt.timeStampTransferEvents,
				}, nil).Times(1)
				tx2.EXPECT().Mint(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				tx2.EXPECT().SetCurrentEvoBlockForOwnershipContract(contract.Address.String(), tt.ownershipContractInitialEvoBlock).Return(nil).Times(1)
			}

			for i, contract := range tt.expectedContracts {
				tx2.EXPECT().GetCollectionAddress(contract).Return(common.HexToAddress(tt.collectionAddressForContract[i]), nil).Times(1)
				tx2.EXPECT().GetMintedWithExternalURIEvents(tt.collectionAddressForContract[i]).
					Return(getMockMintedEvents(tt.blocknumberMintedEvents, tt.timeStampMintedEvents), nil).
					Times(1)
				tx2.EXPECT().GetCurrentEvoBlockForOwnershipContract(contract).Return(uint64(1), nil).Times(1)
				tx2.EXPECT().SetCurrentEvoBlockForOwnershipContract(contract, tt.blocknumberMintedEvents).Return(nil).Times(1)
			}
			if len(tt.scannedEvents) > 0 {
				// TODO remove any
				tx2.EXPECT().Transfer(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}

			tx2.EXPECT().Mint(gomock.Any(), gomock.Any()).Return(nil).Times(tt.expectedTxMintCalls)

			for _, contract := range tt.expectedContracts {
				tx2.EXPECT().IsTreeSetForContract(common.HexToAddress(contract)).Return(false).Times(1)
				tx2.EXPECT().CreateTreesForContract(common.HexToAddress(contract)).Return(nil, nil, nil, nil).Times(1)
				tx2.EXPECT().SetTreesForContract(common.HexToAddress(contract), nil, nil, nil).Times(1)
			}

			for i := range tt.discoveredContracts {
				tx2.EXPECT().IsTreeSetForContract(tt.discoveredContracts[i].Address).Return(true).Times(1)
			}

			if len(tt.discoveredContracts) > 0 {
				tx2.EXPECT().StoreERC721UniversalContracts(tt.discoveredContracts).Return(nil).Times(1)
			}

			tx2.EXPECT().GetLastTaggedBlock(gomock.Any()).Return(int64(0), nil).AnyTimes()
			client.EXPECT().HeaderByNumber(gomock.Any(), big.NewInt(int64(1))).Return(&types.Header{
				Time: 2000,
			}, nil).AnyTimes()
			client.EXPECT().HeaderByNumber(gomock.Any(), big.NewInt(int64(2))).Return(&types.Header{
				Time: 3000,
			}, nil).AnyTimes()

			tx2.EXPECT().TagRoot(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			tx2.EXPECT().DeleteRootTag(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			tx2.EXPECT().SetCurrentOwnershipBlock(tt.newLatestBlock).Return(nil).Times(1)
			tx2.EXPECT().Commit().Return(nil).Times(tt.txCommitTimes)
			tx2.EXPECT().Discard().Times(tt.txDiscardTimes)
			tx2.EXPECT().GetCurrentOwnershipBlock().
				Return(tt.blockNumberDB, nil).
				Times(1)

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.newLatestBlock))).
				Return(&types.Header{ParentHash: common.HexToHash(tt.parentBlockHash)}, nil).
				Times(1)

			tx2.EXPECT().GetEndRangeOwnershipBlockHash().
				Return(common.HexToHash(tt.endRangeBlockHash), nil).
				Times(1)

			err := scanUniversalChain(ctx, &tt.c, client, scanner, mockState)
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
		blockNumberDB               uint64
		blockNumberTimes            int
		scanEventsTimes             int
		scanNewUniversalEventsTimes int
		txCommitTimes               int
		txDiscardTimes              int
		expectedStartingBlock       uint64
		newLatestBlock              string
		expectedContracts           []string
		endRangeBlockHash           string
		parentBlockHash             string
		timeStampLastOwnershipBlock uint64
		timeStampLastEvoBlock       uint64
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
			expectedContracts:           []string{"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},
			endRangeBlockHash:           "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:             "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			timeStampLastOwnershipBlock: 3000,
			timeStampLastEvoBlock:       3000,
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
		//  timeStampLastOwnershipBlock:  3000,
		//	timeStampLastEvoBlock:        3000,
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
		//  timeStampLastOwnershipBlock:  3000,
		//	timeStampLastEvoBlock:        3000,
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
		//  timeStampLastOwnershipBlock:  3000,
		//	timeStampLastEvoBlock:        3000,
		// },
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner, _ := getMocks(t)
			mockState, tx2 := getMocksFromState(t)

			mockState.EXPECT().NewTransaction().Return(tx2).Times(3)
			tx2.EXPECT().Discard().Return().Times(2)
			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, nil).
				Times(tt.blockNumberTimes)

			newLatestBlock, err := strconv.ParseUint(tt.newLatestBlock, 10, 64)
			if err != nil {
				t.Fatalf(`got error "%v" when no error was expeceted`, err)
			}

			client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.l1LatestBlock))).
				Return(&types.Block{}, nil).
				Times(1)

			tx2.EXPECT().SetEndRangeOwnershipBlockHash(common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347")).
				Return(nil).
				Times(1)

			tx2.EXPECT().GetCurrentEvoBlockTimestamp().Return(tt.timeStampLastEvoBlock, nil).Times(1)
			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.l1LatestBlock))).Return(&types.Header{
				Time: tt.timeStampLastOwnershipBlock,
			}, nil).Times(1)
			scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock))).
				Return(nil, nil).
				Times(tt.scanNewUniversalEventsTimes)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock)), tt.expectedContracts).
				Return(nil, big.NewInt(int64(tt.l1LatestBlock)), nil).
				Times(tt.scanEventsTimes)

			if tt.c.Contracts == nil || len(tt.c.Contracts) == 0 {
				tx2.EXPECT().GetAllERC721UniversalContracts().
					Return(tt.expectedContracts).
					Times(1)
			} else {
				tx2.EXPECT().GetExistingERC721UniversalContracts(tt.c.Contracts).
					Return(tt.c.Contracts).
					Times(1)
			}

			for _, contract := range tt.expectedContracts {
				tx2.EXPECT().IsTreeSetForContract(common.HexToAddress(contract)).Return(false).Times(1)
				tx2.EXPECT().CreateTreesForContract(common.HexToAddress(contract)).Return(nil, nil, nil, nil).Times(1)
				tx2.EXPECT().SetTreesForContract(common.HexToAddress(contract), nil, nil, nil).Times(1)
			}
			tx2.EXPECT().GetCollectionAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A").Return(common.HexToAddress("0x0"), nil).Times(1)
			tx2.EXPECT().GetMintedWithExternalURIEvents("0x0000000000000000000000000000000000000000").
				Return(getMockMintedEvents(uint64(0), uint64(0)), nil).Times(1)
			tx2.EXPECT().GetCurrentEvoBlockForOwnershipContract("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A").Return(uint64(1), nil).Times(1)
			tx2.EXPECT().SetCurrentEvoBlockForOwnershipContract("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", uint64(1)).Return(nil).Times(1)
			tx2.EXPECT().SetCurrentOwnershipBlock(newLatestBlock).Return(nil).Times(1)
			tx2.EXPECT().Commit().Return(nil).Times(tt.txCommitTimes)
			tx2.EXPECT().Discard().Times(tt.txDiscardTimes)
			tx2.EXPECT().GetCurrentOwnershipBlock().
				Return(tt.blockNumberDB, nil).
				Times(1)
			tx2.EXPECT().GetLastTaggedBlock(gomock.Any()).Return(int64(0), nil).AnyTimes()
			tx2.EXPECT().TagRoot(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			tx2.EXPECT().DeleteRootTag(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			client.EXPECT().HeaderByNumber(gomock.Any(), big.NewInt(int64(1))).Return(&types.Header{
				Time: 3000,
			}, nil).AnyTimes()

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(newLatestBlock))).
				Return(&types.Header{ParentHash: common.HexToHash(tt.parentBlockHash)}, nil).
				Times(1)

			tx2.EXPECT().GetEndRangeOwnershipBlockHash().
				Return(common.HexToHash(tt.endRangeBlockHash), nil).
				Times(1)

			err = scanUniversalChain(ctx, &tt.c, client, scanner, mockState)
			if err != nil {
				t.Fatalf(`got error "%v" when no error was expeceted`, err)
			}
		})
	}
}

func TestRunScanError(t *testing.T) {
	t.Parallel()
	c := config.Config{
		StartingBlock: 0,
	}
	ctx, cancel := getContext()
	defer cancel()

	client, scanner, _ := getMocks(t)
	state, tx2 := getMocksFromState(t)

	expectedErr := errors.New("block number error")
	state.EXPECT().NewTransaction().Return(tx2)
	tx2.EXPECT().Discard().Times(1)
	client.EXPECT().BlockNumber(ctx).
		Return(uint64(0), expectedErr).
		Times(1)
	tx2.EXPECT().GetCurrentOwnershipBlock().
		Return(uint64(0), nil).
		Times(1)

	err := scanUniversalChain(ctx, &c, client, scanner, state)
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
			_, tx2 := getMocksFromState(t)
			for _, contract := range tt.contracts {
				tx2.EXPECT().HasERC721UniversalContract(contract).Return(tt.mockReturn, tt.mockReturnErr).Times(1)
				if !tt.mockReturn {
					break // break because we only need to check one contract
				}
			}

			result, err := shouldDiscover(tx2, tt.contracts)

			if tt.expectingError {
				if err == nil {
					t.Errorf("got no error when an error was expected")
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
			mockClient, _, storage := getMocks(t)
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
		blockNumberDB          uint64
		blockNumberTimes       int
		txCreatedTimes         int
		endRangeBlockHash      string
		parentBlockHash        string
		expectedFromBlock      uint64
		expectedToBlock        uint64
		expectedNewLatestBlock uint64
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
			txCreatedTimes:         2,
			blockNumberTimes:       1,
			blockNumberDB:          100,
			endRangeBlockHash:      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:        "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: 151,
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
			txCreatedTimes:         2,
			blockNumberTimes:       2,
			blockNumberDB:          0,
			endRangeBlockHash:      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:        "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			expectedFromBlock:      250,
			expectedToBlock:        250,
			expectedNewLatestBlock: 251,
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
			txCreatedTimes:         2,
			blockNumberTimes:       1,
			blockNumberDB:          0,
			endRangeBlockHash:      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:        "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: 151,
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
			txCreatedTimes:         1,
			blockNumberTimes:       1,
			blockNumberDB:          0,
			endRangeBlockHash:      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:        "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			expectedNewLatestBlock: 151,
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
			txCreatedTimes:         1,
			blockNumberTimes:       0,
			blockNumberDB:          0,
			endRangeBlockHash:      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:        "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			expectedNewLatestBlock: 151,
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
			txCreatedTimes:         2,
			blockNumberTimes:       1,
			blockNumberDB:          0,
			endRangeBlockHash:      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:        "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: 151,
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
			l1LatestBlock:     250,
			txCreatedTimes:    1,
			blockNumberTimes:  1,
			blockNumberDB:     100,
			endRangeBlockHash: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:   "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			expectedFromBlock: 100,
			expectedToBlock:   150,
			errorScanEvents:   errors.New("error scanning events"),
			expectedError:     nil, // in this case we break the loop and don't return an error
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner, _ := getMocks(t)
			state, tx2 := getMocksFromState(t)
			state.EXPECT().NewTransaction().Return(tx2).Times(tt.txCreatedTimes)
			tx2.EXPECT().Discard().Times(tt.txCreatedTimes)
			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, tt.errorGetL1LatestBlock).
				Times(tt.blockNumberTimes)

			tx2.EXPECT().GetCurrentEvoBlock().
				Return(tt.blockNumberDB, tt.errorGetBlockNumber).
				Times(1)

			if tt.errorGetL1LatestBlock == nil && tt.errorGetBlockNumber == nil {
				client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.expectedToBlock))).
					Return(&types.Block{}, nil).
					Times(1)

				tx2.EXPECT().SetEndRangeEvoBlockHash(common.HexToHash(tt.endRangeBlockHash)).
					Return(nil).
					Times(1)

				scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedFromBlock)), big.NewInt(int64(tt.expectedToBlock)), nil).
					Return(nil, big.NewInt(int64(tt.expectedToBlock)), tt.errorScanEvents).
					Do(func(_ context.Context, _ *big.Int, _ *big.Int, _ []string) {
						if tt.errorScanEvents != nil {
							cancel() // we cancel the loop since we only want one iteration
						}
					},
					).Times(1)

				if tt.errorScanEvents == nil {
					tx2.EXPECT().SetCurrentEvoBlock(tt.expectedNewLatestBlock).
						Return(tt.errorSaveBlockNumber).Do(
						func(_ uint64) {
							cancel()
						},
					).Times(1)
					if tt.errorSaveBlockNumber == nil {
						client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.expectedToBlock))).Return(&types.Header{Time: 1000}, nil).Times(1)
						tx2.EXPECT().SetCurrentEvoBlockTimestamp(uint64(1000)).Return(nil).Times(1)
						tx2.EXPECT().Commit().Return(nil).Times(1)

						client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.expectedNewLatestBlock))).
							Return(&types.Header{ParentHash: common.HexToHash(tt.parentBlockHash)}, nil).
							Times(1)

						tx2.EXPECT().GetEndRangeEvoBlockHash().
							Return(common.HexToHash(tt.endRangeBlockHash), nil).
							Times(1)
					}
				}
			}

			err := scanEvoChain(ctx, &tt.c, client, scanner, state)
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
		blockNumberDB          uint64
		blockNumberTimes       int
		scanEventsTimes        int
		endRangeBlockHash      string
		parentBlockHash        string
		expectedFromBlock      uint64
		expectedToBlock        uint64
		expectedNewLatestBlock uint64
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
			blockNumberDB:          100,
			endRangeBlockHash:      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			parentBlockHash:        "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: 151,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()
			client, scanner, _ := getMocks(t)
			storage2, tx := getMocksFromState(t)

			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, tt.errorGetL1LatestBlock).
				Times(tt.blockNumberTimes)

			client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.expectedToBlock))).
				Return(&types.Block{}, nil).
				Times(1)

			tx.EXPECT().SetEndRangeEvoBlockHash(common.HexToHash(tt.endRangeBlockHash)).
				Return(nil).
				Times(1)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedFromBlock)), big.NewInt(int64(tt.expectedToBlock)), nil).
				Return([]scan.Event{event}, big.NewInt(int64(tt.expectedToBlock)), tt.errorScanEvents).Times(1)

			tx.EXPECT().GetCurrentEvoBlock().
				Return(tt.blockNumberDB, tt.errorGetBlockNumber).
				Times(1)

			tx.EXPECT().GetMintedWithExternalURIEvents(gomock.Any()).Return(nil, nil).Times(1)
			tx.EXPECT().StoreMintedWithExternalURIEvents("0x0000000000000000000000000000000000000000", gomock.Any()).Return(nil).Times(1)

			storage2.EXPECT().NewTransaction().Return(tx).Times(2)
			tx.EXPECT().Discard().Return().Times(2)
			tx.EXPECT().SetCurrentEvoBlock(tt.expectedNewLatestBlock).Return(tt.errorSaveBlockNumber)
			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.expectedToBlock))).Return(&types.Header{Time: 1000}, nil).Times(1)
			tx.EXPECT().SetCurrentEvoBlockTimestamp(uint64(1000)).Return(nil).Times(1)
			tx.EXPECT().Commit().Return(nil).Do(
				func() {
					cancel() // we cancel the loop since we only want one iteration
				},
			).Times(1)

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.expectedNewLatestBlock))).
				Return(&types.Header{ParentHash: common.HexToHash(tt.parentBlockHash)}, nil).
				Times(1)

			tx.EXPECT().GetEndRangeEvoBlockHash().
				Return(common.HexToHash(tt.endRangeBlockHash), nil).
				Times(1)

			err := scanEvoChain(ctx, &tt.c, client, scanner, storage2)
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
			_, tx := getMocksFromState(t)

			events := make([]scan.Event, 0)
			eventsToStore := make(map[common.Address][]model.MintedWithExternalURI)
			for contract, ev := range tt.scannedEvents {
				events = append(events, ev...)
				scanned := mintEventsToModel(ev)

				tx.EXPECT().GetMintedWithExternalURIEvents(contract.String()).Return(tt.storedEvents[contract], nil)

				eventsToStore[contract] = append(eventsToStore[contract], tt.storedEvents[contract]...)
				eventsToStore[contract] = append(eventsToStore[contract], scanned...)
				tx.EXPECT().StoreMintedWithExternalURIEvents(contract.String(), eventsToStore[contract]).Return(nil)
			}
			if err := storeMintedWithExternalURIEventsByContract(tx, events); err != nil {
				t.Fatalf(`got error "%v", expected error: "%v"`, err, tt.expectedError)
			}
		})
	}
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), 100*time.Millisecond)
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

func getMocks(t *testing.T) (*mock.MockEthClient, *mock.MockScanner, *mockStorage.MockService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockEthClient(ctrl), mock.NewMockScanner(ctrl),
		mockStorage.NewMockService(ctrl)
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
			Address:           common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
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

func getMocksFromState(t *testing.T) (*mockTx.MockService, *mockTx.MockTx) {
	t.Helper()
	ctrl := gomock.NewController(t)
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
