// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/Makovey/gophermart/internal/repository/model"
	gomock "github.com/golang/mock/gomock"
	decimal "github.com/shopspring/decimal"
)

// MockGophermartRepository is a mock of GophermartRepository interface.
type MockGophermartRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGophermartRepositoryMockRecorder
}

// MockGophermartRepositoryMockRecorder is the mock recorder for MockGophermartRepository.
type MockGophermartRepositoryMockRecorder struct {
	mock *MockGophermartRepository
}

// NewMockGophermartRepository creates a new mock instance.
func NewMockGophermartRepository(ctrl *gomock.Controller) *MockGophermartRepository {
	mock := &MockGophermartRepository{ctrl: ctrl}
	mock.recorder = &MockGophermartRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGophermartRepository) EXPECT() *MockGophermartRepositoryMockRecorder {
	return m.recorder
}

// DecreaseUsersBalance mocks base method.
func (m *MockGophermartRepository) DecreaseUsersBalance(ctx context.Context, userID string, withdraw decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecreaseUsersBalance", ctx, userID, withdraw)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecreaseUsersBalance indicates an expected call of DecreaseUsersBalance.
func (mr *MockGophermartRepositoryMockRecorder) DecreaseUsersBalance(ctx, userID, withdraw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecreaseUsersBalance", reflect.TypeOf((*MockGophermartRepository)(nil).DecreaseUsersBalance), ctx, userID, withdraw)
}

// FetchNewOrdersToChan mocks base method.
func (m *MockGophermartRepository) FetchNewOrdersToChan(ctx context.Context, ordersCh chan<- model.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchNewOrdersToChan", ctx, ordersCh)
	ret0, _ := ret[0].(error)
	return ret0
}

// FetchNewOrdersToChan indicates an expected call of FetchNewOrdersToChan.
func (mr *MockGophermartRepositoryMockRecorder) FetchNewOrdersToChan(ctx, ordersCh interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchNewOrdersToChan", reflect.TypeOf((*MockGophermartRepository)(nil).FetchNewOrdersToChan), ctx, ordersCh)
}

// GetOrderByID mocks base method.
func (m *MockGophermartRepository) GetOrderByID(ctx context.Context, orderID string) (model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", ctx, orderID)
	ret0, _ := ret[0].(model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByID indicates an expected call of GetOrderByID.
func (mr *MockGophermartRepositoryMockRecorder) GetOrderByID(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockGophermartRepository)(nil).GetOrderByID), ctx, orderID)
}

// GetOrders mocks base method.
func (m *MockGophermartRepository) GetOrders(ctx context.Context, userID string) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, userID)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockGophermartRepositoryMockRecorder) GetOrders(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockGophermartRepository)(nil).GetOrders), ctx, userID)
}

