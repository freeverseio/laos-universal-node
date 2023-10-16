package rpc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/freeverseio/laos-universal-node/internal/blockchain"
	internalRpc "github.com/freeverseio/laos-universal-node/internal/rpc"
)

type Server struct {
	RPCServer *rpc.Server
}

func NewServer(
	_ context.Context,
	ethcli blockchain.EthClient,
	contractAddr common.Address,
	chainID uint64,
) (*Server, error) {
	server := &Server{rpc.NewServer()}

	eth := internalRpc.NewEthService(ethcli, contractAddr, chainID)
	if err := server.RPCServer.RegisterName("eth", eth); err != nil {
		return nil, err
	}
	net := internalRpc.NewNetService(chainID)
	if err := server.RPCServer.RegisterName("net", net); err != nil {
		return nil, err
	}
	systemHealth := internalRpc.NewSystemHealthService()
	if err := server.RPCServer.RegisterName("health", systemHealth); err != nil {
		return nil, err
	}


	return server, nil
}

// ListenAndServe starts the RPC server to listen and serve incoming requests on the specified address.
// It also handles graceful shutdown on receiving a context cancellation signal.
func (s Server) ListenAndServe(ctx context.Context, addr string) error {
	h := s.RPCServer

	server := &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 20 * time.Second,
	}

	// defer rpcServer.Stop() // nolint:gocritic // TODO: remove or uncomment

	log.Printf("RPC server listening on %s", addr)

	go func() {
		<-ctx.Done()

		// We received an interrupt signal, shut down.
		log.Println("Received server shutdown signal. Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Fatalf("HTTP server Shutdown: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return fmt.Errorf("RPC HTTP server ListenAndServe: %v", err)
	}

	log.Println("RPC server successfully stopped.")
	return nil
}

