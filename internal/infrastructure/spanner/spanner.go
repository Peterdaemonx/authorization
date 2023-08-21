package spanner

import (
	"context"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/option"
)

func NewSpannerClient(ctx context.Context, db string, poolSize int) (*spanner.Client, error) {
	client, err := spanner.NewClient(ctx, db, option.WithGRPCConnectionPool(poolSize))
	if err != nil {
		return nil, err
	}

	return client, nil
}
