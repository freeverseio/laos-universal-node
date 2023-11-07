package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
	"github.com/freeverseio/laos-universal-node/internal/scan"
)

// JSONRPCRequest represents the expected structure of a JSON-RPC request.
type JSONRPCRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
	ID      *json.RawMessage  `json:"id,omitempty"` // Pointer allows for an optional field
}

type ParamsRPCRequest struct {
	Data  string `json:"data,omitempty"`
	To    string `json:"to,omitempty"`
	From  string `json:"from,omitempty"`
	Value string `json:"value,omitempty"`
}

func PostRpcRequestMiddleware(standardHandler, erc721Handler http.Handler, st scan.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a valid JSON-RPC POST request
		if valid, body := validateJSONRPCPostRequest(w, r); valid {
			handleJSONRPCRequest(w, r, body, standardHandler, erc721Handler, st)
		}
	})
}

// validateJSONRPCPostRequest checks if the request is a valid JSON-RPC POST request and reads the body.
func validateJSONRPCPostRequest(w http.ResponseWriter, r *http.Request) (valid bool, body []byte) {
	if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
		reportError(w, "No JSON RPC call or invalid Content-Type", http.StatusBadRequest)
		return false, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		reportError(w, "Error reading request body", http.StatusBadRequest)
		return false, nil
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for further handling
	return true, body
}

// handleJSONRPCRequest processes the JSON-RPC request by forwarding to the appropriate handler.
func handleJSONRPCRequest(w http.ResponseWriter, r *http.Request, body []byte, standardHandler, erc721Handler http.Handler, st scan.Storage) {
	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		reportError(w, "Error parsing JSON request", http.StatusBadRequest)
		return
	}
	if req.JSONRPC != "2.0" {
		http.Error(w, "Invalid JSON-RPC version", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "eth_call":
		handleEthCallMethod(w, r, req, standardHandler, erc721Handler, st)
	default:
		standardHandler.ServeHTTP(w, r)
	}
}

// handleEthCallMethod processes an eth_call method request.
func handleEthCallMethod(w http.ResponseWriter, r *http.Request, req JSONRPCRequest, standardHandler, erc721Handler http.Handler, st scan.Storage) {
	var params ParamsRPCRequest
	if len(req.Params) == 0 || json.Unmarshal(req.Params[0], &params) != nil {
		reportError(w, "Error parsing params or missing params", http.StatusBadRequest)
		return
	}

	if remoteMinting, _ := isUniversalMintingMethod(params.Data); remoteMinting {
		if isInContractList, _ := isContractInList(params.To, st); isInContractList {
			erc721Handler.ServeHTTP(w, r)
			return
		}
	}
	// Fall back to standard handler if not a remote minting method or not in the contract list.
	standardHandler.ServeHTTP(w, r)
}

// reportError is a utility function for reporting errors through HTTP.
func reportError(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}

func isContractInList(contractAddress string, st scan.Storage) (bool, error) {
	list, err := st.ReadAll(context.Background())
	if err != nil {
		return false, err
	}
	for _, contract := range list {
		addr := contract.Address.Hex() // convert to string
		if strings.EqualFold(addr, contractAddress) {
			return true, nil
		}
	}
	return false, nil
}

func isUniversalMintingMethod(data string) (bool, error) {
	calldata, err := erc721.NewCallData(data)
	if err != nil {
		return false, err
	}
	_, exists, err := calldata.UniversalMintingMethod()
	if err != nil {
		return false, err
	}
	return exists, nil
}
