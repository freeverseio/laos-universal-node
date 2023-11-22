package erc721

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain/contract"
)

// calldata represents the data for an ERC721 function call.
type CallData []byte

// NewCalldata creates a new calldata instance from a hexadecimal string.
func NewCallData(s string) (CallData, error) {
	b, err := hexutil.Decode(s)
	if err != nil {
		return CallData{}, err
	}
	return b, nil
}

// erc721method represents the supported ERC721 methods.
type Erc721method int

const (
	ShortAddressLength = 4
	CallDataLength     = 32 // call data must have a length of a multiple of 32 bytes
)

const (
	NotSupported Erc721method = iota
	OwnerOf
	BalanceOf
	TotalSupply
	TokenOfOwnerByIndex
	TokenByIndex
)

// remoteMintingMethodSigs represents the method signatures of the ERC721 methods that are part of the remote minting service.
var remoteMintingMethodSigs = map[string]Erc721method{
	hexutil.Encode(crypto.Keccak256([]byte("ownerOf(uint256)"))[:ShortAddressLength]):                     OwnerOf,
	hexutil.Encode(crypto.Keccak256([]byte("balanceOf(address)"))[:ShortAddressLength]):                   BalanceOf,
	hexutil.Encode(crypto.Keccak256([]byte("totalSupply()"))[:ShortAddressLength]):                        TotalSupply,
	hexutil.Encode(crypto.Keccak256([]byte("tokenOfOwnerByIndex(address,uint256)"))[:ShortAddressLength]): TokenOfOwnerByIndex,
	hexutil.Encode(crypto.Keccak256([]byte("tokenByIndex(uint256)"))[:ShortAddressLength]):                TokenByIndex,
}

// Method returns if the calldata is a supported remote minting ERC721 method and the method.
func (b CallData) UniversalMintingMethod() (Erc721method, bool, error) {
	sig, err := b.methodSignature()
	if err != nil {
		return 0, false, err
	}

	if method, exists := remoteMintingMethodSigs[sig]; exists {
		return method, exists, nil
	}

	return 0, false, nil
}

// methodSignature returns the method signature of the calldata.
func (b CallData) methodSignature() (string, error) {
	if len(b) < ShortAddressLength {
		return "", fmt.Errorf("invalid call data, incomplete method signature (%d bytes < 4)", len(b))
	}
	return hexutil.Encode(b[:ShortAddressLength]), nil
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

	argdata := b[ShortAddressLength:]
	if len(argdata)%CallDataLength != 0 {
		return nil, fmt.Errorf("invalid call data; lengsth should be a multiple of 32 bytes but was %d", len(argdata))
	}
	erc721Abi, err := abi.JSON(strings.NewReader(contract.EnumerableMetaData.ABI))
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
