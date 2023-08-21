package ports

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/platform/http/validator"
)

func TestValidAuthorizationRequest(t *testing.T) {
	tests := []struct {
		name   string
		a      authorizationRequest
		wanted map[string][]string
	}{
		{
			name: "valid_mastercard_authorization_request",
			a: authorizationRequest{
				Amount:                   100,
				Currency:                 "EUR",
				Reference:                "testRef",
				Source:                   "moto",
				LocalTransactionDateTime: data.LocalTransactionDateTime(time.Now()),
				InitialTraceID:           "testTrace123456",
				AuthorizationType:        "finalAuthorization",
				Card: Card{
					Holder: "testName",
					Number: "12345679876198",
					Cvv:    "123",
					Expiry: Expiry{
						Month: "12",
						Year:  strconv.Itoa(time.Now().Year())[2:],
					},
					scheme: "mastercard",
				},
				CardAcceptor: CardAcceptor{
					ID:           "123456789012",
					CategoryCode: "0742",
					Name:         "CMTICKETING",
					City:         "Breda",
					Country:      "NLD",
					PostalCode:   "4825BD",
				},
				CitMitIndicator: CitMitIndicator{
					InitiatedBy: "cardholder",
					SubCategory: "credentialOnFile",
				},
				Exemption: "merchantInitiated",
				ThreeDSecure: ThreeDSecure{
					AuthenticationVerificationValue: "MDA5OTAxMDUyNTExMTEwMDAwMDAwMDc4ODQwMDcwNzAwMDAwMDAwMA==",
					Version:                         "2",
					EcommerceIndicator:              1,
					DirectoryServerTransactionID:    "3bd2137d-08f1-4feb-ba50-3c2d4401c91a",
				},
			},
		},
		{
			name: "valid_visa_authorization_request",
			a: authorizationRequest{
				Amount:                   100,
				Currency:                 "EUR",
				Reference:                "testRef",
				Source:                   "moto",
				LocalTransactionDateTime: data.LocalTransactionDateTime(time.Now()),
				InitialTraceID:           "testTrace123456",
				Card: Card{
					Holder: "testName",
					Number: "12345679876198",
					Cvv:    "123",
					Expiry: Expiry{
						Month: "12",
						Year:  strconv.Itoa(time.Now().Year())[2:],
					},
					scheme: "visa",
				},
				CardAcceptor: CardAcceptor{
					ID:           "123456789012",
					CategoryCode: "0742",
					Name:         "CMTICKETING",
					City:         "Breda",
					Country:      "NLD",
					PostalCode:   "4825BD",
				},
				CitMitIndicator: CitMitIndicator{
					InitiatedBy: "cardholder",
					SubCategory: "credentialOnFile",
				},
				Exemption: "merchantInitiated",
				ThreeDSecure: ThreeDSecure{
					AuthenticationVerificationValue: "MDA5OTAxMDUyNTExMTEwMDAwMDAwMDc4ODQwMDcwNzAwMDAwMDAwMA==",
					Version:                         "2",
					EcommerceIndicator:              1,
					DirectoryServerTransactionID:    "3bd2137d-08f1-4feb-ba50-3c2d4401c91a",
				},
			},
		},
		{
			name: "invalid_currency_authorization_request",
			a: authorizationRequest{
				Amount:                   100,
				Currency:                 "BLABLA",
				Reference:                "testRef",
				Source:                   "moto",
				LocalTransactionDateTime: data.LocalTransactionDateTime(time.Now()),
				InitialTraceID:           "testTrace123456",
				Card: Card{
					Holder: "testName",
					Number: "12345679876198",
					Cvv:    "123",
					Expiry: Expiry{
						Month: "12",
						Year:  strconv.Itoa(time.Now().Year())[2:],
					},
					scheme: "visa",
				},
				CardAcceptor: CardAcceptor{
					ID:           "123456789012",
					CategoryCode: "0742",
					Name:         "CMTICKETING",
					City:         "Breda",
					Country:      "NLD",
					PostalCode:   "4825BD",
				},
				CitMitIndicator: CitMitIndicator{
					InitiatedBy: "cardholder",
					SubCategory: "credentialOnFile",
				},
				Exemption: "merchantInitiated",
				ThreeDSecure: ThreeDSecure{
					AuthenticationVerificationValue: "MDA5OTAxMDUyNTExMTEwMDAwMDAwMDc4ODQwMDcwNzAwMDAwMDAwMA==",
					Version:                         "2",
					EcommerceIndicator:              1,
					DirectoryServerTransactionID:    "3bd2137d-08f1-4feb-ba50-3c2d4401c91a",
				},
			},
			wanted: map[string][]string{
				"currency": {
					0: "unsupported currency",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			tt.a.validate(v)
			if tt.wanted == nil {
				tt.wanted = make(map[string][]string)
			}
			if !reflect.DeepEqual(v.Errors, tt.wanted) {
				t.Errorf("validate(v) = %v, wanted: %v", v.Errors, tt.wanted)
			}
		})
	}
}

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
