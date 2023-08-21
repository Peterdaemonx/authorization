package ports

import (
	"time"

	"gitlab.cmpayments.local/creditcard/platform/http/validator"

	"gitlab.cmpayments.local/creditcard/platform/currencycode"
)

type CapturesResponse map[string][]CaptureResponse

type CaptureResponse struct {
	ID              string     `json:"id"`
	LogID           string     `json:"logId"`
	AuthorizationID string     `json:"authorizationId"`
	Amount          int        `json:"amount"`
	Currency        string     `json:"currency"`
	IsFinal         bool       `json:"isFinal"`
	Reference       string     `json:"reference"`
	IRD             string     `json:"ird,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}

type CaptureRequest struct {
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	IsFinal   bool   `json:"isFinal"`
	Reference string `json:"reference"`
}

func (cr CaptureRequest) validate(v *validator.Validator) {
	v.Check(cr.Amount > 0, "amount", []string{"amount must be greater than 0"})
	_, err := currencycode.GetCurrency(cr.Currency)
	v.Check(err == nil, "currency", []string{"unsupported currency"})
	v.Check(len(cr.Reference) <= 100, "reference", []string{"reference max length 100"})
	v.Check(cr.Reference != "", "reference", []string{"reference can not be empty"})
}

type CaptureRefundResponse struct {
	ID        string     `json:"id"`
	LogID     string     `json:"logId"`
	RefundID  string     `json:"refundId"`
	Amount    int        `json:"amount"`
	Currency  string     `json:"currency"`
	IsFinal   bool       `json:"isFinal"`
	Reference string     `json:"reference"`
	IRD       string     `json:"ird,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
