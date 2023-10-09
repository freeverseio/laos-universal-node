package main

import (
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// System is a type that will be exported as an RPC service.
type System int

type Args struct {
}

// SystemResponse holds the result of the Multiply method.
type SystemResponse struct {
	Up int
}

// System.Up is a method that will be exposed as an RPC method.
func (a *System) Up(args *Args, reply *SystemResponse) error {
	reply.Up = 1
	return nil
}

func main() {
	s := new(System)
	rpc.Register(s)

	// Create an HTTP handler for RPC
	handler := http.NewServeMux()
	handler.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		serverCodec := jsonrpc.NewServerCodec(&httpReadWriteCloser{r.Body, w})
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := rpc.ServeRequest(serverCodec)
		if err != nil {
			fmt.Println("Error while serving JSON request", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	// Start the HTTP server
	fmt.Println("Starting server on port 5001...")
	http.ListenAndServe(":5001", handler)
}

type httpReadWriteCloser struct {
	in  io.Reader
	out io.Writer
}

func (h *httpReadWriteCloser) Read(p []byte) (n int, err error)  { return h.in.Read(p) }
func (h *httpReadWriteCloser) Write(p []byte) (n int, err error) { return h.out.Write(p) }
func (h *httpReadWriteCloser) Close() error                      { return nil }
