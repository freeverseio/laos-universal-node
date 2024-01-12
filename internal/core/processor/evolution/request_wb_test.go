package evolution

import (
	"math/big"
	"testing"
)

func TestHexToDecimal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input     string
		expected  *big.Int
		expectErr bool
	}{
		{"1a", big.NewInt(26), false},
		{"0x1a", big.NewInt(26), false},
		{"abc", big.NewInt(2748), false},
		{"0xabc", big.NewInt(2748), false},
		{"invalid", nil, true},   // Invalid hex
		{"0xinvalid", nil, true}, // Invalid hex
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			result, err := hexToDecimal(tt.input)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("got %T, expected nil error", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if result.Cmp(tt.expected) != 0 {
					t.Fatalf("got %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}
