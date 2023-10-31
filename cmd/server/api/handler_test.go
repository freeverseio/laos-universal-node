package api_test

import (
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
	"go.uber.org/mock/gomock"
)

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mock.NewMockHttpClientInterface(ctrl)
	rpcUrl := "https://polygon-mumbai.test.com/"
	handler := api.NewApiHandler(
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
		t.Fatalf("RpcUrl expected %v, got %v", rpcUrl, handler.GetRpcUrl())
	}
}
