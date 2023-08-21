// //go:build integration
package spanner

//
//import (
//	"context"
//	"gitlab.cmpayments.local/creditcard/platform/currencycode"
//	"math/rand"
//	"strconv"
//	"testing"
//	"time"
//
//	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
//
//	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
//
//	"github.com/stretchr/testify/assert"
//	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
//
//	"github.com/google/uuid"
//	"github.com/googleapis/gax-go/v2/apierror"
//	"github.com/pkg/errors"
//)
//
//var (
//	pspID      = uuid.MustParse("0cd8d732-66c2-4dae-bb99-16494dea7796")
//	refundCard = entity.Card{
//		Number: "4619031141704650",
//		Cvv:    "120",
//		Holder: "C. de vlegter",
//		Expiry: entity.Expiry{
//			Year:  strconv.Itoa(time.Now().Year() + 5),
//			Month: "05",
//		},
//		MaskedPan:  "4619########4650",
//		PanTokenID: uuid.New().String(),
//		Info: cardinfo.Range{
//			Scheme: "mastercard",
//		},
//	}
//	currency, _ = currencycode.FromAlpha("EUR")
//
//	refundTransaction = entity.Refund{
//		ID:                       uuid.New(),
//		Amount:                   1200,
//		Currency:                 currency,
//		Source:                   entity.Ecommerce,
//		LocalTransactionDateTime: data.LocalTransactionDateTime(time.Now()),
//		Status:                   entity.Approved,
//		Stan:                     int(rand.Int63()),
//		Card:                     refundCard,
//		Psp:                      entity.PSP{ID: pspID},
//		CardSchemeData: entity.CardSchemeData{
//			Request: entity.CardSchemeRequest{
//				POSEntryMode: entity.POSEntryMode{
//					PanEntryMode: "81",
//					PinEntryMode: "0",
//				},
//			},
//			Response: entity.CardSchemeResponse{ResponseCode: entity.ResponseCodeFromString("00")},
//		},
//		MastercardSchemeData: entity.MastercardSchemeData{
//			Request: entity.MastercardSchemeRequest{
//				AuthorizationType: "finalAuthorization",
//				ProcessingCode: entity.ProcessingCode{
//					TransactionTypeCode: "20",
//					FromAccountTypeCode: "00",
//					ToAccountTypeCode:   "00",
//				},
//			},
//		},
//	}
//	invalidRefundTransaction = entity.Refund{
//		ID:                       uuid.New(),
//		Amount:                   1200,
//		Currency:                 currency,
//		Source:                   entity.Ecommerce,
//		LocalTransactionDateTime: data.LocalTransactionDateTime(time.Now()),
//		Status:                   entity.Declined,
//		Stan:                     int(rand.Int63()),
//		Card:                     refundCard,
//		Psp:                      entity.PSP{ID: uuid.MustParse("0cd8d732-66c2-4dae-bb99-16494dea7796")},
//	}
//)
//
//func TestRefundRepository_GetAllRefunds(t *testing.T) {
//	ctx := context.Background()
//	repository := connectToRefundRepository()
//	refund := &refundTransaction
//	refund.ID = uuid.New()
//	refund.Amount = 20111
//	refund.ProcessingDate = time.Now()
//
//	t.Run("get_all_refunds", func(t *testing.T) {
//		err := repository.CreateRefund(ctx, *refund)
//		if err != nil {
//			t.Errorf("was unable to create refund %v", err)
//		}
//
//		err = repository.CreateMastercardRefund(ctx, *refund)
//		if err != nil {
//			t.Errorf("was unable to create refund %v", err)
//		}
//
//		err = repository.UpdateRefundResponse(ctx, *refund)
//		if err != nil {
//			t.Errorf("was unable to create refund %v", err)
//		}
//
//		metaData, refunds, err := repository.GetAllRefunds(ctx, pspID, entity.Filters{
//			Page:     1,
//			PageSize: 1,
//			Sort:     "-createdAt",
//		}, map[string]interface{}{
//			"psp_id":         pspID,
//			"amount":         refund.Amount,
//			"responseCode":   refund.CardSchemeData.Response.ResponseCode.Value,
//			"processingDate": refund.ProcessingDate,
//		})
//
//		if err != nil {
//			t.Errorf("was unable to get all refunds refund")
//		}
//
//		if !assert.Equal(t, 1, len(refunds)) {
//			t.Errorf("GetCapturesByAuthorizationIDs() captures length should be 1")
//		}
//
//		if !assert.Equal(t, metaData.PageSize, 1) {
//			t.Errorf("incorrectly set pageSize. want %v, got %v", 1, metaData.PageSize)
//		}
//	})
//}
//
//func TestRefundRepository_CreateRefund(t *testing.T) {
//	var apiError apierror.APIError
//	repository := connectToRefundRepository()
//	tests := []struct {
//		name             string
//		arg              entity.Refund
//		expectedResponse interface{}
//		wantedErr        error
//	}{
//		{
//			name: "create_valid_refund",
//			arg:  refundTransaction,
//		},
//		{
//			name:             "create_invalid_refund",
//			arg:              invalidRefundTransaction,
//			expectedResponse: &apiError,
//			wantedErr:        errors.New("spanner: code = \"FailedPrecondition\", desc = \"Cannot specify a null value for column: refunds.merchant_id in table: refunds"),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			err := repository.CreateRefund(context.Background(), tt.arg)
//			if err != nil && tt.wantedErr != nil && tt.expectedResponse != &apiError {
//				t.Errorf("got: %v, want: %v", err, tt.wantedErr)
//			}
//		})
//	}
//}
//
//func connectToRefundRepository() *RefundRepository {
//	wtimeout, _ := time.ParseDuration("10s")
//	rtimeout, _ := time.ParseDuration("10s")
//	return NewRefundRepository(Client, rtimeout, wtimeout)
//}
