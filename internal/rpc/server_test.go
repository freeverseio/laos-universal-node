package rpc

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/blockchain/mock"
	mockrpc "github.com/freeverseio/laos-universal-node/internal/rpc/mock"
	tmock "github.com/stretchr/testify/mock"
	gomock "go.uber.org/mock/gomock"
)

type mockHTTPServer struct {
	tmock.Mock
}

func (m *mockHTTPServer) ListenAndServe() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockHTTPServer) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockHTTPServer) SetKeepAlivesEnabled(v bool) {
	m.Called(v)
}

func (h *mockHTTPServer) SetAddr(addr string) {
}

func (h *mockHTTPServer) SetHandler(handler http.Handler) {
}

func TestNewServer(t *testing.T) {
	// Preparing mock data
	ctrl := gomock.NewController(t)
	mockcli := mock.NewMockEthClient(ctrl)
	mockcli.EXPECT().Client().Return(nil).AnyTimes()
	mockRPCServer := mockrpc.NewMockRPCServerer(ctrl)

	mockRPCServer.EXPECT().RegisterName("net", gomock.Any()).Return(nil).Times(2)
	mockRPCServer.EXPECT().RegisterName("eth", gomock.Any()).Return(nil).Times(2)
	mockRPCServer.EXPECT().RegisterName("health", gomock.Any()).Return(nil).Times(2)
	// Creating a table of test cases
	tests := []struct {
		name         string
		ethClient    blockchain.EthClient
		contractAddr common.Address
		chainID      uint64
		wantErr      bool
	}{
		{
			name:         "Valid server creation",
			ethClient:    mockcli,
			contractAddr: common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527"),
			chainID:      80001,
			wantErr:      false,
		},
		{
			name:         "Valid server creation",
			ethClient:    mockcli,
			contractAddr: common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527"),
			chainID:      90001,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewServer(
				WithRpcServer(mockRPCServer),
				WithEthService(tt.ethClient),
				WithNetService(tt.chainID),
				WithSystemHealthService(),
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestListenAndServe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRPCServer := mockrpc.NewMockRPCServerer(ctrl)

	// Simulate ServeHTTP behavior (though for this test, it won't be executed)
	mockRPCServer.EXPECT().ServeHTTP(gomock.Any(), gomock.Any()).AnyTimes()
	mockRPCServer.EXPECT().RegisterName("health", gomock.Any()).Return(nil).Times(1)

	// Create a server instance with the mock RPC server
	server, err := NewServer(WithRpcServer(mockRPCServer), WithSystemHealthService())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Use a channel to communicate when ListenAndServe exits
	done := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())

	// Start the server in a goroutine.
	go func() {
		err := server.ListenAndServe(ctx, ":9999") // using a random port, as it won't actually bind
		done <- err
	}()

	// Wait a moment, then cancel the context
	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("unexpected error from ListenAndServe: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ListenAndServe took too long to shut down")
	}
}

func TestListenAndServeWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRPCServer := mockrpc.NewMockRPCServerer(ctrl)
	mockHTTPServer := new(mockHTTPServer)
	mockHTTPServer.On("ListenAndServe").Return(errors.New("mocked error"))

	mockRPCServer.EXPECT().ServeHTTP(gomock.Any(), gomock.Any()).AnyTimes()

	server, err := NewServer(
		WithRpcServer(mockRPCServer),
		WithHTTPServer(mockHTTPServer),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = server.ListenAndServe(context.Background(), ":9999")
	if err == nil || err.Error() != "RPC HTTP server ListenAndServe: mocked error" {
		t.Fatalf("expected a mocked error, got: %v", err)
	}
}
