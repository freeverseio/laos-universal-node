package v1_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
)

func TestLoadMerkleTree(t *testing.T) {
	t.Parallel()
	t.Run("fails when contract is 0x0", func(t *testing.T) {
		t.Parallel()
		mem := memory.New()
		stateService := v1.NewStateService(mem)
		tx := stateService.NewTransaction()
		expectedErr := fmt.Sprintf("contract address is " + common.Address{}.String())

		err := tx.LoadMerkleTrees(common.HexToAddress("0x0"))
		if err == nil {
			t.Errorf("got no error while an error was expected")
		}
		if err != nil && err.Error() != expectedErr {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr)
		}
	})
}
