package main

import (
	"fmt"
	timingwrappers "gitlab.cmpayments.local/creditcard/authorization/internal/timing/wrappers"
	"os"

	authorizationAdapter "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/adapters"
	captureAdapter "gitlab.cmpayments.local/creditcard/authorization/internal/capture/adapters"
	captureService "gitlab.cmpayments.local/creditcard/authorization/internal/capture/app"
	refundAdapter "gitlab.cmpayments.local/creditcard/authorization/internal/refund/adapters"
	reversalAdapter "gitlab.cmpayments.local/creditcard/authorization/internal/reversal/adapters"

	"gitlab.cmpayments.local/creditcard/authorization/internal/app"
	authApp "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app"
	authMock "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app/mock"
	captureMock "gitlab.cmpayments.local/creditcard/authorization/internal/capture/app/mock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/mock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"
	refundApp "gitlab.cmpayments.local/creditcard/authorization/internal/refund/app"
	refundMock "gitlab.cmpayments.local/creditcard/authorization/internal/refund/app/mock"
	reversal "gitlab.cmpayments.local/creditcard/authorization/internal/reversal/app"
	reversalMock "gitlab.cmpayments.local/creditcard/authorization/internal/reversal/app/mock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/web"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/sequences"

	"gitlab.cmpayments.local/libraries-go/logging"
)

const spannerDsnPattern string = `projects/%s/instances/%s/databases/%s`

func onlyMockedStorage(config app.Config) bool {
	if !config.Development.MockPermissionStore {
		return false
	}

	if !config.Development.MockData {
		return false
	}

	return true
}

func storage(app *application) error {
	if onlyMockedStorage(app.conf) {
		app.logger.Debug(app.ctx, "Storage is fully mocked")
		return nil
	}

	dsn := fmt.Sprintf(spannerDsnPattern, app.conf.GCP.ProjectID, app.conf.GCP.Spanner.Instance, app.conf.GCP.Spanner.Database)
	ctx := logging.ContextWithValue(app.ctx, "dsn", dsn)
	app.logger.Info(ctx, fmt.Sprintf("opening spanner connection %s", dsn))

	// The Spanner client can only read emulator settings from ENV,
	// so set those in the configuration defines an emulator address
	if app.conf.Development.SpannerEmulatorAddr != "" {
		os.Setenv("SPANNER_EMULATOR_HOST", app.conf.Development.SpannerEmulatorAddr)
	}
	dbClient, err := spanner.NewSpannerClient(ctx, dsn, app.conf.GCP.Spanner.PoolSize)
	if err != nil {
		ctx = logging.ContextWithError(ctx, err)
		app.logger.Error(ctx, "cannot open spanner connection")
		return err
	}

	app.spannerClient = dbClient

	go func() {
		<-app.ctx.Done()
		app.logger.Info(ctx, "closing spanner connection")
		app.spannerClient.Close()
		app.logger.Info(ctx, "closed spanner connection")
	}()
	return nil
}

func (app *application) AuthorizationStore() authApp.Repository {
	if app.conf.Development.MockData {
		return &authMock.AuthorizationRepo{}
	}

	repo := authorizationAdapter.NewAuthorizationRepository(app.spannerClient,
		app.conf.GCP.Spanner.ReadTimeout, app.conf.GCP.Spanner.WriteTimeout)

	timedRepo := timingwrappers.AuthorizationRepository{Base: repo}

	return timedRepo
}

func (app *application) RefundStore() refundApp.Repository {
	if app.conf.Development.MockData {
		return refundMock.RefundRepo{}
	}
	repo := refundAdapter.NewRefundRepository(app.spannerClient,
		app.conf.GCP.Spanner.ReadTimeout, app.conf.GCP.Spanner.WriteTimeout)

	timedRepo := timingwrappers.RefundRepository{Base: repo}

	return timedRepo
}

func (app *application) CaptureStore() captureService.CaptureRepository {
	if app.conf.Development.MockData {
		return captureMock.CaptureRepo{}
	}
	repo := captureAdapter.NewCaptureRepository(app.spannerClient,
		app.conf.GCP.Spanner.ReadTimeout, app.conf.GCP.Spanner.WriteTimeout)

	timedRepo := timingwrappers.CaptureRepository{Base: repo}

	return timedRepo
}

func (app *application) ReversalStore() reversal.ReversalRepository {
	if app.conf.Development.MockData {
		return reversalMock.ReversalRepo{}
	}

	repo := reversalAdapter.NewReversalRepository(app.spannerClient,
		app.conf.GCP.Spanner.ReadTimeout, app.conf.GCP.Spanner.WriteTimeout)

	timedRepo := timingwrappers.ReversalRepository{Base: repo}

	return timedRepo
}

func (app *application) PaymentServiceProviderStore() web.PspStore {
	if app.conf.Development.MockData {
		return mock.NewMockPaymentServiceProvider()
	}
	return spanner.NewPaymentServiceProviderRepository(app.spannerClient,
		app.conf.GCP.Spanner.ReadTimeout, app.conf.GCP.Spanner.WriteTimeout)
}

func (app *application) SequenceStore(name string) sequences.Store {
	if app.conf.Development.MockData {
		return &mock.SequenceStore{}
	}

	return spanner.NewSequenceRepository(name, app.spannerClient, app.conf.GCP.Spanner.ReadTimeout, app.conf.GCP.Spanner.WriteTimeout)
}
