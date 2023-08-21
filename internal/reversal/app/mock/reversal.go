package mock

import (
	"context"

	"github.com/google/uuid"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

type ReversalRepo struct{}

func (rr ReversalRepo) AuthorizationAlreadyReversed(ctx context.Context, id uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (rr ReversalRepo) CreateReversal(_ context.Context, _ entity.Reversal) error {
	return nil
}

func (rr ReversalRepo) UpdateReversalResponse(_ context.Context, _ entity.Reversal) error {
	return nil
}
