package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"

	gspanner "cloud.google.com/go/spanner"
	"gitlab.cmpayments.local/creditcard/authorization/internal/app"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/connection"
	mastercardScheme "gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/mastercard"
	visaScheme "gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/visa"
	"gitlab.cmpayments.local/creditcard/platform"
	"gitlab.cmpayments.local/libraries-go/logging"
)

var (
	configFile     = flag.String(`config`, `config.yml`, `Configuration file`)
	displayVersion = flag.Bool("version", false, "Display version and exit")
	version        string
)

const (
	loggerName = "authorization"
)

func main() {
	flag.Parse()

	if version == "" {
		version = "develop"
	}

	if *displayVersion {
		fmt.Printf("version:\t%s\n", version)
		os.Exit(0)
	}

	var conf app.Config
	err := app.LoadConfig(*configFile, &conf)
	if err != nil {
		fmt.Printf("failed to retrieve configuration.")
		os.Exit(2)
	}

	logger := app.NewLogger(conf, loggerName)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	logger.Info(ctx, "starting authorization/"+version+" "+*configFile)
	wg := &sync.WaitGroup{}

	app := application{
		ctx:      ctx,
		conf:     conf,
		wg:       wg,
		shutdown: cancel,
		logger:   logger,
		isReady:  &atomic.Value{},
	}

	// TODO Setup similar to Spanner client
	if conf.MasterCard.ConnectionPool.MaxConnections > 0 {
		app.mcConnectionPool = connection.NewPool(conf.MasterCard.ConnectionPool, mastercardScheme.NewResponse, logger, 0)
		app.mcConnectionPool.Start()
		defer app.mcConnectionPool.Stop()
	}

	if conf.Visa.ConnectionPool.MaxConnections > 0 {
		app.visaConnectionPool = connection.NewPool(conf.Visa.ConnectionPool, visaScheme.NewResponse, logger, 2)
		app.visaConnectionPool.Start()
		defer app.visaConnectionPool.Stop()
	}

	app.isReady.Store(false)

	// The application contains multiple services
	// Each service func must:
	//	Run a routine that handles its own responsibility
	//	If that routine stops, call shutdown()
	//	Listen for <-ctx.Done(); when that happens, shut itself down
	services := []service{
		storage,
		cardrange,
		webserver,
		shutdown, // Special case: simply waits for CTRL-C command.
	}

	for _, service := range services {
		err := service(&app)
		if err != nil {
			logger.Emergency(logging.ContextWithError(app.ctx, err), "cannot start application")
			cancel()
			return
		}
	}

	logger.Info(ctx, "Everything should be started, holding")
	app.isReady.Store(true)
	wg.Wait()
	logger.Info(ctx, "Everything has stopped")
}

type application struct {
	ctx                context.Context
	conf               app.Config
	wg                 *sync.WaitGroup
	shutdown           func()
	spannerClient      *gspanner.Client // TODO Interface ?
	logger             platform.Logger
	mcConnectionPool   *connection.Pool
	visaConnectionPool *connection.Pool
	isReady            *atomic.Value
	mip                *mastercardScheme.Mip
	eas                *visaScheme.Eas
	cardinfo           *cardinfo.Collection
}

type service func(*application) error
