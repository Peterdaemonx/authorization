package base1

type F043_CardAcceptorNameLocation struct {
	SF1_CarAcceptorName  string `iso8583:"1=ans-25, justify=left"`
	SF2_CardAcceptorCity string `iso8583:"2=ans-13, justify=left"`
	SF3_CountryCode      string `iso8583:"3=ans-2"`
}
