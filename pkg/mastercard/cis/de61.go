package cis

type DE61_PointOfServiceData struct {
	SF1_TerminalAttendance                        string `iso8583:"1=n-1"`
	_                                             string `iso8583:"2=n-1, autofill"`
	SF3_TerminalLocation                          string `iso8583:"3=n-1"`
	SF4_CardholderPresence                        string `iso8583:"4=n-1"`
	SF5_CardPresence                              string `iso8583:"5=n-1"`
	SF6_CardCaptureCapabilities                   string `iso8583:"6=n-1"`
	SF7_TransactionStatus                         string `iso8583:"7=n-1"`
	SF8_TransactionSecurity                       string `iso8583:"8=n-1"`
	_                                             string `iso8583:"9=n-1, autofill"`
	SF10_CardholderActivatedTerminalLevel         string `iso8583:"10=n-1"`
	SF11_CardDataTerminalInputCapabilityIndicator string `iso8583:"11=n-1"`
	SF12_AuthorizationLifeCycle                   string `iso8583:"12=n-2"`
	SF13_CountryCode                              string `iso8583:"13=n-3"`
	SF14_PostalCode                               string `iso8583:"14=ans-10, justify=left"`
}
