package api

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func (h *ApiHandler) PostRpcHandler(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Read the body of the incoming request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer func() {
		errClose := r.Body.Close()
		if errClose != nil {
			slog.Error("Error closing response body", "error", errClose)
		}
	}() // Check error on Close

	// Prepare the request to the Ethereum node
	proxyReq, err := http.NewRequest(r.Method, h.RpcUrl, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Forward headers (optional)
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
	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			slog.Error("Error closing response body", "error", errClose)
		}
	}() // Check error on Close

	// Forward the response back to the original caller
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isGzipped(responseBody) {
		r2, errReader := gzip.NewReader(bytes.NewReader(responseBody)) // Removed unnecessary conversion
		if errReader != nil {
			fmt.Println("Failed to create GZIP reader:", errReader)
			return
		}
		defer func() {
			errClose := r2.Close()
			if errClose != nil {
				slog.Error("Error closing response body", "error", errClose)
			}
		}() // Check error on Close

		// Read and decompress the data
		decompressedData, errRead := io.ReadAll(r2)
		if errRead != nil {
			slog.Error("Failed to read/decompress GZIP data:", "error", errRead)
			return
		}
		slog.Debug("responseBody", "responseBody", string(decompressedData))
		_, err = w.Write(decompressedData) // Check error on Write
		if err != nil {
			slog.Error("Error writing response body", "error", err)
		}
	} else {
		slog.Debug("responseBody", "responseBody", string(responseBody))
		_, err = w.Write(responseBody) // Check error on Write
		if err != nil {
			slog.Error("Error writing response body", "error", err)
		}
	}
}

// IsGzipped checks if data is GZIP compressed.
func isGzipped(data []byte) bool {
	// GZIP magic number is 0x1f8b
	return len(data) > 1 && data[0] == 0x1f && data[1] == 0x8b
}
