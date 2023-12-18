package api

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type RPCMethod int

type BlockTag string

type FilterObject struct {
	FromBlock string            `json:"fromBlock,omitempty"`
	ToBlock   string            `json:"toBlock,omitempty"`
	Address   string            `json:"address,omitempty"`
	Topics    []json.RawMessage `json:"topics,omitempty"`
	Blockhash *json.RawMessage  `json:"blockhash,omitempty"`
}

// Constants for each RPC method
const (
	RPCMethodEthCall RPCMethod = iota
	RPCMethodEthGetBalance
	RPCMethodEthGetBlockByNumber
	RPCMethodEthGetBlockTransactionCountByNumber
	RPCMethodEthGetCode
	RPCMethodEthGetLogs
	RPCMethodEthGetStorageAt
	RPCMethodEthGetTransactionByBlockHashAndIndex
	RPCMethodEthGetTransactionByBlockNumberAndIndex
	RPCMethodEthGetTransactionCount
	RPCMethodEthGetUncleCountByBlockNumber
	RPCMethodEthNewFilter
	RPCMethodEthGetBlockByHash
	RPCMethodEthGetTransactionReceipt
	RPCMethodEthGetTransactionByHash
	RPCMethodEthGetUncleCountByBlockHash
	RPCMethodEthGetBlockTransactionCountByHash
)

const (
	Latest    BlockTag = "latest"
	Pending   BlockTag = "pending"
	Earliest  BlockTag = "earliest"
	Finalized BlockTag = "finalized"
	Safe      BlockTag = "safe"
)

// Map of RPC method names to their corresponding constants
var rpcMethodsWithBlockNumber = map[string]RPCMethod{
	"eth_call":                                RPCMethodEthCall,
	"eth_getBalance":                          RPCMethodEthGetBalance,
	"eth_getBlockByNumber":                    RPCMethodEthGetBlockByNumber,
	"eth_getBlockTransactionCountByNumber":    RPCMethodEthGetBlockTransactionCountByNumber,
	"eth_getCode":                             RPCMethodEthGetCode,
	"eth_getStorageAt":                        RPCMethodEthGetStorageAt,
	"eth_getTransactionByBlockNumberAndIndex": RPCMethodEthGetTransactionByBlockNumberAndIndex,
	"eth_getTransactionCount":                 RPCMethodEthGetTransactionCount,
	"eth_getUncleCountByBlockNumber":          RPCMethodEthGetUncleCountByBlockNumber,
	"eth_getLogs":                             RPCMethodEthGetLogs,
	"eth_newFilter":                           RPCMethodEthNewFilter,
}

// var rpcMethodsWithHash = map[string]RPCMethod{
// 	"eth_getBlockByHash":                    RPCMethodEthGetBlockByHash,
// 	"eth_getTransactionReceipt":             RPCMethodEthGetTransactionReceipt,
// 	"eth_getTransactionByHash":              RPCMethodEthGetTransactionByHash,
// 	"eth_getTransactionByBlockHashAndIndex": RPCMethodEthGetTransactionByBlockHashAndIndex,
// }

// var rpcMethodsWithCountByHash = map[string]RPCMethod{
// 	"eth_getUncleCountByBlockHash":       RPCMethodEthGetUncleCountByBlockHash,
// 	"eth_getBlockTransactionCountByHash": RPCMethodEthGetBlockTransactionCountByHash,
// }

func HasRPCMethodWithBlocknumber(methodName string) (RPCMethod, bool) {
	method, exists := rpcMethodsWithBlockNumber[methodName]
	return method, exists
}

