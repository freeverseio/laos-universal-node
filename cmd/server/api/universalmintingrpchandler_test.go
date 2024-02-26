package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	"github.com/freeverseio/laos-universal-node/cmd/server/api/mock"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	mockTx "github.com/freeverseio/laos-universal-node/internal/platform/state/mock"
	"go.uber.org/mock/gomock"
)

const (
	hexStringOne  = "0x0000000000000000000000000000000000000000000000000000000000000001"
	hexStringZero = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

func TestUniversalMintingRPCHandlerTableTests(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name       string
		setupMocks func(*mockTx.MockService, *mockTx.MockTx, *mock.MockHTTPClientInterface, *mock.MockRPCMethodManager)
		request    string
		validate   func(*testing.T, api.RPCResponse)
	}{
		{
			name: "Should execute OwnerOf",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				setUpOwnerOfMocks(t, tx, "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":2}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer("0x00000000000000000000000026cb70039fe1bd36b4659858d4c4d0cbcafd743a"), getJsonRawMessagePointer("2"))
			},
		},
		{
			name: "Should execute OwnerOf without id",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				setUpOwnerOfMocks(t, tx, "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"]}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer("0x00000000000000000000000026cb70039fe1bd36b4659858d4c4d0cbcafd743a"), nil)
			},
		},
		{
			name: "Should execute OwnerOf with an error from ownerOf",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				tx.EXPECT().OwnerOf(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), gomock.Any()).Return(common.Address{}, fmt.Errorf("error")).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusBadRequest, getHexJsonRawMessagePointer(""), getJsonRawMessagePointer("1"))
			},
		},
		{
			name: "Should execute OwnerOf with an error from create contract",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				tx.EXPECT().LoadContractTrees(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(fmt.Errorf("error")).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusBadRequest, getHexJsonRawMessagePointer(""), getJsonRawMessagePointer("1"))
			},
		},
		{
			name: "Should execute BalanceOf with 0 assets",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				setUpBalanceOfMocks(t, tx, "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", "0x1b0b4a597c764400ea157ab84358c8788a89cd28", 0)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer(hexStringZero), getJsonRawMessagePointer("1"))
			},
		},
		{
			name: "Should execute BalanceOf with 1 assets",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				setUpBalanceOfMocks(t, tx, "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", "0x1b0b4a597c764400ea157ab84358c8788a89cd28", 1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer(hexStringOne), getJsonRawMessagePointer("1"))
			},
		},
		{
			name: "Should execute BalanceOf with 15455 assets",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				setUpBalanceOfMocks(t, tx, "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", "0x1b0b4a597c764400ea157ab84358c8788a89cd28", 15455)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":111}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer("0x0000000000000000000000000000000000000000000000000000000000003c5f"), getJsonRawMessagePointer("111"))
			},
		},
		{
			name: "Should execute BalanceOf with an error from balanceOf",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				tx.EXPECT().BalanceOf(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), common.HexToAddress("0x1b0b4a597c764400ea157ab84358c8788a89cd28")).Return(nil, fmt.Errorf("error")).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusBadRequest, getHexJsonRawMessagePointer(""), getJsonRawMessagePointer("1"))
			},
		},
		{
			name: "Should execute TokenOfOwnerByIndex",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				tx.EXPECT().TokenOfOwnerByIndex(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), common.HexToAddress("0x1b0b4a597c764400ea157ab84358c8788a89cd28"), 1).Return(big.NewInt(1), nil).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x2f745c590000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd280000000000000000000000000000000000000000000000000000000000000001","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer(hexStringOne), getJsonRawMessagePointer("1"))
			},
		},
		{
			name: "Should execute TokenOfOwnerByIndex with an error from tokenOfOwnerByIndex",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				tx.EXPECT().TokenOfOwnerByIndex(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), common.HexToAddress("0x1b0b4a597c764400ea157ab84358c8788a89cd28"), 1).Return(nil, fmt.Errorf("error")).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x2f745c590000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd280000000000000000000000000000000000000000000000000000000000000001","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusBadRequest, getHexJsonRawMessagePointer(""), getJsonRawMessagePointer("1"))
			},
		},
		{
			name: "Should execute TokenByIndex",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				tx.EXPECT().TokenByIndex(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), 1).Return(big.NewInt(1), nil).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x4f6ccce70000000000000000000000000000000000000000000000000000000000000001","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer(hexStringOne), getJsonRawMessagePointer("1"))
			},
		},
		{
			name: "Should execute TotalSupply",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkleTreeMocks(t, tx)
				tx.EXPECT().TotalSupply(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(int64(1), nil).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x18160ddd","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer(hexStringOne), getJsonRawMessagePointer("1"))
			},
		},
		{
			name: "Should execute TokenURI",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpHttpClientMocks(t, httpClient, `{"jsonrpc":"2.0","id":1,"result":"0x00477777730000000000"}`, nil)
				storage.EXPECT().NewTransaction().Return(tx, nil)
				tx.EXPECT().Discard()
				tx.EXPECT().GetCollectionAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A").Return(common.HexToAddress("0xffFFFFfFfffFFFFffffFFfFe00000000000001bE"), nil)
				tx.EXPECT().GetLastEvoBlock().Return(model.Block{Number: uint64(589660)}, nil)
				// 0x8ff5c is 589660 in base 16
				setUpMockRpcMethodManager(t, rpcMethodManager, api.JSONRPCRequest{}, api.RPCMethodEthCall, "0x8ff5c")
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0xc87b56dd0000000000000000000000000000000000000000000000000000000000000064","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer("0x00477777730000000000"), getJsonRawMessagePointer("1"))
			},
		},

		{
			name: "Should execute blocknumber",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx, httpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager) {
				setUpTransactionMocks(t, storage, tx)
				tx.EXPECT().GetLastOwnershipBlock().Return(model.Block{Number: uint64(42971043)}, nil).Times(1)
			},
			request: `{"method":"eth_blockNumber","params":[],"id":11,"jsonrpc":"2.0"}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, getHexJsonRawMessagePointer("0x28fafa3"), getJsonRawMessagePointer("11"))
			},
		},
	}

	// Run tests
	for _, tt := range testCases {
		tt := tt // Shadow loop variable otherwise it could be overwrittens
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel
			ctrl, storage, mockHttpClient, mockMethodManager := setupMocks(t, tt.setupMocks)
			defer ctrl.Finish()

			jsonRquest := createRequest(t, tt.request)
			request := httptest.NewRequest(http.MethodPost, "/rpc", bytes.NewBufferString(tt.request))

			h := api.NewUniversalMintingRPCHandler(
				api.WithEvoHttpClient(mockHttpClient),
				api.WithRPCMethodManager(mockMethodManager),
			)

			result := h.HandleUniversalMinting(request, jsonRquest, storage)
			tt.validate(t, result)
		})
	}
}

func setupMocks(t *testing.T, mockSetup func(storage *mockTx.MockService, tx *mockTx.MockTx, mockHttpClient *mock.MockHTTPClientInterface, rpcMethodManager *mock.MockRPCMethodManager)) (*gomock.Controller, *mockTx.MockService, *mock.MockHTTPClientInterface, *mock.MockRPCMethodManager) {
	ctrl := gomock.NewController(t)
	storage := mockTx.NewMockService(ctrl)
	tx := mockTx.NewMockTx(ctrl)
	mockHttpClient := mock.NewMockHTTPClientInterface(ctrl)
	mockMethodManager := mock.NewMockRPCMethodManager(ctrl)
	mockSetup(storage, tx, mockHttpClient, mockMethodManager)
	return ctrl, storage, mockHttpClient, mockMethodManager
}

func setupMerkleTreeMocks(t *testing.T, tx *mockTx.MockTx) {
	t.Helper()
	tx.EXPECT().LoadContractTrees(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil).Times(1)
}

func setUpTransactionMocks(t *testing.T, storage *mockTx.MockService, tx *mockTx.MockTx) {
	t.Helper()
	storage.EXPECT().NewTransaction().Return(tx, nil).Times(1)
	tx.EXPECT().Discard().AnyTimes()
}

func setUpOwnerOfMocks(t *testing.T, tx *mockTx.MockTx, addressContract, ownerReturnAddress string) {
	t.Helper()
	tx.EXPECT().OwnerOf(common.HexToAddress(addressContract), gomock.Any()).Return(common.HexToAddress(ownerReturnAddress), nil).Times(1)
}

func setUpBalanceOfMocks(t *testing.T, tx *mockTx.MockTx, addressContract, ownerReturnAddress string, balance int64) {
	t.Helper()
	tx.EXPECT().BalanceOf(common.HexToAddress(addressContract), common.HexToAddress(ownerReturnAddress)).Return(big.NewInt(balance), nil).Times(1)
}

func setUpHttpClientMocks(t *testing.T, httpClient *mock.MockHTTPClientInterface, response string, err error) {
	t.Helper()
	if err != nil {
		httpClient.EXPECT().Do(gomock.Any()).Return(nil, err).Times(1)
	} else {
		httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(response)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil).Times(1)
	}
}

func setUpMockRpcMethodManager(t *testing.T, rpcMethodManager *mock.MockRPCMethodManager, _ api.JSONRPCRequest, method api.RPCMethod, blockNumber string) {
	t.Helper()
	rpcMethodManager.EXPECT().ReplaceBlockTag(gomock.Any(), method, blockNumber).Return(nil).Times(1)
}

func createRequest(t *testing.T, requestBody string) api.JSONRPCRequest {
	t.Helper()
	var jsonRPCRequest api.JSONRPCRequest
	// tt.requestBody to []byte
	body := []byte(requestBody)
	if err := json.Unmarshal(body, &jsonRPCRequest); err != nil {
		t.Fatalf("error unmarshalling request: %v", err)
	}

	return jsonRPCRequest
}

func validateResponse(t *testing.T, rr api.RPCResponse, expectedStatus int, expectedResponse, expectedId *json.RawMessage) {
	t.Helper()
	if expectedStatus == http.StatusOK {
		compareRawMessage(t, rr.Result, expectedResponse)
	} else if rr.Error.Code != api.ErrorCodeInvalidRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Error.Code, api.ErrorCodeInvalidRequest)
	}
	compareRawMessage(t, rr.ID, expectedId)
}
