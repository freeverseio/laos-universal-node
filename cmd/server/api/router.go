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

func Routes(h HandlerInterface, r Router, st scan.Storage) Router {
	router := r.(*mux.Router)
	rpcProxyHandler := http.HandlerFunc(h.PostRPCProxyHandler)
	erc721Handler := http.HandlerFunc(h.PostRPCProxyHandler)

	// Pass both handlers to the middleware and let it decide based on the JSON-RPC method
	router.Handle("/", PostRpcRequestMiddleware(rpcProxyHandler, erc721Handler, st)).Methods("POST")
	router.Use(middleware(h))
	return router
}