func ReplaceBlockTag(req *JSONRPCRequest, method RPCMethod, blockNumberUnode string) (*JSONRPCRequest, error) {
	if len(req.Params) == 0 {
		return req, nil
	}

	switch method {
	case RPCMethodEthGetBlockByNumber,
		RPCMethodEthGetBlockTransactionCountByNumber,
		RPCMethodEthGetTransactionByBlockNumberAndIndex,
		RPCMethodEthGetUncleCountByBlockNumber:
		// blocknumber is the first param for this method
		err := replaceBlockTagWithHash(req, 0, blockNumberUnode)
		if err != nil {
			return nil, err
		}
	case RPCMethodEthGetBalance,
		RPCMethodEthCall,
		RPCMethodEthGetCode,
		RPCMethodEthGetTransactionCount:
		// blocknumber is the second param for this method
		err := replaceBlockTagWithHash(req, 1, blockNumberUnode)
		if err != nil {
			return nil, err
		}
	case RPCMethodEthGetStorageAt:
		// blocknumber is the third param for this method
		err := replaceBlockTagWithHash(req, 2, blockNumberUnode)
		if err != nil {
			return nil, err
		}
	case RPCMethodEthGetLogs, RPCMethodEthNewFilter:
		err := replaceBlockTagFromObject(req, blockNumberUnode)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}

func replaceBlockTagWithHash(req *JSONRPCRequest, position int, blockNumberHash string) error {
	blockNumberRequest, err := rawMessageToString(req.Params[position])
	if err != nil {
		return err
	}
	blockNumber, err := getBlockNumber(blockNumberRequest, blockNumberHash)
	if err != nil {
		return err
	}
	req.Params[position] = stringToRawMessage(blockNumber)
	return nil
}

func replaceBlockTagFromObject(req *JSONRPCRequest, blockNumberHash string) error {
	var filterObject FilterObject
	err := json.Unmarshal(req.Params[0], &filterObject)
	if err != nil {
		return err
	}

	changed := false
	if filterObject.FromBlock != "" {
		blockNumber, errBlock := getBlockNumber(filterObject.FromBlock, blockNumberHash)
		if errBlock != nil {
			return errBlock
		}
		if blockNumber != filterObject.FromBlock {
			filterObject.FromBlock = blockNumber
			changed = true
		}
	}

	if filterObject.ToBlock != "" {
		blockNumber, errBlock := getBlockNumber(filterObject.ToBlock, blockNumberHash)
		if errBlock != nil {
			return errBlock
		}
		if blockNumber != filterObject.ToBlock {
			filterObject.ToBlock = blockNumber
			changed = true
		}
	}

	if changed {
		req.Params[0], err = json.Marshal(filterObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func getBlockNumber(blockNumberRequest, blockNumberHash string) (string, error) {
	// Using switch to handle different cases
	switch {
	case len(blockNumberRequest) > 2 && blockNumberRequest[:2] == "0x":
		c, err := CompareHex(blockNumberRequest, blockNumberHash)
		if err != nil {
			return "", err
		}
		if c == 1 {
			return "", fmt.Errorf("invalid block number: %s", blockNumberRequest)
		}
		return blockNumberRequest, nil

	case BlockTag(blockNumberRequest) == Latest:
		return blockNumberHash, nil

	case BlockTag(blockNumberRequest) == Pending:
		pendingBlockNumber, err := addIntNumberToHex(blockNumberHash, 1)
		if err != nil {
			return "", err
		}
		return pendingBlockNumber, nil

	default:
		return blockNumberRequest, nil
	}
}

func rawMessageToString(raw json.RawMessage) (string, error) {
	var str string
	err := json.Unmarshal(raw, &str)
	if err != nil {
		return "", err
	}
	return str, nil
}

func stringToRawMessage(str string) json.RawMessage {
	quotedResult := fmt.Sprintf(`%q`, str)
	return json.RawMessage(quotedResult)
}

// CompareHex compares two hexadecimal strings and returns:
// -1 if hex1 < hex2
//
//	0 if hex1 == hex2
//	1 if hex1 > hex2
func CompareHex(hex1, hex2 string) (int, error) {
	bigInt1, ok := new(big.Int).SetString(hex1[2:], 16)
	if !ok {
		return 0, fmt.Errorf("invalid hexadecimal number: %s", hex1)
	}

	bigInt2, ok := new(big.Int).SetString(hex2[2:], 16)
	if !ok {
		return 0, fmt.Errorf("invalid hexadecimal number: %s", hex2)
	}

	return bigInt1.Cmp(bigInt2), nil
}

func addIntNumberToHex(hex string, value int) (string, error) {
	// Convert the hex string to a big.Int
	bigInt1, ok := new(big.Int).SetString(hex[2:], 16)
	if !ok {
		return "", fmt.Errorf("invalid hexadecimal number: %s", hex)
	}

	// Convert the int value to big.Int and add it to the first big.Int
	bigInt2 := big.NewInt(int64(value))
	bigInt1.Add(bigInt1, bigInt2)

	// Convert the result back to a hex string
	return fmt.Sprintf("0x%x", bigInt1), nil
}
