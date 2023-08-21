package mock

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"

	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
)

type AuthorizationStore struct {
	auths []entity.Authorization
}

func (a AuthorizationStore) GetCapturesByAuthorizationIDs(ctx context.Context, ids []string) ([]entity.Capture, error) {
	captures := make([]entity.Capture, 1)
	captures = append(captures, createCapture())
	return captures, nil
}

func (a AuthorizationStore) CreateAuthorization(
	_ context.Context,
	_ processing.AuthorizationEntity,
) (uuid.UUID, error) {
	//TODO implement me
	return uuid.New(), nil
}

func (a AuthorizationStore) GetSchemeByAuthorizationID(_ context.Context, _ uuid.UUID) (string, error) {
	return entity.Mastercard, nil
}

func NewMockAuthorizations() *AuthorizationStore {
	return &AuthorizationStore{
		//auths: createAuth(),
	}
}

func createAuth(
	ID string,
	status entity.AuthorizationStatus,
	panTokenID string,
	source entity.Source,
	financialNetworkCode,
	banknetReferenceNumber,
	networkReportingDate,
	traceID string,
	localTime time.Time,
) entity.Authorization {
	return entity.Authorization{
		ID:                       uuid.MustParse(ID),
		LogID:                    uuid.New(),
		Amount:                   100,
		Currency:                 currencycode.Must("EUR"),
		CustomerReference:        "6129484611666145821",
		Source:                   source,
		LocalTransactionDateTime: data.LocalTransactionDateTime(localTime),
		Status:                   "new",
		Stan:                     rand.Int(),
		InstitutionID:            "0",
		Recurring:                entity.Recurring{},
		Card: entity.Card{
			MaskedPan:  "520474######0004",
			PanTokenID: panTokenID,
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
		Exemption: "",
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
				Status: status,
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
					PinServiceCode: "",
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
				TraceID: entity.MTraceID{
					FinancialNetworkCode:   financialNetworkCode,
					BanknetReferenceNumber: banknetReferenceNumber,
					NetworkReportingDate:   networkReportingDate,
				},
			},
		},
	}
}

func createCapture() entity.Capture {
	return entity.Capture{
		ID:              uuid.New(),
		AuthorizationID: uuid.New(),
		Amount:          100,
		Currency:        currencycode.Must("EUR"),
		IsFinal:         false,
		Status:          0,
	}
}

func (a AuthorizationStore) CreateVisaAuthorization(ctx context.Context, authorizationEntity processing.AuthorizationEntity) error {
	return nil
}

func (a AuthorizationStore) AuthorizationByID(_ context.Context, ID string) (entity.Authorization, error) {
	for _, auth := range a.auths {
		if auth.ID.String() == ID {
			return auth, nil
		}
	}

	return entity.Authorization{}, processing.ErrAuthorizationNotFound
}

func (a AuthorizationStore) GetAllAuthorizations(ctx context.Context, pspID uuid.UUID, filters entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Authorization, error) {
	var ps int
	var page int
	if filters.Page*filters.PageSize < len(a.auths) {
		ps = filters.Page * filters.PageSize
	}
	if (ps - filters.PageSize) < 0 {
		page = 0
	} else {
		page = ps - filters.PageSize
	}

	return entity.Metadata{
		CurrentPage: filters.Page,
		PageSize:    filters.PageSize,
		FirstPage:   1,
		LastPage:    false,
	}, a.auths[page:ps], nil
}

func (a AuthorizationStore) CreateMastercardAuthorization(_ context.Context, _ processing.AuthorizationEntity) error {
	return nil
}

func (a AuthorizationStore) GetMastercardAuthorization(ctx context.Context, pspID, authorizationID uuid.UUID, status entity.Status) (entity.Authorization, error) {
	return createAuth(
		"750b17dd-b89b-4991-b3d7-cca78ca7a654",
		entity.AuthorizationStatusFromString("authorized"),
		"cd89ecb2-f50b-466f-a1a5-7b7f4d9a58d1",
		entity.Ecommerce,
		"MWE",
		"01105T",
		"0217",
		fmt.Sprintf("%d%03d", 100, 1),
		time.Time{}), nil
}
