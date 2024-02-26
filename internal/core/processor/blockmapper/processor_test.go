package blockmapper_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	searchMock "github.com/freeverseio/laos-universal-node/internal/core/block/search/mock"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/blockmapper"
	clientMock "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	stateMock "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

type mocks struct {
	ownClient *clientMock.MockEthClient
	evoClient *clientMock.MockEthClient
	state     *stateMock.MockService
	search    *searchMock.MockSearch
	tx        *stateMock.MockTx
}

func TestIsMappingSyncedWithProcessing(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		lastMappedBlock    uint64
		lastProcessedBlock model.Block
		expectedSynced     bool
	}{
		{
			name:               "mapping is synced with processing",
			lastMappedBlock:    10,
			lastProcessedBlock: model.Block{Number: 10},
			expectedSynced:     true,
		},
		{
			name:               "mapping is behind processing",
			lastMappedBlock:    7,
			lastProcessedBlock: model.Block{Number: 10},
			expectedSynced:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ownClient := clientMock.NewMockEthClient(ctrl)
			evoClient := clientMock.NewMockEthClient(ctrl)
			stateService := stateMock.NewMockService(ctrl)
			tx := stateMock.NewMockTx(ctrl)

			stateService.EXPECT().NewTransaction().Return(tx, nil)
			tx.EXPECT().Discard()
			tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(tt.lastMappedBlock, nil)
			tx.EXPECT().GetLastOwnershipBlock().Return(tt.lastProcessedBlock, nil)

			processor := blockmapper.New(ownClient, evoClient, stateService)
			synced, err := processor.IsMappingSyncedWithProcessing()
			if err != nil {
				t.Errorf("got error '%v' while no error was expected", err)
			}

			if synced != tt.expectedSynced {
				t.Errorf("got synced %v, expected %v", synced, tt.expectedSynced)
			}
		})
	}
}

func TestIsMappingSyncedWithProcessingError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                      string
		expectedErr               error
		newTransactionFunc        func(*stateMock.MockService, *stateMock.MockTx)
		getLastMappedBlockFunc    func(*stateMock.MockTx)
		getLastOwnershipBlockFunc func(*stateMock.MockTx)
	}{
		{
			name:        "should handle NewTransaction error",
			expectedErr: fmt.Errorf("error occurred creating transaction: state service failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, fmt.Errorf("state service failed"))
			},
			getLastMappedBlockFunc:    func(*stateMock.MockTx) {},
			getLastOwnershipBlockFunc: func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle GetLastMappedOwnershipBlockNumber error",
			expectedErr: fmt.Errorf("error occurred retrieving the latest mapped ownership block from storage: storage failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				tx.EXPECT().Discard()
				s.EXPECT().NewTransaction().Return(tx, nil)
			},
			getLastMappedBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(0), fmt.Errorf("storage failed"))
			},
			getLastOwnershipBlockFunc: func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle GetLastOwnershipBlock error",
			expectedErr: fmt.Errorf("error occurred retrieving the last processed ownership block from storage: storage failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				tx.EXPECT().Discard()
				s.EXPECT().NewTransaction().Return(tx, nil)
			},
			getLastMappedBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(10), nil)
			},
			getLastOwnershipBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastOwnershipBlock().Return(model.Block{}, fmt.Errorf("storage failed"))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl, mockObjects := getMocks(t)
			defer ctrl.Finish()

			tt.newTransactionFunc(mockObjects.state, mockObjects.tx)
			tt.getLastMappedBlockFunc(mockObjects.tx)
			tt.getLastOwnershipBlockFunc(mockObjects.tx)

			processor := blockmapper.New(mockObjects.ownClient, mockObjects.evoClient, mockObjects.state)
			_, err := processor.IsMappingSyncedWithProcessing()
			if err == nil || err.Error() != tt.expectedErr.Error() {
				t.Errorf("got error '%v', expected '%v'", err, tt.expectedErr)
			}
		})
	}
}

