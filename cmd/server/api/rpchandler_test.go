package api_test

import (
	"bytes"
	"compress/gzip"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mock.NewMockHttpClientInterface(ctrl)
	handler := api.NewApiHandler(
		"https://polygon-mumbai.g.alchemy.com/",
		api.WithHttpClient(mockHttpClient),
	)

	tests := []struct {
		name           string
		requestBody    string
		mockResponse   string
		responseGzip   bool
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful request",
			requestBody:    `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`,
			mockResponse:   `{"jsonrpc":"2.0","result":"1001","id":67}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"jsonrpc":"2.0","result":"1001","id":67}`,
		},
		{
			name:           "successful request with gzip response",
			requestBody:    `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`,
			mockResponse:   `{"jsonrpc":"2.0","result":"1001","id":67}`,
			responseGzip:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"jsonrpc":"2.0","result":"1001","id":67}`,
		},
		{
			name: "successful eth_call request with parrams",
			requestBody: `{
        "jsonrpc": "2.0",
        "method": "eth_call",
        "params": [{
            "to": "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
            "data": "0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000"
        }, "latest"],
        "id": 1
    }`,
			mockResponse:   `{"jsonrpc":"2.0","id":1,"result":"0x00477777730000000000"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"jsonrpc":"2.0","id":1,"result":"0x00477777730000000000"}`,
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/rpc", bytes.NewBufferString(tt.requestBody))
			recorder := httptest.NewRecorder()

			if tt.mockError != nil {
				mockHttpClient.EXPECT().Do(gomock.Any()).Return(nil, tt.mockError).Times(1)
			} else {
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
				}
				if tt.responseGzip {
					mockResponse = &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(gzipString(tt.mockResponse)),
					}
				}
				mockHttpClient.EXPECT().Do(gomock.Any()).Return(mockResponse, nil).Times(1)
			}

			handler.PostRpcHandler(recorder, request)

			response := recorder.Result()
			body, _ := io.ReadAll(response.Body)
			defer func() {
				errClose := response.Body.Close()
				if errClose != nil {
					t.Errorf("Error closing response body: %v", errClose)
				}
			}()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, response.StatusCode)
			}
			if string(body) != tt.expectedBody {
				t.Errorf("expected body %v, got %v", tt.expectedBody, string(body))
			}
		})
	}
}

func gzipString(s string) io.ReadCloser {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(s)); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return io.NopCloser(bytes.NewReader(b.Bytes()))
}
