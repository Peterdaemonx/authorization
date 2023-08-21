package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
)

var (
	ErrAuthorizedAmountExceeded = errors.New("total captured amount exceeds the authorized amount")
	ErrFinalCaptureExists       = errors.New("a final capture has already been performed")
	ErrAuthAlreadyReversed      = errors.New("the authorization is already reversed")
	ErrRecordNotFound           = errors.New("record not found")
	ErrAuthorizationDeclined    = errors.New("can't capture declined authorization")
)

type CaptureStatus int

const (
	CaptureCreated CaptureStatus = iota
	CaptureFailed
	Cleared
	Rejected
)

func (cs CaptureStatus) String() string {
	switch cs {
	case CaptureCreated:
		return "created"
	case CaptureFailed:
		return "failed"
	case Cleared:
		return "cleared"
	case Rejected:
		return "rejected"
	default:
		return ""
	}
}

func CaptureStatusFromString(status string) CaptureStatus {
	switch status {
	case "created":
		return CaptureCreated
	case "failed":
		return CaptureFailed
	case "cleared":
		return Cleared
	case "rejected":
		return Rejected
	}
	return 0
}

type Capture struct {
	ID              uuid.UUID
	LogID           uuid.UUID
	AuthorizationID uuid.UUID
	Amount          int
	Currency        currencycode.Currency
	IsFinal         bool
	Reference       string
	Status          CaptureStatus
	IRD             string
	UpdatedAt       time.Time
}

type CaptureSummary struct {
	Authorization       Authorization
	TotalCapturedAmount int
	HasFinalCapture     bool
}

func (c CaptureSummary) ValidateExtraCapture(amount int) error {
	if c.Authorization.Status == Declined {
		return ErrAuthorizationDeclined
	}

	if c.HasFinalCapture {
		return ErrFinalCaptureExists
	}

	if c.Authorization.Amount < c.TotalCapturedAmount+amount {
		return ErrAuthorizedAmountExceeded
	}

	return nil
}

func (c CaptureSummary) IsFinalizedWith(amount int) bool {
	return c.Authorization.Amount == c.TotalCapturedAmount+amount
}

type CaptureRefundSummary struct {
	Refund              Refund
	TotalCapturedAmount int
	HasFinalCapture     bool
}

func (c CaptureRefundSummary) ValidateExtraCapture(amount int) error {
	if c.Refund.Status == Declined {
		return ErrAuthorizationDeclined
	}

	if c.HasFinalCapture {
		return ErrFinalCaptureExists
	}

	if c.Refund.Amount < c.TotalCapturedAmount+amount {
		return ErrAuthorizedAmountExceeded
	}

	return nil
}

func (c CaptureRefundSummary) IsFinalizedWith(amount int) bool {
	return c.Refund.Amount == c.TotalCapturedAmount+amount
}

type RefundCapture struct {
	ID        uuid.UUID
	LogID     uuid.UUID
	RefundID  uuid.UUID
	Amount    int
	Currency  currencycode.Currency
	IsFinal   bool
	Reference string
	Status    CaptureStatus
	IRD       string
	UpdatedAt time.Time
}
