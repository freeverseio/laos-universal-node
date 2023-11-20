package sparsemt_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree/sparsemt"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
	"gotest.tools/assert"
)

func TestDag(t *testing.T) {
	t.Parallel()
	t.Run(`depth 0 root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		s, err := sparsemt.New(0, tx, "")
		assert.NilError(t, err)
		hash := s.Root()
		assert.Equal(t, hash.String(), "0x0000000000000000000000000000000000000000000000000000000000000000")

		idx := common.Big0
		leaf := common.HexToHash("0x1")

		err = s.SetLeaf(idx, leaf)
		assert.NilError(t, err)

		hash = s.Root()
		assert.Equal(t, hash.String(), "0x0000000000000000000000000000000000000000000000000000000000000001")

		proof, err := s.Proof(idx)
		assert.NilError(t, err)
		assert.Equal(t, len(proof), 0)
	})
	t.Run(`root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		s, err := sparsemt.New(1, tx, "")
		assert.NilError(t, err)
		hash := s.Root()
		assert.Equal(t, hash.String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	})
	t.Run(`proof of tree is long as depth`, func(t *testing.T) {
		t.Parallel()
		for depth := uint(1); depth < 1000; depth++ {
			service := memory.New()
			tx := service.NewTransaction()

			s, err := sparsemt.New(depth, tx, "")
			assert.NilError(t, err)
			proof, err := s.Proof(big.NewInt(1))
			assert.NilError(t, err)
			assert.Equal(t, len(proof), int(depth))
		}
	})
	t.Run(`proof of empty tree`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		depth := uint(256)
		s, err := sparsemt.New(depth, tx, "")
		assert.NilError(t, err)
		proof, err := s.Proof(big.NewInt(1))
		assert.NilError(t, err)
		assert.Equal(t, len(proof), int(depth))

		for _, hash := range proof {
			assert.Equal(t, hash.String(), sparsemt.Null)
		}
	})
	t.Run(`update leaf 0 of empty tree of depth 1`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		s, err := sparsemt.New(1, tx, "")
		assert.NilError(t, err)

		idx := common.Big0
		leaf := common.HexToHash("0x1")

		err = s.SetLeaf(idx, leaf)
		assert.NilError(t, err)

		hash, err := s.Leaf(idx)
		assert.NilError(t, err)
		assert.Equal(t, hash, leaf)

		assert.Equal(t, s.Root().String(), "0xada5013122d395ba3c54772283fb069b10426056ef8ca54750cb9bb552a59e7d")
		assert.NilError(t, s.SetLeaf(idx, leaf))
		assert.Equal(t, s.Root().String(), "0xada5013122d395ba3c54772283fb069b10426056ef8ca54750cb9bb552a59e7d")
	})
	t.Run(`update leaf 1 of empty tree of depth 1`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		s, err := sparsemt.New(1, tx, "")
		assert.NilError(t, err)

		idx := common.Big1
		leaf := common.HexToHash("0x1")

		err = s.SetLeaf(idx, leaf)
		assert.NilError(t, err)

		hash, err := s.Leaf(idx)
		assert.NilError(t, err)
		assert.Equal(t, hash, leaf)
	})
	t.Run(`ProofExt depth 1`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		node := sparsemt.Node{
			L: common.HexToHash("0x1"),
			R: common.HexToHash("0x2"),
		}
		hash, err := sparsemt.PutNode(tx, node, "prefix")
		assert.NilError(t, err)
		mt, err := sparsemt.New(1, tx, "prefix")
		assert.NilError(t, err)
		mt.SetRoot(hash)

		proof, err := mt.Proof(common.Big0)
		assert.NilError(t, err)
		assert.Equal(t, len(proof), 1)
		assert.Equal(t, proof[0], common.HexToHash("0x2"))

		proof, err = mt.Proof(common.Big1)
		assert.NilError(t, err)
		assert.Equal(t, len(proof), 1)
		assert.Equal(t, proof[0], common.HexToHash("0x1"))
	})
	t.Run(`ProofExt depth 2`, func(t *testing.T) {
		t.Parallel()
		mt := createTestingTree(t, "prefix")

		proof, err := mt.Proof(common.Big0)
		assert.NilError(t, err)
		assert.Equal(t, len(proof), 2)
		assert.Equal(t, proof[0], common.HexToHash("0x2e174c10e159ea99b867ce3205125c24a42d128804e4070ed6fcc8cc98166aa0"))
		assert.Equal(t, proof[1], common.HexToHash("0x2"))

		proof, err = mt.Proof(common.Big1)
		assert.NilError(t, err)
		assert.Equal(t, len(proof), 2)
		assert.Equal(t, proof[0], common.HexToHash("0x2e174c10e159ea99b867ce3205125c24a42d128804e4070ed6fcc8cc98166aa0"))
		assert.Equal(t, proof[1], common.HexToHash("0x1"))

		proof, err = mt.Proof(common.Big2)
		assert.NilError(t, err)
		assert.Equal(t, len(proof), 2)
		assert.Equal(t, proof[0], common.HexToHash("0xe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0"))
		assert.Equal(t, proof[1], common.HexToHash("0x4"))

		proof, err = mt.Proof(big.NewInt(3))
		assert.NilError(t, err)
		assert.Equal(t, len(proof), 2)
		assert.Equal(t, proof[0], common.HexToHash("0xe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0"))
		assert.Equal(t, proof[1], common.HexToHash("0x3"))
	})
	t.Run(`Leaf with depth 2`, func(t *testing.T) {
		t.Parallel()
		mt := createTestingTree(t, "prefix")

		hash, err := mt.Leaf(common.Big0)
		assert.NilError(t, err)
		assert.Equal(t, hash.String(), "0x0000000000000000000000000000000000000000000000000000000000000001")

		hash, err = mt.Leaf(common.Big1)
		assert.NilError(t, err)
		assert.Equal(t, hash.String(), "0x0000000000000000000000000000000000000000000000000000000000000002")

		hash, err = mt.Leaf(common.Big2)
		assert.NilError(t, err)
		assert.Equal(t, hash.String(), "0x0000000000000000000000000000000000000000000000000000000000000003")

		hash, err = mt.Leaf(common.Big3)
		assert.NilError(t, err)
		assert.Equal(t, hash.String(), "0x0000000000000000000000000000000000000000000000000000000000000004")
	})
	t.Run(`update leaf`, func(t *testing.T) {
		t.Parallel()
		mt := createTestingTree(t, "prefix")

		assert.NilError(t, mt.SetLeaf(common.Big3, common.HexToHash("0x9")))
		hash, _ := mt.Leaf(common.Big0)
		assert.Equal(t, hash, common.HexToHash("0x1"))
		hash, _ = mt.Leaf(common.Big1)
		assert.Equal(t, hash, common.HexToHash("0x2"))
		hash, _ = mt.Leaf(common.Big2)
		assert.Equal(t, hash, common.HexToHash("0x3"))
		hash, _ = mt.Leaf(common.Big3)
		assert.Equal(t, hash, common.HexToHash("0x9"))

		assert.NilError(t, mt.SetLeaf(common.Big0, common.HexToHash("0x8")))
		hash, _ = mt.Leaf(common.Big0)
		assert.Equal(t, hash, common.HexToHash("0x8"))
		hash, _ = mt.Leaf(common.Big1)
		assert.Equal(t, hash, common.HexToHash("0x2"))
		hash, _ = mt.Leaf(common.Big2)
		assert.Equal(t, hash, common.HexToHash("0x3"))
		hash, _ = mt.Leaf(common.Big3)
		assert.Equal(t, hash, common.HexToHash("0x9"))

		assert.NilError(t, mt.SetLeaf(common.Big1, common.HexToHash("0x7")))
		hash, _ = mt.Leaf(common.Big0)
		assert.Equal(t, hash, common.HexToHash("0x8"))
		hash, _ = mt.Leaf(common.Big1)
		assert.Equal(t, hash, common.HexToHash("0x7"))
		hash, _ = mt.Leaf(common.Big2)
		assert.Equal(t, hash, common.HexToHash("0x3"))
		hash, _ = mt.Leaf(common.Big3)
		assert.Equal(t, hash, common.HexToHash("0x9"))

		assert.NilError(t, mt.SetLeaf(common.Big2, common.HexToHash("0x6")))
		hash, _ = mt.Leaf(common.Big0)
		assert.Equal(t, hash, common.HexToHash("0x8"))
		hash, _ = mt.Leaf(common.Big1)
		assert.Equal(t, hash, common.HexToHash("0x7"))
		hash, _ = mt.Leaf(common.Big2)
		assert.Equal(t, hash, common.HexToHash("0x6"))
		hash, _ = mt.Leaf(common.Big3)
		assert.Equal(t, hash, common.HexToHash("0x9"))
	})
	t.Run(`create 2 smt with same storage`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		smt0, err := sparsemt.New(1, tx, "prefix")
		assert.NilError(t, err)
		smt1, err := sparsemt.New(1, tx, "prefix")
		assert.NilError(t, err)

		assert.Equal(t, smt0.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
		assert.NilError(t, smt0.SetLeaf(common.Big0, common.HexToHash("0x1")))
		assert.Equal(t, smt0.Root().String(), "0xada5013122d395ba3c54772283fb069b10426056ef8ca54750cb9bb552a59e7d")

		assert.Equal(t, smt1.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
		assert.NilError(t, smt1.SetLeaf(common.Big1, common.HexToHash("0x1")))
		assert.Equal(t, smt1.Root().String(), "0xa6eef7e35abe7026729641147f7915573c7e97b47efa546f5f6e3230263bcb49")
	})
	t.Run(`count leaves`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		smt, err := sparsemt.New(2, tx, "")
		assert.NilError(t, err)
		assert.Equal(t, smt.CountLeaves().String(), "4")

		smt, err = sparsemt.New(10, tx, "")
		assert.NilError(t, err)
		assert.Equal(t, smt.CountLeaves().String(), "1024")
	})
	t.Run(`depth`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		smt, err := sparsemt.New(2, tx, "")
		assert.NilError(t, err)
		assert.Equal(t, smt.Depth(), uint(2))

		smt, err = sparsemt.New(20, tx, "")
		assert.NilError(t, err)
		assert.Equal(t, smt.Depth(), uint(20))
	})
	t.Run(`proof of tree with 1 leaf`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		smt, _ := sparsemt.New(0, tx, "")
		proof, err := smt.Proof(common.Big0)
		assert.NilError(t, err)
		assert.Assert(t, proof == nil)
	})
	t.Run(`proof of tree with 2 leaf`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		smt, _ := sparsemt.New(1, tx, "")
		proof, err := smt.Proof(common.Big0)
		assert.NilError(t, err)
		assert.Assert(t, proof != nil)
		assert.Equal(t, len(proof), 1)
		assert.Equal(t, proof[0].String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	})
	t.Run(`proof of unexistent leaf`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		smt, _ := sparsemt.New(1, tx, "")
		_, err := smt.Proof(common.Big3)
		assert.Error(t, err, "Tree:Proof:out of bounds")
	})
}

