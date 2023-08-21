package mock

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
)

type CardInfoApi struct {
}

func (t CardInfoApi) FetchCardInfo(ctx context.Context, pan string) (processing.CardInfo, error) {
	return processing.CardInfo{}, nil
}
