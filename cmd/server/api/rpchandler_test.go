package api_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
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
		expectedBody    string
	}{
		{
			name:           "successful request",
			requestBody:    `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`,
			mockResponse:   `{"jsonrpc":"2.0","result":"1001","id":67}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"jsonrpc":"2.0","result":"1001","id":67}`,
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
			expectedBody:    `{"jsonrpc":"2.0","id":1,"result":"0x00477777730000000000"}`,
		},
		{
			name:           "non successful request",
			requestBody:    `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`,
			mockResponse:   `{"jsonrpc":"2.0","error":{"code":-32601,"message":"The method net_version does not exist/is not available"}}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"jsonrpc":"2.0","error":{"code":-32601,"message":"The method net_version does not exist/is not available"}}`,
		},
		{
			name:           "non successful request with invalid json",
			requestBody:    `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67`,
			mockResponse:   `{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"}}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"}}`,
		},
		{
			name:           "client error",
			requestBody:    `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`,
			mockError:      errors.New("client error"),
			expectedStatus: http.StatusBadGateway,
			expectedBody:   "Bad Gateway\n",
		},
	}

	for _, ttest := range tests {
		tt := ttest // Shadow loop variable otherwise it could be overwrittens
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel
			ctrl := gomock.NewController(t)
			t.Cleanup(func() {
				ctrl.Finish()
			})
			mockHttpClient := mock.NewMockHTTPClientInterface(ctrl)
			handler := api.NewHandler(
				"https://polygon-mumbai.g.alchemy.com/",
				api.WithHttpClient(mockHttpClient),
			)

			request := httptest.NewRequest(http.MethodPost, "/rpc", bytes.NewBufferString(tt.requestBody))
			if tt.requestHeaders != nil && tt.requestHeaders["X-Custom-Header"] != "" {
				// Setting headers in the request
				for key, value := range tt.requestHeaders {
					request.Header.Set(key, value)
				}
			}

			recorder := httptest.NewRecorder()

			if tt.mockError != nil {
				mockHttpClient.EXPECT().Do(gomock.Any()).Return(nil, tt.mockError).Times(1)
			} else {
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
					if tt.requestHeaders != nil && tt.requestHeaders["X-Custom-Header"] != "" {
						customHeaderValue := req.Header.Get("X-Custom-Header")
						if customHeaderValue != tt.expectedHeaders["X-Custom-Header"] {
							t.Fatalf("got %v, expected header %v, ", customHeaderValue, tt.expectedHeaders["X-Custom-Header"])
						}
					}
				}).Return(mockResponse, nil).Times(1)
			}

			handler.PostRPCHandler(recorder, request)

			response := recorder.Result()
			body, _ := io.ReadAll(response.Body)
			defer func() {
				errClose := response.Body.Close()
				if errClose != nil {
					t.Fatalf("got: %v, expected: no error", errClose)
				}
			}()

			if response.StatusCode != tt.expectedStatus {
				t.Fatalf("got %v, expected status %v", response.StatusCode, tt.expectedStatus)
			}
			if string(body) != tt.expectedBody {
				t.Fatalf("got %v, expected body %v", string(body), tt.expectedBody)
			}
		})
	}
}
