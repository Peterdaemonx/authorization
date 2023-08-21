package mock

import (
	"context"
	"math/rand"
	"time"

	"gitlab.cmpayments.local/creditcard/platform/currencycode"

	"gitlab.cmpayments.local/creditcard/authorization/internal/authorization/adapters"
	"gitlab.cmpayments.local/creditcard/authorization/internal/data"

	"github.com/google/uuid"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

type AuthorizationRepo struct{}

func (r AuthorizationRepo) AuthorizationAlreadyReversed(ctx context.Context, id uuid.UUID) (bool, error) {
	return false, nil
}

func (AuthorizationRepo) CreateAuthorization(ctx context.Context, a entity.Authorization) error {
	return nil
}

func (AuthorizationRepo) CreateMastercardAuthorization(ctx context.Context, a entity.Authorization) error {
	return nil
}

func (AuthorizationRepo) CreateVisaAuthorization(ctx context.Context, a entity.Authorization) error {
	return nil
}

func (AuthorizationRepo) UpdateAuthorizationResponse(ctx context.Context, a entity.Authorization) error {
	return nil
}

func (AuthorizationRepo) GetAllAuthorizations(ctx context.Context, pspID uuid.UUID, filters entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Authorization, error) {
	return entity.Metadata{}, []entity.Authorization{}, nil
}

func (r AuthorizationRepo) UpdateAuthorizationStatus(ctx context.Context, authorizationID uuid.UUID, status entity.Status) error {
	return nil
}

func (AuthorizationRepo) GetAuthorizationWithSchemeData(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error) {
	return entity.Authorization{
		ID:                       uuid.New(),
		LogID:                    uuid.New(),
		Amount:                   100,
		Currency:                 currencycode.Must("EUR"),
		CustomerReference:        "6129484611666145821",
		Source:                   entity.Ecommerce,
		LocalTransactionDateTime: data.LocalTransactionDateTime(time.Now()),
		Status:                   "new",
		Stan:                     rand.Int(),
		InstitutionID:            "0",
		Psp: entity.PSP{
			ID:     uuid.MustParse("0cd8d732-66c2-4dae-bb99-16494dea7796"),
			Name:   "mycompany.com",
			Prefix: "001",
		},
		Card: entity.Card{
			MaskedPan: "520474######0004",
		},
		CardAcceptor: entity.CardAcceptor{
			CategoryCode: "3000",
			ID:           "123456789012",
			Name:         "mycompany.com",
			Address: entity.CardAcceptorAddress{
				PostalCode:  "4825BD",
				City:        "Breda",
				CountryCode: "NLD",
			},
		},
		ThreeDSecure: entity.ThreeDSecure{
			Version:                         "1",
			AuthenticationVerificationValue: "jI3JBkkaQ1p8CBAAABy0CHUAAAA",
			EcommerceIndicator:              2,
			DirectoryServerID:               "3bd2137d-08f1-4feb-ba50-3c2d4401c91a",
		},
		CardSchemeData: entity.CardSchemeData{
			Request: entity.CardSchemeRequest{
				ProcessingCode: entity.ProcessingCode{
					TransactionTypeCode: "00",
					FromAccountTypeCode: "00",
					ToAccountTypeCode:   "00",
				},
				POSEntryMode: entity.POSEntryMode{
					PanEntryMode: "81",
					PinEntryMode: "0",
				},
			},
			Response: entity.CardSchemeResponse{
				Status: entity.AuthorizeApproved,
				ResponseCode: entity.ResponseCode{
					Value:       "00",
					Description: "approved",
				},
			}},
		MastercardSchemeData: entity.MastercardSchemeData{
			Request: entity.MastercardSchemeRequest{
				AuthorizationType: entity.FinalAuthorization,
				PosPinCaptureCode: "",
				AdditionalData: entity.AdditionalRequestData{
					TransactionCategoryCode: "T",
					OriginalEcommerceIndicator: entity.SLI{
						SecurityProtocol:         2,
						CardholderAuthentication: 1,
						UCAFCollectionIndicator:  2,
					},
				},
				PointOfServiceData: entity.PointOfServiceData{
					TerminalAttendance:                       1,
					TerminalLocation:                         4,
					CardHolderPresence:                       5,
					CardPresence:                             1,
					CardCaptureCapabilities:                  0,
					CardHolderActivatedTerminalLevel:         6,
					CardDataTerminalInputCapabilityIndicator: 7,
				},
			},
			Response: entity.MastercardSchemeResponse{
				AdditionalData: entity.AdditionalResponseData{
					AppliedEcommerceIndicator: &entity.SLI{
						SecurityProtocol:         9,
						CardholderAuthentication: 1,
						UCAFCollectionIndicator:  7,
					},
					ReasonForUCAFDowngrade: adapters.Ptr(0),
				},
				AdditionalResponseData: "",
				TraceID:                entity.MTraceID{},
			},
		},
	}, nil
}

func (AuthorizationRepo) GetVisaAuthorization(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error) {
	return entity.Authorization{}, nil
}

func (AuthorizationRepo) GetAuthorization(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error) {
	return entity.Authorization{}, nil
}
