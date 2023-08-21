package ports

import (
	"fmt"

	"gitlab.cmpayments.local/creditcard/platform/categorycode"
	"gitlab.cmpayments.local/creditcard/platform/countrycode"
	"gitlab.cmpayments.local/creditcard/platform/http/validator"

	authorizationPorts "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/ports"
	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
)

type Card struct {
	Number string `json:"number"`
	Expiry Expiry `json:"expiry"`
	scheme string
}

type Expiry struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}

func (c Card) validate(v *validator.Validator) {
	v.Check(len(c.Number) >= 12 && len(c.Number) <= 19, "card.number", []string{"card number must be min 12 and max 19 characters long"})
}

type CardResponse struct {
	Scheme string `json:"scheme"`
	Number string `json:"number"`
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
	if !c.EEACountry() {
		v.AddError("cardAcceptor.country", []string{"cardAcceptor country not part of EEA"})
	}

	if countrycode.CountryHasPostalCode(ca.Country) {
		v.Check(ca.PostalCode != "", "cardAcceptor.postalCode", []string{fmt.Sprintf("postal code cannot be empty for %s", ca.Country)})
		v.Check(len(ca.PostalCode) <= 10, "cardAcceptor.postalCode", []string{"postal code must be max 10 characters long"})
	}
}

type refundRequest struct {
	Amount                   int                           `json:"amount"`
	Currency                 string                        `json:"currency"`
	Reference                string                        `json:"reference"`
	Source                   string                        `json:"source"`
	LocalTransactionDateTime data.LocalTransactionDateTime `json:"localTransactionDateTime"`
	AuthorizationType        string                        `json:"authorizationType,omitempty"`
	Card                     Card                          `json:"card"`
	CardAcceptor             CardAcceptor                  `json:"cardAcceptor"`
}

func (a refundRequest) validate(v *validator.Validator) {
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
	v.Check(entity.IsValidSource(a.Source), "source", []string{"invalid source"})
	v.Check(a.AuthorizationType == "" || entity.IsValidAuthorizationType(a.AuthorizationType), "authorizationType", []string{"invalid authorization type"})
	a.Card.validate(v)
	a.CardAcceptor.validate(v)
}

type CardSchemeResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId,omitempty"`
}

type refundResponse struct {
	ID                       string                         `json:"id"`
	LogID                    string                         `json:"logID,omitempty"`
	Amount                   int                            `json:"amount"`
	Currency                 string                         `json:"currency"`
	Reference                string                         `json:"reference"`
	Source                   string                         `json:"source"`
	LocalTransactionDateTime *data.LocalTransactionDateTime `json:"localTransactionDateTime"`
	AuthorizationType        string                         `json:"authorizationType,omitempty"`
	ProcessingDate           string                         `json:"processingDate"`
	Card                     CardResponse                   `json:"card"`
	CardAcceptor             CardAcceptor                   `json:"cardAcceptor"`
	CardSchemeResponse       CardSchemeResponse             `json:"cardSchemeResponse"`
}

type RefundsResponse struct {
	Metadata authorizationPorts.Metadata `json:"metadata"`
	Refunds  []refundResponse            `json:"refunds"`
}
