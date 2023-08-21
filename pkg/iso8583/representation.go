//nolint:gochecknoglobals
package iso8583

import (
	"fmt"
)

const (
	alphabeticRepr representation = 1 << iota // alphabetic characters A–Z and a–z
	numericRepr                               // numeric digits 0–9
	spaceRepr                                 // space
	specialRepr                               // special characters
	binaryRepr                                // binary
)

type representation uint8

func (r representation) Alphabetic() bool {
	return r&alphabeticRepr == alphabeticRepr
}

func (r representation) Numeric() bool {
	return r&numericRepr == numericRepr
}

func (r representation) Space() bool {
	return r&spaceRepr == spaceRepr
}

func (r representation) Special() bool {
	return r&specialRepr == specialRepr
}

func (r representation) Binary() bool {
	return r&binaryRepr == binaryRepr
}

func (r representation) String() string {
	var s string
	if r.Alphabetic() {
		s += "a"
	}

	if r.Binary() {
		s += "b"
	}

	if r.Numeric() {
		s += "n"
	}

	if r.Space() || r.Special() {
		s += "s"
	}

	return s
}

// Assert returns an error if the input contains a character not allowed by the representation
func (r representation) Assert(in []byte) error {
	if r.Binary() {
		// Binary has no limitations
		return nil
	}

	for i, c := range in {
		if cr, ok := charset[c]; ok {
			if r&cr == 0 {
				return fmt.Errorf("character %q is not allowed", c)
			}
		} else {
			return fmt.Errorf("invalid character %#v at position %d", c, i)
		}
	}

	return nil
}

// Data Representation Attributes
// [IPM Clearing Formats • 15 October 2019 187]
//
//nolint:lll
var notation = map[string]representation{
	`a`:   alphabeticRepr,                                         // alphabetic characters A–Z and a–z
	`n`:   numericRepr,                                            // numeric digits 0–9
	`as`:  alphabeticRepr | spaceRepr,                             // alphabetic characters (A–Z and a–z), and space character
	`ns`:  numericRepr | spaceRepr | specialRepr,                  // numeric digits 0–9 and special characters (including space)
	`an`:  alphabeticRepr | numericRepr,                           // alphabetic (A–Z and a–z) and numeric characters
	`ans`: alphabeticRepr | numericRepr | spaceRepr | specialRepr, // alphabetic (A–Z and a–z), numeric, and special characters (including space)
	`b`:   binaryRepr,                                             // binary representation of data in eight-bit bytes
}

// Classification per character
// [Customer Interface Specification • 9 October 2018 64]
var charset = map[byte]representation{
	0x20: spaceRepr,      // <space>
	0x21: specialRepr,    // !
	0x22: specialRepr,    // "
	0x23: specialRepr,    // #
	0x24: specialRepr,    // $
	0x25: specialRepr,    // %
	0x26: specialRepr,    // &
	0x27: specialRepr,    // '
	0x28: specialRepr,    // (
	0x29: specialRepr,    // )
	0x2A: specialRepr,    // *
	0x2B: specialRepr,    // +
	0x2C: specialRepr,    // ,
	0x2D: specialRepr,    // -
	0x2E: specialRepr,    // .
	0x2F: specialRepr,    // /
	0x30: numericRepr,    // 0
	0x31: numericRepr,    // 1
	0x32: numericRepr,    // 2
	0x33: numericRepr,    // 3
	0x34: numericRepr,    // 4
	0x35: numericRepr,    // 5
	0x36: numericRepr,    // 6
	0x37: numericRepr,    // 7
	0x38: numericRepr,    // 8
	0x39: numericRepr,    // 9
	0x3A: specialRepr,    // :
	0x3B: specialRepr,    // ;
	0x3C: specialRepr,    // <
	0x3D: specialRepr,    // =
	0x3E: specialRepr,    // >
	0x3F: specialRepr,    // ?
	0x40: specialRepr,    // @
	0x41: alphabeticRepr, // A
	0x42: alphabeticRepr, // B
	0x43: alphabeticRepr, // C
	0x44: alphabeticRepr, // D
	0x45: alphabeticRepr, // E
	0x46: alphabeticRepr, // F
	0x47: alphabeticRepr, // G
	0x48: alphabeticRepr, // H
	0x49: alphabeticRepr, // I
	0x4A: alphabeticRepr, // J
	0x4B: alphabeticRepr, // K
	0x4C: alphabeticRepr, // L
	0x4D: alphabeticRepr, // M
	0x4E: alphabeticRepr, // N
	0x4F: alphabeticRepr, // O
	0x50: alphabeticRepr, // P
	0x51: alphabeticRepr, // Q
	0x52: alphabeticRepr, // R
	0x53: alphabeticRepr, // S
	0x54: alphabeticRepr, // T
	0x55: alphabeticRepr, // U
	0x56: alphabeticRepr, // V
	0x57: alphabeticRepr, // W
	0x58: alphabeticRepr, // X
	0x59: alphabeticRepr, // Y
	0x5A: alphabeticRepr, // Z
	0x5B: specialRepr,    // [
	0x5C: specialRepr,    // \
	0x5D: specialRepr,    // ]
	0x5E: specialRepr,    // ^
	0x5F: specialRepr,    // _
	0x60: specialRepr,    // `
	0x61: alphabeticRepr, // a
	0x62: alphabeticRepr, // b
	0x63: alphabeticRepr, // c
	0x64: alphabeticRepr, // d
	0x65: alphabeticRepr, // e
	0x66: alphabeticRepr, // f
	0x67: alphabeticRepr, // g
	0x68: alphabeticRepr, // h
	0x69: alphabeticRepr, // i
	0x6A: alphabeticRepr, // j
	0x6B: alphabeticRepr, // k
	0x6C: alphabeticRepr, // l
	0x6D: alphabeticRepr, // m
	0x6E: alphabeticRepr, // n
	0x6F: alphabeticRepr, // o
	0x70: alphabeticRepr, // p
	0x71: alphabeticRepr, // q
	0x72: alphabeticRepr, // r
	0x73: alphabeticRepr, // s
	0x74: alphabeticRepr, // t
	0x75: alphabeticRepr, // u
	0x76: alphabeticRepr, // v
	0x77: alphabeticRepr, // w
	0x78: alphabeticRepr, // x
	0x79: alphabeticRepr, // y
	0x7A: alphabeticRepr, // z
	0x7B: specialRepr,    // {
	0x7C: specialRepr,    // |
	0x7D: specialRepr,    // ]
	0x7E: specialRepr,    // ~
}
