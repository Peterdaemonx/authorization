package probes

import (
	"net/http"
	"sync/atomic"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/platform"
	"gitlab.cmpayments.local/libraries-go/http/jsonresult"
	"gitlab.cmpayments.local/libraries-go/logging"
)

type probesController struct {
	logger platform.Logger
}

func NewProbesController(logger platform.Logger) probesController {
	return probesController{logger: logger}
}

func (p probesController) Ready(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isReady.Load().(bool) {
			id := uuid.New().String()
			ctx := logging.ContextWithValues(r.Context(), map[string]interface{}{"error.id": id})
			p.logger.Error(ctx, "Service Unavailable")
			jsonresult.ServiceUnavailable(w, id)
			return
		}
		jsonresult.Ok(w, nil)
	}
}

func (p probesController) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		jsonresult.Ok(w, nil)
	}
}
