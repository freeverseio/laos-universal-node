// Package erc721 provides a Go implementation of the ERC721 token standard for non-fungible tokens (NFTs)
// on the Ethereum blockchain.
package erc721

import (
	"encoding/hex"
	"errors"
	"fmt"
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
	ethcli blockchain.EthRPCClient,
	contractAddr common.Address,
	chainID uint64,
) (hexutil.Bytes, error) {
	callData, err := NewCallData(data)
	if err != nil {
		return nil, err
	}
	method, err := callData.Method()
	if err != nil {
		return nil, err
	}

	var result string
	err = ethcli.Call(&result, "eth_call", map[string]interface{}{
		"to":   to,
		"data": data,
	}, "latest")

	if err != nil {
		return nil, err
	}

	switch method {
	case OwnerOf:
		return ownerOf(result)

	case TokenURI, URI:
		return tokenURI(result)

	case BalanceOf:
		return balanceOf(result)

	case SupportsInterface:
		return supportsInterface(result)

	case Name:
		return name(result)

	case ContractURI:
		return contractURI(result)

	case Decimals:
		return decimals(result)

	case Symbol:
		return symbol(result)
	}

	return nil, fmt.Errorf("no method found processing erc721 calldata: %s", data)
}

func ownerOf(result string) ([]byte, error) {
	if !common.IsHexAddress(result) {
		return nil, errors.New("not a valid Ethereum address")
	}
	return common.HexToAddress(result).Bytes(), nil
}

func name(result string) ([]byte, error) {
	return nil, fmt.Errorf("not yet implementd")
}

func symbol(result string) ([]byte, error) {
	return nil, fmt.Errorf("not yet implementd")
}

func contractURI(result string) ([]byte, error) {
	stringTy, _ := abi.NewType("string", "string", nil)
	arguments := abi.Arguments{
		{
			Type: stringTy,
		},
	}
	return arguments.Pack(result)
}

func tokenURI(u string) ([]byte, error) {
	stringTy, _ := abi.NewType("string", "string", nil)
	arguments := abi.Arguments{
		{
			Type: stringTy,
		},
	}
	t, err := convert64HexStringToText(u)
	if err != nil {
		return nil, err
	}
	return arguments.Pack(t)
}

func supportsInterface(result string) ([]byte, error) {
	return nil, fmt.Errorf("not yet implementd")
}

func balanceOf(b string) ([]byte, error) {
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	arguments := abi.Arguments{
		{
			Type: uint256Ty,
		},
	}
	return arguments.Pack(convert16HexStringToDecimal(b))
}

func decimals(result string) ([]byte, error) {
	return nil, fmt.Errorf("not yet implementd")
}

func convert64HexStringToText(hexStr string) (string, error) {
	// Remove the "0x" prefix if it exists.
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	dataStart := 64 // 32 bytes offset * 2 (since 1 byte is 2 hex characters)
	dataLengthHex := hexStr[dataStart : dataStart+64]
	dataLength, err := hex.DecodeString(dataLengthHex)
	if err != nil {
		return "", err
	}
	length := int(dataLength[31]) // Last byte is the length
	dataHex := hexStr[dataStart+64 : dataStart+64+length*2]
	dataBytes, err := hex.DecodeString(dataHex)
	if err != nil {
		return "", err
	}
	return string(dataBytes), nil
}

func convert16HexStringToDecimal(hexString string) *big.Int {
	// Create a new Int from a hexadecimal string
	num := new(big.Int)
	num.SetString(hexString[2:], 16) // Omit the "0x" prefix and use base 16
	return num
}
