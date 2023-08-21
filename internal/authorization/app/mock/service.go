// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	entity "gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

// MockTokenizer is a mock of Tokenizer interface.
type MockTokenizer struct {
	ctrl     *gomock.Controller
	recorder *MockTokenizerMockRecorder
}

// MockTokenizerMockRecorder is the mock recorder for MockTokenizer.
type MockTokenizerMockRecorder struct {
	mock *MockTokenizer
}

// NewMockTokenizer creates a new mock instance.
func NewMockTokenizer(ctrl *gomock.Controller) *MockTokenizer {
	mock := &MockTokenizer{ctrl: ctrl}
	mock.recorder = &MockTokenizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenizer) EXPECT() *MockTokenizerMockRecorder {
	return m.recorder
}

// Detokenize mocks base method.
func (m *MockTokenizer) Detokenize(ctx context.Context, merchantID string, card entity.Card) (entity.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Detokenize", ctx, merchantID, card)
	ret0, _ := ret[0].(entity.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Detokenize indicates an expected call of Detokenize.
func (mr *MockTokenizerMockRecorder) Detokenize(ctx, merchantID, card interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Detokenize", reflect.TypeOf((*MockTokenizer)(nil).Detokenize), ctx, merchantID, card)
}

// Tokenize mocks base method.
func (m *MockTokenizer) Tokenize(ctx context.Context, merchantID string, card entity.Card) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tokenize", ctx, merchantID, card)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Tokenize indicates an expected call of Tokenize.
func (mr *MockTokenizerMockRecorder) Tokenize(ctx, merchantID, card interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tokenize", reflect.TypeOf((*MockTokenizer)(nil).Tokenize), ctx, merchantID, card)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AuthorizationAlreadyReversed mocks base method.
func (m *MockRepository) AuthorizationAlreadyReversed(ctx context.Context, id uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthorizationAlreadyReversed", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthorizationAlreadyReversed indicates an expected call of AuthorizationAlreadyReversed.
func (mr *MockRepositoryMockRecorder) AuthorizationAlreadyReversed(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthorizationAlreadyReversed", reflect.TypeOf((*MockRepository)(nil).AuthorizationAlreadyReversed), ctx, id)
}

// CreateAuthorization mocks base method.
func (m *MockRepository) CreateAuthorization(ctx context.Context, a entity.Authorization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAuthorization", ctx, a)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAuthorization indicates an expected call of CreateAuthorization.
func (mr *MockRepositoryMockRecorder) CreateAuthorization(ctx, a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAuthorization", reflect.TypeOf((*MockRepository)(nil).CreateAuthorization), ctx, a)
}

// CreateMastercardAuthorization mocks base method.
func (m *MockRepository) CreateMastercardAuthorization(ctx context.Context, a entity.Authorization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMastercardAuthorization", ctx, a)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMastercardAuthorization indicates an expected call of CreateMastercardAuthorization.
func (mr *MockRepositoryMockRecorder) CreateMastercardAuthorization(ctx, a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMastercardAuthorization", reflect.TypeOf((*MockRepository)(nil).CreateMastercardAuthorization), ctx, a)
}

// CreateVisaAuthorization mocks base method.
func (m *MockRepository) CreateVisaAuthorization(ctx context.Context, a entity.Authorization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVisaAuthorization", ctx, a)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateVisaAuthorization indicates an expected call of CreateVisaAuthorization.
func (mr *MockRepositoryMockRecorder) CreateVisaAuthorization(ctx, a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVisaAuthorization", reflect.TypeOf((*MockRepository)(nil).CreateVisaAuthorization), ctx, a)
}

// GetAllAuthorizations mocks base method.
func (m *MockRepository) GetAllAuthorizations(ctx context.Context, pspID uuid.UUID, filters entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Authorization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllAuthorizations", ctx, pspID, filters, params)
	ret0, _ := ret[0].(entity.Metadata)
	ret1, _ := ret[1].([]entity.Authorization)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllAuthorizations indicates an expected call of GetAllAuthorizations.
func (mr *MockRepositoryMockRecorder) GetAllAuthorizations(ctx, pspID, filters, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllAuthorizations", reflect.TypeOf((*MockRepository)(nil).GetAllAuthorizations), ctx, pspID, filters, params)
}

// GetAuthorization mocks base method.
func (m *MockRepository) GetAuthorization(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthorization", ctx, pspID, authorizationID)
	ret0, _ := ret[0].(entity.Authorization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthorization indicates an expected call of GetAuthorization.
func (mr *MockRepositoryMockRecorder) GetAuthorization(ctx, pspID, authorizationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthorization", reflect.TypeOf((*MockRepository)(nil).GetAuthorization), ctx, pspID, authorizationID)
}

// GetAuthorizationWithSchemeData mocks base method.
func (m *MockRepository) GetAuthorizationWithSchemeData(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthorizationWithSchemeData", ctx, pspID, authorizationID)
	ret0, _ := ret[0].(entity.Authorization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthorizationWithSchemeData indicates an expected call of GetAuthorizationWithSchemeData.
func (mr *MockRepositoryMockRecorder) GetAuthorizationWithSchemeData(ctx, pspID, authorizationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthorizationWithSchemeData", reflect.TypeOf((*MockRepository)(nil).GetAuthorizationWithSchemeData), ctx, pspID, authorizationID)
}

// UpdateAuthorizationResponse mocks base method.
func (m *MockRepository) UpdateAuthorizationResponse(ctx context.Context, a entity.Authorization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAuthorizationResponse", ctx, a)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAuthorizationResponse indicates an expected call of UpdateAuthorizationResponse.
func (mr *MockRepositoryMockRecorder) UpdateAuthorizationResponse(ctx, a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAuthorizationResponse", reflect.TypeOf((*MockRepository)(nil).UpdateAuthorizationResponse), ctx, a)
}

// UpdateAuthorizationStatus mocks base method.
func (m *MockRepository) UpdateAuthorizationStatus(ctx context.Context, authorizationID uuid.UUID, status entity.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAuthorizationStatus", ctx, authorizationID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAuthorizationStatus indicates an expected call of UpdateAuthorizationStatus.
func (mr *MockRepositoryMockRecorder) UpdateAuthorizationStatus(ctx, authorizationID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAuthorizationStatus", reflect.TypeOf((*MockRepository)(nil).UpdateAuthorizationStatus), ctx, authorizationID, status)
}