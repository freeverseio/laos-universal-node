package config_test

import (
	"flag"
	"fmt"
	"math/big"
	"os"
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
			name:                    "sets global consensus and parachain with klaos nova chain id",
			expectedGlobalConsensus: "0:0x4756c4042a431ad2bbe61d8c4b966c1328e7a8daa0110e9bbd3d4013138a0bd4",
			expectedParachain:       2001,
			evoChainId:              big.NewInt(27181),
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

func TestLoadConfig(t *testing.T) {
	// Do not run this test in parallel because it modifies the global state
	t.Run("loads config with default values", func(t *testing.T) {
		resetFlagSet() // Reset the flag set before defining new flags
		args := []string{"cmd", "--evo_blocks_range=1"}
		os.Args = args
		c, err := config.Load()
		if err != nil {
			t.Fatalf("got error %s while no error was expected", err.Error())
		}
		if c == nil {
			t.Fatalf("got nil config while a config was expected")
		}
		if c.EvoBlocksRange != 1 {
			t.Errorf("got evo blocks range %d, expected 1", c.EvoBlocksRange) // Fixed assertion to match expected behavior
		}
	})
	t.Run("fails when evo blocks range is greater than 1", func(t *testing.T) {
		resetFlagSet() // Reset the flag set before defining new flags
		args := []string{"cmd", "--evo_blocks_range=2"}
		os.Args = args
		_, err := config.Load()
		if err == nil {
			t.Fatalf("got no error while an error was expected")
		}
		expectedErr := "evo_blocks_range can not be bigger than 1"
		if err.Error() != expectedErr {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr)
		}
	})
}

func resetFlagSet() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
