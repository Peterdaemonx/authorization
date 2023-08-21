package main

import (
	internalApp "gitlab.cmpayments.local/creditcard/authorization/internal/app"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization"
	_ "gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/mastercard"
	timingwrappers "gitlab.cmpayments.local/creditcard/authorization/internal/timing/wrappers"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/sequences"
)

const (
	mastercard = "mastercard"
	visa       = "visa"
)

func (app application) SchemeMapper(visaSequenceStore sequences.Store, mastercardSequenceStore sequences.Store) *authorization.Mapper {
	app.mip = internalApp.Mip(app.ctx, app.logger, app.mcConnectionPool, app.mip, mastercardSequenceStore, app.shutdown)
	app.eas = internalApp.Eas(app.ctx, app.logger, app.visaConnectionPool, app.eas, visaSequenceStore, app.conf.Visa.SourceStationID, app.conf.Visa.ConnectionPool.TickDelay, app.shutdown)
	sc := authorization.SchemeConnections{
		mastercard: timingwrappers.SchemeConnection{Scheme: mastercard, Connection: app.mip},
		visa:       timingwrappers.SchemeConnection{Scheme: visa, Connection: app.eas},
	}
	return authorization.NewMapper(sc, app.logger)
}
