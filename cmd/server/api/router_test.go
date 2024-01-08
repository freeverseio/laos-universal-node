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

func TestCORS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                       string
		method                     string
		url                        string
		status                     int
		allowOrigin                string
		allowMethods               string
		postRPCRequestHandlerCalls int
	}{
		{"SupportPost", "POST", "/", http.StatusOK, "*", "POST, OPTIONS", 1},
		{"SupportOPTIONS", "OPTIONS", "/", http.StatusOK, "*", "POST, OPTIONS", 0},
		{"SupportGet", "GET", "/", http.StatusMethodNotAllowed, "", "", 0},
	}

	for _, tc := range tests {
		tc := tc // Shadow loop variable otherwise it could be overwrittens
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRPCHandler := mock.NewMockRPCHandler(mockCtrl)
			state := stateMock.NewMockService(mockCtrl)

			mockRPCHandler.EXPECT().PostRPCRequestHandler(gomock.Any(), gomock.Any()).Times(tc.postRPCRequestHandlerCalls)
			mockRPCHandler.EXPECT().SetStateService(gomock.Any()).Times(tc.postRPCRequestHandlerCalls)
			router := mux.NewRouter()
			api.Routes(mockRPCHandler, router, state)
			ts := httptest.NewServer(router)

			client := ts.Client()
			req, err := http.NewRequest(tc.method, ts.URL+tc.url, http.NoBody)
			if err != nil {
				t.Errorf("could not create request: %v", err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("could not send request: %v", err)
			}
			ts.Close()
			if res == nil {
				t.Fatalf("response is nil, expected not nil")
			}
			if err := res.Body.Close(); err != nil {
				t.Errorf("could not close response body: %v", err)
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
		})
	}
}
