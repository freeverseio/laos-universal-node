// Package erc721 provides a Go implementation of the ERC721 token standard for non-fungible tokens (NFTs)
// on the Ethereum blockchain.
package erc721

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/freeverseio/laos-universal-node/internal/blockchain"
)

// ProcessCall processes an ERC721 token contract call.
func ProcessCall(
	data string,
	to common.Address,
	ethcli blockchain.EthRPCClient,
	contractAddr common.Address,
	chainID uint64,
) (hexutil.Bytes, error) {
	callData, err := NewCallData(data)
	if err != nil {
		return nil, err
	}
	// check that only the supported methods are called
	_, err = callData.Method()
	if err != nil {
		return nil, err
	}

	var result string
	err = ethcli.Call(&result, "eth_call", map[string]interface{}{
		"to":   to,
		"data": data,
	}, "latest")

	if err != nil {
		return nil, err
	}
	return hexutil.Decode(result)
}