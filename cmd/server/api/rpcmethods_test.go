package api_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
)

func TestGetRPCMethod(t *testing.T) {
	tests := []struct {
		name          string
		methodName    string
		wantExists    bool
		wantRPCMethod api.RPCMethod
	}{
		{"SupportedMethodEthCall", "eth_call", true, api.RPCMethodEthCall},
		{"SupportedMethodEthGetBalance", "eth_getBalance", true, api.RPCMethodEthGetBalance},
		{"UnsupportedMethod", "eth_unsupportedMethod", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRPCMethod, gotExists := api.HasRPCMethodWithBlocknumber(tt.methodName)

			if gotExists != tt.wantExists {
				t.Errorf("getRPCMethod() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if gotExists && gotRPCMethod != tt.wantRPCMethod {
				t.Errorf("getRPCMethod() gotRPCMethod = %v, want %v", gotRPCMethod, tt.wantRPCMethod)
			}
		})
	}
}

func TestReplaceBlockTag(t *testing.T) {

	// Happy path
	t.Run("valid block tag", func(t *testing.T) {
		req := &api.JSONRPCRequest{
			Params: []json.RawMessage{
				json.RawMessage(`"latest"`),
			},
		}
		method := api.RPCMethodEthGetBlockByNumber
		blockNumber := "0x1b4"

		got, err := api.ReplaceBlockTag(req, method, blockNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		compareRawMessageString(t, got.Params[0], json.RawMessage(`"0x1b4"`))
	})

	t.Run("valid block number", func(t *testing.T) {
		req := &api.JSONRPCRequest{
			Params: []json.RawMessage{
				json.RawMessage(`"0x297e0c2"`),
			},
		}
		method := api.RPCMethodEthGetBlockByNumber
		blockNumber := "0x1b4"

		got, err := api.ReplaceBlockTag(req, method, blockNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		compareRawMessageString(t, got.Params[0], json.RawMessage(`"0x297e0c2"`))
	})

	t.Run("valid block tag for RPCMethodEthGetBalance", func(t *testing.T) {
		req := &api.JSONRPCRequest{
			Params: []json.RawMessage{
				json.RawMessage(`"0x407d73d8a49eeb85d32cf465507dd71d507100c1"`),
				json.RawMessage(`"latest"`),
			},
		}
		method := api.RPCMethodEthGetBalance
		blockNumber := "0x1b4"

		got, err := api.ReplaceBlockTag(req, method, blockNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		compareRawMessageString(t, got.Params[1], json.RawMessage(`"0x1b4"`))
	})
	t.Run("valid block tag for RPCMethodEthCall", func(t *testing.T) {
		req := &api.JSONRPCRequest{
			Params: []json.RawMessage{
				json.RawMessage(`{"to":"0x407d73d8a49eeb85d32cf465507dd71d507100c1","data":"0x0"}`),
				json.RawMessage(`"latest"`),
			},
		}
		method := api.RPCMethodEthCall
		blockNumber := "0x1b4"

		got, err := api.ReplaceBlockTag(req, method, blockNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		compareRawMessageString(t, got.Params[1], json.RawMessage(`"0x1b4"`))
	})

}

func compareRawMessageString(t *testing.T, raw1, raw2 json.RawMessage) {
	t.Helper()
	// Check if both are nil or both are not nil
	if (raw1 == nil) != (raw2 == nil) {
		t.Fatalf("One of the RawMessage is nil and the other is not. Got %v, expected %v", raw1, raw2)
	}

	// Compare the values if both are not nil
	if raw1 != nil && raw2 != nil {
		if !bytes.Equal(raw1, raw2) {
			t.Fatalf("RawMessage values are not equal. Got %v, expected %v", raw1, raw2)
		}
	}
}
