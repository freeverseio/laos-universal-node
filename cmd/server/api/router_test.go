package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
	stateMock "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

// Mocking RPCHandler and state.Service
type MockRPCHandler struct {
	// Add necessary mock methods here
}

func (m *MockRPCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement mock behavior here
}

type MockStateService struct {
	// Add necessary mock methods here
}

// Implement mock state service methods here

func TestCORS(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRPCHandler := mock.NewMockRPCHandler(mockCtrl)
	state := stateMock.NewMockService(mockCtrl)

	mockRPCHandler.EXPECT().PostRPCRequestHandler(gomock.Any(), gomock.Any()).Times(1)
	mockRPCHandler.EXPECT().SetStateService(gomock.Any()).Times(1)
	router := mux.NewRouter()
	api.Routes(mockRPCHandler, router, state)

	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()

	tests := []struct {
		method       string
		url          string
		status       int
		allowOrigin  string
		allowMethods string
	}{
		{"POST", "/", http.StatusOK, "*", "POST, OPTIONS"},
		{"OPTIONS", "/", http.StatusOK, "*", "POST, OPTIONS"},
		{"GET", "/", http.StatusMethodNotAllowed, "", ""},
	}

	for _, tc := range tests {
		req, err := http.NewRequest(tc.method, ts.URL+tc.url, nil)
		if err != nil {
			t.Errorf("could not create request: %v", err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("could not send request: %v", err)
		}

		if res.StatusCode != tc.status {
			t.Errorf("unexpected status: got %v, expected %v", res.StatusCode, tc.status)
		}
		// Check CORS headers
		if origin := res.Header.Get("Access-Control-Allow-Origin"); origin != tc.allowOrigin {
			t.Errorf("unexpected Access-Control-Allow-Origin: got %v, expected %v", origin, tc.allowOrigin)
		}
		if methods := res.Header.Get("Access-Control-Allow-Methods"); methods != tc.allowMethods {
			t.Errorf("unexpected Access-Control-Allow-Methods: got %v, expected %v", methods, tc.allowMethods)
		}
	}
}
