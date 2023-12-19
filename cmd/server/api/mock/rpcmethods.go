// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/server/api/rpcmethods.go
//
// Generated by this command:
//
//	mockgen -source=cmd/server/api/rpcmethods.go -destination=cmd/server/api/mock/rpcmethods.go -package=mock
//
// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	api "github.com/freeverseio/laos-universal-node/cmd/server/api"
	gomock "go.uber.org/mock/gomock"
)

// MockRPCMethodManager is a mock of RPCMethodManager interface.
type MockRPCMethodManager struct {
	ctrl     *gomock.Controller
	recorder *MockRPCMethodManagerMockRecorder
}

// MockRPCMethodManagerMockRecorder is the mock recorder for MockRPCMethodManager.
type MockRPCMethodManagerMockRecorder struct {
	mock *MockRPCMethodManager
}

// NewMockRPCMethodManager creates a new mock instance.
func NewMockRPCMethodManager(ctrl *gomock.Controller) *MockRPCMethodManager {
	mock := &MockRPCMethodManager{ctrl: ctrl}
	mock.recorder = &MockRPCMethodManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRPCMethodManager) EXPECT() *MockRPCMethodManagerMockRecorder {
	return m.recorder
}

// CheckBlockNumberFromResponseFromHashCalls mocks base method.
func (m *MockRPCMethodManager) CheckBlockNumberFromResponseFromHashCalls(resp *api.RPCResponse, method api.RPCMethod, blockNumberUnode string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckBlockNumberFromResponseFromHashCalls", resp, method, blockNumberUnode)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckBlockNumberFromResponseFromHashCalls indicates an expected call of CheckBlockNumberFromResponseFromHashCalls.
func (mr *MockRPCMethodManagerMockRecorder) CheckBlockNumberFromResponseFromHashCalls(resp, method, blockNumberUnode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckBlockNumberFromResponseFromHashCalls", reflect.TypeOf((*MockRPCMethodManager)(nil).CheckBlockNumberFromResponseFromHashCalls), resp, method, blockNumberUnode)
}

// HasRPCMethodWithBlocknumber mocks base method.
func (m *MockRPCMethodManager) HasRPCMethodWithBlocknumber(methodName string) (api.RPCMethod, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasRPCMethodWithBlocknumber", methodName)
	ret0, _ := ret[0].(api.RPCMethod)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// HasRPCMethodWithBlocknumber indicates an expected call of HasRPCMethodWithBlocknumber.
func (mr *MockRPCMethodManagerMockRecorder) HasRPCMethodWithBlocknumber(methodName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasRPCMethodWithBlocknumber", reflect.TypeOf((*MockRPCMethodManager)(nil).HasRPCMethodWithBlocknumber), methodName)
}

// HasRPCMethodWithHash mocks base method.
func (m *MockRPCMethodManager) HasRPCMethodWithHash(methodName string) (api.RPCMethod, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasRPCMethodWithHash", methodName)
	ret0, _ := ret[0].(api.RPCMethod)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// HasRPCMethodWithHash indicates an expected call of HasRPCMethodWithHash.
func (mr *MockRPCMethodManagerMockRecorder) HasRPCMethodWithHash(methodName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasRPCMethodWithHash", reflect.TypeOf((*MockRPCMethodManager)(nil).HasRPCMethodWithHash), methodName)
}

// ReplaceBlockTag mocks base method.
func (m *MockRPCMethodManager) ReplaceBlockTag(req *api.JSONRPCRequest, method api.RPCMethod, blockNumberUnode string) (*api.JSONRPCRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceBlockTag", req, method, blockNumberUnode)
	ret0, _ := ret[0].(*api.JSONRPCRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceBlockTag indicates an expected call of ReplaceBlockTag.
func (mr *MockRPCMethodManagerMockRecorder) ReplaceBlockTag(req, method, blockNumberUnode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceBlockTag", reflect.TypeOf((*MockRPCMethodManager)(nil).ReplaceBlockTag), req, method, blockNumberUnode)
}
