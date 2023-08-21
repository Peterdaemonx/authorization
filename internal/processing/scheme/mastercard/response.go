package mastercard

import (
	"bytes"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/connection"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/mastercard/cis"
)

const (
	networkManagementRequestMTI  = "0800"
	networkManagementResponseMTI = "0810"
)

type Response struct {
	packet  []byte
	message Message
	error   error
}

func NewResponse(packet []byte, error error) connection.Received {
	var msg Message
	if error == nil && packet != nil {
		msg, error = decode(packet[2:])
	}
	return Response{
		packet:  packet,
		message: msg,
		error:   error,
	}
}

func (mcr Response) IsRequestResponse(request connection.Request) bool {
	mcReq, ok := request.(Request)
	if !ok {
		return false
	}
	return mcr.message.DataElements.DE11_SystemTraceAuditNumber == mcReq.message.DataElements.DE11_SystemTraceAuditNumber
}

func (mcr Response) PacketToSend() ([]byte, error) {
	if mcr.message.Mti.String() == networkManagementRequestMTI {
		payload, err := getNetworkManagementResponse(mcr.message)
		return payload, err
	}

	return nil, ErrNotExpectedMessage
}

func (mcr Response) Message() Message {
	return mcr.message
}

func (mcr Response) Error() error {
	return mcr.error
}

func decode(payload []byte) (Message, error) {
	const (
		asciiZero  = byte(48)
		ebcdicZero = byte(240)
	)

	mtiZero := payload[0]

	var decoder Iso8583Decoder
	switch mtiZero {
	case asciiZero:
		decoder = iso8583.NewDecoder(bytes.NewReader(payload), iso8583.FormatAscii, iso8583.NoRdwLayout)
	case ebcdicZero:
		decoder = iso8583.NewDecoder(bytes.NewReader(payload), iso8583.FormatEbcdic, iso8583.NoRdwLayout)
	default:
		return Message{}, ErrUnknownEncoding
	}

	var msg Message
	err := msg.Decode(decoder)
	return msg, err
}

func getNetworkManagementResponse(reqMsg Message) ([]byte, error) {
	resMsg := Message{
		Mti: iso8583.NewMti(networkManagementResponseMTI),
		DataElements: cis.DataElements{
			DE2_PrimaryAccountNumber:              reqMsg.DataElements.DE2_PrimaryAccountNumber,
			DE7_TransmissionDateTime:              reqMsg.DataElements.DE7_TransmissionDateTime,
			DE11_SystemTraceAuditNumber:           reqMsg.DataElements.DE11_SystemTraceAuditNumber,
			DE33_ForwardingInstitutionIDCode:      reqMsg.DataElements.DE33_ForwardingInstitutionIDCode,
			DE39_ResponseCode:                     "00",
			DE63_NetworkData:                      reqMsg.DataElements.DE63_NetworkData,
			DE70_NetworkManagementInformationCode: reqMsg.DataElements.DE70_NetworkManagementInformationCode,
		},
	}

	payload, err := msg2Payload(&resMsg)
	if err != nil {
		return nil, err
	}

	packet, err := payload2Packet(payload)

	return packet, err
}