func TestMapNextBlock(t *testing.T) {
	t.Parallel()
	lastMappedOwnershipBlock := uint64(99)
	nextOwnershipBlock := uint64(100)
	mappedEvoBlock := uint64(9)
	nextOwnershipBlockHeader := types.Header{
		Number: big.NewInt(int64(nextOwnershipBlock)),
		Time:   uint64(123456),
	}
	toMapEvoBlock := uint64(10)

	ctrl, mockObjects := getMocks(t)
	defer ctrl.Finish()

	mockObjects.state.EXPECT().NewTransaction().Return(mockObjects.tx, nil)
	mockObjects.tx.EXPECT().Discard()
	mockObjects.tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(lastMappedOwnershipBlock, nil)
	mockObjects.tx.EXPECT().GetMappedEvoBlockNumber(uint64(99)).Return(mappedEvoBlock, nil)
	mockObjects.ownClient.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(nextOwnershipBlock))).Return(&nextOwnershipBlockHeader, nil)
	mockObjects.search.EXPECT().GetEvolutionBlockByTimestamp(context.Background(), nextOwnershipBlockHeader.Time, mappedEvoBlock).Return(toMapEvoBlock, nil)
	mockObjects.tx.EXPECT().SetOwnershipEvoBlockMapping(nextOwnershipBlock, toMapEvoBlock).Return(nil)
	mockObjects.tx.EXPECT().SetLastMappedOwnershipBlockNumber(nextOwnershipBlock).Return(nil)
	mockObjects.tx.EXPECT().Commit().Return(nil)

	processor := blockmapper.New(mockObjects.ownClient, mockObjects.evoClient, mockObjects.state, blockmapper.WithBlockSearch(mockObjects.search))
	err := processor.MapNextBlock(context.Background())
	if err != nil {
		t.Errorf("got error '%v' while no error was expected", err)
	}
}

