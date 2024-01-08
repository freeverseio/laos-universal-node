package v1_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
)

func TestLoadMerkleTree(t *testing.T) {
	t.Parallel()
	t.Run("fails when contract is 0x0", func(t *testing.T) {
		t.Parallel()
		tx := createTransaction()
		expectedErr := fmt.Sprintf("contract address is " + common.Address{}.String())

		err := tx.LoadMerkleTrees(common.HexToAddress("0x0"))
		if err == nil {
			t.Errorf("got no error while an error was expected")
		}
		if err != nil && err.Error() != expectedErr {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr)
		}
	})
	t.Run("successfully loads merkle trees in memory", func(t *testing.T) {
		t.Parallel()
		tx := createTransaction()

		err := tx.LoadMerkleTrees(common.HexToAddress("0x500"))
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
	})
}

func createTransaction() state.Tx {
	mem := memory.New()
	stateService := v1.NewStateService(mem)
	return stateService.NewTransaction()
}
