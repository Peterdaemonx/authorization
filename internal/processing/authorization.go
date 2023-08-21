package processing

import "gitlab.cmpayments.local/creditcard/authorization/internal/entity"

type AuthorizationEntity struct {
	Authorization entity.Authorization
	SchemeRequest *AuthorizationSchemeRequest
	Response      AuthorizationResult
}

type AuthorizationSchemeRequest struct {
	ProcessingCode           entity.ProcessingCode
	POSEntryMode             entity.POSEntryMode
	PosPinCaptureCode        string // DE 26—Point-of-Service (POS) Personal ID Number (PIN) Capture Code
	AdditionalData           entity.AdditionalRequestData
	PointOfServiceData       entity.PointOfServiceData
	PosConditionCode         string // F025 Point-of-service Condition Code visa-net-authorization-only-online-messages-technical-specifications.pdf page 189
	AdditionalPOSInformation entity.AdditionalPOSInformation
}

// VIPPrivateUsedField F063 VIP Private-Use Field visa-net-authorization-only-online-messages-technical-specifications.pdf page 390
type VIPPrivateUsedField struct {
	// Field 63.0 Bitmap length 24 bit string 3 bytes
	// visa-net-authorization-only-online-messages-technical-specifications.pdf page 392
	Bitmap string
	// Field 63.1 NetworkIdentificationCode length 4, 4 bit BCD (unsigned pack); 2 bytes.
	// visa-net-authorization-only-online-messages-technical-specifications.pdf page 393
	NetworkIdentificationCode string
}
