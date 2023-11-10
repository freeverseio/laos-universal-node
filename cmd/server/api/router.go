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

func Routes(h RPCHandler, r Router, storage scan.Storage) Router {
	router := r.(*mux.Router)

	// middleware decides based on the JSON-RPC method and the contract address
	router.Handle("/", PostRpcRequestMiddleware(h, storage)).Methods("POST")
	return router
}
