package timingwrappers

import (
	"context"
	"fmt"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization"
	"gitlab.cmpayments.local/creditcard/authorization/internal/timing"
)

// SchemeConnection is a handwritten timing wrapper
// Some extra logic was desired to take in some information about the scheme and response-code,
// which would be too hard in a generator
type SchemeConnection struct {
	Scheme     string
	Connection authorization.SchemeConnection
}

func (sc SchemeConnection) Echo(ctx context.Context) error {
	timingLabel := fmt.Sprintf("scheme.%s.echo", sc.Scheme)
	timing.Start(ctx, timingLabel)
	timing.Tag(ctx, "scheme", sc.Scheme)
	res := sc.Connection.Echo(ctx)
	timing.Stop(ctx, timingLabel)
	return res
}

func (sc SchemeConnection) Authorize(ctx context.Context, authorization *entity.Authorization) error {
	timingLabel := fmt.Sprintf("scheme.%s.authorize", sc.Scheme)
	timing.Start(ctx, timingLabel)
	timing.Tag(ctx, "scheme", authorization.Card.Info.Scheme)
	res := sc.Connection.Authorize(ctx, authorization)
	timing.Tag(ctx, "response_code", authorization.CardSchemeData.Response.ResponseCode.Value)
	timing.Stop(ctx, timingLabel)
	return res
}

func (sc SchemeConnection) Reverse(ctx context.Context, reversal *entity.Reversal) error {
	timingLabel := fmt.Sprintf("scheme.%s.reverse", sc.Scheme)
	timing.Start(ctx, timingLabel)
	timing.Tag(ctx, "scheme", reversal.Authorization.Card.Info.Scheme)
	res := sc.Connection.Reverse(ctx, reversal)
	timing.Tag(ctx, "response_code", reversal.CardSchemeData.Response.ResponseCode.Value)
	timing.Stop(ctx, timingLabel)
	return res
}

func (sc SchemeConnection) Refund(ctx context.Context, refund *entity.Refund) error {
	timingLabel := fmt.Sprintf("scheme.%s.refund", sc.Scheme)
	timing.Start(ctx, timingLabel)
	timing.Tag(ctx, "scheme", refund.Card.Info.Scheme)
	res := sc.Connection.Refund(ctx, refund)
	timing.Tag(ctx, "response_code", refund.CardSchemeData.Response.ResponseCode.Value)
	timing.Stop(ctx, timingLabel)
	return res
}
