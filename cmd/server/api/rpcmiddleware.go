package api

import (
	"encoding/json"
	"net/http"

	"github.com/freeverseio/laos-universal-node/internal/state"
)

// JSONRPCRequest represents the expected structure of a JSON-RPC request.
type JSONRPCRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
	ID      *json.RawMessage  `json:"id,omitempty"`
}

type ParamsRPCRequest struct {
	Data  string `json:"data,omitempty"`
	To    string `json:"to,omitempty"`
	From  string `json:"from,omitempty"`
	Value string `json:"value,omitempty"`
}

func PostRpcRequestMiddleware(h RPCHandler, stateService state.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postRPCRequestHandler := http.HandlerFunc(h.PostRPCRequestHandler)
		h.SetStateService(stateService)
		postRPCRequestHandler.ServeHTTP(w, r)
	})
}
