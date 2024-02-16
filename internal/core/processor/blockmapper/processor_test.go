package blockmapper_test

import (
	"fmt"
	"testing"

	"github.com/freeverseio/laos-universal-node/internal/config"
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

			mockOwnershipClient := clientMock.NewMockEthClient(ctrl)
			mockEvoClient := clientMock.NewMockEthClient(ctrl)
			mockStateService := stateMock.NewMockService(ctrl)
			tx := stateMock.NewMockTx(ctrl)

			mockStateService.EXPECT().NewTransaction().Return(tx, nil)
			tx.EXPECT().Discard()
			tx.EXPECT().GetLastMappedOwnershipBlockNumber().Return(tt.lastMappedBlock, nil)
			tx.EXPECT().GetLastOwnershipBlock().Return(tt.lastProcessedBlock, nil)

			processor := blockmapper.New(&config.Config{}, mockOwnershipClient, mockEvoClient, mockStateService)
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

			processor := blockmapper.New(&config.Config{}, mockOwnershipClient, mockEvoClient, mockStateService)
			_, err := processor.IsMappingSyncedWithProcessing()
			if err == nil || err.Error() != tt.expectedErr.Error() {
				t.Errorf("got error '%v', expected '%v'", err, tt.expectedErr)
			}
		})
	}
}
