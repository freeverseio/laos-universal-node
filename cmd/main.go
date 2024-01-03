package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"

	"github.com/freeverseio/laos-universal-node/cmd/server"
	"github.com/freeverseio/laos-universal-node/internal/config"
	evoworker "github.com/freeverseio/laos-universal-node/internal/core/worker/evolution"
	worker "github.com/freeverseio/laos-universal-node/internal/core/worker/ownership"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	v1 "github.com/freeverseio/laos-universal-node/internal/state/v1"
)

var version = "undefined"

const (
	klaosChainID = 2718
)

func main() {
	if err := run(); err != nil {
		slog.Error("error occurred", "err", err)
	}
}

func run() error {
	c := config.Load()
	setLogger(c.Debug)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	evoChainClient, err := ethclient.Dial(c.EvoRpc)
	if err != nil {
		return fmt.Errorf("error instantiating eth client: %w", err)
	}
	evoChainID, err := evoChainClient.ChainID(ctx)
	if err != nil {
		return err
	}
	err = c.SetGlobalConsensusAndParachain(evoChainID)
	if err != nil {
		return err
	}

	ownershipChainClient, err := ethclient.Dial(c.Rpc)
	if err != nil {
		return fmt.Errorf("error instantiating eth client: %w", err)
	}
	ownershipChainID, err := ownershipChainClient.ChainID(ctx)
	if err != nil {
		return err
	}
	dbPath := path.Join(c.Path, fmt.Sprintf("%s-%s", ownershipChainID.String(), evoChainID.String()))

	c.LogFields()

	// "WithMemTableSize" increases MemTableSize to 1GB (1<<30 is 1GB). This increases the transaction size to about 153MB (15% of MemTableSize)
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithLoggingLevel(badger.ERROR).WithMemTableSize(1 << 30))
	if err != nil {
		return fmt.Errorf("error initializing storage: %w", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			slog.Error("error closing db", "err", err)
		}
	}()

	// Disclaimer
	slog.Info("******************************************************************************")
	slog.Info("This is a beta version of the Laos Universal Node. It is not intended for production use. Use at your own risk.")
	slog.Info("You are now running the Universal Node Docker Image. Please be aware that this version currently does not handle blockchain reorganizations (reorgs). As a precaution, we strongly encourage operating with a heightened safety margin in your ownership chain management.")
	slog.Info("******************************************************************************")

	storageService := badgerStorage.NewService(db)
	stateService := v1.NewStateService(storageService)

	group, ctx := errgroup.WithContext(ctx)

	// Badger DB garbage collection
	group.Go(func() error {
		numIterations := 3
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				// garbage collection cleans up at most 1 file per iteration
				// https://dgraph.io/docs/badger/get-started/#garbage-collection
				for i := 0; i < numIterations; i++ {
					err := db.RunValueLogGC(0.5)
					if err != nil {
						if err != badger.ErrNoRewrite {
							slog.Error("error occurred while running badger GC", "err", err.Error())
						}
						break
					}
				}
			}
		}
	})

	// Ownership chain scanner
	group.Go(func() error {
		s := scan.NewScanner(ownershipChainClient, c.Contracts...)
		uWorker := worker.NewWorker(c, ownershipChainClient, s, stateService)
		return uWorker.Run(ctx)
	})

	// Evolution chain scanner
	group.Go(func() error {
		if evoChainID.Cmp(big.NewInt(klaosChainID)) == 0 {
			slog.Info("***********************************************************************************************")
			slog.Info("The KLAOS Parachain on Kusama is a test chain for the LAOS Parachain on Polkadot.")
			slog.Info("KLAOS is not endorsed by the LAOS Foundation nor Freeverse")
			slog.Info("for real-value transactions involving the KLAOS token https://www.laosfoundation.io/disclaimer-klaos")
			slog.Info("***********************************************************************************************")
		}

		s := scan.NewScanner(evoChainClient)
		evoWorker := evoworker.NewWorker(c, evoChainClient, s, stateService)
		return evoWorker.Run(ctx)
	})

	// Universal node RPC server
	group.Go(func() error {
		rpcServer, err := server.New()
		if err != nil {
			return fmt.Errorf("failed to create RPC server: %w", err)
		}
		addr := fmt.Sprintf("0.0.0.0:%v", c.Port)
		slog.Info("starting RPC server", "listen_address", addr)
		return rpcServer.ListenAndServe(ctx, c.Rpc, addr, stateService)
	})

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func setLogger(debug bool) {
	// Default slog.Level is Info (0)
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("version", version),
	}))
	slog.SetDefault(logger)
}
