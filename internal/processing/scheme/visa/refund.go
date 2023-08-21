package visa

import (
	"fmt"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/pos"
	"gitlab.cmpayments.local/creditcard/clearing/pkg/mastercard/countrycode"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/visa/base1"
)

const (
	refundRequestMTI = `0100`
)

func RefundSchemeData(r *entity.Refund) {
	r.CardSchemeData.Request.ProcessingCode = entity.ProcessingCode{
		TransactionTypeCode: "20", // Purchase
		FromAccountTypeCode: "00", // Default Account
		ToAccountTypeCode:   "00", // Default Account
	}
	r.CardSchemeData.Request.POSEntryMode = posEntryMode(false, entity.CitMitIndicator{})
	r.VisaSchemeData.Request.PosConditionCode = posConditionCode(r.Source, false)
	r.VisaSchemeData.Request.AdditionalPOSInformation = additionalPOSInformation(r.Source, entity.ThreeDSecure{EcommerceIndicator: 7})
}

func messageFromRefund(r entity.Refund) *Message {
	return &Message{
		Mti: iso8583.NewMti(refundRequestMTI),
		Fields: base1.Fields{
			F002_PrimaryAccountNumber:                   r.Card.Number,
			F003_ProcessingCode:                         r.CardSchemeData.Request.ProcessingCode.String(),
			F004_TransactionAmount:                      int64(r.Amount),
			F007_TransmissionDateTime:                   base1.F007FromTime(time.Now()),
			F011_SystemTraceAuditNumber:                 fmt.Sprintf("%06d", r.Stan),
			F012_LocalTransactionTime:                   r.LocalTransactionDateTime.Format(`150405`),
			F013_LocalTransactionDate:                   r.LocalTransactionDateTime.Format(`0102`),
			F014_ExpirationDate:                         r.Card.Expiry.String(),
			F018_MerchantType:                           r.CardAcceptor.CategoryCode,
			F019_AcquiringInstituteCountryCode:          countrycode.NLD.Numeric(),
			F022_PosEntryMode:                           fmt.Sprintf("%s%s0", pos.PanEntryCode(r.CardSchemeData.Request.POSEntryMode.PanEntryMode), pos.PinEntryCode(r.CardSchemeData.Request.POSEntryMode.PinEntryMode)),
			F025_PosCondition:                           r.VisaSchemeData.Request.PosConditionCode,
			F032_AcquiringInstitutionIdentificationCode: entity.VisaInstitutionID,
			F037_RetrievalReferenceNumber:               retrievalReferenceNumber(r.Stan),
			F042_CardAcceptorIdentificationCode:         fmt.Sprintf("%s%s", r.Psp.Prefix, r.CardAcceptor.ID),
			F043_CardAcceptorNameLocation: base1.F043_CardAcceptorNameLocation{
				SF1_CarAcceptorName:  r.CardAcceptor.Name,
				SF2_CardAcceptorCity: r.CardAcceptor.Address.City,
				SF3_CountryCode:      countrycode.Must(countrycode.FromAlpha3(r.CardAcceptor.Address.CountryCode)).Alpha2(),
			},
			F049_TransactionCurrencyCode: r.Currency.Numeric(),
			F060_AdditionalPointOfServiceInformation: base1.F060_AdditionalPOSInformation{
				B1: fmt.Sprintf("%s%s", r.VisaSchemeData.Request.AdditionalPOSInformation.TerminalType, r.VisaSchemeData.Request.AdditionalPOSInformation.TerminalEntryCapability),
				B2: fmt.Sprintf("%s%s", r.VisaSchemeData.Request.AdditionalPOSInformation.ChipConditionCode, r.VisaSchemeData.Request.AdditionalPOSInformation.SpecialConditionIndicator),
				B3: "00",
				B4: fmt.Sprintf("%s%s", r.VisaSchemeData.Request.AdditionalPOSInformation.ChipTransactionIndicator, r.VisaSchemeData.Request.AdditionalPOSInformation.ChipCardAuthenticationReliabilityIndicator),
				B5: r.VisaSchemeData.Request.AdditionalPOSInformation.TypeOrLevelIndicator,
				B6: fmt.Sprintf("%s%s", r.VisaSchemeData.Request.AdditionalPOSInformation.CardholderIDMethodIndicator, r.VisaSchemeData.Request.AdditionalPOSInformation.AdditionalAuthorizationIndicators),
			},
			F063_NetworkData: base1.F063_NetworkData{
				SF1_NetworkID: mapNetworkId(),
			},
		},
	}
}
