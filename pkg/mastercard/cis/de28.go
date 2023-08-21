package cis

type DE28_TransactionFeeAmount struct {
	SF1_DebitCreditIndicator string `iso8583:"1=a-1"`
	SF2_Amount               int64  `iso8583:"2=n-8"`
}
