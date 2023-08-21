package ports

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"gitlab.cmpayments.local/creditcard/platform"
	platformhandler "gitlab.cmpayments.local/creditcard/platform/http/handler"
	"gitlab.cmpayments.local/creditcard/platform/http/logging"
	"gitlab.cmpayments.local/creditcard/platform/http/validator"

	platformErr "gitlab.cmpayments.local/creditcard/platform/http/errors"

	captureApp "gitlab.cmpayments.local/creditcard/authorization/internal/capture/app"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
)

func NewHttp(cs captureApp.CaptureService, l platform.Logger) captureHandler {
	return captureHandler{
		captureService: cs,
		logger:         l,
	}
}

type captureHandler struct {
	captureService captureApp.CaptureService
	logger         platform.Logger
}

func (ch captureHandler) CreateCapture(w http.ResponseWriter, r *http.Request) {
	input := CaptureRequest{}
	ctx := r.Context()

	params := httprouter.ParamsFromContext(ctx)
	authorizationID, err := uuid.Parse(params.ByName("authorizationID"))
	if err != nil {
		platformErr.BadRequestResponse(ctx, w, ch.logger, err)
		return
	}

	if err = platformhandler.ReadJSON(w, r, &input); err != nil {
		platformErr.BadRequestResponse(ctx, w, ch.logger, err)
		return
	}

	v := validator.New()

	input.validate(v)
	if !v.Valid() {
		platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
		return
	}

	psp, _ := processing.PSPFromContext(ctx)

	capture := mapCaptureAuthorizationRequest(ctx, authorizationID, input)

	err = ch.captureService.Capture(ctx, psp.ID, capture)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			platformErr.NotFoundResponse(ctx, w, ch.logger)
			return
		case errors.Is(err, entity.ErrAuthorizedAmountExceeded):
			v.AddError("amount", []string{err.Error()})
			platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
			return
		case errors.Is(err, entity.ErrFinalCaptureExists):
			v.AddError("final", []string{err.Error()})
			platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
			return
		case errors.Is(err, entity.ErrAuthorizationDeclined):
			v.AddError("authorization", []string{err.Error()})
			platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
			return
		case errors.Is(err, entity.ErrAuthAlreadyReversed):
			v.AddError("authorization", []string{err.Error()})
			platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
			return
		default:
			platformErr.ServerErrorResponse(ctx, w, ch.logger, err)
			return
		}
	}

	if err = platformhandler.WriteJSON(w, http.StatusCreated, mapCaptureResponse(capture), nil); err != nil {
		platformErr.ServerErrorResponse(ctx, w, ch.logger, err)
	}
}

func (ch captureHandler) CreateRefundCapture(w http.ResponseWriter, r *http.Request) {
	input := CaptureRequest{}
	ctx := r.Context()

	params := httprouter.ParamsFromContext(ctx)
	refundID, err := uuid.Parse(params.ByName("refundID"))
	if err != nil {
		platformErr.BadRequestResponse(ctx, w, ch.logger, err)
		return
	}

	if err = platformhandler.ReadJSON(w, r, &input); err != nil {
		platformErr.BadRequestResponse(ctx, w, ch.logger, err)
		return
	}

	v := validator.New()

	input.validate(v)
	if !v.Valid() {
		platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
		return
	}

	psp, _ := processing.PSPFromContext(ctx)

	capture := mapCaptureRefundRequest(ctx, refundID, input)

	err = ch.captureService.CaptureRefund(ctx, psp.ID, capture)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			platformErr.NotFoundResponse(ctx, w, ch.logger)
			return
		case errors.Is(err, entity.ErrAuthorizedAmountExceeded):
			v.AddError("amount", []string{err.Error()})
			platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
			return
		case errors.Is(err, entity.ErrFinalCaptureExists):
			v.AddError("final", []string{err.Error()})
			platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
			return
		case errors.Is(err, entity.ErrAuthorizationDeclined):
			v.AddError("authorization", []string{err.Error()})
			platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
			return
		default:
			platformErr.ServerErrorResponse(ctx, w, ch.logger, err)
			return
		}
	}

	if err = platformhandler.WriteJSON(w, http.StatusCreated, mapCaptureRefundResponse(capture), nil); err != nil {
		platformErr.ServerErrorResponse(ctx, w, ch.logger, err)
	}
}

