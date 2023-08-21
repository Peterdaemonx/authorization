// nolint: unused
package iso8583

// MTI (Message type indicator) is a four-digit numeric field which indicates the overall function of a message
type (
	MTI             [4]byte
	isoVersion      uint8
	messageClass    uint16
	messageFunction uint16
	messageOrigin   uint8
)

// Bit flags
const (
	ReservedVersion isoVersion = 1 << iota
	Iso8583v1987               // ISO8583:1987
	Iso8583v1993               // ISO8583:1993
	Iso8583v2003               // ISO8583:2003
	NationalUse                // National use
	PrivateUse                 // Private use
)

// nolint:lll
const (
	ReservedClass      messageClass = 1 << iota
	Authorization                   // Determine if funds are available, get an approval but do not post to account for reconciliation
	Financial                       // Determine if funds are available, get an approval and post directly to the account
	FileActions                     // Used for hot-card, TMS and other exchanges
	ReversalChargeback              // Reverses the action of a previous authorization, or charges back a previously cleared financial message
	Reconciliation                  // Settlement information message
	Administrative                  // Administrative advice
	FeeCollection                   //
	NetworkManagement               // Used for secure key exchange, logon, echo test and other network functions
)

//nolint:lll
const (
	ReservedFunction            messageFunction = 1 << iota
	Request                                     // Request from acquirer to issuer to carry out an action; issuer may accept or reject
	RequestResponse                             // Issuer response to a request
	Advice                                      // Advice that an action has taken place; receiver can only accept, not reject
	AdviceResponse                              // Response to an advice
	Notification                                // Notification that an event has taken place; receiver can only accept, not reject
	NotificationAcknowledgement                 // Response to a notification
	Instruction                                 //
	InstructionAcknowledgement                  //
)

const (
	ReservedOriginator messageOrigin = 1 << iota
	Acquirer
	AcquirerRepeat
	Issuer
	IssuerRepeat
	Other
	OtherRepeat
)

func NewMti(mti string) MTI {
	_ = mti[3] // bounds check
	return MTI{mti[0], mti[1], mti[2], mti[3]}
}

// String returns the MTI as a string
func (mti MTI) String() string {
	return string(mti[:])
}

// Bytes returns the MTI as a byte-slice
func (mti MTI) Bytes() []byte {
	return mti[:]
}

// IsoVersion returns the ISO 8583 version as indicated by the first digit
func (mti MTI) isoVersion() isoVersion {
	switch mti[0] {
	case '0':
		return Iso8583v1987
	case '1':
		return Iso8583v1993
	case '2':
		return Iso8583v2003
	case '8':
		return NationalUse
	case '9':
		return PrivateUse
	default:
		return ReservedVersion
	}
}

// MessageClass returns the overall purpose of the message as indicated by the second digit
func (mti MTI) messageClass() messageClass {
	switch mti[1] {
	case '1':
		return Authorization
	case '2':
		return Financial
	case '3':
		return FileActions
	case '4':
		return ReversalChargeback
	case '5':
		return Reconciliation
	case '6':
		return Administrative
	case '7':
		return FeeCollection
	case '8':
		return NetworkManagement
	default:
		return ReservedClass
	}
}

// MessageFunction returns the message function, which defines how the message
// should flow within the system, as indicated by the third digit
func (mti MTI) messageFunction() messageFunction {
	switch mti[2] {
	case '0':
		return Request
	case '1':
		return RequestResponse
	case '2':
		return Advice
	case '3':
		return AdviceResponse
	case '4':
		return Notification
	case '5':
		return NotificationAcknowledgement
	case '6':
		return Instruction
	case '7':
		return InstructionAcknowledgement
	default:
		return ReservedFunction
	}
}

// Originator return the location of the message source as indicated by the fourth digit
func (mti MTI) messageOrigin() messageOrigin {
	switch mti[3] {
	case '0':
		return Acquirer
	case '1':
		return AcquirerRepeat
	case '2':
		return Issuer
	case '3':
		return IssuerRepeat
	case '4':
		return Other
	case '5':
		return OtherRepeat
	default:
		return ReservedOriginator
	}
}

func (v isoVersion) String() string {
	return map[isoVersion]string{
		ReservedVersion: "Reserved by ISO",
		Iso8583v1987:    "ISO8583:1987",
		Iso8583v1993:    "ISO8583:1993",
		Iso8583v2003:    "ISO8583:2003",
		NationalUse:     "National use",
		PrivateUse:      "Private use",
	}[v]
}

func (v messageClass) String() string {
	return map[messageClass]string{
		ReservedClass:      "Reserved by ISO",
		Authorization:      "Authorization",
		Financial:          "Financial",
		FileActions:        "File actions",
		ReversalChargeback: "Reversal or chargeback",
		Reconciliation:     "Reconciliation",
		Administrative:     "Administrative ",
		FeeCollection:      "Fee collection",
		NetworkManagement:  "Network management",
	}[v]
}

func (v messageFunction) String() string {
	return map[messageFunction]string{
		ReservedFunction:            "Reserved for ISO use",
		Request:                     "Request",
		RequestResponse:             "Request response",
		Advice:                      "Advice",
		AdviceResponse:              "Advice response",
		Notification:                "Notification",
		NotificationAcknowledgement: "Notification acknowledgement",
		Instruction:                 "Instruction",
		InstructionAcknowledgement:  "Instruction acknowledgement",
	}[v]
}

func (v messageOrigin) String() string {
	return map[messageOrigin]string{
		ReservedOriginator: "Reserved by ISO",
		Acquirer:           "Acquirer",
		AcquirerRepeat:     "Acquirer repeat",
		Issuer:             "Issuer",
		IssuerRepeat:       "Issuer repeat",
		Other:              "Other",
		OtherRepeat:        "Other repeat",
	}[v]
}
