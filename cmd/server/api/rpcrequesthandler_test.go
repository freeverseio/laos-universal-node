package api_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
	mockStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	v1 "github.com/freeverseio/laos-universal-node/internal/state/v1"
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
		storedContracts                            [][]byte
	}{
		{
			name:           "Good request with eth_call method",
			method:         http.MethodPost,
			contentType:    "application/json",
			requestBody:    `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			expectedStatus: http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 1,
			expectedBody: `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			storedContracts: [][]byte{
				[]byte("contract_0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
		},
		{
			name:                            "Good request with eth_call method and no contract in list",
			method:                          http.MethodPost,
			contentType:                     "application/json",
			requestBody:                     `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			mockResponseProxy:               []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			expectedStatus:                  http.StatusOK,
			expectedProxyHandlerCalledTimes: 1,
			expectedBody:                    `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			storedContracts: [][]byte{
				[]byte(""),
			},
		},
		{
			name:              "Good request with one eth_call method and one eth_getBlockByNumber method",
			method:            http.MethodPost,
			contentType:       "application/json",
			requestBody:       `[{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1},{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest",false],"id":1}]`,
			mockResponse:      []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: 2, Result: "0x00000000000"}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 1,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `[{"jsonrpc":"2.0","id":1,"result":"0x00000000000"},{"jsonrpc":"2.0","id":2,"result":"0x00000000000"}]`,
			storedContracts: [][]byte{
				[]byte("contract_0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
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
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			expectedStatus: http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 1,
			expectedBody: `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			storedContracts: [][]byte{
				[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
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
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			storedContracts: [][]byte{
				[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
		},
		{
			name:              "Good request with eth_call method but no remote minting method",
			method:            http.MethodPost,
			contentType:       "application/json",
			requestBody:       `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x95d89b41","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			storedContracts: [][]byte{
				[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
		},
		{
			name:              "Good request with eth_call method but no remote minting method",
			method:            http.MethodPost,
			contentType:       "application/json",
			requestBody:       `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x95d89b41","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			storedContracts: [][]byte{
				[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
		},
		{
			name:              "Good request with no erc721 method",
			method:            http.MethodPost,
			contentType:       "application/json",
			requestBody:       `{"method":"eth_getBlockByNumber","params":["latest",false],"id":1,"jsonrpc":"2.0"}`,
			mockResponseProxy: []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			expectedStatus:    http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            1,
			expectedBody:                               `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}`,
			storedContracts: [][]byte{
				[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
		},
		{
			name:           "Bad request with GET method",
			method:         http.MethodGet,
			contentType:    "application/json",
			requestBody:    `{"method":"eth_getBlockByNumber","params":["latest",false],"id":1,"jsonrpc":"2.0"}`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			expectedStatus: http.StatusBadRequest,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            0,
			expectedBody:                               `No JSON RPC call or invalid Content-Type`,
			storedContracts: [][]byte{
				[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
		},
		{
			name:           "Bad request with jsonrpc 1.0",
			method:         http.MethodPost,
			contentType:    "application/json",
			requestBody:    `{"method":"eth_getBlockByNumber","params":["latest",false],"id":1,"jsonrpc":"1.0"}`,
			mockResponse:   []api.RPCResponse{{Jsonrpc: "2.0", ID: 1, Result: "0x00000000000"}},
			expectedStatus: http.StatusOK,
			expectedUniversalMintingHandlerCalledTimes: 0,
			expectedProxyHandlerCalledTimes:            0,
			expectedBody:                               `{"jsonrpc":"2.0","id":0,"error":{"code":-32600,"message":"execution reverted"}}`,
			storedContracts: [][]byte{
				[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
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
			storage := mockStorage.NewMockService(ctrl)
			universalHandler := mock.NewMockRPCUniversalHandler(ctrl)
			proxyHandler := mock.NewMockRPCProxyHandler(ctrl)
			mockHttpClient := mock.NewMockHTTPClientInterface(ctrl)
			if len(tc.mockResponse) > 0 {
				universalHandler.EXPECT().HandleUniversalMinting(gomock.Any(), gomock.Any()).Return(tc.mockResponse[0]).Times(tc.expectedUniversalMintingHandlerCalledTimes)
			}
			if len(tc.mockResponseProxy) > 0 {
				proxyHandler.EXPECT().HandleProxyRPC(gomock.Any()).Return(tc.mockResponseProxy[0]).Times(tc.expectedProxyHandlerCalledTimes)
			}
			handler := api.NewGlobalRPCHandler(
				"https://example.com/",
				api.WithHttpClient(mockHttpClient),
				api.WithUniversalMintingRPCHandler(universalHandler),
				api.WithProxyRPCHandler(proxyHandler),
			)

			tx := mockStorage.NewMockTx(ctrl)
			// TODO fix AnyTimes
			storage.EXPECT().NewTransaction().Return(tx).AnyTimes()
			tx.EXPECT().Discard().AnyTimes()

			tx.EXPECT().Get(gomock.Any()).Return(tc.storedContracts[0], nil).AnyTimes()
			stateService := v1.NewStateService(storage)
			handler.SetStateService(stateService)
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
