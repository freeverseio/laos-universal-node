package api_test

import (
	"bytes"
	"encoding/json"
	"reflect"
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

func TestReplaceBlockTagTT(t *testing.T) {
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
					json.RawMessage(`{
		        "fromBlock": "0x1",
		        "toBlock": "0x2",
		        "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
		        "topics": [
		          ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"],
		          null,
		          [
		            "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
		            "0x0000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc"
		          ]
		        ]
		      }`),
				},
			},
			method:      api.RPCMethodEthNewFilter,
			blockNumber: "0x1b4",
			expectedParam: json.RawMessage(`{
		    "fromBlock": "0x1",
		    "toBlock": "0x2",
		    "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
		    "topics": [
		      ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"],
		      null,
		      [
		        "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
		        "0x0000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc"
		      ]
		    ]
		  }`),
			expectError:       false,
			parameterPosition: 0,
		},
		{
			name: "valid block tag for eth_newFilter and toBlock is latest",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					json.RawMessage(`{
		        "fromBlock": "latest",
		        "toBlock": "latest",
		        "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
		        "topics": [
		          "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
		          null,
		          [
		            "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
		            "0x0000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc"
		          ]
		        ]
		      }`),
				},
			},
			method:      api.RPCMethodEthNewFilter,
			blockNumber: "0x1b4",
			expectedParam: json.RawMessage(`{
		    "fromBlock": "0x1b4",
		    "toBlock": "0x1b4",
		    "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
		    "topics": [
		      "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
		      null,
		      [
		        "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
		        "0x0000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc"
		      ]
		    ]
		  }`),
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
					json.RawMessage(`{
            "fromBlock": "latest",
            "toBlock": "latest",
            "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
            "topics": [
              "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"
            ]
          }`),
				},
			},
			method:      api.RPCMethodEthGetLogs,
			blockNumber: "0x1b4",
			expectedParam: json.RawMessage(`{
        "fromBlock": "0x1b4",
        "toBlock": "0x1b4",
        "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
        "topics": [
          "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"
        ]
      }`),
			expectError:       false,
			parameterPosition: 0,
		},
		{
			name: "valid block tag for eth_getLogs with latest and pending",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					json.RawMessage(`{
            "fromBlock": "latest",
            "toBlock": "pending",
            "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
            "topics": [
              "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"
            ]
          }`),
				},
			},
			method:      api.RPCMethodEthGetLogs,
			blockNumber: "0x1b4",
			expectedParam: json.RawMessage(`{
        "fromBlock": "0x1b4",
        "toBlock": "0x1b5",
        "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
        "topics": [
          "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"
        ]
      }`),
			expectError:       false,
			parameterPosition: 0,
		},
		{
			name: "invalid block tag for eth_getLogs (blocknumber too big)",
			req: &api.JSONRPCRequest{
				Params: []json.RawMessage{
					json.RawMessage(`{
            "fromBlock": "0x1b5",
            "toBlock": "0x1b5",
            "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
            "topics": [
              "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"
            ]
          }`),
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
		t.Run(tc.name, func(t *testing.T) {
			got, err := api.ReplaceBlockTag(tc.req, tc.method, tc.blockNumber)

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
				compareRawMessageObject(t, got.Params[tc.parameterPosition], tc.expectedParam)
			}
		})
	}
}

func TestReplaceBlockTagFromObject(t *testing.T) {
	t.Run("valid block tag for eth_newFilter", func(t *testing.T) {
		req := &api.JSONRPCRequest{
			Params: []json.RawMessage{
				json.RawMessage(`{
          "fromBlock": "0x1",
          "toBlock": "0x2",
          "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
          "topics": [
            "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
            null,
            [
              "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
              "0x0000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc"
            ]
          ]
        }`),
			},
		}
		blockNumber := "0x1b4"

		err := api.ReplaceBlockTagFromObject(req, blockNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var filterObject api.FilterObject
		err = json.Unmarshal(req.Params[0], &filterObject)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if filterObject.FromBlock != "0x1" {
			t.Fatalf("unexpected fromBlock: %v", filterObject.FromBlock)
		}
		if filterObject.ToBlock != "0x2" {
			t.Fatalf("unexpected toBlock: %v", filterObject.ToBlock)
		}

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
			t.Fatalf("RawMessage values are not equal. Got %v, expected %v", string(raw1), string(raw2))
		}
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
