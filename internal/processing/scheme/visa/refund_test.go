package visa

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"gitlab.cmpayments.local/creditcard/platform/countrycode"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"

	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/visa/base1"
)

func TestMessageFromRefund(t *testing.T) {
	const (
		acceptorName = "MaxCorp Inc."
		city         = "Breda"
	)
	var (
		countryCode     = countrycode.Must("NLD")
		now             = time.Now()
		transactionTime = data.LocalTransactionDateTime(now)
	)

	type args struct {
		refund entity.Refund
	}
	tests := []struct {
		name string
		args args
		want *Message
	}{
		{
			name: "valid_message_from_refund",
			args: args{
				refund: entity.Refund{
					Card: entity.Card{
						Number: "4619031141704650",
						Expiry: entity.Expiry{Year: time.Now().Add(time.Hour * 24 * 730).Format("06"), Month: "05"},
					},
					CardAcceptor: entity.CardAcceptor{
						Name:         acceptorName,
						Address:      entity.CardAcceptorAddress{City: city, CountryCode: string(countryCode.Alpha3())},
						ID:           "cf78ed34-188a-48cf-ae24-0a2278e8d958",
						CategoryCode: "1334",
					},
					Currency:                 currencycode.Must("EUR"),
					Amount:                   200000,
					Stan:                     0000007,
					LocalTransactionDateTime: data.LocalTransactionDateTime(now),
					Source:                   entity.Ecommerce,
					InstitutionID:            "000002",
					CardSchemeData: entity.CardSchemeData{
						Request: entity.CardSchemeRequest{
							ProcessingCode: entity.ProcessingCode{
								TransactionTypeCode: "20",
								FromAccountTypeCode: "00",
								ToAccountTypeCode:   "00",
							},
							POSEntryMode: entity.POSEntryMode{
								PanEntryMode: entity.PANEntryManual,
								PinEntryMode: entity.PINEntryUnspecified,
							},
						},
					},
					VisaSchemeData: entity.VisaSchemeData{
						Request: entity.VisaSchemeRequest{
							PosConditionCode: "59",
							AdditionalPOSInformation: entity.AdditionalPOSInformation{
								TerminalType:                               "0", // unspecified
								TerminalEntryCapability:                    "1", // Terminal not used
								ChipConditionCode:                          "0",
								SpecialConditionIndicator:                  "0",                                                                                        // Default value
								ChipTransactionIndicator:                   "0",                                                                                        // Not applicable
								ChipCardAuthenticationReliabilityIndicator: "0",                                                                                        // Fill for field 60.7 present, or subsequent subfields that are present.
								TypeOrLevelIndicator:                       mapElectricCommerceIndicator(entity.Ecommerce, entity.ThreeDSecure{EcommerceIndicator: 7}), // Single transaction of a mail/phone order
								CardholderIDMethodIndicator:                "4",                                                                                        // Mail/Telephone/Electronic Commerce
								AdditionalAuthorizationIndicators:          "0",                                                                                        // Not applicable
							},
						},
					},
				},
			},
			want: &Message{
				Mti: iso8583.NewMti(authorizationRequestMTI),
				Fields: base1.Fields{
					F002_PrimaryAccountNumber:                   "4619031141704650",
					F003_ProcessingCode:                         "200000",
					F004_TransactionAmount:                      200000,
					F007_TransmissionDateTime:                   base1.F007FromTime(now),
					F011_SystemTraceAuditNumber:                 fmt.Sprintf("%06d", 0000007),
					F012_LocalTransactionTime:                   transactionTime.Format(`150405`),
					F013_LocalTransactionDate:                   transactionTime.Format(`0102`),
					F014_ExpirationDate:                         fmt.Sprintf("%s%s", time.Now().Add(time.Hour*24*730).Format("06"), "05"),
					F018_MerchantType:                           strconv.Itoa(1334),
					F019_AcquiringInstituteCountryCode:          countryCode.Numeric(),
					F022_PosEntryMode:                           "0100",
					F025_PosCondition:                           "59",
					F032_AcquiringInstitutionIdentificationCode: fmt.Sprintf("%06d", 474537),
					F037_RetrievalReferenceNumber:               now.Format(`0600215`)[1:] + strconv.Itoa(0000007),
					F042_CardAcceptorIdentificationCode:         "cf78ed34-188a-48cf-ae24-0a2278e8d958",
					F043_CardAcceptorNameLocation: base1.F043_CardAcceptorNameLocation{
						SF1_CarAcceptorName:  acceptorName,
						SF2_CardAcceptorCity: city,
						SF3_CountryCode:      countryCode.Alpha2(),
					},
					F049_TransactionCurrencyCode: currencycode.Must(currencycode.EUR).Numeric(),
					F060_AdditionalPointOfServiceInformation: base1.F060_AdditionalPOSInformation{
						B1: "01",
						B2: "00",
						B3: "00",
						B4: "00",
						B5: "07", // This is the default value for sending in the refunds
						// Page 349, visanet-authorization-only-online-messages-technical-specifications.pdf
						// Non-authenticated security transaction: Use to identify an electronic commerce transaction
						// that uses data encryption for security however, cardholder authentication is not performed
						// using a Visa approved protocol, such as 3-D Secure
						B6: "40",
					},
					F063_NetworkData: base1.F063_NetworkData{
						SF1_NetworkID: mapNetworkId(),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := messageFromRefund(tt.args.refund); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot =  %v\nwant = %v", got, tt.want)
			}
		})
	}
}
