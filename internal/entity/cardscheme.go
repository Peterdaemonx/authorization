package entity

import (
	"fmt"
	"strconv"
)

type CardSchemeData struct {
	Request  CardSchemeRequest
	Response CardSchemeResponse
}

type CardSchemeRequest struct {
	ProcessingCode               ProcessingCode
	CardHolderVerificationMethod CardHolderVerificationMethod
	RetrievalReferenceNumber     string
	POSEntryMode                 POSEntryMode
}

type MastercardSchemeData struct {
	Request  MastercardSchemeRequest
	Response MastercardSchemeResponse
}

type VisaSchemeData struct {
	Request  VisaSchemeRequest
	Response VisaSchemeResponse
}

// DE03 Processing-Code
type ProcessingCode struct {
	TransactionTypeCode string // SF1
	FromAccountTypeCode string // SF2
	ToAccountTypeCode   string // SF3
}

func (pc ProcessingCode) String() string {
	return fmt.Sprintf("%s%s%s", pc.TransactionTypeCode, pc.FromAccountTypeCode, pc.ToAccountTypeCode)
}

// DE22 Point-of-Service (POS) Entry Mode
type POSEntryMode struct {
	PanEntryMode PANEntry // SF1 POS Terminal PAN Entry Mode
	PinEntryMode PINEntry // SF2 POS Terminal PIN Entry Mode
}

// PANEntry This field gives information on how the Card was captured at the point of service (magstipe, chip, manual, etc)
//
// The following fields will be mapped accordingly when send to mastercard and visa:
// default 	  		     				= empty string
// PANEntryUnknown							 	= code send to mastercard/visa: 00
// PANEntryManual        				= PAN manual entry, code send to mastercard/visa: 01
// PANEntryAutoMagStripeNotRequired		= PAN auto-entry via magnetic stripe—track data is not required, code send to mastercard/visa: 02
// PANEntryChip          		 				= PAN auto-entry via chip, code send to mastercard/visa: 05
// PANEntryContactless   		 				= PAN auto-entry via contactless M/Chip, code send to mastercard/visa: 07
// PANEntryMagneticStrip 		 				= PAN auto-entry via magnetic strip, code send to mastercard/visa: 90
// PANEntryViaEcomWithOpId  				= PAN/Token entry via electronic commerce with optional Identity Check-AAV or DSRP cryptogram in UCAF, codes send to mastercard: 81
// PANEntryCredentialOnFile 		 				= Credential on File, codes send to mastercard/visa: 10
type PANEntry string

// String returns the defined constants as string
// example: PANEntryUnknown => "unknown"
func (P PANEntry) String() string {
	return string(P)
}

const (
	PANEntryUnknown                  PANEntry = "unknown"
	PANEntryManual                   PANEntry = "panManualEntry"
	PANEntryAutoMagStripeNotRequired PANEntry = "panAutoEntryMagStripeNotRequired"
	PANEntryChip                     PANEntry = "chip"
	PANEntryContactless              PANEntry = "contactless"
	PANEntryCredentialOnFile         PANEntry = "credentialOnFile"
	// PANEntryViaEcomWithOpId unique case for mastercard, we haven't found a case for visa yet.
	// PAN/Token entry via electronic commerce with optional Identity Check-AAV or DSRP cryptogram in UCAF, value send to Mastercard: 81
	PANEntryViaEcomWithOpId PANEntry = "entryViaEcomWithOpId"
	PANEntryMagneticStrip   PANEntry = "magneticStrip"
)

var (
	panEntryMap = map[string]PANEntry{
		`unknown`:                          PANEntryUnknown,
		`panManualEntry`:                   PANEntryManual,
		`panAutoEntryMagStripeNotRequired`: PANEntryAutoMagStripeNotRequired,
		`chip`:                             PANEntryChip,
		`contactless`:                      PANEntryContactless,
		`credentialOnFile`:                 PANEntryCredentialOnFile,
		`entryViaEcomWithOpId`:             PANEntryViaEcomWithOpId,
		`magneticStrip`:                    PANEntryMagneticStrip,
	}
)

func IsValidPANEntry(panEntry string) bool {
	_, exists := panEntryMap[panEntry]
	return exists
}

