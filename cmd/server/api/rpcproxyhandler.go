package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// ProxyRPCHandler
func (h *ProxyRPCHandler) HandleProxyRPC(r *http.Request) RPCResponse {
	// Read the body of the incoming request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error reading request body: %w", err))
	}
	defer func() {
		errClose := r.Body.Close()
		if errClose != nil {
			slog.Error("error closing response body", "err", errClose)
		}
	}() // Check error on Close

	// Prepare the request to the BC node
	proxyReq, err := http.NewRequest(r.Method, h.GetRpcUrl(), io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating request: %w", err))
	}

	// Forward headers the request
	for name, values := range r.Header {
		for _, value := range values {
			// we don't want to forward the Accept-Encoding header because we don't want to receive a gzipped response
			if name != "Accept-Encoding" {
				proxyReq.Header.Set(name, value)
			}
		}
	}

	// Send the request to the Ethereum node
	resp, err := h.GetHttpClient().Do(proxyReq)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error sending request: %w", err))
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			slog.Error("error closing response body", "err", errClose)
		}
	}() // Check error on Close

	response, err := getJsonRPCResponse(resp)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting JSON RPC response: %w", err))
	}
	return *response
}

func getJsonRPCResponse(r *http.Response) (*RPCResponse, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for further handling
	var req RPCResponse
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("error parsing JSON request: %w", err)
	}
	return &req, nil
}
