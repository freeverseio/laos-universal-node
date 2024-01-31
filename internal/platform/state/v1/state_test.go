package v1_test

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/platform/state/v1"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
)

func TestLoadMerkleTrees(t *testing.T) {
	t.Parallel()
	t.Run("fails when contract is 0x0", func(t *testing.T) {
		t.Parallel()
		tx, err := createTransaction()
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		expectedErr := fmt.Sprintf("contract address is " + common.Address{}.String())

		err = tx.LoadContractTrees(common.HexToAddress("0x0"))
		if err == nil {
			t.Errorf("got no error while an error was expected")
		}
		if err != nil && err.Error() != expectedErr {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr)
		}
	})
	t.Run("successfully loads merkle trees in memory, minting and checking balances", func(t *testing.T) {
		t.Parallel()
		tx, err := createTransaction()
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}
		contract := common.HexToAddress("0x500")

		err = tx.LoadContractTrees(contract)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		// if tx.Mint works, we are sure that the merkle trees are correctly loaded in memory
		// as Mint accesses all merkle trees
		mintEvent := model.MintedWithExternalURI{
			Slot:        big.NewInt(1),
			To:          common.HexToAddress("0x3"),
			TokenURI:    "tokenURI",
			TokenId:     big.NewInt(1),
			BlockNumber: 100,
			Timestamp:   1000,
			TxIndex:     1,
		}
		err = tx.Mint(contract, &mintEvent)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		owner, err := tx.OwnerOf(contract, big.NewInt(1))
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		if owner != common.HexToAddress("0x1") {
			t.Errorf(`got owner "%s" when expected "0x1"`, owner.Hex())
		}

		err = tx.UpdateContractState(contract, uint64(10))
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		err = tx.LoadContractTrees(contract) // loading tree the second time does not change anything
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		balance, err := tx.BalanceOf(contract, common.HexToAddress("0x1"))
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		if balance.Cmp( big.NewInt(1)) != 0 {
			t.Errorf(`got balance "%d" when expected "1"`, balance)
		}

		totalSupply, err := tx.TotalSupply(contract)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		if totalSupply != 1 {
			t.Errorf(`got total supply "%d" when expected "1"`, totalSupply)
		}

		tokenOfOwner, err := tx.TokenOfOwnerByIndex(contract, common.HexToAddress("0x1"), 0)
		if err != nil {
			t.Errorf(`got error "%v" when no error was expected`, err)
		}

		if tokenOfOwner.Cmp(big.NewInt(1)) != 0 {
			t.Errorf(`got token of owner "%d" when expected "1"`, tokenOfOwner)
		}
	})
}

func createTransaction() (state.Tx, error) {
	mem := memory.New()
	stateService := v1.NewStateService(mem)
	return stateService.NewTransaction()
}

func formatNumberForSorting(blockNumber uint64, blockNumberDigits uint16) string {
	// Convert the block number to a string
	blockNumberString := strconv.FormatUint(blockNumber, 10)
	// Pad with leading zeros if shorter
	for len(blockNumberString) < int(blockNumberDigits) {
		blockNumberString = "0" + blockNumberString
	}
	return blockNumberString
}