// PINEntry Pin capture capability of terminal.
//
// The following fields will be mapped accordingly when send to mastercard and visa:
// default 						  = 0 unspecified
// PINEntryTerminalCanAcceptOnlinePin     = 1 Indicates terminal can accept and forward online PINs.
// PINEntryTerminalCantAcceptOnlinePin    = 2 Indicates terminal cannot accept and forward online PINs.
// PINEntryMPosSoftwarePinEntryCapability = 3 mPOS Software-based PIN Entry Capability
// PINEntryTerminalPinPadDown             = 8 Terminal PIN pad down.
type PINEntry string

// String returns the defined constants as string
// example: PINEntryUnspecified => "unspecified"
func (P PINEntry) String() string {
	return string(P)
}

const (
	PINEntryUnspecified                    PINEntry = "unspecified"
	PINEntryTerminalCanAcceptOnlinePin     PINEntry = "terminalCanAcceptOnlinePin"
	PINEntryTerminalCantAcceptOnlinePin    PINEntry = "terminalCantAcceptOnlinePin"
	PINEntryMPosSoftwarePinEntryCapability PINEntry = "mPosSoftwarePinEntryCapability"
	PINEntryTerminalPinPadDown             PINEntry = "terminalPinPadDown"
)

var (
	pinEntryMap = map[string]PINEntry{
		`unspecified`:                    PINEntryUnspecified,
		`terminalCanAcceptOnlinePin`:     PINEntryTerminalCanAcceptOnlinePin,
		`terminalCantAcceptOnlinePin`:    PINEntryTerminalCantAcceptOnlinePin,
		`mPosSoftwarePinEntryCapability`: PINEntryMPosSoftwarePinEntryCapability,
		`terminalPinPadDown`:             PINEntryTerminalPinPadDown,
	}
)

func IsValidPINEntry(pinEntry string) bool {
	_, exists := pinEntryMap[pinEntry]
	return exists
}

type AuthenticationData struct {
	ProgramProtocol              string // SF1 3DS version. Either 1 or 2
	DirectoryServerTransactionID string // SF2
}

// DE48 Additional Data—Private Use
type AdditionalRequestData struct {
	TransactionCategoryCode string
	LowRiskIndicator        string
	// SE20 Cardholder Verification Method
	CardholderVerificationMethod CardHolderVerificationMethod
	// SE42 Electronic Commerce Indicators
	OriginalEcommerceIndicator SLI
	// SE66 Authentication Data
	AuthenticationData AuthenticationData
	// SE80 PIN Service Code
	PinServiceCode string
}

// CardHolderVerificationMethod
//
// The following fields will be mapped accordingly when send to mastercard and visa:
// CardHolderVerificationMethodSignature       = visa: 1, mastercard: S
// CardHolderVerificationMethodOnlinePin       = visa: 2, mastercard: P
// CardHolderVerificationMethodUnattendedNoPin = visa: 3, mastercard: S
// CardHolderVerificationMethodMotoEcom        = visa: 4, mastercard: n/a, not possible
type CardHolderVerificationMethod string

const (
	CardHolderVerificationMethodSignature       CardHolderVerificationMethod = "signature"
	CardHolderVerificationMethodOnlinePin       CardHolderVerificationMethod = "onlinePin"
	CardHolderVerificationMethodUnattendedNoPin CardHolderVerificationMethod = "unattendedNoPin"
	CardHolderVerificationMethodMotoEcom        CardHolderVerificationMethod = "motoEcom"
)

var (
	cardHolderVerificationMethodMap = map[string]CardHolderVerificationMethod{
		`signature`:       CardHolderVerificationMethodSignature,
		`onlinePin`:       CardHolderVerificationMethodOnlinePin,
		`unattendedNoPin`: CardHolderVerificationMethodUnattendedNoPin,
		`motoEcom`:        CardHolderVerificationMethodMotoEcom,
	}
)

func IsValidCardHolderVerificationMethod(cardHolderVerificationMethod string) bool {
	_, exists := cardHolderVerificationMethodMap[cardHolderVerificationMethod]
	return exists
}

// DE48 Additional Data—Private Use
type AdditionalResponseData struct {
	// SE42 Electronic Commerce Indicators
	AppliedEcommerceIndicator *SLI
	// SE66 Authentication Data
	ReasonForUCAFDowngrade *int
}

func NewAppliedEcommerceIndicator(sli SLI) *SLI {
	if (sli == SLI{}) {
		return nil
	}

	return &sli
}

