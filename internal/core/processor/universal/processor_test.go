package universal_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/universal"
	mockDiscoverer "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/discoverer/mock"
	mockUpdater "github.com/freeverseio/laos-universal-node/internal/core/processor/universal/updater/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	mockScan "github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"
	"go.uber.org/mock/gomock"

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

			stateService.EXPECT().NewTransaction().Return(tx)
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

func TestVerifyChainConsistency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		startingBlock          uint64
		lastBlockDB            model.Block
		lastBlockDBError       error
		previousBlockData      *types.Block
		previousBlockDataError error
		expectedError          error
	}{
		{
			name:             "No previous block hash in database",
			startingBlock:    100,
			lastBlockDB:      model.Block{},
			lastBlockDBError: nil,
		},
		{
			name:             "error when reading database",
			startingBlock:    100,
			lastBlockDB:      model.Block{},
			lastBlockDBError: errors.New("error from storage"),
			expectedError:    errors.New("error from storage"),
		},

		{
			name:              "Previous block hash matches",
			startingBlock:     100,
			lastBlockDB:       model.Block{Hash: common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e")},
			lastBlockDBError:  nil,
			previousBlockData: types.NewBlockWithHeader(&types.Header{ParentHash: common.HexToHash("0x123")}),
			expectedError:     nil,
		},

		{
			name:                   "error when trying to obtain previous block from chain",
			startingBlock:          100,
			lastBlockDB:            model.Block{Hash: common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e")},
			lastBlockDBError:       nil,
			previousBlockData:      nil,
			previousBlockDataError: errors.New("error retrieving previous block from chain"),
			expectedError:          errors.New("error retrieving previous block from chain"),
		},

		{
			name:                   "Previous block hash does not match",
			startingBlock:          100,
			lastBlockDB:            model.Block{Hash: common.HexToHash("0x123")},
			lastBlockDBError:       nil,
			previousBlockData:      types.NewBlockWithHeader(&types.Header{ParentHash: common.HexToHash("0x123")}),
			previousBlockDataError: nil,
			expectedError: universal.ReorgError{
				Block:       99,
				ChainHash:   common.HexToHash("0x558af54aec2a3b01640511cfc1d2b5772373b7b73ff621225031de3cae9a2c3e"),
				StorageHash: common.HexToHash("0x123"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.TODO()
			stateService, tx, client, _, _, _ := createMocks(t)

			stateService.EXPECT().NewTransaction().Return(tx)
			tx.EXPECT().GetLastOwnershipBlock().Return(tt.lastBlockDB, tt.lastBlockDBError)
			tx.EXPECT().Discard()

			if tt.lastBlockDBError == nil && tt.lastBlockDB.Hash != (common.Hash{}) {
				client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.startingBlock-1))).
					Return(tt.previousBlockData, tt.previousBlockDataError)
			}

			p := universal.NewProcessor(client, stateService, nil, &config.Config{}, nil, nil)
			err := p.VerifyChainConsistency(ctx, tt.startingBlock)
			assertError(t, tt.expectedError, err)
		})
	}
}

func TestProcessUniversalBlockRange(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	stateService, tx, client, scanner, discoverer, updater := createMocks(t)

	p := universal.NewProcessor(client, stateService, scanner, &config.Config{}, discoverer, updater)

	startingBlock := uint64(100)
	stateService.EXPECT().NewTransaction().Return(tx)

	block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(100)})
	blockData := model.Block{Number: 100, Hash: common.HexToHash("0xb07e1289b32edefd8f3c702d016fb73c81d5950b2ebc790ad9d2cb8219066b4c")}

	client.EXPECT().BlockByNumber(ctx, big.NewInt(100)).Return(block, nil)
	tx.EXPECT().SetLastOwnershipBlock(blockData).Return(nil)
	discoverer.EXPECT().ShouldDiscover(tx, startingBlock, blockData.Number).Return(false, nil)
	discoverer.EXPECT().GetContracts(tx).Return([]string{"contract"}, nil)

	updater.EXPECT().GetModelTransferEvents(ctx, startingBlock, blockData.Number, []string{"contract"}).Return(map[string][]model.ERC721Transfer{"contract": {}}, nil)
	updater.EXPECT().UpdateState(ctx, tx, []string{"contract"}, map[string][]model.ERC721Transfer{"contract": {}}, blockData).Return(nil)

	tx.EXPECT().Commit().Return(nil)
	tx.EXPECT().Discard()
	err := p.ProcessUniversalBlockRange(ctx, startingBlock, blockData.Number)
	assertError(t, nil, err)
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
			stateService.EXPECT().NewTransaction().Return(tx)
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

// nolint:gocritic // many return values in function bug it is ok
func createMocks(t *testing.T) (
	*mockTx.MockService,
	*mockTx.MockTx,
	*mockScan.MockEthClient,
	*mockScan.MockScanner,
	*mockDiscoverer.MockDiscoverer,
	*mockUpdater.MockUpdater,
) {
	ctrl := gomock.NewController(t)
	return mockTx.NewMockService(ctrl), mockTx.NewMockTx(ctrl), mockScan.NewMockEthClient(ctrl), mockScan.NewMockScanner(ctrl), mockDiscoverer.NewMockDiscoverer(ctrl), mockUpdater.NewMockUpdater(ctrl)
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
