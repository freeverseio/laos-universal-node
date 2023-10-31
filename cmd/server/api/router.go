package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func Routes(h ApiHandlerInterface, r Router) Router {
	r.HandleFunc("/", h.PostRPCHandler).Methods("POST")
	return r
}
