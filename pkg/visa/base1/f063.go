package base1

type F063_NetworkData struct {
	SF1_NetworkID         string `iso8583:"1=n-4, dataenc=bcd4, bcd4len=4, minlength=0, omitempty"`
	_                     string `iso8583:"2=n-4, minlength=0, omitempty"`
	SF3_MessageReasonCode string `iso8583:"3=n-4, dataenc=bcd4, omitempty"`
}
