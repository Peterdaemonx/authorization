package ports

import (
	"context"
	"errors"
	"net/http"
	"time"

	"gitlab.cmpayments.local/creditcard/platform/http/validator"

	liblogging "gitlab.cmpayments.local/libraries-go/logging"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/tokenization"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/platform"
	platformErr "gitlab.cmpayments.local/creditcard/platform/http/errors"
	platformhandler "gitlab.cmpayments.local/creditcard/platform/http/handler"
	"gitlab.cmpayments.local/creditcard/platform/http/logging"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
	"gitlab.cmpayments.local/creditcard/authorization/internal/refund/app"
	"gitlab.cmpayments.local/creditcard/authorization/internal/web"
)

type refundHandler struct {
	logger                     platform.Logger
	allowProductionCardNumbers bool
	refundService              app.RefundService
	cardRanges                 *cardinfo.Collection
}

func NewRefundHandler(allowProductionCardNumbers bool, logger platform.Logger, refundService app.RefundService, cardRanges *cardinfo.Collection) *refundHandler {
	return &refundHandler{
		allowProductionCardNumbers: allowProductionCardNumbers,
		logger:                     logger,
		refundService:              refundService,
		cardRanges:                 cardRanges,
	}
}

func (h *refundHandler) CreateRefund(w http.ResponseWriter, r *http.Request) {
	input := refundRequest{}
	ctx := r.Context()

	h.logger.Info(ctx, "create refund")

	if err := platformhandler.ReadJSON(w, r, &input); err != nil {
		platformErr.BadRequestResponse(ctx, w, h.logger, err)
		return
	}

	v := validator.New()

	input.validate(v)
	if !v.Valid() {
		platformErr.FailedValidationResponse(ctx, w, h.logger, v.Errors)
		return
	}

	cardInfo, ok := h.cardRanges.Find(input.Card.Number)
	if !ok {
		platformErr.FailedValidationResponse(ctx, w, h.logger, map[string][]string{"card.number": {"unknown card range"}})
		return
	}
	input.Card.scheme = cardInfo.Scheme

	psp, _ := processing.PSPFromContext(ctx)

	refund := mapRefundRequest(ctx, psp, input, cardInfo)

	err := h.refundService.Authorize(ctx, &refund)
	if err != nil {
		switch {
		case errors.Is(err, tokenization.ErrFailedTokenize):
			h.logger.Error(liblogging.ContextWithError(ctx, err), "internal server error")
			platformErr.ServerErrorResponse(ctx, w, h.logger, tokenization.ErrFailedTokenize)
			return
		default:
			platformErr.ServerErrorResponse(ctx, w, h.logger, err)
			return
		}
	}

	if err := platformhandler.WriteJSON(w, http.StatusOK, mapRefundResponse(refund), nil); err != nil {
		platformErr.ServerErrorResponse(ctx, w, h.logger, err)
	}
}

func mapRefundRequest(ctx context.Context, psp entity.PSP, input refundRequest, info cardinfo.Range) entity.Refund {
	if input.AuthorizationType == "" {
		input.AuthorizationType = "finalAuthorization"
	}
	return entity.Refund{
		ID:                       uuid.New(),
		LogID:                    uuid.MustParse(ctx.Value(logging.LogIDKey).(string)),
		Amount:                   input.Amount,
		Currency:                 currencycode.Must(input.Currency),
		CustomerReference:        input.Reference,
		Source:                   entity.Source(input.Source),
		LocalTransactionDateTime: input.LocalTransactionDateTime,
		Card: entity.Card{
			Number: input.Card.Number,
			Expiry: entity.Expiry{
				Year:  input.Card.Expiry.Year,
				Month: input.Card.Expiry.Month,
			},
			MaskedPan: entity.MaskPan(input.Card.Number),
			Info:      info,
		},
		CardAcceptor: entity.CardAcceptor{
			CategoryCode: input.CardAcceptor.CategoryCode,
			ID:           input.CardAcceptor.ID,
			Name:         input.CardAcceptor.Name,
			Address: entity.CardAcceptorAddress{
				PostalCode:  input.CardAcceptor.PostalCode,
				City:        input.CardAcceptor.City,
				CountryCode: input.CardAcceptor.Country,
			},
		},
		Psp: entity.PSP{
			ID: psp.ID,
		},
		MastercardSchemeData: entity.MastercardSchemeData{
			Request: entity.MastercardSchemeRequest{
				AuthorizationType: entity.AuthorizationType(input.AuthorizationType),
			},
		},
	}
}

