package config

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

//go:generate mockgen -package=storage_mocks -source=./config.go -destination=mocks/config.go
type ConfigService interface {
	FetchConfig(ctx context.Context, merchantID string) (entity.CardAcceptor, error)
}

type ConfigSnapshotter interface {
	SnapshotEffectiveConfig(merchantID string) (EffectiveConfig, error)
}
