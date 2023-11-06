package erc721_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
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
		input         erc721.CallData
		expected      erc721.Erc721method
		remoteMinting bool
		err           error
	}{
		{
			input:         hexutil.MustDecode("0x6352211e"),
			expected:      erc721.OwnerOf,
			remoteMinting: true,
			err:           nil,
		},
		{
			input:         hexutil.MustDecode("0x70a08231"),
			expected:      erc721.BalanceOf,
			remoteMinting: true,
			err:           nil,
		},
		{
			input:         hexutil.MustDecode("0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28"),
			expected:      erc721.BalanceOf,
			remoteMinting: true,
			err:           nil,
		},
		{
			input:         hexutil.MustDecode("0xc87b56dd"), //tokenUri
			expected:      0,
			remoteMinting: false,
			err:           nil,
		},
		{
			input:         hexutil.MustDecode("0x95d89b41"), //symbol
			expected:      0,
			remoteMinting: false,
			err:           nil,
		},
		{
			input:         hexutil.MustDecode("0x01ffc9a7"),
			expected:      erc721.SupportsInterface,
			remoteMinting: true,
			err:           nil,
		},
		{
			input:         erc721.CallData{0x00, 0x00, 0x00},
			expected:      0,
			remoteMinting: false,
			err:           fmt.Errorf("invalid call data, incomplete method signature (3 bytes < 4)"),
		},
		{
			input:         erc721.CallData{0x12, 0x34, 0x56, 0x78},
			expected:      0,
			remoteMinting: false,
			err:           nil,
		},
	}

	for _, test := range tests {
		output, exists, err := test.input.UniversalMintingMethod()
		if output != test.expected {
			t.Errorf("got: %v, Expected: %v", output, test.expected)
		}
		if exists != test.remoteMinting {
			t.Errorf("got: %v, Expected: %v", exists, test.remoteMinting)
		}
		if err == nil && test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("got: %v, Expected error: %v, ", err, test.err)
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
