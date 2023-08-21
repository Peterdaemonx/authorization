package iso8583

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

type tlvEncoder struct {
	w   io.Writer
	r   io.Reader
	tag string
}

func (e tlvEncoder) Write(bytes []byte) (n int, err error) {
	return e.w.Write(bytes)
}

func WriteTlvTag(tag string) (packed []byte, err error) {
	tagBytes := make([]byte, hex.DecodedLen(len(tag)))

	_, err = hex.Decode(tagBytes, []byte(tag))
	if err != nil {
		return nil, fmt.Errorf("failed to convert subfield Tag %s to int", tag)
	}
	packed = append(packed, tagBytes...)

	return packed, nil
}

func (e tlvEncoder) Read(value []byte) (n int, err error) {
	err = binary.Read(e.r, binary.BigEndian, e.tag)
	if err != nil {
		return 0, fmt.Errorf("error reading TLV tag %s: %w", e.tag, err)
	}

	err = binary.Read(e.r, binary.BigEndian, value)
	if err != nil {
		return 0, fmt.Errorf("error reading TLV value %s: %w", value, err)
	}

	out := make([]byte, 0)
	length, _ := io.ReadFull(e.r, out)

	return length, nil
}
