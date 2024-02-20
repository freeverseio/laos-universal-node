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
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

func (h *UniversalMintingRPCHandler) HandleUniversalMinting(r *http.Request, jsonRPCRequest JSONRPCRequest, stateService state.Service) RPCResponse {
	// if call is eth_blockNumber we should return the latest block number
	if jsonRPCRequest.Method == "eth_blockNumber" {
		return blockNumber(stateService, jsonRPCRequest.ID)
	}

	var params ethCallParamsRPCRequest
	if len(jsonRPCRequest.Params) == 0 || json.Unmarshal(jsonRPCRequest.Params[0], &params) != nil {
		return getErrorResponse(fmt.Errorf("error parsing params or missing params"), jsonRPCRequest.ID)
	}

	blockNumber := "latest" // if this by chance does not exist in param use the latest block
	if len(jsonRPCRequest.Params) == 2 {
		if errUnmarshal := json.Unmarshal(jsonRPCRequest.Params[1], &blockNumber); errUnmarshal != nil {
			return getErrorResponse(fmt.Errorf("error parsing block number: %w", errUnmarshal), jsonRPCRequest.ID)
		}
	}

	calldata, err := erc721.NewCallData(params.Data)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error parsing calldata: %w", err), jsonRPCRequest.ID)
	}

	if method, exists, err := calldata.UniversalMintingMethod(); err != nil {
		return getErrorResponse(fmt.Errorf("error parsing calldata: %w", err), jsonRPCRequest.ID)
	} else if !exists {
		return getErrorResponse(fmt.Errorf("method not supported"), jsonRPCRequest.ID)
	} else {
		switch method {
		case erc721.OwnerOf:
			return ownerOf(calldata, params, blockNumber, stateService, jsonRPCRequest.ID)
		case erc721.BalanceOf:
			return balanceOf(calldata, params, blockNumber, stateService, jsonRPCRequest.ID)
		case erc721.TotalSupply:
			return totalSupply(params, blockNumber, stateService, jsonRPCRequest.ID)
		case erc721.TokenOfOwnerByIndex:
			return tokenOfOwnerByIndex(calldata, params, blockNumber, stateService, jsonRPCRequest.ID)
		case erc721.TokenByIndex:
			return tokenByIndex(calldata, params, blockNumber, stateService, jsonRPCRequest.ID)
		case erc721.TokenURI:
			return h.tokenURI(r, blockNumber, stateService, jsonRPCRequest)
		case erc721.SupportsInterface:
			return supportsInterface(jsonRPCRequest.ID)
		}
	}
	return getErrorResponse(fmt.Errorf("method not supported"), jsonRPCRequest.ID)
}

func supportsInterface(id *json.RawMessage) RPCResponse {
	// calldata already checked for SupportsInterface 0x780e9d63
	// if we are here it means that the calldata is SupportsInterface(0x780e9d63)
	// so we can return true
	return getResponse("0x0000000000000000000000000000000000000000000000000000000000000001", id, nil)
}

func ownerOf(callData erc721.CallData, params ethCallParamsRPCRequest, blockNumber string, stateService state.Service, id *json.RawMessage) RPCResponse {
	tokenID, err := getParamBigInt(callData, "tokenId")
	if err != nil {
		return getErrorResponse(err, id)
	}
	tx, err := stateService.NewTransaction()
	if err != nil {
		return getErrorResponse(err, id)
	}
	defer tx.Discard()
	err = checkoutBlock(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err), id)
	}

	owner, err := tx.OwnerOf(common.HexToAddress(params.To), tokenID)
	// Format the address to include leading zeros as 40-character (160 bits) hexadecimal string
	// TODO check if there is a better way to do this - there is go-ethereum's abi.Arguments.Pack, but it uses reflection, it might be too slow
	fullAddressString := fmt.Sprintf("0x000000000000000000000000%040x", owner)
	return getResponse(fullAddressString, id, err)
}

func balanceOf(callData erc721.CallData, params ethCallParamsRPCRequest, blockNumber string, stateService state.Service, id *json.RawMessage) RPCResponse {
	ownerAddress, err := getParamAddress(callData, "owner")
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting owner: %w", err), id)
	}
	tx, err := stateService.NewTransaction()
	if err != nil {
		return getErrorResponse(err, id)
	}
	defer tx.Discard()
	err = checkoutBlock(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err), id)
	}

	balance, err := tx.BalanceOf(common.HexToAddress(params.To), ownerAddress)
	return getResponse(fmt.Sprintf("0x%064x", balance), id, err)
}

func totalSupply(params ethCallParamsRPCRequest, blockNumber string, stateService state.Service, id *json.RawMessage) RPCResponse {
	tx, err := stateService.NewTransaction()
	if err != nil {
		return getErrorResponse(err, id)
	}
	defer tx.Discard()
	err = checkoutBlock(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err), id)
	}
	totalSupply, err := tx.TotalSupply(common.HexToAddress(params.To))
	return getResponse(fmt.Sprintf("0x%064x", totalSupply), id, err)
}

