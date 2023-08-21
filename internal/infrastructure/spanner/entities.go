package spanner

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/internal/config"
)

type ConfigSnapshotter interface {
	Configuration(ctx context.Context, merchantID string) (config.Merchant, error)
	WriteConfiguration(ctx context.Context, config config.Merchant) error
}
