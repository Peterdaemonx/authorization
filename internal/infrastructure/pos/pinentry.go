package pos

import (
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

var (
	fromPinEntryCodeMap = map[string]entity.PINEntry{
		`0`: entity.PINEntryUnspecified,
		`1`: entity.PINEntryTerminalCanAcceptOnlinePin,
		`2`: entity.PINEntryTerminalCantAcceptOnlinePin,
		`3`: entity.PINEntryMPosSoftwarePinEntryCapability,
		`8`: entity.PINEntryTerminalPinPadDown,
	}
	toPinEntryCodeMap = map[entity.PINEntry]string{
		entity.PINEntryUnspecified:                    `0`,
		entity.PINEntryTerminalCanAcceptOnlinePin:     `1`,
		entity.PINEntryTerminalCantAcceptOnlinePin:    `2`,
		entity.PINEntryMPosSoftwarePinEntryCapability: `3`,
		entity.PINEntryTerminalPinPadDown:             `8`,
	}
)

// PinEntryCode returns the pin entry code defined as string
// example: Unspecified => `0`
func PinEntryCode(P entity.PINEntry) string {
	pinEntry, exists := toPinEntryCodeMap[P]
	if exists {
		return pinEntry
	}
	return ""
}

// PinEntryFromCode returns entity.PINEntry based of the code
// example: `8` => TerminalPinPadDown
func PinEntryFromCode(pinEntryValue string) entity.PINEntry {
	pinEntry, exists := fromPinEntryCodeMap[pinEntryValue]
	if exists {
		return pinEntry
	}
	return ""
}
