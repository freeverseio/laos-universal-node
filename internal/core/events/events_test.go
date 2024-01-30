package events_test

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/core/events"

	mockClient "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	gomock "go.uber.org/mock/gomock"
)

func TestFilterEventLogsSuccess(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Define your test cases
	tests := []struct {
		name             string
		firstBlock       *big.Int
		lastBlock        *big.Int
		firstBlockEvo    uint64
		lastBlockEvo     uint64
		ownershipLogs    []types.Log
		evoLogs          []types.Log
		expectedLogs     []types.Log
		expectError      bool
		expectedLogCount int
	}{
		{
			name:          "Success case",
			firstBlock:    big.NewInt(100),
			lastBlock:     big.NewInt(200),
			firstBlockEvo: uint64(10),
			lastBlockEvo:  uint64(20),
			ownershipLogs: []types.Log{
				{BlockNumber: 100, Address: common.HexToAddress("0xAddress")},
			},
			evoLogs: []types.Log{
				{BlockNumber: 10, Address: common.HexToAddress("0xAddress")},
			},
			expectedLogs: []types.Log{
				{BlockNumber: 100, Address: common.HexToAddress("0xAddress")},
				{BlockNumber: 10, Address: common.HexToAddress("0xAddress")},
			},
			expectError: false,
		},
		{
			name:          "Success case with no logs",
			firstBlock:    big.NewInt(100),
			lastBlock:     big.NewInt(200),
			firstBlockEvo: uint64(10),
			lastBlockEvo:  uint64(20),
			ownershipLogs: []types.Log{},
			evoLogs:       []types.Log{},
			expectedLogs:  []types.Log{},
			expectError:   false,
		},
		{
			name:          "Success case with several logs",
			firstBlock:    big.NewInt(100),
			lastBlock:     big.NewInt(200),
			firstBlockEvo: uint64(100),
			lastBlockEvo:  uint64(200),
			ownershipLogs: []types.Log{
				{BlockNumber: 100, Address: common.HexToAddress("0xAddress")},
				{BlockNumber: 120, Address: common.HexToAddress("0xAddress")},
				{BlockNumber: 125, Address: common.HexToAddress("0xAddress")},
				{BlockNumber: 128, Address: common.HexToAddress("0xAddress")},
			},
			evoLogs: []types.Log{
				{BlockNumber: 150, Address: common.HexToAddress("0xAddress")},
			},
			expectedLogs: []types.Log{
				{BlockNumber: 100, Address: common.HexToAddress("0xAddress")},
				{BlockNumber: 120, Address: common.HexToAddress("0xAddress")},
				{BlockNumber: 125, Address: common.HexToAddress("0xAddress")},
				{BlockNumber: 128, Address: common.HexToAddress("0xAddress")},
				{BlockNumber: 150, Address: common.HexToAddress("0xAddress")},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockStateService := mockTx.NewMockService(ctrl)
			mockOwnershipClient := mockClient.NewMockEthClient(ctrl)
			mockEvoClient := mockClient.NewMockEthClient(ctrl)
			mockTrans := mockTx.NewMockTx(ctrl)

			e := events.NewEvents(mockOwnershipClient, mockEvoClient, mockStateService, common.Address{})

			ctx := context.TODO()
			contracts := []common.Address{common.HexToAddress("0xAddress")}
			topics := [][]common.Hash{}

			mockOwnershipClient.EXPECT().
				FilterLogs(ctx, ethereum.FilterQuery{
					FromBlock: tt.firstBlock,
					ToBlock:   tt.lastBlock,
					Addresses: contracts,
					Topics:    topics,
				}).
				Return(tt.ownershipLogs, nil)

			mockStateService.EXPECT().NewTransaction().Return(mockTrans)
			mockTrans.EXPECT().GetCorrespondingEvoBlockNumber(uint64(tt.firstBlock.Int64())).Return(tt.firstBlockEvo, nil)
			mockTrans.EXPECT().GetCorrespondingEvoBlockNumber(uint64(tt.lastBlock.Int64())).Return(tt.lastBlockEvo, nil)
			mockEvoClient.EXPECT().
				FilterLogs(ctx, ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(tt.firstBlockEvo)),
					ToBlock:   big.NewInt(int64(tt.lastBlockEvo)),
					Addresses: contracts,
					Topics:    topics,
				}).
				Return(tt.evoLogs, nil)

			logs, err := e.FilterEventLogs(ctx, tt.firstBlock, tt.lastBlock, topics, contracts...)
			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error: %v, expectError %v", err, tt.expectError)
			}

			assertLogsEqual(t, logs, tt.expectedLogs)

		})
	}
}

