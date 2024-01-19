package jellyfish_test

import (
	"fmt"
	"math/big"
	"math/rand"
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
	dbPath := fmt.Sprintf("./tmp/badger-%d", rand.Intn(10000000000))
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithLoggingLevel(badger.ERROR).WithMemTableSize(1 << 30))
	storageService := storage.NewService(db)
	tx := storageService.NewTransaction()
	defer func() {
		err = db.Close()
		if err != nil {
			t.Errorf("error on closing database, got: %v", err)
		}
		os.RemoveAll(dbPath)
	}()

	tree, err := jellyfish.New(tx, "")
	root := tree.Root()

	if root.Cmp(common.Hash{}) != 0 {
		t.Errorf("expected root to be 0, got: %v", root)
	}

}

func TestSetLeafAndRootIsChangedSetRootAndCorrectLeafIsRecalled(t *testing.T) {
	dbPath := fmt.Sprintf("./tmp/badger-%d", rand.Intn(10000000000))
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithLoggingLevel(badger.ERROR).WithMemTableSize(1 << 30))
	storageService := storage.NewService(db)
	tx := storageService.NewTransaction()
	defer func() {
		err = db.Close()
		if err != nil {
			t.Errorf("error on closing database, got: %v", err)
		}
		os.RemoveAll(dbPath)
	}()

	tree, err := jellyfish.New(tx, "")
	tree.SetLeaf(big.NewInt(1), common.HexToHash("0x1"))

	root1 := tree.Root()
	if root1.Cmp(common.HexToHash("0xbe3f2c5bdf4ad79b6b9347ed631296af1017514d13da634ef4d4d5409e8b3603")) != 0 {
		t.Errorf("expected root to be 0, got: %v", root1)
	}

	tree.SetLeaf(big.NewInt(1), common.HexToHash("0x2"))
	root2 := tree.Root()
	if root2.Cmp(common.HexToHash("0x8c88fdd174c9acf5c02eac848a1090e0583cc22c1c55274a3dd9f22be842b663")) != 0 {
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
