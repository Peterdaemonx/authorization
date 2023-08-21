package cis

type DE95_ReplacementAmounts struct {
	SF1_ActualAmountTransaction       int64 `iso8583:"1=n-12, autofill"`
	SF2_ActualAmountSettlement        int64 `iso8583:"2=n-12, autofill"`
	SF3_ActualAmountCardholderBilling int64 `iso8583:"3=n-12, autofill"`
	_                                 int64 `iso8583:"4=n-6, autofill"`
}
