package sparsemt_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/laos-universal-node/internal/platform/merkletree/sparsemt"
	"github.com/freeverseio/laos-universal-node/internal/platform/storage/memory"
	"gotest.tools/assert"
)

func TestNode(t *testing.T) {
	t.Parallel()
	t.Run(`L,R null`, func(t *testing.T) {
		t.Parallel()
		node := sparsemt.Node{}
		assert.Equal(t, node.L.String(), sparsemt.Null)
		assert.Equal(t, node.R.String(), sparsemt.Null)
	})

	t.Run(`new node hash is NullHash`, func(t *testing.T) {
		t.Parallel()
		node := sparsemt.Node{}
		assert.Equal(t, node.Hash().String(), sparsemt.Null)
	})

	t.Run(`L null, R not null`, func(t *testing.T) {
		t.Parallel()
		node := sparsemt.Node{
			R: common.HexToHash("0x1"),
		}
		assert.Equal(t, node.Hash().String(), "0xa6eef7e35abe7026729641147f7915573c7e97b47efa546f5f6e3230263bcb49")
	})
	t.Run(`L not null, R null`, func(t *testing.T) {
		t.Parallel()
		node := sparsemt.Node{
			L: common.HexToHash("0x1"),
		}
		assert.Equal(t, node.Hash().String(), "0xada5013122d395ba3c54772283fb069b10426056ef8ca54750cb9bb552a59e7d")
	})
	t.Run(`L, R not null`, func(t *testing.T) {
		t.Parallel()
		node := sparsemt.Node{
			L: common.HexToHash("0x1"),
			R: common.HexToHash("0x1"),
		}
		assert.Equal(t, node.Hash().String(), "0xcc69885fda6bcc1a4ace058b4a62bf5e179ea78fd58a1ccd71c22cc9b688792f")
	})
	t.Run(`null node`, func(t *testing.T) {
		t.Parallel()
		nullHash := common.HexToHash(sparsemt.Null)
		t.Run(`Get`, func(t *testing.T) {
			service := memory.New()
			tx := service.NewTransaction()

			node, err := sparsemt.GetNode(tx, nullHash, "")
			assert.NilError(t, err)
			assert.Assert(t, node != nil)
			assert.Equal(t, node.Hash().String(), sparsemt.Null)
		})
		t.Run(`Put`, func(t *testing.T) {
			t.Parallel()
			service := memory.New()
			tx := service.NewTransaction()

			node := sparsemt.Node{}
			_, err := sparsemt.PutNode(tx, node, "")
			assert.NilError(t, err)
		})
		t.Run(`Put returns node hash`, func(t *testing.T) {
			t.Parallel()
			service := memory.New()
			tx := service.NewTransaction()

			node := sparsemt.Node{}
			hash, err := sparsemt.PutNode(tx, node, "")
			assert.NilError(t, err)
			assert.Equal(t, hash, node.Hash())
		})
	})
	t.Run(`Put`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		node := sparsemt.Node{
			L: common.HexToHash("0x1"),
		}
		_, err := sparsemt.PutNode(tx, node, "")
		assert.NilError(t, err)
		result, err := sparsemt.GetNode(tx, node.Hash(), "")
		assert.NilError(t, err)
		assert.DeepEqual(t, *result, node)
	})
	t.Run(`GET an unexistent node should return nil`, func(t *testing.T) {
		t.Parallel()
		service := memory.New()
		tx := service.NewTransaction()

		n, err := sparsemt.GetNode(tx,
			common.HexToHash("0xe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0"), "")
		assert.NilError(t, err)
		assert.Assert(t, n == nil)
	})
}
