package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/cmd/server/api"
	mockTx "github.com/freeverseio/laos-universal-node/internal/state/mock"
	"go.uber.org/mock/gomock"
)

func TestUniversalMintingRPCHandler(t *testing.T) {
	t.Run("Should execute OwnerOf", func(t *testing.T) {
		ctrl, storage := setupMocks(t, func(storage *mockTx.MockService, tx *mockTx.MockTx) {
			storage.EXPECT().NewTransaction().Return(tx).Times(1)
			tx.EXPECT().Discard().AnyTimes()
			tx.EXPECT().CreateTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil, nil, nil, nil).Times(1)
			tx.EXPECT().SetTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), nil, nil, nil).Return(nil).Times(1)
			tx.EXPECT().OwnerOf(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), gomock.Any()).Return(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), nil).Times(1)
		})
		defer ctrl.Finish()

		request := createRequest(t, `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`)

		rr := runHandler(t, request, storage)

		// Check the status code and response body
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Should execute OwnerOf with an error from ownerOf", func(t *testing.T) {
		ctrl, storage := setupMocks(t, func(storage *mockTx.MockService, tx *mockTx.MockTx) {
			storage.EXPECT().NewTransaction().Return(tx).Times(1)
			tx.EXPECT().Discard().AnyTimes()
			tx.EXPECT().CreateTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil, nil, nil, nil).Times(1)
			tx.EXPECT().SetTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), nil, nil, nil).Return(nil).Times(1)
			tx.EXPECT().OwnerOf(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), gomock.Any()).Return(common.Address{}, fmt.Errorf("error")).Times(1)
		})
		defer ctrl.Finish()

		request := createRequest(t, `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`)

		rr := runHandler(t, request, storage)

		var response api.JSONRPCErrorResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Error unmarshalling response body: %v", err)
		}

		if response.Error.Code != api.ErrorCodeInvalidRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", response.Error.Code, api.ErrorCodeInvalidRequest)
		}
	})

	t.Run("Should execute OwnerOf with an error from create contract", func(t *testing.T) {
		ctrl, storage := setupMocks(t, func(storage *mockTx.MockService, tx *mockTx.MockTx) {
			storage.EXPECT().NewTransaction().Return(tx).Times(1)
			tx.EXPECT().Discard().AnyTimes()
			tx.EXPECT().CreateTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil, nil, nil, fmt.Errorf("error")).Times(1)
		})
		defer ctrl.Finish()

		request := createRequest(t, `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x6352211e0000000000000000000000021b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`)
		rr := runHandler(t, request, storage)
		var response api.JSONRPCErrorResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Error unmarshalling response body: %v", err)
		}

		if response.Error.Code != api.ErrorCodeInvalidRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", response.Error.Code, api.ErrorCodeInvalidRequest)
		}
	})

	t.Run("Should execute BalanceOf", func(t *testing.T) {
		ctrl, storage := setupMocks(t, func(storage *mockTx.MockService, tx *mockTx.MockTx) {
			storage.EXPECT().NewTransaction().Return(tx).Times(1)
			tx.EXPECT().Discard().AnyTimes()
			tx.EXPECT().CreateTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil, nil, nil, nil).Times(1)
			tx.EXPECT().SetTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), nil, nil, nil).Return(nil).Times(1)
			tx.EXPECT().BalanceOf(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), common.HexToAddress("0x1b0b4a597c764400ea157ab84358c8788a89cd28")).Return(big.NewInt(1), nil).Times(1)
		})
		defer ctrl.Finish()

		request := createRequest(t, `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a082310000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd28","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`)

		rr := runHandler(t, request, storage)

		// Check the status code and response body
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Should execute TokenOfOwnerByIndex", func(t *testing.T) {
		ctrl, storage := setupMocks(t, func(storage *mockTx.MockService, tx *mockTx.MockTx) {
			storage.EXPECT().NewTransaction().Return(tx).Times(1)
			tx.EXPECT().Discard().AnyTimes()
			tx.EXPECT().CreateTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil, nil, nil, nil).Times(1)
			tx.EXPECT().SetTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), nil, nil, nil).Return(nil).Times(1)
			tx.EXPECT().TokenOfOwnerByIndex(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), common.HexToAddress("0x1b0b4a597c764400ea157ab84358c8788a89cd28"), 1).Return(big.NewInt(1), nil).Times(1)
		})
		defer ctrl.Finish()
		request := createRequest(t, `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x2f745c590000000000000000000000001b0b4a597c764400ea157ab84358c8788a89cd280000000000000000000000000000000000000000000000000000000000000001","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`)

		rr := runHandler(t, request, storage)

		// Check the status code and response body
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Should execute TokenByIndex", func(t *testing.T) {
		ctrl, storage := setupMocks(t, func(storage *mockTx.MockService, tx *mockTx.MockTx) {
			storage.EXPECT().NewTransaction().Return(tx).Times(1)
			tx.EXPECT().Discard().AnyTimes()
			tx.EXPECT().CreateTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil, nil, nil, nil).Times(1)
			tx.EXPECT().SetTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), nil, nil, nil).Return(nil).Times(1)
			tx.EXPECT().TokenByIndex(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), 1).Return(big.NewInt(1), nil).Times(1)
		})
		defer ctrl.Finish()

		request := createRequest(t, `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x4f6ccce70000000000000000000000000000000000000000000000000000000000000001","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`)

		rr := runHandler(t, request, storage)

		// Check the status code and response body
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Should execute TotalSupply", func(t *testing.T) {
		ctrl, storage := setupMocks(t, func(storage *mockTx.MockService, tx *mockTx.MockTx) {
			storage.EXPECT().NewTransaction().Return(tx).Times(1)
			tx.EXPECT().Discard().AnyTimes()
			tx.EXPECT().CreateTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(nil, nil, nil, nil).Times(1)
			tx.EXPECT().SetTreesForContract(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"), nil, nil, nil).Return(nil).Times(1)
			tx.EXPECT().TotalSupply(common.HexToAddress("0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A")).Return(int64(1), nil).Times(1)
		})
		defer ctrl.Finish()

		request := createRequest(t, `{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x18160ddd","to":"0x26CB70039FE1bd36b4659858d4c4D0cBcafd743A"}],"id":1}`)

		rr := runHandler(t, request, storage)

		// Check the status code and response body
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}

func setupMocks(t *testing.T, mockSetup func(storage *mockTx.MockService, tx *mockTx.MockTx)) (*gomock.Controller, *mockTx.MockService) {
	ctrl := gomock.NewController(t)
	storage := mockTx.NewMockService(ctrl)
	tx := mockTx.NewMockTx(ctrl)
	mockSetup(storage, tx)
	return ctrl, storage
}

func createRequest(t *testing.T, requestBody string) *http.Request {
	request, err := http.NewRequest("POST", "/your-endpoint", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func runHandler(t *testing.T, request *http.Request, storage *mockTx.MockService) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading request body: %v", err)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for further handling

		var req api.JSONRPCRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Errorf("Error unmarshalling request body: %v", err)
		}

		h := api.GlobalRPCHandler{}
		h.SetJsonRPCRequest(req)
		h.SetStateService(storage)
		h.UniversalMintingRPCHandler(w, r)
	})

	handler.ServeHTTP(rr, request)
	return rr
}
