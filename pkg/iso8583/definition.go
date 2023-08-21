package iso8583

import (
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

// Definition holds the definition for an element
type Definition struct {
	// Field index in struct
	Field int

	// Number is the ISO8583 element number
	Number int

	// LengthIndicator is the size of the length indicator (common are 0, 2 and 3)
	LengthIndicator int

	// LengthEncoding
	LengthEncoding lengthEncoding
	DataEncoding   dataEncoding
	Bcd4Len        int

	// Length is the maximum length of the value
	Length int

	// LengthMin is the minimum length of the value
	LengthMin int

	// Representation is the data representation of value
	Representation representation

	// Justification is the justification of the value
	Justification justification

	// OmitEmpty indicates whether a subfield may be omitted when encoding an empty value
	OmitEmpty bool

	// AutoFill indicates whether a subfield may be autofilled with spaces or zeros when a value is omitted
	AutoFill bool

	// DisableAutofill indicates in which mtis the autofill feature is not going to be enabled
	DisableAutoFill []string

	// SubBitmap indicates if the element starts with a bitmask. Only used for structs/subfields
	SubBitmap int

	// TlvTag indicates the tlv tag number in hexadecimal required for write data in tlv format
	TlvTag string
}

type justification uint8

const (
	_ justification = iota
	JustifyLeft
	JustifyRight
)

type lengthEncoding uint8

type decodeLen func([]byte) (int, error)

const (
	LengthEncodingDefault lengthEncoding = iota
	LengthEncodingAscii
	LengthEncodingHex
	LengthEncodingBin
	LengthEncodingHexBit4
)

func (le lengthEncoding) defDecode(d []byte) (int, error) {
	return le.decode(func(bytes []byte) (int, error) {
		return strconv.Atoi(string(d))
	}, d)
}

func (le lengthEncoding) decode(def decodeLen, d []byte) (int, error) {
	switch le {
	case LengthEncodingDefault:
		return def(d)
	case LengthEncodingHex:
		ascii, err := hex.DecodeString(string(d))
		if err != nil {
			return 0, err
		}

		return int(ascii[0]), nil
	case LengthEncodingAscii:
		return strconv.Atoi(string(d))
	case LengthEncodingBin:
		i := int(d[0])
		return i, nil
	case LengthEncodingHexBit4:
		return len(le.encodeHexBit4(len(d), len(d))), nil
	default:
		return def(d)
	}
}

func (le lengthEncoding) encode(d []byte, i int) []byte {
	switch le {
	default:
		return []byte(fmt.Sprintf("%0*d", i, len(d)))
	case LengthEncodingHex:
		h := hex.EncodeToString([]byte{byte(len(d))})
		return []byte(h)
	case LengthEncodingBin:
		return []byte{byte(len(d))}
	case LengthEncodingHexBit4:
		return le.encodeHexBit4(i, len(d))
	}
}

func (le lengthEncoding) encodeHexBit4(maxLen, dataLen int) []byte {

	res := []byte{uint8(dataLen)}

	// Handle LLLL encode - turn eg 5 to 0005.
	if maxLen == 4 {
		res = append([]uint8{0}, res...)
	}

	return res
}

type dataEncoding uint8

const (
	DataEncodingDefault dataEncoding = iota
	DataEncodingAscii
	DataEncodingEbcdic
	DataEncodingBcd4
	DataEncodingTlv
)

// encoder returns the proper encoder for the given encoding
// The encoder will wrap w; if the configured encoding is Default, d will
// be returned.
// If w is already of the needed implementation (using a v.(type) assertion), it is not wrapped.
func (de dataEncoding) encoder(w, d io.Writer, tag string) io.Writer {
	switch de {
	case DataEncodingEbcdic:
		if _, ok := w.(*writer); ok {
			return w
		}
		return &writer{w, charmap.CodePage1047.NewEncoder()}
	case DataEncodingAscii:
		return w
	case DataEncodingBcd4:
		if _, ok := w.(bcdEncoder); ok {
			return w
		}
		return bcdEncoder{w: w}
	case DataEncodingTlv:
		if _, ok := w.(tlvEncoder); ok {
			return w
		}
		return tlvEncoder{w: w, tag: tag}
	case DataEncodingDefault:
		return d
	default:
		return d
	}
}

// decoder returns the proper decoder for the given encoding
// The decoder will wrap r; if the configured encoding is Default, d will
// be returned.
// If r is already of the needed implementation (using a v.(type) assertion), it is not wrapped.
func (de dataEncoding) decoder(r, d io.Reader, tlvTag string) io.Reader {
	switch de {
	case DataEncodingEbcdic:
		if _, ok := r.(*reader); ok {
			return r
		}
		return &reader{r, charmap.CodePage1047.NewDecoder()}
	case DataEncodingAscii:
		return r
	case DataEncodingBcd4:
		if _, ok := r.(bcdEncoder); ok {
			return r
		}
		return bcdEncoder{r: r}
	case DataEncodingTlv:
		if _, ok := r.(tlvEncoder); ok {
			return r
		}
		return tlvEncoder{r: r, tag: tlvTag}
	case DataEncodingDefault:
		return d
	default:
		return d
	}
}

// StructDefinitions uses reflection to read the element definitions from the tags of struct fields
// The definitions are returned indexed by element number
//
// The specification of a message is implemented using a struct with tagged fields.
// The name of the field is irrelevant and can be freely chosen. The field type and the iso8583-tag determine
// all the behaviour.
//
// FIELD TYPE
// The field type can be;
// - string
// - int, int8, int16, int32 or int64,
// - uint, uint8, uint16, uint32 or uint64,
// - a custom type implementing iso8583.marshaler and/or iso8583.unmarshaler
// - a slice of a custom type
// - a pointer to one of the above
//
// Use the type that is most convenient to use in your code.
// Types that have their default value are omitted when encoding a message. Therefor, if you need a numeric value
// of "0" (zero) do not use `int`, but `*int` or `string`, so the value differs from the default.
// When decoding, fields that are not present keep their default value. So again, do not use `int` if you want to
// differentiate between "0" and not being present.
//
// Slices can only be used with custom types because such types must implement how each slice element is
// encoded/decoded. If a `[]string` would be used, it is impossible to know how many characters make up the
// first element and how many the next.
//
// FIELD TAG
// The field tag can consist of several comma separated attributes. Each attribute is a key(=value).
// The first attribute must always specify the field number and type.
//
// The number matches the field to the bitmaps.
// The type includes both the data type and maximum length of data. Between the type and length is either one
// dash or one or multiple dots. A dash means the data has a fixed length while dots mean variable length.
// The number of dots describe the length of the variable length field. One dot is an LVAR, two dots an LLVAR,
// three an LLLVAR, etc. With mastercard you should only need LLVAR and LLLVAR; 2 and 3 dots.
//
// DATA TYPE
// The following types are supported;
//
//	a    alphabetic characters A–Z and a–z
//	n    numeric digits 0–9
//	as   alphabetic characters (A–Z and a–z), and space character
//	ns   numeric digits 0–9 and special characters (including space)
//	an   alphabetic (A–Z and a–z) and numeric characters
//	ans  alphabetic (A–Z and a–z), numeric, and special characters (including space)
//	b    binary representation of data in eight-bit bytes
//
// FIELD ATTRIBUTES
// Other attributes that can be added in the field tag are;
//
//		minlength=[0-9+]
//		    The minimum length of the field value. The library shall return an error when trying to encode a value
//		    shorter than the minlength. It has no effect on decoding.
//
//		justify=[left|right]
//		    When encoding justifies the value up to the minlength. If there is no minlength then up to the field length.
//		    When decoding padding characters are removed.
//
//		omitempty
//		    All subfields are always encoded, even when they hold their default value. This option can be used to omit
//		    such fields. Must only be used for the last field(s).
//
//		autofill
//		    When a field has it's default value, instead of omitting it fill it with zeros or spaces (depending on
//		    the field type) up to minlength. If there is no minlength then up to the field length.
//
//	 lenenc=[ascii/hex/bin]
//	     Override how the length property of a specific field is encoded.
//
//	 dataenc=[ascii/ebcdic/bcd4]
//	     Override the data encoding of the whole message for a specific field.
//
//	 bcd4len=[1-9+]
//	     The amount of data that was bcd4 encoded.
//
//	 subbitmap=[1-8]
//	     The length of the subfield bitmap; no value means no bitmap at all
//
//nolint:funlen,gocognit,cyclop
func StructDefinitions(v interface{}) (map[int]Definition, error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("value must be (a pointer to) a struct, got %T", v)
	}

	// TODO cache the definitions by indexing on t

	definitions := make(map[int]Definition, t.NumField())

	for idx := 0; idx < t.NumField(); idx++ {
		field := t.Field(idx)

		tag, ok := field.Tag.Lookup("iso8583")
		if !ok {
			// Ignore struct fields without iso8583-tag
			continue
		}

		// Element number, format and optional attributes
		number, format, attributes := parseTag(tag)
		if number == 0 {
			return nil, fmt.Errorf("struct field %q has no number in iso8583-tag: %s", field.Name, tag)
		}

		if format == "" {
			return nil, fmt.Errorf("struct field %q has no format in iso8583-tag: %s", field.Name, tag)
		}

		element := Definition{
			Field:  idx,
			Number: number,
		}

		// Parse the format into the data representation, variable length indication field, and field length
		if i := strings.Index(format, "-"); i != -1 {
			// Representation and length separated by a dash is a fixed length field
			element.Representation = notation[format[:i]]
			format = format[i+1:]
		} else if i := strings.Index(format, "."); i != -1 {
			// Representation and length separated by one or more periods is a variable length field
			element.Representation = notation[format[:i]]
			for format[i] == '.' {
				element.LengthIndicator++
				i++
			}
			format = format[i:]
		}

		if element.Representation == 0 {
			return nil, fmt.Errorf("struct field %q has invalid representation in iso8583-tag %q", field.Name, tag)
		}

		// The remaining format must be the field length
		var err error
		if element.Length, err = strconv.Atoi(format); err != nil {
			return nil, fmt.Errorf(
				"struct field %q has invalid length in iso8583-tag %q: %w",
				field.Name, tag, err,
			)
		}

		// Determine the minimum length of the value
		attr, ok := attributes["minlength"]
		_, tlvOk := attributes["tlvTag"]
		switch {
		case tlvOk:
			// for fields set as tlv we don't need it to modify LengthMin as 1
		case ok:
			if element.LengthMin, err = strconv.Atoi(attr); err != nil {
				return nil, fmt.Errorf(
					"struct field %q has invalid minlength-attribute %q in iso8583-tag: %w",
					field.Name, attr, err,
				)
			}
		case element.LengthIndicator == 0:
			// The field has a fixed length
			element.LengthMin = element.Length
		default:
			// The field has a variable length
			element.LengthMin = 1
		}

		if attr, ok := attributes["justify"]; ok {
			switch attr {
			case "left":
				element.Justification = JustifyLeft
			case "right":
				element.Justification = JustifyRight
			default:
				return nil, fmt.Errorf("struct field %q has invalid justify-attribute %q in iso8583-tag", field.Name, attr)
			}
		}

		if _, ok := attributes["omitempty"]; ok {
			element.OmitEmpty = true
		}

		if _, ok := attributes["autofill"]; ok {
			element.AutoFill = true
		}

		if _, ok := attributes["subbitmap"]; ok {
			element.SubBitmap, err = strconv.Atoi(attributes["subbitmap"])
			if err != nil {
				return nil, fmt.Errorf("struct field %q has invalid subbitmap-attribute %q in iso8583-tag", field.Name, attr)
			}
		}

		if attr, ok := attributes["tlvTag"]; ok {
			if _, err = strconv.ParseInt(attr, 16, 64); err != nil {
				return nil, fmt.Errorf("struct field %q has missing tag-attribute %q for a tlv encoding: %w", field.Name, attr, err)
			}
			element.TlvTag = attr
		}

		attr, ok = attributes["dataenc"]
		switch {
		case !ok:
			element.DataEncoding = DataEncodingDefault
		case attr == "ascii":
			element.DataEncoding = DataEncodingAscii
		case attr == "ebcdic":
			element.DataEncoding = DataEncodingEbcdic
		case attr == "bcd4":
			element.DataEncoding = DataEncodingBcd4
		case attr == "tlv":
			// tlv encoder requires the tlvTag to write out the data
			if element.TlvTag == "" {
				return nil, fmt.Errorf("struct field %q has missing tag-attribute %q for a tlv encoding: %w", field.Name, attr, err)
			}
			element.DataEncoding = DataEncodingTlv
		default:
			return nil, fmt.Errorf("struct field %q has invalid dataenc-attribute %q in iso8583-tag", field.Name, attr)
		}

		attr, ok = attributes["lenenc"]
		switch {
		case !ok:
			element.LengthEncoding = LengthEncodingDefault
		case attr == "ascii":
			element.LengthEncoding = LengthEncodingAscii
		case attr == "hex":
			element.LengthEncoding = LengthEncodingHex
		case attr == "bin":
			element.LengthEncoding = LengthEncodingBin
		case attr == "hexBit4":
			element.LengthEncoding = LengthEncodingHexBit4
		default:
			return nil, fmt.Errorf("struct field %q has invalid lenenc-attribute %q in iso8583-tag", field.Name, attr)
		}

		if attr, ok = attributes["bcd4len"]; ok {
			element.Bcd4Len, err = strconv.Atoi(attr)
			if err != nil {
				return nil, fmt.Errorf("struct field %q has invalid bcd4len-attribute %q in iso8583-tag", field.Name, attr)
			}
		}

		definitions[number] = element
	}

	if len(definitions) == 0 {
		return nil, fmt.Errorf("no iso8583 elements found in %T", v)
	}

	return definitions, nil
}

