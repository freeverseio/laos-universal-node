package api

import (
	"net/http"

	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	"github.com/gorilla/mux"
)

type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")              // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// If it's a preflight OPTIONS request, send an OK status and return
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Routes(h RPCHandler, r Router, stateService state.Service) Router {
	router := r.(*mux.Router)

	router.Use(CORS)

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).Methods("OPTIONS")

	router.Handle("/", PostRpcRequestMiddleware(h, stateService)).Methods("POST")
	return router
}
