package api_test

import (
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
	"go.uber.org/mock/gomock"
)

func TestHandler(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	mockHttpClient := mock.NewMockHTTPClientInterface(ctrl)
	rpcUrl := "https://polygon-mumbai.test.com/"
	handler := api.NewGlobalRPCHandler(
		rpcUrl,
		api.WithHttpClient(mockHttpClient),
	)

	if handler == nil {
		t.Error("handler is nil")
	}
	if handler.GetHttpClient() == nil {
		t.Error("handler.HttpClient is nil")
	}
	if handler.GetRpcUrl() != rpcUrl {
		t.Fatalf("Got RPC URL %v,  expected %v", handler.GetRpcUrl(), rpcUrl)
	}
}