func ReasonForUCAFDowngradeFromString(r string) *int {
	if r == "" {
		return nil
	}

	i, err := strconv.Atoi(r)
	if err != nil {
		panic(err)
	}

	return &i
}

// DE61 Point-of-Service (POS) Data
type PointOfServiceData struct {
	TerminalAttendance                       int    // SF1 POS Terminal Attendance
	TerminalLocation                         int    // SF3 POS Terminal Location
	CardHolderPresence                       int    // SF4 POS Cardholder Presence
	CardPresence                             int    // SF5 POS Card Presence
	CardCaptureCapabilities                  int    // SF6 POS Card Capture Capabilities
	TransactionStatus                        int    // SF7 POS Transaction CardSchemeResponseStatus (not store in Spanner)
	TransactionSecurity                      int    // SF8 POS Transaction Security (not store in Spanner)
	CardHolderActivatedTerminalLevel         int    // SF10 Cardholder-Activated Terminal Level
	CardDataTerminalInputCapabilityIndicator int    // SF11 POS Card Data Terminal Input Capability Indicator
	AuthorizationLifeCycle                   string // SF12 POS Authorization Life Cycle (not store in Spanner)
	CountryCode                              string // SF13 POS Country Code (or Sub-CardAcceptor Information, if applicable) (not store in Spanner)
	PostalCode                               string // SF14 POS Postal Code (or Sub-CardAcceptor Information, if applicable) (not store in Spanner)
}

// AdditionalPOSInformation F060 Additional POS information visa-net-authorization-only-online-messages-technical-specifications.pdf page 338
type AdditionalPOSInformation struct {
	TerminalType                               string // SF1 length 1 byte
	TerminalEntryCapability                    string // SF2 length 1 byte
	ChipConditionCode                          string // SF3 length 1 byte
	ExistingDebtIndicator                      string // SF4 length 1 byte
	MerchantGroupIndicator                     string // SF5 length 2 byte
	ChipTransactionIndicator                   string // SF6 length 1 byte
	ChipCardAuthenticationReliabilityIndicator string // SF7 length 1 byte
	TypeOrLevelIndicator                       string // SF8 length 2 byte
	CardholderIDMethodIndicator                string // SF9 length 1 byte
	PartialAuthorizationIndicator              string // SF10 length 1 byte
	SpecialConditionIndicator                  string // Visa F60.4
	AdditionalAuthorizationIndicators          string // Visa F60.10
}

type PrivateUseFields struct {
	MerchantIdentifier                string
	CardholderCertificateSerialNumber string
	MerchantCertificateSerialNumber   string
	TransactionId                     string
	CavvData                          string
	Cvv2AuthorizationRequestData      string
	NotApplicable                     string
	ServiceIndicators                 string
	POSEnvironment                    string
}

type MastercardSchemeRequest struct {
	AuthorizationType        AuthorizationType
	PosPinCaptureCode        string // DE 26—Point-of-Service (POS) Personal ID Number (PIN) Capture Code
	AdditionalData           AdditionalRequestData
	PointOfServiceData       PointOfServiceData
	PosConditionCode         string // F025 Point-of-service Condition Code visa-net-authorization-only-online-messages-technical-specifications.pdf page 189
	AdditionalPOSInformation AdditionalPOSInformation
}

type VisaSchemeRequest struct {
	PosConditionCode         string
	AdditionalPOSInformation AdditionalPOSInformation
	PrivateUseFields         PrivateUseFields
}

type ResponseCode struct {
	Value       string
	Description string
}

type SLI struct {
	SecurityProtocol         int
	CardholderAuthentication int
	UCAFCollectionIndicator  int
}

func SLIFromString(sli string) SLI {
	if len(sli) != 3 {
		return SLI{}
	}

	securityProtocol, err := strconv.Atoi(sli[:1])
	if err != nil {
		panic("failed to convert security level indicator")
	}

	cardholderAuthentication, err := strconv.Atoi(sli[1:2])
	if err != nil {
		panic("failed to convert security level indicator")
	}

	ucafCollectionIndicator, err := strconv.Atoi(sli[2:])
	if err != nil {
		panic("failed to convert security level indicator")
	}

	return SLI{
		SecurityProtocol:         securityProtocol,
		CardholderAuthentication: cardholderAuthentication,
		UCAFCollectionIndicator:  ucafCollectionIndicator,
	}
}

