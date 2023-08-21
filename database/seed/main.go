package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	gspanner "cloud.google.com/go/spanner"
	"gitlab.cmpayments.local/creditcard/authorization/database/fixture"
	"gitlab.cmpayments.local/creditcard/authorization/internal/app"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"
)

func main() {
	var cfg app.Config
	action := flag.String("action", "seed", "command to be executed (seed)")
	configFile := flag.String("config", "config.yml", "Configuration file")

	flag.Parse()
	err := app.LoadConfig(*configFile, &cfg)
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	ctx := context.Background()

	client, err := newSpannerClient(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to setup spanner client: %s", err.Error())
	}
	defer client.Close()

	switch *action {
	case "seed":
		err := fixture.Seed(ctx, client)
		if err != nil {
			log.Fatalf("failed to seed data: %s", err.Error())
		}
	default:
		log.Fatalf("invalid argument: %s", *action)
	}

	os.Exit(0)
}

func newSpannerClient(ctx context.Context, conf app.Config) (*gspanner.Client, error) {
	// The Spanner client can only read emulator settings from ENV,
	// so set those in the configuration defines an emulator address
	if conf.Development.SpannerEmulatorAddr != "" {
		if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
			os.Setenv("SPANNER_EMULATOR_HOST", conf.Development.SpannerEmulatorAddr)
		}
	}

	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", conf.GCP.ProjectID, conf.GCP.Spanner.Instance, conf.GCP.Spanner.Database)
	dbClient, err := spanner.NewSpannerClient(ctx, dsn, conf.GCP.Spanner.PoolSize)
	if err != nil {
		return nil, err
	}

	return dbClient, nil
}
