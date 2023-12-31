// Code generated by MockGen. DO NOT EDIT.
// Source: ./psp.go

// Package storage_mocks is a generated GoMock package.
package storage_mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	entity "gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

// MockPspStore is a mock of PspStore interface.
type MockPspStore struct {
	ctrl     *gomock.Controller
	recorder *MockPspStoreMockRecorder
}

// MockPspStoreMockRecorder is the mock recorder for MockPspStore.
type MockPspStoreMockRecorder struct {
	mock *MockPspStore
}

// NewMockPspStore creates a new mock instance.
func NewMockPspStore(ctrl *gomock.Controller) *MockPspStore {
	mock := &MockPspStore{ctrl: ctrl}
	mock.recorder = &MockPspStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPspStore) EXPECT() *MockPspStoreMockRecorder {
	return m.recorder
}

// GetPspByAPIKey mocks base method.
func (m *MockPspStore) GetPspByAPIKey(ctx context.Context, apiKey string) (entity.PSP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPspByAPIKey", ctx, apiKey)
	ret0, _ := ret[0].(entity.PSP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPspByAPIKey indicates an expected call of GetPspByAPIKey.
func (mr *MockPspStoreMockRecorder) GetPspByAPIKey(ctx, apiKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPspByAPIKey", reflect.TypeOf((*MockPspStore)(nil).GetPspByAPIKey), ctx, apiKey)
}

// GetPspByID mocks base method.
func (m *MockPspStore) GetPspByID(ctx context.Context, id uuid.UUID) (entity.PSP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPspByID", ctx, id)
	ret0, _ := ret[0].(entity.PSP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPspByID indicates an expected call of GetPspByID.
func (mr *MockPspStoreMockRecorder) GetPspByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPspByID", reflect.TypeOf((*MockPspStore)(nil).GetPspByID), ctx, id)
}
