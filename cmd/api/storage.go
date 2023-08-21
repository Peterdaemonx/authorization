package main

import (
	"os"

	gstorage "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func newStorageClient(app *application) (*gstorage.Client, error) {
	if app.conf.Development.MockStorageBucket {
		storageEmuHost := os.Getenv("STORAGE_EMULATOR_HOST")
		if storageEmuHost == "" {
			storageEmuHost = "http://storage:4443/storage/v1/"
		}
		return gstorage.NewClient(app.ctx, option.WithEndpoint(storageEmuHost))
	}

	return gstorage.NewClient(app.ctx)
}
