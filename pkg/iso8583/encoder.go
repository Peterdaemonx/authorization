//nolint:gomnd
package iso8583

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"golang.org/x/text/encoding/charmap"
)

type format uint8

const (
	_ format = iota
	ASCII
	EBCDIC
)

type layout uint8

const (
	_ layout = iota
	NoRDW
	WithRDW
)

func NewEncoder(w io.Writer, opts ...Opt) *encoder {
	e := encoder{w: w}
	e.coder = applyOpts(e.coder, opts...)

	return &e
}

type encoder struct {
	coder
	w io.Writer
}

// EncodeIso8583 writes an iso8583-encoded message to its writer
//
//nolint:funlen,cyclop,gocognit,wrapcheck
func (e *encoder) EncodeIso8583(mti MTI, v interface{}) error {
	var err error

	definitions, err := StructDefinitions(v)
	if err != nil {
		return err
	}

	values := reflect.Indirect(reflect.ValueOf(v))

	// There will always be at least 1 bit map, the primary bit map
	bitmaps := make([]bitMap, 1)

	// Populate bit maps based on the presence of values
	for number, element := range definitions {
		// Determine the bit for this field number
		set := number / 64
		bit := number % 64

		if bit == 1 {
			return fmt.Errorf("field %d must not be defined; bit maps are handled internally", number)
		}

		if bit == 64 {
			return fmt.Errorf("field %d must not be defined; the message checksum is handled internally", number)
		}

		// Test if a value is provided for this field
		if values.Field(element.Field).IsZero() && !element.AutoFill {
			continue
		}

		// Create all the bit maps we need to flag this field as being present
		for len(bitmaps) <= set {
			// Set the first bit in the now currently last bit map, indicating there will be another bit map following
			bitmaps[len(bitmaps)-1].SetField(1)
			// Add the next bit map
			bitmaps = append(bitmaps, bitMap{}) //nolint:makezero
		}

		// Set the bit in the appropriate bit map indicating presence of this field
		bitmaps[set].SetField(bit)
	}

	w := &bytes.Buffer{}

	var transformer io.Writer
	if e.format == EBCDIC {
		transformer = &writer{w, charmap.CodePage1047.NewEncoder()}
	} else {
		transformer = w
	}

	// Write MTI
	if _, err := e.mtiFormat.encoder(w, transformer, "").Write(mti.Bytes()); err != nil {
		return fmt.Errorf("could not write MTI; %w", err)
	}

	// Write primary bit map
	// As this is binary data, it must not be transcoded
	if _, err := w.Write(bitmaps[0][:]); err != nil {
		return fmt.Errorf("could not write primary bit map; %w", err)
	}

	// Loop over bit maps and marshal all provided values
	// We loop over the bit maps. Not over `values` because of calculated fields like bit maps and checksums, and not
	// over `definitions` because it is unordered and has gaps
	for set, bitmap := range bitmaps {
		for bit := 1; bit <= 64; bit++ {
			// element number
			number := set*64 + bit

			if !bitmap.Field(bit) {
				// Field is not present
				continue
			}

			if bit == 1 {
				// The first field of a set is the bit map for the next set
				// As this is binary data, it must not be transcoded
				if _, err := w.Write(bitmaps[set+1].Bytes()); err != nil {
					return fmt.Errorf("could not write value for element %d; %w", number, err)
				}

				continue
			}

			if bit == 64 {
				// Field 64 is the message checksum.
				// It can only be present in the last set, as the last field of the message
				// TODO Create Message Checksum field?
				continue
			}

			element := definitions[number]
			value := values.Field(element.Field).Interface()

			var w2 io.Writer
			if element.Representation.Binary() {
				// Binary elements must not be transformed
				w2 = w
			} else {
				w2 = element.DataEncoding.encoder(transformer, transformer, "")
			}

			if err = Encode(value, element, w2); err != nil {
				return fmt.Errorf("could not encode; %w", NewElementError(number, err))
			}
		}
	}

	if e.layout == WithRDW {
		// Write the RDW with the length of the buffer before the actual data
		rdw := make([]byte, 4)
		binary.BigEndian.PutUint32(rdw, uint32(w.Len()))

		if _, err := e.w.Write(rdw); err != nil {
			return err
		}
	}

	// Flush the buffer to the actual output
	if _, err := io.Copy(e.w, w); err != nil {
		return err
	}

	return nil
}
