package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	stateMock "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

func TestPostRpcHandler(t *testing.T) {
	t.Parallel() // Run tests in parallel

	tests := []struct {
		name            string
		requestBody     string
		requestHeaders  map[string]string
		expectedHeaders map[string]string
		mockResponse    string
		mockError       error
		expectedStatus  int
		expectedBody    api.RPCResponse
	}{
		{
			name:           "successful request",
			requestBody:    `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`,
			mockResponse:   `{"jsonrpc":"2.0","result":"1001","id":67}`,
			expectedStatus: http.StatusOK,
			expectedBody: api.RPCResponse{
				Jsonrpc: "2.0",
				ID:      getJsonRawMessagePointer("67"),
				Result:  getHexJsonRawMessagePointer("1001"),
			},
		},
		{
			name:           "successful request with result null",
			requestBody:    `{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0x1",false],"id":67}`,
			mockResponse:   `{"jsonrpc":"2.0","result":null,"id":67}`,
			expectedStatus: http.StatusOK,
			expectedBody: api.RPCResponse{
				Jsonrpc: "2.0",
				ID:      getJsonRawMessagePointer("67"),
				Result:  nil,
			},
		},
		{
			name: "successful eth_call request with params and headers",
			requestBody: `{
		    "jsonrpc": "2.0",
		    "method": "eth_call",
		    "params": [{
		        "to": "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
		        "data": "0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000"
		    }, "latest"],
		    "id": 1
		}`,
			requestHeaders:  map[string]string{"X-Custom-Header": "custom_value"},
			expectedHeaders: map[string]string{"X-Custom-Header": "custom_value"},
			mockResponse:    `{"jsonrpc":"2.0","id":1,"result":"0x00477777730000000000"}`,
			expectedStatus:  http.StatusOK,
			expectedBody: api.RPCResponse{
				Jsonrpc: "2.0",
				ID:      getJsonRawMessagePointer("1"),
				Result:  getHexJsonRawMessagePointer("0x00477777730000000000"),
			},
		},
		{
			name: "successful eth_call request with params and Accept-Encoding headers",
			requestBody: `{
		    "jsonrpc": "2.0",
		    "method": "eth_call",
		    "params": [{
		        "to": "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
		        "data": "0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000"
		    }, "latest"],
		    "id": 1
		}`,
			requestHeaders:  map[string]string{"X-Custom-Header": "custom_value", "Accept-Encoding": "gzip, deflate, br"},
			expectedHeaders: map[string]string{"X-Custom-Header": "custom_value"},
			mockResponse:    `{"jsonrpc":"2.0","id":1,"result":"0x00477777730000000000"}`,
			expectedStatus:  http.StatusOK,
			expectedBody: api.RPCResponse{
				Jsonrpc: "2.0",
				ID:      getJsonRawMessagePointer("1"),
				Result:  getHexJsonRawMessagePointer("0x00477777730000000000"),
			},
		},
		{
			name:           "non successful request",
			requestBody:    `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`,
			mockResponse:   `{"jsonrpc":"2.0","error":{"code":-32601,"message":"The method net_version does not exist/is not available"}}`,
			expectedStatus: http.StatusOK,
			expectedBody: api.RPCResponse{
				Jsonrpc: "2.0",
				ID:      nil,
				Error: &api.RPCError{
					Code:    -32601,
					Message: "The method net_version does not exist/is not available",
				},
			},
		},
		{
			name:           "client error",
			requestBody:    `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`,
			mockError:      errors.New("client error"),
			expectedStatus: http.StatusBadGateway,
			expectedBody: api.RPCResponse{
				Jsonrpc: "2.0",
				ID:      getJsonRawMessagePointer("67"),
				Error: &api.RPCError{
					Code:    -32600,
					Message: "execution reverted",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Shadow loop variable otherwise it could be overwrittens
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHttpClient := mock.NewMockHTTPClientInterface(ctrl)
			mockMethodManager := mock.NewMockRPCMethodManager(ctrl)
			mockMethodManager.EXPECT().HasRPCMethodWithBlockNumber(gomock.Any()).Return(api.RPCMethodEthCall, false).Times(1)

			proxyHandler := api.NewProxyHandler(
				api.WithHttpClientProxyHandler(mockHttpClient),
				api.WithProxyRPCMethodManager(mockMethodManager),
			)

			state := stateMock.NewMockService(ctrl)
			tx := stateMock.NewMockTx(ctrl)
			state.EXPECT().NewTransaction().Return(tx, nil).AnyTimes()
			tx.EXPECT().Discard().AnyTimes()
			tx.EXPECT().GetLastOwnershipBlock().Return(model.Block{Number: uint64(1001)}, nil).AnyTimes()

			request := httptest.NewRequest(http.MethodPost, "/rpc", bytes.NewBufferString(tt.requestBody))
			if tt.requestHeaders != nil {
				// Setting headers in the request
				for key, value := range tt.requestHeaders {
					request.Header.Set(key, value)
				}
			}

			var jsonRPCRequest api.JSONRPCRequest
			// tt.requestBody to []byte
			body := []byte(tt.requestBody)
			if err := json.Unmarshal(body, &jsonRPCRequest); err != nil {
				t.Fatalf("error unmarshalling request: %v", err)
			}

			// Mock the HTTP client behavior
			if tt.mockError != nil {
				mockHttpClient.EXPECT().Do(gomock.Any()).Return(nil, tt.mockError)
			} else {
				mockMethodManager.EXPECT().HasRPCMethodWithHash(gomock.Any()).Return(api.RPCMethodEthCall, false).Times(1)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}
				mockHttpClient.EXPECT().Do(gomock.Any()).Do(func(arg interface{}) {
					req, ok := arg.(*http.Request)
					if !ok {
						t.Fatalf("got %T, expected *http.Request", arg)
					}
					if tt.requestHeaders != nil {
						if tt.requestHeaders["Accept-Encoding"] != "" {
							if req.Header.Get("Accept-Encoding") != "" {
								t.Fatalf("got %v, expected header %v, ", req.Header.Get("Accept-Encoding"), "")
							}
						}
						for name, values := range req.Header {
							for _, value := range values {
								// we don't want to forward the Accept-Encoding header because we don't want to receive a gzipped response
								if name != "Accept-Encoding" {
									if value != tt.expectedHeaders[name] {
										t.Fatalf("got %v, expected header %v, ", value, tt.expectedHeaders[name])
									}
								}
							}
						}
					}
				}).Return(mockResponse, nil)
			}

			apiResponse := proxyHandler.HandleProxyRPC(request, jsonRPCRequest, state)
			// compare apiResponse.ID with tt.expectedBody.ID
			compareRawMessage(t, apiResponse.ID, tt.expectedBody.ID)

			if apiResponse.Jsonrpc != tt.expectedBody.Jsonrpc {
				t.Fatalf("got %v, expected %v", apiResponse.Jsonrpc, tt.expectedBody.Jsonrpc)
			}

			compareRawMessage(t, apiResponse.Result, tt.expectedBody.Result)
			if tt.expectedBody.Error != nil && apiResponse.Error.Code != tt.expectedBody.Error.Code {
				t.Fatalf("got %v, expected %v", apiResponse.Error.Code, tt.expectedBody.Error.Code)
			}
			if tt.expectedBody.Error != nil && apiResponse.Error.Message != tt.expectedBody.Error.Message {
				t.Fatalf("got %v, expected %v", apiResponse.Error.Message, tt.expectedBody.Error.Message)
			}
		})
	}
}

func compareRawMessage(t *testing.T, raw1, raw2 *json.RawMessage) {
	t.Helper()
	// Check if both are nil or both are not nil
	if (raw1 == nil) != (raw2 == nil) {
		t.Fatalf("One of the RawMessage is nil and the other is not. Got %v, expected %v", raw1, raw2)
	}

	// Compare the values if both are not nil
	if raw1 != nil && raw2 != nil {
		if !bytes.Equal(*raw1, *raw2) {
			t.Fatalf("RawMessage values are not equal. Got %v, expected %v", *raw1, *raw2)
		}
	}
}
