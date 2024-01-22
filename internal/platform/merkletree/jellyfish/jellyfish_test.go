package jellyfish_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree/jellyfish"
	storage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
)

// these tests are integration tests that use badgerDB. They are fast though and there are not many of them
// so we use badgerDB instead of memory or mocking the storage
func TestInitialRoot(t *testing.T) {
	randInt, err := rand.Int(rand.Reader, big.NewInt(10000000000))
	if err != nil {
		t.Errorf("error on generating random int, got: %v", err)
	}
	dbPath := fmt.Sprintf("./tmp/badger-%d", randInt)
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithLoggingLevel(badger.ERROR).WithMemTableSize(1 << 30))
	storageService := storage.NewService(db)
	tx := storageService.NewTransaction()
	defer func() {
		err = db.Close()
		if err != nil {
			t.Errorf("error on closing database, got: %v", err)
		}
		err = os.RemoveAll(dbPath)
		if err != nil {
			t.Errorf("error on removing database, got: %v", err)
		}
	}()

	tree, err := jellyfish.New(tx, "")
	root := tree.Root()

	if root.Cmp(common.Hash{}) != 0 {
		t.Errorf("expected root to be 0, got: %v", root)
	}
}

func TestSetLeafAndRootIsChangedSetRootAndCorrectLeafIsRecalled(t *testing.T) {
	randInt, err := rand.Int(rand.Reader, big.NewInt(10000000000))
	if err != nil {
		t.Errorf("error on generating random int, got: %v", err)
	}
	dbPath := fmt.Sprintf("./tmp/badger-%d", randInt)
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithLoggingLevel(badger.ERROR).WithMemTableSize(1 << 30))
	storageService := storage.NewService(db)
	tx := storageService.NewTransaction()
	defer func() {
		err = db.Close()
		if err != nil {
			t.Errorf("error on closing database, got: %v", err)
		}
		err = os.RemoveAll(dbPath)
		if err != nil {
			t.Errorf("error on removing database, got: %v", err)
		}
	}()

	tree, err := jellyfish.New(tx, "")
	err = tree.SetLeaf(big.NewInt(1), common.HexToHash("0x1"))
	if err != nil {
		t.Errorf("error when setting leaf: %v", err)
	}

	root1 := tree.Root()
	if root1.Cmp(common.HexToHash("0xfe88890366165b69d5a030600b98e8645ebebd242049ee6fec678b0a13fa021e")) != 0 {
		t.Errorf("expected root to be 0, got: %v", root1)
	}

	err = tree.SetLeaf(big.NewInt(1), common.HexToHash("0x2"))
	if err != nil {
		t.Errorf("error when setting leaf: %v", err)
	}
	root2 := tree.Root()
	if root2.Cmp(common.HexToHash("0x8ff34634ce84f5a17a1f08d72a53de94df8ccaa2b30a3b2711a9b38b24ea336b")) != 0 {
		t.Errorf("expected root to be 0, got: %v", root2)
	}

	leaf, err := tree.Leaf(big.NewInt(1))
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if leaf.Cmp(common.HexToHash("0x2")) != 0 {
		t.Errorf("expected leaf to be 0x2, got: %v", leaf)
	}

	tree.SetRoot(root1)
	root3 := tree.Root()
	if root3.Cmp(root1) != 0 {
		t.Errorf("expected root to be 0, got: %v", root2)
	}

	leaf, err = tree.Leaf(big.NewInt(1))
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if leaf.Cmp(common.HexToHash("0x1")) != 0 {
		t.Errorf("expected leaf to be 0x1, got: %v", leaf)
	}
}
