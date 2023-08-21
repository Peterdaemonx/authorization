package risk

import (
	"context"

	"gitlab.cmpayments.local/creditcard/platform/http/errors"
	"gitlab.cmpayments.local/libraries-go/http/jsonresult"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

type Rule interface {
	Support(trx interface{}) bool
	Eval(ctx context.Context, trx interface{}, assessment *Assessment) error
}

func Rules() []Rule {
	return []Rule{
		&lowValueExemption{},
	}
}

var errLowValueMaxAmountAllowed = errors.ProcessorError{
	Code:    jsonresult.BadRequestErrorCode,
	Message: "validation error",
	Details: []string{"the transaction amount for the low value exemption is not allowed for this merchant"},
}

type lowValueExemption struct {
}

func (v lowValueExemption) Support(trx interface{}) bool {
	authorization, ok := trx.(entity.Authorization)
	return ok && authorization.Exemption == entity.LowValueExemption
}

func (v lowValueExemption) Eval(_ context.Context, trx interface{}, _ *Assessment) error {
	authorization, _ := trx.(entity.Authorization)
	if authorization.Amount > 300 {
		return errLowValueMaxAmountAllowed
	}

	return nil
}
