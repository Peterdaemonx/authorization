package visa

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/pos"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/connection"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/visa/base1"
	"gitlab.cmpayments.local/creditcard/platform/countrycode"
)

const (
	reversalRequestMTI = `0400`
)

func messageFromReversal(r entity.Reversal) *Message {
	return &Message{
		Mti: iso8583.NewMti(reversalRequestMTI),
		Fields: base1.Fields{
			F002_PrimaryAccountNumber:                   r.Authorization.Card.Number,
			F003_ProcessingCode:                         r.Authorization.CardSchemeData.Request.ProcessingCode.String(),
			F004_TransactionAmount:                      int64(r.Amount),
			F007_TransmissionDateTime:                   base1.F007FromTime(r.Authorization.ProcessingDate),
			F011_SystemTraceAuditNumber:                 fmt.Sprintf("%06d", r.Authorization.Stan),
			F012_LocalTransactionTime:                   time.Now().Format("150405"),
			F013_LocalTransactionDate:                   time.Now().Format("0102"),
			F014_ExpirationDate:                         r.Authorization.Card.Expiry.String(),
			F018_MerchantType:                           r.Authorization.CardAcceptor.CategoryCode,
			F019_AcquiringInstituteCountryCode:          countrycode.Must("NLD").Mastercard.Numeric(),
			F022_PosEntryMode:                           fmt.Sprintf("%s%s0", pos.PanEntryCode(r.CardSchemeData.Request.POSEntryMode.PanEntryMode), pos.PinEntryCode(r.CardSchemeData.Request.POSEntryMode.PinEntryMode)),
			F025_PosCondition:                           r.Authorization.VisaSchemeData.Request.PosConditionCode,
			F032_AcquiringInstitutionIdentificationCode: entity.VisaInstitutionID,
			F037_RetrievalReferenceNumber:               r.Authorization.CardSchemeData.Request.RetrievalReferenceNumber,
			F038_AuthorizationIdenticationResponse:      r.Authorization.CardSchemeData.Response.AuthorizationIDResponse,
			F042_CardAcceptorIdentificationCode:         fmt.Sprintf("%s%s", r.Authorization.Psp.Prefix, r.Authorization.CardAcceptor.ID),
			F043_CardAcceptorNameLocation: base1.F043_CardAcceptorNameLocation{
				SF1_CarAcceptorName:  r.Authorization.CardAcceptor.Name,
				SF2_CardAcceptorCity: r.Authorization.CardAcceptor.Address.City,
				SF3_CountryCode:      countrycode.Must(r.Authorization.CardAcceptor.Address.CountryCode).Mastercard.Alpha2(),
			},
			F049_TransactionCurrencyCode: r.Authorization.Currency.Numeric(),
			F060_AdditionalPointOfServiceInformation: base1.F060_AdditionalPOSInformation{
				B1: fmt.Sprintf("%s%s", r.Authorization.VisaSchemeData.Request.AdditionalPOSInformation.TerminalType, r.Authorization.VisaSchemeData.Request.AdditionalPOSInformation.TerminalEntryCapability),
				B2: fmt.Sprintf("%s%s", r.Authorization.VisaSchemeData.Request.AdditionalPOSInformation.ChipConditionCode, r.Authorization.VisaSchemeData.Request.AdditionalPOSInformation.SpecialConditionIndicator),
				B3: "00",
				B4: fmt.Sprintf("%s%s", r.Authorization.VisaSchemeData.Request.AdditionalPOSInformation.ChipTransactionIndicator, r.Authorization.VisaSchemeData.Request.AdditionalPOSInformation.ChipCardAuthenticationReliabilityIndicator),
				B5: r.Authorization.VisaSchemeData.Request.AdditionalPOSInformation.TypeOrLevelIndicator,
				B6: fmt.Sprintf("%s%s", r.Authorization.VisaSchemeData.Request.AdditionalPOSInformation.CardholderIDMethodIndicator, r.Authorization.VisaSchemeData.Request.AdditionalPOSInformation.AdditionalAuthorizationIndicators),
			},
			F062_CustomPaymentServiceFields: base1.F062_CustomPaymentService{
				SF2_TransactionIdentifier: r.Authorization.VisaSchemeData.Response.TransactionId,
			},
			F063_NetworkData: base1.F063_NetworkData{
				SF1_NetworkID: mapNetworkId(),
				// TODO other scenario's still need to be implemented.
				//		2503: No confirmation from point of service
				//		2504: Partial dispense by ATM (misdispense) or POS partial reversal
				SF3_MessageReasonCode: mapReasonCode(r.Reason),
			},
			F090_OriginalDataElements: base1.F090_OriginalDataElements{
				SF1_OriginalMessageType:          authorizationRequestMTI,
				SF2_OriginalTraceNumber:          fmt.Sprintf("%06d", r.Authorization.Stan),
				SF3_OriginalTransmissionDateTime: r.Authorization.LocalTransactionDateTime.Format("0201150405"),
				SF4_OriginalAcquirerID:           "20814",
			},
		},
	}
}

func mapReasonCode(reason error) string {
	switch {
	// We've received a timeout or no free connection, result:
	// Transaction not completed, return 2502
	case errors.Is(reason, connection.ErrRequestTimeout):
		return "2502"
	case errors.Is(reason, connection.ErrNoFreeConnections):
		return "2502"
	// No error received. Transaction voided by customer, return 2501
	case reason == nil:
		return "2501"
	default:
		return "2501"
	}
}

func reversalResultFromMessage(msg Message) (entity.CardSchemeResponse, entity.VisaSchemeResponse) {
	return entity.CardSchemeResponse{
		Status: entity.AuthorizationStatusFromCardSchemeResponseCode(msg.Fields.F039_ResponseCode),
		ResponseCode: entity.ResponseCode{
			Value:       msg.Fields.F039_ResponseCode,
			Description: entity.ResponseDescriptionFromCode(msg.Fields.F039_ResponseCode),
		},
	}, entity.VisaSchemeResponse{TransactionId: msg.Fields.F062_CustomPaymentServiceFields.SF2_TransactionIdentifier}
}
