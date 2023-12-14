package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
	"github.com/freeverseio/laos-universal-node/internal/state"
)

type RPCResponse struct {
	Jsonrpc string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Result  string           `json:"result,omitempty"`
	Error   *RPCError        `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *GlobalRPCHandler) PostRPCRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "No JSON RPC call or invalid Content-Type", http.StatusBadRequest)
		return
	}

	// Read the body of the incoming request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrMsgBadRequest, http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for further handling

	defer func() {
		errClose := r.Body.Close()
		if errClose != nil {
			slog.Error("error closing response body", "err", errClose)
		}
	}()

	rpcRequests, isArrayRequest, err := parseBody(body)
	if err != nil {
		http.Error(w, ErrMsgBadRequest, http.StatusBadRequest)
		return
	}
	responseBody := make([]RPCResponse, 0, len(rpcRequests))
	for _, rpcRequest := range rpcRequests {
		responseBody = append(responseBody, h.GetRPCResponse(r, rpcRequest))
	}
	w.Header().Set("Content-Type", "application/json")

	if isArrayRequest {
		err = json.NewEncoder(w).Encode(responseBody)
	} else {
		err = json.NewEncoder(w).Encode(responseBody[0])
	}

	if err != nil {
		slog.Error("Failed to send response", "err", err)
	}
}

func (h *GlobalRPCHandler) HandleUniversalMinting(req JSONRPCRequest, stateService state.Service) RPCResponse {
	return h.universalMintingRPCHandler.HandleUniversalMinting(req, stateService)
}

func (h *GlobalRPCHandler) HandleProxyRPC(r *http.Request, req JSONRPCRequest) RPCResponse {
	return h.rpcProxyHandler.HandleProxyRPC(r, req)
}

func (h *GlobalRPCHandler) GetRPCResponse(r *http.Request, req JSONRPCRequest) RPCResponse {
	if req.JSONRPC != "2.0" {
		return getErrorResponse(fmt.Errorf("invalid JSON-RPC version"), req.ID)
	}
	switch req.Method {
	case "eth_call":
		return h.handleEthCallMethod(r, req)
	case "eth_blockNumber":
		return h.HandleUniversalMinting(req, h.stateService)
	default:
		return h.HandleProxyRPC(r, req)
	}
}

func (h *GlobalRPCHandler) handleEthCallMethod(r *http.Request, req JSONRPCRequest) RPCResponse {
	var params ParamsRPCRequest
	if len(req.Params) == 0 || json.Unmarshal(req.Params[0], &params) != nil {
		return getErrorResponse(fmt.Errorf("error parsing params or missing params"), req.ID)
	}

	// Check for universal minting method.
	isUniversalMinting, err := isUniversalMintingMethod(params.Data)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error checking for universal minting method: %w", err), req.ID)
	}

	// If not related to remote minting, delegate to standard handler.
	if !isUniversalMinting {
		return h.HandleProxyRPC(r, req)
	}

	// Check if contract is stored
	contractExists, err := isContractStored(params.To, h.stateService)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error checking contract list: %w", err), req.ID)
	}

	// If contract is stored, use the specific handler for ERC721 universal minting.
	if contractExists {
		return h.HandleUniversalMinting(req, h.stateService)
	} else {
		return h.HandleProxyRPC(r, req)
	}
}

func isContractStored(contractAddress string, stateService state.Service) (bool, error) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	return tx.HasERC721UniversalContract(contractAddress)
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

func parseBody(body []byte) (request []JSONRPCRequest, isArray bool, err error) {
	// First, try to unmarshal as a single JSONRPCRequest
	var singleReq JSONRPCRequest
	if err := json.Unmarshal(body, &singleReq); err == nil {
		return []JSONRPCRequest{singleReq}, false, nil
	}

	// If single unmarshalling fails, try as an array
	var multiReq []JSONRPCRequest
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields() // Prevent unknown fields
	if err := decoder.Decode(&multiReq); err != nil {
		return nil, false, fmt.Errorf("error parsing JSON request: %w", err)
	}

	return multiReq, true, nil
}

func getResponse(result string, id *json.RawMessage, err error) RPCResponse {
	if err != nil {
		return getErrorResponse(err, id)
	}
	return RPCResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  result,
	}
}

func getErrorResponse(err error, id *json.RawMessage) RPCResponse {
	slog.Error("Failed to send response", "err", err)

	errorResponse := RPCResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    ErrorCodeInvalidRequest,
			Message: "execution reverted",
		},
	}

	return errorResponse
}
