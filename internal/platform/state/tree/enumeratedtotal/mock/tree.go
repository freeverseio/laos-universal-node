// Code generated by MockGen. DO NOT EDIT.
// Source: internal/platform/state/tree/enumeratedtotal/tree.go
//
// Generated by this command:
//
//	mockgen -source=internal/platform/state/tree/enumeratedtotal/tree.go -destination=internal/platform/state/tree/enumeratedtotal/mock/tree.go -package=mock
//
// Package mock is a generated GoMock package.
package mock

import (
	big "math/big"
	reflect "reflect"

	common "github.com/ethereum/go-ethereum/common"
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

// Burn mocks base method.
func (m *MockTree) Burn(idx int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Burn", idx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Burn indicates an expected call of Burn.
func (mr *MockTreeMockRecorder) Burn(idx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Burn", reflect.TypeOf((*MockTree)(nil).Burn), idx)
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

// GetLastTaggedBlock mocks base method.
func (m *MockTree) GetLastTaggedBlock() (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastTaggedBlock")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastTaggedBlock indicates an expected call of GetLastTaggedBlock.
func (mr *MockTreeMockRecorder) GetLastTaggedBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastTaggedBlock", reflect.TypeOf((*MockTree)(nil).GetLastTaggedBlock))
}

// Mint mocks base method.
func (m *MockTree) Mint(tokenId *big.Int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mint", tokenId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Mint indicates an expected call of Mint.
func (mr *MockTreeMockRecorder) Mint(tokenId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mint", reflect.TypeOf((*MockTree)(nil).Mint), tokenId)
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

// TokenByIndex mocks base method.
func (m *MockTree) TokenByIndex(idx int) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TokenByIndex", idx)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TokenByIndex indicates an expected call of TokenByIndex.
func (mr *MockTreeMockRecorder) TokenByIndex(idx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TokenByIndex", reflect.TypeOf((*MockTree)(nil).TokenByIndex), idx)
}

// TotalSupply mocks base method.
func (m *MockTree) TotalSupply() (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TotalSupply")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TotalSupply indicates an expected call of TotalSupply.
func (mr *MockTreeMockRecorder) TotalSupply() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TotalSupply", reflect.TypeOf((*MockTree)(nil).TotalSupply))
}
