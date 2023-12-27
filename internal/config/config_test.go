package config_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/freeverseio/laos-universal-node/internal/config"
)

// do not run tests in parallel to avoid data races with os.Args and flag.CommandLine

func TestValidEvoRPC(t *testing.T) {
	tests := []struct {
		name                    string
		expectedGlobalConsensus string
		expectedParachain       uint64
		evoRPC                  string
	}{
		{
			name:                    "set valid global consensus and parachain with caladan rpc",
			expectedGlobalConsensus: "0:0x22c48a576c33970622a2b4686a8aa5e4b58350247d69fb5d8015f12a8c8e1e4c",
			expectedParachain:       2900,
			evoRPC:                  "caladan",
		},
		{
			name:                    "set valid global consensus and parachain with klaos rpc",
			expectedGlobalConsensus: "3",
			expectedParachain:       3336,
			evoRPC:                  "klaos",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			initialFlags := cloneFlagSet()
			defer resetCommandLineFlags(initialFlags)
			initialOsArgs := os.Args
			defer resetOsArgs(initialOsArgs)

			os.Args = append(os.Args, "-evo_rpc", tt.evoRPC)
			c, err := config.Load()
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

func TestInvalidEvoRPC(t *testing.T) {
	t.Run("fails when evo rpc is not recognized", func(t *testing.T) {
		initialFlags := cloneFlagSet()
		defer resetCommandLineFlags(initialFlags)
		initialOsArgs := os.Args
		defer resetOsArgs(initialOsArgs)

		expectedErr := fmt.Errorf("unknown evolution chain rpc provided: RPCURL")
		os.Args = append(os.Args, "-evo_rpc", "RPCURL")
		_, err := config.Load()
		if err == nil {
			t.Fatalf("got no error while an error was expected")
		}
		if err.Error() != expectedErr.Error() {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr.Error())
		}
	})
}

func cloneFlagSet() *flag.FlagSet {
	copyOfFlagSet := flag.NewFlagSet(flag.CommandLine.Name(), flag.CommandLine.ErrorHandling())

	flag.CommandLine.VisitAll(func(flag *flag.Flag) {
		copyOfFlagSet.Var(flag.Value, flag.Name, flag.Usage)
	})

	return copyOfFlagSet
}

func resetCommandLineFlags(original *flag.FlagSet) {
	flag.CommandLine = original
}

func resetOsArgs(original []string) {
	os.Args = original
}
