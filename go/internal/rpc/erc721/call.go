// Package erc721 provides a Go implementation of the ERC721 token standard for non-fungible tokens (NFTs)
// on the Ethereum blockchain.
package erc721

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/blockchain"
)

// ProcessCall processes an ERC721 token contract call.
func ProcessCall( 
	data string, 
	to common.Address,
	ethcli	blockchain.EthClient,	
	contractAddr common.Address,	
	chainID uint64) (hexutil.Bytes, error) {

	callData, err := NewCallData(data)
	if err != nil {
		return nil, err
	}
	method, err := callData.Method()
	log.Println(method)
	if err != nil {
		return nil, err
	}

	var result string
	err = ethcli.Client().Call(&result, "eth_call", map[string]interface{}{
		"to":   to,
		"data": data,
	}, "latest")
	if err != nil {
		return nil, err
	}

	switch method {
	case OwnerOf:
		return common.HexToAddress(result).Bytes(), nil

	case TokenURI, URI:
		return tokenURI(result)

	case BalanceOf:
		return balanceOf(callData)

	case SupportsInterface:
		return supportsInterface(callData)

	case Name:
		return name()

	case ContractURI:
		return contractURI()

	case Decimals:
		return decimals()

	case Symbol:
		return symbol()
	}

	return nil, fmt.Errorf("no method found processing erc721 calldata: %s", data)
}

// name is an internal function that retrieves the name of the ERC721 token.
func name() ([]byte, error) {
	stringTy, _ := abi.NewType("string", "string", nil)
	arguments := abi.Arguments{
		{
			Type: stringTy,
		},
	}
	return arguments.Pack("Living Assets")
}

// symbol is an internal function that retrieves the symbol of the ERC721 token.
func symbol() ([]byte, error) {
	stringTy, _ := abi.NewType("string", "string", nil)
	arguments := abi.Arguments{
		{
			Type: stringTy,
		},
	}

	return arguments.Pack("LA")
}

// contractURI is an internal function that retrieves the contract URI of the ERC721 token.
func contractURI() ([]byte, error) {
	stringTy, _ := abi.NewType("string", "string", nil)
	arguments := abi.Arguments{
		{
			Type: stringTy,
		},
	}
	return arguments.Pack("https://livingassets.io/contractUri")
}

// tokenURI retrieves the token URI for a given token ID.
func tokenURI(u string) ([]byte, error) {
	stringTy, _ := abi.NewType("string", "string", nil)
	arguments := abi.Arguments{
		{
			Type: stringTy,
		},
	}
	t, err := convertHexStringToText(u)
	if err != nil {
		return nil, err
	}
	return arguments.Pack(string(t))
}

// supportsInterface checks if the contract supports a given interface according to the ERC721 standard.
func supportsInterface(calldata CallData) ([]byte, error) {
	interfaceIDI, err := calldata.GetParam("interfaceId")
	if err != nil {
		return nil, err
	}
	interfaceID, ok := interfaceIDI.([4]uint8)
	if !ok {
		return nil, fmt.Errorf("invalid interfaceId %s", interfaceIDI)
	}
	supports := false
	if hexutil.Encode(interfaceID[:]) == "0x5b5e139f" {
		supports = true
	}
	boolTy, _ := abi.NewType("bool", "bool", nil)
	arguments := abi.Arguments{
		{
			Type: boolTy,
		},
	}
	return arguments.Pack(supports)
}

// balanceOf retrieves the balance of a given address according to the ERC721 standard.
func balanceOf( _ CallData) ([]byte, error) {
	log.Println("balanceOf..")
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	arguments := abi.Arguments{
		{
			Type: uint256Ty,
		},
	}
	return arguments.Pack(big.NewInt(1000000))
}

// decimals packs a uint8 value of 0 into a []byte slice using the Ethereum ABI encoding.
func decimals() ([]byte, error) {
	uint8Ty, _ := abi.NewType("uint8", "uint8", nil)
	arguments := abi.Arguments{
		{
			Type: uint8Ty,
		},
	}
	return arguments.Pack(uint8(0))
}



func convertHexStringToText(hexStr string) ([]byte, error) {
	// Remove the "0x" prefix if it exists.
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	// Extract the relevant portion based on the ABI encoding.
	dataStart := 64 // 32 bytes offset * 2 (since 1 byte is 2 hex characters)
	dataLengthHex := hexStr[dataStart : dataStart+64]
	dataLength, err := hex.DecodeString(dataLengthHex)
	if err != nil {
		return nil, err
	}
	length := int(dataLength[31]) // Last byte is the length

	dataHex := hexStr[dataStart+64 : dataStart+64+length*2]
	dataBytes, err := hex.DecodeString(dataHex)
	if err != nil {
		return nil, err
	}
	log.Println(string(dataBytes))
	return dataBytes, nil
}
