package server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/freeverseio/laos-universal-node/cmd/server"
	"github.com/freeverseio/laos-universal-node/cmd/server/mock"

	gomock "go.uber.org/mock/gomock"
)

func TestListenAndServe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHTTPServer := mock.NewMockHTTPServerController(ctrl)
	mockHTTPServer.EXPECT().SetAddr("localhost:8080")
	mockHTTPServer.EXPECT().SetHandler(gomock.Any()).Times(1)
	mockHTTPServer.EXPECT().ListenAndServe().Return(http.ErrServerClosed)
	mockHTTPServer.EXPECT().Shutdown(gomock.Any()).Return(nil).AnyTimes()
	mockHTTPServer.EXPECT().SetKeepAlivesEnabled(false).AnyTimes()

	s, err := server.New(server.WithHTTPServer(mockHTTPServer))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err = s.ListenAndServe(ctx, "rpcUrl", "localhost:8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListenAndServeWithCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHTTPServer := mock.NewMockHTTPServerController(ctrl)
	mockHTTPServer.EXPECT().SetAddr(":9999")
	mockHTTPServer.EXPECT().SetHandler(gomock.Any()).Times(1)
	mockHTTPServer.EXPECT().ListenAndServe().Return(http.ErrServerClosed)
	mockHTTPServer.EXPECT().Shutdown(gomock.Any()).Return(nil).AnyTimes()
	mockHTTPServer.EXPECT().SetKeepAlivesEnabled(false).AnyTimes()

	s, err := server.New(server.WithHTTPServer(mockHTTPServer))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Use a channel to communicate when ListenAndServe exits
	done := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())

	// Start the server in a goroutine.
	go func() {
		err := s.ListenAndServe(ctx, "rpcUrl", ":9999") // using a random port, as it won't actually bind
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

	mockHTTPServer := mock.NewMockHTTPServerController(ctrl)
	mockHTTPServer.EXPECT().SetAddr(":9999")
	mockHTTPServer.EXPECT().SetHandler(gomock.Any()).Times(1)
	mockHTTPServer.EXPECT().ListenAndServe().Return(nil)
	mockHTTPServer.EXPECT().Shutdown(gomock.Any()).Return(nil).AnyTimes()
	mockHTTPServer.EXPECT().SetKeepAlivesEnabled(false).AnyTimes()

	s, err := server.New(server.WithHTTPServer(mockHTTPServer))
	if err != nil {
		t.Fatalf("got unexpected error: %v, expected: no error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err = s.ListenAndServe(ctx, "rpcUrl", ":9999")
	if err == nil {
		t.Fatalf("got nil, expected error")
	}
}
