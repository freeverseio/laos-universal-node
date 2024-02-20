package blockmapper

import (
	"context"
	"fmt"
	"testing"

	searchMock "github.com/freeverseio/laos-universal-node/internal/core/block/search/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	stateMock "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

type mocks struct {
	tx     *stateMock.MockTx
	search *searchMock.MockSearch
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
	t.Run("no mapping exists and GetFirstOwnershipBlock fails", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		lastMappedOwnershipBlock := uint64(0)
		errMsg := fmt.Errorf("err")
		expectedErr := fmt.Errorf("error occurred retrieving the first ownership block from storage: %w", errMsg)

		tx := stateMock.NewMockTx(ctrl)
		tx.EXPECT().GetFirstOwnershipBlock().Return(model.Block{}, errMsg)

		p := processor{}
		_, err := p.getInitialEvoBlock(context.Background(), lastMappedOwnershipBlock, tx)
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
	ownershipBlock := model.Block{
		Number:    ownershipStartingBlock,
		Timestamp: 1000,
	}
	evoBlock := model.Block{
		Number:    evoStartingBlock,
		Timestamp: 1500,
	}
	expectedOldestBlock := uint64(10)

	mockObjects.tx.EXPECT().GetFirstOwnershipBlock().Return(ownershipBlock, nil)
	mockObjects.tx.EXPECT().GetFirstEvoBlock().Return(evoBlock, nil)
	mockObjects.search.EXPECT().GetEvolutionBlockByTimestamp(context.Background(), ownershipBlock.Timestamp).Return(expectedOldestBlock, nil)

	p := processor{
		blockSearch: mockObjects.search,
	}
	gotOldestBlock, err := p.getOldestUserDefinedBlock(context.Background(), mockObjects.tx)
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
		name                       string
		expectedErr                error
		getFirstOwnershipBlockFunc func(tx *stateMock.MockTx)
		getFirstEvoBlockFunc       func(tx *stateMock.MockTx)
		getEvoBlockByTimestamp     func(search *searchMock.MockSearch)
	}{
		{
			name:        "GetFirstEvoBlock fails",
			expectedErr: fmt.Errorf("error occurred retrieving the first evolution block from storage: err"),
			getFirstOwnershipBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetFirstOwnershipBlock().Return(model.Block{Timestamp: uint64(1500), Number: uint64(100)}, nil)
			},
			getFirstEvoBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetFirstEvoBlock().Return(model.Block{}, fmt.Errorf("err"))
			},
			getEvoBlockByTimestamp: func(search *searchMock.MockSearch) {},
		},
		{
			name:        "search GetEvolutionBlockByTimestamp fails",
			expectedErr: fmt.Errorf("error occurred searching for evolution block number by target timestamp 1000 (ownership block number 100): err"),
			getFirstOwnershipBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetFirstOwnershipBlock().Return(model.Block{Timestamp: uint64(1000), Number: uint64(100)}, nil)
			},
			getFirstEvoBlockFunc: func(tx *stateMock.MockTx) {
				tx.EXPECT().GetFirstEvoBlock().Return(model.Block{Timestamp: uint64(1500), Number: uint64(20)}, nil)
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
			tt.getFirstOwnershipBlockFunc(mockObjects.tx)
			tt.getFirstEvoBlockFunc(mockObjects.tx)
			tt.getEvoBlockByTimestamp(mockObjects.search)

			p := processor{
				blockSearch: mockObjects.search,
			}
			_, err := p.getOldestUserDefinedBlock(context.Background(), mockObjects.tx)
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
		tx:     stateMock.NewMockTx(ctrl),
		search: searchMock.NewMockSearch(ctrl),
	}
}
