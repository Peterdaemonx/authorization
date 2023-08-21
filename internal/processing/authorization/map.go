package authorization

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/platform"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
)

var FormatError = errors.New("message format error")

type SchemeConnections map[string]SchemeConnection

func (sc SchemeConnections) For(s string) SchemeConnection {
	return sc[strings.ToLower(s)]
}

func NewMapper(sc SchemeConnections, logger platform.Logger) *Mapper {
	return &Mapper{sc: sc, logger: logger}
}

type Mapper struct {
	sc     SchemeConnections
	logger platform.Logger
}

//go:generate mockgen -package=mocks -source=./map.go -destination=./mocks/map.go
type SchemeConnection interface {
	Authorize(ctx context.Context, authorization *entity.Authorization) error
	Reverse(ctx context.Context, reversal *entity.Reversal) error
	Refund(ctx context.Context, refund *entity.Refund) error
	Echo(ctx context.Context) error
}

func (m Mapper) SendEcho(ctx context.Context, scheme string) error {
	sc := m.sc.For(scheme)
	if sc == nil {
		return errors.New("no connection for scheme " + scheme)
	}

	err := sc.Echo(ctx)
	if err != nil {
		return fmt.Errorf("failed to send echo with cardscheme: %w", err)
	}
	return nil
}

func (m Mapper) SendAuthorization(ctx context.Context, a *entity.Authorization) error {
	sc := m.sc.For(a.Card.Info.Scheme)
	if sc == nil {
		return errors.New("no connection for scheme " + a.Card.Info.Scheme)
	}

	// this function should only return an error and should enrich the passed in authorize
	err := sc.Authorize(ctx, a)
	if err != nil {
		return fmt.Errorf("failed to authorize with the cardscheme: %w", err)
	}

	if a.CardSchemeData.Response.ResponseCode.Value == processing.FormatError {
		return FormatError
	}

	return nil
}

func (m Mapper) SendReversal(ctx context.Context, r *entity.Reversal) error {
	sc := m.sc.For(r.Authorization.Card.Info.Scheme)
	if sc == nil {
		return errors.New("no connection for scheme " + r.Authorization.Card.Info.Scheme)
	}

	err := sc.Reverse(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to reverse with the cardscheme: %w", err)
	}

	if r.Authorization.CardSchemeData.Response.ResponseCode.Value == processing.FormatError {
		return FormatError
	}

	return nil
}

func (m Mapper) SendRefund(ctx context.Context, r *entity.Refund) error {
	sc := m.sc.For(r.Card.Info.Scheme)
	if sc == nil {
		return errors.New("no connection for scheme " + r.Card.Info.Scheme)
	}

	// this function should only return an error and should enrich the passed in authorize
	err := sc.Refund(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to refund with the cardscheme: %w", err)
	}

	if r.CardSchemeData.Response.ResponseCode.Value == processing.FormatError {
		return FormatError
	}

	return nil
}
