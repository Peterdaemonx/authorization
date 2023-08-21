package cis

type DE54_AmountsAdditional struct {
	SF1_AccountType  string `iso8583:"1=n-2"`
	SF2_AmountType   string `iso8583:"2=n-2"`
	SF3_CurrencyCode string `iso8583:"3=n-3"`
	SF4_AmountSign   string `iso8583:"4=a-1"`
	SF5_Amount       int64  `iso8583:"5=n-12, justify=right"`
}
