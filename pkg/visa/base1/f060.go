package base1

type F060_AdditionalPOSInformation struct {
	// B1 Subfield 1 & 2
	B1 string `iso8583:"1=n-2, dataenc=bcd4, minlength=2"`
	// B2 Subfield 3 & 4
	B2 string `iso8583:"2=n-2, dataenc=bcd4, minlength=2,omitempty"`
	// B3 Subfield 5
	B3 string `iso8583:"3=n-2, dataenc=bcd4, minlength=2,omitempty"`
	// B4 Subfield 6 & 7
	B4 string `iso8583:"4=n-2, dataenc=bcd4, minlength=2,omitempty"`
	// B5 Subfield 8
	B5 string `iso8583:"5=n-2, dataenc=bcd4, minlength=2,omitempty"`
	// B6 Subfield 9 & 10
	B6 string `iso8583:"6=n-2, dataenc=bcd4, minlength=2,omitempty"`
}
