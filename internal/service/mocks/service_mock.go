// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/Makovey/gophermart/internal/transport/http/model"
	gomock "github.com/golang/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// LoginUser mocks base method.
func (m *MockUserService) LoginUser(ctx context.Context, request model.AuthRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", ctx, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockUserServiceMockRecorder) LoginUser(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockUserService)(nil).LoginUser), ctx, request)
}

// RegisterNewUser mocks base method.
func (m *MockUserService) RegisterNewUser(ctx context.Context, request model.AuthRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterNewUser", ctx, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterNewUser indicates an expected call of RegisterNewUser.
func (mr *MockUserServiceMockRecorder) RegisterNewUser(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterNewUser", reflect.TypeOf((*MockUserService)(nil).RegisterNewUser), ctx, request)
}

// MockOrderService is a mock of OrderService interface.
type MockOrderService struct {
	ctrl     *gomock.Controller
	recorder *MockOrderServiceMockRecorder
}

// MockOrderServiceMockRecorder is the mock recorder for MockOrderService.
type MockOrderServiceMockRecorder struct {
	mock *MockOrderService
}

// NewMockOrderService creates a new mock instance.
func NewMockOrderService(ctrl *gomock.Controller) *MockOrderService {
	mock := &MockOrderService{ctrl: ctrl}
	mock.recorder = &MockOrderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderService) EXPECT() *MockOrderServiceMockRecorder {
	return m.recorder
}

// GetOrders mocks base method.
func (m *MockOrderService) GetOrders(ctx context.Context, userID string) ([]model.OrderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, userID)
	ret0, _ := ret[0].([]model.OrderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockOrderServiceMockRecorder) GetOrders(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockOrderService)(nil).GetOrders), ctx, userID)
}

// ProcessNewOrder mocks base method.
func (m *MockOrderService) ProcessNewOrder(ctx context.Context, userID, orderID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessNewOrder", ctx, userID, orderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessNewOrder indicates an expected call of ProcessNewOrder.
func (mr *MockOrderServiceMockRecorder) ProcessNewOrder(ctx, userID, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessNewOrder", reflect.TypeOf((*MockOrderService)(nil).ProcessNewOrder), ctx, userID, orderID)
}

// ValidateOrderID mocks base method.
func (m *MockOrderService) ValidateOrderID(orderID string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateOrderID", orderID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ValidateOrderID indicates an expected call of ValidateOrderID.
func (mr *MockOrderServiceMockRecorder) ValidateOrderID(orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateOrderID", reflect.TypeOf((*MockOrderService)(nil).ValidateOrderID), orderID)
}

// MockBalanceService is a mock of BalanceService interface.
type MockBalanceService struct {
	ctrl     *gomock.Controller
	recorder *MockBalanceServiceMockRecorder
}

// MockBalanceServiceMockRecorder is the mock recorder for MockBalanceService.
type MockBalanceServiceMockRecorder struct {
	mock *MockBalanceService
}

// NewMockBalanceService creates a new mock instance.
func NewMockBalanceService(ctrl *gomock.Controller) *MockBalanceService {
	mock := &MockBalanceService{ctrl: ctrl}
	mock.recorder = &MockBalanceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalanceService) EXPECT() *MockBalanceServiceMockRecorder {
	return m.recorder
}

// GetUsersBalance mocks base method.
func (m *MockBalanceService) GetUsersBalance(ctx context.Context, userID string) (model.BalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersBalance", ctx, userID)
	ret0, _ := ret[0].(model.BalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersBalance indicates an expected call of GetUsersBalance.
func (mr *MockBalanceServiceMockRecorder) GetUsersBalance(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersBalance", reflect.TypeOf((*MockBalanceService)(nil).GetUsersBalance), ctx, userID)
}

// WithdrawUsersBalance mocks base method.
func (m *MockBalanceService) WithdrawUsersBalance(ctx context.Context, userID string, withdraw model.WithdrawRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithdrawUsersBalance", ctx, userID, withdraw)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithdrawUsersBalance indicates an expected call of WithdrawUsersBalance.
func (mr *MockBalanceServiceMockRecorder) WithdrawUsersBalance(ctx, userID, withdraw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithdrawUsersBalance", reflect.TypeOf((*MockBalanceService)(nil).WithdrawUsersBalance), ctx, userID, withdraw)
}

// MockHistoryService is a mock of HistoryService interface.
type MockHistoryService struct {
	ctrl     *gomock.Controller
	recorder *MockHistoryServiceMockRecorder
}

// MockHistoryServiceMockRecorder is the mock recorder for MockHistoryService.
type MockHistoryServiceMockRecorder struct {
	mock *MockHistoryService
}

// NewMockHistoryService creates a new mock instance.
func NewMockHistoryService(ctrl *gomock.Controller) *MockHistoryService {
	mock := &MockHistoryService{ctrl: ctrl}
	mock.recorder = &MockHistoryServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHistoryService) EXPECT() *MockHistoryServiceMockRecorder {
	return m.recorder
}

// GetUsersWithdrawHistory mocks base method.
func (m *MockHistoryService) GetUsersWithdrawHistory(ctx context.Context, userID string) ([]model.WithdrawHistoryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersWithdrawHistory", ctx, userID)
	ret0, _ := ret[0].([]model.WithdrawHistoryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersWithdrawHistory indicates an expected call of GetUsersWithdrawHistory.
func (mr *MockHistoryServiceMockRecorder) GetUsersWithdrawHistory(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersWithdrawHistory", reflect.TypeOf((*MockHistoryService)(nil).GetUsersWithdrawHistory), ctx, userID)
}
