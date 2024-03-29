package evolution_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/evolution"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
	"go.uber.org/mock/gomock"
)

func TestProcessEvoBlockRangeWithBadger(t *testing.T) {
	t.Run("obtained one event, event processed and last block updated successfully with badger", func(t *testing.T) {
		ctx := context.TODO()
		_, _, client, scanner, laosRpc := createMocks(t)

		db := createBadger(t)
		badgerService := badgerStorage.NewService(db)
		stateService := v1.NewStateService(badgerService)

		lastBlockData := model.Block{
			Number:    120,
			Hash:      common.HexToHash("0x7ea18f6be7115ddbb51aa052f2780a1501847f4b3a444f1a6066982b7dbab6fc"),
			Timestamp: 150,
		}
		startingBlock := uint64(100)
		startingBlockData := model.Block{
			Number:    100,
			Hash:      common.HexToHash("0xb72b31eb84c4bbbbd62aff06a3c8c88991ac7c118c47aa6fba3609ed1baa8fd3"),
			Timestamp: 110,
		}
		contract := common.HexToAddress("0x555")
		event, _ := createEventMintedWithExternalURI(lastBlockData.Number, contract)

		laosRpc.EXPECT().LatestFinalizedBlockHash().Return(latestFinalizedBlockHash, nil).Times(1)
		laosRpc.EXPECT().BlockNumber(latestFinalizedBlockHash).Return(big.NewInt(125), nil).Times(1)

		scanner.EXPECT().
			ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlockData.Number)), nil).
			Return([]scan.Event{event}, nil)

		client.EXPECT().
			BlockByNumber(ctx, big.NewInt(int64(lastBlockData.Number))).
			Return(types.NewBlockWithHeader(&types.Header{
				Time:   lastBlockData.Timestamp,
				Number: big.NewInt(int64(lastBlockData.Number)),
			}), nil)
		client.EXPECT().
			BlockByNumber(ctx, big.NewInt(int64(startingBlockData.Number))).
			Return(types.NewBlockWithHeader(&types.Header{
				Time:   startingBlockData.Timestamp,
				Number: big.NewInt(int64(startingBlockData.Number)),
			}), nil)

		p := evolution.NewProcessor(client, stateService, scanner, laosRpc, &config.Config{})
		err := p.ProcessEvoBlockRange(ctx, startingBlock, lastBlockData.Number)
		assertError(t, nil, err)

		tx, err := stateService.NewTransaction()
		assertError(t, nil, err)
		events, err := tx.GetMintedWithExternalURIEvents(contract.Hex(), 120)
		assertError(t, nil, err)

		if len(events) != 1 {
			t.Fatalf("expected 1 event, got %d", len(events))
		}
		if events[0].BlockNumber != lastBlockData.Number {
			t.Fatalf("expected block number %d, got %d", lastBlockData.Number, events[0].BlockNumber)
		}
	})
}

func TestProcessEvoBlockRangeWithBadger100Events(t *testing.T) {
	t.Run("obtained 100 events, each event processed and last block updated successfully with badger", func(t *testing.T) {
		ctx := context.TODO()
		_, _, client, scanner, laosRpc := createMocks(t)

		db := createBadger(t)
		badgerService := badgerStorage.NewService(db)
		stateService := v1.NewStateService(badgerService)

		lastBlockData := model.Block{
			Number:    120,
			Hash:      common.HexToHash("0x7ea18f6be7115ddbb51aa052f2780a1501847f4b3a444f1a6066982b7dbab6fc"),
			Timestamp: 150,
		}
		startingBlock := uint64(100)
		contract := common.HexToAddress("0x555")

		// Create 100 minted events
		events := make([]scan.Event, 100)
		for i := 0; i < 100; i++ {
			event, _ := createEventMintedWithExternalURIWithIndex(lastBlockData.Number, contract, uint64(i))
			events[i] = event
		}

		laosRpc.EXPECT().LatestFinalizedBlockHash().Return(latestFinalizedBlockHash, nil).AnyTimes()
		laosRpc.EXPECT().BlockNumber(latestFinalizedBlockHash).Return(big.NewInt(125), nil).AnyTimes()

		scanner.EXPECT().
			ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlockData.Number)), nil).
			Return(events, nil).AnyTimes()

		client.EXPECT().
			BlockByNumber(ctx, gomock.Any()).
			Return(types.NewBlockWithHeader(&types.Header{
				Time:   lastBlockData.Timestamp,
				Number: big.NewInt(int64(lastBlockData.Number)),
			}), nil).AnyTimes()

		p := evolution.NewProcessor(client, stateService, scanner, laosRpc, &config.Config{})

		err := p.ProcessEvoBlockRange(ctx, startingBlock, lastBlockData.Number)
		assertError(t, nil, err)

		tx, err := stateService.NewTransaction()
		assertError(t, nil, err)
		e, err := tx.GetMintedWithExternalURIEvents(contract.Hex(), 120)
		assertError(t, nil, err)
		if len(e) != 100 {
			t.Fatalf("expected 100 events, got %d", len(e))
		}
	})
}

func createBadger(t *testing.T) *badger.DB {
	t.Helper()
	db, err := badger.Open(
		badger.DefaultOptions("").
			WithInMemory(true).
			WithLoggingLevel(badger.ERROR))
	if err != nil {
		t.Fatalf("error initializing storage: %v", err)
	}

	return db
}

func createEventMintedWithExternalURIWithIndex(blockNumber uint64, contract common.Address, txIndex uint64) (scan.EventMintedWithExternalURI, model.MintedWithExternalURI) {
	event := scan.EventMintedWithExternalURI{
		Slot:        big.NewInt(5),
		To:          common.HexToAddress("0x123"),
		TokenURI:    "https://www.google.com",
		TokenId:     big.NewInt(10),
		Contract:    contract,
		BlockNumber: blockNumber,
		Timestamp:   100,
		TxIndex:     txIndex,
	}

	adjustedEvent := model.MintedWithExternalURI{
		Slot:        big.NewInt(5),
		To:          common.HexToAddress("0x123"),
		TokenURI:    "https://www.google.com",
		TokenId:     big.NewInt(10),
		BlockNumber: blockNumber,
		Timestamp:   100,
		TxIndex:     txIndex,
	}
	return event, adjustedEvent
}
