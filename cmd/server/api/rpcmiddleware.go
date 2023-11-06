package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/freeverseio/laos-universal-node/internal/scan"
)

// JSONRPCRequest represents the expected structure of a JSON-RPC request.
type JSONRPCRequest struct {
	JSONRPC string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params"`
	ID      *json.RawMessage `json:"id,omitempty"` // Pointer allows for an optional field
}

type ParamsRPCRequest struct {
	Data  string `json:"data"`
	To    string `json:"to"`
	From  string `json:"from"`
	Value string `json:"value"`
}

// Adjust the middleware to handle different JSON-RPC methods
func PostRpcRequestMiddleware(standardHandler, erc721Handler http.Handler, st scan.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.Header.Get("Content-Type") == "application/json" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}

			// It's important to restore the body so that the next handler can read it
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			var req JSONRPCRequest
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
				return
			}

			// Direct the request to the appropriate handler based on the method
			if req.Method == "eth_call" {
				var params ParamsRPCRequest
				if err := json.Unmarshal(req.Params, &params); err != nil {
					http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
					return
				}
				isInContractList, err := checkContractInList(params.To, st)
				if err != nil {
					http.Error(w, "Error checking contract in list", http.StatusBadRequest)
					return
				}
				if isInContractList {
					erc721Handler.ServeHTTP(w, r)
				} else {
					standardHandler.ServeHTTP(w, r)
				}
			} else {
				standardHandler.ServeHTTP(w, r)
			}

			return // Stop the middleware chain here
		}

		// Call the next middleware/handler if it's not a JSON-RPC call
		standardHandler.ServeHTTP(w, r)
	})
}

func checkContractInList(contractAddress string, st scan.Storage) (bool, error) {
	list, err := st.ReadAll(context.Background())
	if err != nil {
		return false, err
	}
	for _, contract := range list {
		if contract.Address.Hex() == contractAddress {
			return true, nil
		}
	}
	return false, nil
}

// func checkingMethodFromCallData(data string) (bool, error) {
// 	erc721.NewCallData(data)

// 	return true, nil
// }

func middleware(h HandlerInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && r.Header.Get("Content-Type") == "application/json" {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Error reading request body", http.StatusInternalServerError)
					return
				}

				errCloser := r.Body.Close() // Must close the original body
				if errCloser != nil {       // Check for errors closing the body
					http.Error(w, "Error closing request body", http.StatusInternalServerError)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body)) // Create a new body with the same data

				var req JSONRPCRequest
				if err := json.Unmarshal(body, &req); err != nil {
					http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
					return
				}

				if req.JSONRPC != "2.0" {
					http.Error(w, "Invalid JSON-RPC version", http.StatusBadRequest)
					return
				}
				h.SetJsonRPCRequest(req)
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
