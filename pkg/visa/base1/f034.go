package base1

type F034_ElectronicCommerceData struct {
	HEX01_AuthenticationData                  HEX01_AuthenticationData                  `iso8583:"1=b....9, dataenc=tlv, lenenc=hexBit4, tlvTag=01, omitempty"` // length = dataset ID + data
	HEX02_AcceptanceEnvironmentAdditionalData HEX02_AcceptanceEnvironmentAdditionalData `iso8583:"2=b....3, dataenc=tlv, lenenc=hexBit4, tlvTag=02, omitempty"`
	HEX4A_StrongConsumerAuthentication        HEX4A_StrongConsumerAuthentication        `iso8583:"3=b....12, dataenc=tlv, lenenc=hexBit4, tlvTag=4A, omitempty"`
}

// HEX01_AuthenticationData Dataset ID Hex 01
type HEX01_AuthenticationData struct {
	T86_3DSecureProtocolVersionNumber string `iso8583:"1=b..8, dataenc=ebcdic, lenenc=hexBit4, minlength=5, tlvTag=86, omitempty"` // Tag 86
}

// HEX02_AcceptanceEnvironmentAdditionalData Dataset ID Hex 02
type HEX02_AcceptanceEnvironmentAdditionalData struct {
	T80_InitiatingPartyIndicator string `iso8583:"1=b..1, dataenc=ebcdic, lenenc=hexBit4, tlvTag=80, omitempty"` // Tag 80
}

// HEX4A_StrongConsumerAuthentication Dataset ID Hex 4A
type HEX4A_StrongConsumerAuthentication struct {
	T87_LowValueExemptionIndicator                string `iso8583:"1=b..1, dataenc=ebcdic, lenenc=hexBit4, tlvTag=87, omitempty"` // Tag 87
	T88_SecureCorporatePaymentIndicator           string `iso8583:"2=b..1, dataenc=ebcdic, lenenc=hexBit4, tlvTag=88, omitempty"` // Tag 88
	T89_TransactionRiskAnalysisExemptionIndicator string `iso8583:"3=b..1, dataenc=ebcdic, lenenc=hexBit4, tlvTag=89, omitempty"` // Tag 89
	T8A_DelegatedAuthenticationIndicator          string `iso8583:"4=b..1, dataenc=ebcdic, lenenc=hexBit4, tlvTag=8A, omitempty"` // Tag 8A
}