func tokenOfOwnerByIndex(callData erc721.CallData, params ethCallParamsRPCRequest, blockNumber string, stateService state.Service, id *json.RawMessage) RPCResponse {
	index, err := getParamBigInt(callData, "index")
	if err != nil {
		slog.Error("Error getting tokenId", "err", err)
		return getErrorResponse(fmt.Errorf("error getting tokenId: %w", err), id)
	}
	ownerAddress, err := getParamAddress(callData, "owner")
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting owner: %w", err), id)
	}
	tx, err := stateService.NewTransaction()
	if err != nil {
		return getErrorResponse(err, id)
	}
	defer tx.Discard()
	err = checkoutBlock(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err), id)
	}
	tokenId, err := tx.TokenOfOwnerByIndex(common.HexToAddress(params.To), ownerAddress, int(index.Int64()))
	return getResponse(fmt.Sprintf("0x%064x", tokenId), id, err)
}

func tokenByIndex(callData erc721.CallData, params ethCallParamsRPCRequest, blockNumber string, stateService state.Service, id *json.RawMessage) RPCResponse {
	index, err := getParamBigInt(callData, "index")
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting tokenId: %w", err), id)
	}

	tx, err := stateService.NewTransaction()
	if err != nil {
		return getErrorResponse(err, id)
	}
	defer tx.Discard()
	err = checkoutBlock(tx, common.HexToAddress(params.To), blockNumber)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error creating merkle trees: %w", err), id)
	}
	tokenId, err := tx.TokenByIndex(common.HexToAddress(params.To), int(index.Int64()))
	return getResponse(fmt.Sprintf("0x%064x", tokenId), id, err)
}

func (h *UniversalMintingRPCHandler) tokenURI(r *http.Request, blockNumber string, stateService state.Service, req JSONRPCRequest) RPCResponse {
	// get the evo block corresponding to the ownership block in time from storage
	evoBlock, err := getMappedEvoBlockNumber(blockNumber, stateService)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting evo block number: %w", err), req.ID)
	}
	// change the original request to use the corresponding evo block
	req.Params[1] = stringToRawMessage(evoBlock)
	// retrieve the latest evo block number from storage in case the request is for the latest block
	lastEvoBlock, err := getLastEvoBlockNumberFromStorage(stateService)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting last evo block number: %w", err), req.ID)
	}

	// TODO get contract address from req.Params[0] (Params[0] must be parsed!), call "tx.GetCollectionAddress(contract)" and
	// replace the contract address of the request with the collection address

	errBlockTag := h.rpcMethodManager.ReplaceBlockTag(&req, RPCMethodEthCall, lastEvoBlock)
	if errBlockTag != nil {
		return getErrorResponse(fmt.Errorf("error replacing block tag: %w", errBlockTag), req.ID)
	}
	// JSONRPCRequest to []byte
	body, err := json.Marshal(req)
	if err != nil {
		return getErrorResponse(fmt.Errorf("error marshalling request: %w", err), req.ID)
	}
	// Prepare the request to the BC node
	proxyReq, err := http.NewRequest(r.Method, h.rpcUrl, io.NopCloser(bytes.NewReader(body)))
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

	return *response
}

func blockNumber(stateService state.Service, id *json.RawMessage) RPCResponse {
	tx, err := stateService.NewTransaction()
	if err != nil {
		return getErrorResponse(err, id)
	}
	defer tx.Discard()
	block, err := tx.GetLastOwnershipBlock()
	if err != nil {
		return getErrorResponse(fmt.Errorf("error getting current block number: %w", err), id)
	}

	return getResponse(fmt.Sprintf("0x%x", block.Number), id, nil)
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

func checkoutBlock(tx state.Tx, contractAddress common.Address, blockNumber string) error {
	// if block is not latest we should checkout tree for that tag
	// it is important that this transaction is not commit which is always the case for this transaction
	if blockNumber != "latest" {
		num, err := strconv.ParseInt(strings.Replace(blockNumber, "0x", "", 1), 16, 64)
		if err != nil {
			slog.Error("wrong block number", "err", err)
			return err
		}

		err = tx.Checkout(num)
		if err != nil {
			slog.Error("error occurred checking out merkle tree at block number", "block_number", num, "err", err)
			return err
		}
	}

	// loading trees after checkout
	err := tx.LoadContractTrees(contractAddress)
	if err != nil {
		return err
	}
	return nil
}

func isBlockTag(blockNumber string) bool {
	switch blockTag(blockNumber) {
	case latest, pending, earliest, finalized, safe:
		return true
	}
	return false
}

// TODO maybe pass tx instead of stateService
func getMappedEvoBlockNumber(ownershipBlock string, stateService state.Service) (string, error) {
	evoBlock := ownershipBlock
	if !isBlockTag(ownershipBlock) {
		tx, err := stateService.NewTransaction()
		if err != nil {
			return "", err
		}
		defer tx.Discard()
		parsedBlockNumber, err := strconv.ParseUint(ownershipBlock, 10, 64)
		if err != nil {
			return "", err
		}
		evoBlockNumber, err := tx.GetMappedEvoBlockNumber(parsedBlockNumber)
		if err != nil {
			return "", err
		}
		evoBlock = fmt.Sprintf("0x%x", evoBlockNumber)
	}
	return evoBlock, nil
}

func getLastEvoBlockNumberFromStorage(stateService state.Service) (string, error) {
	tx, err := stateService.NewTransaction()
	if err != nil {
		return "", fmt.Errorf("error creating a new transaction: %w", err)
	}
	defer tx.Discard()
	block, err := tx.GetLastEvoBlock()
	if err != nil {
		return "", fmt.Errorf("error getting evo block number: %w", err)
	}

	return fmt.Sprintf("0x%x", block.Number), nil
}
