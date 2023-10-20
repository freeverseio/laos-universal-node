package erc721

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/stretchr/testify/mock"
)

type MockRPCClient struct {
	mock.Mock
}

func (m *MockRPCClient) Call(result interface{}, method string, args ...interface{}) error {
	argsIn := []interface{}{result, method}
	argsIn = append(argsIn, args...)
	retValues := m.Called(argsIn...)
	// Check if a result is provided, if so, set the result value
	if retValues.Get(0) != nil {
		switch v := result.(type) {
		case *map[string]interface{}:
			*v = retValues.Get(0).(map[string]interface{})
		case *string:
			*v = retValues.Get(0).(string)
		// Add more types if needed
		default:
			return fmt.Errorf("unsupported type: %T", v)
		}
	}
	return retValues.Error(1)
}

// Test for ChainId method
func TestChainId(t *testing.T) {
	mockClient := new(MockRPCClient)
	// Mock behavior & inject result
	mockClient.On("Call", mock.Anything, "eth_chainId").Return("0x13881", nil)
	service := EthService{
		Ethcli: mockClient,
	}
	chainId := service.ChainId()
	expectedChainId := (*hexutil.Big)(big.NewInt(80001))
	if chainId.ToInt().Cmp(expectedChainId.ToInt()) != 0 {
		t.Errorf("Expected chain ID to be %v but got %v", expectedChainId, chainId)
	}
}

