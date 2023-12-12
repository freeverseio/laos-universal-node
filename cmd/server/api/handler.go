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
	HandleProxyRPC(r *http.Request, req JSONRPCRequest) RPCResponse
	HandleUniversalMinting(req JSONRPCRequest, stateService state.Service) RPCResponse
	PostRPCRequestHandler(w http.ResponseWriter, r *http.Request)
	SetStateService(stateService state.Service)
}

type RPCUniversalHandler interface {
	HandleUniversalMinting(req JSONRPCRequest, stateService state.Service) RPCResponse
}
type RPCProxyHandler interface {
	HandleProxyRPC(r *http.Request, req JSONRPCRequest) RPCResponse
	GetRpcUrl() string
	GetHttpClient() HTTPClientInterface
	SetHttpClient(client HTTPClientInterface)
}
type GlobalRPCHandler struct {
	rpcUrl                     string
	httpClient                 HTTPClientInterface // Inject the HTTP client interface here
	stateService               state.Service
	UniversalMintingRPCHandler RPCUniversalHandler
	ProxyRPCHandler            RPCProxyHandler
}

func (h *GlobalRPCHandler) GetUniversalMintingRPCHandler() RPCUniversalHandler {
	return h.UniversalMintingRPCHandler
}

func (h *GlobalRPCHandler) GetProxyRPCHandler() RPCProxyHandler {
	return h.ProxyRPCHandler
}

type UniversalMintingRPCHandler struct{}

type ProxyRPCHandler struct {
	rpcUrl     string
	httpClient HTTPClientInterface // Inject the HTTP client interface here
}

func (h *ProxyRPCHandler) SetHttpClient(client HTTPClientInterface) {
	h.httpClient = client
}

func (h *ProxyRPCHandler) GetRpcUrl() string {
	return h.rpcUrl
}

func (h *ProxyRPCHandler) GetHttpClient() HTTPClientInterface {
	return h.httpClient
}

func (h *GlobalRPCHandler) SetStateService(stateService state.Service) {
	h.stateService = stateService
}

type HandlerOption func(*GlobalRPCHandler)

func WithHttpClient(client HTTPClientInterface) HandlerOption {
	return func(h *GlobalRPCHandler) {
		h.httpClient = client
		h.ProxyRPCHandler.SetHttpClient(client)
	}
}

func WithUniversalMintingRPCHandler(handler RPCUniversalHandler) HandlerOption {
	return func(h *GlobalRPCHandler) {
		h.UniversalMintingRPCHandler = handler
	}
}

func WithProxyRPCHandler(handler RPCProxyHandler) HandlerOption {
	return func(h *GlobalRPCHandler) {
		h.ProxyRPCHandler = handler
	}
}

func NewGlobalRPCHandler(rpcUrl string, opts ...HandlerOption) *GlobalRPCHandler {
	httpClient := &HTTPClientWrapper{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	handler := &GlobalRPCHandler{
		rpcUrl:                     rpcUrl,
		httpClient:                 httpClient,
		UniversalMintingRPCHandler: &UniversalMintingRPCHandler{},
		ProxyRPCHandler: &ProxyRPCHandler{
			rpcUrl:     rpcUrl,
			httpClient: httpClient,
		},
	}

	for _, opt := range opts {
		opt(handler)
	}

	return handler
}
