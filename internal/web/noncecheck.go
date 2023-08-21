package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
	"gitlab.cmpayments.local/creditcard/platform"
	platformlogging "gitlab.cmpayments.local/creditcard/platform/http/logging"
	"gitlab.cmpayments.local/libraries-go/http/jsonresult"
	"gitlab.cmpayments.local/libraries-go/logging"
)

type nonceService interface {
	ValidateNonce(ctx context.Context, pspID uuid.UUID, nonce string) error
}

var (
	ErrNonceAlreadyUsed = errors.New("nonce has already been used")
)

type NonceCheck struct {
	ns     nonceService
	logger platform.Logger
}

func NewNonceCheck(ns nonceService, logger platform.Logger) NonceCheck {
	return NonceCheck{ns, logger}
}

func (h NonceCheck) WithNonceCheck(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx      = r.Context()
			psp, ok  = processing.PSPFromContext(r.Context())
			nonce    = r.Header.Get("nonce")
			logID, _ = platformlogging.TraceID(ctx, platformlogging.LogIDKey)
		)

		if !ok {
			h.logger.Error(ctx, "failed to get psp from the context")
			jsonresult.InternalServerError(w, logID, "failed to get psp from the context")
			return
		}

		err := h.ns.ValidateNonce(r.Context(), psp.ID, nonce)
		if err != nil {
			switch {
			case errors.Is(err, ErrNonceAlreadyUsed):
				jsonresult.BadRequest(w, map[string][]string{"nonce": {ErrNonceAlreadyUsed.Error()}})
			case errors.Is(err, ErrNonceNotFound):
				jsonresult.BadRequest(w, map[string][]string{"nonce": {ErrNonceNotFound.Error()}})
			default:
				jsonresult.InternalServerError(w, logID, "error validating nonce")
			}

			h.logger.Error(logging.ContextWithError(ctx, err), "failed to validate nonce")
			return
		}

		handler.ServeHTTP(w, r.WithContext(processing.ContextWithNonce(ctx, nonce)))
	}
}
