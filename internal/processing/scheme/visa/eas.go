package visa

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/yerden/go-util/bcd"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/connection"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
)

func NewEas(pool pool, sp StanProvider, ssid string) Eas {
	return Eas{
		pool:            pool,
		stanProvider:    sp,
		sourceStationID: ssid,
	}
}

//go:generate mockgen -package=mock -source=./eas.go -destination=./mock/eas.go
type pool interface {
	Send(ctx context.Context, request connection.Request) connection.Received
}

type StanProvider interface {
	Next() int
}

type Eas struct {
	pool            pool
	stanProvider    StanProvider
	sourceStationID string
}

var (
	ErrNotExpectedMessage = errors.New("not expected message")
	ErrMessageTooLong     = errors.New("message too long")
)

func (m Eas) Authorize(ctx context.Context, a *entity.Authorization) error {
	a.Stan = m.nextStan()
	a.ProcessingDate = time.Now()
	AuthorizationSchemeData(a)
	message, err := MessageFromAuthorization(*a)
	if err != nil {
		return err
	}
	req := NewRequest(message)
	req.message.SourceStationID = m.sourceStationID

	fmt.Printf("send in authorization with RRN: %s\n", req.message.Fields.F037_RetrievalReferenceNumber)

	res := m.pool.Send(ctx, req).(Response)
	if res.Error() != nil {
		return fmt.Errorf("failed to send authorization request to EAS: %w", res.Error())
	}

	a.CardSchemeData.Request.RetrievalReferenceNumber = req.message.Fields.F037_RetrievalReferenceNumber
	a.CardSchemeData.Response, a.VisaSchemeData.Response = authorizationResultFromMessage(res.Message())

	// This field is set here because we aren't sending it in right now for Visa (and thus it isn't part of the message yet).
	// Unlike Mastercard, Visa doesn't downgrade the ecommerce indicator, so therefor we use the ecommerce we send to the cardscheme here.
	a.CardSchemeData.Response.EcommerceIndicator = a.ThreeDSecure.EcommerceIndicator

	return nil
}

func (m Eas) Reverse(ctx context.Context, r *entity.Reversal) error {
	r.Authorization.Stan = m.nextStan()
	r.ProcessingDate = time.Now()

	req := NewRequest(messageFromReversal(*r))
	req.message.SourceStationID = m.sourceStationID

	fmt.Printf("send in reversal with RRN: %s\n", req.message.Fields.F037_RetrievalReferenceNumber)

	res := m.pool.Send(ctx, req).(Response)
	if res.Error() != nil {
		return fmt.Errorf("failed to send reversal request to EAS: %w", res.Error())
	}

	r.CardSchemeData.Request.RetrievalReferenceNumber = req.message.Fields.F037_RetrievalReferenceNumber
	r.CardSchemeData.Response, r.Authorization.VisaSchemeData.Response = reversalResultFromMessage(res.message)

	return nil
}

func (m Eas) Refund(ctx context.Context, r *entity.Refund) error {
	r.Stan = m.nextStan()
	r.ProcessingDate = time.Now()
	RefundSchemeData(r)

	req := NewRequest(messageFromRefund(*r))
	req.message.SourceStationID = m.sourceStationID

	fmt.Printf("send in refund with RRN: %s\n", req.message.Fields.F037_RetrievalReferenceNumber)

	res := m.pool.Send(ctx, req).(Response)
	if res.Error() != nil {
		return fmt.Errorf("failed to send refund request to EAS: %w", res.Error())
	}

	r.CardSchemeData.Request.RetrievalReferenceNumber = req.message.Fields.F037_RetrievalReferenceNumber
	r.CardSchemeData.Response, r.VisaSchemeData.Response = authorizationResultFromMessage(res.Message())

	return nil
}

func (m Eas) Echo(ctx context.Context) error {
	message := Echo(m.sourceStationID)

	res := m.pool.Send(ctx, NewRequest(message)).(Response)
	if res.Error() != nil {
		return fmt.Errorf("failed to send echo request to EAS: %w", res.Error())
	}

	return nil
}

func NewHeader(msgLength int, ssid string) []byte {
	buf := make([]byte, 22)

	buf[0] = uint8(22)
	buf[1] = uint8(1)
	buf[2] = uint8(2)
	binary.BigEndian.PutUint16(buf[3:5], uint16(msgLength+22))

	if ssid != "" {
		enc := bcd.NewEncoder(bcd.Standard)
		_, err := enc.Encode(buf[8:11], []byte(ssid))
		if err != nil {
			panic(err)
		}
	}

	return buf
}

func (m Eas) nextStan() int {
	return m.stanProvider.Next()
}

func payload2Packet(msg *Message, payload []byte) ([]byte, error) {
	length := len(payload)

	if length > math.MaxUint16 {
		return nil, ErrMessageTooLong
	}
	header := NewHeader(length, msg.SourceStationID)

	lengthBuff := make([]byte, 4)
	binary.BigEndian.PutUint16(lengthBuff, uint16(length+22))
	header = append(lengthBuff, header[:]...)
	payload = append(header, payload...)
	return payload, nil
}

func msg2Payload(msg *Message) ([]byte, error) {
	var buf bytes.Buffer

	if err := msg.Encode(iso8583.NewEncoder(&buf, iso8583.FormatAscii, iso8583.NoRdwLayout, iso8583.MtiBcd4, iso8583.HexLen)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
