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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOwnershipClient := clientMock.NewMockEthClient(ctrl)
			mockEvoClient := clientMock.NewMockEthClient(ctrl)
			mockStateService := stateMock.NewMockService(ctrl)
			tx := stateMock.NewMockTx(ctrl)

			tt.newTransactionFunc(mockStateService, tx)
			tt.getLastMappedBlockFunc(tx)
			tt.getLastOwnershipBlockFunc(tx)

			processor := blockmapper.New(mockOwnershipClient, mockEvoClient, mockStateService)
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ownClient := clientMock.NewMockEthClient(ctrl)
	evoClient := clientMock.NewMockEthClient(ctrl)
	stateService := stateMock.NewMockService(ctrl)
	search := searchMock.NewMockSearch(ctrl)
	tx := stateMock.NewMockTx(ctrl)

	stateService.EXPECT().NewTransaction().Return(tx, nil)
	tx.EXPECT().Discard()
	tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(lastMappedOwnershipBlock, nil)
	tx.EXPECT().GetMappedEvoBlockNumber(uint64(99)).Return(mappedEvoBlock, nil)
	ownClient.EXPECT().HeaderByNumber(context.Background(), big.NewInt(int64(nextOwnershipBlock))).Return(&nextOwnershipBlockHeader, nil)
	search.EXPECT().GetEvolutionBlockByTimestamp(context.Background(), nextOwnershipBlockHeader.Time, mappedEvoBlock).Return(toMapEvoBlock, nil)
	tx.EXPECT().SetOwnershipEvoBlockMapping(nextOwnershipBlock, toMapEvoBlock).Return(nil)
	tx.EXPECT().SetLastMappedOwnershipBlockNumber(nextOwnershipBlock).Return(nil)
	tx.EXPECT().Commit().Return(nil)

	processor := blockmapper.New(ownClient, evoClient, stateService, blockmapper.WithBlockSearch(search))
	err := processor.MapNextBlock(context.Background())
	if err != nil {
		t.Errorf("got error '%v' while no error was expected", err)
	}
}