func TestMapNextBlockError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                             string
		expectedErr                      error
		newTransactionFunc               func(*stateMock.MockService, *stateMock.MockTx)
		getLastMappedOwnBlockFunc        func(*stateMock.MockTx)
		getFirstOwnershipBlockFunc       func(*stateMock.MockTx)
		getMappedEvoBlockFunc            func(*stateMock.MockTx)
		headerByNumberFunc               func(*clientMock.MockEthClient)
		getEvolutionBlockByTimestampFunc func(*searchMock.MockSearch)
		setOwnershipEvoBlockMappingFunc  func(*stateMock.MockTx)
		setLastMappedOwnBlockFunc        func(*stateMock.MockTx)
		commitFunc                       func(*stateMock.MockTx)
	}{
		{
			name:        "should handle NewTransaction error",
			expectedErr: fmt.Errorf("error occurred creating transaction: state service failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, fmt.Errorf("state service failed"))
			},
			getLastMappedOwnBlockFunc:        func(*stateMock.MockTx) {},
			getFirstOwnershipBlockFunc:       func(*stateMock.MockTx) {},
			getMappedEvoBlockFunc:            func(*stateMock.MockTx) {},
			headerByNumberFunc:               func(*clientMock.MockEthClient) {},
			getEvolutionBlockByTimestampFunc: func(*searchMock.MockSearch) {},
			setOwnershipEvoBlockMappingFunc:  func(*stateMock.MockTx) {},
			setLastMappedOwnBlockFunc:        func(*stateMock.MockTx) {},
			commitFunc:                       func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle GetLastMappedOwnershipBlockNumber error",
			expectedErr: fmt.Errorf("error occurred retrieving the latest mapped ownership block from storage: storage failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
			},
			getLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(0), fmt.Errorf("storage failed"))
			},
			getFirstOwnershipBlockFunc:       func(*stateMock.MockTx) {},
			getMappedEvoBlockFunc:            func(*stateMock.MockTx) {},
			headerByNumberFunc:               func(*clientMock.MockEthClient) {},
			getEvolutionBlockByTimestampFunc: func(*searchMock.MockSearch) {},
			setOwnershipEvoBlockMappingFunc:  func(*stateMock.MockTx) {},
			setLastMappedOwnBlockFunc:        func(*stateMock.MockTx) {},
			commitFunc:                       func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle GetFirstOwnershipBlock error",
			expectedErr: fmt.Errorf("error occurred retrieving the first ownership block from storage: storage failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
			},
			getLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(0), nil)
			},
			getFirstOwnershipBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetFirstOwnershipBlock().Return(model.Block{}, fmt.Errorf("storage failed"))
			},
			getMappedEvoBlockFunc:            func(*stateMock.MockTx) {},
			headerByNumberFunc:               func(*clientMock.MockEthClient) {},
			getEvolutionBlockByTimestampFunc: func(*searchMock.MockSearch) {},
			setOwnershipEvoBlockMappingFunc:  func(*stateMock.MockTx) {},
			setLastMappedOwnBlockFunc:        func(*stateMock.MockTx) {},
			commitFunc:                       func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle GetMappedEvoBlockNumber error",
			expectedErr: fmt.Errorf("error occurred retrieving the mapped evolution block number by ownership block 99 from storage: storage failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
			},
			getLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(99), nil)
			},
			getFirstOwnershipBlockFunc: func(*stateMock.MockTx) {},
			getMappedEvoBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetMappedEvoBlockNumber(uint64(99)).Return(uint64(0), fmt.Errorf("storage failed"))
			},
			headerByNumberFunc:               func(*clientMock.MockEthClient) {},
			getEvolutionBlockByTimestampFunc: func(*searchMock.MockSearch) {},
			setOwnershipEvoBlockMappingFunc:  func(*stateMock.MockTx) {},
			setLastMappedOwnBlockFunc:        func(*stateMock.MockTx) {},
			commitFunc:                       func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle HeaderByNumber error",
			expectedErr: fmt.Errorf("error occurred retrieving block number 100 from ownership chain: blockchain failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
			},
			getLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(99), nil)
			},
			getFirstOwnershipBlockFunc: func(*stateMock.MockTx) {},
			getMappedEvoBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetMappedEvoBlockNumber(uint64(99)).Return(uint64(9), nil)
			},
			headerByNumberFunc: func(ownClient *clientMock.MockEthClient) {
				ownClient.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(100))).Return(nil, fmt.Errorf("blockchain failed"))
			},
			getEvolutionBlockByTimestampFunc: func(*searchMock.MockSearch) {},
			setOwnershipEvoBlockMappingFunc:  func(*stateMock.MockTx) {},
			setLastMappedOwnBlockFunc:        func(*stateMock.MockTx) {},
			commitFunc:                       func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle GetEvolutionBlockByTimestamp error",
			expectedErr: fmt.Errorf("error occurred searching for evolution block number by target timestamp 123456 (ownership block number 100): search failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
			},
			getLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(99), nil)
			},
			getFirstOwnershipBlockFunc: func(*stateMock.MockTx) {},
			getMappedEvoBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetMappedEvoBlockNumber(uint64(99)).Return(uint64(9), nil)
			},
			headerByNumberFunc: func(ownClient *clientMock.MockEthClient) {
				ownClient.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(100))).Return(&types.Header{
					Number: big.NewInt(int64(100)),
					Time:   uint64(123456),
				}, nil)
			},
			getEvolutionBlockByTimestampFunc: func(search *searchMock.MockSearch) {
				search.EXPECT().GetEvolutionBlockByTimestamp(context.Background(), uint64(123456), uint64(9)).Return(uint64(0), fmt.Errorf("search failed"))
			},
			setOwnershipEvoBlockMappingFunc: func(*stateMock.MockTx) {},
			setLastMappedOwnBlockFunc:       func(*stateMock.MockTx) {},
			commitFunc:                      func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle SetOwnershipEvoBlockMapping error",
			expectedErr: fmt.Errorf("error setting ownership block number 100 (key) to evo block number 10 (value) in storage: storage failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
			},
			getLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(99), nil)
			},
			getFirstOwnershipBlockFunc: func(*stateMock.MockTx) {},
			getMappedEvoBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetMappedEvoBlockNumber(uint64(99)).Return(uint64(9), nil)
			},
			headerByNumberFunc: func(ownClient *clientMock.MockEthClient) {
				ownClient.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(100))).Return(&types.Header{
					Number: big.NewInt(int64(100)),
					Time:   uint64(123456),
				}, nil)
			},
			getEvolutionBlockByTimestampFunc: func(search *searchMock.MockSearch) {
				search.EXPECT().GetEvolutionBlockByTimestamp(context.Background(), uint64(123456), uint64(9)).Return(uint64(10), nil)
			},
			setOwnershipEvoBlockMappingFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().SetOwnershipEvoBlockMapping(uint64(100), uint64(10)).Return(fmt.Errorf("storage failed"))
			},
			setLastMappedOwnBlockFunc: func(*stateMock.MockTx) {},
			commitFunc:                func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle SetLastMappedOwnershipBlockNumber error",
			expectedErr: fmt.Errorf("error setting the last mapped ownership block number 100 in storage: storage failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
			},
			getLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(99), nil)
			},
			getFirstOwnershipBlockFunc: func(*stateMock.MockTx) {},
			getMappedEvoBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetMappedEvoBlockNumber(uint64(99)).Return(uint64(9), nil)
			},
			headerByNumberFunc: func(ownClient *clientMock.MockEthClient) {
				ownClient.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(100))).Return(&types.Header{
					Number: big.NewInt(int64(100)),
					Time:   uint64(123456),
				}, nil)
			},
			getEvolutionBlockByTimestampFunc: func(search *searchMock.MockSearch) {
				search.EXPECT().GetEvolutionBlockByTimestamp(context.Background(), uint64(123456), uint64(9)).Return(uint64(10), nil)
			},
			setOwnershipEvoBlockMappingFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().SetOwnershipEvoBlockMapping(uint64(100), uint64(10)).Return(nil)
			},
			setLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().SetLastMappedOwnershipBlockNumber(uint64(100)).Return(fmt.Errorf("storage failed"))
			},
			commitFunc: func(*stateMock.MockTx) {},
		},
		{
			name:        "should handle Commit error",
			expectedErr: fmt.Errorf("error committing transaction: storage failed"),
			newTransactionFunc: func(s *stateMock.MockService, tx *stateMock.MockTx) {
				s.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
			},
			getLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(uint64(99), nil)
			},
			getFirstOwnershipBlockFunc: func(*stateMock.MockTx) {},
			getMappedEvoBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetMappedEvoBlockNumber(uint64(99)).Return(uint64(9), nil)
			},
			headerByNumberFunc: func(ownClient *clientMock.MockEthClient) {
				ownClient.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(100))).Return(&types.Header{
					Number: big.NewInt(int64(100)),
					Time:   uint64(123456),
				}, nil)
			},
			getEvolutionBlockByTimestampFunc: func(search *searchMock.MockSearch) {
				search.EXPECT().GetEvolutionBlockByTimestamp(context.Background(), uint64(123456), uint64(9)).Return(uint64(10), nil)
			},
			setOwnershipEvoBlockMappingFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().SetOwnershipEvoBlockMapping(uint64(100), uint64(10)).Return(nil)
			},
			setLastMappedOwnBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().SetLastMappedOwnershipBlockNumber(uint64(100)).Return(nil)
			},
			commitFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().Commit().Return(fmt.Errorf("storage failed"))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl, mockObjects := getMocks(t)
			defer ctrl.Finish()

			tt.newTransactionFunc(mockObjects.state, mockObjects.tx)
			tt.getLastMappedOwnBlockFunc(mockObjects.tx)
			tt.getFirstOwnershipBlockFunc(mockObjects.tx)
			tt.getMappedEvoBlockFunc(mockObjects.tx)
			tt.headerByNumberFunc(mockObjects.ownClient)
			tt.getEvolutionBlockByTimestampFunc(mockObjects.search)
			tt.setOwnershipEvoBlockMappingFunc(mockObjects.tx)
			tt.setLastMappedOwnBlockFunc(mockObjects.tx)
			tt.commitFunc(mockObjects.tx)

			processor := blockmapper.New(mockObjects.ownClient, mockObjects.evoClient, mockObjects.state, blockmapper.WithBlockSearch(mockObjects.search))
			err := processor.MapNextBlock(context.Background())
			if err == nil || err.Error() != tt.expectedErr.Error() {
				t.Fatalf("got error '%v', want '%v'", err, tt.expectedErr)
			}
		})
	}
}

func getMocks(t *testing.T) (ctrl *gomock.Controller, mockObjects mocks) {
	ctrl = gomock.NewController(t)

	ownClient := clientMock.NewMockEthClient(ctrl)
	evoClient := clientMock.NewMockEthClient(ctrl)
	stateService := stateMock.NewMockService(ctrl)
	search := searchMock.NewMockSearch(ctrl)
	tx := stateMock.NewMockTx(ctrl)

	return ctrl, mocks{
		ownClient: ownClient,
		evoClient: evoClient,
		state:     stateService,
		search:    search,
		tx:        tx,
	}
}
