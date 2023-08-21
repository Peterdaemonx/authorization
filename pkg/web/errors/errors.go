package errors

import (
	"context"
	"net/http"

	"gitlab.cmpayments.local/creditcard/platform"
	platformhandler "gitlab.cmpayments.local/creditcard/platform/http/handler"
	platformLogging "gitlab.cmpayments.local/creditcard/platform/http/logging"
	"gitlab.cmpayments.local/libraries-go/http/jsonresult"
	"gitlab.cmpayments.local/libraries-go/logging"
)

func ServerErrorResponse(ctx context.Context, w http.ResponseWriter, logger platform.Logger, err error) {
	logger.Error(logging.ContextWithError(ctx, err), "internal server error")

	message := "there was an error processing your request: " + err.Error()
	errorResponse(ctx, w, logger, http.StatusInternalServerError, message, nil)
}

func FailedValidationResponse(ctx context.Context, w http.ResponseWriter, logger platform.Logger, details map[string][]string) {
	logger.Error(logging.ContextWithValue(ctx, "details", details), "unprocessable entity error")

	message := "input validation error"
	errorResponse(ctx, w, logger, http.StatusUnprocessableEntity, message, details)
}

func BadRequestResponse(ctx context.Context, w http.ResponseWriter, logger platform.Logger, err error) {
	logger.Error(logging.ContextWithError(ctx, err), "bad request error")

	message := "input validation error"
	errorResponse(ctx, w, logger, http.StatusBadRequest, message, map[string][]string{"err": {err.Error()}})
}

func NotFoundResponse(ctx context.Context, w http.ResponseWriter, logger platform.Logger) {
	logger.Error(ctx, "not found")

	message := "resource not found"
	errorResponse(ctx, w, logger, http.StatusNotFound, message, map[string][]string{})
}

func errorResponse(
	ctx context.Context,
	w http.ResponseWriter,
	logger platform.Logger,
	status int,
	message string,
	details map[string][]string,
) {
	logId, _ := platformLogging.TraceID(ctx, platformLogging.LogIDKey)
	if err := platformhandler.WriteJSON(w,
		status,
		jsonresult.MappedErrorResponse{
			Error: jsonresult.MappedError{
				Code:    status,
				LogId:   logId,
				Message: message,
				Details: details,
			}},
		nil,
	); err != nil {
		logger.Error(logging.ContextWithError(ctx, err), message)
		w.WriteHeader(status)
	}
}
