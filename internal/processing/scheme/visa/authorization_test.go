package visa

import (
	"reflect"
	"testing"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

func Test3dSecure(t *testing.T) {
	tests := []struct {
		name               string
		threeDSecure       entity.ThreeDSecure
		source             entity.Source
		expected60_8return string
	}{
		{
			name:               "MOTO_auth",
			threeDSecure:       entity.ThreeDSecure{},
			source:             entity.Moto,
			expected60_8return: "01",
		},
		{
			name:               "is_threedscure",
			threeDSecure:       entity.ThreeDSecure{},
			expected60_8return: "07",
		},
		{
			name:               "is_not_threedsecure",
			threeDSecure:       entity.ThreeDSecure{Version: "2"},
			expected60_8return: "00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := is3dSecure(tt.source, tt.threeDSecure)
			if !reflect.DeepEqual(got, tt.expected60_8return) {
				t.Errorf("got: %s want: %s", got, tt.expected60_8return)
			}
		})
	}
}

func TestCvvRequestData(t *testing.T) {
	tests := []struct {
		name        string
		cvv         string
		ssr         bool
		expectedCvv string
	}{
		{
			name:        "test_correct_cvv",
			cvv:         "936",
			ssr:         false,
			expectedCvv: "11 936",
		},
		{
			name:        "test_omit_cvv",
			cvv:         "",
			ssr:         false,
			expectedCvv: "01    ",
		},
		{
			name:        "test_cvv_length_4",
			cvv:         "0123",
			ssr:         false,
			expectedCvv: "110123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cvvRequestData(tt.cvv, tt.ssr)
			if !reflect.DeepEqual(got, tt.expectedCvv) {
				t.Errorf("got: %s, wanted: %s", got, tt.expectedCvv)
			}
		})
	}
}
