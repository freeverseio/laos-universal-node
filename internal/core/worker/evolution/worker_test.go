package worker

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain/contract"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
	mockTx "github.com/freeverseio/laos-universal-node/internal/state/mock"
	"go.uber.org/mock/gomock"
)

func TestScanEvoChainOnce(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		c                      config.Config
		l1LatestBlock          uint64
		blockNumberDB          uint64
		blockNumberTimes       int
		txCreatedTimes         int
		endRangeBlockHash      common.Hash
		endRangeBlockHeader    types.Header
		expectedFromBlock      uint64
		expectedToBlock        uint64
		expectedNewLatestBlock uint64
		errorScanEvents        error
		errorSaveBlockNumber   error
		errorGetBlockNumber    error
		errorGetL1LatestBlock  error
		expectedError          error
	}{
		{
			name: "scan evo chain OK",
			c: config.Config{
				StartingBlock:   0,
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:     250,
			txCreatedTimes:    1,
			blockNumberTimes:  1,
			blockNumberDB:     100,
			endRangeBlockHash: common.HexToHash("0x2825a9d4c85d8342c6ef38070763db379f58e3948e3b6978d10c5d415b1dd385"),
			endRangeBlockHeader: types.Header{
				ParentHash:  common.HexToHash("0x8ef4db2b2081c0516426eba21c941bfc989a6e93e39c1c34ae24a7e372d02f57"),
				UncleHash:   common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
				Coinbase:    common.HexToAddress("0x045c57b46dede60001623105d351c7941c90149e"),
				Root:        common.HexToHash("0xbf24678b95e5152e321267902070da3f8a63200b7f41a922b8d325b8864e68c7"),
				TxHash:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				ReceiptHash: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				Number:      big.NewInt(99),
			},
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: 151,
		},
		{
			name: "scan evo chain OK with no starting block in DB and 0 in config",
			c: config.Config{
				StartingBlock:   0,
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:     250,
			txCreatedTimes:    1,
			blockNumberTimes:  2,
			blockNumberDB:     0,
			endRangeBlockHash: common.HexToHash("0x9caeabf605936ee6653159c2cd7bcef6d0e7e3dbb7384c23be71f4dcd54f7716"),
			endRangeBlockHeader: types.Header{
				ParentHash:  common.HexToHash("0x8ef4db2b2081c0516426eba21c941bfc989a6e93e39c1c34ae24a7e372d02f57"),
				UncleHash:   common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
				Coinbase:    common.HexToAddress("0x045c57b46dede60001623105d351c7941c90149e"),
				Root:        common.HexToHash("0xbf24678b95e5152e321267902070da3f8a63200b7f41a922b8d325b8864e68c7"),
				TxHash:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				ReceiptHash: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				Number:      big.NewInt(249),
			},
			expectedFromBlock:      250,
			expectedToBlock:        250,
			expectedNewLatestBlock: 251,
		},
		{
			name: "scan evo chain OK with no starting block in DB and 100 in config",
			c: config.Config{
				EvoStartingBlock: 100,
				EvoBlocksMargin:  0,
				EvoBlocksRange:   50,
				WaitingTime:      1 * time.Second,
				Contracts:        []string{},
			},
			l1LatestBlock:     250,
			txCreatedTimes:    1,
			blockNumberTimes:  1,
			blockNumberDB:     0,
			endRangeBlockHash: common.HexToHash("0x2825a9d4c85d8342c6ef38070763db379f58e3948e3b6978d10c5d415b1dd385"),
			endRangeBlockHeader: types.Header{
				ParentHash:  common.HexToHash("0x8ef4db2b2081c0516426eba21c941bfc989a6e93e39c1c34ae24a7e372d02f57"),
				UncleHash:   common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
				Coinbase:    common.HexToAddress("0x045c57b46dede60001623105d351c7941c90149e"),
				Root:        common.HexToHash("0xbf24678b95e5152e321267902070da3f8a63200b7f41a922b8d325b8864e68c7"),
				TxHash:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				ReceiptHash: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				Number:      big.NewInt(99),
			},
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: 151,
		},
		{
			name: "scan evo chain with an error getting the block number from L1",
			c: config.Config{
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:          250,
			txCreatedTimes:         1,
			blockNumberTimes:       1,
			blockNumberDB:          0,
			endRangeBlockHash:      common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
			expectedNewLatestBlock: 151,
			errorGetL1LatestBlock:  errors.New("error getting block number from L1"),
			expectedError:          errors.New("error retrieving the latest block from chain: error getting block number from L1"),
		},
		{
			name: "scan evo chain with an error getting the block number from DB",
			c: config.Config{
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:          250,
			txCreatedTimes:         1,
			blockNumberTimes:       0,
			blockNumberDB:          0,
			endRangeBlockHash:      common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
			expectedNewLatestBlock: 151,
			errorGetBlockNumber:    errors.New("error getting block number from DB"),
			expectedError:          errors.New("error retrieving the current block from storage: error getting block number from DB"),
		},
		{
			name: "scan evo chain and getting an error saving the block number",
			c: config.Config{
				EvoStartingBlock: 100,
				EvoBlocksMargin:  0,
				EvoBlocksRange:   50,
				WaitingTime:      1 * time.Second,
				Contracts:        []string{},
			},
			l1LatestBlock:     250,
			txCreatedTimes:    1,
			blockNumberTimes:  1,
			blockNumberDB:     0,
			endRangeBlockHash: common.HexToHash("0x2825a9d4c85d8342c6ef38070763db379f58e3948e3b6978d10c5d415b1dd385"),
			endRangeBlockHeader: types.Header{
				ParentHash:  common.HexToHash("0x8ef4db2b2081c0516426eba21c941bfc989a6e93e39c1c34ae24a7e372d02f57"),
				UncleHash:   common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
				Coinbase:    common.HexToAddress("0x045c57b46dede60001623105d351c7941c90149e"),
				Root:        common.HexToHash("0xbf24678b95e5152e321267902070da3f8a63200b7f41a922b8d325b8864e68c7"),
				TxHash:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				ReceiptHash: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				Number:      big.NewInt(99),
			},
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: 151,
			errorSaveBlockNumber:   errors.New("error saving block number"),
			expectedError:          nil, // in this case we break the loop and don't return an error
		},
		{
			name: "scan evo chain and getting an error from scan events",
			c: config.Config{
				EvoStartingBlock: 100,
				EvoBlocksMargin:  0,
				EvoBlocksRange:   50,
				WaitingTime:      1 * time.Second,
				Contracts:        []string{},
			},
			l1LatestBlock:     250,
			txCreatedTimes:    1,
			blockNumberTimes:  1,
			blockNumberDB:     100,
			endRangeBlockHash: common.HexToHash("0x2825a9d4c85d8342c6ef38070763db379f58e3948e3b6978d10c5d415b1dd385"),
			endRangeBlockHeader: types.Header{
				ParentHash:  common.HexToHash("0x8ef4db2b2081c0516426eba21c941bfc989a6e93e39c1c34ae24a7e372d02f57"),
				UncleHash:   common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
				Coinbase:    common.HexToAddress("0x045c57b46dede60001623105d351c7941c90149e"),
				Root:        common.HexToHash("0xbf24678b95e5152e321267902070da3f8a63200b7f41a922b8d325b8864e68c7"),
				TxHash:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				ReceiptHash: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				Number:      big.NewInt(99),
			},
			expectedFromBlock: 100,
			expectedToBlock:   150,
			errorScanEvents:   errors.New("error scanning events"),
			expectedError:     nil, // in this case we break the loop and don't return an error
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()

			client, scanner := getMocks(t)
			state, tx2 := getMocksFromState(t)

            evoWorker := NewWorker(&tt.c, client, scanner, state)

			state.EXPECT().NewTransaction().Return(tx2).Times(tt.txCreatedTimes)
			tx2.EXPECT().Discard().Times(tt.txCreatedTimes)
			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, tt.errorGetL1LatestBlock).
				Times(tt.blockNumberTimes)

			tx2.EXPECT().GetCurrentEvoBlock().
				Return(tt.blockNumberDB, tt.errorGetBlockNumber).
				Times(1)

			if tt.errorGetL1LatestBlock == nil && tt.errorGetBlockNumber == nil {
				state.EXPECT().NewTransaction().Return(tx2).Times(1)
				tx2.EXPECT().Discard().Times(1)
				tx2.EXPECT().GetEvoEndRangeBlockHash().Return(tt.endRangeBlockHash, nil).Times(1)

				block := types.NewBlockWithHeader(&tt.endRangeBlockHeader)
				if tt.endRangeBlockHash != (common.Hash{}) {
					client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.expectedFromBlock-1))).
						Return(block, nil).
						Times(1)
				}

				client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.expectedToBlock))).
					Return(block, nil).
					Times(1)

				tx2.EXPECT().SetEvoEndRangeBlockHash(tt.endRangeBlockHash).
					Return(nil).
					Times(1)

				scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedFromBlock)), big.NewInt(int64(tt.expectedToBlock)), nil).
					Return(nil, tt.errorScanEvents).
					Do(func(_ context.Context, _ *big.Int, _ *big.Int, _ []string) {
						if tt.errorScanEvents != nil {
							cancel() // we cancel the loop since we only want one iteration
						}
					},
					).Times(1)

				if tt.errorScanEvents == nil {
					tx2.EXPECT().SetCurrentEvoBlock(tt.expectedNewLatestBlock).
						Return(tt.errorSaveBlockNumber).Do(
						func(_ uint64) {
							cancel()
						},
					).Times(1)
					if tt.errorSaveBlockNumber == nil {
						client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.expectedToBlock))).Return(&types.Header{Time: 1000}, nil).Times(1)
						tx2.EXPECT().SetCurrentEvoBlockTimestamp(uint64(1000)).Return(nil).Times(1)
						tx2.EXPECT().Commit().Return(nil).Times(1)
					}
				}
			}

            err := evoWorker.Run(ctx)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Fatalf(`got error "%v", expected error: "%v"`, err, tt.expectedError)
			}
		})
	}
}

