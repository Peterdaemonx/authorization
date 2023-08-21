package web

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrNonceNotFound = errors.New("nonce header parameter not found in the request")
)

type NonceStore interface {
	StoreNonce(ctx context.Context, pspID uuid.UUID, nonce string) error
}

type service struct {
	nonceStore NonceStore
}

func NewNonceService(nonceStore NonceStore) service {
	return service{nonceStore: nonceStore}
}

func (s service) ValidateNonce(ctx context.Context, pspID uuid.UUID, nonce string) error {

	if nonce == "" {
		return ErrNonceNotFound
	}

	err := s.nonceStore.StoreNonce(ctx, pspID, nonce)
	if err != nil {
		return fmt.Errorf("failed to validate given nonce (%s) and psp ID (%s), error: %w", nonce, pspID.String(), err)
	}

	return nil
}
