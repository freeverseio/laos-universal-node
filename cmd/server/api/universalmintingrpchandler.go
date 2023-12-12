package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"
	"strings"

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
	jsonRPCRequest, err := getJsonRPCRequest(r)
	if err != nil {
		http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
		return
	}

	// if call is eth_blockNumber we should return the latest block number
	if jsonRPCRequest.Method == "eth_blockNumber" {
		blockNumber(w, h.stateService)
		return
	}

	var params ParamsRPCRequest
	if len(jsonRPCRequest.Params) == 0 || json.Unmarshal(jsonRPCRequest.Params[0], &params) != nil {
		http.Error(w, "Error parsing params or missing params", http.StatusBadRequest)
		return
	}

	blockNumber := "latest" // if this by chance does not exist in param use the latest block
	if len(jsonRPCRequest.Params) == 2 {
		if errUnmarshal := json.Unmarshal(jsonRPCRequest.Params[1], &blockNumber); errUnmarshal != nil {
			http.Error(w, "Error parsing block number", http.StatusBadRequest)
			return
		}
	}
	slog.Debug("block number", "blockNumber", blockNumber)

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
			ownerOf(calldata, params, blockNumber, h.stateService, w)
		case erc721.BalanceOf:
			balanceOf(calldata, params, blockNumber, h.stateService, w)
		case erc721.TotalSupply:
			totalSupply(params, blockNumber, h.stateService, w)
		case erc721.TokenOfOwnerByIndex:
			tokenOfOwnerByIndex(calldata, params, blockNumber, h.stateService, w)
		case erc721.TokenByIndex:
			tokenByIndex(calldata, params, blockNumber, h.stateService, w)
		case erc721.TokenURI:
			h.tokenURI(calldata, params, blockNumber, w)
		case erc721.SupportsInterface:
			supportsInterface(w)
		}
	}
}

func supportsInterface(w http.ResponseWriter) {
	// calldata already checked for SupportsInterface 0x780e9d63
	// if we are here it means that the calldata is SupportsInterface(0x780e9d63)
	// so we can return true
	sendResponse(w, "0x0000000000000000000000000000000000000000000000000000000000000001", nil)
}

func getJsonRPCRequest(r *http.Request) (*JSONRPCRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for further handling
	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("error parsing JSON request: %w", err)
	}
	return &req, nil
}

func ownerOf(callData erc721.CallData, params ParamsRPCRequest, blockNumber string, stateService state.Service, w http.ResponseWriter) {
	tokenID, err := getParamBigInt(callData, "tokenId")
	if err != nil {
		slog.Error("Error getting tokenId", "err", err)
		sendErrorResponse(w, err)
		return
	}
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}

	owner, err := tx.OwnerOf(common.HexToAddress(params.To), tokenID)
	// Format the address to include leading zeros as 40-character (160 bits) hexadecimal string
	// TODO check if there is a better way to do this
	fullAddressString := fmt.Sprintf("0x000000000000000000000000%040x", owner)
	sendResponse(w, fullAddressString, err)
}

func balanceOf(callData erc721.CallData, params ParamsRPCRequest, blockNumber string, stateService state.Service, w http.ResponseWriter) {
	ownerAddress, err := getParamAddress(callData, "owner")
	if err != nil {
		slog.Error("Error getting owner", "err", err)
		sendErrorResponse(w, err)
		return
	}
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}

	balance, err := tx.BalanceOf(common.HexToAddress(params.To), ownerAddress)
	// TODO check if there is a better way to format the balance
	sendResponse(w, fmt.Sprintf("0x%064x", balance), err)
}

func totalSupply(params ParamsRPCRequest, blockNumber string, stateService state.Service, w http.ResponseWriter) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err := loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}
	totalSupply, err := tx.TotalSupply(common.HexToAddress(params.To))
	sendResponse(w, fmt.Sprintf("0x%064x", totalSupply), err)
}

func tokenOfOwnerByIndex(callData erc721.CallData, params ParamsRPCRequest, blockNumber string, stateService state.Service, w http.ResponseWriter) {
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
	tx, err = loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}
	tokenId, err := tx.TokenOfOwnerByIndex(common.HexToAddress(params.To), ownerAddress, int(index.Int64()))
	sendResponse(w, fmt.Sprintf("0x%064x", tokenId), err)
}

func tokenByIndex(callData erc721.CallData, params ParamsRPCRequest, blockNumber string, stateService state.Service, w http.ResponseWriter) {
	index, err := getParamBigInt(callData, "index")
	if err != nil {
		slog.Error("Error getting tokenId", "err", err)
		sendErrorResponse(w, err)
		return
	}

	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		slog.Error("Error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}
	tokenId, err := tx.TokenByIndex(common.HexToAddress(params.To), int(index.Int64()))
	sendResponse(w, fmt.Sprintf("0x%064x", tokenId), err)
}

func (h *GlobalRPCHandler) tokenURI(callData erc721.CallData, params ParamsRPCRequest, blockNumber string, w http.ResponseWriter) {
	// TODO test me
	tokenID, err := getParamBigInt(callData, "tokenId")
	if err != nil {
		slog.Error("error getting tokenId", "err", err)
		sendErrorResponse(w, err)
		return
	}

	tx := h.stateService.NewTransaction()
	defer tx.Discard()
	tx, err = loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		slog.Error("error creating merkle trees", "err", err)
		sendErrorResponse(w, err)
		return
	}
	tokenURI, err := tx.TokenURI(common.HexToAddress(params.To), tokenID)
	if err != nil {
		slog.Error("error retrieving token URI", "err", err)
		sendErrorResponse(w, err)
		return
	}
	encodedValue, err := erc721.AbiEncodeString(tokenURI)
	sendResponse(w, encodedValue, err)
}

func blockNumber(w http.ResponseWriter, stateService state.Service) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	blockNumber, err := tx.GetCurrentOwnershipBlock()
	if err != nil {
		slog.Error("Error getting current block number", "err", err)
		sendErrorResponse(w, err)
		return
	}
	// minus 1 because we want to return the last tagged block
	sendResponse(w, fmt.Sprintf("0x%x", blockNumber-1), nil)
}

func loadMerkleTree(tx state.Tx, contractAddress common.Address, blockNumber string) (state.Tx, error) {
	ownershipTree, enumeratedTree, enumeratedtotalTree, err := tx.CreateTreesForContract(contractAddress)
	if err != nil {
		return nil, err
	}

	tx.SetTreesForContract(contractAddress, ownershipTree, enumeratedTree, enumeratedtotalTree)
	// if block is not latest we should checkout tree for that tag
	// it is important that this transaction is not commit which is always the case for this transaction
	if blockNumber != "latest" {
		num, err := strconv.ParseInt(strings.Replace(blockNumber, "0x", "", 1), 16, 64)
		if err != nil {
			slog.Error("wrong block number", "err", err)
			return nil, err
		}

		err = tx.Checkout(contractAddress, num)
		if err != nil {
			slog.Error("error occurred checking out merkle tree at block number", "block_number", num,
				"contract_address", contractAddress, "err", err)
			return nil, err
		}
	}
	return tx, nil
}

func sendResponse(w http.ResponseWriter, result string, err error) {
	if err != nil {
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
		ID:      errorId,
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    ErrorCodeInvalidRequest,
			Message: "execution reverted",
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
