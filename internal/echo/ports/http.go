package ports

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization"
	httpErrors "gitlab.cmpayments.local/creditcard/authorization/pkg/web/errors"
	"gitlab.cmpayments.local/creditcard/platform"
	"gitlab.cmpayments.local/libraries-go/http/jsonresult"
)

type echoHandler struct {
	logger platform.Logger
	mapper *authorization.Mapper
}

func NewEchoHandler(logger platform.Logger, mapper *authorization.Mapper) *echoHandler {
	return &echoHandler{
		logger: logger,
		mapper: mapper,
	}
}

func (e echoHandler) SendEchoFn(scheme string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := e.mapper.SendEcho(ctx, scheme)
		if err != nil {
			httpErrors.ServerErrorResponse(ctx, w, e.logger, errors.New(fmt.Sprintf("echoing %s failed, %s", scheme, err.Error())))
			return
		}

		jsonresult.Ok(w, nil)
	}
}
