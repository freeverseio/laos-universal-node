package blockmapper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/worker/blockmapper"
	"go.uber.org/mock/gomock"
)

func TestSearchBlockByTimestamp(t *testing.T) {
	// t.Run("target timestamp is in the first half of the block range", func(t *testing.T) {
	// 	ctrl := gomock.NewController(t)
	// 	defer ctrl.Finish()

	// 	client := mockClient.NewMockEthClient(ctrl)

	// 	// Setup your worker with the mock client
	// 	conf := &config.Config{WaitingTime: 5 * time.Second}
	// 	w := blockmapper.New(conf, client)

	// 	// Mock responses
	// 	latestBlock := uint64(100)
	// 	latestBlockHeader := &types.Header{
	// 		Number: big.NewInt(int64(latestBlock)),
	// 		Time:   100000,
	// 	}
	// 	midBlockHeader := &types.Header{
	// 		Number: big.NewInt(50),
	// 		Time:   95000,
	// 	}
	// 	client.EXPECT().HeaderByNumber(context.Background(), nil).Return(latestBlockHeader, nil).Times(1)
	// 	client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(50)).Return(midBlockHeader, nil).Times(1)

	// 	// Test for a timestamp that should return the mid block
	// 	targetTimestamp := int64(95000)
	// 	expectedBlockNumber := uint64(50)

	// 	blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp)
	// 	if err != nil {
	// 		t.Errorf("Unexpected error: %v", err)
	// 	}
	// 	if blockNumber != expectedBlockNumber {
	// 		t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
	// 	}
	// })
	t.Run("with timestamp 1707487609", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client, err := ethclient.Dial("https://polygon-mainnet.gorengine.com/")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Setup your worker with the mock client
		conf := &config.Config{WaitingTime: 5 * time.Second}
		w := blockmapper.New(conf, client)

		// Test for a timestamp that should return the mid block
		// 1707491201 ts own -> 53316112
		// 1707491202 ts evo -> 510371
		targetTimestamp := int64(1707491201)
		expectedBlockNumber := uint64(53316113)
		startTime := time.Now()
		blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp, client, blockmapper.OwershipBlockCorrectionFunc)
		duration := time.Since(startTime)
		fmt.Println("Duration: ", duration)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if blockNumber != expectedBlockNumber {
			t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
		}
	})
	t.Run("with timestamp 1707487609", func(t *testing.T) {

		client, err := ethclient.Dial("https://rpc.klaos.laosfoundation.io/")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Setup your worker with the mock client
		conf := &config.Config{WaitingTime: 5 * time.Second}
		w := blockmapper.New(conf, client)

		// Test for a timestamp that should return the mid block
		targetTimestamp := int64(1707491202)
		expectedBlockNumber := uint64(510371)
		startTime := time.Now()
		blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp, client, blockmapper.EvoChainBlockCorrectionFunc)
		duration := time.Since(startTime)
		fmt.Println("Duration: ", duration)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if blockNumber != expectedBlockNumber {
			t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
		}
	})
}