// parseTag parses the attributes from a struct field's iso8583-tag.
// Each attribute is a "key" or "key=value", separated by a comma and optional whitespace.
// The first attribute must be the element number with its format (representation + length) as value. Separate
// representation and length with a dash for fixed length fields,
// or use dots to indicate the size of the length indication
// field
// example: 1=n-3
// example: 2=ans...30, justify=left, minlength=3
//
//nolint:funlen,gocognit
func parseTag(tag string) (int, string, map[string]string) {
	var (
		number     int
		format     string
		attributes = make(map[string]string)
		// Find the "key" and "key=value" pairs
		start, i int
		first    = true
		l        = len(tag)
	)

	for {
		// Trim spaces
		for i < l && tag[i] == ' ' {
			i++
		}

		// Stop on EOL
		if i == l {
			break
		}

		// Find the key by scanning until EOL, ',' or '='
		start = i

		for i < l && tag[i] != ',' && tag[i] != '=' {
			if tag[i] == '\\' {
				// Skip next char
				i++
			}
			i++
		}

		if i == start {
			// Expected a name, but got EOL, a ',' or '='
			panic(fmt.Errorf("invalid syntax in tag %q near %q", tag, tag[i:]))
		}

		var (
			name  = tag[start:i]
			value string
		)

		if i < l && tag[i] == '=' {
			// Skip the '='
			i++

			// Find the value by scanning until EOL or ','
			start = i

			for i < len(tag) && tag[i] != ',' {
				if tag[i] == '\\' {
					i++
				}
				i++
			}

			value = tag[start:i]
		}

		// We are either at EOF or a ','
		if i < l {
			// Skip the ','
			i++
		}

		if first {
			// The first set is the element number and format
			number, _ = strconv.Atoi(name)
			format = value
			first = false
		} else {
			attributes[name] = value
		}
	}

	return number, format, attributes
}