func TestScanEvoChainWithEvents(t *testing.T) {
	t.Parallel()
	mintedWithExternalURIEventHash := "0xa7135052b348b0b4e9943bae82d8ef1c5ac225e594ef4271d12f0744cfc98348"
	eventLog := &types.Log{
		Topics: []common.Hash{
			common.HexToHash(mintedWithExternalURIEventHash),
			common.HexToHash("0x000000000000000000000000c112bde959080c5b46e73749e3e170f47123e85a"),
		},
		Data:        common.Hex2Bytes("00000000000000000000000000000000000000003d5b1313de887a00000000003d5b1313de887a0000000000c112bde959080c5b46e73749e3e170f47123e85a0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000002e516d4e5247426d7272724862754b4558375354544d326f68325077324d757438674863537048706a367a7a637375000000000000000000000000000000000000"),
		BlockNumber: 100,
	}
	evoAbi, err := abi.JSON(strings.NewReader(contract.EvolutionMetaData.ABI))
	if err != nil {
		t.Fatalf(`got error "%v", expected error: "%v"`, err, nil)
	}
	event, err := parseMintedWithExternalURI(eventLog, &evoAbi)
	if err != nil {
		t.Fatalf(`got error "%v", expected error: "%v"`, err, nil)
	}
	tests := []struct {
		name                   string
		c                      config.Config
		l1LatestBlock          uint64
		blockNumberDB          uint64
		blockNumberTimes       int
		scanEventsTimes        int
		endRangeBlockHash      common.Hash
		endRangeBlockHeader    types.Header
		expectedFromBlock      uint64
		expectedToBlock        uint64
		expectedNewLatestBlock uint64
		errorScanEvents        error
		errorSaveBlockNumber   error
		errorGetBlockNumber    error
		errorGetL1LatestBlock  error
		expectedError          error
	}{
		{
			name: "scan evo chain OK",
			c: config.Config{
				StartingBlock:   0,
				EvoBlocksMargin: 0,
				EvoBlocksRange:  50,
				WaitingTime:     1 * time.Second,
				Contracts:       []string{},
			},
			l1LatestBlock:     250,
			blockNumberTimes:  1,
			blockNumberDB:     100,
			endRangeBlockHash: common.HexToHash("0x2825a9d4c85d8342c6ef38070763db379f58e3948e3b6978d10c5d415b1dd385"),
			endRangeBlockHeader: types.Header{
				ParentHash:  common.HexToHash("0x8ef4db2b2081c0516426eba21c941bfc989a6e93e39c1c34ae24a7e372d02f57"),
				UncleHash:   common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
				Coinbase:    common.HexToAddress("0x045c57b46dede60001623105d351c7941c90149e"),
				Root:        common.HexToHash("0xbf24678b95e5152e321267902070da3f8a63200b7f41a922b8d325b8864e68c7"),
				TxHash:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				ReceiptHash: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
				Number:      big.NewInt(99),
			},
			expectedFromBlock:      100,
			expectedToBlock:        150,
			expectedNewLatestBlock: 151,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := getContext()
			defer cancel()
			client, scanner := getMocks(t)
			storage2, tx := getMocksFromState(t)

            evoWorker := NewWorker(&tt.c, client, scanner, storage2)

			client.EXPECT().BlockNumber(ctx).
				Return(tt.l1LatestBlock, tt.errorGetL1LatestBlock).
				Times(tt.blockNumberTimes)

			tx.EXPECT().GetEvoEndRangeBlockHash().Return(tt.endRangeBlockHash, nil).Times(1)

			block := types.NewBlockWithHeader(&tt.endRangeBlockHeader)
			if tt.endRangeBlockHash != (common.Hash{}) {
				client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.expectedFromBlock-1))).
					Return(block, nil).
					Times(1)
			}

			client.EXPECT().BlockByNumber(ctx, big.NewInt(int64(tt.expectedToBlock))).
				Return(block, nil).
				Times(1)

			tx.EXPECT().SetEvoEndRangeBlockHash(tt.endRangeBlockHash).
				Return(nil).
				Times(1)

			scanner.EXPECT().ScanEvents(ctx, big.NewInt(int64(tt.expectedFromBlock)), big.NewInt(int64(tt.expectedToBlock)), nil).
				Return([]scan.Event{event}, tt.errorScanEvents).Times(1)

			tx.EXPECT().GetCurrentEvoBlock().
				Return(tt.blockNumberDB, tt.errorGetBlockNumber).
				Times(1)

			tx.EXPECT().GetMintedWithExternalURIEvents(gomock.Any()).Return(nil, nil).Times(1)
			tx.EXPECT().StoreMintedWithExternalURIEvents("0x0000000000000000000000000000000000000000", gomock.Any()).Return(nil).Times(1)

			storage2.EXPECT().NewTransaction().Return(tx).Times(2)
			tx.EXPECT().Discard().Return().Times(2)
			tx.EXPECT().SetCurrentEvoBlock(tt.expectedNewLatestBlock).Return(tt.errorSaveBlockNumber)
			client.EXPECT().HeaderByNumber(ctx, big.NewInt(int64(tt.expectedToBlock))).Return(&types.Header{Time: 1000}, nil).Times(1)
			tx.EXPECT().SetCurrentEvoBlockTimestamp(uint64(1000)).Return(nil).Times(1)
			tx.EXPECT().Commit().Return(nil).Do(
				func() {
					cancel() // we cancel the loop since we only want one iteration
				},
			).Times(1)

			err := evoWorker.Run(ctx)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Fatalf(`got error "%v", expected error: "%v"`, err, tt.expectedError)
			}
		})
	}
}

