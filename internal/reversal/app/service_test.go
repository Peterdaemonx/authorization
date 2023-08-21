package app

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"

	captureMock "gitlab.cmpayments.local/creditcard/authorization/internal/capture/app/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"gitlab.cmpayments.local/libraries-go/logging"

	authMock "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app/mock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization/mocks"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
	reversalMock "gitlab.cmpayments.local/creditcard/authorization/internal/reversal/app/mock"
)

func TestReversalService_Reverse(t *testing.T) {
	tests := []struct {
		name           string
		scheme         string
		expectedError  string
		expectedStatus entity.Status
		mocks          func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
			reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
			id uuid.UUID) uuid.UUID
	}{
		{
			name:           "mastercard_reversal_not_approved",
			scheme:         mastercard,
			expectedError:  "failed to reverse authorization, authorization was not approved",
			expectedStatus: entity.Declined,
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Declined}, nil)
				return id
			},
		},
		{
			name:           "mastercard_reversal_approved",
			expectedStatus: entity.Approved,
			scheme:         mastercard,
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()

				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				capture.EXPECT().FinalCaptureExists(ctx, gomock.Any()).Return(false, nil)
				capture.EXPECT().GetCaptureSummary(ctx, gomock.Any()).Return(entity.CaptureSummary{TotalCapturedAmount: 0}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(false, nil)
				tokenizer.EXPECT().Detokenize(ctx, gomock.Any(), gomock.Any()).Return(entity.Card{Info: cardinfo.Range{Scheme: mastercard}}, nil)
				reversal.EXPECT().CreateReversal(ctx, gomock.Any()).Return(nil)
				scheme.EXPECT().Reverse(ctx, gomock.Any()).Return(nil)
				reversal.EXPECT().UpdateReversalResponse(ctx, gomock.Any()).Return(nil)

				return id
			},
		},
		{
			name:           "visa_reversal_not_approved",
			expectedStatus: entity.Declined,
			scheme:         visa,
			expectedError:  "failed to reverse authorization, authorization was not approved",
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Declined}, nil)
				return id
			},
		},
		{
			name:           "visa_reversal_approved",
			expectedStatus: entity.Approved,
			scheme:         visa,
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				capture.EXPECT().FinalCaptureExists(ctx, gomock.Any()).Return(false, nil)
				capture.EXPECT().GetCaptureSummary(ctx, gomock.Any()).Return(entity.CaptureSummary{TotalCapturedAmount: 0}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(false, nil)
				tokenizer.EXPECT().Detokenize(ctx, gomock.Any(), gomock.Any()).Return(entity.Card{Info: cardinfo.Range{Scheme: visa}}, nil)
				reversal.EXPECT().CreateReversal(ctx, gomock.Any()).Return(nil)
				reversal.EXPECT().UpdateReversalResponse(ctx, gomock.Any()).Return(nil)
				scheme.EXPECT().Reverse(ctx, gomock.Any()).Return(nil)
				return id
			},
		},
		{
			name:           "no_connection_for_scheme",
			expectedStatus: entity.Approved,
			scheme:         visa,
			expectedError:  "failed to send reversal to card scheme: no connection for scheme",
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				capture.EXPECT().FinalCaptureExists(ctx, gomock.Any()).Return(false, nil)
				capture.EXPECT().GetCaptureSummary(ctx, gomock.Any()).Return(entity.CaptureSummary{TotalCapturedAmount: 0}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(false, nil)
				tokenizer.EXPECT().Detokenize(ctx, gomock.Any(), gomock.Any()).Return(entity.Card{}, nil)
				reversal.EXPECT().CreateReversal(ctx, gomock.Any()).Return(nil)
				return id
			},
		},
		{
			name:           "detokenize_returns_error_expired_card",
			expectedStatus: entity.Approved,
			scheme:         visa,
			expectedError:  "failed to reverse",
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				expiredYear := strconv.Itoa(time.Now().Year() - 1)
				card := entity.Card{Expiry: entity.Expiry{Year: expiredYear, Month: time.Now().Month().String()}}
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				tokenizer.EXPECT().Detokenize(ctx, id.String(), entity.Card{}).Return(card, errors.New("failed to reverse"))
				return id
			},
		},
		{
			name:           "final_capture_exists_returns_error",
			expectedStatus: entity.Approved,
			scheme:         visa,
			expectedError:  "failed to fetch final capture for authorization",
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				tokenizer.EXPECT().Detokenize(ctx, gomock.Any(), gomock.Any()).Return(entity.Card{}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(false, nil)
				capture.EXPECT().FinalCaptureExists(ctx, gomock.Any()).Return(false, errors.New("failed to fetch final capture for authorization"))
				return id
			},
		},
		{
			name:           "final_capture_exists_returns_true",
			expectedStatus: entity.Approved,
			scheme:         visa,
			expectedError:  "a final capture has already been performed",
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				tokenizer.EXPECT().Detokenize(ctx, gomock.Any(), gomock.Any()).Return(entity.Card{}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(false, nil)
				capture.EXPECT().FinalCaptureExists(ctx, gomock.Any()).Return(true, nil)
				return id
			},
		},
		{
			name:           "authorization_already_reserved_returns_error",
			expectedStatus: entity.Approved,
			scheme:         visa,
			expectedError:  "failed to fetch reversals for authorization",
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				tokenizer.EXPECT().Detokenize(ctx, gomock.Any(), gomock.Any()).Return(entity.Card{}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(false, errors.New("failed to fetch reversals for authorization"))
				return id
			},
		},
		{
			name:           "authorization_already_reserved_returns_true",
			expectedStatus: entity.Approved,
			scheme:         visa,
			expectedError:  "the authorization is already reversed",
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				tokenizer.EXPECT().Detokenize(ctx, gomock.Any(), gomock.Any()).Return(entity.Card{}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(true, nil)
				return id
			},
		},
		{
			name:           "authorization_create_duplicate_val_on_index",
			expectedStatus: entity.Approved,
			scheme:         visa,
			expectedError:  "already exists",
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				capture.EXPECT().FinalCaptureExists(ctx, gomock.Any()).Return(false, nil)
				capture.EXPECT().GetCaptureSummary(ctx, gomock.Any()).Return(entity.CaptureSummary{TotalCapturedAmount: 0}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(false, nil)
				tokenizer.EXPECT().Detokenize(ctx, id.String(), entity.Card{}).Return(entity.Card{}, nil)
				reversal.EXPECT().CreateReversal(ctx, gomock.Any()).Return(entity.ErrDupValOnIndex)
				return id
			},
		},
		{
			// This case could never come to pass. If we have a wrong value in a reversal it would
			// be impossible to even call the GetAuthorizationWithSchemeData() func.
			// this test is here just in case anything changes in the future.
			name:           "create_wrong_reversal",
			scheme:         visa,
			expectedError:  "failed to reverse",
			expectedStatus: entity.Approved,
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				capture.EXPECT().FinalCaptureExists(ctx, gomock.Any()).Return(false, nil)
				capture.EXPECT().GetCaptureSummary(ctx, gomock.Any()).Return(entity.CaptureSummary{TotalCapturedAmount: 0}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(false, nil)
				tokenizer.EXPECT().Detokenize(ctx, id.String(), entity.Card{}).Return(entity.Card{}, nil)
				reversal.EXPECT().CreateReversal(ctx, entity.Reversal{Authorization: entity.Authorization{Status: entity.Approved}}).Return(errors.New("Failed to stored reversal"))
				return id
			},
		},
		{
			// This case could never come to pass. If we have a wrong value in a reversal it would
			// be impossible to even call the GetAuthorizationWithSchemeData() func.
			// this test is here just in case anything changes in the future.
			name:           "updated_reversal_return_err",
			scheme:         visa,
			expectedError:  "failed to updated reversal",
			expectedStatus: entity.Approved,
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{Status: entity.Approved}, nil)
				capture.EXPECT().FinalCaptureExists(ctx, gomock.Any()).Return(false, nil)
				capture.EXPECT().GetCaptureSummary(ctx, gomock.Any()).Return(entity.CaptureSummary{TotalCapturedAmount: 0}, nil)
				auth.EXPECT().AuthorizationAlreadyReversed(ctx, gomock.Any()).Return(false, nil)
				tokenizer.EXPECT().Detokenize(ctx, gomock.Any(), gomock.Any()).Return(entity.Card{Info: cardinfo.Range{Scheme: visa}}, nil)
				reversal.EXPECT().CreateReversal(ctx, gomock.Any()).Return(nil)
				scheme.EXPECT().Reverse(ctx, gomock.Any()).Return(nil)
				reversal.EXPECT().UpdateReversalResponse(ctx, entity.Reversal{Authorization: entity.Authorization{Status: entity.Approved, Card: entity.Card{Info: cardinfo.Range{Scheme: visa}}}}).Return(errors.New("failed to updated reversal"))
				return id
			},
		},
		{
			name:          "getAuthorizationWithSchemeData_returns_error_record_not_found",
			scheme:        visa,
			expectedError: entity.ErrRecordNotFound.Error(),
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{}, entity.ErrRecordNotFound)
				return id
			},
		},
		{
			name:          "getAuthorizationWithSchemeData_returns_error",
			scheme:        visa,
			expectedError: entity.ErrRecordNotFound.Error(),
			mocks: func(scheme *mocks.MockSchemeConnection, auth *authMock.MockRepository, capture *captureMock.MockCaptureRepository,
				reversal *reversalMock.MockReversalRepository, tokenizer *reversalMock.MockTokenizer,
				id uuid.UUID) uuid.UUID {
				ctx := context.Background()
				auth.EXPECT().GetAuthorizationWithSchemeData(ctx, gomock.Any(), gomock.Any()).Return(entity.Authorization{}, errors.New("any default error"))
				return id
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)

			schemeMock := mocks.NewMockSchemeConnection(ctrl)
			authRepo := authMock.NewMockRepository(ctrl)
			reversalRepo := reversalMock.NewMockReversalRepository(ctrl)
			captureRepo := captureMock.NewMockCaptureRepository(ctrl)
			tokenizer := reversalMock.NewMockTokenizer(ctrl)

			mapper := authorization.NewMapper(authorization.SchemeConnections{mastercard: schemeMock, visa: schemeMock}, logging.Logger{})
			service := NewReversalService(logging.Logger{}, authRepo, captureRepo, reversalRepo, tokenizer, mapper)

			expectedId := tt.mocks(schemeMock, authRepo, captureRepo, reversalRepo, tokenizer, uuid.New())
			rev := entity.Reversal{}
			err := service.Reverse(ctx, expectedId, &rev)

			if err != nil && tt.expectedError == "" {
				t.Errorf("wanted: %s got: %s", tt.expectedError, err)
			}
			if rev.Authorization.Status != tt.expectedStatus {
				t.Errorf("wanted: %s got: %s", tt.expectedStatus, rev.Authorization.Status)
			}
		})
	}

}
