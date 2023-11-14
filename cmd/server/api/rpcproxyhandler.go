package api

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

func (h *GlobalRPCHandler) PostRPCProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Read the body of the incoming request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrMsgBadRequest, http.StatusBadRequest)
		return
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
		http.Error(w, ErrMsgInternalError, http.StatusInternalServerError)
		return
	}

	// Forward headers the request
	proxyReq.Header = r.Header

	// Send the request to the Ethereum node
	resp, err := h.GetHttpClient().Do(proxyReq)
	if err != nil {
		http.Error(w, ErrMsgBadGateway, http.StatusBadGateway)
		return
	}
	// Forward headers to the response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Set(name, value)
		}
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			slog.Error("error closing response body", "err", errClose)
		}
	}() // Check error on Close

	// Forward the response back to the original caller
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, ErrMsgInternalError, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseBody) // Check error on Write
	if err != nil {
		slog.Error("error writing response body", "err", err)
	}
}
