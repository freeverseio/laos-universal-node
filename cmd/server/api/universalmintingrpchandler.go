package api

import (
	"encoding/json"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
	"github.com/freeverseio/laos-universal-node/internal/state"
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
	josonRPCRequest := h.GetJsonRPCRequest()
	var params ParamsRPCRequest
	if len(josonRPCRequest.Params) == 0 || json.Unmarshal(josonRPCRequest.Params[0], &params) != nil {
		http.Error(w, "Error parsing params or missing params", http.StatusBadRequest)
		return
	}

	var blockNumber string
	if len(josonRPCRequest.Params) == 2 {
		if err := json.Unmarshal(josonRPCRequest.Params[1], &blockNumber); err != nil {
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
	tokenIDParam, err := callData.GetParam("tokenId")
	if err != nil {
		slog.Error("Error getting tokenId", "err", err)
		http.Error(w, "Error getting tokenId", http.StatusBadRequest)
		return
	}
	//
	tokenID, ok := tokenIDParam.(*big.Int)
	if !ok {
		slog.Error("Invalid tokenId", "tokenID", tokenID)
		http.Error(w, "Invalid tokenId ", http.StatusBadRequest)
		return
	}
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		http.Error(w, "Error creating merkle trees", http.StatusInternalServerError)
		return
	}
	owner, err := tx.OwnerOf(common.HexToAddress(params.To), tokenID)
	if err != nil {
		slog.Error("Error getting owner", "err", err)
		http.Error(w, "Error getting owner", http.StatusInternalServerError)
		return
	}

	response := RPCResponse{
		Jsonrpc: "2.0",
		ID:      1,
		Result:  owner.Hex(),
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode and send the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error in case the response couldn't be sent
		slog.Error("Failed to send response", "err", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func balanceOf(callData erc721.CallData, params ParamsRPCRequest, stateService state.Service, w http.ResponseWriter) {
	ownerAddressParam, err := callData.GetParam("owner")
	if err != nil {
		slog.Error("Error getting owner", "err", err)
		http.Error(w, "Error getting owner", http.StatusBadRequest)
		return
	}
	ownerAddress := ownerAddressParam.(common.Address)

	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		http.Error(w, "Error creating merkle trees", http.StatusInternalServerError)
		return
	}
	balance, err := tx.BalanceOf(common.HexToAddress(params.To), ownerAddress)
	if err != nil {
		slog.Error("Error getting balance", "err", err)
		http.Error(w, "Error getting balance", http.StatusInternalServerError)
		return
	}

	response := RPCResponse{
		Jsonrpc: "2.0",
		ID:      1,
		Result:  balance.String(),
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode and send the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error in case the response couldn't be sent
		slog.Error("Failed to send response", "err", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func totalSupply(params ParamsRPCRequest, stateService state.Service, w http.ResponseWriter) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err := createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		http.Error(w, "Error creating merkle trees", http.StatusInternalServerError)
		return
	}
	totalSupply, err := tx.TotalSupply(common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error getting balance", "err", err)
		// TODO format error
		http.Error(w, "Error getting balance", http.StatusInternalServerError)
		return
	}

	response := RPCResponse{
		Jsonrpc: "2.0",
		ID:      1,
		Result:  strconv.FormatInt(totalSupply, 10),
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode and send the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error in case the response couldn't be sent
		slog.Error("Failed to send response", "err", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func tokenOfOwnerByIndex(callData erc721.CallData, params ParamsRPCRequest, stateService state.Service, w http.ResponseWriter) {
	indexParam, err := callData.GetParam("index")
	if err != nil {
		slog.Error("Error getting owner", "err", err)
		http.Error(w, "Error getting owner", http.StatusBadRequest)
		return
	}
	index := indexParam.(*big.Int)

	ownerAddressParam, err := callData.GetParam("owner")
	if err != nil {
		slog.Error("Error getting owner", "err", err)
		http.Error(w, "Error getting owner", http.StatusBadRequest)
		return
	}
	ownerAddress := ownerAddressParam.(common.Address)

	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		// TODO format error
		http.Error(w, "Error creating merkle trees", http.StatusInternalServerError)
		return
	}
	tokenId, err := tx.TokenOfOwnerByIndex(common.HexToAddress(params.To), ownerAddress, int(index.Int64()))
	if err != nil {
		slog.Error("Error getting balance", "err", err)
		// TODO format error
		http.Error(w, "Error what the fuck", http.StatusInternalServerError)
		return
	}

	response := RPCResponse{
		Jsonrpc: "2.0",
		ID:      1,
		Result:  tokenId.String(),
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode and send the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error in case the response couldn't be sent
		slog.Error("Failed to send response", "err", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func tokenByIndex(callData erc721.CallData, params ParamsRPCRequest, stateService state.Service, w http.ResponseWriter) {
	indexParam, err := callData.GetParam("index")
	if err != nil {
		slog.Error("Error getting owner", "err", err)
		http.Error(w, "Error getting owner", http.StatusBadRequest)
		return
	}
	index := indexParam.(*big.Int)

	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = createMerkleTrees(tx, common.HexToAddress(params.To))
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		// TODO format error
		http.Error(w, "Error creating merkle trees", http.StatusInternalServerError)
		return
	}
	tokenId, err := tx.TokenByIndex(common.HexToAddress(params.To), int(index.Int64()))
	if err != nil {
		slog.Error("Error getting balance", "err", err)
		// TODO format error
		http.Error(w, "Error getting balance", http.StatusInternalServerError)
		return
	}

	response := RPCResponse{
		Jsonrpc: "2.0",
		ID:      1,
		Result:  tokenId.String(),
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode and send the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error in case the response couldn't be sent
		slog.Error("Failed to send response", "err", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
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
