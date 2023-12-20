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

type Block struct {
	Difficulty       string            `json:"difficulty,omitempty"`
	ExtraData        string            `json:"extraData,omitempty"`
	GasLimit         string            `json:"gasLimit,omitempty"`
	GasUsed          string            `json:"gasUsed,omitempty"`
	Hash             string            `json:"hash,omitempty"`
	LogsBloom        string            `json:"logsBloom,omitempty"`
	Miner            string            `json:"miner,omitempty"`
	MixHash          string            `json:"mixHash,omitempty"`
	Nonce            string            `json:"nonce,omitempty"`
	Number           string            `json:"number,omitempty"`
	ParentHash       string            `json:"parentHash,omitempty"`
	ReceiptsRoot     string            `json:"receiptsRoot,omitempty"`
	Sha3Uncles       string            `json:"sha3Uncles,omitempty"`
	Size             string            `json:"size,omitempty"`
	StateRoot        string            `json:"stateRoot,omitempty"`
	Timestamp        string            `json:"timestamp,omitempty"`
	TotalDifficulty  string            `json:"totalDifficulty,omitempty"`
	Transactions     []json.RawMessage `json:"transactions,omitempty"`
	TransactionsRoot string            `json:"transactionsRoot,omitempty"`
	Uncles           []json.RawMessage `json:"uncles,omitempty"`
}

type Transaction struct {
	BlockHash            string            `json:"blockHash,omitempty"`
	BlockNumber          string            `json:"blockNumber,omitempty"`
	From                 string            `json:"from,omitempty"`
	Gas                  string            `json:"gas,omitempty"`
	GasPrice             string            `json:"gasPrice,omitempty"`
	MaxFeePerGas         string            `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string            `json:"maxPriorityFeePerGas,omitempty"`
	Hash                 string            `json:"hash,omitempty"`
	Input                string            `json:"input,omitempty"`
	Nonce                string            `json:"nonce,omitempty"`
	To                   string            `json:"to,omitempty"`
	TransactionIndex     string            `json:"transactionIndex,omitempty"`
	Value                string            `json:"value,omitempty"`
	Type                 string            `json:"type,omitempty"`
	AccessList           []json.RawMessage `json:"accessList,omitempty"`
	ChainId              string            `json:"chainId,omitempty"`
	V                    string            `json:"v,omitempty"`
	R                    string            `json:"r,omitempty"`
	S                    string            `json:"s,omitempty"`
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

var rpcMethodsWithHash = map[string]RPCMethod{
	"eth_getBlockByHash":                    RPCMethodEthGetBlockByHash,
	"eth_getTransactionReceipt":             RPCMethodEthGetTransactionReceipt,
	"eth_getTransactionByHash":              RPCMethodEthGetTransactionByHash,
	"eth_getTransactionByBlockHashAndIndex": RPCMethodEthGetTransactionByBlockHashAndIndex,
}

// we don't handle these methods for now but we have to handle this in the future
//
//	var rpcMethodsWithCountByHash = map[string]RPCMethod{
//		"eth_getUncleCountByBlockHash":       RPCMethodEthGetUncleCountByBlockHash,
//		"eth_getBlockTransactionCountByHash": RPCMethodEthGetBlockTransactionCountByHash,
//	}
type RPCMethodManager interface {
	HasRPCMethodWithBlockNumber(methodName string) (RPCMethod, bool)
	HasRPCMethodWithHash(methodName string) (RPCMethod, bool)
	CheckBlockNumberFromResponseFromHashCalls(resp *RPCResponse, method RPCMethod, blockNumberUnode string) error
	ReplaceBlockTag(req *JSONRPCRequest, method RPCMethod, blockNumberUnode string) (*JSONRPCRequest, error)
}

type ProxyRPCMethodManager struct{}

func NewProxyRPCMethodManager() RPCMethodManager {
	return &ProxyRPCMethodManager{}
}

func (b *ProxyRPCMethodManager) HasRPCMethodWithBlockNumber(methodName string) (RPCMethod, bool) {
	method, exists := rpcMethodsWithBlockNumber[methodName]
	return method, exists
}

func (b *ProxyRPCMethodManager) HasRPCMethodWithHash(methodName string) (RPCMethod, bool) {
	method, exists := rpcMethodsWithHash[methodName]
	return method, exists
}

func (b *ProxyRPCMethodManager) CheckBlockNumberFromResponseFromHashCalls(resp *RPCResponse, method RPCMethod, blockNumberUnode string) error {
	var blockNumber string
	var err error

	switch method {
	case RPCMethodEthGetBlockByHash:
		var block Block
		blockNumber, err = unmarshalAndGetBlockNumber(resp, &block)
	case RPCMethodEthGetTransactionByHash, RPCMethodEthGetTransactionReceipt, RPCMethodEthGetTransactionByBlockHashAndIndex:
		var tx Transaction
		blockNumber, err = unmarshalAndGetBlockNumber(resp, &tx)
	}

	if err != nil {
		return err
	}

	c, err := CompareHex(blockNumber, blockNumberUnode)
	if err != nil {
		return err
	}
	if c > 0 { // blockNumber > blockNumberUnode
		return fmt.Errorf("invalid block number: %s", blockNumber)
	}
	return nil
}

func (b *ProxyRPCMethodManager) ReplaceBlockTag(req *JSONRPCRequest, method RPCMethod, blockNumberUnode string) (*JSONRPCRequest, error) {
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

func unmarshalAndGetBlockNumber(resp *RPCResponse, v interface{}) (string, error) {
	err := json.Unmarshal(*resp.Result, v)
	if err != nil {
		return "", err
	}

	switch value := v.(type) {
	case *Block:
		return value.Number, nil
	case *Transaction:
		return value.BlockNumber, nil
	}

	return "", fmt.Errorf("unknown type for unmarshalling")
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
