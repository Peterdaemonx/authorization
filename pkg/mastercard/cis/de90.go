package cis

type DE90_OriginalDataElements struct {
	SF1_OriginalMessageTypeIdentifier     string                       `iso8583:"1=n-4"`
	SF2_OriginalSystemTraceAuditNumber    string                       `iso8583:"2=n-6"`
	SF3_OriginalTransmissionDateAndTime   *DE7_TransmissionDateAndTime `iso8583:"3=n-10"`
	SF4_OriginalAcquiringInstituteIdCode  string                       `iso8583:"4=n-11, justify=right"`
	SF5_OriginalForwardingInstituteIdCode string                       `iso8583:"5=n-11, justify=right"`
}
