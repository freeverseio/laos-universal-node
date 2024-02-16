package blockmapper_test

import (
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

			tx.EXPECT().Discard()
			mockStateService.EXPECT().NewTransaction().Return(tx, nil)
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
