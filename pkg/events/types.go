package events

import (
	"time"
)

type PSP struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CardAcceptor struct {
	Name           string `json:"name"`
	City           string `json:"city"`
	Country        string `json:"country"`
	CategoryCode   string `json:"categoryCode"`
	CardAcceptorID string `json:"cardAcceptorID"`
}

type MastercardSchemeData struct {
	AuthorizationType      string             `json:"authorizationType"`
	SystemTraceAuditNumber int                `json:"systemTraceAuditNumber"`
	TraceID                string             `json:"traceId"`
	FinancialNetworkCode   string             `json:"financialNetworkCode"`
	PosPinCaptureCode      string             `json:"posPinCaptureCode,omitempty"` // CIS DE26
	POSEntryMode           POSEntryMode       `json:"posEntryMode,omitempty"`
	AdditionalData         AdditionalData     `json:"additionalData"`
	PointOfServiceData     PointOfServiceData `json:"pointOfServiceData"`
	ResponseReference      string             `json:"responseReference,omitempty"`
	BanknetReferenceNumber string             `json:"banknetReferenceNumber"`
	NetworkReportingDate   string             `json:"networkReportingDate"`
}

// ProcessingCode CIS DE03
type ProcessingCode struct {
	// SF1, SF2 and SF3
	TransactionTypeCode string `json:"transactionTypeCode"`
	FromAccountTypeCode string `json:"fromAccountTypeCode"`
	ToAccountTypeCode   string `json:"toAccountTypeCode"`
}

// POSEntryMode DE 22—Point-of-Service (POS) Entry Mode
type POSEntryMode struct {
	// SF 1 POS Terminal PAN Entry Mode
	PanEntryMode string `json:"panEntryMode,omitempty"`
	// SF 2—POS Terminal PIN Entry Mode
	PinEntryMode string `json:"pinEntryMode,omitempty"`
}

// SE 42 Electronic Commerce Indicators
type EcommerceIndicators struct {
	// SF 1 = Security Protocol
	SecurityProtocol int `json:"securityProtocol"`
	// SF 2 Cardholder Authentication
	CardholderAuthentication int `json:"cardholderAuthentication"`
	// SF 3 UCAF Collection Indicator
	UCAFCollectionIndicator int `json:"ucafCollectionIndicator"`
}

// AdditionalData DE 48—Additional Data—Private Use
type AdditionalData struct {
	// SE 80 PIN Service Code
	PinServiceCode string `json:"pinServiceCode,omitempty"`
	// SE 42 Electronic Commerce Indicators
	EcommerceIndicators EcommerceIndicators `json:"eCommerceIndicators"`
}

// PointOfServiceData DE 61—Point-of-Service (POS) Data
type PointOfServiceData struct {
	// SF 1—POS Terminal Attendance
	TerminalAttendance int `json:"terminalAttendance"`
	// SF 3—POS Terminal Location
	TerminalLocation int `json:"terminalLocation"`
	// SF 4—POS Cardholder Presence
	CardHolderPresence int `json:"cardHolderPresence"`
	// SF 5—POS Card Presence
	CardPresence int `json:"cardPresence"`
	// SF 6—POS Card Capture Capabilities
	CardCaptureCapabilities int `json:"cardCaptureCapabilities"`
	// SF 10—Cardholder-Activated Terminal Level
	CardHolderActivatedTerminalLevel int `json:"cardHolderActivatedTerminalLevel"`
	// SF 11—POS Card Data Terminal Input Capability Indicator
	CardDataTerminalInputCapabilityIndicator int `json:"cardDataTerminalInputCapabilityIndicator"`
}

type AuthorizationCaptureMessageV1 struct {
	AuthorizationID                        string                `json:"authorizationID"`
	CaptureID                              string                `json:"captureID"`
	PSP                                    PSP                   `json:"psp"`
	PanTokenID                             string                `json:"panTokenID"`
	MaskedPan                              string                `json:"maskedPan"`
	CardScheme                             string                `json:"cardScheme"`
	CardIssuerCountry                      string                `json:"cardIssuerCountry"`
	CardProductID                          string                `json:"cardProductID"`
	CardProgramID                          string                `json:"cardProgramID"`
	Amount                                 int                   `json:"amount"`
	Currency                               string                `json:"currency"`
	LocalTransactionDateTime               time.Time             `json:"localTransactionDateTime"`
	PartialCapture                         bool                  `json:"partialCapture"`
	CapturedAt                             time.Time             `json:"capturedAt"`
	ProcessingDate                         time.Time             `json:"processingDate"`
	ThreeDSVersion                         string                `json:"threeDSVersion,omitempty"`
	ThreeDSAuthenticationVerificationValue string                `json:"threeDSAuthenticationVerificationValue,omitempty"`
	ThreeDSDirectoryServerTransactionID    string                `json:"threeDSDirectoryServerTransactionID,omitempty"`
	CustomerReference                      string                `json:"customerReference,omitempty"`
	ResponseCode                           string                `json:"responseCode,omitempty"`
	Source                                 string                `json:"source"`
	CardAcceptor                           CardAcceptor          `json:"cardAcceptor"`
	ProcessingCode                         ProcessingCode        `json:"processingCode"`
	MastercardSchemeData                   *MastercardSchemeData `json:"mastercardSchemeData,omitempty"`
}

type RefundCaptureMessageV1 struct {
	RefundID                 string                `json:"refundID"`
	CaptureID                string                `json:"captureID"`
	PSP                      PSP                   `json:"psp"`
	PanTokenID               string                `json:"panTokenID"`
	MaskedPan                string                `json:"maskedPan"`
	CardScheme               string                `json:"cardScheme"`
	CardIssuerCountry        string                `json:"cardIssuerCountry"`
	CardProductID            string                `json:"cardProductID"`
	CardProgramID            string                `json:"cardProgramID"`
	Amount                   int                   `json:"amount"`
	Currency                 string                `json:"currency"`
	LocalTransactionDateTime time.Time             `json:"localTransactionDateTime"`
	PartialCapture           bool                  `json:"partialCapture"`
	CapturedAt               time.Time             `json:"capturedAt"`
	ProcessingDate           time.Time             `json:"processingDate"`
	CustomerReference        string                `json:"customerReference,omitempty"`
	ResponseCode             string                `json:"responseCode,omitempty"`
	Source                   string                `json:"source"`
	CardAcceptor             CardAcceptor          `json:"cardAcceptor"`
	ProcessingCode           ProcessingCode        `json:"processingCode"`
	MastercardSchemeData     *MastercardSchemeData `json:"mastercardSchemeData,omitempty"`
}

type CaptureUpdateMessage struct {
	AuthorizationID string    `json:"authorizationID"`
	CaptureID       string    `json:"captureID"`
	Status          string    `json:"status"`
	IRD             string    `json:"ird"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
