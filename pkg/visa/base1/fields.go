package base1

type Fields struct {
	F002_PrimaryAccountNumber                   string                        `iso8583:"2=n..19, dataenc=bcd4, lenenc=hex, justify=right"`
	F003_ProcessingCode                         string                        `iso8583:"3=n-6, dataenc=bcd4"`
	F004_TransactionAmount                      int64                         `iso8583:"4=n-12, dataenc=bcd4, lenenc=hex, justify=right"`
	F007_TransmissionDateTime                   F007_TransmissionDateAndTime  `iso8583:"7=n-10, dataenc=bcd4"`
	F011_SystemTraceAuditNumber                 string                        `iso8583:"11=n-6, dataenc=bcd4"`
	F012_LocalTransactionTime                   string                        `iso8583:"12=n-6, dataenc=bcd4"`
	F013_LocalTransactionDate                   string                        `iso8583:"13=n-4, dataenc=bcd4"`
	F014_ExpirationDate                         string                        `iso8583:"14=n-4, dataenc=bcd4"`
	F018_MerchantType                           string                        `iso8583:"18=n-4, dataenc=bcd4"`
	F019_AcquiringInstituteCountryCode          string                        `iso8583:"19=n-3, dataenc=bcd4"` // Docs say it is 3, but the test tool gives us 4
	F022_PosEntryMode                           string                        `iso8583:"22=n-4, dataenc=bcd4"`
	F025_PosCondition                           string                        `iso8583:"25=n-2, dataenc=bcd4"`
	F028_TransactionFeeAmount                   string                        `iso8583:"28=an-9, dataenc=ebcdic"`                           //	Should be an-9 - seems hex encoded?
	F032_AcquiringInstitutionIdentificationCode string                        `iso8583:"32=n..12, dataenc=bcd4, lenenc=hex, justify=right"` // variable length 1 byte, binary + 11n, 4-bit BCD (unsigned packed); maximum 7 bytes
	F034_ElectronicCommerceData                 F034_ElectronicCommerceData   `iso8583:"34=b....65537, lenenc=hexBit4, omitempty"`
	F037_RetrievalReferenceNumber               string                        `iso8583:"37=an-12, dataenc=ebcdic"` // format: ydddnnnnnnnn
	F038_AuthorizationIdenticationResponse      string                        `iso8583:"38=an-6, dataenc=ebcdic"`
	F039_ResponseCode                           string                        `iso8583:"39=an-2, dataenc=ebcdic"`
	F041_CardAcceptorTerminalIdentification     string                        `iso8583:"41=ans-8, dataenc=ebcdic"` // . 1 byte for length .. 2 bytes for length ... 3 bytes for length
	F042_CardAcceptorIdentificationCode         string                        `iso8583:"42=ans-15, dataenc=ebcdic, justify=left"`
	F043_CardAcceptorNameLocation               F043_CardAcceptorNameLocation `iso8583:"43=ans-40, dataenc=ebcdic"`
	F044_AdditionalResponseData                 F044_AdditionalResponseData   `iso8583:"44=ans.25, dataenc=ebcdic, lenenc=bin"`
	F049_TransactionCurrencyCode                string                        `iso8583:"49=n-3, dataenc=bcd4"`
	F060_AdditionalPointOfServiceInformation    F060_AdditionalPOSInformation `iso8583:"60=b.7, lenenc=bin"`
	F062_CustomPaymentServiceFields             F062_CustomPaymentService     `iso8583:"62=b.255, lenenc=bin, subbitmap=8"`
	F063_NetworkData                            F063_NetworkData              `iso8583:"63=b.79, lenenc=bin, subbitmap=3"`
	F070_NetworkManagementInformationCode       string                        `iso8583:"70=n-3, dataenc=bcd4"`
	F090_OriginalDataElements                   F090_OriginalDataElements     `iso8583:"90=n-42, dataenc=bcd4"`
	F126_PrivateUseFields                       F126_PrivateUseFields         `iso8583:"126=b.255, lenenc=bin, subbitmap=8"`
}
