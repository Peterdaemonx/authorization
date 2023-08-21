package pos

import (
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

var (
	fromPanEntryCodeMap = map[string]entity.PANEntry{
		`00`: entity.PANEntryUnknown,
		`01`: entity.PANEntryManual,
		`02`: entity.PANEntryAutoMagStripeNotRequired,
		`05`: entity.PANEntryChip,
		`07`: entity.PANEntryContactless,
		`10`: entity.PANEntryCredentialOnFile,
		`81`: entity.PANEntryViaEcomWithOpId,
		`90`: entity.PANEntryMagneticStrip,
	}
	toPanEntryCodeMap = map[entity.PANEntry]string{
		entity.PANEntryUnknown:                  `00`,
		entity.PANEntryManual:                   `01`,
		entity.PANEntryAutoMagStripeNotRequired: `02`,
		entity.PANEntryChip:                     `05`,
		entity.PANEntryContactless:              `07`,
		entity.PANEntryCredentialOnFile:         `10`,
		entity.PANEntryViaEcomWithOpId:          `81`,
		entity.PANEntryMagneticStrip:            `90`,
	}
)

// PanEntryFromCode returns entity.PANEntry based of the code
// example: `05` => Contactless
func PanEntryFromCode(panEntryCode string) entity.PANEntry {
	panEntry, exists := fromPanEntryCodeMap[panEntryCode]
	if exists {
		return panEntry
	}
	return ""
}

// PanEntryCode returns the pan entry code defined as string
// example: Unknown => `00`
func PanEntryCode(panEntry entity.PANEntry) string {
	panentry, exists := toPanEntryCodeMap[panEntry]
	if exists {
		return panentry
	}
	return ""
}
