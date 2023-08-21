package main

import (
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/mock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/web"
)

func (app *application) WebNonce() web.NonceCheck {
	var ns web.NonceStore

	if app.conf.Development.MockNonceStore {
		ns = mock.NonceStore{}
	} else {
		ns = spanner.NewNonceRepository(app.spannerClient, app.conf.GCP.Spanner.ReadTimeout, app.conf.GCP.Spanner.WriteTimeout)
	}

	return web.NewNonceCheck(web.NewNonceService(ns), app.logger)
}
