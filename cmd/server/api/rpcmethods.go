package api

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type RPCMethod int

type BlockTag string

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
	"eth_getLogs":                             RPCMethodEthGetLogs,
	"eth_getStorageAt":                        RPCMethodEthGetStorageAt,
	"eth_getTransactionByBlockHashAndIndex":   RPCMethodEthGetTransactionByBlockHashAndIndex,
	"eth_getTransactionByBlockNumberAndIndex": RPCMethodEthGetTransactionByBlockNumberAndIndex,
	"eth_getTransactionCount":                 RPCMethodEthGetTransactionCount,
	"eth_getUncleCountByBlockNumber":          RPCMethodEthGetUncleCountByBlockNumber,
	"eth_newFilter":                           RPCMethodEthNewFilter,
}

var rpcMethodsWithHash = map[string]RPCMethod{
	"eth_getBlockByHash":        RPCMethodEthGetBlockByHash,
	"eth_getTransactionReceipt": RPCMethodEthGetTransactionReceipt,
	"eth_getTransactionByHash":  RPCMethodEthGetTransactionByHash,
}

var rpcMethodsWithCountByHash = map[string]RPCMethod{
	"eth_getUncleCountByBlockHash":       RPCMethodEthGetUncleCountByBlockHash,
	"eth_getBlockTransactionCountByHash": RPCMethodEthGetBlockTransactionCountByHash,
}

func HasRPCMethodWithBlocknumber(methodName string) (RPCMethod, bool) {
	method, exists := rpcMethodsWithBlockNumber[methodName]
	return method, exists
}

func ReplaceBlockTag(req *JSONRPCRequest, method RPCMethod, blockNumberUnode string) (*JSONRPCRequest, error) {

	if len(req.Params) == 0 {
		return req, nil
	}

	switch method {
	case RPCMethodEthGetBlockByNumber:
		// blocknumber is the first param for this method
		ReplaceBlockTagWithHash(req, 0, blockNumberUnode)
	case RPCMethodEthGetBalance, RPCMethodEthCall:
		// blocknumber is the second param for this method
		ReplaceBlockTagWithHash(req, 1, blockNumberUnode)
		// case RPCMethodEthGetBlockTransactionCountByNumber:
		// 	params.Params = []string{blockNumber}
		// case RPCMethodEthGetTransactionByBlockHashAndIndex:
		// 	params.Params = []string{blockNumber, params.Params[1]}
		// case RPCMethodEthGetTransactionByBlockNumberAndIndex:
		// 	params.Params = []string{blockNumber, params.Params[1]}
		// case RPCMethodEthGetTransactionCount:
		// 	params.Params = []string{params.Params[0], blockNumber}
		// case RPCMethodEthGetUncleCountByBlockNumber:
		// 	params.Params = []string{blockNumber}
		// case RPCMethodEthNewFilter:
		// 	params.Params = []string{params.Params[0], blockNumber}
	}

	return req, nil

}

func ReplaceBlockTagWithHash(req *JSONRPCRequest, position int, blockNumberHash string) error {
	blockNumberRequest, err := rawMessageToString(req.Params[position])
	if err != nil {
		return err
	}
	// check if blockNumberRequest starts with 0x
	if len(blockNumberRequest) > 2 && blockNumberRequest[:2] == "0x" {
		c, err := CompareHex(blockNumberRequest, blockNumberHash)
		if err != nil {
			return err
		}
		if c == 1 {
			return fmt.Errorf("invalid block number: %s", blockNumberRequest)
		}
	} else if BlockTag(blockNumberRequest) == Latest || BlockTag(blockNumberRequest) == Pending {
		req.Params[position] = stringToRawMessage(blockNumberHash)
	}
	return nil
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
	bigInt1, ok := new(big.Int).SetString(hex1, 16)
	if !ok {
		return 0, fmt.Errorf("invalid hexadecimal number: %s", hex1)
	}

	bigInt2, ok := new(big.Int).SetString(hex2, 16)
	if !ok {
		return 0, fmt.Errorf("invalid hexadecimal number: %s", hex2)
	}

	return bigInt1.Cmp(bigInt2), nil
}
