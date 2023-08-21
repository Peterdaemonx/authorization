package entity

type Terminal struct {
	TerminalId         string
	TerminalCapability TerminalCapability
	TerminalLevel      TerminalLevel
}

// TerminalCapability is used to indicate what type of transaction a terminal can accept (chip, mag, contactless)
//
// The following fields will be mapped accordingly when send to mastercard and visa:
// MagneticStripeRead                                 = Magnetic stripe read capability only, visa: 2, mastercard: 2
// TerminalContactlessEMVinput                        = Terminal supports contactless EMV input, visa: 5, mastercard: 3
// EMVProximityReadCapableOnly                        = EMV Proximity-read-capable only, visa: 5, mastercard: 9
// TerminalEMVContactChipAndMagneticStripeAndKeyEntry = Terminal supports EMV contact chip input and magnetic stripe input and key entry input, visa: 5, mastercard: 8
type TerminalCapability string

const (
	MagneticStripeRead                                 TerminalCapability = "magneticStripeRead"
	TerminalContactlessEMVinput                        TerminalCapability = "terminalContactlessEmvInput"
	EMVProximityReadCapableOnly                        TerminalCapability = "emvProximityReadCapableOnly"
	TerminalEMVContactChipAndMagneticStripeAndKeyEntry TerminalCapability = "terminalEmvContactChipAndMagneticStripeAndKeyEntry"
)

var (
	terminalCapabilityMap = map[string]TerminalCapability{
		`magneticStripeRead`:                                 MagneticStripeRead,
		`terminalContactlessEmvInput`:                        TerminalContactlessEMVinput,
		`emvProximityReadCapableOnly`:                        EMVProximityReadCapableOnly,
		`terminalEmvContactChipAndMagneticStripeAndKeyEntry`: TerminalEMVContactChipAndMagneticStripeAndKeyEntry,
	}
)

func IsValidTerminalCapability(terminalCapability string) bool {
	_, exists := terminalCapabilityMap[terminalCapability]
	return exists
}

// TerminalLevel used to indicate that the terminal is unattended (CAT) currently we only accept automatedDispensingMachine which is used for vending machines and such.
type TerminalLevel string

const (
	AutomatedDispensingMachine TerminalLevel = "automatedDispensingMachine"
)

var (
	terminalLevelMap = map[string]TerminalLevel{
		`automatedDispensingMachine`: AutomatedDispensingMachine,
	}
)

func IsValidTerminalLevel(terminalLevel string) bool {
	_, exists := terminalLevelMap[terminalLevel]
	return exists
}