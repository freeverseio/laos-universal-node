package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	ErrUniversalMintingNotReady = "universal minting not supported yet"
	ErrorInvalidRequest         = -32600 // Invalid Request
	ErrorId                     = 1
)

func (h *Handler) UniversalMintingRPCHandler(w http.ResponseWriter, r *http.Request) {
	// Set the header to application/json for the response
	w.Header().Set("Content-Type", "application/json")

	// Define the error structure as per the JSON-RPC 2.0 specification
	errorResponse := struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{
		JSONRPC: "2.0",
		ID:      ErrorId,
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    ErrorInvalidRequest,
			Message: ErrUniversalMintingNotReady,
		},
	}

	// Set the HTTP status code
	w.WriteHeader(http.StatusBadRequest)

	// Marshal the error structure to JSON
	errorJSON, err := json.Marshal(errorResponse)
	if err != nil {
		// Send an internal server error if marshaling fails
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the JSON error message to the response writer
	_, writeErr := w.Write(errorJSON)
	if writeErr != nil {
		slog.Error("error writing response", "err", writeErr)
	}
}
