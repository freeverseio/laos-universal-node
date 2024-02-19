package universal_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/mock/gomock"

	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/universal"
	mockDiscoverer "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer/mock"
	mockUpdater "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/updater/mock"
	mockClient "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	mockScan "github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
)

func TestGetInitStartingBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                           string
		startingBlockData              model.Block
		startingBlockError             error
		userProvidedBlock              uint64
		lastBlockNumberFromClient      uint64
		lastBlockNumberFromClientError error
		expectedResult                 uint64
		expectedError                  error
	}{
		{
			name:               "starting block exists in storage",
			startingBlockData:  model.Block{Number: 10},
			startingBlockError: nil,
			userProvidedBlock:  0,
			expectedResult:     11,
			expectedError:      nil,
		},

		{
			name:               "starting block does not exist in storage, it returns error",
			startingBlockData:  model.Block{},
			startingBlockError: errors.New("error from storage"),
			userProvidedBlock:  0,
			expectedResult:     0,
			expectedError:      errors.New("error retrieving the current block from storage: error from storage"),
		},

		{
			name:              "starting block does not exist in storage, user provided starting block",
			startingBlockData: model.Block{},
			userProvidedBlock: 20,
			expectedResult:    20,
			expectedError:     nil,
		},
		{
			name:                           "starting block does not exist in storage, user provided starting block is zero",
			startingBlockData:              model.Block{},
			startingBlockError:             nil,
			userProvidedBlock:              0,
			lastBlockNumberFromClient:      30,
			lastBlockNumberFromClientError: nil,
			expectedResult:                 30,
			expectedError:                  nil,
		},

		{
			name:                           "starting block does not exist in storage, user provided starting block is zero, error from client",
			startingBlockData:              model.Block{},
			startingBlockError:             nil,
			userProvidedBlock:              0,
			lastBlockNumberFromClient:      0,
			lastBlockNumberFromClientError: errors.New("error from client"),
			expectedResult:                 0,
			expectedError:                  errors.New("error retrieving the latest block from chain: error from client"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()
			stateService, tx, client, _, _, _ := createMocks(t)

			stateService.EXPECT().NewTransaction().Return(tx, nil)
			tx.EXPECT().GetLastOwnershipBlock().Return(tt.startingBlockData, tt.startingBlockError)
			tx.EXPECT().Discard()
			if tt.userProvidedBlock == 0 && tt.startingBlockData.Number == 0 && tt.startingBlockError == nil {
				client.EXPECT().BlockNumber(ctx).Return(tt.lastBlockNumberFromClient, tt.lastBlockNumberFromClientError)
			}

			p := universal.NewProcessor(client, stateService, nil, &config.Config{StartingBlock: tt.userProvidedBlock}, nil, nil)
			result, err := p.GetInitStartingBlock(ctx)
			assertError(t, tt.expectedError, err)
			if result != tt.expectedResult {
				t.Fatalf(`got result "%v", expected "%v"`, result, tt.expectedResult)
			}
		})
	}
}

func TestGetLastBlock(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name               string
		startingBlock      uint64
		l1LatestBlock      uint64
		configBlocksRange  uint64
		configBlocksMargin uint64
		expectedResult     uint64
		expectedError      error
	}{
		{
			name:               "Starting block within range",
			startingBlock:      100,
			l1LatestBlock:      200,
			configBlocksRange:  10,
			configBlocksMargin: 5,
			expectedResult:     110,
			expectedError:      nil,
		},
		{
			name:               "Starting block exceeds range",
			startingBlock:      195,
			l1LatestBlock:      200,
			configBlocksRange:  10,
			configBlocksMargin: 5,
			expectedResult:     195,
			expectedError:      nil,
		},
		{
			name:               "Error getting latest block",
			startingBlock:      100,
			l1LatestBlock:      0,
			configBlocksRange:  10,
			configBlocksMargin: 5,
			expectedResult:     0,
			expectedError:      errors.New("error getting latest block"),
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()
			_, _, client, _, _, _ := createMocks(t)

			client.EXPECT().BlockNumber(ctx).Return(tt.l1LatestBlock, tt.expectedError)

			p := universal.NewProcessor(client, nil, nil, &config.Config{BlocksMargin: uint(tt.configBlocksMargin), BlocksRange: uint(tt.configBlocksRange)}, nil, nil)

			result, err := p.GetLastBlock(ctx, tt.startingBlock)
			assertError(t, tt.expectedError, err)
			if result != tt.expectedResult {
				t.Fatalf(`got result "%v", expected "%v"`, result, tt.expectedResult)
			}
		})
	}
}