// Test for BlockNumber method
func TestBlockNumber(t *testing.T) {
	mockClient := new(MockRPCClient)
	// Mock behavior & inject result
	mockClient.On("Call", mock.Anything, "eth_blockNumber").Return("0x277f60e", nil)
	service := EthService{
		Ethcli: mockClient,
	}
	blockNumber, err := service.BlockNumber(context.Background())
	if blockNumber != hexutil.Uint64(41416206) {
		t.Errorf("Expected block number to be 0 but got %v", blockNumber)
	}
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

// Test for GetBlockByNumber method
func TestGetBlockByNumber(t *testing.T) {
	mockClient := new(MockRPCClient)
	mockResult := map[string]interface{}{
		"baseFeePerGas":   "0x10",
		"difficulty":      "0x7",
		"extraData":       "0xd78301000683626f7288676f312e32302e35856c696e757800000000000000006116f632ba8263ce3d3e570e3689d6fbfe4f40c9aa9b8b05b2be3a0c3da6c6b74a8d915e78970f28dd7653ee2b1c3ad3e8dd86b0a1afd077beab9f83690b8ef900",
		"gasLimit":        "0x1c2324b",
		"gasUsed":         "0x586e3",
		"hash":            "0x1486afd8523cb57dcf6d11659dddf3f2123618614d12ecb015a104c30ef4ef20",
		"logsBloom":       "0x000000000420000000000100002000000000000800000000000100802200012000000200060000000400801200000030000080000000040010000000002000000000000000000100000200081000008000000000008010040041002000008800000800080200000200000004004008000000000000c8000880400810000000000004080000000000000000240000000000000000000000020004000200000000220020400000100000000000000000000000000002000000000100080000000004000000002000000000001001000021008000000000000808000108008000020000210000000000080000000000040000020040002004000000080000800100004",
		"miner":           "0x0000000000000000000000000000000000000000",
		"mixHash":         "0x0000000000000000000000000000000000000000000000000000000000000000",
		"nonce":           "0x0000000000000000",
		"number":          "0x277f60e",
		"parentHash":      "0x72e288bbacd6c55fd28ad660e26265a414c68ae805db408907e45cbd7ebac0a8",
		"receiptsRoot":    "0x0a7d86ed4c6c6b4e47645cb600f1e013896783186c1d528dbd3aa34258881ab5",
		"sha3Uncles":      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"size":            "0x76f",
		"stateRoot":       "0x823a4c5b853eac0208d83ec4d31b6464fd8eb7d674baefce658713d25df9d5b1",
		"timestamp":       "0x653218f4",
		"totalDifficulty": "0xf19a0fd",
		"transactions": []string{
			"0x17ccefdbb21e7b93e621cf975f95a939b22c9e659057deb881f731dcb9b2f82c",
		},
		"transactionsRoot": "0x8b3c3bd57474ac3c60ced349928994ea30930519968478eac592268f024caf9b",
		"uncles":           []interface{}{},
	}
	// Mock behavior & inject result
	mockClient.On("Call", mock.Anything, "eth_getBlockByNumber", "0x123", true).Return(mockResult, nil)
	service := EthService{
		Ethcli: mockClient,
	}
	block, err := service.GetBlockByNumber("0x123", true)
	if block == nil {
		t.Errorf("Expected block not to be nil")
	}
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if block == nil {
		t.Errorf("Expected block not to be nil")
	}
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	baseFeePerGas, ok := block["baseFeePerGas"].(string)
	if !ok || baseFeePerGas != "0x10" {
		t.Errorf("Expected baseFeePerGas to be '0x10', got %v", baseFeePerGas)
	}
	difficulty, ok := block["difficulty"].(string)
	if !ok || difficulty != "0x7" {
		t.Errorf("Expected difficulty to be '0x7', got %v", difficulty)
	}
}

// Test for Call method
func TestCall(t *testing.T) {
	t.Run("Could execute Call TokenURI without an error", func(t *testing.T) {
		mockClient := new(MockRPCClient)
		// Mock behavior & inject result
		mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ExpectedData, nil)

		service := EthService{
			Ethcli: mockClient,
		}
		// Define the test transaction
		tx := blockchain.Transaction{
			To:   "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
			Data: "0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000",
		}

		// Call the Call method
		res, err := service.Call(tx, "1")
		if err != nil {
			t.Fatalf("Expected no error but got %v", err)
		}
		if res.String() != ExpectedData {
			t.Errorf("Expected data to be %v but got %v", ExpectedData, res.String())
		}
	})

	t.Run("Could execute Call OwnerOf without an error", func(t *testing.T) {
		mockClient := new(MockRPCClient)
		// Mock behavior & inject result
		expectedResult := "0x1b0b4a597c764400ea157ab84358c8788a89cd28"
		mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, nil)

		service := EthService{
			Ethcli: mockClient,
		}
		// Define the test transaction
		tx := blockchain.Transaction{
			To:   "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
			Data: "0x6352211e0000000000000000000000000000000000000000000000000000000000000000",
		}

		// Call the Call method
		res, err := service.Call(tx, "1")
		if err != nil {
			t.Fatalf("Expected no error but got %v", err)
		}
		if res.String() != expectedResult {
			t.Errorf("Expected result to be %v but got %v", expectedResult, res.String())
		}
	})

	t.Run("Could execute Call BalanceOf without an error", func(t *testing.T) {
		//  0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28

		mockClient := new(MockRPCClient)
		// Mock behavior & inject result
		expectedResult := "0x0000000000000000000000000000000000000000000000000000000000000001"
		mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, nil)

		service := EthService{
			Ethcli: mockClient,
		}
		// Define the test transaction
		tx := blockchain.Transaction{
			To:   "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
			Data: "0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28",
		}

		// Call the Call method
		res, err := service.Call(tx, "1")
		if err != nil {
			t.Fatalf("Expected no error but got %v", err)
		}
		if res.String() != expectedResult {
			t.Errorf("Expected result to be %v but got %v", expectedResult, res.String())
		}
	})

	t.Run("Could execute Call TokenURI with an error", func(t *testing.T) {
		mockClient := new(MockRPCClient)
		// Mock behavior & inject result
		mockClient.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ExpectedData, fmt.Errorf("error from call"))

		service := EthService{
			Ethcli: mockClient,
		}
		// Define the test transaction
		tx := blockchain.Transaction{
			To:   "0xc4d9faef49ec1e604a76ee78bc992abadaa29527",
			Data: "0xc87b56dd0000000000000000000000000000000000000000000000000000000000000000",
		}

		// Call the Call method
		_, err := service.Call(tx, "1")
		if err == nil || err.Error() != "error from call" {
			t.Fatalf("Expected error 'error from call' but got %v", err)
		}
	})
}

func TestGetBalance(t *testing.T) {
	mockClient := new(MockRPCClient)
	// Mock behavior & inject result
	mockClient.On("Call", mock.Anything, "eth_getBalance", common.HexToAddress("0x1B0b4a597C764400Ea157aB84358c8788A89cd28"), "latest").Return("0x5cec30275aa9343c", nil)
	service := EthService{
		Ethcli: mockClient,
	}
	balance, err := service.GetBalance(common.HexToAddress("0x1B0b4a597C764400Ea157aB84358c8788A89cd28"), "latest")
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if balance != hexutil.Uint64(6695779691575981116) {
		t.Errorf("Expected block number to be 0 but got %v", balance)
	}
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}
