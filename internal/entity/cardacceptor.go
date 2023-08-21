package entity

type CardAcceptor struct {
	CategoryCode string
	ID           string
	Name         string
	Address      CardAcceptorAddress
}

type CardAcceptorAddress struct {
	PostalCode  string
	City        string
	CountryCode string
}
