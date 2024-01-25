package events_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/core/events"

	mockClient "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	gomock "go.uber.org/mock/gomock"
)

func TestFilterEventLogs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStateService := mockTx.NewMockService(ctrl)
	mockOwnershipClient := mockClient.NewMockEthClient(ctrl)
	mockEvoClient := mockClient.NewMockEthClient(ctrl)
	mockTrans := mockTx.NewMockTx(ctrl)

	e := events.NewEvents(mockOwnershipClient, mockEvoClient, mockStateService, common.Address{})

	ctx := context.TODO()
	firstBlock := big.NewInt(100)
	lastBlock := big.NewInt(200)
	topics := [][]common.Hash{}
	contracts := []common.Address{common.HexToAddress("0xAddress")}
	expectedLogs := []types.Log{{}, {}} // Assuming 2 logs for simplicity

	mockOwnershipClient.EXPECT().
		FilterLogs(ctx, gomock.Any()).
		Return(expectedLogs[:1], nil)

	mockStateService.EXPECT().NewTransaction().Return(mockTrans)

	mockTrans.EXPECT().GetCorrespondingEvoBlockNumber(gomock.Any()).
		Return(uint64(0), nil).Times(2)

	mockEvoClient.EXPECT().
		FilterLogs(ctx, gomock.Any()).
		Return(expectedLogs[1:], nil)

	logs, err := e.FilterEventLogs(ctx, firstBlock, lastBlock, topics, contracts...)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(logs) != len(expectedLogs) {
		t.Errorf("expected %d logs, got %d", len(expectedLogs), len(logs))
	}
}