func TestStoreMintedWithExternalURIEventsByContract(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		scannedEvents map[common.Address][]scan.Event
		storedEvents  map[common.Address][]model.MintedWithExternalURI
		expectedError error
	}{
		{
			name:          "does not stored when no new events",
			scannedEvents: map[common.Address][]scan.Event{},
			expectedError: nil,
		},
		{
			name: "Store new event to non stored contract",
			scannedEvents: map[common.Address][]scan.Event{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(2),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Store  different contract events to existing contracts",
			scannedEvents: map[common.Address][]scan.Event{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(2),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(3),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
				},
				common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(43),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"),
					},
				},
			},
			storedEvents: map[common.Address][]model.MintedWithExternalURI{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					model.MintedWithExternalURI{
						Slot: big.NewInt(1),
						To:   common.HexToAddress("0xc52293dCf9bE02C6EF7e6F840BF0F0E2FC45c646"),
					},
				},
				common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"): {
					model.MintedWithExternalURI{
						Slot: big.NewInt(40),
						To:   common.HexToAddress("0xc52293dCf9bE02C6EF7e6F840BF0F0E2FC45c646"),
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Store  different contract events to existing and non existing contracts",
			scannedEvents: map[common.Address][]scan.Event{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(2),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(3),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"),
					},
				},
				common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"): {
					scan.EventMintedWithExternalURI{
						Slot:     big.NewInt(43),
						To:       common.HexToAddress("0x789abcdCf9bE02C6EF7e6F840BF0F0E2FC45c123"),
						Contract: common.HexToAddress("0xeabcd86cd26373efbF15d0ec69F39c9E77Beff91"),
					},
				},
			},
			storedEvents: map[common.Address][]model.MintedWithExternalURI{
				common.HexToAddress("0xeB28886cd26373efbF15d0ec69F39c9E77Be54A5"): {
					model.MintedWithExternalURI{
						Slot: big.NewInt(1),
						To:   common.HexToAddress("0xc52293dCf9bE02C6EF7e6F840BF0F0E2FC45c646"),
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, tx := getMocksFromState(t)

			events := make([]scan.Event, 0)
			eventsToStore := make(map[common.Address][]model.MintedWithExternalURI)
			for contract, ev := range tt.scannedEvents {
				events = append(events, ev...)
				scanned := mintEventsToModel(ev)

				tx.EXPECT().GetMintedWithExternalURIEvents(contract.String()).Return(tt.storedEvents[contract], nil)

				eventsToStore[contract] = append(eventsToStore[contract], tt.storedEvents[contract]...)
				eventsToStore[contract] = append(eventsToStore[contract], scanned...)
				tx.EXPECT().StoreMintedWithExternalURIEvents(contract.String(), eventsToStore[contract]).Return(nil)
			}
			if err := storeMintedWithExternalURIEventsByContract(tx, events); err != nil {
				t.Fatalf(`got error "%v", expected error: "%v"`, err, tt.expectedError)
			}
		})
	}
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), 100*time.Millisecond)
}

