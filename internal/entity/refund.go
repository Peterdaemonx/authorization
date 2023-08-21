package entity

import (
	"time"

	"gitlab.cmpayments.local/creditcard/platform/currencycode"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
)

type Refund struct {
	ID                       uuid.UUID
	LogID                    uuid.UUID
	Amount                   int
	Currency                 currencycode.Currency
	CustomerReference        string
	Source                   Source
	LocalTransactionDateTime data.LocalTransactionDateTime
	Status                   Status
	Stan                     int
	InstitutionID            string
	ProcessingDate           time.Time
	CreatedAt                time.Time
	Card                     Card
	CardAcceptor             CardAcceptor
	Psp                      PSP
	CardSchemeData           CardSchemeData
	MastercardSchemeData     MastercardSchemeData
	VisaSchemeData           VisaSchemeData
}