func mapRefundResponse(r entity.Refund) refundResponse {
	refund := refundResponse{
		ID:                       r.ID.String(),
		LogID:                    r.LogID.String(),
		Amount:                   r.Amount,
		Currency:                 r.Currency.Alpha3(),
		Reference:                r.CustomerReference,
		Source:                   string(r.Source),
		LocalTransactionDateTime: &r.LocalTransactionDateTime,
		AuthorizationType:        string(r.MastercardSchemeData.Request.AuthorizationType),
		ProcessingDate:           r.ProcessingDate.Format(time.RFC3339),
		Card: CardResponse{
			Scheme: r.Card.Info.Scheme,
			Number: r.Card.MaskedPan,
		},
		CardAcceptor: CardAcceptor{
			ID:           r.CardAcceptor.ID,
			CategoryCode: r.CardAcceptor.CategoryCode,
			Name:         r.CardAcceptor.Name,
			City:         r.CardAcceptor.Address.City,
			Country:      r.CardAcceptor.Address.CountryCode,
			PostalCode:   r.CardAcceptor.Address.PostalCode,
		},
		CardSchemeResponse: CardSchemeResponse{
			Status:  r.CardSchemeData.Response.Status.String(),
			Code:    r.CardSchemeData.Response.ResponseCode.Value,
			Message: r.CardSchemeData.Response.ResponseCode.Description,
			TraceID: r.CardSchemeData.Response.TraceId,
		},
	}

	return refund
}

func (h *refundHandler) GetRefunds(w http.ResponseWriter, r *http.Request) {
	var input struct {
		params map[string]interface{}
		entity.Filters
	}

	ctx := r.Context()
	v := validator.New()
	qs := r.URL.Query()

	reference := web.ReadString(qs, "reference", "")
	amount := web.ReadInt(qs, "amount", -1, v)
	processingDate := web.ReadDate(qs, "processingDate", time.Time{}, v)
	pan := web.ReadString(qs, "pan", "")
	status := web.ReadString(qs, "status", "")
	responseCode := web.ReadString(qs, "responseCode", "")

	input.Filters.Page = web.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = web.ReadInt(qs, "pageSize", 15, v)
	input.Filters.Sort = web.ReadString(qs, "sort", "-processingDate")
	input.Filters.SortSafelist = []string{"amount", "-amount", "createdAt", "-createdAt", "processingDate", "-processingDate"}

	v.Check(len(pan) <= 4, "pan", []string{"cannot be longer than 4 digits"})
	v.Check(len(responseCode) <= 2, "responseCode", []string{"cannot be longer than 2 characters"})

	// No need to validate the PSPID. This is done in the permissionstore middleware
	parameters := map[string]interface{}{
		"reference":      reference,
		"amount":         amount,
		"processingDate": processingDate,
		"pan":            pan,
		"status":         status,
		"responseCode":   responseCode,
	}

	if ValidateFilters(v, input.Filters); !v.Valid() {
		platformErr.FailedValidationResponse(ctx, w, h.logger, v.Errors)
		return
	}

	psp, _ := processing.PSPFromContext(ctx)

	metadata, refunds, err := h.refundService.GetRefunds(ctx, psp.ID, input.Filters, parameters)
	if err != nil {
		platformErr.ServerErrorResponse(ctx, w, h.logger, err)
		return
	}

	response := RefundsResponse{
		Refunds: []refundResponse{},
	}

	for _, rf := range refunds {
		response.Refunds = append(response.Refunds, mapRefundResponse(rf))
	}

	response.Metadata.CurrentPage = metadata.CurrentPage
	response.Metadata.PageSize = metadata.PageSize
	response.Metadata.FirstPage = metadata.FirstPage
	response.Metadata.LastPage = metadata.LastPage

	if err = platformhandler.WriteJSON(w, http.StatusOK, response, nil); err != nil {
		platformErr.ServerErrorResponse(ctx, w, h.logger, err)
	}
}

func ValidateFilters(v *validator.Validator, f entity.Filters) {
	// Check that the page and page_size parameters contain sensible values. v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10000000, "page", []string{"must be a maximum of 10 million"})
	v.Check(f.PageSize > 0, "pageSize", []string{"must be greater than zero"})
	v.Check(f.Page > 0, "page", []string{"must be greater than zero"})
	v.Check(f.PageSize <= 100, "pageSize", []string{"must be a maximum of 100"})
	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", []string{"invalid sort value"})
}