// GetUsersBalance mocks base method.
func (m *MockGophermartRepository) GetUsersBalance(ctx context.Context, userID string) (model.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersBalance", ctx, userID)
	ret0, _ := ret[0].(model.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersBalance indicates an expected call of GetUsersBalance.
func (mr *MockGophermartRepositoryMockRecorder) GetUsersBalance(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersBalance", reflect.TypeOf((*MockGophermartRepository)(nil).GetUsersBalance), ctx, userID)
}

// IncreaseUsersBalance mocks base method.
func (m *MockGophermartRepository) IncreaseUsersBalance(ctx context.Context, userID string, reward decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncreaseUsersBalance", ctx, userID, reward)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncreaseUsersBalance indicates an expected call of IncreaseUsersBalance.
func (mr *MockGophermartRepositoryMockRecorder) IncreaseUsersBalance(ctx, userID, reward interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseUsersBalance", reflect.TypeOf((*MockGophermartRepository)(nil).IncreaseUsersBalance), ctx, userID, reward)
}

// LoginUser mocks base method.
func (m *MockGophermartRepository) LoginUser(ctx context.Context, login string) (model.RegisterUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", ctx, login)
	ret0, _ := ret[0].(model.RegisterUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockGophermartRepositoryMockRecorder) LoginUser(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockGophermartRepository)(nil).LoginUser), ctx, login)
}

// PostNewOrder mocks base method.
func (m *MockGophermartRepository) PostNewOrder(ctx context.Context, orderID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostNewOrder", ctx, orderID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostNewOrder indicates an expected call of PostNewOrder.
func (mr *MockGophermartRepositoryMockRecorder) PostNewOrder(ctx, orderID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostNewOrder", reflect.TypeOf((*MockGophermartRepository)(nil).PostNewOrder), ctx, orderID, userID)
}

// RecordUsersWithdraw mocks base method.
func (m *MockGophermartRepository) RecordUsersWithdraw(ctx context.Context, userID, orderID string, amount decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordUsersWithdraw", ctx, userID, orderID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecordUsersWithdraw indicates an expected call of RecordUsersWithdraw.
func (mr *MockGophermartRepositoryMockRecorder) RecordUsersWithdraw(ctx, userID, orderID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordUsersWithdraw", reflect.TypeOf((*MockGophermartRepository)(nil).RecordUsersWithdraw), ctx, userID, orderID, amount)
}

// RegisterNewUser mocks base method.
func (m *MockGophermartRepository) RegisterNewUser(ctx context.Context, user model.RegisterUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterNewUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterNewUser indicates an expected call of RegisterNewUser.
func (mr *MockGophermartRepositoryMockRecorder) RegisterNewUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterNewUser", reflect.TypeOf((*MockGophermartRepository)(nil).RegisterNewUser), ctx, user)
}

// UpdateOrder mocks base method.
func (m *MockGophermartRepository) UpdateOrder(ctx context.Context, status model.OrderStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", ctx, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockGophermartRepositoryMockRecorder) UpdateOrder(ctx, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockGophermartRepository)(nil).UpdateOrder), ctx, status)
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// LoginUser mocks base method.
func (m *MockUserRepository) LoginUser(ctx context.Context, login string) (model.RegisterUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", ctx, login)
	ret0, _ := ret[0].(model.RegisterUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockUserRepositoryMockRecorder) LoginUser(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockUserRepository)(nil).LoginUser), ctx, login)
}

// RegisterNewUser mocks base method.
func (m *MockUserRepository) RegisterNewUser(ctx context.Context, user model.RegisterUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterNewUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterNewUser indicates an expected call of RegisterNewUser.
func (mr *MockUserRepositoryMockRecorder) RegisterNewUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterNewUser", reflect.TypeOf((*MockUserRepository)(nil).RegisterNewUser), ctx, user)
}

// MockOrderRepository is a mock of OrderRepository interface.
type MockOrderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryMockRecorder
}

// MockOrderRepositoryMockRecorder is the mock recorder for MockOrderRepository.
type MockOrderRepositoryMockRecorder struct {
	mock *MockOrderRepository
}

// NewMockOrderRepository creates a new mock instance.
func NewMockOrderRepository(ctrl *gomock.Controller) *MockOrderRepository {
	mock := &MockOrderRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepository) EXPECT() *MockOrderRepositoryMockRecorder {
	return m.recorder
}

// FetchNewOrdersToChan mocks base method.
func (m *MockOrderRepository) FetchNewOrdersToChan(ctx context.Context, ordersCh chan<- model.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchNewOrdersToChan", ctx, ordersCh)
	ret0, _ := ret[0].(error)
	return ret0
}

// FetchNewOrdersToChan indicates an expected call of FetchNewOrdersToChan.
func (mr *MockOrderRepositoryMockRecorder) FetchNewOrdersToChan(ctx, ordersCh interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchNewOrdersToChan", reflect.TypeOf((*MockOrderRepository)(nil).FetchNewOrdersToChan), ctx, ordersCh)
}

// GetOrderByID mocks base method.
func (m *MockOrderRepository) GetOrderByID(ctx context.Context, orderID string) (model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", ctx, orderID)
	ret0, _ := ret[0].(model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByID indicates an expected call of GetOrderByID.
func (mr *MockOrderRepositoryMockRecorder) GetOrderByID(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockOrderRepository)(nil).GetOrderByID), ctx, orderID)
}

// GetOrders mocks base method.
func (m *MockOrderRepository) GetOrders(ctx context.Context, userID string) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, userID)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockOrderRepositoryMockRecorder) GetOrders(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockOrderRepository)(nil).GetOrders), ctx, userID)
}

