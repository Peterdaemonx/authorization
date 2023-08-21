package main

import (
	"os"
	"os/signal"
	"syscall"
)

func shutdown(app *application) error {
	app.wg.Add(1)
	go func() {
		app.logger.Info(app.ctx, "shutdown: waiting for signal")

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		select {
		case <-c:
			app.logger.Info(app.ctx, "shutdown: got signal")
			app.shutdown()
		case <-app.ctx.Done():
			app.logger.Info(app.ctx, "shutdown: context canceled, all done")
		}

		app.wg.Done()
	}()

	return nil
}
