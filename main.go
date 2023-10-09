package main

import (
	"log/slog"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/laos-universal-node/scanner"
)

func main() {
	cli, err := ethclient.Dial("https://eth.llamarpc.com")
	if err != nil {
		slog.Error("error instantiating eth client", "msg", err.Error())
		os.Exit(1)
	}

	// Bored Ape ERC721 Ethereum contract
	contract := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
	scanner.ScanEvents(cli, contract, big.NewInt(18288287), nil)
}
