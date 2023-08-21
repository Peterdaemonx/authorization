package pos

import (
	"testing"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

func TestPANEntry_Value(t *testing.T) {
	tests := []struct {
		name     string
		PanEntry entity.PANEntry
		want     string
	}{
		{
			name: "00_is_unknown",
			want: `00`,
			PanEntry: entity.PANEntryUnknown,
		},
		{
			name: "01_is_PanManualEntry",
			want: `01`,
			PanEntry: entity.PANEntryManual,
		},
		{
			name: "05_is_chip",
			want: `05`,
			PanEntry: entity.PANEntryChip,
		},
		{
			name: "07_is_contactless",
			want: `07`,
			PanEntry: entity.PANEntryContactless,
		},
		{
			name: "10_is_credentialFile",
			want: `10`,
			PanEntry: entity.PANEntryCredentialOnFile,
		},
		{
			name: "81_is_entryViaEcomWithOpId",
			want: `81`,
			PanEntry: entity.PANEntryViaEcomWithOpId,
		},
		{
			name: "90_is_magneticStrip",
			want: `90`,
			PanEntry: entity.PANEntryMagneticStrip,
		},
		{
			name:"nothing",
			want: "",
			PanEntry: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PanEntryCode(tt.PanEntry); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPanEntryString(t *testing.T) {
	tests := []struct {
		name          string
		panEntryValue string
		want          entity.PANEntry
	}{
		{
			name:          "unknown_is_00",
			panEntryValue: `00`,
			want:          entity.PANEntryUnknown,
		},
		{
			name:          "PanManualEntry_is_01",
			panEntryValue: `01`,
			want:          entity.PANEntryManual,
		},
		{
			name:          "chip_is_05",
			panEntryValue: `05`,
			want:          entity.PANEntryChip,
		},
		{
			name:          "contactless_is_07",
			panEntryValue: `07`,
			want:          entity.PANEntryContactless,
		},
		{
			name:          "credentialFile_is_10",
			panEntryValue: `10`,
			want:          entity.PANEntryCredentialOnFile,
		},
		{
			name:          "entryViaEcomWithOpId_is_81",
			panEntryValue: `81`,
			want:          entity.PANEntryViaEcomWithOpId,
		},
		{
			name:          "magneticStrip_is_90",
			panEntryValue: `90`,
			want:          entity.PANEntryMagneticStrip,
		},
		{
			name:          "nothing",
			panEntryValue: "",
			want:          "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PanEntryFromCode(tt.panEntryValue); got != tt.want {
				t.Errorf("PanEntryString() = %v, want %v", got, tt.want)
			}
		})
	}
}

