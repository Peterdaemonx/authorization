package ports

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"gitlab.cmpayments.local/creditcard/platform/categorycode"
	"gitlab.cmpayments.local/creditcard/platform/countrycode"
	"gitlab.cmpayments.local/creditcard/platform/http/validator"

	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
)

type Card struct {
	Holder string `json:"holder"`
	Number string `json:"number"`
	Cvv    string `json:"cvv"`
	Expiry Expiry `json:"expiry"`
	scheme string
}

func (c Card) validate(v *validator.Validator) {
	v.Check(len(c.Holder) >= 2 && len(c.Holder) <= 26, "card.holder", []string{"card holder must be min 2 and max 26 characters long"})
	v.Check(len(c.Number) >= 12 && len(c.Number) <= 19, "card.number", []string{"card number must be min 12 and max 19 characters long"})
	v.Check(c.Cvv == "" || len(c.Cvv) == 3, "card.cvv", []string{"length cvv must be 3 digits"})
	c.Expiry.validate(v)
}

type CardResponse struct {
	Number string `json:"number"`
	Scheme string `json:"scheme"`
}

type Expiry struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}

func (e Expiry) validate(v *validator.Validator) {
	t, err := time.Parse("0601", fmt.Sprintf("%s%s", e.Year, e.Month))
	v.Check(err == nil, "card.expiry", []string{"invalid expiry format"})
	v.Check(t.After(time.Now()), "card.expiry", []string{"card expired"})
}

type CardAcceptor struct {
	ID           string `json:"id"`
	CategoryCode string `json:"categoryCode"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Country      string `json:"country"`
	PostalCode   string `json:"postalCode"`
}

func (ca CardAcceptor) validate(v *validator.Validator) {
	_, ok := categorycode.MCCS[ca.CategoryCode]
	v.Check(ca.ID != "", "cardAcceptor.id", []string{"cardAcceptor id cannot be empty"})
	v.Check(len(ca.ID) <= 12, "cardAcceptor.id", []string{"cardAcceptor id must be max 12 characters long"})
	v.Check(ok, "cardAcceptor.categoryCode", []string{"merchant category code not found."})
	v.Check(ca.ID != "", "cardAcceptor.name", []string{"cardAcceptor name cannot be empty"})
	v.Check(len(ca.Name) <= 22, "cardAcceptor.name", []string{"cardAcceptor name must be max 22 characters long"})
	v.Check(len(ca.City) <= 13, "cardAcceptor.city", []string{"cardAcceptor city must be max 13 characters long"})
	c, err := countrycode.GetCountry(ca.Country)
	if err != nil {
		v.AddError("cardAcceptor.country", []string{"invalid cardAcceptor country"})
	}

	if !c.EEACountry() && err == nil {
		v.AddError("cardAcceptor.country", []string{"cardAcceptor country not allowed, must be in EEA"})
	}

	if countrycode.CountryHasPostalCode(ca.Country) {
		v.Check(ca.PostalCode != "", "cardAcceptor.postalCode", []string{fmt.Sprintf("postal code cannot be empty for %s", ca.Country)})
		v.Check(len(ca.PostalCode) <= 10, "cardAcceptor.postalCode", []string{"postal code must be max 10 characters long"})
	}
}

type CitMitIndicator struct {
	InitiatedBy string `json:"initiatedBy,omitempty"`
	SubCategory string `json:"subCategory,omitempty"`
}

func (cof CitMitIndicator) validate(scheme string, v *validator.Validator) {
	if (cof != CitMitIndicator{}) {
		v.Check(entity.IsValidInitiatedBy(cof.InitiatedBy), "citMitIndicator.initiatedBy", []string{"invalid initiated by"})
		v.Check(entity.IsValidSubCategory(cof.SubCategory), "citMitIndicator.subCategory", []string{"invalid sub category"})
		v.Check(entity.IsValidCitMit(scheme, cof.InitiatedBy, cof.SubCategory), "citMitIndicator", []string{"invalid combination"})
	}
}

type ThreeDSecure struct {
	AuthenticationVerificationValue string                  `json:"authenticationVerificationValue"`
	Version                         string                  `json:"version"`
	EcommerceIndicator              data.EcommerceIndicator `json:"ecommerceIndicator"`
	DirectoryServerTransactionID    string                  `json:"directoryServerTransactionId"`
}

func (tds ThreeDSecure) validate(v *validator.Validator) {
	if (tds != ThreeDSecure{}) {
		v.Check(validator.IntIn(int(tds.EcommerceIndicator), 0, 1, 2, 5, 6, 7), "threeDSecure.ecommerceIndicator", []string{"invalid ecommerce indicator"})
		v.Check(tds.AuthenticationVerificationValue != "", "threeDSecure.AuthenticationVerificationValue", []string{"authentication verification value cannot be empty"})
		v.Check(len(strings.Trim(tds.DirectoryServerTransactionID, " ")) == 36, "threeDSecure.directoryServerTransactionId", []string{"directory server transaction id must be 36 characters long"})
	}
}

type ThreeDSecureResponse struct {
	Version                      string `json:"version,omitempty"`
	EcommerceIndicator           string `json:"ecommerceIndicator"`
	DirectoryServerTransactionID string `json:"directoryServerTransactionId,omitempty"`
}

type authorizationRequest struct {
	Amount                   int                           `json:"amount"`
	Currency                 string                        `json:"currency"`
	Reference                string                        `json:"reference"`
	Source                   string                        `json:"source"`
	LocalTransactionDateTime data.LocalTransactionDateTime `json:"localTransactionDateTime"`
	AuthorizationType        string                        `json:"authorizationType,omitempty"`
	InitialRecurring         bool                          `json:"initialRecurring"`
	InitialTraceID           string                        `json:"initialTraceId"`
	Card                     Card                          `json:"card"`
	CardAcceptor             CardAcceptor                  `json:"cardAcceptor"`
	CitMitIndicator          CitMitIndicator               `json:"citMitIndicator"`
	Exemption                string                        `json:"exemption"`
	ThreeDSecure             ThreeDSecure                  `json:"threeDSecure"`
}

func (a authorizationRequest) validate(v *validator.Validator) {
	v.Check(a.Amount > 0, "amount", []string{"amount must be greater than 0"})
	v.Check(a.Amount < 3000000, "amount", []string{"amount must be less than 3000000"})
	c, err := currencycode.GetCurrency(a.Currency)
	v.Check(err == nil, "currency", []string{"unsupported currency"})
	if a.Card.scheme == entity.Visa && !c.AllowedByVisa() {
		v.AddError("currency", []string{"currency is not allowed by visa"})
	}
	if a.Card.scheme == entity.Mastercard && !c.AllowedByMastercard() {
		v.AddError("currency", []string{"currency is not allowed by mastercard"})
	}
	v.Check(len(a.Reference) <= 100, "reference", []string{"reference max length 100"})
	v.Check(entity.IsValidSource(a.Source), "source", []string{"invalid source"})
	v.Check(a.AuthorizationType == "" || entity.IsValidAuthorizationType(a.AuthorizationType), "authorizationType", []string{"invalid authorization type"})
	v.Check(a.InitialTraceID == "" || (len(a.InitialTraceID) >= 13 && len(a.InitialTraceID) <= 15), "initialTraceId", []string{"initial trace id must be min 13 and max 15 characters long"})
	v.Check(a.InitialTraceID == "" || (regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(a.InitialTraceID)), "initialTraceId", []string{"initial trace id must be alpha numeric"})
	a.Card.validate(v)
	a.CardAcceptor.validate(v)
	a.CitMitIndicator.validate(a.Card.scheme, v)
	v.Check(a.Exemption == "" || entity.IsValidExemption(a.Exemption), "exemption", []string{"invalid exemption"})
	a.ThreeDSecure.validate(v)
}

type CardSchemeResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId,omitempty"`
}