func TestProcessUniversalBlockRange(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	tests := []struct {
		name                         string
		startingBlock                uint64
		previousBlockHeaderFromChain *types.Header
		previousBlockDataFromDB      model.Block
		blockHeaderFromChain         *types.Header
		blockDataFromDB              model.Block
		discoverReturn               bool
		updateReturn                 map[uint64]map[string][]model.ERC721Transfer
		expectedError                error
		expectedTxCommit             int
		expectedNumberOfReorgCheck   int
	}{
		{
			name:          "successful processing with discovery and update",
			startingBlock: 100,
			previousBlockHeaderFromChain: &types.Header{
				Number: big.NewInt(90),
			},
			previousBlockDataFromDB: model.Block{
				Number: 90,
				Hash:   common.HexToHash("0x62055b9639cbed71604205301891afe40ae0fe4f57ceadbf35d9a476361c48ed"),
			},
			blockHeaderFromChain: &types.Header{
				Number: big.NewInt(100),
			},
			blockDataFromDB: model.Block{
				Number: 100,
				Hash:   common.HexToHash("0xb07e1289b32edefd8f3c702d016fb73c81d5950b2ebc790ad9d2cb8219066b4c"),
			},
			discoverReturn:             false,
			updateReturn:               make(map[uint64]map[string][]model.ERC721Transfer),
			expectedError:              nil,
			expectedTxCommit:           1,
			expectedNumberOfReorgCheck: 1,
		},
		{
			name:          "processing with reorg",
			startingBlock: 100,
			previousBlockHeaderFromChain: &types.Header{
				Number: big.NewInt(90),
			},
			previousBlockDataFromDB: model.Block{
				Number: 90,
				Hash:   common.HexToHash("0x123"),
			},
			blockHeaderFromChain: &types.Header{
				Number: big.NewInt(100),
			},
			blockDataFromDB: model.Block{
				Number: 100,
				Hash:   common.HexToHash("0xb07e1289b32edefd8f3c702d016fb73c81d5950b2ebc790ad9d2cb8219066b4c"),
			},
			discoverReturn: false,
			updateReturn:   make(map[uint64]map[string][]model.ERC721Transfer),
			expectedError: universal.ReorgError{
				Block:       90,
				ChainHash:   common.HexToHash("0x62055b9639cbed71604205301891afe40ae0fe4f57ceadbf35d9a476361c48ed"),
				StorageHash: common.HexToHash("0x123"),
			},
			expectedTxCommit:           0,
			expectedNumberOfReorgCheck: 1,
		},
		{
			name:          "successful processing with no hash in storage",
			startingBlock: 100,
			previousBlockHeaderFromChain: &types.Header{
				Number: big.NewInt(90),
			},
			previousBlockDataFromDB: model.Block{
				Number: 90,
			},
			blockHeaderFromChain: &types.Header{
				Number: big.NewInt(100),
			},
			blockDataFromDB: model.Block{
				Number: 100,
				Hash:   common.HexToHash("0xb07e1289b32edefd8f3c702d016fb73c81d5950b2ebc790ad9d2cb8219066b4c"),
			},
			discoverReturn:             false,
			updateReturn:               make(map[uint64]map[string][]model.ERC721Transfer),
			expectedError:              nil,
			expectedTxCommit:           1,
			expectedNumberOfReorgCheck: 0,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stateService, tx, client, scanner, discoverer, updater := createMocks(t)
			p := universal.NewProcessor(client, stateService, scanner, &config.Config{}, discoverer, updater)
			startingBlockData := model.Block{
				Number:    100,
				Hash:      common.HexToHash("0xb72b31eb84c4bbbbd62aff06a3c8c88991ac7c118c47aa6fba3609ed1baa8fd3"),
				Timestamp: 110,
			}

			stateService.EXPECT().NewTransaction().Return(tx, nil)

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.startingBlock))).Return(tt.blockHeaderFromChain, nil)

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.previousBlockDataFromDB.Number))).Return(tt.previousBlockHeaderFromChain, nil).Times(tt.expectedNumberOfReorgCheck)

			tx.EXPECT().SetLastOwnershipBlock(tt.blockDataFromDB).Return(nil)
			// the last saved OwnershipBlock
			tx.EXPECT().GetLastOwnershipBlock().Return(tt.previousBlockDataFromDB, nil)
			discoverer.EXPECT().ShouldDiscover(tx, tt.startingBlock, tt.blockDataFromDB.Number).Return(tt.discoverReturn, nil)
			discoverer.EXPECT().GetContracts(tx).Return([]string{"contract"}, nil)

			updater.EXPECT().GetModelTransferEvents(ctx, tt.startingBlock, tt.blockDataFromDB.Number, []string{"contract"}).Return(tt.updateReturn, nil)

			tx.EXPECT().GetFirstOwnershipBlock().Return(model.Block{}, nil)
			client.EXPECT().
				HeaderByNumber(ctx, big.NewInt(int64(startingBlockData.Number))).
				Return(&types.Header{
					Time:   startingBlockData.Timestamp,
					Number: big.NewInt(int64(startingBlockData.Number)),
				}, nil)
			tx.EXPECT().SetFirstOwnershipBlock(startingBlockData).Return(nil)

			updater.EXPECT().UpdateState(ctx, tx, []string{"contract"}, make(map[common.Address]uint64), tt.updateReturn, tt.startingBlock, tt.blockDataFromDB).Return(nil)

			tx.EXPECT().Commit().Return(nil).Times(tt.expectedTxCommit)
			tx.EXPECT().Discard()

			err := p.ProcessUniversalBlockRange(ctx, tt.startingBlock, tt.blockDataFromDB.Number)
			assertError(t, tt.expectedError, err)
		})
	}
}

