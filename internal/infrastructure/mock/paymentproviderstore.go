package mock

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"github.com/google/uuid"
)

type PaymentProviderStore struct {
}

func NewPaymentProviderMock() *PaymentProviderStore {
	return &PaymentProviderStore{}
}

func (p PaymentProviderStore) GetPspByID(ctx context.Context, id uuid.UUID) (entity.PSP, error) {
	return entity.PSP{}, nil
}

func (p PaymentProviderStore) GetPspByAPIKey(ctx context.Context, apiKey string) (entity.PSP, error) {
	return entity.PSP{}, nil
}
