package processing

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

type (
	keyContext string
)

const (
	pspKey   keyContext = "psp_id"
	nonceKey keyContext = "nonce"
)

func ContextWithPsp(ctx context.Context, psp entity.PSP) context.Context {
	return context.WithValue(ctx, pspKey, psp)
}

func PSPFromContext(ctx context.Context) (entity.PSP, bool) {
	psp, ok := ctx.Value(pspKey).(entity.PSP)
	return psp, ok
}

func ContextWithNonce(ctx context.Context, nonce string) context.Context {
	return context.WithValue(ctx, nonceKey, nonce)
}

func NonceFromContext(ctx context.Context) (entity.PSP, bool) {
	psp, ok := ctx.Value(nonceKey).(entity.PSP)
	return psp, ok
}
