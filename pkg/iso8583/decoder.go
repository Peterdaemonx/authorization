//nolint:gomnd
package iso8583

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"

	"golang.org/x/text/encoding/charmap"
)

type decoder struct {
	r io.Reader

	coder
}

func NewDecoder(r io.Reader, opts ...Opt) *decoder {
	d := decoder{
		r: r,
	}

	for _, opt := range opts {
		d.coder = opt(d.coder)
	}

	return &d
}

// DecodeIso8583 reads the next iso8583-encoded message from the reader
// The MTI and values are assigned to the passed variables
//
//nolint:nestif,funlen,gocognit,cyclop,wrapcheck
func (d *decoder) DecodeIso8583(mti *MTI, v interface{}) error {
	var err error

	definitions, err := StructDefinitions(v)
	if err != nil {
		return err
	}

	var r io.Reader

	if d.layout == WithRDW {
		// Read the 4-byte Record Descriptor Word
		rdw := make([]byte, 4)
		if _, err := io.ReadFull(d.r, rdw); err != nil {
			if errors.Is(err, io.EOF) {
				// This is the only moment we expect to get an EOF sooner or later
				return err //nolint:wrapcheck
			}

			return fmt.Errorf("could not read RDW; %w", err)
		}

		length := binary.BigEndian.Uint32(rdw)

		// Read the entire record into a buffer
		buf := &bytes.Buffer{}
		if _, err := io.CopyN(buf, d.r, int64(length)); err != nil {
			return fmt.Errorf("could not read message; %w", err)
		}
		// Use the buffer for further reading. EOF means end of the block¾¾
		r = buf
	} else {
		// Read data directly from the source
		r = d.r
	}

	var transformer io.Reader
	if d.format == EBCDIC {
		transformer = &reader{r, charmap.CodePage1047.NewDecoder()}
	} else {
		transformer = r
	}

	// The first 4 bytes are the MTI
	if _, err := io.ReadFull(d.mtiFormat.decoder(r, transformer, ""), mti[:]); err != nil {
		if errors.Is(err, io.EOF) && d.layout != WithRDW {
			// This is the only moment we expect to get an EOF sooner or later
			return err
		}
		// We're done with this file
		if errors.Is(err, io.EOF) {
			return err
		}

		return fmt.Errorf("could not read MTI; %w", err)
	}

	bitmaps := make([]bitMap, 1)

	// The next 8-bytes are the primary bit map
	// As this is binary data, it must not be transcoded
	if _, err := io.ReadFull(r, bitmaps[0][:]); err != nil {
		return fmt.Errorf("could not read primary bit map; %w", err)
	}

	values := reflect.Indirect(reflect.ValueOf(v))

	// Go over each bit and unmarshal the fields that are present
	// New bit maps are added as we find them
	for set := 0; set < len(bitmaps); set++ {
		// Loop over bits of current bit map, looking for field presence
		bitmap := bitmaps[set]

		for bit := 1; bit <= 64; bit++ {
			if !bitmap.Field(bit) {
				// Field is not present
				continue
			}

			// element number
			number := set*64 + bit

			if bit == 1 {
				// Field 1 of any set is the bit map for the next set
				// As this is binary data, it must not be transcoded
				var bitmap bitMap
				if _, err = io.ReadFull(r, bitmap[:]); err != nil {
					return fmt.Errorf("could not read value for element %d; %w", number, err)
				}

				// Append bit map to the list of bit maps to be processed
				//nolint:makezero
				bitmaps = append(bitmaps, bitmap)

				continue
			}

			if bit == 64 {
				// Field 64 in a set is the message checksum. It can only be present in the last set,
				// as the last field of the message
				// TODO FetchConfig Message Checksum field
				continue
			}

			// Obtain definition for this field so we know how to read it
			element, ok := definitions[number]
			if !ok {
				// We do not have the definition for this field
				return fmt.Errorf("missing definition for element %d", number)
			}

			// Get the variable we must decode into
			field := StructField(values, element.Field)

			var r2 io.Reader
			if element.Representation.Binary() {
				// Binary elements must not be transcoded
				r2 = r
			} else {
				r2 = element.DataEncoding.decoder(transformer, transformer, element.TlvTag)
			}

			// Pass a pointer to the variable where the data must be decoded into
			if err := Decode(field.Addr().Interface(), element, r2, d.lenFormat); err != nil {
				return NewElementError(number, err)
			}
		}
	}

	return nil
}

// StructField prepares a struct field to be decoded into.
// If the field is not initialized it will be. If the field is a pointer it is dereferenced
func StructField(values reflect.Value, idx int) reflect.Value {
	field := values.Field(idx)

	if !field.CanInterface() {
		// If we can't use the real value, create a fake one
		field = reflect.New(field.Type()).Elem()
	}

	if field.Kind() == reflect.Ptr {
		// Replace a nil-pointer with a pointer to an actual variable
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		// Use the actual variable instead of the pointer
		field = field.Elem()
	}

	return field
}