func createTestingTree(t *testing.T, prefix string) merkletree.MerkleTree {
	service := memory.New()
	tx := service.NewTransaction()

	mt, err := sparsemt.New(2, tx, prefix)
	assert.NilError(t, err)
	err = mt.SetLeaf(common.Big0, common.HexToHash("0x1"))
	assert.NilError(t, err)
	err = mt.SetLeaf(common.Big1, common.HexToHash("0x2"))
	assert.NilError(t, err)
	err = mt.SetLeaf(common.Big2, common.HexToHash("0x3"))
	assert.NilError(t, err)
	err = mt.SetLeaf(common.Big3, common.HexToHash("0x4"))
	assert.NilError(t, err)
	return mt
}

func TestTreeWithDepth0(t *testing.T) {
	t.Parallel()
	t.Run(`instantiate`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		_, err := sparsemt.New(0, tx, "")
		assert.NilError(t, err)
	})
	t.Run(`root`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tree, _ := sparsemt.New(0, tx, "")
		assert.Equal(t, tree.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	})
	t.Run(`path`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tree, _ := sparsemt.New(0, tx, "")
		_, err := tree.Path(common.Big1)
		assert.Error(t, err, "Tree:Path:out of bounds")
		path, err := tree.Path(common.Big0)
		assert.NilError(t, err)
		assert.Equal(t, len(path), 0)
	})
	t.Run(`adding an out of bounds leaf`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tree, _ := sparsemt.New(0, tx, "")
		assert.Error(t, tree.SetLeaf(common.Big1, common.HexToHash("0x1")), "Tree:Proof:out of bounds")
	})
	t.Run(`root adding a leaf`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tree, _ := sparsemt.New(0, tx, "")
		assert.NilError(t, tree.SetLeaf(common.Big0, common.HexToHash("0x1")))
		assert.Equal(t, tree.Root().String(), "0x0000000000000000000000000000000000000000000000000000000000000001")
		leaf, err := tree.Leaf(common.Big0)
		assert.NilError(t, err)
		assert.Equal(t, leaf.String(), "0x0000000000000000000000000000000000000000000000000000000000000001")
	})
	t.Run(`get leaf out of bounds`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		tree, _ := sparsemt.New(0, tx, "")
		_, err := tree.Leaf(common.Big1)
		assert.Error(t, err, "Tree:Path:out of bounds")
	})
}
