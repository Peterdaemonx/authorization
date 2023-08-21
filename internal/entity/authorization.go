package entity

import (
	"time"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
)

type Status string

const (
	Approved     Status = "approved"
	Declined     Status = "declined"
	Failed       Status = "failed"
	RiskDeclined Status = "riskDeclined"
)

type ThreeDSecure struct {
	Version                         string
	AuthenticationVerificationValue string
	DirectoryServerID               string
	EcommerceIndicator              int
}

func (tds ThreeDSecure) NotSet() bool {
	return tds == (ThreeDSecure{})
}

type Authorization struct {
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
	Recurring                Recurring
	Card                     Card
	CardAcceptor             CardAcceptor
	Psp                      PSP
	Exemption                ExemptionType
	ThreeDSecure             ThreeDSecure
	CardSchemeData           CardSchemeData
	CitMitIndicator          CitMitIndicator
	MastercardSchemeData     MastercardSchemeData
	VisaSchemeData           VisaSchemeData
	Terminal                 Terminal
}

type AuthorizationType string

const (
	PreAuthorization   AuthorizationType = "preAuthorization"
	FinalAuthorization AuthorizationType = "finalAuthorization"
)

var mapAuthorizationType = map[string]AuthorizationType{
	`preAuthorization`:   PreAuthorization,
	`finalAuthorization`: FinalAuthorization,
}

func IsValidAuthorizationType(at string) bool {
	_, ok := mapAuthorizationType[at]
	return ok
}

type AuthorizationStatus int

const (
	AuthorizeApproved AuthorizationStatus = iota
	AuthorizeDeclined
	AuthorizeFailed
	AuthorizeStatusUnknown
)

func (a AuthorizationStatus) String() string {
	switch a {
	case AuthorizeApproved:
		return "approved"
	case AuthorizeDeclined:
		return "declined"
	case AuthorizeFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func AuthorizationStatusFromCardSchemeResponseCode(code string) AuthorizationStatus {
	switch code {
	case "00", "08", "10":
		return AuthorizeApproved
	case "91", "92", "96":
		return AuthorizeFailed
	default:
		return AuthorizeDeclined
	}
}

func AuthorizationStatusFromString(code string) AuthorizationStatus {
	switch code {
	case "approved":
		return AuthorizeApproved
	case "declined":
		return AuthorizeDeclined
	case "failed":
		return AuthorizeFailed
	default:
		return AuthorizeStatusUnknown
	}
}
