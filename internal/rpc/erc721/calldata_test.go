package erc721_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/rpc/erc721"
)

func TestNewCalldata(t *testing.T) {
	tests := []struct {
		input         string
		expected      erc721.CallData
		expectedError string
	}{
		{
			input:         "0x1234",
			expected:      erc721.CallData{0x12, 0x34},
			expectedError: "",
		},
		{
			input:         "invalid",
			expected:      erc721.CallData{},
			expectedError: "hex string without 0x prefix",
		},
	}

	for _, test := range tests {
		output, err := erc721.NewCallData(test.input)

		if !slicesEqual(output, test.expected) {
			t.Errorf("Expected: %v, got: %v", test.expected, output)
		}

		if test.expectedError == "" {
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		} else {
			if err == nil || err.Error() != test.expectedError {
				t.Errorf("Expected error: %v, got: %v", test.expectedError, err)
			}
		}
	}
}

func TestMethod(t *testing.T) {
	tests := []struct {
		input    erc721.CallData
		expected erc721.Erc721method
		err      error
	}{
		{
			input:    hexutil.MustDecode("0x6352211e"),
			expected: erc721.OwnerOf,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0x70a08231"),
			expected: erc721.BalanceOf,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0xc87b56dd"),
			expected: erc721.TokenURI,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0x01ffc9a7"),
			expected: erc721.SupportsInterface,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0x06fdde03"),
			expected: erc721.Name,
			err:      nil,
		},
		{
			input:    hexutil.MustDecode("0x313ce567"),
			expected: erc721.Decimals,
			err:      nil,
		},
		{
			input:    erc721.CallData{0x00, 0x00, 0x00},
			expected: 0,
			err:      fmt.Errorf("invalid call data, incomplete method signature (3 bytes < 4)"),
		},
		{
			input:    erc721.CallData{0x12, 0x34, 0x56, 0x78},
			expected: 0,
			err:      fmt.Errorf("unallowed method: 0x12345678"),
		},
	}

	for _, test := range tests {
		output, err := test.input.Method()
		if output != test.expected {
			t.Errorf("Expected: %v, got: %v", test.expected, output)
		}
		if err == nil && test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("Expected error: %v, got: %v", test.err, err)
			}
		}
	}
}

func TestGetParam(t *testing.T) {
	// This is a bit more complex since it requires correct ABI encoding.
	// For the sake of example, I'll provide a simple framework.
	tests := []struct {
		input    erc721.CallData
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
		if output != test.expected {
			t.Errorf("Expected: %v, got: %v", test.expected, output)
		}
		if err != test.err {
			t.Errorf("Expected error: %v, got: %v", test.err, err)
		}
	}
}

func slicesEqual(a, b erc721.CallData) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
