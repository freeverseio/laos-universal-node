package blockmapper

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	searchMock "github.com/freeverseio/laos-universal-node/internal/core/block/search/mock"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/mock"
	clientMock "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	stateMock "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

type mocks struct {
	ownClient      *clientMock.MockEthClient
	evoClient      *clientMock.MockEthClient
	ownBlockHelper *mock.MockBlockHelper
	evoBlockHelper *mock.MockBlockHelper
	search         *searchMock.MockSearch
}

func TestGetInitialEvoBlockError(t *testing.T) {
	t.Parallel()
	t.Run("GetMappedEvoBlockNumber fails", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		lastMappedOwnershipBlock := uint64(10)
		errMsg := fmt.Errorf("err")
		expectedErr := fmt.Errorf("error occurred retrieving the mapped evolution block number by ownership block %d from storage: %w",
			lastMappedOwnershipBlock, errMsg)

		tx := stateMock.NewMockTx(ctrl)
		tx.EXPECT().GetMappedEvoBlockNumber(lastMappedOwnershipBlock).Return(uint64(0), errMsg)

		p := processor{}
		_, err := p.getInitialEvoBlock(context.Background(), lastMappedOwnershipBlock, tx)
		if err == nil || err.Error() != expectedErr.Error() {
			t.Errorf("got error '%v' while error '%v' was expected", err, expectedErr)
		}
	})
	t.Run("no mapping exists and GetOwnershipInitStartingBlock fails", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		ownershipBlockHelper := mock.NewMockBlockHelper(ctrl)
		defer ctrl.Finish()
		lastMappedOwnershipBlock := uint64(0)
		errMsg := fmt.Errorf("err")
		expectedErr := fmt.Errorf("error occurred retrieving the ownership init starting block: %w", errMsg)

		ownershipBlockHelper.EXPECT().GetOwnershipInitStartingBlock(context.Background()).Return(uint64(0), errMsg)

		p := processor{
			ownershipBlockHelper: ownershipBlockHelper,
		}
		_, err := p.getInitialEvoBlock(context.Background(), lastMappedOwnershipBlock, stateMock.NewMockTx(ctrl))
		if err == nil || err.Error() != expectedErr.Error() {
			t.Errorf("got error '%v' while error '%v' was expected", err, expectedErr)
		}
	})
}

func TestGetOldestUserDefinedBlock(t *testing.T) {
	t.Parallel()
	ctrl, mockObjects := getMocks(t)
	defer ctrl.Finish()
	ownershipStartingBlock := uint64(100)
	evoStartingBlock := uint64(20)
	ownershipHeader := &types.Header{
		Number: big.NewInt(int64(ownershipStartingBlock)),
		Time:   1000,
	}
	evoHeader := &types.Header{
		Number: big.NewInt(int64(evoStartingBlock)),
		Time:   1500,
	}
	expectedOldestBlock := uint64(10)

	mockObjects.ownBlockHelper.EXPECT().GetOwnershipInitStartingBlock(context.Background()).Return(ownershipStartingBlock, nil)
	mockObjects.evoBlockHelper.EXPECT().GetEvoInitStartingBlock(context.Background()).Return(evoStartingBlock, nil)
	mockObjects.ownClient.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(ownershipStartingBlock))).Return(ownershipHeader, nil)
	mockObjects.evoClient.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(evoStartingBlock))).Return(evoHeader, nil)
	mockObjects.search.EXPECT().GetEvolutionBlockByTimestamp(context.Background(), ownershipHeader.Time).Return(expectedOldestBlock, nil)

	p := processor{
		ownershipClient:      mockObjects.ownClient,
		evoClient:            mockObjects.evoClient,
		ownershipBlockHelper: mockObjects.ownBlockHelper,
		evoBlockHelper:       mockObjects.evoBlockHelper,
		blockSearch:          mockObjects.search,
	}
	gotOldestBlock, err := p.getOldestUserDefinedBlock(context.Background())
	if err != nil {
		t.Errorf("got error '%v' while no error was expected", err)
	}
	if gotOldestBlock != expectedOldestBlock {
		t.Errorf("got oldest block '%d', expected '%d'", gotOldestBlock, expectedOldestBlock)
	}
}

