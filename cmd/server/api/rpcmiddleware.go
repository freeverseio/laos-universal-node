package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
	"github.com/freeverseio/laos-universal-node/internal/repository"
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

func PostRpcRequestMiddleware(h RPCHandler, repositoryService repository.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we pass both handlers and decide which one to call based on the request
		proxyRPCHandler := http.HandlerFunc(h.PostRPCProxyHandler)                   // proxy handler for standard requests
		universalMintingRPCHandler := http.HandlerFunc(h.UniversalMintingRPCHandler) // handler for universal minting requests
		// Check for a valid JSON-RPC POST request
		if valid, body := validateJSONRPCPostRequest(w, r); valid {
			handleJSONRPCRequest(w, r, body, proxyRPCHandler, universalMintingRPCHandler, repositoryService)
		}
	})
}

// validateJSONRPCPostRequest checks if the request is a valid JSON-RPC POST request and reads the body.
func validateJSONRPCPostRequest(w http.ResponseWriter, r *http.Request) (valid bool, request *JSONRPCRequest) {
	if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "No JSON RPC call or invalid Content-Type", http.StatusBadRequest)
		return false, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		slog.Error("error reading request body", "err", err)
		return false, nil
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for further handling

	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
		slog.Error("error parsing JSON request", "err", err)
		return
	}
	if req.JSONRPC != "2.0" {
		http.Error(w, "Invalid JSON-RPC version", http.StatusBadRequest)
		return
	}

	return true, &req
}

// handleJSONRPCRequest processes the JSON-RPC request by forwarding to the appropriate handler.
func handleJSONRPCRequest(w http.ResponseWriter, r *http.Request, jsonRequest *JSONRPCRequest, proxyRPCHandler, universalMintingHandler http.Handler, repositoryService repository.Service) {
	switch jsonRequest.Method {
	case "eth_call":
		handleEthCallMethod(w, r, jsonRequest, proxyRPCHandler, universalMintingHandler, repositoryService)
	default:
		proxyRPCHandler.ServeHTTP(w, r)
	}
}

func handleEthCallMethod(w http.ResponseWriter, r *http.Request, req *JSONRPCRequest, proxyRPCHandler, universalMintingHandler http.Handler, repositoryService repository.Service) {
	var params ParamsRPCRequest
	if len(req.Params) == 0 || json.Unmarshal(req.Params[0], &params) != nil {
		http.Error(w, "Error parsing params or missing params", http.StatusBadRequest)
		return
	}

	// Check for universal minting method.
	isRemoteMinting, err := isUniversalMintingMethod(params.Data)
	if err != nil {
		http.Error(w, "Error checking for universal minting method: "+err.Error(), http.StatusBadRequest)
		return
	}

	// If not related to remote minting, delegate to standard handler.
	if !isRemoteMinting {
		proxyRPCHandler.ServeHTTP(w, r)
		return
	}

	// Check if contract is in the list.
	isInContractList, err := isContractInList(params.To, repositoryService)
	if err != nil {
		http.Error(w, "Error checking contract list: "+err.Error(), http.StatusBadRequest)
		return
	}

	// If contract is in the list, use the specific handler for ERC721 universal minting.
	if isInContractList {
		universalMintingHandler.ServeHTTP(w, r)
		return
	} else {
		proxyRPCHandler.ServeHTTP(w, r)
		return
	}
}

func isContractInList(contractAddress string, repositoryService repository.Service) (bool, error) {
	list, err := repositoryService.GetAllERC721UniversalContracts()
	if err != nil {
		return false, err
	}
	for _, contract := range list {
		if strings.EqualFold(contract, contractAddress) {
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
