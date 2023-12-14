package api

import (
	"net/http"

	"github.com/freeverseio/laos-universal-node/internal/state"
	"github.com/gorilla/mux"
)

type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func Routes(h RPCHandler, r Router, stateService state.Service) Router {
	router := r.(*mux.Router)
	router.Handle("/", PostRpcRequestMiddleware(h, stateService)).Methods("POST")
	return router
}
