package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
	"github.com/freeverseio/laos-universal-node/internal/state"
)

const (
	RpcId = 1
)

type RPCResponder interface{}

type RPCResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  string    `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
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

	rpcRequests, err := parseBody(body)
	if err != nil {
		http.Error(w, ErrMsgBadRequest, http.StatusBadRequest)
		return
	}
	responseBody := make([]RPCResponse, 0, len(rpcRequests))
	for _, rpcRequest := range rpcRequests {
		responseBody = append(responseBody, h.GetRPCResponse(r, rpcRequest))
	}
	w.Header().Set("Content-Type", "application/json")
	if len(responseBody) == 1 {
		err = json.NewEncoder(w).Encode(responseBody[0])
	} else {
		err = json.NewEncoder(w).Encode(responseBody)
	}

	if err != nil {
		slog.Error("Failed to send response", "err", err)
	}
}

func (h *GlobalRPCHandler) HandleUniversalMinting(r *http.Request, stateService state.Service) RPCResponse {
	return h.UniversalMintingRPCHandler.HandleUniversalMinting(r, stateService)
}

func (h *GlobalRPCHandler) HandleProxyRPC(r *http.Request) RPCResponse {
	return h.ProxyRPCHandler.HandleProxyRPC(r)
}

func (h *GlobalRPCHandler) GetRPCResponse(r *http.Request, req JSONRPCRequest) RPCResponse {
	if req.JSONRPC != "2.0" {
		return getErrorResponse(fmt.Errorf("invalid JSON-RPC version"))
	}
	switch req.Method {
	case "eth_call":
		return h.handleEthCallMethod(r, &req)
	case "eth_blockNumber":
		return h.HandleUniversalMinting(r, h.stateService)
	default:
		return h.HandleProxyRPC(r)
	}
}

func (h *GlobalRPCHandler) handleEthCallMethod(r *http.Request, req *JSONRPCRequest) RPCResponse {
	var params ParamsRPCRequest
	if len(req.Params) == 0 || json.Unmarshal(req.Params[0], &params) != nil {
		return getErrorResponse(fmt.Errorf("error parsing params or missing params"))
	}

	// Check for universal minting method.
	isRemoteMinting, err := isUniversalMintingMethod(params.Data)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error checking for universal minting method: %w", err))
	}

	// If not related to remote minting, delegate to standard handler.
	if !isRemoteMinting {
		return h.HandleProxyRPC(r)
	}

	// Check if contract is stored
	contractExists, err := isContractStored(params.To, h.stateService)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error checking contract list: %w", err))
	}

	// If contract is stored, use the specific handler for ERC721 universal minting.
	if contractExists {
		return h.HandleUniversalMinting(r, h.stateService)
	} else {
		return h.HandleProxyRPC(r)
	}
}

func isContractStored(contractAddress string, stateService state.Service) (bool, error) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	lowerCaseContractAddress := strings.ToLower(contractAddress)
	contract, err := tx.Get(state.ContractPrefix + lowerCaseContractAddress)
	if err != nil {
		return false, err
	}

	if contract != nil {
		lowerCaseContract := strings.ToLower(string(contract))
		if lowerCaseContract != "" {
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

func parseBody(body []byte) ([]JSONRPCRequest, error) {
	// First, try to unmarshal as a single JSONRPCRequest
	var singleReq JSONRPCRequest
	if err := json.Unmarshal(body, &singleReq); err == nil {
		return []JSONRPCRequest{singleReq}, nil
	}

	// If single unmarshalling fails, try as an array
	var multiReq []JSONRPCRequest
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields() // Prevent unknown fields
	if err := decoder.Decode(&multiReq); err != nil {
		return nil, fmt.Errorf("error parsing JSON request: %w", err)
	}

	return multiReq, nil
}
