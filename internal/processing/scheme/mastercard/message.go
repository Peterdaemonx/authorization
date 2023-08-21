package mastercard

import (
	"fmt"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/mastercard/cis"
)

type Iso8583Decoder interface {
	// DecodeIso8583 stores the data elements for the next message in its input
	// in the value pointed to by v and returns the MTI
	DecodeIso8583(*iso8583.MTI, interface{}) error
}

type Iso8583Encoder interface {
	// EncodeIso8583 writes the MTI and data elements in v to its output
	EncodeIso8583(iso8583.MTI, interface{}) error
}

type Message struct {
	Mti          iso8583.MTI
	DataElements cis.DataElements
}

// Encode writes the message MTI and Data Elements, including the Private Data Subelements, to an encoder
func (m *Message) Encode(encoder Iso8583Encoder) error {
	return encoder.EncodeIso8583(m.Mti, m.DataElements)
}

// Decode reads a message from a decoder
func (m *Message) Decode(decoder Iso8583Decoder) error {
	return decoder.DecodeIso8583(&m.Mti, &m.DataElements)
}

func Echo() *Message {
	hours, minutes, _ := time.Now().Clock()
	return &Message{
		Mti: iso8583.NewMti(networkManagementRequestMTI),
		DataElements: cis.DataElements{
			DE2_PrimaryAccountNumber:              "20706",
			DE7_TransmissionDateTime:              cis.DE7FromTime(time.Now()),
			DE11_SystemTraceAuditNumber:           fmt.Sprintf("00%02d%02d", hours, minutes),
			DE33_ForwardingInstitutionIDCode:      "200706",
			DE70_NetworkManagementInformationCode: "270",
		}}
}
