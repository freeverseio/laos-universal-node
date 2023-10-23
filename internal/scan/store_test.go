package scan_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/scan"
)

func TestReadAll(t *testing.T) {
	s, err := scan.NewFSStorage("../../erc721_contracts.txt")
	if err != nil {
		t.Fatalf("error initializing NewFSStorage : %v", err.Error())
	}

	contracts, err := s.ReadAll(context.Background())
	if err != nil {
		t.Fatalf("error reading file: %v", err.Error())
	}

	if len(contracts) != 1 {
		t.Fatalf("expecting 1 contract, got %d", len(contracts))
	}

	expected := scan.ERC721BridgelessContract{
		Address: common.HexToAddress("0x26cb70039fe1bd36b4659858d4c4d0cbcafd743a"),
		Block:   41346532,
		BaseURI: "evochain1/collectionId/",
	}
	if contracts[0] != expected {
		t.Fatalf("expecting %v, got %v", expected, contracts[0])
	}
}
