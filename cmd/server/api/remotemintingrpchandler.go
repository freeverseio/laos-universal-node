package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type JSONRPCErrorResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// sendJSONRPCError sends a JSON RPC error response with the given message and code.
func sendJSONRPCError(w http.ResponseWriter, code int, message string) {
	errorResponse := JSONRPCErrorResponse{
		JSONRPC: "2.0",
		ID:      ErrorId,
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	errorJSON, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if _, writeErr := w.Write(errorJSON); writeErr != nil {
		slog.Error("error writing response", "err", writeErr)
	}
}

func (h *Handler) UniversalMintingRPCHandler(w http.ResponseWriter, r *http.Request) {
	sendJSONRPCError(w, ErrorCodeInvalidRequest, ErrMsgUniversalMintingNotReady)
}
