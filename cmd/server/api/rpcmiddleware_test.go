package api_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/scan/mock"
	"go.uber.org/mock/gomock"
)

func TestPostRpcRequestMiddleware(t *testing.T) {
	t.Parallel() // Run tests in parallel

	// Create a test handler that will be wrapped by the middleware
	standardHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("standardHandler called"))
		if err != nil {
			t.Fatal(err)
		}
	})
	erc721Handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("erc721Handler called"))
		if err != nil {
			t.Fatal(err)
		}
	})

	// Define test cases
	tests := []struct {
		name               string
		body               string
		contentType        string
		method             string
		expectedStatusCode int
		expectedResponse   string
		handlerToBeCalled  string
		storedContracts    []scan.ERC721UniversalContract
	}{
		{
			name:               "Good request with eth_call method",
			body:               `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			contentType:        "application/json",
			method:             "POST",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "erc721Handler called",
			handlerToBeCalled:  "erc721",
			storedContracts: []scan.ERC721UniversalContract{
				{
					Address: common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
					Block:   uint64(0),
					BaseURI: "evochain1/collectionId/",
				},
			},
		},
		{
			name: "Good request with eth_call method",
			body: `{
		    "jsonrpc": "2.0",
		    "method": "eth_call",
		    "params": [{
		        "to": "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A",
		        "data": "0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28"
		    }, "latest"],
		    "id": 1
		}`,
			contentType:        "application/json",
			method:             "POST",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "erc721Handler called",
			handlerToBeCalled:  "erc721",
			storedContracts: []scan.ERC721UniversalContract{
				{
					Address: common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
					Block:   uint64(0),
					BaseURI: "evochain1/collectionId/",
				},
			},
		},
		{
			name:               "Good request with eth_call method but no remote minting method",
			body:               `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x95d89b41","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			contentType:        "application/json",
			method:             "POST",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "standardHandler called",
			handlerToBeCalled:  "standard",
			storedContracts: []scan.ERC721UniversalContract{
				{
					Address: common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
					Block:   uint64(0),
					BaseURI: "evochain1/collectionId/",
				},
			},
		},
		{
			name:               "Good request with no erc721 method",
			body:               `{"method":"eth_getBlockByNumber","params":["latest",false],"id":1,"jsonrpc":"2.0"}`,
			contentType:        "application/json",
			method:             "POST",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "standardHandler called",
			handlerToBeCalled:  "standard",
		},
		{
			name:               "Bad request with GET method",
			body:               `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			contentType:        "application/json",
			method:             "GET",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "No JSON RPC call\n",
			handlerToBeCalled:  "none",
		},
		{
			name:               "Bad request with jsonrpc 1.0",
			body:               `{"jsonrpc":"1.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			contentType:        "application/json",
			method:             "POST",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Invalid JSON-RPC version\n",
			handlerToBeCalled:  "none",
		},
	}

	// Run tests
	for _, ttest := range tests {
		tc := ttest // Shadow loop variable otherwise it could be overwrittens
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel
			ctrl := gomock.NewController(t)
			storageMock := mock.NewMockStorage(ctrl)
			t.Cleanup(func() {
				ctrl.Finish()
			})
			req := httptest.NewRequest(tc.method, "/rpc", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)

			// Record responses
			w := httptest.NewRecorder()
			storageMock.EXPECT().ReadAll(context.Background()).Return(tc.storedContracts, nil).AnyTimes()
			// Create the middleware and serve using the test handlers
			middleware := api.PostRpcRequestMiddleware(standardHandler, erc721Handler, storageMock)
			middleware.ServeHTTP(w, req)

			// Check the status code and body
			resp := w.Result()

			defer func() {
				errClose := resp.Body.Close()
				if errClose != nil {
					t.Fatalf("got: %v, expected: no error", errClose)
				}
			}()

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("got %d, Expected status code %d", resp.StatusCode, tc.expectedStatusCode)
			}

			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(resp.Body)
			if err != nil {
				t.Errorf("got %v, expected no error", err)
			}
			if !strings.Contains(buf.String(), tc.expectedResponse) {
				t.Errorf("got %q, Expected response to contain %q", buf.String(), tc.expectedResponse)
			}
		})
	}
}
