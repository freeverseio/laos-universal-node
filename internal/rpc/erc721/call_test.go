package erc721

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessCallUnallowedMethod(t *testing.T) {

	mockEthClient := new(MockRPCClient)
	data := "0x12345678" // Example data
	to := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	contractAddr := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdef")
	chainID := uint64(1)

	_, err := ProcessCall(data, to, mockEthClient, contractAddr, chainID)
	assert.Error(t, err)
	if assert.Error(t, err) {
		assert.Equal(t, "unallowed method: 0x12345678", err.Error())
	}
}

func TestProcessCall(t *testing.T) {
	mockClient := new(MockRPCClient)
	// Mock behavior & inject result
	expectedResult := "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047777773000000000000000000000000000000000000000000000000000000000"
	mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, nil)

	data := "0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000" // Example data
	to := common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527")
	contractAddr := common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527")
	chainID := uint64(1)

	res, err := ProcessCall(data, to, mockClient, contractAddr, chainID)
	assert.NoError(t, err)

	// Assert that the returned result matches the expected one
	assert.Equal(t, expectedResult, hexutil.Encode(res))

	// Additional validation for the decoded data (if necessary)
	decodedData, err := hexutil.Decode(hexutil.Encode(res))
	assert.NoError(t, err)
	// Extract the actual string from the ABI-encoded data
	actualStr := string(decodedData[64 : 64+4]) // Assuming that "www0" is always 4 bytes long

	assert.Equal(t, "www0", actualStr)

}
