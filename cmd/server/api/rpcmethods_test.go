package api_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
)

func TestGetRPCMethod(t *testing.T) {
	t.Parallel()
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
		tt := tt // Shadow loop variable otherwise it could be overwrittens
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			methodManager := api.NewProxyRPCMethodManager()
			gotRPCMethod, gotExists := methodManager.HasRPCMethodWithBlockNumber(tt.methodName)

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
	t.Parallel()
	type test struct {
		name              string
		req               *api.JSONRPCRequest
		method            api.RPCMethod
		blockNumber       string
		expectedParam     json.RawMessage
		expectError       bool
		expectedError     string
		parameterPosition int
	}

	tests := []test{
		{
			name: "valid block tag",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{json.RawMessage(`"latest"`)},
			},
			method:        api.RPCMethodEthGetBlockByNumber,
			blockNumber:   "0x1b4",
			expectedParam: json.RawMessage(`"0x1b4"`),
			expectError:   false,
		},
		{
			name: "valid block number",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{json.RawMessage(`"0x297e0c2"`)},
			},
			method:        api.RPCMethodEthGetBlockByNumber,
			blockNumber:   "0x297e0c2",
			expectedParam: json.RawMessage(`"0x297e0c2"`),
			expectError:   false,
		},
		{
			name: "invalid block number",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{json.RawMessage(`"0x297e0c3"`)},
			},
			method:        api.RPCMethodEthGetBlockByNumber,
			blockNumber:   "0x297e0c2",
			expectError:   true,
			expectedError: "invalid block number: 0x297e0c3",
		},
		{
			name: "valid block tag for RPCMethodEthGetBalance",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					json.RawMessage(`"0x407d73d8a49eeb85d32cf465507dd71d507100c1"`),
					json.RawMessage(`"latest"`),
				},
			},
			method:            api.RPCMethodEthGetBalance,
			blockNumber:       "0x1b4",
			expectedParam:     json.RawMessage(`"0x1b4"`),
			expectError:       false,
			parameterPosition: 1,
		},
		{
			name: "valid block tag for RPCMethodEthCall",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					json.RawMessage(`{"to":"0x407d73d8a49eeb85d32cf465507dd71d507100c1","data":"0x0"}`),
					json.RawMessage(`"latest"`),
				},
			},
			method:            api.RPCMethodEthCall,
			blockNumber:       "0x1b4",
			expectedParam:     json.RawMessage(`"0x1b4"`),
			expectError:       false,
			parameterPosition: 1,
		},
		{
			name: "valid block tag for eth_newFilter",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					mock.GetFilterRequest("0x1", "0x2"),
				},
			},
			method:            api.RPCMethodEthNewFilter,
			blockNumber:       "0x1b4",
			expectedParam:     mock.GetFilterRequest("0x1", "0x2"),
			expectError:       false,
			parameterPosition: 0,
		},
		{
			name: "valid block tag for eth_newFilter and toBlock is latest",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					mock.GetFilterRequest("latest", "latest"),
				},
			},
			method:            api.RPCMethodEthNewFilter,
			blockNumber:       "0x1b4",
			expectedParam:     mock.GetFilterRequest("0x1b4", "0x1b4"),
			expectError:       false,
			parameterPosition: 0,
		},
		{
			name: "valid block tag for eth_getLogs with only topics",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					json.RawMessage(`{
            "topics": [
              "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"
            ]
          }`),
				},
			},
			method:      api.RPCMethodEthGetLogs,
			blockNumber: "0x1b4",
			expectedParam: json.RawMessage(`{
        "topics": [
          "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"
        ]
      }`),
			expectError:       false,
			parameterPosition: 0,
		},
		{
			name: "valid block tag for eth_getLogs with latest",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					mock.GetLogsRequest("latest", "latest"),
				},
			},
			method:            api.RPCMethodEthGetLogs,
			blockNumber:       "0x1b4",
			expectedParam:     mock.GetLogsRequest("0x1b4", "0x1b4"),
			expectError:       false,
			parameterPosition: 0,
		},
		{
			name: "invalid block tag for eth_getLogs (blocknumber too big)",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					mock.GetLogsRequest("0x1b5", "0x1b5"),
				},
			},
			method:            api.RPCMethodEthGetLogs,
			blockNumber:       "0x1b4",
			expectError:       true,
			expectedError:     "invalid block number: 0x1b5",
			parameterPosition: 0,
		},
	}

	for _, tc := range tests {
		tc := tc // Shadow loop variable otherwise it could be overwrittens
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			methodsManager := api.NewProxyRPCMethodManager()
			err := methodsManager.ReplaceBlockTag(tc.req, tc.method, tc.blockNumber)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if err.Error() != tc.expectedError {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				compareRawMessageObject(t, tc.req.Params[tc.parameterPosition], tc.expectedParam)
			}
		})
	}
}

