//nolint:gomnd
package iso8583

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// Notes
// omitempty must not be used within slices. Each slice item must be of fixed length

// marshaler is the interface implemented by types that can marshal themselves.
type marshaler interface {
	MarshalIso8583() ([]byte, error)
}

// unmarshaler is the interface implemented by types that can unmarshal themselves.
type unmarshaler interface {
	UnmarshalIso8583([]byte) error
}

// Encode writes the iso8583 representation of the given field
//
//nolint:funlen
func Encode(v interface{}, element Definition, w io.Writer) error {
	data, err := Marshal(v, element)
	if err != nil {
		return err
	}

	// TlvTag must be written in the first position in the buffer
	if element.TlvTag != "" {
		tagBytes, err := WriteTlvTag(element.TlvTag)
		if err != nil {
			return err
		}

		_, err = w.Write(tagBytes)
		if err != nil {
			return fmt.Errorf("could not write value; %w", err)
		}
	}

	if data == nil && element.AutoFill {
		// If there is no value
		// Zero the field out
		var filler []byte

		switch {
		case element.Representation.Alphabetic():
			filler = []byte(` `)
		case element.Representation.Numeric():
			filler = []byte(`0`)
		default:
			filler = []byte(nil)
		}

		data = bytes.Repeat(filler, element.LengthMin)
	}

	// Justification of data
	if pad := element.LengthMin - len(data); pad > 0 {
		switch element.Justification {
		case JustifyLeft:
			// Add spaces on the right, so the value aligns left
			data = append(data, bytes.Repeat([]byte(` `), pad)...)
		case JustifyRight:
			// Add zero's on the left, so the value aligns right
			data = append(bytes.Repeat([]byte(`0`), pad), data...)
		}
	}

	// Assert data complies with element.Representation
	if err := element.Representation.Assert(data); err != nil {
		return fmt.Errorf("invalid value %q; %w", data, err)
	}

	// Assert length of value
	length := len(data)
	if length > element.Length {
		return fmt.Errorf("value length %d (%s) exceeds element length %d", length, string(data), element.Length)
	}

	if length < element.LengthMin {
		return fmt.Errorf("value length %d (%s) subceeds element minlength %d", length, string(data), element.LengthMin)
	}

	if element.LengthIndicator > 0 {
		// Data has a variable length indication field prepended
		li := element.LengthEncoding.encode(data, element.LengthIndicator)

		// Append length indication
		if _, err := w.Write(li); err != nil {
			return fmt.Errorf("could not write value length; %w", err)
		}
	}

	// BCD4 appends a 0 to the field. This needs to be appended to the data and not the length of the field
	// This is the reason for two writes.
	elemWr := element.DataEncoding.encoder(w, w, element.TlvTag)
	_, err = elemWr.Write(data)
	if err != nil {
		return fmt.Errorf("could not write value; %w", err)
	}

	return nil
}

// Decode is given the v to read element into from r
//
//nolint:wrapcheck
func Decode(v interface{}, element Definition, r io.Reader, lenEnc lengthEncoding) error {
	// Determine how much data must be read
	length := element.Length

	if element.LengthIndicator > 0 {
		// Value has a variable length indication field prepended
		li := make([]byte, element.LengthIndicator)
		if _, err := io.ReadFull(r, li); err != nil {
			return err
		}

		var err error

		length, err = element.LengthEncoding.decode(lenEnc.defDecode, li)
		if err != nil {
			return fmt.Errorf("invalid value in length field; %w", err)
		}
	}

	data := make([]byte, length)

	elemRdr := element.DataEncoding.decoder(r, r, element.TlvTag)

	// Read data
	if _, err := io.ReadFull(elemRdr, data); err != nil {
		return err
	}

	// Justification
	switch element.Justification {
	case JustifyLeft:
		data = bytes.TrimRight(data, ` `)
	case JustifyRight:
		data = bytes.TrimLeft(data, `0`)
	}

	if data == nil && element.AutoFill {
		// If there is no value
		// Zero the field out
		var filler []byte

		switch {
		case element.Representation.Alphabetic():
			filler = []byte(` `)
		case element.Representation.Numeric():
			filler = []byte(`0`)
		default:
			filler = []byte(nil)
		}

		data = bytes.Repeat(filler, element.LengthMin)
	}

	// Unmarshal the []byte's into the value
	// It is possible data is now empty
	return Unmarshal(v, data, element, lenEnc)
}