func (ch *captureHandler) GetCapturesByAuthorizationIDs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, ok := r.URL.Query()["ids"]
	if !ok {
		platformErr.BadRequestResponse(ctx, w, ch.logger, errors.New("failed to parse URL query string parameters"))
		return
	}

	ch.logger.Info(ctx, fmt.Sprintf("fetching data for the following IDs: %s", strings.Join(ids, ", ")))

	// Here we only want to check whether we have valid UUID's. Spanner need strings so no need to actually convert to UUID
	v := validator.New()

	for _, id := range ids {
		v.Check(v.IsValidUUID(id), "ids", []string{fmt.Sprintf("failed to parse %s", id)})
	}
	if !v.Valid() {
		platformErr.FailedValidationResponse(ctx, w, ch.logger, v.Errors)
		return
	}

	captures, err := ch.captureService.GetCapturesByAuthorizationIDs(ctx, ids)
	if err != nil {
		platformErr.ServerErrorResponse(ctx, w, ch.logger, fmt.Errorf("failed to get captures: %s", err))
		return
	}

	if err = platformhandler.WriteJSON(w, http.StatusOK, mapCapturesResponse(captures), nil); err != nil {
		platformErr.ServerErrorResponse(ctx, w, ch.logger, err)
	}
}

func mapCapturesResponse(captures []entity.Capture) CapturesResponse {
	responseMap := map[string][]CaptureResponse{}

	for _, capture := range captures {

		if _, ok := responseMap[capture.AuthorizationID.String()]; ok {
			responseMap[capture.AuthorizationID.String()] = append(responseMap[capture.AuthorizationID.String()], mapCaptureResponse(capture))
		} else {
			responseMap[capture.AuthorizationID.String()] = []CaptureResponse{mapCaptureResponse(capture)}
		}
	}
	return responseMap
}

func mapCaptureAuthorizationRequest(ctx context.Context, authorizationID uuid.UUID, r CaptureRequest) entity.Capture {
	return entity.Capture{
		ID:              uuid.New(),
		LogID:           uuid.MustParse(ctx.Value(logging.LogIDKey).(string)),
		AuthorizationID: authorizationID,
		Amount:          r.Amount,
		Currency:        currencycode.Must(r.Currency),
		IsFinal:         r.IsFinal,
		Reference:       r.Reference,
	}
}

func mapCaptureResponse(c entity.Capture) CaptureResponse {
	updatedAt := &c.UpdatedAt
	if c.UpdatedAt.IsZero() {
		updatedAt = nil
	}
	return CaptureResponse{
		ID:              c.ID.String(),
		LogID:           c.LogID.String(),
		AuthorizationID: c.AuthorizationID.String(),
		Amount:          c.Amount,
		Currency:        c.Currency.Alpha3(),
		IsFinal:         c.IsFinal,
		Reference:       c.Reference,
		IRD:             c.IRD,
		UpdatedAt:       updatedAt,
	}
}

func mapCaptureRefundRequest(ctx context.Context, refundID uuid.UUID, r CaptureRequest) entity.RefundCapture {
	return entity.RefundCapture{
		ID:        uuid.New(),
		LogID:     uuid.MustParse(ctx.Value(logging.LogIDKey).(string)),
		RefundID:  refundID,
		Amount:    r.Amount,
		Currency:  currencycode.Must(r.Currency),
		IsFinal:   r.IsFinal,
		Reference: r.Reference,
	}
}

func mapCaptureRefundResponse(c entity.RefundCapture) CaptureRefundResponse {
	updatedAt := &c.UpdatedAt
	if c.UpdatedAt.IsZero() {
		updatedAt = nil
	}
	return CaptureRefundResponse{
		ID:        c.ID.String(),
		LogID:     c.LogID.String(),
		RefundID:  c.RefundID.String(),
		Amount:    c.Amount,
		Currency:  c.Currency.Alpha3(),
		IsFinal:   c.IsFinal,
		Reference: c.Reference,
		IRD:       c.IRD,
		UpdatedAt: updatedAt,
	}
}
