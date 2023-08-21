package fetcher

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"gitlab.cmpayments.local/creditcard/authorization/internal/config"
)

type configFetcher struct {
	snapshotter config.ConfigSnapshotter
}

// NewConfigFetcher returns a processor that creates a snapshot of the
// Effective Config of the CardAcceptor and sets it on the AuthorizeCommand
func NewConfigFetcher(snapshotter config.ConfigSnapshotter) *configFetcher {
	return &configFetcher{
		snapshotter: snapshotter,
	}
}

func (pt configFetcher) FetchConfig(_ context.Context, merchantID string) (entity.CardAcceptor, error) {
	conf, err := pt.snapshotter.SnapshotEffectiveConfig(merchantID)
	if err != nil {
		return entity.CardAcceptor{}, err
	}

	return conf.Merchant, nil
}
