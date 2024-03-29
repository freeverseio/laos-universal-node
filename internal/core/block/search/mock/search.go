// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/block/search/search.go
//
// Generated by this command:
//
//	mockgen -source=internal/core/block/search/search.go -destination=internal/core/block/search/mock/search.go -package=mock
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSearch is a mock of Search interface.
type MockSearch struct {
	ctrl     *gomock.Controller
	recorder *MockSearchMockRecorder
}

// MockSearchMockRecorder is the mock recorder for MockSearch.
type MockSearchMockRecorder struct {
	mock *MockSearch
}

// NewMockSearch creates a new mock instance.
func NewMockSearch(ctrl *gomock.Controller) *MockSearch {
	mock := &MockSearch{ctrl: ctrl}
	mock.recorder = &MockSearchMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearch) EXPECT() *MockSearchMockRecorder {
	return m.recorder
}

// GetEvolutionBlockByTimestamp mocks base method.
func (m *MockSearch) GetEvolutionBlockByTimestamp(ctx context.Context, targetTimestamp uint64, startingPoint ...uint64) (uint64, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, targetTimestamp}
	for _, a := range startingPoint {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetEvolutionBlockByTimestamp", varargs...)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvolutionBlockByTimestamp indicates an expected call of GetEvolutionBlockByTimestamp.
func (mr *MockSearchMockRecorder) GetEvolutionBlockByTimestamp(ctx, targetTimestamp any, startingPoint ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, targetTimestamp}, startingPoint...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvolutionBlockByTimestamp", reflect.TypeOf((*MockSearch)(nil).GetEvolutionBlockByTimestamp), varargs...)
}

// GetOwnershipBlockByTimestamp mocks base method.
func (m *MockSearch) GetOwnershipBlockByTimestamp(ctx context.Context, targetTimestamp uint64, startingPoint ...uint64) (uint64, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, targetTimestamp}
	for _, a := range startingPoint {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetOwnershipBlockByTimestamp", varargs...)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOwnershipBlockByTimestamp indicates an expected call of GetOwnershipBlockByTimestamp.
func (mr *MockSearchMockRecorder) GetOwnershipBlockByTimestamp(ctx, targetTimestamp any, startingPoint ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, targetTimestamp}, startingPoint...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwnershipBlockByTimestamp", reflect.TypeOf((*MockSearch)(nil).GetOwnershipBlockByTimestamp), varargs...)
}
