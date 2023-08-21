package visa

type Request struct {
	message *Message
}

func NewRequest(message *Message) Request {
	return Request{
		message: message,
	}
}

func (mcr Request) Packet() ([]byte, error) {
	payload, err := msg2Payload(mcr.message)
	if err != nil {
		return nil, err
	}

	packet, err := payload2Packet(mcr.message, payload)

	return packet, err
}
