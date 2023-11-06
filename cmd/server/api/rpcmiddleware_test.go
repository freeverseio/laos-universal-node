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
	ctrl := gomock.NewController(t)

	storageMock := mock.NewMockStorage(ctrl)

	// Create a test handler that will be wrapped by the middleware
	standardHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("standardHandler called"))
	})
	erc721Handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("erc721Handler called"))
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
		// {
		// 	name:               "Bad Content-Type",
		// 	body:               `{"jsonrpc":"2.0","method":"eth_call","params":{"data":"0x...","to":"0x..."},"id":1}`,
		// 	contentType:        "text/plain",
		// 	method:             "POST",
		// 	expectedStatusCode: http.StatusOK, // Should route to standardHandler because of bad Content-Type
		// 	expectedResponse:   "standardHandler called",
		// 	handlerToBeCalled:  "standard",
		// },
		{
			name:               "Good request with eth_call method",
			body:               `{"jsonrpc":"2.0","method":"eth_call","params":{"data":"0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"},"id":1}`,
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
		// {
		// 	name:               "Bad JSON",
		// 	body:               `{"jsonrpc":"2.0","method":"eth_call","params":bad_json}`,
		// 	contentType:        "application/json",
		// 	method:             "POST",
		// 	expectedStatusCode: http.StatusBadRequest, // Should return BadRequest because of bad JSON
		// 	expectedResponse:   "Error parsing JSON request",
		// 	handlerToBeCalled:  "none",
		// },

	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/rpc", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)

			// Record responses
			w := httptest.NewRecorder()
			storageMock.EXPECT().ReadAll(context.Background()).Return(tc.storedContracts, nil).Times(1)
			// Create the middleware and serve using the test handlers
			middleware := api.PostRpcRequestMiddleware(standardHandler, erc721Handler, storageMock)
			middleware.ServeHTTP(w, req)

			// Check the status code and body
			resp := w.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, resp.StatusCode)
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			if !strings.Contains(buf.String(), tc.expectedResponse) {
				t.Errorf("Expected response to contain %q, got %q", tc.expectedResponse, buf.String())
			}
		})
	}
}