type VisaSchemeResponse struct {
	TransactionId int
}

type MastercardSchemeResponse struct {
	AdditionalData         AdditionalResponseData
	AdditionalResponseData string
	TraceID                MTraceID
}

type CardSchemeResponse struct {
	Status                  AuthorizationStatus
	ResponseCode            ResponseCode
	AuthorizationIDResponse string
	EcommerceIndicator      int
	TraceId                 string
}

func ResponseDescriptionFromCode(code string) string {
	return ResponseCodeFromString(code).Description
}

func ResponseCodeFromString(code string) ResponseCode {
	switch code {
	case "00", "10":
		return ResponseCode{"approved", "Approved"}
	case "01":
		return ResponseCode{"issuer_declined", "Issuer Declined"}
	case "03":
		return ResponseCode{"invalid_merchant", "Invalid merchant"}
	case "04":
		return ResponseCode{"capture_card", "Capture card"}
	case "05":
		return ResponseCode{"do_not_honor", "Do not honor"}
	case "08":
		return ResponseCode{"honor_with_id", "Honor with ID"}
	case "12":
		return ResponseCode{"invalid_transaction", "Invalid transaction"}
	case "13":
		return ResponseCode{"Invalid amount", "Invalid amount"}
	case "14":
		return ResponseCode{"invalid_cardnumber", "Invalid card number"}
	case "15":
		return ResponseCode{"invalid_issuer", "Invalid issuer"}
	case "30":
		return ResponseCode{"format_error", "Format error"}
	case "41":
		return ResponseCode{"lost_card", "Lost card"}
	case "43":
		return ResponseCode{"stolen_card", "Stolen card"}
	case "51":
		return ResponseCode{"insufficient_funds", "Insufficient funds"}
	case "54":
		return ResponseCode{"expired_card", "Expired card"}
	case "55":
		return ResponseCode{"invalid_pin", "Invalid PIN"}
	case "57":
		return ResponseCode{"transaction_not_permitted_to_issuer", "Transaction not permitted to issuer/cardholder"}
	case "58":
		return ResponseCode{"transaction_not_permitted_to_terminal", "Transaction not permitted to acquirer/terminal"}
	case "61":
		return ResponseCode{"exceeds_limit", "Exceeds withdrawal amount limit"}
	case "62":
		return ResponseCode{"restricted_card", "Restricted card"}
	case "63":
		return ResponseCode{"security_violation", "Security violation"}
	case "65":
		return ResponseCode{"soft_decline", "Soft Decline"}
	case "70":
		return ResponseCode{"issuer_declined", "Contact Card Issuer"}
	case "71":
		return ResponseCode{"pin_not_changed", "PIN Not Changed"}
	case "75":
		return ResponseCode{"pin_tries_exceeded", "Allowable number of PIN tries exceeded"}
	case "76":
		return ResponseCode{"invalid_transaction", "Invalid/nonexistent “To Account” specified"}
	case "77":
		return ResponseCode{"invalid_transaction", "Invalid/nonexistent “From Account” specified"}
	case "78":
		return ResponseCode{"invalid_transaction", "Invalid/nonexistent account specified (general)"}
	case "81":
		return ResponseCode{"invalid_transaction", "Domestic Debit Transaction Not Allowed (Regional use only)"}
	case "84":
		return ResponseCode{"invalid_transaction", "Invalid Authorization Life Cycle"}
	case "85":
		return ResponseCode{"valid ", "Not declined Valid for all zero amount transactions."}
	case "86":
		return ResponseCode{"pin_not_possible", "PIN Validation not possible"}
	case "87":
		return ResponseCode{"no_cash_back_allowed ", "Purchase Amount Only, No Cash Back Allowed"}
	case "88":
		return ResponseCode{"cryptographic_failure", "Cryptographic failure"}
	case "89":
		return ResponseCode{"invalid_pin", "Unacceptable PIN-Transaction Declined-Retry"}
	case "91":
		return ResponseCode{"issuer_declined", "Authorization System or issuer system inoperative"}
	case "92":
		return ResponseCode{"invalid_issuer", "Unable to route transaction"}
	case "94":
		return ResponseCode{"duplicate", "Duplicate transmission detected"}
	case "96":
		return ResponseCode{"system_error", "System error"}
	default:
		return ResponseCode{"unknown_code", "received unknown code " + code}
	}
}
