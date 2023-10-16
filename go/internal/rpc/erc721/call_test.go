package erc721

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertHexStringToText(t *testing.T) {
	t.Run("converts hex string to text", func(t *testing.T) {
		b, err := convert64HexStringToText("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000037572690000000000000000000000000000000000000000000000000000000000")
		assert.Nil(t, err, "Expected error to be nil, but got: %v", err)
		assert.Equal(t, b, "uri", "should be equal")
	})

	t.Run("converts hex string to text", func(t *testing.T) {
		b := convert16HexStringToDecimal("0x0000000000000000000000000000000000000000000000000000000000000001")
		assert.Equal(t, b.Int64(), int64(1), "should be equal")
	})
	t.Run("converts hex string to text", func(t *testing.T) {
		b := convert16HexStringToDecimal("0x00000000000000000000000000000000000000000000000000000000000008A7")
		assert.Equal(t, b.Int64(), int64(2215), "should be equal")
	})
}

func TestTokenURI(t *testing.T) {
	t.Run("returns token URI", func(t *testing.T) {
		_, err := tokenURI("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000037572690000000000000000000000000000000000000000000000000000000000")
		assert.Nil(t, err, "Expected error to be nil, but got: %v", err)
	})
}
