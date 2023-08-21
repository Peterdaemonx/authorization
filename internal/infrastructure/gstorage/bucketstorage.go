package gcstorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	gstorage "cloud.google.com/go/storage"
	"gitlab.cmpayments.local/creditcard/platform"
	"google.golang.org/api/iterator"
)

type bucket struct {
	logger platform.Logger
	handle *gstorage.BucketHandle
}

func NewBucket(logger platform.Logger, handle *gstorage.BucketHandle) (bucket, error) {
	if handle == nil {
		return bucket{}, errors.New("no BucketHandle instance have been set")
	}
	return bucket{
		logger: logger,
		handle: handle,
	}, nil
}

func (m bucket) LastFile(ctx context.Context, scheme string, fileType string) (io.ReadCloser, error) {
	day := time.Now()
	lookbackLimit := day.Add(-1 * time.Hour * 24 * 365)

	for {

		q := &gstorage.Query{
			Prefix:    scheme + day.Format("/2006/01/02/") + fileType,
			Delimiter: "/",
		}

		dir := m.handle.Objects(ctx, q)

		var newestAttrs *gstorage.ObjectAttrs

		for {
			attrs, err := dir.Next()
			if err == iterator.Done {
				break
			}
			//m.logger.Info(ctx, fmt.Sprintf("BucketItem: %#v", attrs))

			if err != nil {
				return nil, fmt.Errorf("Bucket().Objects(%v): %w", q, err)
			}

			if newestAttrs == nil || attrs.Created.After(newestAttrs.Created) {
				newestAttrs = attrs
			}
		}

		if newestAttrs != nil {
			obj := m.handle.Object(newestAttrs.Name)
			r, err := obj.NewReader(ctx)
			if err != nil {
				return nil, fmt.Errorf("Bucket().Object(%s).NewReader(): %w", newestAttrs.Name, err)
			}
			return r, nil
		}

		// If nothing was found, go back one day.
		day = day.Add(-1 * time.Hour * 24)

		if day.Before(lookbackLimit) {
			return nil, nil
		}
	}
}