type authorizationResponse struct {
	ID                       string                         `json:"id"`
	LogID                    string                         `json:"logID,omitempty"`
	Amount                   int                            `json:"amount"`
	Currency                 string                         `json:"currency"`
	Reference                string                         `json:"reference"`
	Source                   string                         `json:"source"`
	LocalTransactionDateTime *data.LocalTransactionDateTime `json:"localTransactionDateTime"`
	AuthorizationType        string                         `json:"authorizationType,omitempty"`
	InitialRecurring         bool                           `json:"initialRecurring"`
	TraceID                  string                         `json:"traceId,omitempty"`
	ProcessingDate           string                         `json:"processingDate"`
	Card                     CardResponse                   `json:"card"`
	CardAcceptor             CardAcceptor                   `json:"cardAcceptor"`
	Exemption                string                         `json:"exemption,omitempty"`
	ThreeDSecure             ThreeDSecureResponse           `json:"threeDSecure"`
	CardSchemeResponse       CardSchemeResponse             `json:"cardSchemeResponse"`
	CitMitIndicator          *CitMitIndicator               `json:"citMitIndicator,omitempty"`
	// The mastercard specific data
	//
	// required: false
	MasterCardData *MasterCardDataResponse `json:"masterCardData,omitempty"`
}

type MasterCardDataResponse struct {
	DE61 DE61   `json:"de61"`
	DE22 DE22   `json:"de22"`
	DE3  string `json:"de3"`
}

type DE61 struct {
	AuthorizationLifeCycle                   string `json:"authorizationLifeCycle"`
	CardCaptureCapabilities                  int    `json:"cardCaptureCapabilities"`
	CardDataTerminalInputCapabilityIndicator int    `json:"cardDataTerminalInputCapabilityIndicator"`
	CardHolderActivatedTerminalLevel         int    `json:"cardHolderActivatedTerminalLevel"`
	CardHolderPresence                       int    `json:"cardHolderPresence"`
	CardPresence                             int    `json:"cardPresence"`
	CountryCode                              string `json:"countryCode"`
	PostalCode                               string `json:"postalCode"`
	TerminalAttendance                       int    `json:"terminalAttendance"`
	TerminalLocation                         int    `json:"terminalLocation"`
	TransactionSecurity                      int    `json:"transactionSecurity"`
	TransactionStatus                        int    `json:"transactionStatus"`
}

type DE22 struct {
	SF1 string `json:"sf1"`
}

type Metadata struct {
	CurrentPage int  `json:"currentPage,omitempty"`
	PageSize    int  `json:"pageSize,omitempty"`
	FirstPage   int  `json:"firstPage,omitempty"`
	LastPage    bool `json:"lastPage,omitempty"`
}

type AuthorizationsResponse struct {
	Metadata       Metadata                `json:"metadata"`
	Authorizations []authorizationResponse `json:"authorizations"`
}
