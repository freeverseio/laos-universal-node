// Code generated by MockGen. DO NOT EDIT.
// Source: internal/platform/storage/storage.go
//
// Generated by this command:
//
//	mockgen -source=internal/platform/storage/storage.go -destination=internal/platform/storage/mock/storage.go -package=mock
//
// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	storage "github.com/freeverseio/laos-universal-node/internal/platform/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockTx is a mock of Tx interface.
type MockTx struct {
	ctrl     *gomock.Controller
	recorder *MockTxMockRecorder
}

// MockTxMockRecorder is the mock recorder for MockTx.
type MockTxMockRecorder struct {
	mock *MockTx
}

// NewMockTx creates a new mock instance.
func NewMockTx(ctrl *gomock.Controller) *MockTx {
	mock := &MockTx{ctrl: ctrl}
	mock.recorder = &MockTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTx) EXPECT() *MockTxMockRecorder {
	return m.recorder
}

// ClearAll mocks base method.
func (m *MockTx) ClearAll() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearAll")
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearAll indicates an expected call of ClearAll.
func (mr *MockTxMockRecorder) ClearAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearAll", reflect.TypeOf((*MockTx)(nil).ClearAll))
}

// Commit mocks base method.
func (m *MockTx) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTxMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTx)(nil).Commit))
}

// Delete mocks base method.
func (m *MockTx) Delete(key []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTxMockRecorder) Delete(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTx)(nil).Delete), key)
}

// Discard mocks base method.
func (m *MockTx) Discard() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Discard")
}

// Discard indicates an expected call of Discard.
func (mr *MockTxMockRecorder) Discard() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Discard", reflect.TypeOf((*MockTx)(nil).Discard))
}

// Get mocks base method.
func (m *MockTx) Get(key []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockTxMockRecorder) Get(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTx)(nil).Get), key)
}

// GetKeysWithPrefix mocks base method.
func (m *MockTx) GetKeysWithPrefix(prefix []byte, reverse ...bool) [][]byte {
	m.ctrl.T.Helper()
	varargs := []any{prefix}
	for _, a := range reverse {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetKeysWithPrefix", varargs...)
	ret0, _ := ret[0].([][]byte)
	return ret0
}

// GetKeysWithPrefix indicates an expected call of GetKeysWithPrefix.
func (mr *MockTxMockRecorder) GetKeysWithPrefix(prefix any, reverse ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{prefix}, reverse...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeysWithPrefix", reflect.TypeOf((*MockTx)(nil).GetKeysWithPrefix), varargs...)
}

// Len mocks base method.
func (m *MockTx) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *MockTxMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockTx)(nil).Len))
}

// Set mocks base method.
func (m *MockTx) Set(key, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockTxMockRecorder) Set(key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockTx)(nil).Set), key, value)
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockService) Get(key []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockServiceMockRecorder) Get(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockService)(nil).Get), key)
}

// GetKeysWithPrefix mocks base method.
func (m *MockService) GetKeysWithPrefix(prefix []byte, reverse ...bool) ([][]byte, error) {
	m.ctrl.T.Helper()
	varargs := []any{prefix}
	for _, a := range reverse {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetKeysWithPrefix", varargs...)
	ret0, _ := ret[0].([][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeysWithPrefix indicates an expected call of GetKeysWithPrefix.
func (mr *MockServiceMockRecorder) GetKeysWithPrefix(prefix any, reverse ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{prefix}, reverse...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeysWithPrefix", reflect.TypeOf((*MockService)(nil).GetKeysWithPrefix), varargs...)
}

// NewTransaction mocks base method.
func (m *MockService) NewTransaction() storage.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTransaction")
	ret0, _ := ret[0].(storage.Tx)
	return ret0
}

// NewTransaction indicates an expected call of NewTransaction.
func (mr *MockServiceMockRecorder) NewTransaction() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTransaction", reflect.TypeOf((*MockService)(nil).NewTransaction))
}

// Set mocks base method.
func (m *MockService) Set(key, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockServiceMockRecorder) Set(key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockService)(nil).Set), key, value)
}
