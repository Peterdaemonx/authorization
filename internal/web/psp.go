package web

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"github.com/google/uuid"
)

//go:generate mockgen -package=storage_mocks -source=./psp.go -destination=../processing/mocks/paymentserviceproviderrepository.go
type PspStore interface {
	GetPspByAPIKey(ctx context.Context, apiKey string) (entity.PSP, error)
	GetPspByID(ctx context.Context, id uuid.UUID) (entity.PSP, error)
}
