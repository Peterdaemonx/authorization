package main

import (
	"net/http"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"

	"gitlab.cmpayments.local/creditcard/authorization/internal/cmplatform"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/mock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/web"
)

func (app *application) RequestAuthenticator() web.RequestAuthenticator {
	var pc web.PlatformClient
	var ps web.PermissionStore

	if app.conf.Development.MockCmPlatform {
		pc = mock.PlatformClient{}
	} else {
		pc = cmplatform.NewIdentityClient(app.conf.CmPlatform.BaseDomain, http.Client{})
	}

	if app.conf.Development.MockPermissionStore {
		ps = mock.PermissionStore{}
	} else {
		ps = spanner.NewPermissionRepository(app.spannerClient, app.conf.GCP.Spanner.ReadTimeout, app.conf.GCP.Spanner.WriteTimeout)
	}

	return web.NewRequestAuthenticator(app.PaymentServiceProviderStore(), ps, pc, app.logger)
}
