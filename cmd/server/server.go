package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

type HTTPServerController interface {
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

type Server struct {
	httpServer HTTPServerController
}

type ServerOption func(*Server) error

// WithHTTPServer allows you to provide a custom HTTPServerer implementation.
func WithHTTPServer(httpServer HTTPServerController) ServerOption {
	return func(s *Server) error {
		s.httpServer = httpServer
		return nil
	}
}

func New(opts ...ServerOption) (*Server, error) {
	server := &Server{
		httpServer: &HTTPServer{
			server: &http.Server{
				ReadHeaderTimeout: 20 * time.Second,
				WriteTimeout:      20 * time.Second,
				ReadTimeout:       20 * time.Second,
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
func (s Server) ListenAndServe(ctx context.Context, rpcUrl, addr string) error {
	s.httpServer.SetAddr(addr)

	handler := api.NewHandler(rpcUrl)
	router := mux.NewRouter()
	s.httpServer.SetHandler(api.Routes(handler, router))
	slog.Info("server listening", "address", addr)

	go func() {
		<-ctx.Done()
		// We received an interrupt signal, shut down.
		slog.Info("received server shutdown signal. Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		s.httpServer.SetKeepAlivesEnabled(false)
		if err := s.httpServer.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			slog.Error("HTTP server Shutdown", "err", err)
		}
	}()

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return fmt.Errorf("server ListenAndServe: %w", err)
	}

	slog.Info("server successfully stopped.")
	return nil
}