func TestGetOldestUserDefinedBlockError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                     string
		expectedErr              error
		getOwnStartingBlockFunc  func(blockHelper *mock.MockBlockHelper)
		getEvoStartingBlockFunc  func(blockHelper *mock.MockBlockHelper)
		getOwnHeaderByNumberFunc func(client *clientMock.MockEthClient)
		getEvoHeaderByNumberFunc func(client *clientMock.MockEthClient)
		getEvoBlockByTimestamp   func(search *searchMock.MockSearch)
	}{
		{
			name:        "GetEvoInitStartingBlock fails",
			expectedErr: fmt.Errorf("error occurred retrieving the evolution init starting block: err"),
			getOwnStartingBlockFunc: func(blockHelper *mock.MockBlockHelper) {
				blockHelper.EXPECT().GetOwnershipInitStartingBlock(context.Background()).Return(uint64(100), nil)
			},
			getEvoStartingBlockFunc: func(blockHelper *mock.MockBlockHelper) {
				blockHelper.EXPECT().GetEvoInitStartingBlock(context.Background()).Return(uint64(0), fmt.Errorf("err"))
			},
			getOwnHeaderByNumberFunc: func(client *clientMock.MockEthClient) {},
			getEvoHeaderByNumberFunc: func(client *clientMock.MockEthClient) {},
			getEvoBlockByTimestamp:   func(search *searchMock.MockSearch) {},
		},
		{
			name:        "ownership client HeaderByNumber fails",
			expectedErr: fmt.Errorf("error occurred retrieving block number 100 from ownership chain: err"),
			getOwnStartingBlockFunc: func(blockHelper *mock.MockBlockHelper) {
				blockHelper.EXPECT().GetOwnershipInitStartingBlock(context.Background()).Return(uint64(100), nil)
			},
			getEvoStartingBlockFunc: func(blockHelper *mock.MockBlockHelper) {
				blockHelper.EXPECT().GetEvoInitStartingBlock(context.Background()).Return(uint64(20), nil)
			},
			getOwnHeaderByNumberFunc: func(client *clientMock.MockEthClient) {
				client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(100))).Return(nil, fmt.Errorf("err"))
			},
			getEvoHeaderByNumberFunc: func(client *clientMock.MockEthClient) {},
			getEvoBlockByTimestamp:   func(search *searchMock.MockSearch) {},
		},
		{
			name:        "evo client HeaderByNumber fails",
			expectedErr: fmt.Errorf("error occurred retrieving block number 20 from evolution chain: err"),
			getOwnStartingBlockFunc: func(blockHelper *mock.MockBlockHelper) {
				blockHelper.EXPECT().GetOwnershipInitStartingBlock(context.Background()).Return(uint64(100), nil)
			},
			getEvoStartingBlockFunc: func(blockHelper *mock.MockBlockHelper) {
				blockHelper.EXPECT().GetEvoInitStartingBlock(context.Background()).Return(uint64(20), nil)
			},
			getOwnHeaderByNumberFunc: func(client *clientMock.MockEthClient) {
				client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(100))).Return(&types.Header{}, nil)
			},
			getEvoHeaderByNumberFunc: func(client *clientMock.MockEthClient) {
				client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(20))).Return(nil, fmt.Errorf("err"))
			},
			getEvoBlockByTimestamp: func(search *searchMock.MockSearch) {},
		},
		{
			name:        "search GetEvolutionBlockByTimestamp fails",
			expectedErr: fmt.Errorf("error occurred searching for evolution block number by target timestamp 1000 (ownership block number 100): err"),
			getOwnStartingBlockFunc: func(blockHelper *mock.MockBlockHelper) {
				blockHelper.EXPECT().GetOwnershipInitStartingBlock(context.Background()).Return(uint64(100), nil)
			},
			getEvoStartingBlockFunc: func(blockHelper *mock.MockBlockHelper) {
				blockHelper.EXPECT().GetEvoInitStartingBlock(context.Background()).Return(uint64(20), nil)
			},
			getOwnHeaderByNumberFunc: func(client *clientMock.MockEthClient) {
				client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(100))).Return(&types.Header{Time: uint64(1000)}, nil)
			},
			getEvoHeaderByNumberFunc: func(client *clientMock.MockEthClient) {
				client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(20))).Return(&types.Header{Time: uint64(1500)}, nil)
			},
			getEvoBlockByTimestamp: func(search *searchMock.MockSearch) {
				search.EXPECT().GetEvolutionBlockByTimestamp(context.Background(), uint64(1000)).Return(uint64(0), fmt.Errorf("err"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl, mockObjects := getMocks(t)
			defer ctrl.Finish()
			tt.getOwnStartingBlockFunc(mockObjects.ownBlockHelper)
			tt.getEvoStartingBlockFunc(mockObjects.evoBlockHelper)
			tt.getOwnHeaderByNumberFunc(mockObjects.ownClient)
			tt.getEvoHeaderByNumberFunc(mockObjects.evoClient)
			tt.getEvoBlockByTimestamp(mockObjects.search)

			p := processor{
				ownershipClient:      mockObjects.ownClient,
				evoClient:            mockObjects.evoClient,
				ownershipBlockHelper: mockObjects.ownBlockHelper,
				evoBlockHelper:       mockObjects.evoBlockHelper,
				blockSearch:          mockObjects.search,
			}
			_, err := p.getOldestUserDefinedBlock(context.Background())
			if err == nil || err.Error() != tt.expectedErr.Error() {
				t.Errorf("got error '%v', expected '%v'", err, tt.expectedErr)
			}
		})
	}
}

func getMocks(t *testing.T) (ctrl *gomock.Controller, objects *mocks) {
	t.Helper()
	ctrl = gomock.NewController(t)
	return ctrl, &mocks{
		ownClient:      clientMock.NewMockEthClient(ctrl),
		evoClient:      clientMock.NewMockEthClient(ctrl),
		ownBlockHelper: mock.NewMockBlockHelper(ctrl),
		evoBlockHelper: mock.NewMockBlockHelper(ctrl),
		search:         searchMock.NewMockSearch(ctrl),
	}
}
