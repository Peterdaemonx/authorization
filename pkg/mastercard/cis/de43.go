package cis

type DE43_CardAcceptorNameAndLocation struct {
	SF1_Name               string `iso8583:"1=ans-22, justify=left"`
	_                      string `iso8583:"2=ans-1, autofill"`
	SF3_City               string `iso8583:"3=ans-13, justify=left"`
	_                      string `iso8583:"4=ans-1, autofill"`
	SF5_StateOrCountryCode string `iso8583:"5=ans-3"`
}
