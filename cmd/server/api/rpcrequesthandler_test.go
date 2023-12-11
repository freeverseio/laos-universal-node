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
		name            string
		method          string
		contentType     string
		requestBody     string
		mockResponse    []api.RPCResponse
		expectedStatus  int
		expectedBody    string
		storedContracts [][]byte
	}{
		{
			name:           "Good request with eth_call method",
			method:         http.MethodPost,
			contentType:    "application/json",
			requestBody:    `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"jsonrpc":"2.0","id":1,"result":"0x00000000000"}\n`,
			storedContracts: [][]byte{
				[]byte("contract_0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"),
			},
		},

		// Additional test cases...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/rpc", bytes.NewBufferString(tc.requestBody))
			request.Header.Set("Content-Type", tc.contentType)
			recorder := httptest.NewRecorder()
			ctrl := gomock.NewController(t)

			mockHttpClient := mock.NewMockHTTPClientInterface(ctrl)
			handler := api.NewGlobalRPCHandler(
				"https://example.com/",
				api.WithHttpClient(mockHttpClient),
			)

			storage := mockStorage.NewMockService(ctrl)

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
			defer response.Body.Close()

			if response.StatusCode != tc.expectedStatus {
				t.Errorf("got %v, want %v", response.StatusCode, tc.expectedStatus)
			}
			if string(body) != tc.expectedBody {
				t.Errorf("got %v, want %v", string(body), tc.expectedBody)
			}
		})
	}
}
