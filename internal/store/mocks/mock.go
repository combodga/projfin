// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	reflect "reflect"

	projfin "github.com/combodga/projfin"
	gomock "github.com/golang/mock/gomock"
)

// MockOrder is a mock of Order interface.
type MockOrder struct {
	ctrl     *gomock.Controller
	recorder *MockOrderMockRecorder
}

// MockOrderMockRecorder is the mock recorder for MockOrder.
type MockOrderMockRecorder struct {
	mock *MockOrder
}

// NewMockOrder creates a new mock instance.
func NewMockOrder(ctrl *gomock.Controller) *MockOrder {
	mock := &MockOrder{ctrl: ctrl}
	mock.recorder = &MockOrderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrder) EXPECT() *MockOrderMockRecorder {
	return m.recorder
}

// CheckOrder mocks base method.
func (m *MockOrder) CheckOrder(username, orderNumber string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckOrder", username, orderNumber)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckOrder indicates an expected call of CheckOrder.
func (mr *MockOrderMockRecorder) CheckOrder(username, orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckOrder", reflect.TypeOf((*MockOrder)(nil).CheckOrder), username, orderNumber)
}

// GetOrdersUser mocks base method.
func (m *MockOrder) GetOrdersUser(orderNumber string) (projfin.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersUser", orderNumber)
	ret0, _ := ret[0].(projfin.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersUser indicates an expected call of GetOrdersUser.
func (mr *MockOrderMockRecorder) GetOrdersUser(orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersUser", reflect.TypeOf((*MockOrder)(nil).GetOrdersUser), orderNumber)
}

// GetUserBalance mocks base method.
func (m *MockOrder) GetUserBalance(username string) (projfin.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", username)
	ret0, _ := ret[0].(projfin.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockOrderMockRecorder) GetUserBalance(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockOrder)(nil).GetUserBalance), username)
}

// InvalidateOrder mocks base method.
func (m *MockOrder) InvalidateOrder(orderNumber string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateOrder", orderNumber)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateOrder indicates an expected call of InvalidateOrder.
func (mr *MockOrderMockRecorder) InvalidateOrder(orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateOrder", reflect.TypeOf((*MockOrder)(nil).InvalidateOrder), orderNumber)
}

// ListOrders mocks base method.
func (m *MockOrder) ListOrders(username string) ([]projfin.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOrders", username)
	ret0, _ := ret[0].([]projfin.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOrders indicates an expected call of ListOrders.
func (mr *MockOrderMockRecorder) ListOrders(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOrders", reflect.TypeOf((*MockOrder)(nil).ListOrders), username)
}

// MakeOrder mocks base method.
func (m *MockOrder) MakeOrder(username, orderNumber string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeOrder", username, orderNumber)
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeOrder indicates an expected call of MakeOrder.
func (mr *MockOrderMockRecorder) MakeOrder(username, orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeOrder", reflect.TypeOf((*MockOrder)(nil).MakeOrder), username, orderNumber)
}

// OrdersProcessing mocks base method.
func (m *MockOrder) OrdersProcessing() ([]projfin.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrdersProcessing")
	ret0, _ := ret[0].([]projfin.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OrdersProcessing indicates an expected call of OrdersProcessing.
func (mr *MockOrderMockRecorder) OrdersProcessing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrdersProcessing", reflect.TypeOf((*MockOrder)(nil).OrdersProcessing))
}

// ProcessOrder mocks base method.
func (m *MockOrder) ProcessOrder(orderNumber string, accrual float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessOrder", orderNumber, accrual)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessOrder indicates an expected call of ProcessOrder.
func (mr *MockOrderMockRecorder) ProcessOrder(orderNumber, accrual interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessOrder", reflect.TypeOf((*MockOrder)(nil).ProcessOrder), orderNumber, accrual)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// DoLogin mocks base method.
func (m *MockUser) DoLogin(username, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoLogin", username, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// DoLogin indicates an expected call of DoLogin.
func (mr *MockUserMockRecorder) DoLogin(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoLogin", reflect.TypeOf((*MockUser)(nil).DoLogin), username, password)
}

// DoRegister mocks base method.
func (m *MockUser) DoRegister(username, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoRegister", username, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// DoRegister indicates an expected call of DoRegister.
func (mr *MockUserMockRecorder) DoRegister(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoRegister", reflect.TypeOf((*MockUser)(nil).DoRegister), username, password)
}

// MockWithdraw is a mock of Withdraw interface.
type MockWithdraw struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawMockRecorder
}

// MockWithdrawMockRecorder is the mock recorder for MockWithdraw.
type MockWithdrawMockRecorder struct {
	mock *MockWithdraw
}

// NewMockWithdraw creates a new mock instance.
func NewMockWithdraw(ctrl *gomock.Controller) *MockWithdraw {
	mock := &MockWithdraw{ctrl: ctrl}
	mock.recorder = &MockWithdrawMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdraw) EXPECT() *MockWithdrawMockRecorder {
	return m.recorder
}

// ListWithdrawals mocks base method.
func (m *MockWithdraw) ListWithdrawals(username string) ([]projfin.Withdraw, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWithdrawals", username)
	ret0, _ := ret[0].([]projfin.Withdraw)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWithdrawals indicates an expected call of ListWithdrawals.
func (mr *MockWithdrawMockRecorder) ListWithdrawals(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWithdrawals", reflect.TypeOf((*MockWithdraw)(nil).ListWithdrawals), username)
}

// Withdraw mocks base method.
func (m *MockWithdraw) Withdraw(username, orderNumber string, sum float64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", username, orderNumber, sum)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockWithdrawMockRecorder) Withdraw(username, orderNumber, sum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockWithdraw)(nil).Withdraw), username, orderNumber, sum)
}
