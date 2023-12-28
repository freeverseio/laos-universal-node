package config_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/freeverseio/laos-universal-node/internal/config"
)

func TestValidChainID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		expectedGlobalConsensus string
		expectedParachain       uint64
		evoChainId              *big.Int
	}{
		{
			name:                    "sets global consensus and parachain with caladan chain id",
			expectedGlobalConsensus: "0:0x22c48a576c33970622a2b4686a8aa5e4b58350247d69fb5d8015f12a8c8e1e4c",
			expectedParachain:       2900,
			evoChainId:              big.NewInt(667),
		},
		{
			name:                    "sets global consensus and parachain with klaos chain id",
			expectedGlobalConsensus: "3",
			expectedParachain:       3336,
			evoChainId:              big.NewInt(2718),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &config.Config{}

			err := c.SetGlobalConsensusAndParachain(tt.evoChainId)
			if err != nil {
				t.Fatalf("got error %s while no error was expected", err.Error())
			}
			if c.GlobalConsensus != tt.expectedGlobalConsensus {
				t.Fatalf("got global consensus %s, expected %s", c.GlobalConsensus, tt.expectedGlobalConsensus)
			}
			if c.Parachain != tt.expectedParachain {
				t.Fatalf("got parachain %d, expected %d", c.Parachain, tt.expectedParachain)
			}
		})
	}
}

func TestInvalidChainID(t *testing.T) {
	t.Parallel()
	t.Run("fails when evo chain id is not recognized", func(t *testing.T) {
		t.Parallel()
		expectedErr := fmt.Errorf("unknown evolution chain id: 0")
		c := &config.Config{}
		err := c.SetGlobalConsensusAndParachain(big.NewInt(0))
		if err == nil {
			t.Fatalf("got no error while an error was expected")
		}
		if err.Error() != expectedErr.Error() {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr.Error())
		}
	})
}
