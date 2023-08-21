package ports

import (
	"reflect"
	"testing"

	"gitlab.cmpayments.local/creditcard/platform/http/validator"
)

func TestValidateCardAcceptor(t *testing.T) {
	tests := []struct {
		name   string
		c      CardAcceptor
		wanted map[string][]string
	}{
		{
			name: "valid card acceptor",
			c: CardAcceptor{
				ID:           "987288762212",
				CategoryCode: "0742",
				Name:         "Max Crop",
				City:         "Breda",
				Country:      "NLD",
				PostalCode:   "4876CA",
			},
			wanted: nil,
		},
		{
			name: "card acceptor no postalcode",
			c: CardAcceptor{
				ID:           "987288762212",
				CategoryCode: "0742",
				Name:         "Max Crop",
				City:         "Breda",
				Country:      "NLD",
				PostalCode:   "",
			},
			wanted: map[string][]string{
				"cardAcceptor.postalCode": {
					0: "postal code cannot be empty for NLD",
				},
			},
		},
		{
			name: "card_accaptor_country_not_EEA",
			c: CardAcceptor{
				ID:           "987288762212",
				CategoryCode: "0742",
				Name:         "Max Crop",
				City:         "Breda",
				Country:      "USA",
				PostalCode:   "4876CA",
			},
			wanted: map[string][]string{
				"cardAcceptor.country": {
					0: "cardAcceptor country not part of EEA",
				},
			},
		},
		{
			name: "invalid card acceptor country",
			c: CardAcceptor{
				ID:           "987288762212",
				CategoryCode: "0742",
				Name:         "Max Crop",
				City:         "Breda",
				Country:      "NL",
				PostalCode:   "4876CA",
			},
			wanted: map[string][]string{
				"cardAcceptor.country": {
					0: "invalid cardAcceptor country",
				},
			},
		},
		{
			name: "invalid card acceptor categorycode",
			c: CardAcceptor{
				ID:           "987288762212",
				CategoryCode: "9872",
				Name:         "Max Crop",
				City:         "Breda",
				Country:      "NLD",
				PostalCode:   "4876CA",
			},
			wanted: map[string][]string{
				"cardAcceptor.categoryCode": {
					0: "merchant category code not found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			tt.c.validate(v)
			if tt.wanted == nil {
				tt.wanted = make(map[string][]string)
			}
			if !reflect.DeepEqual(v.Errors, tt.wanted) {
				t.Errorf("validate(v) = %v, wanted: %v", v.Errors, tt.wanted)
			}
		})
	}
}