func getMocks(t *testing.T) (*mock.MockEthClient, *mock.MockScanner) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockEthClient(ctrl), mock.NewMockScanner(ctrl)
}

func getMocksFromState(t *testing.T) (*mockTx.MockService, *mockTx.MockTx) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mockTx.NewMockService(ctrl), mockTx.NewMockTx(ctrl)
}

func parseMintedWithExternalURI(eL *types.Log, contractAbi *abi.ABI) (scan.EventMintedWithExternalURI, error) {
	var mintWithExternalURI scan.EventMintedWithExternalURI
	eventMintedWithExternalURI := "MintedWithExternalURI"
	err := unpackIntoInterface(&mintWithExternalURI, eventMintedWithExternalURI, contractAbi, eL)
	if err != nil {
		return mintWithExternalURI, err
	}
	mintWithExternalURI.To = common.HexToAddress(eL.Topics[1].Hex())

	return mintWithExternalURI, nil
}

func mintEventsToModel(events []scan.Event) []model.MintedWithExternalURI {
	modelMintEvents := make([]model.MintedWithExternalURI, 0)
	for _, e := range events {
		scanMintEvent := e.(scan.EventMintedWithExternalURI)
		modelMintEvents = append(modelMintEvents, model.MintedWithExternalURI{
			Slot:        scanMintEvent.Slot,
			To:          scanMintEvent.To,
			TokenId:     scanMintEvent.TokenId,
			TokenURI:    scanMintEvent.TokenURI,
			BlockNumber: scanMintEvent.BlockNumber,
			Timestamp:   scanMintEvent.Timestamp,
		})
	}

	return modelMintEvents
}

func unpackIntoInterface(e scan.Event, eventName string, contractAbi *abi.ABI, eL *types.Log) error {
	err := contractAbi.UnpackIntoInterface(e, eventName, eL.Data)
	if err != nil {
		return fmt.Errorf("error unpacking the event %s: %w", eventName, err)
	}

	return nil
}
