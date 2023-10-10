package main

import (
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

// System is a type that will be exported as an RPC service.
type System int

type Args struct{}

// SystemResponse holds the result of the Multiply method.
type SystemResponse struct {
	Up int
}

// Up increments the Up field of the SystemResponse struct by 1.
func (a *System) Up(_ *Args, reply *SystemResponse) error {
	reply.Up = 1
	return nil
}

func main() {
	s := new(System)
	err := rpc.Register(s)
	if err != nil {
		fmt.Println("Error while registering RPC service", err)
		return
	}

	// Create an HTTP handler for RPC
	handler := http.NewServeMux()
	handler.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		serverCodec := jsonrpc.NewServerCodec(&httpReadWriteCloser{r.Body, w})
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		rpcErr := rpc.ServeRequest(serverCodec)
		if rpcErr != nil {
			fmt.Println("Error while serving JSON request", rpcErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	// Create an HTTP server with timeouts
	server := &http.Server{
		Addr:         ":5001",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the HTTP server
	fmt.Println("Starting server on port 5001...")
	servErr := server.ListenAndServe()
	if servErr != nil {
		fmt.Println("Error while starting server", servErr)
		return
	}
}

type httpReadWriteCloser struct {
	in  io.Reader
	out io.Writer
}

func (h *httpReadWriteCloser) Read(p []byte) (n int, err error)  { return h.in.Read(p) }
func (h *httpReadWriteCloser) Write(p []byte) (n int, err error) { return h.out.Write(p) }
func (h *httpReadWriteCloser) Close() error                      { return nil }
