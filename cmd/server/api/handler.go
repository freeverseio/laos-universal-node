package api

import (
	"net/http"
	"time"

	"github.com/freeverseio/laos-universal-node/internal/state"
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

type RPCHandler interface {
	PostRPCProxyHandler(w http.ResponseWriter, r *http.Request)
	UniversalMintingRPCHandler(w http.ResponseWriter, r *http.Request)
	SetStateService(stateService state.Service)
}

type GlobalRPCHandler struct {
	rpcUrl       string
	evoRpcUrl    string
	httpClient   HTTPClientInterface // Inject the HTTP client interface here
	stateService state.Service
}

func (h *GlobalRPCHandler) GetRpcUrl() string {
	return h.rpcUrl
}

func (h *GlobalRPCHandler) GetEvoRpcUrl() string {
	return h.evoRpcUrl
}

func (h *GlobalRPCHandler) GetHttpClient() HTTPClientInterface {
	return h.httpClient
}

func (h *GlobalRPCHandler) SetStateService(stateService state.Service) {
	h.stateService = stateService
}

type HandlerOption func(*GlobalRPCHandler)

func WithHttpClient(client HTTPClientInterface) HandlerOption {
	return func(h *GlobalRPCHandler) {
		h.httpClient = client
	}
}

func NewGlobalRPCHandler(rpcUrl, evoRpcUrl string, opts ...HandlerOption) *GlobalRPCHandler {
	handler := &GlobalRPCHandler{
		rpcUrl:    rpcUrl,
		evoRpcUrl: evoRpcUrl,
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
