package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
)

// TestUniversalMintingRPCHandler checks if the handler returns the correct error response and status code.
func TestUniversalMintingRPCHandler(t *testing.T) {
	t.Parallel() // Run tests in parallel
	t.Run("Should return the correct error response and status code", func(t *testing.T) {
		// Create a request to pass to our handler.
		req, err := http.NewRequest("POST", "/path", http.NoBody) // Use the appropriate method and path for your application
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := api.Handler{} // Assuming Handler is properly initialized for a real test
			h.UniversalMintingRPCHandler(w, r)
		})

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusBadRequest)
		}

		// Check the response body is what we expect.
		expected := struct {
			JSONRPC string `json:"jsonrpc"`
			ID      int    `json:"id"`
			Error   struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}{
			JSONRPC: "2.0",
			ID:      api.ErrorId,
			Error: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{
				Code:    api.ErrorInvalidRequest,
				Message: api.ErrUniversalMintingNotReady,
			},
		}
		expectedJSON, _ := json.Marshal(expected)

		if !bytes.Equal(rr.Body.Bytes(), expectedJSON) {
			t.Errorf("handler returned unexpected body: got %v expected %v", rr.Body.String(), string(expectedJSON))
		}
	})

}
