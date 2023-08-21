package mastercard

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/connection"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
)

func NewMip(pool pool, sp StanProvider) Mip {
	return Mip{
		pool:         pool,
		stanProvider: sp,
	}
}

//go:generate mockgen -package=mock -source=./mip.go -destination=./mock/mip.go
type pool interface {
	Send(ctx context.Context, request connection.Request) connection.Received
}

type StanProvider interface {
	Next() int
}

type Mip struct {
	pool         pool
	stanProvider StanProvider
}

var (
	ErrNotExpectedMessage = errors.New("not expected message")
	ErrUnknownEncoding    = errors.New("unknown encoding")
	ErrMessageTooLong     = errors.New("message too long")
)

func (m Mip) Authorize(ctx context.Context, a *entity.Authorization) error {
	a.Stan = m.nextStan()
	a.ProcessingDate = time.Now()
	authorizationSchemeData(a)

	req := NewRequest(messageFromAuthorization(*a))

	res := m.pool.Send(ctx, req).(Response)
	if res.Error() != nil {
		return fmt.Errorf("failed to send authorization request to MIP: %w", res.Error())
	}

	a.CardSchemeData.Response, a.MastercardSchemeData.Response = authorizationResultFromMessage(res.Message())

	return nil
}

func (m Mip) Reverse(ctx context.Context, r *entity.Reversal) error {
	r.Authorization.Stan = m.nextStan()
	r.ProcessingDate = time.Now()
	reversalSchemeData(r)

	req := NewRequest(messageFromReversal(*r))

	res := m.pool.Send(ctx, req).(Response)
	if res.Error() != nil {
		return fmt.Errorf("failed to send reversal request to MIP: %w", res.Error())
	}

	r.CardSchemeData.Response = reversalResultFromMessage(res.Message())

	return nil
}

func (m Mip) Refund(ctx context.Context, r *entity.Refund) error {
	r.Stan = m.nextStan()
	r.ProcessingDate = time.Now()
	refundSchemeData(r)

	req := NewRequest(messageFromRefund(*r))

	res := m.pool.Send(ctx, req).(Response)
	if res.Error() != nil {
		return fmt.Errorf("failed to send refund request to MIP: %w", res.Error())
	}

	r.CardSchemeData.Response, r.MastercardSchemeData.Response = refundResultFromMessage(res.Message())

	return nil
}

func (m Mip) Echo(ctx context.Context) error {
	echo := Echo()

	res := m.pool.Send(ctx, NewRequest(echo)).(Response)
	if res.Error() != nil {
		return fmt.Errorf("%w", res.Error())
	}

	return nil
}

func (m Mip) nextStan() int {
	return m.stanProvider.Next()
}

func payload2Packet(payload []byte) ([]byte, error) {
	length := len(payload)
	if length > math.MaxUint16 {
		return nil, ErrMessageTooLong
	}

	lengthBuff := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBuff, uint16(length))
	return append(lengthBuff, payload...), nil
}

func msg2Payload(msg *Message) ([]byte, error) {
	var buf bytes.Buffer

	if err := msg.Encode(iso8583.NewEncoder(&buf, iso8583.FormatEbcdic, iso8583.NoRdwLayout)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
