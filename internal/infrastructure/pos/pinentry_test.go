package pos

import (
	"testing"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

func TestPinEntryCode(t *testing.T) {
	tests := []struct {
		name     string
		PinEntry entity.PINEntry
		want     string
	}{
		{
			name:     "unspecified_is_0",
			PinEntry: entity.PINEntryUnspecified,
			want:     "0",
		},
		{
			name:     "MPosSoftwarePinEntryCapability_is_3",
			PinEntry: entity.PINEntryMPosSoftwarePinEntryCapability,
			want:     "3",
		},
		{
			name:     "TerminalCanAcceptOnlinePin_is_1",
			PinEntry: entity.PINEntryTerminalCanAcceptOnlinePin,
			want:     "1",
		},
		{
			name:     "TerminalCantAcceptOnlinePin_is_2",
			PinEntry: entity.PINEntryTerminalCantAcceptOnlinePin,
			want:     "2",
		},
		{
			name:     "TerminalPinPadDown_is_8",
			PinEntry: entity.PINEntryTerminalPinPadDown,
			want:     "8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PinEntryCode(tt.PinEntry); got != tt.want {
				t.Errorf("PinEntryCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPinEntryFromValue(t *testing.T) {
	tests := []struct {
		name         string
		pinEntryCode string
		want         entity.PINEntry
	}{
		{
			name:         "0_is_unspecified",
			pinEntryCode: "0",
			want:         entity.PINEntryUnspecified,
		},
		{
			name:         "3_is_MPosSoftwarePinEntryCapability",
			pinEntryCode: "3",
			want:         entity.PINEntryMPosSoftwarePinEntryCapability,
		},
		{
			name:         "1_is_TerminalCanAcceptOnlinePin",
			pinEntryCode: "1",
			want:         entity.PINEntryTerminalCanAcceptOnlinePin,
		},
		{
			name:         "2_is_TerminalCantAcceptOnlinePin",
			pinEntryCode: "2",
			want:         entity.PINEntryTerminalCantAcceptOnlinePin,
		},
		{
			name:         "8_is_TerminalPinPadDown",
			pinEntryCode: "8",
			want:         entity.PINEntryTerminalPinPadDown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PinEntryFromCode(tt.pinEntryCode); got != tt.want {
				t.Errorf("PinEntryFromCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
