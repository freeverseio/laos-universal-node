// Code generated by MockGen. DO NOT EDIT.
// Source: internal/state/enumerated/tree.go
//
// Generated by this command:
//
//	mockgen -source=internal/state/enumerated/tree.go -destination=internal/state/enumerated/mock/tree.go -package=mock
//
// Package mock is a generated GoMock package.
package mock

import (
	big "math/big"
	reflect "reflect"

	common "github.com/ethereum/go-ethereum/common"
	model "github.com/freeverseio/laos-universal-node/internal/platform/model"
	gomock "go.uber.org/mock/gomock"
)

// MockTree is a mock of Tree interface.
type MockTree struct {
	ctrl     *gomock.Controller
	recorder *MockTreeMockRecorder
}

// MockTreeMockRecorder is the mock recorder for MockTree.
type MockTreeMockRecorder struct {
	mock *MockTree
}

// NewMockTree creates a new mock instance.
func NewMockTree(ctrl *gomock.Controller) *MockTree {
	mock := &MockTree{ctrl: ctrl}
	mock.recorder = &MockTreeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTree) EXPECT() *MockTreeMockRecorder {
	return m.recorder
}

// Checkout mocks base method.
func (m *MockTree) Checkout(blockNumber int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checkout", blockNumber)
	ret0, _ := ret[0].(error)
	return ret0
}

// Checkout indicates an expected call of Checkout.
func (mr *MockTreeMockRecorder) Checkout(blockNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checkout", reflect.TypeOf((*MockTree)(nil).Checkout), blockNumber)
}

// FindBlockWithTag mocks base method.
func (m *MockTree) FindBlockWithTag(blockNumber int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBlockWithTag", blockNumber)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBlockWithTag indicates an expected call of FindBlockWithTag.
func (mr *MockTreeMockRecorder) FindBlockWithTag(blockNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBlockWithTag", reflect.TypeOf((*MockTree)(nil).FindBlockWithTag), blockNumber)
}

// Mint mocks base method.
func (m *MockTree) Mint(tokenId *big.Int, owner common.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mint", tokenId, owner)
	ret0, _ := ret[0].(error)
	return ret0
}

// Mint indicates an expected call of Mint.
func (mr *MockTreeMockRecorder) Mint(tokenId, owner any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mint", reflect.TypeOf((*MockTree)(nil).Mint), tokenId, owner)
}

// Root mocks base method.
func (m *MockTree) Root() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Root")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// Root indicates an expected call of Root.
func (mr *MockTreeMockRecorder) Root() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Root", reflect.TypeOf((*MockTree)(nil).Root))
}

// TagRoot mocks base method.
func (m *MockTree) TagRoot(blockNumber int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TagRoot", blockNumber)
	ret0, _ := ret[0].(error)
	return ret0
}

// TagRoot indicates an expected call of TagRoot.
func (mr *MockTreeMockRecorder) TagRoot(blockNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TagRoot", reflect.TypeOf((*MockTree)(nil).TagRoot), blockNumber)
}

// TokensOf mocks base method.
func (m *MockTree) TokensOf(owner common.Address) ([]big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TokensOf", owner)
	ret0, _ := ret[0].([]big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TokensOf indicates an expected call of TokensOf.
func (mr *MockTreeMockRecorder) TokensOf(owner any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TokensOf", reflect.TypeOf((*MockTree)(nil).TokensOf), owner)
}

// Transfer mocks base method.
func (m *MockTree) Transfer(minted bool, eventTransfer *model.ERC721Transfer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", minted, eventTransfer)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transfer indicates an expected call of Transfer.
func (mr *MockTreeMockRecorder) Transfer(minted, eventTransfer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockTree)(nil).Transfer), minted, eventTransfer)
}
