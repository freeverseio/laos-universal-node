package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func Routes(h HandlerInterface, r Router) Router {
	r.HandleFunc("/", h.PostRPCHandler).Methods("POST")
	return r
}