func TestFilterEventLogsError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		setupMock   func(mockOwnershipClient *mockClient.MockEthClient, mockEvoClient *mockClient.MockEthClient, mockStateService *mockTx.MockService, mockTrans *mockTx.MockTx)
		firstBlock  *big.Int
		lastBlock   *big.Int
		expectError bool
	}{
		{
			name: "Error from ownership client FilterLogs",
			setupMock: func(mockOwnershipClient *mockClient.MockEthClient, mockEvoClient *mockClient.MockEthClient, mockStateService *mockTx.MockService, mockTrans *mockTx.MockTx) {
				mockOwnershipClient.EXPECT().
					FilterLogs(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("filter logs error"))
			},
			firstBlock:  big.NewInt(100),
			lastBlock:   big.NewInt(200),
			expectError: true,
		},
		{
			name: "Error from evo client FilterLogs",
			setupMock: func(mockOwnershipClient *mockClient.MockEthClient, mockEvoClient *mockClient.MockEthClient, mockStateService *mockTx.MockService, mockTrans *mockTx.MockTx) {
				mockOwnershipClient.EXPECT().
					FilterLogs(gomock.Any(), gomock.Any()).
					Return([]types.Log{}, nil)
				mockStateService.EXPECT().NewTransaction().Return(mockTrans)
				mockTrans.EXPECT().GetCorrespondingEvoBlockNumber(gomock.Any()).Return(uint64(10), nil).AnyTimes()
				mockEvoClient.EXPECT().
					FilterLogs(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("filter logs error"))
			},
			firstBlock:  big.NewInt(100),
			lastBlock:   big.NewInt(200),
			expectError: true,
		},
		{
			name: "Error from stateService GetCorrespondingEvoBlockNumber",
			setupMock: func(mockOwnershipClient *mockClient.MockEthClient, mockEvoClient *mockClient.MockEthClient, mockStateService *mockTx.MockService, mockTrans *mockTx.MockTx) {
				mockOwnershipClient.EXPECT().
					FilterLogs(gomock.Any(), gomock.Any()).
					Return([]types.Log{}, nil)
				mockStateService.EXPECT().NewTransaction().Return(mockTrans)
				mockTrans.EXPECT().GetCorrespondingEvoBlockNumber(gomock.Any()).Return(uint64(10), errors.New("error getting corresponding evo block number"))
			},
			firstBlock:  big.NewInt(100),
			lastBlock:   big.NewInt(200),
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockStateService := mockTx.NewMockService(ctrl)
			mockOwnershipClient := mockClient.NewMockEthClient(ctrl)
			mockEvoClient := mockClient.NewMockEthClient(ctrl)
			mockTrans := mockTx.NewMockTx(ctrl)

			tt.setupMock(mockOwnershipClient, mockEvoClient, mockStateService, mockTrans)

			e := events.NewEvents(mockOwnershipClient, mockEvoClient, mockStateService, common.Address{})

			ctx := context.TODO()
			topics := [][]common.Hash{}
			contracts := []common.Address{common.HexToAddress("0xAddress")}

			_, err := e.FilterEventLogs(ctx, tt.firstBlock, tt.lastBlock, topics, contracts...)
			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error result: %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func assertLogsEqual(t *testing.T, logs []types.Log, expectedLogs []types.Log) {
	if len(logs) != len(expectedLogs) {
		t.Errorf("expected %d logs, got %d", len(expectedLogs), len(logs))
		return
	}

	for i := range logs {
		if logs[i].BlockNumber != expectedLogs[i].BlockNumber ||
			logs[i].Address != expectedLogs[i].Address ||
			!reflect.DeepEqual(logs[i].Topics, expectedLogs[i].Topics) ||
			logs[i].TxHash != expectedLogs[i].TxHash {
			t.Errorf("logs[%d] = %+v, expected %+v", i, logs[i], expectedLogs[i])
		}
	}
}
