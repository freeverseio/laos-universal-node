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

type HandlerInterface interface {
	PostRPCHandler(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	rpcUrl     string
	httpClient HTTPClientInterface // Inject the HTTP client interface here
}

func (h *Handler) GetRpcUrl() string {
	return h.rpcUrl
}

func (h *Handler) GetHttpClient() HTTPClientInterface {
	return h.httpClient
}

type HandlerOption func(*Handler)

func WithHttpClient(client HTTPClientInterface) HandlerOption {
	return func(h *Handler) {
		h.httpClient = client
	}
}

func NewHandler(rpcUrl string, opts ...HandlerOption) *Handler {
	handler := &Handler{
		rpcUrl: rpcUrl,
		httpClient: &HTTPClientWrapper{
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
