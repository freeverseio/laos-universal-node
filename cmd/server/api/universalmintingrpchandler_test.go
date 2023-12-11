package api_test

import (
	"bytes"
	"fmt"
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	mockTx "github.com/freeverseio/laos-universal-node/internal/state/mock"
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
		setupMocks func(*mockTx.MockService, *mockTx.MockTx)
		request    string
		validate   func(*testing.T, api.RPCResponse)
	}{
		{
			name: "Should execute OwnerOf",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				setUpOwnerOfMocks(t, tx, "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, "0x00000000000000000000000026cb70039fe1bd36b4659858d4c4d0cbcafd743a")
			},
		},
		{
			name: "Should execute OwnerOf with an error from ownerOf",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				tx.EXPECT().OwnerOf(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), gomock.Any()).Return(common.Address{}, fmt.Errorf("error")).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusBadRequest, "")
			},
		},
		{
			name: "Should execute OwnerOf with an error from create contract",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				tx.EXPECT().CreateTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil, nil, nil, fmt.Errorf("error")).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusBadRequest, "")
			},
		},
		{
			name: "Should execute BalanceOf with 0 assets",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				setUpBalanceOfMocks(t, tx, "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", "0x1b0b4a597c764400ea157ab84358c8788a89cd28", 0)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, hexStringZero)
			},
		},
		{
			name: "Should execute BalanceOf with 1 assets",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				setUpBalanceOfMocks(t, tx, "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", "0x1b0b4a597c764400ea157ab84358c8788a89cd28", 1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, hexStringOne)
			},
		},
		{
			name: "Should execute BalanceOf with 15455 assets",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				setUpBalanceOfMocks(t, tx, "0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A", "0x1b0b4a597c764400ea157ab84358c8788a89cd28", 15455)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, "0x0000000000000000000000000000000000000000000000000000000000003c5f")
			},
		},
		{
			name: "Should execute BalanceOf with an error from balanceOf",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				tx.EXPECT().BalanceOf(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), common.HexToAddress("0x1b0b4a597c764400ea157ab84358c8788a89cd28")).Return(nil, fmt.Errorf("error")).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusBadRequest, "")
			},
		},
		{
			name: "Should execute TokenOfOwnerByIndex",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				tx.EXPECT().TokenOfOwnerByIndex(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), common.HexToAddress("0x1b0b4a597c764400ea157ab84358c8788a89cd28"), 1).Return(big.NewInt(1), nil).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x2f745c590000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd280000000000000000000000000000000000000000000000000000000000000001","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, hexStringOne)
			},
		},
		{
			name: "Should execute TokenOfOwnerByIndex with an error from tokenOfOwnerByIndex",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				tx.EXPECT().TokenOfOwnerByIndex(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), common.HexToAddress("0x1b0b4a597c764400ea157ab84358c8788a89cd28"), 1).Return(nil, fmt.Errorf("error")).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x2f745c590000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd280000000000000000000000000000000000000000000000000000000000000001","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusBadRequest, "")
			},
		},
		{
			name: "Should execute TokenByIndex",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				tx.EXPECT().TokenByIndex(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), 1).Return(big.NewInt(1), nil).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x4f6ccce70000000000000000000000000000000000000000000000000000000000000001","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, hexStringOne)
			},
		},
		{
			name: "Should execute TotalSupply",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				setupMerkeTreeMocks(t, tx)
				tx.EXPECT().TotalSupply(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(int64(1), nil).Times(1)
			},
			request: `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x18160ddd","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}, "latest"],"id":1}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, hexStringOne)
			},
		},

		{
			name: "Should execute blocknumber",
			setupMocks: func(storage *mockTx.MockService, tx *mockTx.MockTx) {
				setUpTransactionMocks(t, storage, tx)
				tx.EXPECT().GetCurrentOwnershipBlock().Return(uint64(42971043), nil).Times(1)
			},
			request: `{"method":"eth_blockNumber","params":[],"id":1,"jsonrpc":"2.0"}`,
			validate: func(t *testing.T, rr api.RPCResponse) {
				validateResponse(t, rr, http.StatusOK, "0x28fafa2")
			},
		},
	}

	// Run tests
	for _, tt := range testCases {
		tt := tt // Shadow loop variable otherwise it could be overwrittens
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel
			ctrl, storage := setupMocks(t, tt.setupMocks)
			defer ctrl.Finish()

			request := createRequest(t, tt.request)
			h := api.UniversalMintingRPCHandler{}
			result := h.HandleUniversalMinting(request, storage)
			fmt.Println(result)
			tt.validate(t, result)
		})
	}
}

func setupMocks(t *testing.T, mockSetup func(storage *mockTx.MockService, tx *mockTx.MockTx)) (*gomock.Controller, *mockTx.MockService) {
	ctrl := gomock.NewController(t)
	storage := mockTx.NewMockService(ctrl)
	tx := mockTx.NewMockTx(ctrl)
	mockSetup(storage, tx)
	return ctrl, storage
}

func setupMerkeTreeMocks(t *testing.T, tx *mockTx.MockTx) {
	t.Helper()
	tx.EXPECT().CreateTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil, nil, nil, nil).Times(1)
	tx.EXPECT().SetTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), nil, nil, nil).Times(1)
}

func setUpTransactionMocks(t *testing.T, storage *mockTx.MockService, tx *mockTx.MockTx) {
	t.Helper()
	storage.EXPECT().NewTransaction().Return(tx).Times(1)
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

func createRequest(t *testing.T, requestBody string) *http.Request {
	request, err := http.NewRequest("POST", "/your-endpoint", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func validateResponse(t *testing.T, rr api.RPCResponse, expectedStatus int, expectedResponse string) {
	if expectedStatus == http.StatusOK {
		if rr.Result != expectedResponse {
			t.Errorf("handler returned unexpected result: got %v want %v", rr.Result, expectedResponse)
		}
	} else {
		if rr.Error.Code != api.ErrorCodeInvalidRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Error.Code, api.ErrorCodeInvalidRequest)
		}
	}
}