// Marshal returns the []bytes for the given variable
func Marshal(v interface{}, element Definition) ([]byte, error) {
	if marshaler, ok := v.(marshaler); ok {
		// this variable marshals itself
		return marshaler.MarshalIso8583() //nolint:wrapcheck
	}
	vv := reflect.Indirect(reflect.ValueOf(v))

	if !vv.IsValid() || vv.IsZero() {
		// If the field is empty, return nil to indicate there is no value
		return nil, nil
	}

	switch vv.Kind() {
	case reflect.String:
		return []byte(vv.String()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(vv.Int(), 10)), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(vv.Uint(), 10)), nil

	case reflect.Struct:
		var buf bytes.Buffer
		err := MarshalSubfields(v, &buf, element)

		return buf.Bytes(), err

	case reflect.Slice:
		var buf bytes.Buffer
		// Marshal each slice element
		for i := 0; i < vv.Len(); i++ {
			if err := MarshalSubfields(vv.Index(i).Interface(), &buf, element); err != nil {
				return nil, err
			}
		}

		return buf.Bytes(), nil

	default:
		return nil, fmt.Errorf("can not marhsal field of type %T", vv.Interface())
	}
}

// Unmarshal tries to load the []bytes in the given variable.
// For this function to have any effect, the variable must be a pointer.
//
//nolint:funlen,wrapcheck
func Unmarshal(v interface{}, data []byte, element Definition, lenEnc lengthEncoding) error {
	if unmarshaler, ok := v.(unmarshaler); ok {
		// This value can unmarshal itself
		return unmarshaler.UnmarshalIso8583(data)
	}

	if len(data) == 0 {
		// Nothing to do for empty fields.
		return nil
	}

	// Get the value pointed to
	vv := reflect.Indirect(reflect.ValueOf(v).Elem())

	// Determine the type of variable
	switch vv.Kind() {
	case reflect.String:
		vv.SetString(string(data))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(string(data), 10, 0)
		if err != nil {
			return err
		}

		vv.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(string(data), 10, 0)
		if err != nil {
			return err
		}

		vv.SetUint(i)

	case reflect.Struct:
		buf := bytes.NewBuffer(data)
		// Pass a pointer to the struct
		if err := UnmarshalSubfields(vv.Addr().Interface(), buf, element, lenEnc); err != nil {
			return err
		}

		if buf.Len() > 0 {
			return fmt.Errorf("%d bytes remain after unmarshaling struct: %x from %x", buf.Len(), buf.Bytes(), data)
		}

	case reflect.Slice:
		// vv is a slice of vt's
		vt := vv.Type().Elem()

		buf := bytes.NewBuffer(data)
		for buf.Len() > 0 {
			// Create a new variable of the correct type
			field := reflect.New(vt).Elem()
			value := field.Addr().Interface()

			// Replace a nil-pointer with a pointer to an actual variable
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.New(vt.Elem()))
				// Use the pointer
				value = field.Interface()
			}

			// Unmarshal into new variable
			if err := UnmarshalSubfields(value, buf, element, lenEnc); err != nil {
				return err
			}

			// Append new variable to the slice
			vv.Set(reflect.Append(vv, field))
		}

		if buf.Len() > 0 {
			return fmt.Errorf("%d bytes remain after unmarshaling slice", buf.Len())
		}

	default:
		return fmt.Errorf("can not unmarhsal field of type %T", vv.Interface())
	}

	return nil
}

type bitMap [8]byte

func (b *bitMap) Bytes() []byte {
	return b[:]
}

// SetField sets a field as present in the bit map
func (b *bitMap) SetField(field int) {
	// SendAuthorization field number (1 - 64) to a bit in a byte
	byteIndex := (field - 1) / 8
	bitIndex := (field - 1) % 8
	shift := uint(7 - bitIndex)
	// OR the bit to make it high
	b[byteIndex] |= 1 << shift
}

// Field returns whether the field is present in the bit map
func (b *bitMap) Field(field int) bool {
	// Map field number (1 - 64) to a bit in a byte
	byteIndex := (field - 1) / 8
	bitIndex := (field - 1) % 8
	shift := uint(7 - bitIndex)
	// AND the bit to check it is not low
	return b[byteIndex]&(1<<shift) != 0
}
