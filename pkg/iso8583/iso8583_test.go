package iso8583_test

import (
	"bytes"
	"reflect"
	"testing"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
)

type customType string

func (m customType) MarshalIso8583() ([]byte, error) {
	// Wrap the data
	d := []byte("Hello " + string(m) + "!")
	return d, nil
}

func (m *customType) UnmarshalIso8583(d []byte) error {
	// Unwrap the data
	*m = customType(d[6 : len(d)-1])
	return nil
}

//nolint:funlen
func TestIso8583(t *testing.T) {
	type subfields struct {
		Subfield1 string `iso8583:"1=n-1"`
		//nolint:structcheck,unused
		unexported int    `iso8583:"2=n-1, autofill"`
		Subfield3  string `iso8583:"3=n-1"`
		_          string `iso8583:"4=as-1, autofill"`
		Subfield5  int    `iso8583:"5=n-1, autofill, omitempty"`
		Subfield6  uint   `iso8583:"6=n-1"`
		Subfield7  int    `iso8583:"7=n-1, omitempty"`
	}

	type item struct {
		Key   string `iso8583:"1=an..99"`
		Value string `iso8583:"2=ans..99"`
	}

	type elements struct {
		String_DE10          string  `iso8583:"10=as.....30, minlength=5, justify=left"`
		StringPtr_DE11       *string `iso8583:"11=an-8, justify=right"`
		Int_DE20             int     `iso8583:"20=n-3,   justify=right"`
		IntPtr_DE21          *int    `iso8583:"21=ns.3,   justify=left, minlength=3"`
		Uint64Ptr_DE22       *uint64 `iso8583:"22=n...20, justify=left"`
		AutofillZeroInt_DE23 int     `iso8583:"23=n-3, autofill"`
		Fubar                string
		SubfieldsPtr_DE30    *subfields `iso8583:"30=ans.7, justify=right"`
		PtrSlicePtr_DE40     *[]*item   `iso8583:"40=ans..999"`
		MyType_DE50          customType `iso8583:"127=ans...999"`
	}

	var (
		MTI       = `0987`
		vString   = "een"
		vInt      = 12
		vUint64   = uint64(18446744073709551612)
		vZeroInt  = 0
		Subfields = subfields{
			Subfield1: "1",
			Subfield3: `3`,
			Subfield6: 6,
		}
		PtrSlicePtr = []*item{{"a", "first letter of alphabeth"}, {"Zzzzzz", "sleeping"}}
		customType  = customType("world")
	)

	in := elements{
		String_DE10:          vString,
		StringPtr_DE11:       &vString,
		Int_DE20:             vInt,
		IntPtr_DE21:          &vInt,
		Uint64Ptr_DE22:       &vUint64,
		AutofillZeroInt_DE23: vZeroInt,
		SubfieldsPtr_DE30:    &Subfields,
		PtrSlicePtr_DE40:     &PtrSlicePtr,
		MyType_DE50:          customType,
	}

	var buf bytes.Buffer
	encoder := iso8583.NewEncoder(&buf, iso8583.FormatAscii, iso8583.NoRdwLayout)
	decoder := iso8583.NewDecoder(&buf, iso8583.FormatAscii, iso8583.NoRdwLayout)

	// Encode
	if err := encoder.EncodeIso8583(iso8583.NewMti(MTI), &in); err != nil {
		t.Fatal(err)
	}

	encoded := buf.String()

	var actually string

	// First 4 bytes must be the mti
	actually, encoded = encoded[:4], encoded[4:]
	if actually != MTI {
		t.Fatalf("Expected MTI to be encoded as %q, got %q", MTI, actually)
	}

	// The primary bit map should have all the right fields set
	primary, encoded := []byte(encoded[:8]), encoded[8:]
	if expected := []byte{0x80, 0x60, 0x1e, 0x04, 0x1, 0x0, 0x0, 0x0}; !bytes.Equal(primary, expected) {
		t.Fatalf("Primary bit map does not represent the correct fields; %08b", primary)
	}

	// The secondary bit map should have all the right fields set
	secondary, encoded := []byte(encoded[:8]), encoded[8:]
	if expected := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2}; !bytes.Equal(secondary, expected) {
		t.Fatalf("Secondary bit map does not represent the correct fields; %08b", secondary)
	}

	// Expectations
	expectations := []struct {
		field    string
		length   int
		expected string
	}{
		{"DE-10", 10, `00005een  `},
		{"DE-11", 8, `00000een`},
		{"DE-20", 3, `012`},
		{"DE-21", 4, `312 `},
		{"DE-22", 23, `02018446744073709551612`},
		{"DE-23", 3, `000`},
		// LI=6, SF1=1, SF2={default}, SF3=3, SF4={default}, SF5=default, SF6=6, SF7={omitted}
		{"DE-30", 7, `6103 06`},
		{"DE-40", 50, `4801a25first letter of alphabeth06Zzzzzz08sleeping`},
		{"DE-50", 15, `012Hello world!`},
	}

	for _, expectation := range expectations {
		actually, encoded = encoded[:expectation.length], encoded[expectation.length:]
		if actually != expectation.expected {
			t.Fatalf("Expected %s to be encoded as %q, got %q", expectation.field, expectation.expected, actually)
		}
	}

	if len(encoded) != 0 {
		t.Fatalf("Expected no other data, got %q", encoded)
	}

	// Decode
	var (
		mti iso8583.MTI
		out elements
	)

	if err := decoder.DecodeIso8583(&mti, &out); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Errorf("Expected: %#v", in)
		t.Fatalf("Decoded:  %#v", out)
	}
}
