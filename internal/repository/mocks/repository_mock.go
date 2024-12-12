// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/Makovey/gophermart/internal/repository/model"
	gomock "github.com/golang/mock/gomock"
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
