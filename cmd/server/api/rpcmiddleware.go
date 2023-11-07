package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
	"github.com/freeverseio/laos-universal-node/internal/scan"
)

// JSONRPCRequest represents the expected structure of a JSON-RPC request.
type JSONRPCRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
	ID      *json.RawMessage  `json:"id,omitempty"` // Pointer allows for an optional field
}

type ParamsRPCRequest struct {
	Data  string `json:"data,omitempty"`
	To    string `json:"to,omitempty"`
	From  string `json:"from,omitempty"`
	Value string `json:"value,omitempty"`
}

// Adjust the middleware to handle different JSON-RPC methods
func PostRpcRequestMiddleware(standardHandler, erc721UniversalMintingHandler http.Handler, st scan.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "No JSON RPC call", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body

		var req JSONRPCRequest
		if errParsing := json.Unmarshal(body, &req); errParsing != nil {
			http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
			return
		}

		if req.JSONRPC != "2.0" {
			http.Error(w, "Invalid JSON-RPC version", http.StatusBadRequest)
			return
		}

		// if it's not an eth_call (non erc721), just pass it through
		if req.Method != "eth_call" {
			standardHandler.ServeHTTP(w, r)
			return
		}

		var params ParamsRPCRequest
		if errParsingParams := json.Unmarshal(req.Params[0], &params); errParsingParams != nil {
			slog.Error("error parsing params", "err", err)
			return
		}

		remoteMinting, err := isUniversalMintingMethod(params.Data)
		if err != nil {
			http.Error(w, "Error checking remote minting method", http.StatusBadRequest)
			return
		}

		if !remoteMinting {
			standardHandler.ServeHTTP(w, r)
			return
		}

		isInContractList, err := isContractInList(params.To, st)
		if err != nil {
			http.Error(w, "Error checking contract in list", http.StatusBadRequest)
			return
		}

		if isInContractList {
			erc721UniversalMintingHandler.ServeHTTP(w, r)
		} else {
			standardHandler.ServeHTTP(w, r)
		}
	})
}

func isContractInList(contractAddress string, st scan.Storage) (bool, error) {
	list, err := st.ReadAll(context.Background())
	if err != nil {
		return false, err
	}

	for _, contract := range list {
		addr := contract.Address.Hex() // convert to string
		if strings.EqualFold(addr, contractAddress) {
			return true, nil
		}
	}
	return false, nil
}

func isUniversalMintingMethod(data string) (bool, error) {
	calldata, err := erc721.NewCallData(data)
	if err != nil {
		return false, err
	}
	_, exists, err := calldata.UniversalMintingMethod()
	if err != nil {
		return false, err
	}

	return exists, nil
}
