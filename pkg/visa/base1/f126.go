package base1

type F126_PrivateUseFields struct {
	SF1_UnusedReserved                      string                        `iso8583:"1=ans-25, minlength=0, omitempty"` // fields 1 to 4 are the same, total length 155 (workaround for subfield definitions constrain)
	SF2_UnusedReserved                      string                        `iso8583:"2=ans-57, minlength=0, omitempty"`
	SF3_UnusedReserved                      string                        `iso8583:"3=ans-57, minlength=0, omitempty"`
	SF4_UnusedReserved                      string                        `iso8583:"4=ans-17, minlength=0, omitempty"`
	SF5_MerchantIdentifier                  string                        `iso8583:"5=ans-8, minlength=0, omitempty"`
	SF6_CardholderCertificateSerialNumber   string                        `iso8583:"6=b-17, minlength=0, omitempty"`
	SF7_MerchantCertificateSerialNumber     string                        `iso8583:"7=b-17, minlength=0, omitempty"`
	SF8_TransactionID                       string                        `iso8583:"8=b-20,  minlength=0, omitempty"`
	SF9_CAVVData                            string                        `iso8583:"9=n-40, dataenc=bcd4, omitempty, minlength=0"`
	SF10_CVV2AuthorizationRequestData       string                        `iso8583:"10=ans-6, dataenc=ebcdic, minlength=0, omitempty"`
	_                                       string                        `iso8583:"11=ans-0, minlength=0, omitempty"`
	SF12_ServiceIndicators                  string                        `iso8583:"12=n-24, minlength=0, omitempty"`
	SF13_POSEnvironment                     string                        `iso8583:"13=an-1, dataenc=ebcdic, minlength=0, omitempty"`
	_                                       string                        `iso8583:"14=an-1, minlength=0, omitempty"`
	SF15_MastercardUCAFCollectionIndicator  string                        `iso8583:"15=ans-1, dataenc=ebcdic, minlength=0, omitempty"`
	SF16_MastercardUCAFField                string                        `iso8583:"16=ans-33, dataenc=ebcdic, minlength=0, omitempty"`
	_                                       string                        `iso8583:"17=an-1, minlength=0, omitempty"`
	SF18_AgentUniqueAccountResult           SF18_AgentUniqueAccountResult `iso8583:"18=b-12, minlength=0, omitempty"`
	SF19_DynamicCurrencyConversionIndicator string                        `iso8583:"19=ans-1, dataenc=ebcdic, minlength=0, omitempty"`
	SF20_3DSecureIndicator                  string                        `iso8583:"20=an-1, dataenc=ebcdic, minlength=0, omitempty"`
	_                                       string                        `iso8583:"21=ans-1, minlength=0, omitempty"`
	_                                       string                        `iso8583:"22=ans-1, minlength=0, omitempty"`
	_                                       string                        `iso8583:"23=ans-1, minlength=0, omitempty"`
	_                                       string                        `iso8583:"24=ans-1, minlength=0, omitempty"`
}

type SF18_AgentUniqueAccountResult struct {
	// Position 1 contains a fixed binary value: 11 (0B)
	PS1           string `iso8583:"1=b-1"`
	AgentUniqueId string `iso8583:"2=b-5, dataenc=ebcdic"`
	// Reserved, filled with zeros
	_ string `iso8583:"3=b-6, justify=right"`
}
