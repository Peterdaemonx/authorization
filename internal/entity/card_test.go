package entity

import (
	"testing"

	_ "gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
)

func TestExpiry_MustMonthToInt(t *testing.T) {
	tests := []struct {
		name     string
		arg      Card
		expected int
	}{
		{
			name:     "valid month",
			arg:      Card{Expiry: Expiry{Year: "05", Month: "05"}},
			expected: 05,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			month := tt.arg.Expiry.MustMonthToInt()
			if tt.expected != month {
				t.Errorf("expected:%d\n\rgot:%d", tt.expected, month)
			}
		})
	}
}

func TestExpiry_MustYearToInt(t *testing.T) {
	tests := []struct {
		name     string
		arg      Card
		expected int
	}{
		{
			name:     "valid year",
			arg:      Card{Expiry: Expiry{Year: "05", Month: "05"}},
			expected: 2005,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year := tt.arg.Expiry.MustYearToInt()
			if tt.expected != year {
				t.Errorf("expected:%d\n\rgot:%d", tt.expected, year)
			}
		})
	}
}

func TestMaskPan(t *testing.T) {
	type args struct {
		pan string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "pan lower than 10 digits",
			args: args{pan: "0123456789"},
			want: "##########",
		},
		{
			name: "pan exceeding 10 digits",
			args: args{pan: "5204740000001002"},
			want: "52047400####1002",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskPan(tt.args.pan); got != tt.want {
				t.Errorf("MaskPan() = %v, want %v", got, tt.want)
			}
		})
	}
}
