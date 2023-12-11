package api_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	apiMock "github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
	mockStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/mock"
	v1 "github.com/freeverseio/laos-universal-node/internal/state/v1"
	"go.uber.org/mock/gomock"
)

func TestPostRpcRequestMiddleware(t *testing.T) {
	t.Parallel() // Run tests in parallel

	// Define test cases
	tests := []struct {
		name                                  string
		body                                  string
		contentType                           string
		method                                string
		expectedStatusCode                    int
		expectedResponse                      string
		proxyHandlerCalledTimes               int
		ercUniversalMintingHandlerCalledTimes int
		storedContracts                       [][]byte
	}{
		{
			name:                                  "Good request with eth_call method",
			body:                                  `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			contentType:                           "application/json",
			method:                                "POST",
			expectedStatusCode:                    http.StatusOK,
			expectedResponse:                      "universalMintingHandler called",
			ercUniversalMintingHandlerCalledTimes: 1,
			storedContracts: [][]byte{
				[]byte("contract_0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
		},
		// {
		// 	name:                    "Good request with eth_call method but contract not in list",
		// 	body:                    `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
		// 	contentType:             "application/json",
		// 	method:                  "POST",
		// 	expectedStatusCode:      http.StatusOK,
		// 	expectedResponse:        "proxyHandler called",
		// 	proxyHandlerCalledTimes: 1,
		// 	storedContracts: [][]byte{
		// 		[]byte(""),
		// 	},
		// },
		// {
		// 	name: "Good request with eth_call method",
		// 	body: `{
		//     "jsonrpc": "2.0",
		//     "method": "eth_call",
		//     "params": [{
		//         "to": "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A",
		//         "data": "0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28"
		//     }, "latest"],
		//     "id": 1
		// }`,
		// 	contentType:                           "application/json",
		// 	method:                                "POST",
		// 	expectedStatusCode:                    http.StatusOK,
		// 	expectedResponse:                      "universalMintingHandler called",
		// 	ercUniversalMintingHandlerCalledTimes: 1,
		// 	storedContracts: [][]byte{
		// 		[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
		// 	},
		// },
		// {
		// 	name: "Good request with eth_call method supportsInterface 0x780e9d63",
		// 	body: `{
		//     "jsonrpc": "2.0",
		//     "method": "eth_call",
		//     "params": [{
		//         "to": "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A",
		//         "data": "0x01ffc9a7780e9d6300000000000000000000000000000000000000000000000000000000"
		//     }, "latest"],
		//     "id": 1
		// }`,
		// 	contentType:                           "application/json",
		// 	method:                                "POST",
		// 	expectedStatusCode:                    http.StatusOK,
		// 	expectedResponse:                      "universalMintingHandler called",
		// 	ercUniversalMintingHandlerCalledTimes: 1,
		// 	storedContracts: [][]byte{
		// 		[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
		// 	},
		// },
		// {
		// 	name: "Good request with eth_call method supportsInterface 0x80ac58cd",
		// 	body: `{
		//     "jsonrpc": "2.0",
		//     "method": "eth_call",
		//     "params": [{
		//         "to": "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A",
		//         "data": "0x01ffc9a780ac58cd00000000000000000000000000000000000000000000000000000000"
		//     }, "latest"],
		//     "id": 1
		// }`,
		// 	contentType:                           "application/json",
		// 	method:                                "POST",
		// 	expectedStatusCode:                    http.StatusOK,
		// 	expectedResponse:                      "proxyHandler called",
		// 	ercUniversalMintingHandlerCalledTimes: 0,
		// 	proxyHandlerCalledTimes:               1,
		// 	storedContracts: [][]byte{
		// 		[]byte("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
		// 	},
		// },
		// {
		// 	name:                    "Good request with eth_call method but no remote minting method",
		// 	body:                    `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x95d89b41","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
		// 	contentType:             "application/json",
		// 	method:                  "POST",
		// 	expectedStatusCode:      http.StatusOK,
		// 	expectedResponse:        "proxyHandler called",
		// 	proxyHandlerCalledTimes: 1,
		// 	storedContracts:         [][]byte{[]byte("contract_0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")},
		// },
		// {
		// 	name:                    "Good request with no erc721 method",
		// 	body:                    `{"method":"eth_getBlockByNumber","params":["latest",false],"id":1,"jsonrpc":"2.0"}`,
		// 	contentType:             "application/json",
		// 	method:                  "POST",
		// 	expectedStatusCode:      http.StatusOK,
		// 	expectedResponse:        "proxyHandler called",
		// 	proxyHandlerCalledTimes: 1,
		// 	storedContracts: [][]byte{
		// 		[]byte(""),
		// 	},
		// },
		// {
		// 	name:               "Bad request with GET method",
		// 	body:               `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
		// 	contentType:        "application/json",
		// 	method:             "GET",
		// 	expectedStatusCode: http.StatusBadRequest,
		// 	expectedResponse:   "No JSON RPC call or invalid Content-Type\n",
		// 	storedContracts: [][]byte{
		// 		[]byte(""),
		// 	},
		// },
		// {
		// 	name:               "Bad request with jsonrpc 1.0",
		// 	body:               `{"jsonrpc":"1.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
		// 	contentType:        "application/json",
		// 	method:             "POST",
		// 	expectedStatusCode: http.StatusBadRequest,
		// 	expectedResponse:   "Invalid JSON-RPC version\n",
		// 	storedContracts: [][]byte{
		// 		[]byte(""),
		// 	},
		// },
	}

	// Run tests
	for _, tt := range tests {
		tt := tt // Shadow loop variable otherwise it could be overwrittens
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel
			ctrl := gomock.NewController(t)
			storage := mockStorage.NewMockService(ctrl)

			tx := mockStorage.NewMockTx(ctrl)
			// TODO fix AnyTimes
			storage.EXPECT().NewTransaction().Return(tx).AnyTimes()
			tx.EXPECT().Discard().AnyTimes()
			handlerMock := apiMock.NewMockRPCHandler(ctrl)
			t.Cleanup(func() {
				ctrl.Finish()
			})
			req := httptest.NewRequest(tt.method, "/rpc", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			// Record responses
			w := httptest.NewRecorder()
			tx.EXPECT().Get(gomock.Any()).Return(tt.storedContracts[0], nil).AnyTimes()
			stateService := v1.NewStateService(storage)

			handlerMock.EXPECT().HandleProxyRPC(req).Do(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("proxyHandler called"))
				if err != nil {
					t.Fatal(err)
				}
			}).Times(tt.proxyHandlerCalledTimes)
			handlerMock.EXPECT().HandleUniversalMinting(req, storage).Do(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("universalMintingHandler called"))
				if err != nil {
					t.Fatal(err)
				}
			}).Times(tt.ercUniversalMintingHandlerCalledTimes)

			if tt.expectedStatusCode == http.StatusOK {
				var req api.JSONRPCRequest
				if err := json.Unmarshal([]byte(tt.body), &req); err != nil {
					http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
					slog.Error("error parsing JSON request", "err", err)
					return
				}
				handlerMock.EXPECT().SetStateService(stateService).Times(1)
			}

			// Create the middleware and serve using the test handlers
			middleware := api.PostRpcRequestMiddleware(handlerMock, stateService)
			middleware.ServeHTTP(w, req)

			// Check the status code and body
			resp := w.Result()

			defer func() {
				errClose := resp.Body.Close()
				if errClose != nil {
					t.Fatalf("got: %v, expected: no error", errClose)
				}
			}()

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("got %d, Expected status code %d", resp.StatusCode, tt.expectedStatusCode)
			}

			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(resp.Body)
			if err != nil {
				t.Errorf("got %v, expected no error", err)
			}
			if !strings.Contains(buf.String(), tt.expectedResponse) {
				t.Errorf("got %q, Expected response to contain %q", buf.String(), tt.expectedResponse)
			}
		})
	}
}
