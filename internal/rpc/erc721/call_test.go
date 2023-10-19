package erc721

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
)

const ExpectedData = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047777773000000000000000000000000000000000000000000000000000000000"

func TestProcessCallUnallowedMethod(t *testing.T) {
	mockEthClient := new(MockRPCClient)
	data := "0x12345678" // Example data
	to := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	_, err := ProcessCall(data, to, mockEthClient)
	if err == nil || err.Error() != "unallowed method: 0x12345678" {
		t.Errorf("Expected error: unallowed method: 0x12345678, got: %v", err)
	}
}

func TestProcessCall(t *testing.T) {
	mockClient := new(MockRPCClient)
	// Mock behavior & inject result
	mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ExpectedData, nil)

	data := "0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000" // Example data
	to := common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527")

	res, err := ProcessCall(data, to, mockClient)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Assert that the returned result matches the expected one
	if hexutil.Encode(res) != ExpectedData {
		t.Errorf("Expected result: %v, got: %v", ExpectedData, hexutil.Encode(res))
	}

	// Additional validation for the decoded data (if necessary)
	decodedData, err := hexutil.Decode(hexutil.Encode(res))
	if err != nil {
		t.Fatalf("Expected no error during decode, got: %v", err)
	}

	// Extract the actual string from the ABI-encoded data
	actualStr := string(decodedData[64 : 64+4]) // Assuming that "www0" is always 4 bytes long
	if actualStr != "www0" {
		t.Errorf("Expected string: www0, got: %v", actualStr)
	}
}
