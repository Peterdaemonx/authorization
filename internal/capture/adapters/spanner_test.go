//go:build integration

package adapters

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
	"google.golang.org/api/iterator"

	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	infraSpanner "gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"

	"github.com/google/uuid"
	"github.com/googleapis/gax-go/v2/apierror"

	authorizationAdapter "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/adapters"
)

var (
	captureID    = uuid.New()
	authID       = uuid.New()
	validCapture = entity.Capture{
		ID:              captureID,
		AuthorizationID: authID,
		Amount:          1200,
		Currency:        currencycode.Must("EUR"),
		IsFinal:         true,
		Status:          0,
		IRD:             "1234",
	}
	psp = entity.PSP{
		ID:     uuid.New(),
		Name:   "mycompany.com",
		Prefix: "001",
	}
	authorization = entity.Authorization{
		ID:                       authID,
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
		},
	}

	timeout, _  = time.ParseDuration("10s")
	db          = "projects/cc-acquiring-development/instances/acquiring-instance/databases/authorizations"
	authRepo    *authorizationAdapter.AuthorizationRepository
	captureRepo *captureRepository
	Client      *spanner.Client
)

func TestMain(m *testing.M) {
	os.Setenv("SPANNER_EMULATOR_HOST", fmt.Sprintf("localhost:%v", 10010))
	var err error
	ctx := context.Background()
	Client, err = infraSpanner.NewSpannerClient(ctx, db, 1)
	if err != nil {
		log.Fatal(err)
	}

	defer Client.Close()
	err = resetDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}

	authRepo = authorizationAdapter.NewAuthorizationRepository(Client, timeout, timeout)
	captureRepo = NewCaptureRepository(Client, timeout, timeout)

	err = createPsp(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = authRepo.CreateAuthorization(ctx, authorization)
	if err != nil {
		log.Fatal(err)
	}
	// Run all the tests.
	retCode := m.Run()

	os.Exit(retCode)
}

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

func Test_captureRepository_CreateCapture(t *testing.T) {

	tests := []struct {
		name      string
		arg       entity.Capture
		wantedErr error
	}{
		{
			name: "created capture",
			arg: entity.Capture{
				ID:              uuid.New(),
				AuthorizationID: authID,
				Amount:          100,
				Currency:        currencycode.Must("EUR"),
				IsFinal:         true,
				Status:          1,
			},
			wantedErr: nil,
		},
		{
			name: "created invalid capture without amount",
			arg: entity.Capture{
				ID:              uuid.New(),
				AuthorizationID: authID,
				Currency:        currencycode.Must("EUR"),
				IsFinal:         true,
				Status:          1,
			},
			wantedErr: errors.New("spanner: code = \"FailedPrecondition\", desc = \"Cannot specify a null value for column: authorizations.amount in table: authorizations\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var apiError *apierror.APIError
			err := captureRepo.CreateCapture(ctx, tt.arg)
			if tt.wantedErr != nil && err != nil {
				errors.As(err, &apiError)
				if err.Error() != tt.wantedErr.Error() {
					t.Errorf("got: %v, wanted: %v", err, tt.wantedErr)
				}
			}
		})
	}
}

func Test_captureRepository_UpdateCapture(t *testing.T) {
	capture := validCapture
	capture.ID = uuid.New()

	err := captureRepo.CreateCapture(context.Background(), capture)
	if err != nil {
		t.Errorf("Cloudn't create capture: %v", err)
	}
	capture.Status = entity.CaptureCreated

	type args struct {
		ctx     context.Context
		capture entity.Capture
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update_capture_entity",
			args: args{
				ctx:     context.Background(),
				capture: capture,
			},
			wantErr: false,
		},
		{
			name: "update_not_existing_capture_entity",
			args: args{
				ctx: context.Background(),
				capture: entity.Capture{
					ID:              uuid.New(),
					AuthorizationID: authID,
					Amount:          1200,
					Currency:        currencycode.Must("EUR"),
					IsFinal:         true,
					Status:          0,
					IRD:             "1234",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := captureRepo.UpdateCapture(tt.args.ctx, tt.args.capture); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCapture() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_captureRepository_GetCaptureSummary(t *testing.T) {
	capture := validCapture
	authorization.ID = uuid.New()
	capture.ID = uuid.New()
	capture.AuthorizationID = authorization.ID
	anotherAuth := entity.Authorization{ID: uuid.New()}

	err := authRepo.CreateAuthorization(context.Background(), authorization)
	if err != nil {
		t.Errorf("Cloudn't create authorization: %v", err)
	}

	err = captureRepo.CreateCapture(context.Background(), capture)
	if err != nil {
		t.Errorf("Cloudn't create capture: %v", err)
	}

	type args struct {
		ctx           context.Context
		authorization entity.Authorization
	}
	tests := []struct {
		name    string
		args    args
		want    entity.CaptureSummary
		wantErr bool
	}{
		{
			name: "get_capture_summary",
			args: args{
				ctx:           context.Background(),
				authorization: authorization,
			},
			want: entity.CaptureSummary{
				Authorization:       authorization,
				TotalCapturedAmount: 1200,
				HasFinalCapture:     true,
			},
			wantErr: false,
		},
		{
			name: "not_found_summary",
			args: args{
				ctx:           context.Background(),
				authorization: anotherAuth,
			},
			want: entity.CaptureSummary{
				Authorization:       anotherAuth,
				TotalCapturedAmount: 0,
				HasFinalCapture:     false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary, err := captureRepo.GetCaptureSummary(tt.args.ctx, tt.args.authorization)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCaptureSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(summary, tt.want) {
				t.Errorf("GetCaptureSummary() got = %v, want %v", summary, tt.want)
			}
		})
	}
}

func Test_captureRepository_GetCapturesByAuthorizationIDs(t *testing.T) {
	capture := validCapture
	capture.ID = uuid.New()

	err := captureRepo.CreateCapture(context.Background(), capture)
	if err != nil {
		t.Errorf("Cloudn't create capture: %v", err)
	}

	type args struct {
		ctx context.Context
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "get_capture_by_authorization_id",
			args: args{
				ctx: context.Background(),
				ids: []string{authorization.ID.String()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			captures, err := captureRepo.GetCapturesByAuthorizationIDs(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCapturesByAuthorizationIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(captures) == 0 {
				t.Errorf("captures not found from authorization ID: %s", authorization.ID)
			}
		})
	}
}
