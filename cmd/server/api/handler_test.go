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
	if handler == nil || handler.GetUniversalMintingRPCHandler() == nil {
		t.Error("handler.UniversalMintingRPCHandler is nil")
	}
	if handler.GetProxyRPCHandler() == nil {
		t.Error("handler.ProxyRPCHandler is nil")
	}
	if handler.GetProxyRPCHandler().GetHttpClient() == nil {
		t.Error("handler.HttpClient is nil")
	}
	if handler.GetProxyRPCHandler().GetRpcUrl() != rpcUrl {
		t.Fatalf("RpcUrl got %v,  expected %v", handler.GetProxyRPCHandler().GetRpcUrl(), rpcUrl)
	}
}
