package worker

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
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

func TestRunScanWithStoredContracts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		c                                              config.Config
		l1LatestBlock                                  uint64
		name                                           string
		blockNumberDB                                  uint64
		blockNumberTimes                               int
		scanEventsTimes                                int
		scanNewUniversalEventsTimes                    int
		txCommitTimes                                  int
		txDiscardTimes                                 int
		expectedStartingBlock                          uint64
		newLatestBlock                                 uint64
		collectionAddressForContract                   []string
		expectedContracts                              []string
		erc721UniversalEvents                          []scan.EventNewERC721Universal
		discoveredContracts                            []model.ERC721UniversalContract
		scannedEvents                                  []scan.Event
		blockNumberTransferEvents                      uint64
		timeStampTransferEvents                        uint64
		blocknumberMintedEvents                        uint64
		timeStampMintedEvents                          uint64
		timeStampLastOwnershipBlock                    uint64
		ownershipContractInitialEvoIndexBeforeDiscover uint64
		ownershipContractInitialEvoIndex               uint64
		timeStampLastEvoBlock                          uint64
		expectedTxMintCalls                            int
		endRangeBlockHash                              common.Hash
		endRangeBlockHeader                            types.Header
	}{
		{
			c: config.Config{
				StartingBlock:   1,
				BlocksMargin:    0,
				BlocksRange:     100,
				WaitingTime:     1 * time.Second,
				GlobalConsensus: "3",
				Parachain:       9999,
			},
			l1LatestBlock:                101,
			expectedStartingBlock:        1,
			name:                         "scan events one time with stored contracts and updateStateWithTransfer",
			blockNumberTimes:             2,
			scanEventsTimes:              1,
			scanNewUniversalEventsTimes:  1,
			txCommitTimes:                1,
			txDiscardTimes:               1,
			newLatestBlock:               102,
			collectionAddressForContract: []string{"0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000"},
			expectedContracts:            []string{"0xc3dd09d5387fa0ab798e0adc152d15b8d1a299df", "0x26cb70039fe1bd36b4659858d4c4d0cbcafd743a"},
			erc721UniversalEvents:        getNewERC721UniversalEvents(),
			discoveredContracts:          getERC721UniversalContract(),
			scannedEvents:                createERC721TransferEvents(),
			blockNumberTransferEvents:    1,
			timeStampTransferEvents:      1000,
			blocknumberMintedEvents:      1,
			ownershipContractInitialEvoIndexBeforeDiscover: 1,
			ownershipContractInitialEvoIndex:               1,
			timeStampLastOwnershipBlock:                    3000,
			timeStampLastEvoBlock:                          3000,
			timeStampMintedEvents:                          0,
			endRangeBlockHash:                              common.HexToHash("0xd7edd5f44a9864a419452ff790d0194cab0fbc5b664f2de41f57b1a6ef3a474d"),
			endRangeBlockHeader: types.Header{
				ParentHash:  common.HexToHash("0x8ef4db2b2081c0516426eba21c941bfc989a6e93e39c1c34ae24a7e372d02f57"),
				UncleHash:   common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
				Coinbase:    common.HexToAddress("0x045c57b46dede60001623105d351c7941c90149e"),
				Root:        common.HexToHash("0xbf24678b95e5152e321267902070da3f8a63200b7f41a922b8d325b8864e68c7"),
				TxHash:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				ReceiptHash: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				Number:      big.NewInt(0),
			},
		},
		{
			c: config.Config{
				StartingBlock:   1,
				BlocksMargin:    0,
				BlocksRange:     100,
				WaitingTime:     1 * time.Second,
				GlobalConsensus: "3",
				Parachain:       9999,
			},
			l1LatestBlock:                101,
			expectedStartingBlock:        1,
			name:                         "scan events one time with stored contracts and updateStateWithTransfer",
			blockNumberTimes:             2,
			scanEventsTimes:              1,
			scanNewUniversalEventsTimes:  1,
			txCommitTimes:                1,
			txDiscardTimes:               1,
			newLatestBlock:               102,
			collectionAddressForContract: []string{"0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000"},
			expectedContracts:            []string{"0xc3dd09d5387fa0ab798e0adc152d15b8d1a299df", "0x26cb70039fe1bd36b4659858d4c4d0cbcafd743a"},
			erc721UniversalEvents:        getNewERC721UniversalEvents(),
			discoveredContracts:          getERC721UniversalContract(),
			scannedEvents:                createERC721TransferEvents(),
			blockNumberTransferEvents:    1,
			timeStampTransferEvents:      1000,
			blocknumberMintedEvents:      101,
			ownershipContractInitialEvoIndexBeforeDiscover: 0,
			ownershipContractInitialEvoIndex:               1,
			timeStampMintedEvents:                          2000,
			timeStampLastOwnershipBlock:                    3000,
			timeStampLastEvoBlock:                          3000,
			expectedTxMintCalls:                            1,
			endRangeBlockHash:                              common.HexToHash("0xd7edd5f44a9864a419452ff790d0194cab0fbc5b664f2de41f57b1a6ef3a474d"),
			endRangeBlockHeader: types.Header{
				ParentHash:  common.HexToHash("0x8ef4db2b2081c0516426eba21c941bfc989a6e93e39c1c34ae24a7e372d02f57"),
				UncleHash:   common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
				Coinbase:    common.HexToAddress("0x045c57b46dede60001623105d351c7941c90149e"),
				Root:        common.HexToHash("0xbf24678b95e5152e321267902070da3f8a63200b7f41a922b8d325b8864e68c7"),
				TxHash:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				ReceiptHash: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				Number:      big.NewInt(0),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner := getMocks(t)
			mockState, tx2 := getMocksFromState(t)
			mockState.EXPECT().NewTransaction().Return(tx2).Times(3)
			tx2.EXPECT().Discard().Return().Times(2)

			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, nil).
				Times(tt.blockNumberTimes)

			tx2.EXPECT().GetLastEvoBlock().Return(model.Block{Timestamp: tt.timeStampLastEvoBlock}, nil).Times(1)
			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.l1LatestBlock))).Return(&types.Header{
				Time: tt.timeStampLastOwnershipBlock,
			}, nil).Times(1)

			tx2.EXPECT().GetOwnershipEndRangeBlockHash().Return(tt.endRangeBlockHash, nil).Times(1)

			block := types.NewBlockWithHeader(&tt.endRangeBlockHeader)
			if tt.endRangeBlockHash != (common.Hash{}) {
				client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.expectedStartingBlock-1))).
					Return(block, nil).
					Times(1)
			}

			client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.l1LatestBlock))).
				Return(block, nil).
				Times(1)

			tx2.EXPECT().SetOwnershipEndRangeBlockHash(tt.endRangeBlockHash).
				Return(nil).
				Times(1)

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.blockNumberTransferEvents))).Return(&types.Header{
				Time: tt.timeStampTransferEvents,
			}, nil).Times(1)

			scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock))).
				Return(tt.erc721UniversalEvents, nil).
				Times(tt.scanNewUniversalEventsTimes)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock)), tt.expectedContracts).
				Return(tt.scannedEvents, nil).
				Times(tt.scanEventsTimes)

			tx2.EXPECT().GetAllERC721UniversalContracts().
				Return(tt.expectedContracts).
				Times(1)

			for _, contract := range tt.erc721UniversalEvents {
				collectionAddress, _ := contract.CollectionAddress()
				tx2.EXPECT().GetMintedWithExternalURIEvents(collectionAddress.String()).
					Return(getMockMintedEvents(tt.blocknumberMintedEvents, tt.timeStampMintedEvents), nil).
					Times(1)
				client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(contract.BlockNumber))).Return(&types.Header{
					Time: tt.timeStampTransferEvents,
				}, nil).Times(1)
				tx2.EXPECT().Mint(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				tx2.EXPECT().SetCurrentEvoEventsIndexForOwnershipContract(contract.NewContractAddress.String(), tt.ownershipContractInitialEvoIndexBeforeDiscover).Return(nil).Times(1)
			}

			for i, contract := range tt.expectedContracts {
				tx2.EXPECT().GetCollectionAddress(contract).Return(common.HexToAddress(tt.collectionAddressForContract[i]), nil).Times(1)
				tx2.EXPECT().GetMintedWithExternalURIEvents(tt.collectionAddressForContract[i]).
					Return(getMockMintedEvents(tt.blocknumberMintedEvents, tt.timeStampMintedEvents), nil).
					Times(1)
				tx2.EXPECT().GetCurrentEvoEventsIndexForOwnershipContract(contract).Return(tt.ownershipContractInitialEvoIndexBeforeDiscover, nil).Times(1)
				tx2.EXPECT().SetCurrentEvoEventsIndexForOwnershipContract(contract, tt.ownershipContractInitialEvoIndex).Return(nil).Times(1)
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

			for i := range tt.erc721UniversalEvents {
				tx2.EXPECT().IsTreeSetForContract(tt.erc721UniversalEvents[i].NewContractAddress).Return(true).Times(1)
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

			worker := NewWorker(&tt.c, client, scanner, mockState)
			err := worker.Run(ctx)
			if err != nil {
				t.Fatalf(`got error "%v" when no error was expected`, err)
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
		endRangeBlockHash           common.Hash
		endRangeBlockHeader         types.Header
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
			endRangeBlockHash:           common.HexToHash("0xd7edd5f44a9864a419452ff790d0194cab0fbc5b664f2de41f57b1a6ef3a474d"),
			endRangeBlockHeader: types.Header{
				ParentHash:  common.HexToHash("0x8ef4db2b2081c0516426eba21c941bfc989a6e93e39c1c34ae24a7e372d02f57"),
				UncleHash:   common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
				Coinbase:    common.HexToAddress("0x045c57b46dede60001623105d351c7941c90149e"),
				Root:        common.HexToHash("0xbf24678b95e5152e321267902070da3f8a63200b7f41a922b8d325b8864e68c7"),
				TxHash:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				ReceiptHash: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				Number:      big.NewInt(0),
			},
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

			client, scanner := getMocks(t)
			mockState, tx2 := getMocksFromState(t)

			mockState.EXPECT().NewTransaction().Return(tx2).Times(3)
			tx2.EXPECT().Discard().Return().Times(2)
			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, nil).
				Times(tt.blockNumberTimes)

			newLatestBlock, err := strconv.ParseUint(tt.newLatestBlock, 10, 64)
			if err != nil {
				t.Fatalf(`got error "%v" when no error was expected`, err)
			}

			tx2.EXPECT().GetOwnershipEndRangeBlockHash().Return(tt.endRangeBlockHash, nil).Times(1)

			block := types.NewBlockWithHeader(&tt.endRangeBlockHeader)
			if tt.endRangeBlockHash != (common.Hash{}) {
				client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.expectedStartingBlock-1))).
					Return(block, nil).
					Times(1)
			}

			client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.l1LatestBlock))).
				Return(block, nil).
				Times(1)

			tx2.EXPECT().SetOwnershipEndRangeBlockHash(tt.endRangeBlockHash).
				Return(nil).
				Times(1)

			tx2.EXPECT().GetLastEvoBlock().Return(model.Block{Timestamp: tt.timeStampLastEvoBlock}, nil).Times(1)
			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.l1LatestBlock))).Return(&types.Header{
				Time: tt.timeStampLastOwnershipBlock,
			}, nil).Times(1)
			scanner.EXPECT().ScanNewUniversalEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock))).
				Return(nil, nil).
				Times(tt.scanNewUniversalEventsTimes)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedStartingBlock)), big.NewInt(int64(tt.l1LatestBlock)), tt.expectedContracts).
				Return(nil, nil).
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
			tx2.EXPECT().GetCurrentEvoEventsIndexForOwnershipContract("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A").Return(uint64(1), nil).Times(1)
			tx2.EXPECT().SetCurrentEvoEventsIndexForOwnershipContract("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", uint64(1)).Return(nil).Times(1)
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

			worker := NewWorker(&tt.c, client, scanner, mockState)
			err = worker.Run(ctx)

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

	client, scanner := getMocks(t)
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

	worker := NewWorker(&c, client, scanner, state)
	err := worker.Run(ctx)

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

func getNewERC721UniversalEvents() []scan.EventNewERC721Universal {
	return []scan.EventNewERC721Universal{
		{
			BaseURI:            "https://uloc.io/GlobalConsensus(3)/Parachain(9999)/AccountKey20(0x0000000000000000000000000000000000000000)/",
			BlockNumber:        100,
			NewContractAddress: common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
		},
	}
}

func getERC721UniversalContract() []model.ERC721UniversalContract {
	return []model.ERC721UniversalContract{
		{
			Address:           common.HexToAddress("0xC3dd09D5387FA0Ab798e0ADC152d15b8d1a299DF"),
			CollectionAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			BlockNumber:       100,
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

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), 100*time.Millisecond)
}

func getMocks(t *testing.T) (*mock.MockEthClient, *mock.MockScanner) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockEthClient(ctrl), mock.NewMockScanner(ctrl)
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

func getMocksFromState(t *testing.T) (*mockTx.MockService, *mockTx.MockTx) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mockTx.NewMockService(ctrl), mockTx.NewMockTx(ctrl)
}
