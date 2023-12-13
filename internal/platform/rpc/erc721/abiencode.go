package erc721

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func AbiEncodeString(value string) (string, error) {
	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return "", err
	}
	return abiEncodeValue(value, &stringType)
}

func abiEncodeValue(value any, abiType *abi.Type) (string, error) {
	abiArguments := abi.Arguments{
		{
			Type: *abiType,
			Name: "return",
		},
	}
	abiEncodedValue, err := abiArguments.Pack(value)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(abiEncodedValue), nil
}
