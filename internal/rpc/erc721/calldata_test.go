package erc721

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestNewCalldata(t *testing.T) {
	tests := []struct {
		input         string
		expected      CallData
		expectedError string
	}{
		{
			input:         "0x1234",
			expected:      CallData{0x12, 0x34},
			expectedError: "",
		},
		{
			input:         "invalid",
			expected:      CallData{},
			expectedError: "hex string without 0x prefix",
		},
	}

	for _, test := range tests {
		output, err := NewCallData(test.input)
		assert.Equal(t, test.expected, output)
		if test.expectedError == "" {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), test.expectedError)
		}
	}
}

func TestMethod(t *testing.T) {
	tests := []struct {
		input    CallData
		expected erc721method
		err      error
	}{
		{
			input:    hexutil.MustDecode("0x6352211e"),
			expected: OwnerOf,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0x70a08231"),
			expected: BalanceOf,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0xc87b56dd"),
			expected: TokenURI,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0x01ffc9a7"),
			expected: SupportsInterface,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0x06fdde03"),
			expected: Name,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0x313ce567"),
			expected: Decimals,
			err:      nil,
		},
		{
			input:    CallData{0x00, 0x00, 0x00},
			expected: 0,
			err:      fmt.Errorf("invalid call data, incomplete method signature (3 bytes < 4)"),
		},
		{
			input:    CallData{0x12, 0x34, 0x56, 0x78},
			expected: 0,
			err:      fmt.Errorf("unallowed method: 0x12345678"),
		},
	}

	for _, test := range tests {
		output, err := test.input.Method()
		assert.Equal(t, test.expected, output)
		assert.Equal(t, test.err, err)
	}
}

func TestGetParam(t *testing.T) {
	// This is a bit more complex since it requires correct ABI encoding.
	// For the sake of example, I'll provide a simple framework.
	tests := []struct {
		input    CallData
		param    string
		expected interface{}
		err      error
	}{
		{
			input:    hexutil.MustDecode("0x70a08231000000000000000000000000bd7931f025ecf360b21e1ab92ec34b49084bca5b"),
			param:    "owner",
			expected: common.HexToAddress("0xbD7931f025ecF360b21E1aB92ec34b49084bcA5B"),
			err:      nil,
		},
	}

	for _, test := range tests {
		output, err := test.input.GetParam(test.param)
		assert.Equal(t, test.expected, output)
		assert.Equal(t, test.err, err)
	}
}
