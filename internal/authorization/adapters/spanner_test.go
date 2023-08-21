//go:build integration

package adapters

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
	"google.golang.org/api/iterator"

	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	spannerInfra "gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
)

var (
	project  = "cc-acquiring-development"
	instance = "acquiring-instance"
	database = "authorizations"
	wtimeout = time.Second * 5
	rtimeout = time.Second * 5
	db       = fmt.Sprintf("projects/%v/instances/%v/databases/%v", project, instance, database)
	pspID    = uuid.New()
	psp      = entity.PSP{
		ID:     pspID,
		Name:   "mycompany.com",
		Prefix: "001",
	}
	visaSchemeData = entity.VisaSchemeData{
		Request: entity.VisaSchemeRequest{},
		Response: entity.VisaSchemeResponse{
			TransactionId: 98797,
		},
	}
	mastercardSchemeData = entity.MastercardSchemeData{
		Response: entity.MastercardSchemeResponse{
			AdditionalData: entity.AdditionalResponseData{
				AppliedEcommerceIndicator: entity.NewAppliedEcommerceIndicator(entity.SLI{
					SecurityProtocol:         2,
					CardholderAuthentication: 1,
					UCAFCollectionIndicator:  2,
				}),
				ReasonForUCAFDowngrade: nil,
			},
			AdditionalResponseData: "",
			TraceID: entity.MTraceID{
				FinancialNetworkCode:   "00",
				BanknetReferenceNumber: "00",
				NetworkReportingDate:   "0812",
			},
		},
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
	}
	authorization = entity.Authorization{
		LogID:                    uuid.New(),
		Amount:                   100,
		Currency:                 currencycode.Must("EUR"),
		CustomerReference:        "6129484611666145821",
		Source:                   entity.Ecommerce,
		LocalTransactionDateTime: data.LocalTransactionDateTime(time.Now()),
		Status:                   "new",
		Stan:                     rand.Int(),
		InstitutionID:            "0",
		Psp:                      psp,
		Card: entity.Card{
			MaskedPan: "520474######0004",
			Info:      cardinfo.Range{Scheme: "mastercard", ProgramID: "DEV", ProductID: "TST"},
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
				AuthorizationIDResponse: "TstID1",
			}},
	}
)

func createPsp(ctx context.Context) error {
	_, err := Client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.NewStatement(`
        	INSERT INTO psp(psp_id, name, prefix)
			VALUES (@psp_id, @name, @prefix)
		`)
		stmt.Params["psp_id"] = psp.ID.String()
		stmt.Params["name"] = psp.Name
		stmt.Params["prefix"] = psp.Prefix
		_, err := txn.Update(ctx, stmt)
		return err
	})
	return err
}

func TestAuthorizationRepository_GetAuthorizationWithSchemeData(t *testing.T) {
	ctx := context.Background()
	repo := NewAuthorizationRepository(Client, rtimeout, wtimeout)

	// Create PSP/ Fetch PSP
	err := createPsp(ctx)
	if err != nil {
		t.Errorf("createPsp(ctx): %s", err)
	}

	tests := []struct {
		name string
		arg  func() entity.Authorization
	}{
		{
			name: "test_mastercard_authorization",
			arg: func() entity.Authorization {
				auth := authorization

				auth.ID = uuid.New()
				auth.Card.Info.Scheme = "mastercard"

				auth.MastercardSchemeData = mastercardSchemeData
				auth.CardSchemeData.Request.POSEntryMode.PanEntryMode = "02"
				auth.CardSchemeData.Request.POSEntryMode.PinEntryMode = "00"
				return auth
			},
		},
		{
			name: "test_visa_authorization",
			arg: func() entity.Authorization {
				auth := authorization

				auth.ID = uuid.New()
				auth.Card.Info.Scheme = "visa"
				auth.CardSchemeData.Request.POSEntryMode.PanEntryMode = "02"
				auth.CardSchemeData.Request.POSEntryMode.PinEntryMode = "00"
				auth.VisaSchemeData = visaSchemeData
				return auth
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := tt.arg()

			err = repo.CreateAuthorization(ctx, auth)
			if err != nil {
				t.Errorf("repo.CreateAuthorization(ctx, auth): %s", err)
			}

			switch auth.Card.Info.Scheme {
			case string(entity.Mastercard):
				err = repo.CreateMastercardAuthorization(ctx, auth)
				if err != nil {
					t.Errorf("repo.CreateMastercardAuthorization(ctx, auth): %s", err)
				}
			case string(entity.Visa):
				// TODO implement
				err = repo.CreateVisaAuthorization(ctx, auth)
				if err != nil {
					t.Errorf("repo.CreateVisaAuthorization(ctx, auth): %s", err)
				}
			default:
				t.Errorf("failed determine auth.Card.Info.Scheme")
			}

			err = repo.UpdateAuthorizationResponse(ctx, auth)
			if err != nil {
				t.Errorf("repo.UpdateAuthorizationResponse(ctx, auth): %s", err)
			}

			_, err = repo.GetAuthorizationWithSchemeData(ctx, pspID, auth.ID)
			if err != nil {
				t.Errorf("repo.GetAuthorizationWithSchemeData(ctx, pspID, authId): %s", err)
			}
		})
	}
}

var Client *spanner.Client

func TestMain(m *testing.M) {
	os.Setenv("SPANNER_EMULATOR_HOST", fmt.Sprintf("localhost:%v", 10010))
	ctx := context.Background()
	var err error

	Client, err = spannerInfra.NewSpannerClient(ctx, db, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer Client.Close()
	err = resetDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Run all the tests.
	retCode := m.Run()

	os.Exit(retCode)
}

func resetDatabase(ctx context.Context) error {
	// Fetch all the user defined table names.
	iter := Client.Single().Query(ctx, spanner.Statement{SQL: `
		SELECT
			t.table_name
		FROM INFORMATION_SCHEMA.tables AS t
		WHERE t.table_schema = ''
		AND t.table_name != 'SchemaMigrations';
   `})
	defer iter.Stop()
	var ms []*spanner.Mutation
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return err
		}

		var tableName string
		if err := row.ColumnByName("table_name", &tableName); err != nil {
			return err
		}

		// Create mutation for deleting all data in table.
		ms = append(ms, spanner.Delete(tableName, spanner.AllKeys()))
	}

	_, err := Client.Apply(ctx, ms)
	return err
}
