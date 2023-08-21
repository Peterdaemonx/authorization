package visa

import (
	"fmt"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/visa/base1"
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
	SourceStationID string
	Mti             iso8583.MTI
	Fields          base1.Fields
}

// Encode writes the message MTI and Data Elements, including the Private Data Subelements, to an encoder
func (m *Message) Encode(encoder Iso8583Encoder) error {
	return encoder.EncodeIso8583(m.Mti, m.Fields)
}

// Decode reads a message from a decoder
func (m *Message) Decode(decoder Iso8583Decoder) error {
	return decoder.DecodeIso8583(&m.Mti, &m.Fields)
}

func Echo(sourceStationId string) *Message {
	hours, minutes, _ := time.Now().Clock()
	return &Message{
		SourceStationID: sourceStationId,
		Mti:             iso8583.NewMti(networkManagementRequestMTI),
		Fields: base1.Fields{
			F007_TransmissionDateTime:             base1.F007FromTime(time.Now()),
			F011_SystemTraceAuditNumber:           fmt.Sprintf("00%02d%02d", hours, minutes),
			F063_NetworkData:                      base1.F063_NetworkData{SF1_NetworkID: "0002"},
			F070_NetworkManagementInformationCode: "301",
		},
	}
}
