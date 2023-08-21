package mastercard

type ThreeDomainSecureVersion int

const (
	ThreeDomainSecureNone ThreeDomainSecureVersion = iota // 3D Secure not used
	ThreeDomainSecureV1                                   // 3D Secure Version 1.0 / SecureCode program
	ThreeDomainSecureV2                                   // EMV 3-D Secure v2 / Mastercard Identity Check program
)

type ThreeDomainSecure struct {
	// Used version of 3dsecure
	Version ThreeDomainSecureVersion
	// Universal Cardholder Authentication Field security protocol
	ElectronicCommerceSecurityLevelIndicator string
	// If 3dsecure is used, the authentication value
	AccountholderAuthenticationValue string
	// If 3DSv2 is used then a directoryservertransactionID is mandatory
	DirectoryServerTransactionID string
}
