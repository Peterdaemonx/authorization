package main

import (
	"github.com/julienschmidt/httprouter"
	authorizationApp "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app"
	authorizationPorts "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/ports"
	captureApp "gitlab.cmpayments.local/creditcard/authorization/internal/capture/app"
	capturePorts "gitlab.cmpayments.local/creditcard/authorization/internal/capture/ports"
	"gitlab.cmpayments.local/creditcard/authorization/internal/echo/ports"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/probes"
	refundApp "gitlab.cmpayments.local/creditcard/authorization/internal/refund/app"
	refundPorts "gitlab.cmpayments.local/creditcard/authorization/internal/refund/ports"
	reversalApp "gitlab.cmpayments.local/creditcard/authorization/internal/reversal/app"
	reversalPorts "gitlab.cmpayments.local/creditcard/authorization/internal/reversal/ports"
	"gitlab.cmpayments.local/creditcard/authorization/internal/timing"
	"gitlab.cmpayments.local/creditcard/platform/http/logging"
	"gitlab.cmpayments.local/libraries-go/logging/httplog"
	"net/http"
	"time"
)

func (app *application) routes() (http.Handler, error) {
	router := httprouter.New()

	p := probes.NewProbesController(app.logger)

	timeTrack := timing.NewTracker(
		timing.Slice{Interval: time.Second * 5, Duration: time.Minute},
		timing.Slice{Interval: time.Minute, Duration: time.Minute * 15},
		timing.Slice{Interval: time.Minute * 15, Duration: time.Hour},
		timing.Slice{Interval: time.Hour, Duration: time.Hour * 24},
		//timing.Slice{Interval: time.Second, Duration: time.Second * 10},
		//timing.Slice{Interval: time.Second, Duration: time.Second * 20},
		//timing.Slice{Interval: time.Second, Duration: time.Second * 30},
		//timing.Slice{Interval: time.Second, Duration: time.Second * 40},
		//timing.Slice{Interval: time.Second, Duration: time.Second * 50},
		//timing.Slice{Interval: time.Second, Duration: time.Second * 60},
	)
	timeTrack.Process(app.ctx)

	webReqAuthz := app.RequestAuthenticator()
	webNonce := app.WebNonce()

	authRepo := app.AuthorizationStore()
	captureRepo := app.CaptureStore()
	reversalRepo := app.ReversalStore()
	refundRepo := app.RefundStore()

	mastercardss := app.SequenceStore("visa_stan")
	visass := app.SequenceStore("mastercard_stan")
	schemeMapper := app.SchemeMapper(visass, mastercardss)
	tokenization, err := app.TokenizationService()
	if err != nil {
		return nil, err
	}

	publisher, err := app.MessagePublisher()
	if err != nil {
		return nil, err
	}

	reversalService := reversalApp.NewReversalService(app.logger, authRepo, captureRepo, reversalRepo, tokenization, schemeMapper)
	authorizationService := authorizationApp.NewAuthorizationService(app.logger, authRepo, tokenization, reversalService, schemeMapper)
	captureService := captureApp.NewCaptureService(authRepo, refundRepo, captureRepo, publisher, app.conf.GCP.PubSub.AuthorizationCapturedTopicID, app.conf.GCP.PubSub.RefundCapturedTopicID)
	refundService := refundApp.NewRefundService(app.logger, refundRepo, tokenization, schemeMapper)

	authorizationHandler := authorizationPorts.NewAuthorizationHandler(app.conf.AllowProductionCardNumbers, app.cardinfo, app.logger, authorizationService)

	captureHandler := capturePorts.NewHttp(
		captureService,
		app.logger,
	)

	refundHandler := refundPorts.NewRefundHandler(
		app.conf.AllowProductionCardNumbers,
		app.logger,
		refundService,
		app.cardinfo,
	)

	reversalHandler := reversalPorts.NewReversalHandler(
		app.conf.AllowProductionCardNumbers,
		app.cardinfo,
		app.logger,
		reversalService,
	)

	echoHandler := ports.NewEchoHandler(app.logger, schemeMapper)

	fs := http.FileServer(http.Dir("./docs"))

	// Endpoints
	router.HandlerFunc(http.MethodGet, "/docs/*everything", http.StripPrefix("/docs", fs).ServeHTTP)
	router.HandlerFunc(http.MethodGet, "/v1/probe/liveness", p.Ready(app.isReady))
	router.HandlerFunc(http.MethodGet, "/v1/probe/readiness", p.Health())
	router.HandlerFunc(http.MethodGet, "/v1/echo/mastercard", echoHandler.SendEchoFn("mastercard"))
	router.HandlerFunc(http.MethodGet, "/v1/echo/visa", echoHandler.SendEchoFn("visa"))

	router.HandlerFunc(http.MethodGet, `/v1/metrics`, timing.MetricsHandler(timeTrack))

	router.HandlerFunc(http.MethodPost, "/v1/authorizations",
		timeTrack.Http("create_authorization",
			webReqAuthz.WithPermission("create_authorization",
				webNonce.WithNonceCheck(authorizationHandler.CreateAuthorization))))

	router.HandlerFunc(http.MethodGet, "/v1/authorizations",
		timeTrack.Http("get_authorization",
			webReqAuthz.WithPermission("get_authorizations",
				authorizationHandler.GetAuthorizations)))

	router.HandlerFunc(http.MethodPost, "/v1/authorizations/:authorizationID/captures",
		timeTrack.Http("capture_authorization",
			webReqAuthz.WithPermission("create_capture",
				webNonce.WithNonceCheck(captureHandler.CreateCapture))))

	router.HandlerFunc(http.MethodGet, "/v1/authorizations/:authorization",
		timeTrack.Http("get_authorization",
			webReqAuthz.WithPermission("get_authorizations", authorizationHandler.GetAuthorization)))

	router.HandlerFunc(http.MethodGet, "/v1/captures",
		timeTrack.Http("get_captures",
			webReqAuthz.WithPermission("get_captures", captureHandler.GetCapturesByAuthorizationIDs)))

	router.HandlerFunc(http.MethodPost, "/v1/authorizations/:authorizationID/reversals",
		timeTrack.Http("create_reversal",
			webReqAuthz.WithPermission("create_reversal",
				webNonce.WithNonceCheck(reversalHandler.ReverseAuthorization))))

	router.HandlerFunc(http.MethodPost, "/v1/refunds",
		timeTrack.Http("create_refund",
			webReqAuthz.WithPermission("create_refund",
				webNonce.WithNonceCheck(refundHandler.CreateRefund))))

	router.HandlerFunc(http.MethodPost, "/v1/refunds/:refundID/captures",
		timeTrack.Http("capture_refund",
			webReqAuthz.WithPermission("create_capture",
				webNonce.WithNonceCheck(captureHandler.CreateRefundCapture))))

	router.HandlerFunc(http.MethodGet, "/v1/refunds",
		timeTrack.Http("get_refunds",
			webReqAuthz.WithPermission("get_refunds", refundHandler.GetRefunds)))

	// Wrap the router with all the middlewares
	return logging.NewTraceIDMiddlewareFunc()(
		httplog.NewHandler(app.logger,
			logging.EnableCORS(router, app.conf.Cors.AllowedOrigins))), nil
}
