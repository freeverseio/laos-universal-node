package search_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/core/block/search"
	mockClient "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/mock"
	"go.uber.org/mock/gomock"
)

func TestGetBlockByTimestamp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                string
		blockHeaders        []*types.Header
		targetTimestamp     uint64
		targetFunc          func(s search.Search, ctx context.Context, ts uint64) (uint64, error)
		evoClientCalls      int
		ownClientCalls      int
		expectedBlockNumber uint64
	}{
		{
			name: "target timestamp is in the first half of the block range with exact timestamp for ownership",
			blockHeaders: []*types.Header{
				{Number: big.NewInt(100), Time: 100000}, // Latest block header
				{Number: big.NewInt(50), Time: 95000},   // Mid block header
			},
			targetTimestamp: 95000,
			targetFunc: func(s search.Search, ctx context.Context, ts uint64) (uint64, error) {
				return s.GetOwnershipBlockByTimestamp(ctx, ts)
			},
			ownClientCalls:      1,
			evoClientCalls:      0,
			expectedBlockNumber: 50,
		},
		{
			name: "target timestamp is in the first half of the block range with exact timestamp for evo",
			blockHeaders: []*types.Header{
				{Number: big.NewInt(100), Time: 100000}, // Latest block header
				{Number: big.NewInt(50), Time: 95000},   // Mid block header
			},
			targetTimestamp: 95000,
			targetFunc: func(s search.Search, ctx context.Context, ts uint64) (uint64, error) {
				return s.GetEvolutionBlockByTimestamp(ctx, ts)
			},
			ownClientCalls:      0,
			evoClientCalls:      1,
			expectedBlockNumber: 49,
		},
		{
			name: "should find blocknumber in evo chain from ownership timestamp with timestamp on the left side",
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
			targetTimestamp: 95000,
			targetFunc: func(s search.Search, ctx context.Context, ts uint64) (uint64, error) {
				return s.GetEvolutionBlockByTimestamp(ctx, ts)
			},
			ownClientCalls:      0,
			evoClientCalls:      1,
			expectedBlockNumber: 73,
		},
		{
			name: "should find blocknumber in evo chain from ownership timestamp with timestamp on the right side",
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
			targetTimestamp: 180000,
			targetFunc: func(s search.Search, ctx context.Context, ts uint64) (uint64, error) {
				return s.GetEvolutionBlockByTimestamp(ctx, ts)
			},
			ownClientCalls:      0,
			evoClientCalls:      1,
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
			targetTimestamp: 95000,
			targetFunc: func(s search.Search, ctx context.Context, ts uint64) (uint64, error) {
				return s.GetOwnershipBlockByTimestamp(ctx, ts)
			},
			ownClientCalls:      1,
			evoClientCalls:      0,
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
			targetTimestamp: 180000,
			targetFunc: func(s search.Search, ctx context.Context, ts uint64) (uint64, error) {
				return s.GetOwnershipBlockByTimestamp(ctx, ts)
			},
			ownClientCalls:      1,
			evoClientCalls:      0,
			expectedBlockNumber: 174,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ownClient, evoClient := mockClient.NewMockEthClient(ctrl), mockClient.NewMockEthClient(ctrl)

			// Setup mock responses for HeaderByNumber
			for i, header := range tt.blockHeaders {
				numberArg := big.NewInt(header.Number.Int64()) // nil for the latest header
				if i == 0 {
					numberArg = nil // The first call expects the latest block header
				}
				ownClient.EXPECT().HeaderByNumber(context.Background(), numberArg).Return(header, nil).Times(tt.ownClientCalls)
				evoClient.EXPECT().HeaderByNumber(context.Background(), numberArg).Return(header, nil).Times(tt.evoClientCalls)
			}

			s := search.New(ownClient, evoClient)

			blockNumber, err := tt.targetFunc(s, context.Background(), tt.targetTimestamp)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if blockNumber != tt.expectedBlockNumber {
				t.Errorf("got %v, expected %v", blockNumber, tt.expectedBlockNumber)
			}
		})
	}
}
