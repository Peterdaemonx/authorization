package visa

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/yerden/go-util/bcd"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/connection"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/visa/base1"
)

const (
	networkManagementRequestMTI  = "0800"
	networkManagementResponseMTI = "0810"
)

const (
	rejectionHeaderLength = 26
	vmlhReserved          = 2
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
	return mcr.message.Fields.F011_SystemTraceAuditNumber == mcReq.message.Fields.F011_SystemTraceAuditNumber
}

func (mcr Response) PacketToSend() ([]byte, error) {
	if mcr.message.Mti.String() == networkManagementRequestMTI {
		payload, err := getNetworkManagementResponse(&mcr.message)
		return payload, err
	}

	log.Printf("Unexpected MTI: %s", mcr.message.Mti.String())
	log.Printf("Message: %#v", mcr.message.Fields)

	return nil, ErrNotExpectedMessage
}

func (mcr Response) Message() Message {
	return mcr.message
}

func (mcr Response) Error() error {
	return mcr.error
}

func decode(payload []byte) (Message, error) {
	var msg Message

	header, message := splitHeaderAndMessage(payload)

	if isRejectionHeader(header) {
		rejectionCode, err := decodeRejectionCode(header)
		if err != nil {
			return msg, fmt.Errorf("visa/decode(): %w", err)
		}

		return msg, fmt.Errorf("decode(): rejection code: %s", rejectionCode)
	}

	decoder := iso8583.NewDecoder(bytes.NewReader(message), iso8583.FormatAscii, iso8583.NoRdwLayout, iso8583.MtiBcd4, iso8583.HexLen)

	err := msg.Decode(decoder)
	return msg, err
}

func getNetworkManagementResponse(reqMsg *Message) ([]byte, error) {
	resMsg := &Message{
		Mti: iso8583.NewMti(networkManagementResponseMTI),
		Fields: base1.Fields{
			F007_TransmissionDateTime:             reqMsg.Fields.F007_TransmissionDateTime,
			F011_SystemTraceAuditNumber:           reqMsg.Fields.F011_SystemTraceAuditNumber,
			F070_NetworkManagementInformationCode: reqMsg.Fields.F070_NetworkManagementInformationCode,
		},
	}

	payload, err := msg2Payload(resMsg)
	if err != nil {
		return nil, err
	}

	packet, err := payload2Packet(resMsg, payload)

	return packet, err
}

func splitHeaderAndMessage(buffer []byte) (headerBuffer []byte, messageBuffer []byte) {
	// buff[0] and buff[1] make VMLH (Visa Message Length Header) - therefore we skip those
	primaryHeaderLength := int(buffer[vmlhReserved])
	secondaryHeaderLength := 0
	if primaryHeaderLength >= rejectionHeaderLength {
		secondaryHeaderLength = int(buffer[vmlhReserved+primaryHeaderLength])
	}
	splitOffset := primaryHeaderLength + secondaryHeaderLength
	return buffer[vmlhReserved:splitOffset], buffer[vmlhReserved+splitOffset:]
}

func isRejectionHeader(headerBuff []byte) bool {
	return headerBuff[0] >= rejectionHeaderLength && headerBuff[22]&0x80 != 0
}

func decodeRejectionCode(headerBuff []byte) (string, error) {

	decoded := make([]byte, 4)
	decoder := bcd.NewDecoder(bcd.Standard)
	_, err := decoder.Decode(decoded, headerBuff[24:26])
	if err != nil {
		return "", fmt.Errorf("decodeRejectionCode(): %w", err)
	}
	return leftPad2Len(string(decoded), "0", 4), nil
}

func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}
