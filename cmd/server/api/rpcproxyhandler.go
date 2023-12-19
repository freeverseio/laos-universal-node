package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/freeverseio/laos-universal-node/internal/state"
)

// RPCProxyHandler
func (h *RPCProxyHandler) HandleProxyRPC(r *http.Request, req JSONRPCRequest, stateService state.Service) RPCResponse {
	// check if we have to replace the block tag
	method, hasBlockNumber := h.proxyRPCMethodManager.HasRPCMethodWithBlocknumber(req.Method)
	if hasBlockNumber {
		blockNumber, errBlock := getBlockNumberFromDb(stateService)
		if errBlock != nil {
			return getErrorResponse(fmt.Errorf("error getting block number from db: %w", errBlock), req.ID)
		}
		req, errBlockTag := h.proxyRPCMethodManager.ReplaceBlockTag(&req, method, blockNumber)
		if errBlockTag != nil {
			return getErrorResponse(fmt.Errorf("error replacing block tag: %w", errBlockTag), req.ID)
		}
	}
	// JSONRPCRequest to []byte
	body, err := json.Marshal(req)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error marshalling request: %w", err), req.ID)
	}

	// Prepare the request to the BC node
	proxyReq, err := http.NewRequest(r.Method, h.GetRpcUrl(), io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating request: %w", err), req.ID)
	}

	// Forward headers the request
	for name, values := range r.Header {
		for _, value := range values {
			// we don't want to forward the Accept-Encoding header because we don't want to receive a encoded response (e.g. gzip)
			if name != "Accept-Encoding" {
				proxyReq.Header.Set(name, value)
			}
		}
	}

	// Send the request to the Ethereum node
	resp, err := h.GetHttpClient().Do(proxyReq)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error sending request: %w", err), req.ID)
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			slog.Error("error closing response body", "err", errClose)
		}
	}() // Check error on Close

	response, err := getJsonRPCResponse(resp)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting JSON RPC response: %w", err), req.ID)
	}
	// check if we have to replace the block tag
	method, hasBlockHash := h.proxyRPCMethodManager.HasRPCMethodWithHash(req.Method)
	if hasBlockHash {
		blockNumber, errBlock := getBlockNumberFromDb(stateService)
		if errBlock != nil {
			return getErrorResponse(fmt.Errorf("error getting block number from db: %w", errBlock), req.ID)
		}
		errCheck := h.proxyRPCMethodManager.CheckBlockNumberFromResponseFromHashCalls(response, method, blockNumber)
		if errCheck != nil {
			return getErrorResponse(errCheck, req.ID)
		}
	}

	return *response
}

func getJsonRPCResponse(r *http.Response) (*RPCResponse, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for further handling
	var req RPCResponse
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("error parsing JSON request: %w", err)
	}
	return &req, nil
}

func getBlockNumberFromDb(stateService state.Service) (string, error) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	blockNumber, err := tx.GetCurrentOwnershipBlock()
	if err != nil {
		return "", fmt.Errorf("error getting current block number: %w", err)
	}
	// minus 1 because we want to return the last tagged block
	return fmt.Sprintf("0x%x", blockNumber-1), nil
}