// PostNewOrder mocks base method.
func (m *MockOrderRepository) PostNewOrder(ctx context.Context, orderID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostNewOrder", ctx, orderID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostNewOrder indicates an expected call of PostNewOrder.
func (mr *MockOrderRepositoryMockRecorder) PostNewOrder(ctx, orderID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostNewOrder", reflect.TypeOf((*MockOrderRepository)(nil).PostNewOrder), ctx, orderID, userID)
}

// UpdateOrder mocks base method.
func (m *MockOrderRepository) UpdateOrder(ctx context.Context, status model.OrderStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", ctx, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockOrderRepositoryMockRecorder) UpdateOrder(ctx, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockOrderRepository)(nil).UpdateOrder), ctx, status)
}

// MockBalancesRepository is a mock of BalancesRepository interface.
type MockBalancesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBalancesRepositoryMockRecorder
}

// MockBalancesRepositoryMockRecorder is the mock recorder for MockBalancesRepository.
type MockBalancesRepositoryMockRecorder struct {
	mock *MockBalancesRepository
}

// NewMockBalancesRepository creates a new mock instance.
func NewMockBalancesRepository(ctrl *gomock.Controller) *MockBalancesRepository {
	mock := &MockBalancesRepository{ctrl: ctrl}
	mock.recorder = &MockBalancesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalancesRepository) EXPECT() *MockBalancesRepositoryMockRecorder {
	return m.recorder
}

// DecreaseUsersBalance mocks base method.
func (m *MockBalancesRepository) DecreaseUsersBalance(ctx context.Context, userID string, withdraw decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecreaseUsersBalance", ctx, userID, withdraw)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecreaseUsersBalance indicates an expected call of DecreaseUsersBalance.
func (mr *MockBalancesRepositoryMockRecorder) DecreaseUsersBalance(ctx, userID, withdraw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecreaseUsersBalance", reflect.TypeOf((*MockBalancesRepository)(nil).DecreaseUsersBalance), ctx, userID, withdraw)
}

// GetUsersBalance mocks base method.
func (m *MockBalancesRepository) GetUsersBalance(ctx context.Context, userID string) (model.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersBalance", ctx, userID)
	ret0, _ := ret[0].(model.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersBalance indicates an expected call of GetUsersBalance.
func (mr *MockBalancesRepositoryMockRecorder) GetUsersBalance(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersBalance", reflect.TypeOf((*MockBalancesRepository)(nil).GetUsersBalance), ctx, userID)
}

// IncreaseUsersBalance mocks base method.
func (m *MockBalancesRepository) IncreaseUsersBalance(ctx context.Context, userID string, reward decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncreaseUsersBalance", ctx, userID, reward)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncreaseUsersBalance indicates an expected call of IncreaseUsersBalance.
func (mr *MockBalancesRepositoryMockRecorder) IncreaseUsersBalance(ctx, userID, reward interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseUsersBalance", reflect.TypeOf((*MockBalancesRepository)(nil).IncreaseUsersBalance), ctx, userID, reward)
}

// MockHistoryRepository is a mock of HistoryRepository interface.
type MockHistoryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockHistoryRepositoryMockRecorder
}

// MockHistoryRepositoryMockRecorder is the mock recorder for MockHistoryRepository.
type MockHistoryRepositoryMockRecorder struct {
	mock *MockHistoryRepository
}

// NewMockHistoryRepository creates a new mock instance.
func NewMockHistoryRepository(ctrl *gomock.Controller) *MockHistoryRepository {
	mock := &MockHistoryRepository{ctrl: ctrl}
	mock.recorder = &MockHistoryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHistoryRepository) EXPECT() *MockHistoryRepositoryMockRecorder {
	return m.recorder
}

// RecordUsersWithdraw mocks base method.
func (m *MockHistoryRepository) RecordUsersWithdraw(ctx context.Context, userID, orderID string, amount decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordUsersWithdraw", ctx, userID, orderID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecordUsersWithdraw indicates an expected call of RecordUsersWithdraw.
func (mr *MockHistoryRepositoryMockRecorder) RecordUsersWithdraw(ctx, userID, orderID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordUsersWithdraw", reflect.TypeOf((*MockHistoryRepository)(nil).RecordUsersWithdraw), ctx, userID, orderID, amount)
}
