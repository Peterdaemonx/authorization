package ports

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.cmpayments.local/creditcard/platform"
	platformErr "gitlab.cmpayments.local/creditcard/platform/http/errors"
	platformhandler "gitlab.cmpayments.local/creditcard/platform/http/handler"
	"gitlab.cmpayments.local/creditcard/platform/http/logging"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/tokenization"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
	reversalApp "gitlab.cmpayments.local/creditcard/authorization/internal/reversal/app"
)

type reversalHandler struct {
	logger                     platform.Logger
	allowProductionCardNumbers bool
	cardRanges                 *cardinfo.Collection
	reversalService            reversalApp.ReversalService
}

func NewReversalHandler(allowProductionCardNumbers bool, cardRanges *cardinfo.Collection, logger platform.Logger, reversalService reversalApp.ReversalService) reversalHandler {
	return reversalHandler{
		allowProductionCardNumbers: allowProductionCardNumbers,
		cardRanges:                 cardRanges,
		logger:                     logger,
		reversalService:            reversalService,
	}
}

func (h reversalHandler) ReverseAuthorization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := httprouter.ParamsFromContext(ctx)
	authorizationID, err := uuid.Parse(params.ByName("authorizationID"))
	if err != nil {
		platformErr.BadRequestResponse(ctx, w, h.logger, err)
		return
	}

	psp, _ := processing.PSPFromContext(ctx)

	reversal := mapReversalRequest(ctx, authorizationID)

	err = h.reversalService.Reverse(ctx, psp.ID, &reversal)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			platformErr.NotFoundResponse(ctx, w, h.logger)
			return
		case errors.Is(err, entity.ErrAuthorizationNotApproved):
			platformErr.FailedValidationResponse(ctx, w, h.logger, map[string][]string{"authorizationId": {err.Error()}})
			return
		case errors.Is(err, entity.ErrFinalCaptureExists):
			platformErr.FailedValidationResponse(ctx, w, h.logger, map[string][]string{"authorizationId": {err.Error()}})
			return
		case errors.Is(err, entity.ErrDupValOnIndex):
			platformErr.FailedValidationResponse(ctx, w, h.logger, map[string][]string{"authorizationId": {err.Error()}})
			return
		case errors.Is(err, entity.ErrAuthAlreadyReversed):
			platformErr.FailedValidationResponse(ctx, w, h.logger, map[string][]string{"authorizationId": {err.Error()}})
			return
		case errors.Is(err, tokenization.ErrFailedDetokenize):
			platformErr.ServerErrorResponse(ctx, w, h.logger, tokenization.ErrFailedDetokenize)
			return
		default:
			platformErr.ServerErrorResponse(ctx, w, h.logger, err)
			return
		}
	}

	if err = platformhandler.WriteJSON(w, http.StatusCreated, mapReversalResponse(reversal), nil); err != nil {
		platformErr.ServerErrorResponse(ctx, w, h.logger, err)
	}
}

func mapReversalRequest(ctx context.Context, authorizationID uuid.UUID) entity.Reversal {
	return entity.Reversal{
		ID:              uuid.New(),
		LogID:           uuid.MustParse(ctx.Value(logging.LogIDKey).(string)),
		AuthorizationID: authorizationID,
		Status:          entity.ReversalNew,
	}
}

func mapReversalResponse(r entity.Reversal) ReversalResponse {
	return ReversalResponse{
		ID:              r.ID.String(),
		LogID:           r.LogID.String(),
		AuthorizationID: r.AuthorizationID.String(),
		CardSchemeResponse: CardSchemeResponse{
			Status:  r.CardSchemeData.Response.Status.String(),
			Code:    r.CardSchemeData.Response.ResponseCode.Value,
			Message: r.CardSchemeData.Response.ResponseCode.Description,
		},
	}
}
