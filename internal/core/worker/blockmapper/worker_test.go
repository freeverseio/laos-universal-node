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
	t.Run("target timestamp is in the first half of the block range", func(t *testing.T) {
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
		expectedBlockNumber := uint64(51)

		blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp, client, blockmapper.OwershipBlockCorrectionFunc)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if blockNumber != expectedBlockNumber {
			t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
		}
	})

	// t.Run("should find blocknumber in evo chain", func(t *testing.T) {
	// 	ctrl := gomock.NewController(t)
	// 	defer ctrl.Finish()

	// 	client := mockClient.NewMockEthClient(ctrl)

	// 	// Setup your worker with the mock client
	// 	conf := &config.Config{WaitingTime: 5 * time.Second}
	// 	w := blockmapper.New(conf, client)

	// 	// Mock responses
	// 	latestBlock := uint64(200)
	// 	latestBlockHeader := &types.Header{
	// 		Number: big.NewInt(int64(latestBlock)),
	// 		Time:   200000,
	// 	}
	// 	midBlockHeader := &types.Header{
	// 		Number: big.NewInt(100),
	// 		Time:   100000,
	// 	}
	// 	client.EXPECT().HeaderByNumber(context.Background(), nil).Return(latestBlockHeader, nil).Times(1)
	// 	client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(100)).Return(midBlockHeader, nil).Times(1)
	// 	client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(49)).Return(&types.Header{
	// 		Number: big.NewInt(49),
	// 		Time:   50000,
	// 	}, nil).Times(1)

	// 	client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(74)).Return(&types.Header{
	// 		Number: big.NewInt(74),
	// 		Time:   95000,
	// 	}, nil).Times(1)

	// 	// Test for a timestamp that should return the mid block
	// 	targetTimestamp := int64(95000)
	// 	expectedBlockNumber := uint64(74)

	// 	blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp, client, blockmapper.EvoChainBlockCorrectionFunc)
	// 	if err != nil {
	// 		t.Errorf("Unexpected error: %v", err)
	// 	}
	// 	if blockNumber != expectedBlockNumber {
	// 		t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
	// 	}
	// })
	t.Run("should find blocknumber in ownership chain with timestamp of evo on the left side", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mockClient.NewMockEthClient(ctrl)

		// Setup your worker with the mock client
		conf := &config.Config{WaitingTime: 5 * time.Second}
		w := blockmapper.New(conf, client)

		// Mock responses
		latestBlock := uint64(200)
		latestBlockHeader := &types.Header{
			Number: big.NewInt(int64(latestBlock)),
			Time:   200000,
		}

		expectedBlockHeader := &types.Header{
			Number: big.NewInt(74),
			Time:   95001,
		}
		// Test for a timestamp that should return the mid block
		targetTimestamp := int64(95000)
		expectedBlockNumber := expectedBlockHeader.Number.Uint64()

		client.EXPECT().HeaderByNumber(context.Background(), nil).Return(latestBlockHeader, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(100)).Return(&types.Header{
			Number: big.NewInt(100),
			Time:   100000,
		}, nil).Times(1)
		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(49)).Return(&types.Header{
			Number: big.NewInt(49),
			Time:   50000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(74)).Return(expectedBlockHeader, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(61)).Return(&types.Header{
			Number: big.NewInt(61),
			Time:   75000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(67)).Return(&types.Header{
			Number: big.NewInt(67),
			Time:   85000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(70)).Return(&types.Header{
			Number: big.NewInt(70),
			Time:   88000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(72)).Return(&types.Header{
			Number: big.NewInt(72),
			Time:   89000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(73)).Return(&types.Header{
			Number: big.NewInt(73),
			Time:   94991,
		}, nil).Times(1)

		blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp, client, blockmapper.OwershipBlockCorrectionFunc)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if blockNumber != expectedBlockNumber {
			t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
		}
	})

	t.Run("should find blocknumber in ownership chain with timestamp of evo right side", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mockClient.NewMockEthClient(ctrl)

		// Setup your worker with the mock client
		conf := &config.Config{WaitingTime: 5 * time.Second}
		w := blockmapper.New(conf, client)

		// Mock responses
		latestBlock := uint64(200)
		latestBlockHeader := &types.Header{
			Number: big.NewInt(int64(latestBlock)),
			Time:   200000,
		}

		// own block 174 ts 180900
		targetTimestampFromEvo := int64(180000)
		expectedBlockNumberFromOwn := uint64(174)

		client.EXPECT().HeaderByNumber(context.Background(), nil).Return(latestBlockHeader, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(100)).Return(&types.Header{
			Number: big.NewInt(100),
			Time:   100000,
		}, nil).Times(1)
		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(150)).Return(&types.Header{
			Number: big.NewInt(150),
			Time:   160000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(175)).Return(&types.Header{
			Number: big.NewInt(175),
			Time:   185000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(162)).Return(&types.Header{
			Number: big.NewInt(162),
			Time:   170000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(168)).Return(&types.Header{
			Number: big.NewInt(168),
			Time:   175000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(171)).Return(&types.Header{
			Number: big.NewInt(171),
			Time:   178000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(173)).Return(&types.Header{
			Number: big.NewInt(173),
			Time:   179900,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(174)).Return(&types.Header{
			Number: big.NewInt(174),
			Time:   180900,
		}, nil).Times(1)

		blockNumber, err := w.SearchBlockByTimestamp(targetTimestampFromEvo, client, blockmapper.OwershipBlockCorrectionFunc)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if blockNumber != expectedBlockNumberFromOwn {
			t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumberFromOwn)
		}
	})

	t.Run("should find blocknumber in evo chain with timestamp of own left side", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mockClient.NewMockEthClient(ctrl)

		// Setup your worker with the mock client
		conf := &config.Config{WaitingTime: 5 * time.Second}
		w := blockmapper.New(conf, client)

		// Mock responses
		latestBlock := uint64(200)
		latestBlockHeader := &types.Header{
			Number: big.NewInt(int64(latestBlock)),
			Time:   200000,
		}

		client.EXPECT().HeaderByNumber(context.Background(), nil).Return(latestBlockHeader, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(100)).Return(&types.Header{
			Number: big.NewInt(100),
			Time:   100000,
		}, nil).Times(1)
		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(49)).Return(&types.Header{
			Number: big.NewInt(49),
			Time:   50000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(74)).Return(&types.Header{
			Number: big.NewInt(74),
			Time:   95001,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(61)).Return(&types.Header{
			Number: big.NewInt(61),
			Time:   75000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(67)).Return(&types.Header{
			Number: big.NewInt(67),
			Time:   85000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(70)).Return(&types.Header{
			Number: big.NewInt(70),
			Time:   88000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(72)).Return(&types.Header{
			Number: big.NewInt(72),
			Time:   89000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(73)).Return(&types.Header{
			Number: big.NewInt(73),
			Time:   94991,
		}, nil).Times(1)

		// Test for a timestamp that should return the mid block
		targetTimestamp := int64(95000)
		expectedBlockNumber := uint64(73)

		blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp, client, blockmapper.EvoChainBlockCorrectionFunc)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if blockNumber != expectedBlockNumber {
			t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
		}
	})

	t.Run("should find blocknumber in evo chain with timestamp of own right side", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mockClient.NewMockEthClient(ctrl)

		// Setup your worker with the mock client
		conf := &config.Config{WaitingTime: 5 * time.Second}
		w := blockmapper.New(conf, client)

		// Mock responses
		latestBlock := uint64(200)
		latestBlockHeader := &types.Header{
			Number: big.NewInt(int64(latestBlock)),
			Time:   200000,
		}

		// Test for a timestamp that should return the mid block
		// evo block 173 = ts 179900
		// ownblock 200 ts 180000
		// => evo 173 => 200 own
		targetTimestampFromOwnership := int64(180000)
		expectedBlockNumberFromEvo := uint64(173)

		client.EXPECT().HeaderByNumber(context.Background(), nil).Return(latestBlockHeader, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(100)).Return(&types.Header{
			Number: big.NewInt(100),
			Time:   100000,
		}, nil).Times(1)
		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(150)).Return(&types.Header{
			Number: big.NewInt(150),
			Time:   160000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(175)).Return(&types.Header{
			Number: big.NewInt(175),
			Time:   185000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(162)).Return(&types.Header{
			Number: big.NewInt(162),
			Time:   170000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(168)).Return(&types.Header{
			Number: big.NewInt(168),
			Time:   175000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(171)).Return(&types.Header{
			Number: big.NewInt(171),
			Time:   178000,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(173)).Return(&types.Header{
			Number: big.NewInt(173),
			Time:   179900,
		}, nil).Times(1)

		client.EXPECT().HeaderByNumber(context.Background(), big.NewInt(174)).Return(&types.Header{
			Number: big.NewInt(174),
			Time:   180900,
		}, nil).Times(1)

		blockNumber, err := w.SearchBlockByTimestamp(targetTimestampFromOwnership, client, blockmapper.EvoChainBlockCorrectionFunc)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if blockNumber != expectedBlockNumberFromEvo {
			t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumberFromEvo)
		}
	})

	// t.Run("with timestamp 1707487609", func(t *testing.T) {
	// 	ctrl := gomock.NewController(t)
	// 	defer ctrl.Finish()

	// 	client, err := ethclient.Dial("https://polygon-mainnet.gorengine.com/")
	// 	if err != nil {
	// 		t.Errorf("Unexpected error: %v", err)
	// 	}
	// 	// Setup your worker with the mock client
	// 	conf := &config.Config{WaitingTime: 5 * time.Second}
	// 	w := blockmapper.New(conf, client)

	// 	// Test for a timestamp that should return the mid block
	// 	// 1707491203 ts own -> 53316113
	// 	// 1707491202 ts evo -> 510371
	// 	targetTimestamp := int64(1707491203)
	// 	expectedBlockNumber := uint64(53316114)
	// 	startTime := time.Now()
	// 	blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp, client, blockmapper.OwershipBlockCorrectionFunc)
	// 	duration := time.Since(startTime)
	// 	fmt.Println("Duration: ", duration)
	// 	if err != nil {
	// 		t.Errorf("Unexpected error: %v", err)
	// 	}
	// 	if blockNumber != expectedBlockNumber {
	// 		t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
	// 	}
	// })
	// t.Run("with timestamp 1707487609", func(t *testing.T) {

	// 	client, err := ethclient.Dial("https://rpc.klaos.laosfoundation.io/")
	// 	if err != nil {
	// 		t.Errorf("Unexpected error: %v", err)
	// 	}
	// 	// Setup your worker with the mock client
	// 	conf := &config.Config{WaitingTime: 5 * time.Second}
	// 	w := blockmapper.New(conf, client)

	// 	// Test for a timestamp that should return the mid block
	// 	// block evo 510371 => ts 1707491202
	// 	targetTimestamp := int64(1707491202)
	// 	expectedBlockNumber := uint64(510370)
	// 	startTime := time.Now()
	// 	blockNumber, err := w.SearchBlockByTimestamp(targetTimestamp, client, blockmapper.EvoChainBlockCorrectionFunc)
	// 	duration := time.Since(startTime)
	// 	fmt.Println("Duration: ", duration)
	// 	if err != nil {
	// 		t.Errorf("Unexpected error: %v", err)
	// 	}
	// 	if blockNumber != expectedBlockNumber {
	// 		t.Errorf("got %v, expected %v", blockNumber, expectedBlockNumber)
	// 	}
	// })
}
