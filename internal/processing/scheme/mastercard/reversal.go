package mastercard

import (
	"fmt"
	"strconv"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/mastercard/cis"
	"gitlab.cmpayments.local/creditcard/platform/countrycode"
)

const (
	reversalRequestMTI  = `0400`
	reversalResponseMTI = `0410`
)

func reversalSchemeData(r *entity.Reversal) {
	r.Authorization.InstitutionID = fmt.Sprintf("%06d", 20814)
	r.Authorization.CardSchemeData.Request.CardHolderVerificationMethod = entity.CardHolderVerificationMethodOnlinePin
	r.Authorization.MastercardSchemeData.Request.AdditionalData.TransactionCategoryCode = transactionCategoryCode(r.Authorization.Source)
}

func reversalResultFromMessage(msg Message) entity.CardSchemeResponse {
	return entity.CardSchemeResponse{
		Status: entity.AuthorizationStatusFromCardSchemeResponseCode(msg.DataElements.DE39_ResponseCode),
		ResponseCode: entity.ResponseCode{
			Value:       msg.DataElements.DE39_ResponseCode,
			Description: entity.ResponseDescriptionFromCode(msg.DataElements.DE39_ResponseCode),
		},
	}
}

func messageFromReversal(r entity.Reversal) *Message {
	return &Message{
		Mti: iso8583.NewMti(reversalRequestMTI),
		DataElements: cis.DataElements{
			DE2_PrimaryAccountNumber:             r.Authorization.Card.Number,
			DE3_ProcessingCode:                   r.Authorization.CardSchemeData.Request.ProcessingCode.String(),
			DE4_TransactionAmount:                int64(r.Amount),
			DE7_TransmissionDateTime:             cis.DE7FromTime(r.ProcessingDate),
			DE11_SystemTraceAuditNumber:          fmt.Sprintf("%06d", r.Authorization.Stan),
			DE12_LocalTransactionTime:            r.Authorization.LocalTransactionDateTime.Format("150405"),
			DE13_LocalTransactionDate:            r.Authorization.LocalTransactionDateTime.Format("0102"),
			DE14_ExpirationDate:                  r.Authorization.Card.Expiry.String(),
			DE18_MerchantType:                    r.Authorization.CardAcceptor.CategoryCode,
			DE20_PrimaryAccountNumberCountryCode: r.Authorization.Card.Info.IssuerCountryCode,
			DE22_PointOfServiceEntryMode:         fmt.Sprintf("%s%s", r.Authorization.CardSchemeData.Request.POSEntryMode.PanEntryMode, r.Authorization.CardSchemeData.Request.POSEntryMode.PinEntryMode),
			DE32_AcquringInstitutionCode:         entity.MastercardInstitutionID,
			DE38_AuthorizationIdResponse:         r.Authorization.CardSchemeData.Response.AuthorizationIDResponse,
			DE39_ResponseCode:                    r.Authorization.CardSchemeData.Response.ResponseCode.Value,
			DE42_CardAcceptorCodeId:              r.Authorization.CardAcceptor.ID,
			DE43_CardAcceptorNameAndLocation: &cis.DE43_CardAcceptorNameAndLocation{
				SF1_Name:               r.Authorization.CardAcceptor.Name,
				SF3_City:               r.Authorization.CardAcceptor.Address.PostalCode,
				SF5_StateOrCountryCode: r.Authorization.CardAcceptor.Address.CountryCode,
			},
			DE48_AdditionalData: &cis.DE48_AdditionalData{
				TransactionCategoryCode:           r.Authorization.MastercardSchemeData.Request.AdditionalData.TransactionCategoryCode,
				SE20_CardholderVerificationMethod: mapCardHolderVerificationMethod(r.Authorization.CardSchemeData.Request.CardHolderVerificationMethod),
				SE63_TraceId:                      cis.NewDE48_SE63_TraceId(r.Authorization.MastercardSchemeData.Response.TraceID),
			},
			DE49_TransactionCurrencyCode: r.Authorization.Currency.Numeric(),
			DE61_PointOfServiceData: &cis.DE61_PointOfServiceData{
				SF1_TerminalAttendance:                        fmt.Sprintf("%d", r.Authorization.MastercardSchemeData.Request.PointOfServiceData.TerminalAttendance),
				SF3_TerminalLocation:                          fmt.Sprintf("%d", r.Authorization.MastercardSchemeData.Request.PointOfServiceData.TerminalLocation),
				SF4_CardholderPresence:                        fmt.Sprintf("%d", r.Authorization.MastercardSchemeData.Request.PointOfServiceData.CardHolderPresence),
				SF5_CardPresence:                              fmt.Sprintf("%d", r.Authorization.MastercardSchemeData.Request.PointOfServiceData.CardPresence),
				SF6_CardCaptureCapabilities:                   fmt.Sprintf("%d", r.Authorization.MastercardSchemeData.Request.PointOfServiceData.CardCaptureCapabilities),
				SF7_TransactionStatus:                         fmt.Sprintf("%d", r.Authorization.MastercardSchemeData.Request.PointOfServiceData.TransactionStatus),
				SF8_TransactionSecurity:                       fmt.Sprintf("%d", r.Authorization.MastercardSchemeData.Request.PointOfServiceData.TransactionSecurity),
				SF10_CardholderActivatedTerminalLevel:         fmt.Sprintf("%d", r.Authorization.MastercardSchemeData.Request.PointOfServiceData.CardHolderActivatedTerminalLevel),
				SF11_CardDataTerminalInputCapabilityIndicator: fmt.Sprintf("%d", r.Authorization.MastercardSchemeData.Request.PointOfServiceData.CardDataTerminalInputCapabilityIndicator),
				SF12_AuthorizationLifeCycle:                   r.Authorization.MastercardSchemeData.Request.PointOfServiceData.AuthorizationLifeCycle,
				SF13_CountryCode:                              countrycode.Must(r.Authorization.MastercardSchemeData.Request.PointOfServiceData.CountryCode).Mastercard.Numeric(),
				SF14_PostalCode:                               r.Authorization.MastercardSchemeData.Request.PointOfServiceData.PostalCode,
			},
			DE90_OriginalDataElements: &cis.DE90_OriginalDataElements{
				SF1_OriginalMessageTypeIdentifier:     `0100`,
				SF2_OriginalSystemTraceAuditNumber:    strconv.Itoa(r.Authorization.Stan),
				SF3_OriginalTransmissionDateAndTime:   cis.DE7FromTime(r.Authorization.ProcessingDate),
				SF4_OriginalAcquiringInstituteIdCode:  entity.MastercardInstitutionID,
				SF5_OriginalForwardingInstituteIdCode: entity.MastercardInstitutionID,
			},
		},
	}
}