func TestCheckBlockNumberFromResponseFromHashCalls(t *testing.T) {
	t.Parallel()
	type test struct {
		name               string
		method             api.RPCMethod
		blockNumber        string
		response           *json.RawMessage
		expectedBlockError string
	}

	tests := []test{
		{
			name:               "eth_getBlockByHash with correct block number",
			method:             api.RPCMethodEthGetBlockByHash,
			blockNumber:        "0x29b8ef5",
			response:           &mock.MockResponseBlock,
			expectedBlockError: "",
		},
		{
			name:               "eth_getBlockByHash with correct block number",
			method:             api.RPCMethodEthGetBlockByHash,
			blockNumber:        "0x29b8ef4",
			response:           &mock.MockResponseBlock,
			expectedBlockError: "invalid block number: 0x29b8ef5",
		},
		{
			name:               "eth_getBlockByHash with inexistent block number",
			method:             api.RPCMethodEthGetBlockByHash,
			blockNumber:        "0x29b8ef4",
			response:           nil,
			expectedBlockError: "",
		},
		{
			name:               "eth_getTransactionByHash with correct block number",
			method:             api.RPCMethodEthGetTransactionByHash,
			blockNumber:        "0x29b8ef5",
			response:           &mock.MockResponseTransaction,
			expectedBlockError: "",
		},
		{
			name:               "RPCMethodEthGetTransactionReceipt with correct block number",
			method:             api.RPCMethodEthGetTransactionReceipt,
			blockNumber:        "0x29b8ef5",
			response:           &mock.MockResponseTransaction,
			expectedBlockError: "",
		},
		{
			name:               "RPCMethodEthGetTransactionByBlockHashAndIndex with correct block number",
			method:             api.RPCMethodEthGetTransactionByBlockHashAndIndex,
			blockNumber:        "0x29b8ef5",
			response:           &mock.MockResponseTransaction,
			expectedBlockError: "",
		},
		{
			name:               "eth_getTransactionByHash with wrong block number",
			method:             api.RPCMethodEthGetTransactionByHash,
			blockNumber:        "0x29b8ef4",
			response:           &mock.MockResponseTransaction,
			expectedBlockError: "invalid block number: 0x29b8ef5",
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rpcResponse := api.RPCResponse{
				Result: tt.response,
			}
			methodsManager := api.NewProxyRPCMethodManager()
			err := methodsManager.CheckBlockNumberFromResponseFromHashCalls(&rpcResponse, tt.method, tt.blockNumber)
			if err != nil {
				if err.Error() != tt.expectedBlockError {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func compareRawMessageObject(t *testing.T, raw1, raw2 json.RawMessage) {
	t.Helper()
	var obj1, obj2 interface{}

	err := json.Unmarshal(raw1, &obj1)
	if err != nil {
		t.Fatalf("Error unmarshaling raw1: %v", err)
	}

	err = json.Unmarshal(raw2, &obj2)
	if err != nil {
		t.Fatalf("Error unmarshaling raw2: %v", err)
	}

	if !reflect.DeepEqual(obj1, obj2) {
		t.Fatalf("JSON objects are not equal. Got %+v, expected %+v", obj1, obj2)
	}
}
