package mock

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"github.com/google/uuid"
)

type paymentServiceProviderStore struct {
}

func NewMockPaymentServiceProvider() paymentServiceProviderStore {
	return paymentServiceProviderStore{}
}

func (p paymentServiceProviderStore) GetPspByAPIKey(ctx context.Context, apiKey string) (entity.PSP, error) {
	return entity.PSP{
		ID:     uuid.New(),
		Name:   "mycompany.com",
		Prefix: "001",
	}, nil
}

func (p paymentServiceProviderStore) GetPspByID(ctx context.Context, id uuid.UUID) (entity.PSP, error) {
	return entity.PSP{
		ID:     uuid.New(),
		Name:   "mycompany.com",
		Prefix: "001",
	}, nil
}
