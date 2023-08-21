// Code generated by MockGen. DO NOT EDIT.
// Source: ./map.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

// MockSchemeConnection is a mock of SchemeConnection interface.
type MockSchemeConnection struct {
	ctrl     *gomock.Controller
	recorder *MockSchemeConnectionMockRecorder
}

// MockSchemeConnectionMockRecorder is the mock recorder for MockSchemeConnection.
type MockSchemeConnectionMockRecorder struct {
	mock *MockSchemeConnection
}

// NewMockSchemeConnection creates a new mock instance.
func NewMockSchemeConnection(ctrl *gomock.Controller) *MockSchemeConnection {
	mock := &MockSchemeConnection{ctrl: ctrl}
	mock.recorder = &MockSchemeConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchemeConnection) EXPECT() *MockSchemeConnectionMockRecorder {
	return m.recorder
}

// Authorize mocks base method.
func (m *MockSchemeConnection) Authorize(ctx context.Context, authorization *entity.Authorization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", ctx, authorization)
	ret0, _ := ret[0].(error)
	return ret0
}

// Authorize indicates an expected call of Authorize.
func (mr *MockSchemeConnectionMockRecorder) Authorize(ctx, authorization interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockSchemeConnection)(nil).Authorize), ctx, authorization)
}

// Echo mocks base method.
func (m *MockSchemeConnection) Echo(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Echo", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Echo indicates an expected call of Echo.
func (mr *MockSchemeConnectionMockRecorder) Echo(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Echo", reflect.TypeOf((*MockSchemeConnection)(nil).Echo), ctx)
}

// Refund mocks base method.
func (m *MockSchemeConnection) Refund(ctx context.Context, refund *entity.Refund) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refund", ctx, refund)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refund indicates an expected call of Refund.
func (mr *MockSchemeConnectionMockRecorder) Refund(ctx, refund interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refund", reflect.TypeOf((*MockSchemeConnection)(nil).Refund), ctx, refund)
}

// Reverse mocks base method.
func (m *MockSchemeConnection) Reverse(ctx context.Context, reversal *entity.Reversal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reverse", ctx, reversal)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reverse indicates an expected call of Reverse.
func (mr *MockSchemeConnectionMockRecorder) Reverse(ctx, reversal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reverse", reflect.TypeOf((*MockSchemeConnection)(nil).Reverse), ctx, reversal)
}
