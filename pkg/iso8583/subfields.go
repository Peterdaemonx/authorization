package iso8583

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

// MarshalSubfields marshals all values for iso8583-fields into w
func MarshalSubfields(v interface{}, w io.Writer, parentElement Definition) error {
	// Read field definitions from struct
	definitions, err := StructDefinitions(v)
	if err != nil {
		return err
	}

	values := reflect.Indirect(reflect.ValueOf(v))

	// If an optional subfield is presented without data, but a succeeding subfield has data, the empty
	// optional subfields has to be provided anyhow. We therefor need to know the highest provided field
	// [IPM Clearing Formats 285]
	var numberMax int

	bitmap := bitMap{}

	for number, element := range definitions {
		if element.OmitEmpty && values.Field(element.Field).IsZero() {
			// ignore empty fields that may be omitted
			continue
		}

		if number > numberMax {
			numberMax = number
		}

		bitmap.SetField(number)
	}

	if parentElement.SubBitmap > 0 {
		if _, err := w.Write(bitmap.Bytes()[0:parentElement.SubBitmap]); err != nil {
			return fmt.Errorf("could not subfield bitmap element %d; %w", parentElement.Number, err)
		}
	}

	// Marshal each field
	for number := 1; number <= numberMax; number++ {
		element, ok := definitions[number]
		if !ok {
			// We do not have a definition for this field
			return NewElementError(number, fmt.Errorf("missing definition"))
		}

		if !bitmap.Field(number) && element.TlvTag != "" {
			continue
		}

		field := values.Field(element.Field)
		if !field.CanInterface() {
			// If we can not obtain the actual value, use the fields zero-value
			field = reflect.New(field.Type()).Elem()
		}

		if err = Encode(field.Interface(), element, w); err != nil {
			return NewElementError(number, err)
		}
	}

	return nil
}

// UnmarshalSubfields unmarshals all values from r into iso8583-fields in v
func UnmarshalSubfields(v interface{}, r io.Reader, parentElement Definition, lenEnc lengthEncoding) error {
	// Read field definitions from struct
	definitions, err := StructDefinitions(v)
	if err != nil {
		return err
	}

	// If a bitmap is present, use it. Otherwise, assume all fields will be present.
	bitmap := bitMap{}
	if parentElement.SubBitmap > 0 {
		// Subfield bitmaps are of dynamic length, so only read the required amount of bytes
		if _, err := io.ReadFull(r, bitmap[0:parentElement.SubBitmap]); err != nil {
			return fmt.Errorf("could not read subfield bit map; %w", err)
		}
	} else {
		bitmap = bitMap{0xFF, 0xFF, 0xFF}
	}

	// Find the highest known field number
	var numberMax int
	for number := range definitions {
		if number > numberMax {
			numberMax = number
		}
	}

	values := reflect.Indirect(reflect.ValueOf(v))
	// If subfield 3 is presented, so must 1 and 2. But 4 and higher don't have to be
	for number := 1; number <= numberMax; number++ {
		element, ok := definitions[number]
		if !ok {
			return NewElementError(number, fmt.Errorf("missing definition"))
		}

		if !bitmap.Field(number) {
			continue
		}

		// Get the variable we must decode into
		field := StructField(values, element.Field)

		// Pass a pointer to the variable where the data must be decoded into
		if err := Decode(field.Addr().Interface(), element, r, lenEnc); err != nil {
			if errors.Is(err, io.EOF) {
				// Not all subfields have to be present (ex PDS 158)
				break
			}

			return NewElementError(number, err)
		}
	}

	return nil
}
