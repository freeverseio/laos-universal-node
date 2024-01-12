package evolution

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
)

const (
	jsonContentType = "application/json"
	postMethod      = "POST"
)

// Define an interface for HTTP client operations
type LaosRPCRequests interface {
	LatestFinalizedBlockHash() (string, error)
	BlockNumber(blockHash string) (*big.Int, error)
}

type LaosHTTP struct {
	client http.Client
	url    string
}

func NewLaosHTTP(url string) LaosRPCRequests {
	return LaosHTTP{
		client: http.Client{},
		url:    url,
	}
}

// BlockHeader struct to represent the block header in the JSON response
type BlockHeader struct {
	ParentHash     string `json:"parentHash"`
	Number         string `json:"number"`
	StateRoot      string `json:"stateRoot"`
	ExtrinsicsRoot string `json:"extrinsicsRoot"`
}

// Block struct to represent the block in the JSON response
type Block struct {
	Header BlockHeader `json:"header"`
}

// Response struct to represent the entire JSON response
type ChainGetBlock struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Block Block `json:"block"`
	} `json:"result"`
	ID int `json:"id"`
}

// Response struct to represent the entire JSON response
type ChainGetFinalizedHead struct {
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
}

func (l LaosHTTP) makeRequest(payload []byte, result interface{}) error {
	req, err := http.NewRequest(postMethod, l.url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", jsonContentType)

	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to %s: %w", l.url, err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			slog.Error("error closing response body", "err", err)
		}
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("error in request to %s, got status code: %d", l.url, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	return nil
}

func (l LaosHTTP) LatestFinalizedBlockHash() (string, error) {
	payload := []byte(`{"jsonrpc":"2.0","method":"chain_getFinalizedHead","params":[],"id":1}`)

	var response ChainGetFinalizedHead
	if err := l.makeRequest(payload, &response); err != nil {
		return "", err
	}

	return response.Result, nil
}

func (l LaosHTTP) BlockNumber(blockHash string) (*big.Int, error) {
	payload := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","method":"chain_getBlock","params":[%q],"id":1}`, blockHash))

	var response ChainGetBlock
	if err := l.makeRequest(payload, &response); err != nil {
		return nil, err
	}

	hexNumber := response.Result.Block.Header.Number

	blockNumber, err := hexToDecimal(hexNumber)
	if err != nil {
		return nil, fmt.Errorf("error converting latest finalized block number: %w", err)
	}

	return blockNumber, nil
}

func hexToDecimal(hex string) (*big.Int, error) {
	// Remove the "0x" prefix
	if len(hex) > 2 && hex[:2] == "0x" {
		hex = hex[2:]
	}

	decimal, ok := new(big.Int).SetString(hex, 16)
	if !ok {
		return nil, fmt.Errorf("could not convert %v to decimal value", hex)
	}

	return decimal, nil
}
