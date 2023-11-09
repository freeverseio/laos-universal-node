// Code generated by MockGen. DO NOT EDIT.
// Source: internal/scan/scan.go
//
// Generated by this command:
//
//	mockgen -source=internal/scan/scan.go -destination=internal/scan/mock/scan.go -package=mock
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	big "math/big"
	reflect "reflect"

	ethereum "github.com/ethereum/go-ethereum"
	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	scan "github.com/freeverseio/laos-universal-node/internal/scan"
	gomock "go.uber.org/mock/gomock"
)

// MockEthClient is a mock of EthClient interface.
type MockEthClient struct {
	ctrl     *gomock.Controller
	recorder *MockEthClientMockRecorder
}

// MockEthClientMockRecorder is the mock recorder for MockEthClient.
type MockEthClientMockRecorder struct {
	mock *MockEthClient
}

// NewMockEthClient creates a new mock instance.
func NewMockEthClient(ctrl *gomock.Controller) *MockEthClient {
	mock := &MockEthClient{ctrl: ctrl}
	mock.recorder = &MockEthClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEthClient) EXPECT() *MockEthClientMockRecorder {
	return m.recorder
}

// BlockByHash mocks base method.
func (m *MockEthClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockByHash", ctx, hash)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockByHash indicates an expected call of BlockByHash.
func (mr *MockEthClientMockRecorder) BlockByHash(ctx, hash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockByHash", reflect.TypeOf((*MockEthClient)(nil).BlockByHash), ctx, hash)
}

// BlockByNumber mocks base method.
func (m *MockEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockByNumber", ctx, number)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockByNumber indicates an expected call of BlockByNumber.
func (mr *MockEthClientMockRecorder) BlockByNumber(ctx, number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockByNumber", reflect.TypeOf((*MockEthClient)(nil).BlockByNumber), ctx, number)
}

// BlockNumber mocks base method.
func (m *MockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockNumber", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockNumber indicates an expected call of BlockNumber.
func (mr *MockEthClientMockRecorder) BlockNumber(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockNumber", reflect.TypeOf((*MockEthClient)(nil).BlockNumber), ctx)
}

// CallContract mocks base method.
func (m *MockEthClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CallContract", ctx, call, blockNumber)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CallContract indicates an expected call of CallContract.
func (mr *MockEthClientMockRecorder) CallContract(ctx, call, blockNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CallContract", reflect.TypeOf((*MockEthClient)(nil).CallContract), ctx, call, blockNumber)
}

// ChainID mocks base method.
func (m *MockEthClient) ChainID(ctx context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChainID", ctx)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChainID indicates an expected call of ChainID.
func (mr *MockEthClientMockRecorder) ChainID(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChainID", reflect.TypeOf((*MockEthClient)(nil).ChainID), ctx)
}

// Close mocks base method.
func (m *MockEthClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockEthClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockEthClient)(nil).Close))
}

// CodeAt mocks base method.
func (m *MockEthClient) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CodeAt", ctx, contract, blockNumber)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CodeAt indicates an expected call of CodeAt.
func (mr *MockEthClientMockRecorder) CodeAt(ctx, contract, blockNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CodeAt", reflect.TypeOf((*MockEthClient)(nil).CodeAt), ctx, contract, blockNumber)
}

// EstimateGas mocks base method.
func (m *MockEthClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimateGas", ctx, call)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EstimateGas indicates an expected call of EstimateGas.
func (mr *MockEthClientMockRecorder) EstimateGas(ctx, call any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimateGas", reflect.TypeOf((*MockEthClient)(nil).EstimateGas), ctx, call)
}

// FilterLogs mocks base method.
func (m *MockEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterLogs", ctx, q)
	ret0, _ := ret[0].([]types.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterLogs indicates an expected call of FilterLogs.
func (mr *MockEthClientMockRecorder) FilterLogs(ctx, q any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterLogs", reflect.TypeOf((*MockEthClient)(nil).FilterLogs), ctx, q)
}

// HeaderByHash mocks base method.
func (m *MockEthClient) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeaderByHash", ctx, hash)
	ret0, _ := ret[0].(*types.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HeaderByHash indicates an expected call of HeaderByHash.
func (mr *MockEthClientMockRecorder) HeaderByHash(ctx, hash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeaderByHash", reflect.TypeOf((*MockEthClient)(nil).HeaderByHash), ctx, hash)
}

// HeaderByNumber mocks base method.
func (m *MockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeaderByNumber", ctx, number)
	ret0, _ := ret[0].(*types.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HeaderByNumber indicates an expected call of HeaderByNumber.
func (mr *MockEthClientMockRecorder) HeaderByNumber(ctx, number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeaderByNumber", reflect.TypeOf((*MockEthClient)(nil).HeaderByNumber), ctx, number)
}

// PendingCallContract mocks base method.
func (m *MockEthClient) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingCallContract", ctx, call)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingCallContract indicates an expected call of PendingCallContract.
func (mr *MockEthClientMockRecorder) PendingCallContract(ctx, call any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingCallContract", reflect.TypeOf((*MockEthClient)(nil).PendingCallContract), ctx, call)
}

// PendingCodeAt mocks base method.
func (m *MockEthClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingCodeAt", ctx, account)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingCodeAt indicates an expected call of PendingCodeAt.
func (mr *MockEthClientMockRecorder) PendingCodeAt(ctx, account any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingCodeAt", reflect.TypeOf((*MockEthClient)(nil).PendingCodeAt), ctx, account)
}

// PendingNonceAt mocks base method.
func (m *MockEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingNonceAt", ctx, account)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingNonceAt indicates an expected call of PendingNonceAt.
func (mr *MockEthClientMockRecorder) PendingNonceAt(ctx, account any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingNonceAt", reflect.TypeOf((*MockEthClient)(nil).PendingNonceAt), ctx, account)
}

// SendTransaction mocks base method.
func (m *MockEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTransaction", ctx, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendTransaction indicates an expected call of SendTransaction.
func (mr *MockEthClientMockRecorder) SendTransaction(ctx, tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTransaction", reflect.TypeOf((*MockEthClient)(nil).SendTransaction), ctx, tx)
}

// SubscribeFilterLogs mocks base method.
func (m *MockEthClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeFilterLogs", ctx, q, ch)
	ret0, _ := ret[0].(ethereum.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeFilterLogs indicates an expected call of SubscribeFilterLogs.
func (mr *MockEthClientMockRecorder) SubscribeFilterLogs(ctx, q, ch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeFilterLogs", reflect.TypeOf((*MockEthClient)(nil).SubscribeFilterLogs), ctx, q, ch)
}

// SubscribeNewHead mocks base method.
func (m *MockEthClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeNewHead", ctx, ch)
	ret0, _ := ret[0].(ethereum.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeNewHead indicates an expected call of SubscribeNewHead.
func (mr *MockEthClientMockRecorder) SubscribeNewHead(ctx, ch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeNewHead", reflect.TypeOf((*MockEthClient)(nil).SubscribeNewHead), ctx, ch)
}

// SuggestGasPrice mocks base method.
func (m *MockEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuggestGasPrice", ctx)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuggestGasPrice indicates an expected call of SuggestGasPrice.
func (mr *MockEthClientMockRecorder) SuggestGasPrice(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuggestGasPrice", reflect.TypeOf((*MockEthClient)(nil).SuggestGasPrice), ctx)
}

// SuggestGasTipCap mocks base method.
func (m *MockEthClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuggestGasTipCap", ctx)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuggestGasTipCap indicates an expected call of SuggestGasTipCap.
func (mr *MockEthClientMockRecorder) SuggestGasTipCap(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuggestGasTipCap", reflect.TypeOf((*MockEthClient)(nil).SuggestGasTipCap), ctx)
}

// SyncProgress mocks base method.
func (m *MockEthClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncProgress", ctx)
	ret0, _ := ret[0].(*ethereum.SyncProgress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncProgress indicates an expected call of SyncProgress.
func (mr *MockEthClientMockRecorder) SyncProgress(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncProgress", reflect.TypeOf((*MockEthClient)(nil).SyncProgress), ctx)
}

// TransactionByHash mocks base method.
func (m *MockEthClient) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionByHash", ctx, txHash)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// TransactionByHash indicates an expected call of TransactionByHash.
func (mr *MockEthClientMockRecorder) TransactionByHash(ctx, txHash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionByHash", reflect.TypeOf((*MockEthClient)(nil).TransactionByHash), ctx, txHash)
}

// TransactionCount mocks base method.
func (m *MockEthClient) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionCount", ctx, blockHash)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransactionCount indicates an expected call of TransactionCount.
func (mr *MockEthClientMockRecorder) TransactionCount(ctx, blockHash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionCount", reflect.TypeOf((*MockEthClient)(nil).TransactionCount), ctx, blockHash)
}

// TransactionInBlock mocks base method.
func (m *MockEthClient) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionInBlock", ctx, blockHash, index)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransactionInBlock indicates an expected call of TransactionInBlock.
func (mr *MockEthClientMockRecorder) TransactionInBlock(ctx, blockHash, index any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionInBlock", reflect.TypeOf((*MockEthClient)(nil).TransactionInBlock), ctx, blockHash, index)
}

// TransactionReceipt mocks base method.
func (m *MockEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionReceipt", ctx, txHash)
	ret0, _ := ret[0].(*types.Receipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransactionReceipt indicates an expected call of TransactionReceipt.
func (mr *MockEthClientMockRecorder) TransactionReceipt(ctx, txHash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionReceipt", reflect.TypeOf((*MockEthClient)(nil).TransactionReceipt), ctx, txHash)
}

// MockEvent is a mock of Event interface.
type MockEvent struct {
	ctrl     *gomock.Controller
	recorder *MockEventMockRecorder
}

// MockEventMockRecorder is the mock recorder for MockEvent.
type MockEventMockRecorder struct {
	mock *MockEvent
}

// NewMockEvent creates a new mock instance.
func NewMockEvent(ctrl *gomock.Controller) *MockEvent {
	mock := &MockEvent{ctrl: ctrl}
	mock.recorder = &MockEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEvent) EXPECT() *MockEventMockRecorder {
	return m.recorder
}

// MockScanner is a mock of Scanner interface.
type MockScanner struct {
	ctrl     *gomock.Controller
	recorder *MockScannerMockRecorder
}

// MockScannerMockRecorder is the mock recorder for MockScanner.
type MockScannerMockRecorder struct {
	mock *MockScanner
}

// NewMockScanner creates a new mock instance.
func NewMockScanner(ctrl *gomock.Controller) *MockScanner {
	mock := &MockScanner{ctrl: ctrl}
	mock.recorder = &MockScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScanner) EXPECT() *MockScannerMockRecorder {
	return m.recorder
}

// ScanEvents mocks base method.
func (m *MockScanner) ScanEvents(ctx context.Context, fromBlock, toBlock *big.Int, contracts []string) ([]scan.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScanEvents", ctx, fromBlock, toBlock, contracts)
	ret0, _ := ret[0].([]scan.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScanEvents indicates an expected call of ScanEvents.
func (mr *MockScannerMockRecorder) ScanEvents(ctx, fromBlock, toBlock, contracts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScanEvents", reflect.TypeOf((*MockScanner)(nil).ScanEvents), ctx, fromBlock, toBlock, contracts)
}

// ScanNewUniversalEvents mocks base method.
func (m *MockScanner) ScanNewUniversalEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]scan.ERC721UniversalContract, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScanNewUniversalEvents", ctx, fromBlock, toBlock)
	ret0, _ := ret[0].([]scan.ERC721UniversalContract)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScanNewUniversalEvents indicates an expected call of ScanNewUniversalEvents.
func (mr *MockScannerMockRecorder) ScanNewUniversalEvents(ctx, fromBlock, toBlock any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScanNewUniversalEvents", reflect.TypeOf((*MockScanner)(nil).ScanNewUniversalEvents), ctx, fromBlock, toBlock)
}
