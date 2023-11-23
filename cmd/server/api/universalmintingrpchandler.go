package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
	"github.com/freeverseio/laos-universal-node/internal/state"
)

const (
	RpcId = 1
)

type RPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

type JSONRPCErrorResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (h *GlobalRPCHandler) UniversalMintingRPCHandler(w http.ResponseWriter, r *http.Request) {
	jsonRPCRequest := h.GetJsonRPCRequest()
	var params ParamsRPCRequest
	if len(jsonRPCRequest.Params) == 0 || json.Unmarshal(jsonRPCRequest.Params[0], &params) != nil {
		http.Error(w, "Error parsing params or missing params", http.StatusBadRequest)
		return
	}

	var blockNumber string
	if len(jsonRPCRequest.Params) == 2 {
		if err := json.Unmarshal(jsonRPCRequest.Params[1], &blockNumber); err != nil {
			http.Error(w, "Error parsing block number", http.StatusBadRequest)
			return
		}
	}
	slog.Info("block number", "blockNumber", blockNumber)

	calldata, err := erc721.NewCallData(params.Data)
	if err != nil {
		http.Error(w, "Error parsing calldata", http.StatusBadRequest)
		return
	}

	if method, exists, err := calldata.UniversalMintingMethod(); err != nil {
		http.Error(w, "Error parsing calldata", http.StatusBadRequest)
		return
	} else if !exists {
		http.Error(w, "Method not supported", http.StatusBadRequest)
		return
	} else {
		switch method {
		case erc721.OwnerOf:
			ownerOf(calldata, params, h.stateService, w)
		case erc721.BalanceOf:
			balanceOf(calldata, params, h.stateService, w)
		case erc721.TotalSupply:
			totalSupply(params, h.stateService, w)
		case erc721.TokenOfOwnerByIndex:
			tokenOfOwnerByIndex(calldata, params, h.stateService, w)
		case erc721.TokenByIndex:
			tokenByIndex(calldata, params, h.stateService, w)
		}
	}
}

func ownerOf(callData erc721.CallData, params ParamsRPCRequest, stateService state.Service, w http.ResponseWriter) {
	tokenID, err := getParamBigInt(callData, "tokenId")
	if err != nil {
		slog.Error("Error getting tokenId", "err", err)
		sendErrorResponse(w, err)
		return
	}
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}

	owner, err := tx.OwnerOf(common.HexToAddress(params.To), tokenID)
	sendResponse(w, owner.Hex(), err)
}

func balanceOf(callData erc721.CallData, params ParamsRPCRequest, stateService state.Service, w http.ResponseWriter) {
	ownerAddress, err := getParamAddress(callData, "owner")
	if err != nil {
		slog.Error("Error getting owner", "err", err)
		sendErrorResponse(w, err)
		return
	}

	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}

	balance, err := tx.BalanceOf(common.HexToAddress(params.To), ownerAddress)
	sendResponse(w, balance.String(), err)
}

func totalSupply(params ParamsRPCRequest, stateService state.Service, w http.ResponseWriter) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err := createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}
	totalSupply, err := tx.TotalSupply(common.HexToAddress(params.To))
	sendResponse(w, strconv.FormatInt(totalSupply, 10), err)
}

func tokenOfOwnerByIndex(callData erc721.CallData, params ParamsRPCRequest, stateService state.Service, w http.ResponseWriter) {
	index, err := getParamBigInt(callData, "index")
	if err != nil {
		slog.Error("Error getting tokenId", "err", err)
		sendErrorResponse(w, err)
		return
	}
	ownerAddress, err := getParamAddress(callData, "owner")
	if err != nil {
		slog.Error("Error getting owner", "err", err)
		sendErrorResponse(w, err)
		return
	}
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}
	tokenId, err := tx.TokenOfOwnerByIndex(common.HexToAddress(params.To), ownerAddress, int(index.Int64()))
	sendResponse(w, tokenId.String(), err)
}

func tokenByIndex(callData erc721.CallData, params ParamsRPCRequest, stateService state.Service, w http.ResponseWriter) {
	index, err := getParamBigInt(callData, "index")
	if err != nil {
		slog.Error("Error getting tokenId", "err", err)
		sendErrorResponse(w, err)
		return
	}

	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}
	tokenId, err := tx.TokenByIndex(common.HexToAddress(params.To), int(index.Int64()))
	sendResponse(w, tokenId.String(), err)
}

func createMerkleTrees(tx state.Tx, contactAddress common.Address) (state.Tx, error) {
	ownershipTree, enumeratedTree, enumeratedtotalTree, err := tx.CreateTreesForContract(contactAddress)
	if err != nil {
		return nil, err
	}

	err = tx.SetTreesForContract(contactAddress, ownershipTree, enumeratedTree, enumeratedtotalTree)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func sendResponse(w http.ResponseWriter, result string, err error) {
	if err != nil {
		slog.Error("Failed to send response", "err", err)
		sendErrorResponse(w, err)
		return
	}

	response := RPCResponse{
		Jsonrpc: "2.0",
		ID:      RpcId,
		Result:  result,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Failed to send response", "err", err)
	}
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	slog.Error("Failed to send response", "err", err)

	errorResponse := JSONRPCErrorResponse{
		JSONRPC: "2.0",
		ID:      ErrorId,
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    ErrorCodeInvalidRequest,
			Message: err.Error(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(errorResponse)
	if err != nil {
		slog.Error("Failed to send response", "err", err)
	}
}

func getParamBigInt(callData erc721.CallData, paramName string) (*big.Int, error) {
	param, err := callData.GetParam(paramName)
	if err != nil {
		return nil, err
	}

	bigIntParam, ok := param.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid %s", paramName)
	}

	return bigIntParam, nil
}

func getParamAddress(callData erc721.CallData, paramName string) (common.Address, error) {
	param, err := callData.GetParam(paramName)
	if err != nil {
		return common.Address{}, err
	}

	addressParam, ok := param.(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf("invalid %s", paramName)
	}

	return addressParam, nil
}
