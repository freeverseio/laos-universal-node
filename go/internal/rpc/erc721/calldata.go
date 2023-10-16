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
	// OwnerOf represents the 'ownerOf' ERC721 method.
	OwnerOf erc721method = iota
	// BalanceOf represents the 'balanceOf' ERC721 method.
	BalanceOf
	// TokenURI represents the 'tokenURI' ERC721 method.
	TokenURI
	// URI represents the 'uri' ERC721 method.
	URI
	// SupportsInterface represents the 'supportsInterface' ERC721 method.
	SupportsInterface
	// Name represents the 'name' ERC721 method.
	Name
	// ContractURI represents the 'contractURI' ERC721 method.
	ContractURI
	// Decimals represents the 'decimals' ERC721 method.
	Decimals
	// Symbol represents the 'symbol' ERC721 method.
	Symbol
)

var (
	ownerOfSig           = hexutil.Encode(crypto.Keccak256Hash([]byte("ownerOf(uint256)")).Bytes()[:4])
	balanceOfSig         = hexutil.Encode(crypto.Keccak256Hash([]byte("balanceOf(address)")).Bytes()[:4])
	tokenURISig          = hexutil.Encode(crypto.Keccak256Hash([]byte("tokenURI(uint256)")).Bytes()[:4])
	supportsInterfaceSig = hexutil.Encode(crypto.Keccak256Hash([]byte("supportsInterface(bytes4)")).Bytes()[:4])
	nameSig              = hexutil.Encode(crypto.Keccak256Hash([]byte("name()")).Bytes()[:4])
	contractURISig       = hexutil.Encode(crypto.Keccak256Hash([]byte("contractURI()")).Bytes()[:4])
	decimalsSig          = hexutil.Encode(crypto.Keccak256Hash([]byte("decimals()")).Bytes()[:4])
	symbolSig            = hexutil.Encode(crypto.Keccak256Hash([]byte("symbol()")).Bytes()[:4])
)

// Method returns the ERC721 method invoked by the calldata.
func (b CallData) Method() (erc721method, error) {
	sig, err := b.methodSignature()
	if err != nil {
		return 0, err
	}

	switch sig {
	case ownerOfSig:
		return OwnerOf, nil
	case balanceOfSig:
		return BalanceOf, nil
	case tokenURISig:
		return TokenURI, nil
	case supportsInterfaceSig:
		return SupportsInterface, nil
	case nameSig:
		return Name, nil
	case contractURISig:
		return ContractURI, nil
	case decimalsSig:
		return Decimals, nil
	case symbolSig:
		return Symbol, nil
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
