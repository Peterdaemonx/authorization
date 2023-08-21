package app

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	authMock "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app/mock"
	captureMock "gitlab.cmpayments.local/creditcard/authorization/internal/capture/app/mock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization/mocks"
	"gitlab.cmpayments.local/creditcard/authorization/internal/reversal/app"
	reversalMock "gitlab.cmpayments.local/creditcard/authorization/internal/reversal/app/mock"
	platformlogging "gitlab.cmpayments.local/creditcard/platform/http/logging"
	"gitlab.cmpayments.local/libraries-go/logging"
)

const (
	mastercard = "mastercard"
	visa       = "visa"
)

func TestAuthorizationService_InitialAndSubsequentAuthorization(t *testing.T) {
	tests := []struct {
		name           string
		scheme         string
		expectedError  string
		expectedStatus entity.Status
		mockedAuth1    func() entity.Authorization
		mockedAuth2    func(traceID string) entity.Authorization
		mocks          func(ctx context.Context, scheme *mocks.MockSchemeConnection, mockAuthRepository *authMock.MockRepository,
			reversal *reversalMock.MockReversalRepository, tokenizer *authMock.MockTokenizer, captureMock *captureMock.MockCaptureRepository,
			auth entity.Authorization)
	}{
		{
			name:   "valid_initial_and_subsequent_authorization",
			scheme: visa,
			mockedAuth1: func() entity.Authorization {
				auth := entity.Authorization{}
				auth.Recurring.Initial = true
				auth.Recurring.Subsequent = false
				auth.Card.Number = "1230981230981234"
				auth.Card.PanTokenID = "valid_authorization_happy_flow_test"
				auth.Card.Info.Scheme = "visa"
				auth.Psp.ID = uuid.New()
				return auth
			},
			mockedAuth2: func(traceId string) entity.Authorization {
				auth := entity.Authorization{}
				auth.Recurring.Initial = false
				auth.Recurring.Subsequent = true
				auth.Recurring.TraceID = traceId
				auth.Card.Number = "1230981230981234"
				auth.Card.PanTokenID = "valid_authorization_happy_flow_test"
				auth.Card.Info.Scheme = "visa"
				auth.Psp.ID = uuid.New()
				return auth
			},
			mocks: func(ctx context.Context, scheme *mocks.MockSchemeConnection, mockAuthRepository *authMock.MockRepository, reversal *reversalMock.MockReversalRepository, tokenizer *authMock.MockTokenizer, captureMock *captureMock.MockCaptureRepository, auth entity.Authorization) {
				tokenizer.EXPECT().Tokenize(ctx, auth.Psp.ID.String(), auth.Card).Return(auth.Card.PanTokenID, nil)
				mockAuthRepository.EXPECT().CreateAuthorization(ctx, auth).Return(nil)
				scheme.EXPECT().Authorize(ctx, &auth).Return(nil)
				mockAuthRepository.EXPECT().CreateVisaAuthorization(ctx, auth).Return(nil)
				mockAuthRepository.EXPECT().UpdateAuthorizationResponse(ctx, auth).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestAuthorizationService_Authorize(t *testing.T) {
	tests := []struct {
		name           string
		scheme         string
		expectedError  string
		expectedStatus entity.Status
		mockedAuth     func() entity.Authorization
		mocks          func(ctx context.Context, scheme *mocks.MockSchemeConnection, mockAuthRepository *authMock.MockRepository,
			reversal *reversalMock.MockReversalRepository, tokenizer *authMock.MockTokenizer, captureMock *captureMock.MockCaptureRepository,
			auth entity.Authorization)
	}{
		{
			name:   "valid_authorization_happy_flow",
			scheme: visa,
			mockedAuth: func() entity.Authorization {
				auth := entity.Authorization{}
				auth.Card.Number = "1230981230981234"
				auth.Card.PanTokenID = "valid_authorization_happy_flow_test"
				auth.Card.Info.Scheme = "visa"
				auth.Psp.ID = uuid.New()
				return auth
			},
			mocks: func(ctx context.Context, scheme *mocks.MockSchemeConnection, mockAuthRepository *authMock.MockRepository, reversal *reversalMock.MockReversalRepository, tokenizer *authMock.MockTokenizer, captureMock *captureMock.MockCaptureRepository, auth entity.Authorization) {
				tokenizer.EXPECT().Tokenize(ctx, auth.Psp.ID.String(), auth.Card).Return(auth.Card.PanTokenID, nil)
				mockAuthRepository.EXPECT().CreateAuthorization(ctx, auth).Return(nil)
				scheme.EXPECT().Authorize(ctx, &auth).Return(nil)
				mockAuthRepository.EXPECT().CreateVisaAuthorization(ctx, auth).Return(nil)
				mockAuthRepository.EXPECT().UpdateAuthorizationResponse(ctx, auth).Return(nil)
			},
		},
		{
			name:   "authorization_tokenization_failed",
			scheme: visa,
			mockedAuth: func() entity.Authorization {
				auth := entity.Authorization{}
				auth.Card.Number = "1230981230981234"
				auth.Card.PanTokenID = "valid_authorization_happy_flow_test"
				auth.Card.Info.Scheme = "visa"
				return auth
			},
			expectedError: "failed to tokenize consumer token: merchantID cannot be empty",
			mocks: func(ctx context.Context, scheme *mocks.MockSchemeConnection, mockAuthRepository *authMock.MockRepository, reversal *reversalMock.MockReversalRepository, tokenizer *authMock.MockTokenizer, captureMock *captureMock.MockCaptureRepository, auth entity.Authorization) {
				tokenizer.EXPECT().Tokenize(ctx, uuid.Nil.String(), auth.Card).Return("", errors.New("merchantID cannot be empty"))
			},
		},
		{
			name:   "authorization_amount_too_long",
			scheme: visa,
			mockedAuth: func() entity.Authorization {
				auth := entity.Authorization{}
				auth.Card.Number = "1230981230981234"
				auth.Card.PanTokenID = "authorization_amount_too_long_test"
				auth.Card.Info.Scheme = "visa"
				auth.Amount = -9223372036854775807
				return auth
			},
			expectedError: "failed to store authorization: Authorization amount is too long",
			mocks: func(ctx context.Context, scheme *mocks.MockSchemeConnection, mockAuthRepository *authMock.MockRepository, reversal *reversalMock.MockReversalRepository, tokenizer *authMock.MockTokenizer, captureMock *captureMock.MockCaptureRepository, auth entity.Authorization) {
				tokenizer.EXPECT().Tokenize(ctx, auth.Psp.ID.String(), auth.Card).Return(auth.Card.PanTokenID, nil)
				mockAuthRepository.EXPECT().CreateAuthorization(ctx, auth).Return(errors.New("Authorization amount is too long"))
			},
		},
		{
			name:   "authorization_amount_min",
			scheme: visa,
			mockedAuth: func() entity.Authorization {
				auth := entity.Authorization{}
				auth.Card.Number = "1230981230981234"
				auth.Card.PanTokenID = "authorization_amount_min_test"
				auth.Card.Info.Scheme = "visa"
				auth.Amount = -9223372036854775807
				return auth
			},
			expectedError: "failed to store authorization: Authorization amount is too short",
			mocks: func(ctx context.Context, scheme *mocks.MockSchemeConnection, mockAuthRepository *authMock.MockRepository, reversal *reversalMock.MockReversalRepository, tokenizer *authMock.MockTokenizer, captureMock *captureMock.MockCaptureRepository, auth entity.Authorization) {
				tokenizer.EXPECT().Tokenize(ctx, auth.Psp.ID.String(), auth.Card).Return(auth.Card.PanTokenID, nil)
				mockAuthRepository.EXPECT().CreateAuthorization(ctx, auth).Return(errors.New("Authorization amount is too short"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			schemeMock := mocks.NewMockSchemeConnection(ctrl)
			authRepo := authMock.NewMockRepository(ctrl)
			reversalRepo := reversalMock.NewMockReversalRepository(ctrl)
			captureRepo := captureMock.NewMockCaptureRepository(ctrl)
			tokenizer := authMock.NewMockTokenizer(ctrl)
			mapper := authorization.NewMapper(authorization.SchemeConnections{mastercard: schemeMock, visa: schemeMock}, logging.Logger{})

			reversalService := app.NewReversalService(logging.Logger{}, authRepo, captureRepo, reversalRepo, tokenizer, mapper)

			service := NewAuthorizationService(logging.Logger{}, authRepo, tokenizer, reversalService, mapper)

			auth := tt.mockedAuth()
			ctx := context.Background()
			ctx = platformlogging.NewTraceID(ctx, platformlogging.LogIDKey)
			tt.mocks(ctx, schemeMock, authRepo, reversalRepo, tokenizer, captureRepo, auth)
			err := service.Authorize(ctx, &auth)

			if err != nil && tt.expectedError == "" {
				t.Errorf("wanted: %s got: %s", tt.expectedError, err)
			}
			if auth.Status != tt.expectedStatus {
				t.Errorf("wanted: %s got: %s", tt.expectedStatus, auth.Status)
			}
		})
	}
}
