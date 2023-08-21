package cis

type DE108_AdditionalTransactionReferenceData struct {
	_ string `iso8583:"1=ans...322"` // no idea why this is used //nolint:gofmt
	_ string `iso8583:"2=ans...322"` // no idea why this is used //nolint:gofmt
	// actual data
	SE03_TransactionReferenceData *DE108_SE03_TransactionReferenceData `iso8583:"3=ans...138"`
}

type DE108_SE03_TransactionReferenceData struct {
	SF01_UniqueTransactionReference string `iso8583:"1=ans..19"`
}
