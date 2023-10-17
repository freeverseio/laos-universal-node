package rpc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/freeverseio/laos-universal-node/internal/blockchain"
	internalRpc "github.com/freeverseio/laos-universal-node/internal/rpc"
)

type HTTPServerer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
	SetKeepAlivesEnabled(bool)
	SetAddr(string)
	SetHandler(http.Handler)
}

// The real implementation that wraps http.Server
type HTTPServer struct {
	server *http.Server
}

func (h *HTTPServer) ListenAndServe() error {
	return h.server.ListenAndServe()
}

func (h *HTTPServer) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

func (h *HTTPServer) SetKeepAlivesEnabled(v bool) {
	h.server.SetKeepAlivesEnabled(v)
}

func (h *HTTPServer) SetAddr(addr string) {
	h.server.Addr = addr
}

func (h *HTTPServer) SetHandler(handler http.Handler) {
	h.server.Handler = handler
}

type RPCServerer interface {
	RegisterName(name string, receiver interface{}) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Server struct {
	RPCServer  RPCServerer
	HTTPServer HTTPServerer
}

type ServerOption func(*Server) error

// WithEthService initializes and registers the eth service with the server.
func WithEthService(ethcli blockchain.EthClient, contractAddr common.Address, chainID uint64) ServerOption {
	return func(s *Server) error {
		eth := internalRpc.NewEthService(ethcli.Client(), contractAddr, chainID)
		return s.RPCServer.RegisterName("eth", eth)
	}
}

// WithNetService initializes and registers the net service with the server.
func WithNetService(chainID uint64) ServerOption {
	return func(s *Server) error {
		net := internalRpc.NewNetService(chainID)
		return s.RPCServer.RegisterName("net", net)
	}
}

// WithSystemHealthService initializes and registers the system health service with the server.
func WithSystemHealthService() ServerOption {
	return func(s *Server) error {
		systemHealth := internalRpc.NewSystemHealthService()
		return s.RPCServer.RegisterName("health", systemHealth)
	}
}

// WithRPCServer allows you to provide a custom RPCServerer implementation.
func WithRpcServer(rpcServer RPCServerer) ServerOption {
	return func(s *Server) error {
		s.RPCServer = rpcServer
		return nil
	}
}

// WithHTTPServer allows you to provide a custom HTTPServerer implementation.
func WithHTTPServer(httpServer HTTPServerer) ServerOption {
	return func(s *Server) error {
		s.HTTPServer = httpServer
		return nil
	}
}

func NewServer(opts ...ServerOption) (*Server, error) {
	// Default to rpc.NewServer() unless an option overwrites it.
	server := &Server{
		RPCServer: rpc.NewServer(),
		HTTPServer: &HTTPServer{
			server: &http.Server{
				ReadHeaderTimeout: 20 * time.Second,
			},
		},
	}

	for _, opt := range opts {
		if err := opt(server); err != nil {
			return nil, err
		}
	}

	return server, nil
}

// ListenAndServe starts the RPC server to listen and serve incoming requests on the specified address.
// It also handles graceful shutdown on receiving a context cancellation signal.
func (s Server) ListenAndServe(ctx context.Context, addr string) error {
	s.HTTPServer.SetAddr(addr)
	s.HTTPServer.SetHandler(s.RPCServer)

	slog.Info(
		"RPC server listening",
		slog.String("address", addr),
	)

	go func() {
		<-ctx.Done()
		// We received an interrupt signal, shut down.
		slog.Info("Received server shutdown signal. Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		s.HTTPServer.SetKeepAlivesEnabled(false)
		if err := s.HTTPServer.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			slog.Error("HTTP server Shutdown: %v", err)
		}
	}()

	if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("RPC HTTP server ListenAndServe: %v", err)
		// Error starting or closing listener:
		return fmt.Errorf("RPC HTTP server ListenAndServe: %v", err)
	}

	slog.Info("RPC server successfully stopped.")
	return nil
}
