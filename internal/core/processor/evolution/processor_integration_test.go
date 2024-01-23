package evolution_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/evolution"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
)

func TestProcessEvoBlockRangeWithBadger(t *testing.T) {
	t.Run("obtained one event, event processed and last block updated successfully with badger", func(t *testing.T) {
		t.Parallel()
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

		p := evolution.NewProcessor(client, stateService, scanner, laosRpc, &config.Config{})
		err := p.ProcessEvoBlockRange(ctx, startingBlock, lastBlockData.Number)
		assertError(t, nil, err)

		tx := stateService.NewTransaction()
		events, err := tx.GetMintedWithExternalURIEvents(contract.Hex())
		assertError(t, nil, err)
		fmt.Println(events)
		if len(events) != 1 {
			t.Fatalf("expected 1 event, got %d", len(events))
		}
		if events[0].BlockNumber != lastBlockData.Number {
			t.Fatalf("expected block number %d, got %d", lastBlockData.Number, events[0].BlockNumber)
		}
		if events[0].BlockHash != event.BlockHash {
			t.Fatalf("expected block hash %s, got %s", event.BlockHash, events[0].BlockHash)
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

func createBadgerTransaction(t *testing.T, db *badger.DB) state.Tx {
	t.Helper()
	badgerService := badgerStorage.NewService(db)
	stateService := v1.NewStateService(badgerService)
	return stateService.NewTransaction()
}
