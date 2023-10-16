package erc721

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertHexStringToText(t *testing.T) {
	t.Run("converts hex string to text", func(t *testing.T) {
		b, err := convertHexStringToText("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000037572690000000000000000000000000000000000000000000000000000000000")
		assert.Nil(t, err, "Expected error to be nil, but got: %v", err)
		assert.Equal(t, string(b), "uri", "should be equal")
	})
}

func TestTokenURI(t *testing.T) {
	t.Run("returns token URI", func(t *testing.T) {
		_, err := tokenURI("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000037572690000000000000000000000000000000000000000000000000000000000")
		assert.Nil(t, err, "Expected error to be nil, but got: %v", err)
	})
}