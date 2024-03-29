package api

import (
	"net/http"
	"time"

	"github.com/freeverseio/laos-universal-node/internal/platform/state"
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
	HandleUniversalMinting(r *http.Request, req JSONRPCRequest) RPCResponse
	PostRPCRequestHandler(w http.ResponseWriter, r *http.Request)
	SetStateService(stateService state.Service)
}

type RPCUniversalHandler interface {
	HandleUniversalMinting(r *http.Request, req JSONRPCRequest, stateService state.Service) RPCResponse
}

type ProxyHandler interface {
	HandleProxyRPC(r *http.Request, req JSONRPCRequest, stateService state.Service) RPCResponse
	GetRpcUrl() string
	GetHttpClient() HTTPClientInterface
	SetHttpClient(client HTTPClientInterface)
}

type GlobalRPCHandler struct {
	stateService               state.Service
	universalMintingRPCHandler RPCUniversalHandler
	rpcProxyHandler            ProxyHandler
}

func (h *GlobalRPCHandler) GetUniversalMintingRPCHandler() RPCUniversalHandler {
	return h.universalMintingRPCHandler
}

func (h *GlobalRPCHandler) GetRPCProxyHandler() ProxyHandler {
	return h.rpcProxyHandler
}

type UniversalMintingRPCHandler struct {
	rpcUrl           string
	rpcMethodManager RPCMethodManager
	httpClient       HTTPClientInterface // Inject the HTTP client interface here
}

type UniversalMintingRPCHandlerOption func(*UniversalMintingRPCHandler)

func NewUniversalMintingRPCHandler(ops ...UniversalMintingRPCHandlerOption) RPCUniversalHandler {
	h := &UniversalMintingRPCHandler{}
	for _, op := range ops {
		op(h)
	}
	return h
}

func WithRPCMethodManager(rpcMethodManager RPCMethodManager) UniversalMintingRPCHandlerOption {
	return func(h *UniversalMintingRPCHandler) {
		h.rpcMethodManager = rpcMethodManager
	}
}

func WithEvoHttpClient(client HTTPClientInterface) UniversalMintingRPCHandlerOption {
	return func(h *UniversalMintingRPCHandler) {
		h.httpClient = client
	}
}

type RPCProxyHandler struct {
	rpcUrl                string
	proxyRPCMethodManager RPCMethodManager
	httpClient            HTTPClientInterface // Inject the HTTP client interface here
}

func NewProxyHandler(ops ...ProxyHandlerOption) ProxyHandler {
	h := &RPCProxyHandler{}
	for _, op := range ops {
		op(h)
	}
	return h
}

type ProxyHandlerOption func(*RPCProxyHandler)

func WithProxyRPCMethodManager(proxyRPCMethodManager RPCMethodManager) ProxyHandlerOption {
	return func(h *RPCProxyHandler) {
		h.proxyRPCMethodManager = proxyRPCMethodManager
	}
}

func WithHttpClientProxyHandler(client HTTPClientInterface) ProxyHandlerOption {
	return func(h *RPCProxyHandler) {
		h.httpClient = client
	}
}

func (h *RPCProxyHandler) SetHttpClient(client HTTPClientInterface) {
	h.httpClient = client
}

func (h *RPCProxyHandler) GetRpcUrl() string {
	return h.rpcUrl
}

func (h *RPCProxyHandler) GetHttpClient() HTTPClientInterface {
	return h.httpClient
}

func (h *UniversalMintingRPCHandler) SetHttpClient(client HTTPClientInterface) {
	h.httpClient = client
}

func (h *UniversalMintingRPCHandler) GetRpcUrl() string {
	return h.rpcUrl
}

func (h *UniversalMintingRPCHandler) GetHttpClient() HTTPClientInterface {
	return h.httpClient
}

func (h *GlobalRPCHandler) SetStateService(stateService state.Service) {
	h.stateService = stateService
}

type HandlerOption func(*GlobalRPCHandler)

func WithHttpClient(client HTTPClientInterface) HandlerOption {
	return func(h *GlobalRPCHandler) {
		h.rpcProxyHandler.SetHttpClient(client)
	}
}

func WithUniversalMintingRPCHandler(handler RPCUniversalHandler) HandlerOption {
	return func(h *GlobalRPCHandler) {
		h.universalMintingRPCHandler = handler
	}
}

func WithRPCProxyHandler(handler ProxyHandler) HandlerOption {
	return func(h *GlobalRPCHandler) {
		h.rpcProxyHandler = handler
	}
}

func NewGlobalRPCHandler(rpcUrl, evoRpcUrl string, opts ...HandlerOption) *GlobalRPCHandler {
	httpClient := &HTTPClientWrapper{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	handler := &GlobalRPCHandler{
		universalMintingRPCHandler: &UniversalMintingRPCHandler{
			rpcUrl:           evoRpcUrl,
			rpcMethodManager: NewProxyRPCMethodManager(),
			httpClient:       httpClient,
		},

		rpcProxyHandler: &RPCProxyHandler{
			rpcUrl:                rpcUrl,
			httpClient:            httpClient,
			proxyRPCMethodManager: NewProxyRPCMethodManager(),
		},
	}

	for _, opt := range opts {
		opt(handler)
	}

	return handler
}
