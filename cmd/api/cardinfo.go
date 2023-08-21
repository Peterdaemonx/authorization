package main

import (
	"fmt"
	"time"

	gcstorage "gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/gstorage"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
	"gitlab.cmpayments.local/libraries-go/logging"
)

func cardrange(app *application) error {

	app.cardinfo = cardinfo.NewCollection(app.conf.BlockedBins)

	storageClient, err := newStorageClient(app)
	if err != nil {
		message := fmt.Sprintf("storage failed to setup: %v", err)
		app.logger.Error(app.ctx, message)
		return fmt.Errorf(message)
	}

	bucket, err := gcstorage.NewBucket(app.logger, storageClient.Bucket(app.conf.GCP.Storage.BucketName))
	if err != nil {
		message := fmt.Sprintf("failed to create new bucket: %v", err)
		app.logger.Error(app.ctx, message)
		return fmt.Errorf(message)
	}

	binRangeFiles := map[string]string{
		"mastercard": app.conf.MasterCard.BinrangeFiletype,
		"visa":       app.conf.Visa.BinrangeFiletype,
	}
	serv := cardinfo.NewService(app.cardinfo, bucket, app.logger, binRangeFiles, app.conf.Visa.AddTestPans)

	if !app.conf.Development.MockCardInfo {
		go func() {
			for {
				ticker := time.NewTicker(time.Hour)
				select {
				case <-ticker.C:
					err := serv.LoadBinRanges(app.ctx)
					if err != nil {
						ctx := logging.ContextWithError(app.ctx, err)
						app.logger.Error(ctx, "cannot re-load BIN range tables")
					}
				case <-app.ctx.Done():
					return
				}
			}
		}()
	}

	if app.conf.Development.MockCardInfo {
		return serv.LoadTest()
	}

	return serv.LoadBinRanges(app.ctx)
}
