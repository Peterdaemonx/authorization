package mock

import (
	"context"

	"github.com/google/uuid"
)

type NonceStore struct {
}

func (s NonceStore) StoreNonce(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}
