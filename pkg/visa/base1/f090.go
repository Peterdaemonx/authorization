package base1

type F090_OriginalDataElements struct {
	SF1_OriginalMessageType             string `iso8583:"1=n-4, omitempty"`
	SF2_OriginalTraceNumber             string `iso8583:"2=n-6, justify=right"`
	SF3_OriginalTransmissionDateTime    string `iso8583:"3=n-10, justify=right"`
	SF4_OriginalAcquirerID              string `iso8583:"4=n-11, justify=right"`
	SF5_OriginalForwardingInstitutionID string `iso8583:"5=n-11, justify=right"`
}
