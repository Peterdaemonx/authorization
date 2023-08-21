package main

import (
	"context"
	"net"
	"net/http"

	"gitlab.cmpayments.local/libraries-go/logging"
)

func webserver(app *application) error {
	ctx := app.ctx

	routes, err := app.routes()
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:    app.conf.Listen,
		Handler: routes,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	app.wg.Add(1)
	go func() {
		app.logger.Info(ctx, "webserver starting on "+app.conf.Listen)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			ctx = logging.ContextWithError(ctx, err)
			app.logger.Error(ctx, "webserver crashed")
		} else {
			app.logger.Info(ctx, "webserver stopped")
		}
		app.wg.Done()
		app.shutdown()
	}()

	go func() {
		<-app.ctx.Done()
		_ = server.Close()
	}()

	return nil
}
