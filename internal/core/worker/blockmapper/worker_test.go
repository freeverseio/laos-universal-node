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

func TestSearchBlockByTimestampTableTest(t *testing.T) {
	tests := []struct {
		name                string
		blockHeaders        []*types.Header // Add block headers slice
		targetTimestamp     int64
		correctionFunction  func(uint64) uint64
		expectedBlockNumber uint64
	}{
		{
			name: "target timestamp is in the first half of the block range with exact timestamp for ownership",
			blockHeaders: []*types.Header{
				{Number: big.NewInt(100), Time: 100000}, // Latest block header
				{Number: big.NewInt(50), Time: 95000},   // Mid block header
			},
			targetTimestamp:     95000,
			correctionFunction:  blockmapper.OwershipBlockCorrectionFunc,
			expectedBlockNumber: 50,
		},
		{
			name: "target timestamp is in the first half of the block range with exact timestamp for evo",
			blockHeaders: []*types.Header{
				{Number: big.NewInt(100), Time: 100000}, // Latest block header
				{Number: big.NewInt(50), Time: 95000},   // Mid block header
			},
			targetTimestamp:     95000,
			correctionFunction:  blockmapper.EvoChainBlockCorrectionFunc,
			expectedBlockNumber: 49,
		},
		{
			name: "should find blocknumber in evo chain from ownship timestamp with timestamp on the left side",
			blockHeaders: []*types.Header{
				{Number: big.NewInt(200), Time: 200000}, // Latest block header
				{Number: big.NewInt(100), Time: 100000},
				{Number: big.NewInt(49), Time: 50000},
				{Number: big.NewInt(74), Time: 95001},
				{Number: big.NewInt(61), Time: 75000},
				{Number: big.NewInt(67), Time: 85000},
				{Number: big.NewInt(70), Time: 88000},
				{Number: big.NewInt(72), Time: 89000},
				{Number: big.NewInt(73), Time: 94991},
			},
			targetTimestamp:     95000,
			correctionFunction:  blockmapper.EvoChainBlockCorrectionFunc,
			expectedBlockNumber: 73,
		},
		{
			name: "should find blocknumber in evo chain from ownship timestamp with timestamp on the right side",
			blockHeaders: []*types.Header{
				{Number: big.NewInt(200), Time: 200000}, // Latest block header
				{Number: big.NewInt(100), Time: 100000},
				{Number: big.NewInt(150), Time: 160000},
				{Number: big.NewInt(175), Time: 185000},
				{Number: big.NewInt(162), Time: 170000},
				{Number: big.NewInt(168), Time: 175000},
				{Number: big.NewInt(171), Time: 178000},
				{Number: big.NewInt(173), Time: 179900}, // expected block for ts 180000
				{Number: big.NewInt(174), Time: 180900},
			},
			targetTimestamp:     180000,
			correctionFunction:  blockmapper.EvoChainBlockCorrectionFunc,
			expectedBlockNumber: 173,
		},
		{
			name: "should find blocknumber in own chain from evo timestamp with timestamp on the left side",
			blockHeaders: []*types.Header{
				{Number: big.NewInt(200), Time: 200000}, // Latest block header
				{Number: big.NewInt(100), Time: 100000},
				{Number: big.NewInt(49), Time: 50000},
				{Number: big.NewInt(74), Time: 95001},
				{Number: big.NewInt(61), Time: 75000},
				{Number: big.NewInt(67), Time: 85000},
				{Number: big.NewInt(70), Time: 88000},
				{Number: big.NewInt(72), Time: 89000},
				{Number: big.NewInt(73), Time: 94991},
			},
			targetTimestamp:     95000,
			correctionFunction:  blockmapper.OwershipBlockCorrectionFunc,
			expectedBlockNumber: 74,
		},
		{
			name: "should find blocknumber in ownership chain with timestamp of evo right side",
			blockHeaders: []*types.Header{
				{Number: big.NewInt(200), Time: 200000}, // Latest block header
				{Number: big.NewInt(100), Time: 100000},
				{Number: big.NewInt(150), Time: 160000},
				{Number: big.NewInt(175), Time: 185000},
				{Number: big.NewInt(162), Time: 170000},
				{Number: big.NewInt(168), Time: 175000},
				{Number: big.NewInt(171), Time: 178000},
				{Number: big.NewInt(173), Time: 179900},
				{Number: big.NewInt(174), Time: 180900}, // expected block for ts 180000
			},
			targetTimestamp:     180000,
			correctionFunction:  blockmapper.OwershipBlockCorrectionFunc,
			expectedBlockNumber: 174,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := mockClient.NewMockEthClient(ctrl)

			// Setup mock responses for HeaderByNumber
			for i, header := range tt.blockHeaders {
				numberArg := big.NewInt(header.Number.Int64()) // nil for the latest header
				if i == 0 {
					numberArg = nil // The first call expects the latest block header
				}
				client.EXPECT().HeaderByNumber(context.Background(), numberArg).Return(header, nil).Times(1)
			}

			conf := &config.Config{WaitingTime: 5 * time.Second}
			w := blockmapper.New(conf, client)

			blockNumber, err := w.SearchBlockByTimestamp(tt.targetTimestamp, client, tt.correctionFunction)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if blockNumber != tt.expectedBlockNumber {
				t.Errorf("got %v, expected %v", blockNumber, tt.expectedBlockNumber)
			}
		})
	}
}
