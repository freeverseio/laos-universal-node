package api

import (
	"net/http"
	"time"
)

// Define an interface for HTTP client operations
type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Create a wrapper for the net/http client
type HttpClientWrapper struct {
	Client *http.Client
}

func (h *HttpClientWrapper) Do(req *http.Request) (*http.Response, error) {
	return h.Client.Do(req)
}

type ApiHandlerInterface interface {
	PostRpcHandler(w http.ResponseWriter, r *http.Request)
}

type ApiHandler struct {
	RpcUrl     string
	HttpClient HttpClientInterface // Inject the HTTP client interface here
}

type ApiHandlerOption func(*ApiHandler)

func WithHttpClient(client HttpClientInterface) ApiHandlerOption {
	return func(h *ApiHandler) {
		h.HttpClient = client
	}
}

func NewApiHandler(rpcUrl string, opts ...ApiHandlerOption) *ApiHandler {
	handler := &ApiHandler{
		RpcUrl: rpcUrl,
		HttpClient: &HttpClientWrapper{
			Client: &http.Client{
				Timeout: 10 * time.Second,
			},
		}, // Default HttpClient
	}

	for _, opt := range opts {
		opt(handler)
	}

	return handler
}
