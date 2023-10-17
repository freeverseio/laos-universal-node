package erc721

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/laos-universal-node/internal/blockchain/erc721"
)

// calldata represents the data for an ERC721 function call.
type CallData []byte

// NewCalldata creates a new calldata instance from a hexadecimal string.
func NewCallData(s string) (CallData, error) { // nolint:revive // TODO refactor: unexported return type
	b, err := hexutil.Decode(s)
	if err != nil {
		return CallData{}, err
	}
	return b, nil
}

// erc721method represents the supported ERC721 methods.
type erc721method int

const (
	OwnerOf erc721method = iota
	BalanceOf
	TokenURI
	URI
	SupportsInterface
	Name
	ContractURI
	Decimals
	Symbol
)

var methodSigs = map[string]erc721method{
	hexutil.Encode(crypto.Keccak256([]byte("ownerOf(uint256)"))[:4]):          OwnerOf,
	hexutil.Encode(crypto.Keccak256([]byte("balanceOf(address)"))[:4]):        BalanceOf,
	hexutil.Encode(crypto.Keccak256([]byte("tokenURI(uint256)"))[:4]):         TokenURI,
	hexutil.Encode(crypto.Keccak256([]byte("supportsInterface(bytes4)"))[:4]): SupportsInterface,
	hexutil.Encode(crypto.Keccak256([]byte("name()"))[:4]):                    Name,
	hexutil.Encode(crypto.Keccak256([]byte("contractURI()"))[:4]):             ContractURI,
	hexutil.Encode(crypto.Keccak256([]byte("decimals()"))[:4]):                Decimals,
	hexutil.Encode(crypto.Keccak256([]byte("symbol()"))[:4]):                  Symbol,
}

// Method returns the ERC721 method invoked by the calldata.
func (b CallData) Method() (erc721method, error) {
	sig, err := b.methodSignature()
	if err != nil {
		return 0, err
	}

	if method, exists := methodSigs[sig]; exists {
		return method, nil
	}

	return 0, fmt.Errorf("unallowed method: %s", sig)
}

// methodSignature returns the method signature of the calldata.
func (b CallData) methodSignature() (string, error) {
	if len(b) < 4 {
		return "", fmt.Errorf("invalid call data, incomplete method signature (%d bytes < 4)", len(b))
	}
	return hexutil.Encode(b[:4]), nil
}

// GetParam returns the value of a specific parameter in the input arguments of the calldata.
func (b CallData) GetParam(param string) (interface{}, error) {
	inputArgs, err := b.getInputArgs()
	if err != nil {
		return nil, err
	}
	if inputArgs[param] == nil {
		return nil, fmt.Errorf("param %s not found into the object %v", param, inputArgs)
	}
	return inputArgs[param], nil
}

// getInputArgs returns a map of input arguments extracted from the calldata.
func (b CallData) getInputArgs() (map[string]interface{}, error) {
	sig, err := b.methodSignature()
	if err != nil {
		return nil, err
	}

	argdata := b[4:]
	if len(argdata)%32 != 0 {
		return nil, fmt.Errorf("invalid call data; length should be a multiple of 32 bytes (was %d)", len(argdata))
	}
	erc721Abi, err := abi.JSON(strings.NewReader(erc721.Erc721ABI))
	if err != nil {
		return nil, err
	}
	id, err := hexutil.Decode(sig)
	if err != nil {
		return nil, err
	}
	method, err := erc721Abi.MethodById(id)
	if err != nil {
		return nil, err
	}
	inputArgs := map[string]interface{}{}
	err = method.Inputs.UnpackIntoMap(inputArgs, argdata)
	if err != nil {
		return nil, err
	}
	return inputArgs, nil
}
