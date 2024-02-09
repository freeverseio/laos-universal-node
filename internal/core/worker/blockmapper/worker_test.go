package blockmapper_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/worker/blockmapper"
	mockClient "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	"go.uber.org/mock/gomock"
)

func TestSearchBlockByTimestamp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mockClient.NewMockEthClient(ctrl)

	// Setup your worker with the mock client
	conf := &config.Config{WaitingTime: 5 * time.Second}
	w := blockmapper.New(conf, client)

	// Mock responses
	latestBlock := uint64(100)
	latestBlockHeader := &types.Header{
		Number: big.NewInt(int64(latestBlock)),
		Time:   100000,
	}
	midBlockHeader := &types.Header{
		Number: big.NewInt(50),
		Time:   95000,
	}
	client.EXPECT().HeaderByNumber(context.Background(), nil).Return(latestBlockHeader, nil).Times(1)
	client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(50)).Return(midBlockHeader, nil).Times(1)

	// Test for a timestamp that should return the mid block
	targetTimestamp := int64(95000)
	expectedBlockNumber := uint64(50)

	blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if blockNumber != expectedBlockNumber {
		t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
	}
}
