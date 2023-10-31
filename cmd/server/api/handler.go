package api

import (
	"net/http"
	"time"
)

// Define an interface for HTTP client operations
type HTTPClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Create a wrapper for the net/http client
type HTTPClientWrapper struct {
	client *http.Client
}

func (h *HTTPClientWrapper) Do(req *http.Request) (*http.Response, error) {
	return h.client.Do(req)
}

type ApiHandlerInterface interface {
	PostRpcHandler(w http.ResponseWriter, r *http.Request)
}

type ApiHandler struct {
	RpcUrl     string
	HttpClient HTTPClientInterface // Inject the HTTP client interface here
}

type ApiHandlerOption func(*ApiHandler)

func WithHttpClient(client HTTPClientInterface) ApiHandlerOption {
	return func(h *ApiHandler) {
		h.HttpClient = client
	}
}

func NewApiHandler(rpcUrl string, opts ...ApiHandlerOption) *ApiHandler {
	handler := &ApiHandler{
		RpcUrl: rpcUrl,
		HttpClient: &HTTPClientWrapper{
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
		}, // Default HttpClient
	}

	for _, opt := range opts {
		opt(handler)
	}

	return handler
}
