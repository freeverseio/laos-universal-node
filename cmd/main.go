package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"golang.org/x/sync/errgroup"
)

var version = "undefined"

// System is a type that will be exported as an RPC service.
type System int

type Args struct{}

// SystemResponse holds the result of the Multiply method.
type SystemResponse struct {
	Up int
}

// nolint:unparam // for the first version this implementation is enough
func (a *System) Up(_ *Args, reply *SystemResponse) error {
	reply.Up = 1
	return nil
}

type httpReadWriteCloser struct {
	in  io.Reader
	out io.Writer
}

func (h *httpReadWriteCloser) Read(p []byte) (n int, err error)  { return h.in.Read(p) }
func (h *httpReadWriteCloser) Write(p []byte) (n int, err error) { return h.out.Write(p) }
func (h *httpReadWriteCloser) Close() error                      { return nil }

func main() {
	if err := run(); err != nil {
		slog.Error("error occurred", "err", err)
	}
}

func run() error {
	c := config.Load()

	setLogger(c.Debug)
	c.LogFields()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		client, err := ethclient.Dial(c.Rpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}
		storage, err := scan.NewFSStorage("erc721_contracts.txt")
		if err != nil {
			return fmt.Errorf("error initializing storage: %w", err)
		}
		s := scan.NewScanner(client, storage, c.Contracts...)
		return runScan(ctx, c, client, s)
	})

	// Create an HTTP handler for RPC
	handler := http.NewServeMux()
	handler.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		serverCodec := jsonrpc.NewServerCodec(&httpReadWriteCloser{r.Body, w})
		w.Header().Set("Content-type", "application/json")
		rpcErr := rpc.ServeRequest(serverCodec)
		if rpcErr != nil {
			slog.Warn("error while serving JSON request", "err", rpcErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	// Create an HTTP server with timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("starting universal node RPC server", "port", c.Port)

	group.Go(func() error {
		<-ctx.Done()
		slog.Info("shutting down the RPC server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
			return fmt.Errorf("error shutting down the RPC server: %w", shutdownErr)
		}
		return nil
	})
	group.Go(func() error {
		sys := new(System)
		err := rpc.Register(sys)
		if err != nil {
			return fmt.Errorf("error registering RPC service: %w", err)
		}
		if srvErr := server.ListenAndServe(); srvErr != nil && srvErr != http.ErrServerClosed {
			return srvErr
		}
		return nil
	})
	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func runScan(ctx context.Context, c *config.Config, client scan.EthClient, s scan.Scanner) error {
	var err error
	startingBlock := c.StartingBlock
	if startingBlock == 0 {
		startingBlock, err = getL1LatestBlock(ctx, client)
		if err != nil {
			return fmt.Errorf("error retrieving the latest block: %w", err)
		}
		slog.Debug("latest block found", "latest_block", startingBlock)
	}
	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return nil
		default:
			l1LatestBlock, err := getL1LatestBlock(ctx, client)
			if err != nil {
				slog.Error("error retrieving the latest block", "err", err.Error())
				break
			}
			lastBlock := calculateLastBlock(startingBlock, l1LatestBlock, c.BlocksRange, c.BlocksMargin)
			if lastBlock < startingBlock {
				slog.Debug("last calculated block is behind starting block, waiting...",
					"last_block", lastBlock, "starting_block", startingBlock)
				waitBeforeNextScan(ctx, c.WaitingTime)
				break
			}

			if err = s.ScanNewBridgelessMintingEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock))); err != nil {
				slog.Error("error occurred while discovering new bridgeless minting events", "err", err.Error())
				break
			}

			_, err = s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)))
			if err != nil {
				slog.Error("error occurred while scanning events", "err", err.Error())
				break
			}
			startingBlock = lastBlock + 1
		}
	}
}

func waitBeforeNextScan(ctx context.Context, waitingTime time.Duration) {
	timer := time.NewTimer(waitingTime)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}

func getL1LatestBlock(ctx context.Context, client scan.EthClient) (uint64, error) {
	l1LatestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return l1LatestBlock, nil
}

func calculateLastBlock(startingBlock, l1LatestBlock uint64, blocksRange, blocksMargin uint) uint64 {
	return min(startingBlock+uint64(blocksRange), l1LatestBlock-uint64(blocksMargin))
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
