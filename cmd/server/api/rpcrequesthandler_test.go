package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
	stateMock "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

func TestPostRPCRequestHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                                       string
		method                                     string
		contentType                                string
		requestBody                                string
		mockResponse                               []api.RPCResponse
		mockResponseProxy                          []api.RPCResponse
		expectedStatus                             int
		expectedUniversalMintingHandlerCalledTimes int
		expectedProxyHandlerCalledTimes            int
		expectedBody                               string
		txCalledTimes                              int
		hasERC721UniversalContractReturn           bool
	}{
		{
			name:           "Good request with eth_call method",
			method:         http.MethodPost,
			contentType:    "application/json",
			requestBody:    `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("33"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus: http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 1,
			expectedBody:                     `{"jsonrpc":"2.0","id":33,"result":"0x00000000000"}`,
			hasERC721UniversalContractReturn: true,
			txCalledTimes:                    1,
		},
		{
			name:           "Good request with eth_call method and no id",
			method:         http.MethodPost,
			contentType:    "application/json",
			requestBody:    `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}]}`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: nil, Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus: http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 1,
			expectedBody:                     `{"jsonrpc":"2.0","id":null,"result":"0x00000000000"}`,
			hasERC721UniversalContractReturn: true,
			txCalledTimes:                    1,
		},
		{
			name:           "Good request with eth_call method as an array",
			method:         http.MethodPost,
			contentType:    "application/json",
			requestBody:    `[{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}]`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus: http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 1,
			expectedBody:                     `[{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}]`,
			hasERC721UniversalContractReturn: true,
			txCalledTimes:                    1,
		},
		{
			name:                             "Good request with eth_call method and no contract in list",
			method:                           http.MethodPost,
			contentType:                      "application/json",
			requestBody:                      `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			mockResponseProxy:                []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus:                   http.StatusOK,
			expectedProxyHandlerCalledTimes:  1,
			expectedBody:                     `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			hasERC721UniversalContractReturn: false,
			txCalledTimes:                    1,
		},
		{
			name:              "Good request with one eth_call method and one eth_getBlockByNumber method",
			method:            http.MethodPost,
			contentType:       "application/json",
			requestBody:       `[{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1},{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest",false],"id":1}]`,
			mockResponse:      []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("2"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 1,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `[{"jsonrpc":"2.0","id":1,"result":"0x00000000000"},{"jsonrpc":"2.0","id":2,"result":"0x00000000000"}]`,
			hasERC721UniversalContractReturn:           true,
			txCalledTimes:                              1,
		},
		{
			name:        "Good request with eth_call method supportsInterface 0x780e9d63",
			method:      http.MethodPost,
			contentType: "application/json",
			requestBody: `{
		      "jsonrpc": "2.0",
		      "method": "eth_call",
		      "params": [{
		          "to": "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A",
		          "data": "0x01ffc9a7780e9d6300000000000000000000000000000000000000000000000000000000"
		      }, "latest"],
		      "id": 1
		  }`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus: http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 1,
			expectedBody:                     `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			hasERC721UniversalContractReturn: true,
			txCalledTimes:                    1,
		},
		{
			name:        "Good request with eth_call method supportsInterface 0x80ac58cd",
			method:      http.MethodPost,
			contentType: "application/json",
			requestBody: `{
		    "jsonrpc": "2.0",
		    "method": "eth_call",
		    "params": [{
		        "to": "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A",
		        "data": "0x01ffc9a780ac58cd00000000000000000000000000000000000000000000000000000000"
		    }, "latest"],
		    "id": 1
		}`,
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			hasERC721UniversalContractReturn:           true,
		},
		{
			name:              "Good request with eth_call method but no remote minting method",
			method:            http.MethodPost,
			contentType:       "application/json",
			requestBody:       `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x95d89b41","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			hasERC721UniversalContractReturn:           true,
		},
		{
			name:              "Good request with no erc721 method",
			method:            http.MethodPost,
			contentType:       "application/json",
			requestBody:       `{"method":"eth_getBlockByNumber","params":["latest",false],"id":1,"jsonrpc":"2.0"}`,
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			hasERC721UniversalContractReturn:           true,
		},
		{
			name:           "Good request with no erc721 method but eth_blockNumber call",
			method:         http.MethodPost,
			contentType:    "application/json",
			requestBody:    `{"method":"eth_blockNumber","params":[],"id":1,"jsonrpc":"2.0"}`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus: http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 1,
			expectedProxyHandlerCalledTimes:            0,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			hasERC721UniversalContractReturn:           true,
		},
		{
			name:              "Good request with no erc721 method but eth_getBlockByHash and result null",
			method:            http.MethodPost,
			contentType:       "application/json",
			requestBody:       `{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0x1",false],"id":1}`,
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: nil}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"result":null}`,
			hasERC721UniversalContractReturn:           false,
			txCalledTimes:                              0,
		},
		{
			name:           "Bad request with GET method",
			method:         http.MethodGet,
			contentType:    "application/json",
			requestBody:    `{"method":"eth_getBlockByNumber","params":["latest",false],"id":1,"jsonrpc":"2.0"}`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus: http.StatusBadRequest,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            0,
			expectedBody:                               `No JSON RPC call or invalid Content-Type`,
			hasERC721UniversalContractReturn:           true,
		},
		{
			name:           "Bad request with jsonrpc 1.0",
			method:         http.MethodPost,
			contentType:    "application/json",
			requestBody:    `{"method":"eth_getBlockByNumber","params":["latest",false],"id":1,"jsonrpc":"1.0"}`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: getJsonRawMessagePointer("1"), Result: getHexJsonRawMessagePointer("0x00000000000")}},
			expectedStatus: http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            0,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"execution reverted"}}`,
			hasERC721UniversalContractReturn:           true,
		},
	}

	for _, tc := range tests {
		tc := tc // Shadow loop variable otherwise it could be overwrittens
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request := httptest.NewRequest(tc.method, "/rpc", bytes.NewBufferString(tc.requestBody))
			request.Header.Set("Content-Type", tc.contentType)
			recorder := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			state := stateMock.NewMockService(ctrl)

			tx := stateMock.NewMockTx(ctrl)

			universalHandler := mock.NewMockRPCUniversalHandler(ctrl)
			proxyHandler := mock.NewMockProxyHandler(ctrl)
			mockHttpClient := mock.NewMockHTTPClientInterface(ctrl)

			if len(tc.mockResponse) > 0 {
				universalHandler.EXPECT().HandleUniversalMinting(gomock.Any(), gomock.Any()).Return(tc.mockResponse[0]).Times(tc.expectedUniversalMintingHandlerCalledTimes)
			}
			if len(tc.mockResponseProxy) > 0 {
				proxyHandler.EXPECT().HandleProxyRPC(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockResponseProxy[0]).Times(tc.expectedProxyHandlerCalledTimes)
			}

			handler := api.NewGlobalRPCHandler(
				"https://example.com/",
				api.WithHttpClient(mockHttpClient),
				api.WithUniversalMintingRPCHandler(universalHandler),
				api.WithRPCProxyHandler(proxyHandler),
			)

			state.EXPECT().NewTransaction().Return(tx).Times(tc.txCalledTimes)
			tx.EXPECT().Discard().Times(tc.txCalledTimes)
			tx.EXPECT().
				HasERC721UniversalContract("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A").
				Return(tc.hasERC721UniversalContractReturn, nil).
				Times(tc.txCalledTimes)

			handler.SetStateService(state)
			http.HandlerFunc(handler.PostRPCRequestHandler).ServeHTTP(recorder, request)

			response := recorder.Result()
			body, _ := io.ReadAll(response.Body)
			defer func() {
				err := response.Body.Close()
				if err != nil {
					t.Errorf("got %v, want %v", err, nil)
				}
			}()

			if response.StatusCode != tc.expectedStatus {
				t.Errorf("got %v, want %v", response.StatusCode, tc.expectedStatus)
			}
			for i := 0; i < len(body) && i < len(tc.expectedBody); i++ {
				if string(body)[i] != tc.expectedBody[i] {
					t.Errorf("got %v, want %v", string(body), tc.expectedBody)
				}
			}
		})
	}
}

func getJsonRawMessagePointer(idStr string) *json.RawMessage {
	rawMsg := json.RawMessage(idStr)
	return &rawMsg
}

func getHexJsonRawMessagePointer(idStr string) *json.RawMessage {
	quotedResult := fmt.Sprintf(`%q`, idStr)
	r := json.RawMessage(quotedResult)
	return &r
}