func TestIsEvoSyncedWithOwnership(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		TimeEvo       uint64
		TimeOwnership uint64
		expect        bool
	}{
		{
			name:          "evo is synced with ownership",
			TimeEvo:       200,
			TimeOwnership: 100,
			expect:        true,
		},
		{
			name:          "evo is not synced with ownership",
			TimeEvo:       100,
			TimeOwnership: 200,
			expect:        false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.TODO()

			stateService, tx, client, scanner, discoverer, updater := createMocks(t)

			p := universal.NewProcessor(client, stateService, scanner, &config.Config{}, discoverer, updater)

			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.TimeOwnership))).
				Return(&types.Header{Number: big.NewInt(100), Time: tt.TimeOwnership}, nil)

			stateService.EXPECT().NewTransaction().Return(tx, nil)
			tx.EXPECT().GetLastEvoBlock().Return(model.Block{Number: tt.TimeEvo, Timestamp: tt.TimeEvo}, nil)
			tx.EXPECT().Discard()

			result, err := p.IsEvoSyncedWithOwnership(ctx, tt.TimeOwnership)
			assertError(t, nil, err)
			if result != tt.expect {
				t.Fatalf(`got result "%v", expected "%v"`, result, tt.expect)
			}
		})
	}
}

func TestRecoverFromReorg(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                     string
		startingBlock            uint64
		safeBlockNumber          uint64
		numberOfRecursions       uint64
		checkReorgError          error
		getAllStoredBlockNumbers []uint64
		getAllContracts          []string
		getBlockHeadersDB        []*types.Header
		getBlockHeadersL1        []*types.Header
		checkoutError            error
		expectedError            error
	}{
		{
			name:                     "successful reorg recovery",
			startingBlock:            100,
			safeBlockNumber:          99,
			numberOfRecursions:       1,
			checkReorgError:          nil,
			getAllStoredBlockNumbers: []uint64{100, 99, 98},
			getAllContracts:          []string{"contract1", "contract2"},
			getBlockHeadersDB: []*types.Header{{
				Number: big.NewInt(99),
				Time:   100,
			}},
			getBlockHeadersL1: []*types.Header{{
				Number: big.NewInt(99),
				Time:   100,
			}},
			checkoutError: nil,
			expectedError: nil,
		},
		{
			name:                     "successful reorg recovery",
			startingBlock:            100,
			safeBlockNumber:          95,
			numberOfRecursions:       1,
			checkReorgError:          nil,
			getAllStoredBlockNumbers: []uint64{100, 95, 94},
			getAllContracts:          []string{"contract1", "contract2"},
			getBlockHeadersDB: []*types.Header{{
				Number: big.NewInt(95),
				Time:   100,
			}},
			getBlockHeadersL1: []*types.Header{{
				Number: big.NewInt(95),
				Time:   100,
			}},
			checkoutError: nil,
			expectedError: nil,
		},
		{
			name:                     "successful reorg recovery with no blocks to checkout",
			startingBlock:            100,
			numberOfRecursions:       0,
			safeBlockNumber:          0,
			checkReorgError:          nil,
			getAllStoredBlockNumbers: []uint64{100},
			getAllContracts:          []string{"contract1", "contract2"},
			checkoutError:            nil,
			expectedError:            nil,
		},
		{
			name:                     "successful reorg recovery with 2 recursions",
			startingBlock:            100,
			safeBlockNumber:          95,
			numberOfRecursions:       2,
			checkReorgError:          nil,
			getAllStoredBlockNumbers: []uint64{100, 98, 95},
			getAllContracts:          []string{"contract1", "contract2"},
			getBlockHeadersDB: []*types.Header{
				{
					Number: big.NewInt(98),
					Time:   99,
					Root:   common.HexToHash("0x123"),
				}, {
					Number: big.NewInt(95),
					Time:   100,
				},
			},
			getBlockHeadersL1: []*types.Header{
				{
					Number: big.NewInt(98),
					Time:   88,
					Root:   common.HexToHash("0x123"),
				}, {
					Number: big.NewInt(95),
					Time:   100,
				},
			},
			checkoutError: nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.TODO()
			stateService, tx, client, _, _, _ := createMocks(t)
			for _, header := range tt.getBlockHeadersL1 {
				client.EXPECT().HeaderByNumber(ctx, header.Number).Return(header, nil).Times(1)
			}
			stateService.EXPECT().NewTransaction().Return(tx, nil).Times(1)
			tx.EXPECT().Discard().Times(1)
			tx.EXPECT().Commit().Times(1)
			tx.EXPECT().GetAllStoredBlockNumbers().Return(tt.getAllStoredBlockNumbers, nil).Times(1)
			for i := 0; i < int(tt.numberOfRecursions); i++ {
				block := tt.getBlockHeadersDB[i]
				tx.EXPECT().GetOwnershipBlock(block.Number.Uint64()).Return(model.Block{
					Number: block.Number.Uint64(),
					Hash:   block.Hash(),
				}, nil).Times(1)
			}

			tx.EXPECT().SetLastOwnershipBlock(gomock.Any()).Return(nil).Times(1)
			tx.EXPECT().DeleteOrphanBlockData(tt.safeBlockNumber).Return(nil).Times(1)
			tx.EXPECT().DeleteOrphanRootTags(int64(tt.safeBlockNumber)+1, int64(tt.startingBlock)).Return(nil).Times(1)

			tx.EXPECT().Checkout(int64(tt.safeBlockNumber)).Return(tt.checkoutError).Times(1)

			p := universal.NewProcessor(client, stateService, nil, &config.Config{}, nil, nil)
			block, err := p.RecoverFromReorg(ctx, tt.startingBlock)
			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("RecoverFromReorg() error = %v, wantErr %v", err, tt.expectedError)
			}
			if block.Number != tt.safeBlockNumber {
				t.Errorf("RecoverFromReorg() block = %v, want %v", block, tt.safeBlockNumber)
			}
		})
	}
}

// nolint:gocritic // many return values in function => we accept this for this test helper
func createMocks(t *testing.T) (
	*mockTx.MockService,
	*mockTx.MockTx,
	*mockClient.MockEthClient,
	*mockScan.MockScanner,
	*mockDiscoverer.MockDiscoverer,
	*mockUpdater.MockUpdater,
) {
	ctrl := gomock.NewController(t)
	return mockTx.NewMockService(ctrl), mockTx.NewMockTx(ctrl), mockClient.NewMockEthClient(ctrl), mockScan.NewMockScanner(ctrl), mockDiscoverer.NewMockDiscoverer(ctrl), mockUpdater.NewMockUpdater(ctrl)
}

func assertError(t *testing.T, expectedError, err error) {
	t.Helper()
	if expectedError != nil {
		if err == nil || err.Error() != expectedError.Error() {
			t.Fatalf(`got error "%v", expected error: "%v"`, err, expectedError)
		}
	} else {
		if err != expectedError {
			t.Fatalf(`got error "%v", expected error: "%v"`, err, expectedError)
		}
	}
}
