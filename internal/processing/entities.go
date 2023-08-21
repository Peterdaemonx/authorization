package processing

import (
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/data"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

const (
	AuthorizedStatus entity.Status = "authorized"
	FailedStatus     entity.Status = "failed"
	DeclinedStatus   entity.Status = "declined"
	ReversedStatus   entity.Status = "reversed"
)

type Refund struct {
	ID                       uuid.UUID `json:"id"`
	Card                     entity.Card
	Merchant                 entity.CardAcceptor
	Psp                      entity.PSP
	Status                   entity.Status
	Amount                   int
	Currency                 string
	LocalTransactionDateTime time.Time
	Source                   entity.Source
	AuthorizationType        entity.AuthorizationType
	CustomerReference        string
	Stan                     int
	InstitutionID            int
	TraceID                  string
}

type Permission struct {
	ID     string
	Code   string
	Label  string
	ApiKey string
}

type Info struct {
	Low      string
	High     string
	Scheme   string
	Programs []Program
}

// Program
type Program struct {
	ID       string
	Priority int
	Issuer   Issuer
}

// Issuer.
type Issuer struct {
	Name        string
	Country     string
	Website     string
	Phonenumber string
}

type CardInfo struct {
	Info Info
}

// ThreeDSecure contains the output of a 3DS PARes check for any version
type ThreeDSecure struct {
	AuthenticationVerificationValue string
	EcommerceIndicator              data.EcommerceIndicator
	DirectoryServerID               string
	Version                         string
}

func StringToIota(p string) int {
	switch p {
	case "approved":
		return 0
	case "declined":
		return 1
	case "authorized":
		return 3
	default:
		return 2
	}
}

type AuthorizationResult struct {
	Status                               entity.AuthorizationStatus
	ResponseCode                         string
	AdditionalResponseData               string
	AuthorizationIDResponse              string
	FinancialNetworkCode                 string
	BanknetReferenceNumber               string
	NetworkReportingDate                 string
	TraceID                              string
	TransmissionDate                     time.Time
	ThreeDomainSecureSLI                 data.EcommerceIndicator
	ThreeDomainSecureOriginalSLI         data.EcommerceIndicator
	ThreeDomainSecureReasonUCAFDowngrade string
}

type RefundResult struct {
	Status                 entity.AuthorizationStatus
	ResponseCode           string
	Reference              string
	FinancialNetworkCode   string
	BanknetReferenceNumber string
	NetworkReportingDate   string
	TraceID                string
	TransmissionDate       time.Time
}

type CardDataInput string
