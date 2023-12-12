package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
	"github.com/freeverseio/laos-universal-node/internal/state"
)

func (h *UniversalMintingRPCHandler) HandleUniversalMinting(jsonRPCRequest JSONRPCRequest, stateService state.Service) RPCResponse {

	// if call is eth_blockNumber we should return the latest block number
	if jsonRPCRequest.Method == "eth_blockNumber" {
		return blockNumber(stateService)
	}

	var params ParamsRPCRequest
	if len(jsonRPCRequest.Params) == 0 || json.Unmarshal(jsonRPCRequest.Params[0], &params) != nil {
		return getErrorResponse(fmt.Errorf("error parsing params or missing params"))
	}

	blockNumber := "latest" // if this by chance does not exist in param use the latest block
	if len(jsonRPCRequest.Params) == 2 {
		if errUnmarshal := json.Unmarshal(jsonRPCRequest.Params[1], &blockNumber); errUnmarshal != nil {
			return getErrorResponse(fmt.Errorf("error parsing block number: %w", errUnmarshal))
		}
	}

	calldata, err := erc721.NewCallData(params.Data)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error parsing calldata: %w", err))
	}

	if method, exists, err := calldata.UniversalMintingMethod(); err != nil {
		return getErrorResponse(fmt.Errorf("error parsing calldata: %w", err))
	} else if !exists {
		return getErrorResponse(fmt.Errorf("method not supported"))
	} else {
		switch method {
		case erc721.OwnerOf:
			return ownerOf(calldata, params, blockNumber, stateService)
		case erc721.BalanceOf:
			return balanceOf(calldata, params, blockNumber, stateService)
		case erc721.TotalSupply:
			return totalSupply(params, blockNumber, stateService)
		case erc721.TokenOfOwnerByIndex:
			return tokenOfOwnerByIndex(calldata, params, blockNumber, stateService)
		case erc721.TokenByIndex:
			return tokenByIndex(calldata, params, blockNumber, stateService)
		case erc721.SupportsInterface:
			return supportsInterface()
		}
	}
	return getErrorResponse(fmt.Errorf("method not supported"))
}

func supportsInterface() RPCResponse {
	// calldata already checked for SupportsInterface 0x780e9d63
	// if we are here it means that the calldata is SupportsInterface(0x780e9d63)
	// so we can return true
	return getResponse("0x0000000000000000000000000000000000000000000000000000000000000001", nil)
}

func ownerOf(callData erc721.CallData, params ParamsRPCRequest, blockNumber string, stateService state.Service) RPCResponse {
	tokenID, err := getParamBigInt(callData, "tokenId")
	if err != nil {
		return getErrorResponse(err)
	}
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err))
	}

	owner, err := tx.OwnerOf(common.HexToAddress(params.To), tokenID)
	// Format the address to include leading zeros as 40-character (160 bits) hexadecimal string
	// TODO check if there is a better way to do this
	fullAddressString := fmt.Sprintf("0x000000000000000000000000%040x", owner)
	return getResponse(fullAddressString, err)
}

func balanceOf(callData erc721.CallData, params ParamsRPCRequest, blockNumber string, stateService state.Service) RPCResponse {
	ownerAddress, err := getParamAddress(callData, "owner")
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting owner: %w", err))
	}
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err))
	}

	balance, err := tx.BalanceOf(common.HexToAddress(params.To), ownerAddress)
	return getResponse(fmt.Sprintf("0x%064x", balance), err)
}

func totalSupply(params ParamsRPCRequest, blockNumber string, stateService state.Service) RPCResponse {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err := loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err))
	}
	totalSupply, err := tx.TotalSupply(common.HexToAddress(params.To))
	return getResponse(fmt.Sprintf("0x%064x", totalSupply), err)
}

func tokenOfOwnerByIndex(callData erc721.CallData, params ParamsRPCRequest, blockNumber string, stateService state.Service) RPCResponse {
	index, err := getParamBigInt(callData, "index")
	if err != nil {
		slog.Error("Error getting tokenId", "err", err)
		return getErrorResponse(fmt.Errorf("error getting tokenId: %w", err))
	}
	ownerAddress, err := getParamAddress(callData, "owner")
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting owner: %w", err))
	}
	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err))
	}
	tokenId, err := tx.TokenOfOwnerByIndex(common.HexToAddress(params.To), ownerAddress, int(index.Int64()))
	return getResponse(fmt.Sprintf("0x%064x", tokenId), err)
}

func tokenByIndex(callData erc721.CallData, params ParamsRPCRequest, blockNumber string, stateService state.Service) RPCResponse {
	index, err := getParamBigInt(callData, "index")
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting tokenId: %w", err))
	}

	tx := stateService.NewTransaction()
	defer tx.Discard()
	tx, err = loadMerkleTree(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err))
	}
	tokenId, err := tx.TokenByIndex(common.HexToAddress(params.To), int(index.Int64()))
	return getResponse(fmt.Sprintf("0x%064x", tokenId), err)
}

func blockNumber(stateService state.Service) RPCResponse {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	blockNumber, err := tx.GetCurrentOwnershipBlock()
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting current block number: %w", err))
	}
	// minus 1 because we want to return the last tagged block
	return getResponse(fmt.Sprintf("0x%x", blockNumber-1), nil)
}

func getResponse(result string, err error) RPCResponse {
	if err != nil {
		return getErrorResponse(err)
	}

	return RPCResponse{
		Jsonrpc: "2.0",
		ID:      RpcId,
		Result:  result,
	}
}

func getErrorResponse(err error) RPCResponse {
	slog.Error("Failed to send response", "err", err)
	errorResponse := RPCResponse{
		Jsonrpc: "2.0",
		ID:      errorId,
		Error: &RPCError{
			Code:    ErrorCodeInvalidRequest,
			Message: "execution reverted",
		},
	}

	return errorResponse
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
