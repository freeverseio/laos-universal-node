package evolution_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/evolution"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	mockScan "github.com/freeverseio/laos-universal-node/internal/platform/scan/mock"

	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
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
			stateService, tx, client, _ := createMocks(t)

			stateService.EXPECT().NewTransaction().Return(tx)
			tx.EXPECT().GetLastEvoBlock().Return(tt.startingBlockData, tt.startingBlockError)
			tx.EXPECT().Discard()
			if tt.userProvidedBlock == 0 && tt.startingBlockData.Number == 0 && tt.startingBlockError == nil {
				client.EXPECT().BlockNumber(ctx).Return(tt.lastBlockNumberFromClient, tt.lastBlockNumberFromClientError)
			}

			p := evolution.NewProcessor(client, stateService, nil, tt.userProvidedBlock, 0, 0)
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
			_, _, client, _ := createMocks(t)

			client.EXPECT().BlockNumber(ctx).Return(tt.l1LatestBlock, tt.expectedError)

			p := evolution.NewProcessor(client, nil, nil, 0, tt.configBlocksMargin, tt.configBlocksRange)
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
			expectedError: evolution.ReorgError{
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
			stateService, tx, client, _ := createMocks(t)

			stateService.EXPECT().NewTransaction().Return(tx)
			tx.EXPECT().GetLastEvoBlock().Return(tt.lastBlockDB, tt.lastBlockDBError)
			tx.EXPECT().Discard()

			if tt.lastBlockDBError == nil && tt.lastBlockDB.Hash != (common.Hash{}) {
				client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.startingBlock-1))).
					Return(tt.previousBlockData, tt.previousBlockDataError)
			}

			p := evolution.NewProcessor(client, stateService, nil, 0, 0, 0)
			err := p.VerifyChainConsistency(ctx, tt.startingBlock)
			assertError(t, tt.expectedError, err)
		})
	}
}

