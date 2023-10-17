package erc721

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/blockchain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Implementing the methods of blockchain.EthClient with mock methods

// ... Add other mocked methods as needed

// Test for ChainId method
func TestChainId(t *testing.T) {
	service := EthService{
		ChainID: 42, // Sample chain ID
	}
	chainId := service.ChainId()
	expectedChainId := (*hexutil.Big)(big.NewInt(42))
	assert.Equal(t, expectedChainId, chainId)
}

// Test for BlockNumber method
func TestBlockNumber(t *testing.T) {
	service := EthService{}
	blockNumber, err := service.BlockNumber(context.Background())
	assert.Equal(t, hexutil.Uint64(0), blockNumber)
	assert.NoError(t, err)
}

// Test for GetBlockByNumber method
func TestGetBlockByNumber(t *testing.T) {
	service := EthService{}
	block, err := service.GetBlockByNumber("0x123", true)
	assert.NotNil(t, block)
	assert.NoError(t, err)
}

type MockRPCClient struct {
	mock.Mock
}

func (m *MockRPCClient) Call(result interface{}, method string, args ...interface{}) error {
	argsIn := []interface{}{result, method}
	argsIn = append(argsIn, args...)
	retValues := m.Called(argsIn...)
	// Check if a result is provided, if so, set the result value
	if retValues.Get(0) != nil {
		*result.(*string) = retValues.Get(0).(string)
	}
	return retValues.Error(1)
}

// Test for Call method
func TestCall(t *testing.T) {
	t.Run("Could execute Call TokenURI without an error", func(t *testing.T) {
		mockClient := new(MockRPCClient)
		// Mock behavior & inject result
		expectedResult := "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047777773000000000000000000000000000000000000000000000000000000000"
		mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, nil)

		service := EthService{
			Ethcli:       mockClient,
			ContractAddr: common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527"),
		}
		// Define the test transaction
		tx := blockchain.Transaction{
			To:   "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
			Data: "0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000",
		}

		// Call the Call method
		res, err := service.Call(tx, "1")
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, res.String())
	})

	t.Run("Could execute Call OwnerOf without an error", func(t *testing.T) {
		mockClient := new(MockRPCClient)
		// Mock behavior & inject result
		expectedResult := "0x1b0b4a597c764400ea157ab84358c8788a89cd28"
		mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, nil)

		service := EthService{
			Ethcli:       mockClient,
			ContractAddr: common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527"),
		}
		// Define the test transaction
		tx := blockchain.Transaction{
			To:   "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
			Data: "0x6352211e0000000000000000000000000000000000000000000000000000000000000000",
		}

		// Call the Call method
		res, err := service.Call(tx, "1")
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, res.String())
	})

	t.Run("Could execute Call BalanceOf without an error", func(t *testing.T) {
		//  0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28

		mockClient := new(MockRPCClient)
		// Mock behavior & inject result
		expectedResult := "0x0000000000000000000000000000000000000000000000000000000000000001"
		mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, nil)

		service := EthService{
			Ethcli:       mockClient,
			ContractAddr: common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527"),
		}
		// Define the test transaction
		tx := blockchain.Transaction{
			To:   "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
			Data: "0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28",
		}

		// Call the Call method
		res, err := service.Call(tx, "1")
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, res.String())
	})

	t.Run("Could execute Call TokenURI with an error", func(t *testing.T) {
		mockClient := new(MockRPCClient)
		// Mock behavior & inject result
		expectedResult := "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047777773000000000000000000000000000000000000000000000000000000000"
		mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, fmt.Errorf("error from call"))

		service := EthService{
			Ethcli:       mockClient,
			ContractAddr: common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527"),
		}
		// Define the test transaction
		tx := blockchain.Transaction{
			To:   "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
			Data: "0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000",
		}

		// Call the Call method
		_, err := service.Call(tx, "1")
		assert.Error(t, err, "error from call")

	})

}
