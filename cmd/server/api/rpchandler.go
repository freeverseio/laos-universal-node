package api

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

func (h *ApiHandler) PostRpcHandler(w http.ResponseWriter, r *http.Request) {
	// Read the body of the incoming request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer closeBodyReader(r.Body)

	// Prepare the request to the BC node
	proxyReq, err := http.NewRequest(r.Method, h.RpcUrl, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Forward headers the request
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Set(name, value)
		}
	}

	// Send the request to the Ethereum node
	resp, err := h.HttpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	// Forward headers to the response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Set(name, value)
		}
	}

	defer closeBodyReader(resp.Body)

	// Forward the response back to the original caller
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Debug("responseBody", "responseBody", string(responseBody))
	_, err = w.Write(responseBody) // Check error on Write
	if err != nil {
		slog.Error("Error writing response body", "error", err)
	}
}

func closeBodyReader(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		slog.Error("Error closing response body", "error", err)
	}
}
