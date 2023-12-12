package erc721_test

import (
	"testing"

	"github.com/freeverseio/laos-universal-node/internal/platform/rpc/erc721"
)

func TestAbiEncodeString(t *testing.T) {
	t.Parallel()
	t.Run("abi encode string success", func(t *testing.T) {
		t.Parallel()
		input := "ipfs://Qmdt3BvDYb4r4ZiMdjq8D3jExzqprKphcejZ6mhdwP14d4"
		expected := "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000035697066733a2f2f516d64743342764459623472345a694d646a713844336a45787a7170724b706863656a5a366d68647750313464340000000000000000000000"
		got, err := erc721.AbiEncodeString(input)
		if err != nil {
			t.Errorf("got error %s while no error was expected", err.Error())
		}
		if got != expected {
			t.Fatalf("got abi-encoded string %s, expected %s", got, expected)
		}
	})
}
