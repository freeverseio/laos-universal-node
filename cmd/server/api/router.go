package api

import (
	"net/http"

	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/gorilla/mux"
)

type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func Routes(h HandlerInterface, r Router, storage scan.Storage) Router {
	router := r.(*mux.Router)
	rpcProxyHandler := http.HandlerFunc(h.PostRPCProxyHandler)
	erc721UniversalMintingHandler := http.HandlerFunc(h.UniversalMintingRPCHandler)

	// Pass both handlers to the middleware and let it decide based on the JSON-RPC method and the contract address
	router.Handle("/", PostRpcRequestMiddleware(rpcProxyHandler, erc721UniversalMintingHandler, storage)).Methods("POST")
	return router
}