func TestProcessEvoBlockRange(t *testing.T) {
	t.Parallel()

	t.Run("error when scanning for events", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		stateService, tx, client, scanner := createMocks(t)

		stateService.EXPECT().NewTransaction().Return(tx)
		tx.EXPECT().Discard()

		lastBlockData := model.Block{Number: 120, Hash: common.HexToHash("0x123"), Timestamp: 150}
		startingBlock := uint64(100)

		scanner.EXPECT().
			ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlockData.Number)), nil).
			Return(make([]scan.Event, 0), errors.New("error scanning events"))

		p := evolution.NewProcessor(client, stateService, scanner, 0, 0, 0)
		err := p.ProcessEvoBlockRange(ctx, startingBlock, lastBlockData.Number)
		assertError(t, errors.New("error scanning events"), err)
	})

	t.Run("obtained one event, error on getting events from db ", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		stateService, tx, client, scanner := createMocks(t)

		stateService.EXPECT().NewTransaction().Return(tx)
		tx.EXPECT().Discard()

		lastBlockData := model.Block{Number: 120, Hash: common.HexToHash("0x123"), Timestamp: 150}
		startingBlock := uint64(100)
		contract := common.HexToAddress("0x555")
		event, _ := createEventMintedWithExternalURI(lastBlockData.Number, contract)
		scanner.EXPECT().
			ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlockData.Number)), nil).
			Return([]scan.Event{event}, nil)

		tx.EXPECT().
			GetMintedWithExternalURIEvents(contract.String()).
			Return(nil, errors.New("error getting events from db"))

		p := evolution.NewProcessor(client, stateService, scanner, 0, 0, 0)
		err := p.ProcessEvoBlockRange(ctx, startingBlock, lastBlockData.Number)
		assertError(t, errors.New("error getting events from db"), err)
	})

	t.Run("obtained one event, error on storing events in db", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		stateService, tx, client, scanner := createMocks(t)

		stateService.EXPECT().NewTransaction().Return(tx)
		tx.EXPECT().Discard()

		lastBlockData := model.Block{Number: 120, Hash: common.HexToHash("0x123"), Timestamp: 150}
		startingBlock := uint64(100)
		contract := common.HexToAddress("0x555")
		event, adjustedEvent := createEventMintedWithExternalURI(lastBlockData.Number, contract)
		scanner.EXPECT().
			ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlockData.Number)), nil).
			Return([]scan.Event{event}, nil)

		tx.EXPECT().
			GetMintedWithExternalURIEvents(contract.String()).
			Return(nil, nil)

		tx.EXPECT().
			StoreMintedWithExternalURIEvents(contract.String(), []model.MintedWithExternalURI{adjustedEvent}).
			Return(errors.New("error storing events to db"))

		p := evolution.NewProcessor(client, stateService, scanner, 0, 0, 0)
		err := p.ProcessEvoBlockRange(ctx, startingBlock, lastBlockData.Number)
		assertError(t, errors.New("error storing events to db"), err)
	})

	t.Run("obtained one event, error when getting last block info", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		stateService, tx, client, scanner := createMocks(t)

		stateService.EXPECT().NewTransaction().Return(tx)
		tx.EXPECT().Discard()

		lastBlockData := model.Block{Number: 120, Hash: common.HexToHash("0x123"), Timestamp: 150}
		startingBlock := uint64(100)
		contract := common.HexToAddress("0x555")
		event, adjustedEvent := createEventMintedWithExternalURI(lastBlockData.Number, contract)
		scanner.EXPECT().
			ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlockData.Number)), nil).
			Return([]scan.Event{event}, nil)

		tx.EXPECT().
			GetMintedWithExternalURIEvents(contract.String()).
			Return(nil, nil)

		tx.EXPECT().
			StoreMintedWithExternalURIEvents(contract.String(), []model.MintedWithExternalURI{adjustedEvent}).
			Return(nil)

		client.EXPECT().
			BlockByNumber(ctx, big.NewInt(int64(lastBlockData.Number))).
			Return(nil, errors.New("error getting last block info"))

		p := evolution.NewProcessor(client, stateService, scanner, 0, 0, 0)
		err := p.ProcessEvoBlockRange(ctx, startingBlock, lastBlockData.Number)
		assertError(t, errors.New("error getting last block info"), err)
	})

	t.Run("obtained one event, error when storing last block info", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		stateService, tx, client, scanner := createMocks(t)

		stateService.EXPECT().NewTransaction().Return(tx)
		tx.EXPECT().Discard()

		lastBlockData := model.Block{
			Number:    120,
			Hash:      common.HexToHash("0x7ea18f6be7115ddbb51aa052f2780a1501847f4b3a444f1a6066982b7dbab6fc"),
			Timestamp: 150,
		}
		startingBlock := uint64(100)
		contract := common.HexToAddress("0x555")
		event, adjustedEvent := createEventMintedWithExternalURI(lastBlockData.Number, contract)
		scanner.EXPECT().
			ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlockData.Number)), nil).
			Return([]scan.Event{event}, nil)

		tx.EXPECT().
			GetMintedWithExternalURIEvents(contract.String()).
			Return(nil, nil)

		tx.EXPECT().
			StoreMintedWithExternalURIEvents(contract.String(), []model.MintedWithExternalURI{adjustedEvent}).
			Return(nil)

		client.EXPECT().
			BlockByNumber(ctx, big.NewInt(int64(lastBlockData.Number))).
			Return(types.NewBlockWithHeader(&types.Header{
				Time:   lastBlockData.Timestamp,
				Number: big.NewInt(int64(lastBlockData.Number)),
			}), nil)

		tx.EXPECT().SetLastEvoBlock(lastBlockData).Return(errors.New("error storing last block info"))

		p := evolution.NewProcessor(client, stateService, scanner, 0, 0, 0)
		err := p.ProcessEvoBlockRange(ctx, startingBlock, lastBlockData.Number)
		assertError(t, errors.New("error storing last block info"), err)
	})

	t.Run("obtained one event, event processed and last block updated successfully", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		stateService, tx, client, scanner := createMocks(t)

		stateService.EXPECT().NewTransaction().Return(tx)
		tx.EXPECT().Discard()

		lastBlockData := model.Block{
			Number:    120,
			Hash:      common.HexToHash("0x7ea18f6be7115ddbb51aa052f2780a1501847f4b3a444f1a6066982b7dbab6fc"),
			Timestamp: 150,
		}
		startingBlock := uint64(100)
		contract := common.HexToAddress("0x555")
		event, adjustedEvent := createEventMintedWithExternalURI(lastBlockData.Number, contract)
		scanner.EXPECT().
			ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlockData.Number)), nil).
			Return([]scan.Event{event}, nil)

		tx.EXPECT().
			GetMintedWithExternalURIEvents(contract.String()).
			Return(nil, nil)

		tx.EXPECT().
			StoreMintedWithExternalURIEvents(contract.String(), []model.MintedWithExternalURI{adjustedEvent}).
			Return(nil)

		client.EXPECT().
			BlockByNumber(ctx, big.NewInt(int64(lastBlockData.Number))).
			Return(types.NewBlockWithHeader(&types.Header{
				Time:   lastBlockData.Timestamp,
				Number: big.NewInt(int64(lastBlockData.Number)),
			}), nil)

		tx.EXPECT().SetLastEvoBlock(lastBlockData).Return(nil)
		tx.EXPECT().Commit().Return(nil)

		p := evolution.NewProcessor(client, stateService, scanner, 0, 0, 0)
		err := p.ProcessEvoBlockRange(ctx, startingBlock, lastBlockData.Number)
		assertError(t, nil, err)
	})
}

func createMocks(t *testing.T) (*mockTx.MockService, *mockTx.MockTx, *mockScan.MockEthClient, *mockScan.MockScanner) {
	ctrl := gomock.NewController(t)
	return mockTx.NewMockService(ctrl), mockTx.NewMockTx(ctrl), mockScan.NewMockEthClient(ctrl), mockScan.NewMockScanner(ctrl)
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

func createEventMintedWithExternalURI(blockNumber uint64, contract common.Address) (scan.EventMintedWithExternalURI, model.MintedWithExternalURI) {
	event := scan.EventMintedWithExternalURI{
		Slot:        big.NewInt(5),
		To:          common.HexToAddress("0x123"),
		TokenURI:    "https://www.google.com",
		TokenId:     big.NewInt(10),
		Contract:    contract,
		BlockNumber: blockNumber,
		Timestamp:   100,
	}

	adjustedEvent := model.MintedWithExternalURI{
		Slot:        big.NewInt(5),
		To:          common.HexToAddress("0x123"),
		TokenURI:    "https://www.google.com",
		TokenId:     big.NewInt(10),
		BlockNumber: blockNumber,
		Timestamp:   100,
	}
	return event, adjustedEvent
}
