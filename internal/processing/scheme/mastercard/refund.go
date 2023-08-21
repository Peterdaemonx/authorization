package mastercard

import (
	"fmt"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/mastercard/cis"
	"gitlab.cmpayments.local/creditcard/platform/countrycode"
)

func refundSchemeData(r *entity.Refund) {
	r.CardSchemeData.Request.ProcessingCode = entity.ProcessingCode{
		TransactionTypeCode: "20", // Refund
		FromAccountTypeCode: "00", // Default Account
		ToAccountTypeCode:   "00", // Default Account
	}
	r.CardSchemeData.Request.POSEntryMode = posEntryMode(r.Source, false)
	r.MastercardSchemeData.Request.AdditionalData = refundAdditionalData(r.Source)
	r.MastercardSchemeData.Request.PointOfServiceData = pointOfServiceData(r.Source, false, r.CardAcceptor.Address)
}

func refundResultFromMessage(msg Message) (entity.CardSchemeResponse, entity.MastercardSchemeResponse) {
	return entity.CardSchemeResponse{
			Status: entity.AuthorizationStatusFromCardSchemeResponseCode(msg.DataElements.DE39_ResponseCode),
			ResponseCode: entity.ResponseCode{
				Value:       msg.DataElements.DE39_ResponseCode,
				Description: entity.ResponseDescriptionFromCode(msg.DataElements.DE39_ResponseCode),
			},
			AuthorizationIDResponse: msg.DataElements.DE38_AuthorizationIdResponse,
			TraceId:                 entity.TraceIDFromString(fmt.Sprintf("%s%s", msg.DataElements.DE63_NetworkData, msg.DataElements.DE15_SettlementDate)).String(),
		},
		entity.MastercardSchemeResponse{
			AdditionalResponseData: parseAdditionalResponseData(msg.DataElements.DE44_AdditionalResponseData, msg.DataElements.DE39_ResponseCode),
			TraceID:                entity.TraceIDFromString(fmt.Sprintf("%s%s", msg.DataElements.DE63_NetworkData, msg.DataElements.DE15_SettlementDate)),
		}
}

func refundAdditionalData(source entity.Source) entity.AdditionalRequestData {
	return entity.AdditionalRequestData{
		TransactionCategoryCode: transactionCategoryCode(source),
		OriginalEcommerceIndicator: entity.SLI{
			SecurityProtocol:         9,
			CardholderAuthentication: 1,
			UCAFCollectionIndicator:  0,
		},
	}
}

func messageFromRefund(r entity.Refund) *Message {
	return &Message{
		Mti: iso8583.NewMti(AuthorizationRequestMTI),
		DataElements: cis.DataElements{
			DE2_PrimaryAccountNumber:             r.Card.Number,
			DE3_ProcessingCode:                   r.CardSchemeData.Request.ProcessingCode.String(),
			DE4_TransactionAmount:                int64(r.Amount),
			DE7_TransmissionDateTime:             cis.DE7FromTime(r.ProcessingDate),
			DE11_SystemTraceAuditNumber:          fmt.Sprintf("%06d", r.Stan),
			DE12_LocalTransactionTime:            r.LocalTransactionDateTime.Format("150405"),
			DE13_LocalTransactionDate:            r.LocalTransactionDateTime.Format("0102"),
			DE18_MerchantType:                    r.CardAcceptor.CategoryCode,
			DE20_PrimaryAccountNumberCountryCode: r.Card.Info.IssuerCountryCode,
			DE22_PointOfServiceEntryMode:         fmt.Sprintf("%s%s", r.CardSchemeData.Request.POSEntryMode.PanEntryMode, r.CardSchemeData.Request.POSEntryMode.PinEntryMode),
			DE32_AcquringInstitutionCode:         entity.MastercardInstitutionID,
			DE42_CardAcceptorCodeId:              r.CardAcceptor.ID,
			DE43_CardAcceptorNameAndLocation: &cis.DE43_CardAcceptorNameAndLocation{
				SF1_Name:               r.CardAcceptor.Name,
				SF3_City:               r.CardAcceptor.Address.City,
				SF5_StateOrCountryCode: r.CardAcceptor.Address.CountryCode,
			},
			DE48_AdditionalData: &cis.DE48_AdditionalData{
				TransactionCategoryCode: r.MastercardSchemeData.Request.AdditionalData.TransactionCategoryCode,
				SE42_ElectronicCommerceIndicators: &cis.DE48_SE42_ElectronicCommerceIndicators{
					SF1_SecurityLevelIndicatorAndUCAFCollectionIndicator: "910",
				},
			},
			DE49_TransactionCurrencyCode: r.Currency.Numeric(),
			DE61_PointOfServiceData: &cis.DE61_PointOfServiceData{
				SF1_TerminalAttendance:                        fmt.Sprintf("%d", r.MastercardSchemeData.Request.PointOfServiceData.TerminalAttendance),
				SF3_TerminalLocation:                          fmt.Sprintf("%d", r.MastercardSchemeData.Request.PointOfServiceData.TerminalLocation),
				SF4_CardholderPresence:                        fmt.Sprintf("%d", r.MastercardSchemeData.Request.PointOfServiceData.CardHolderPresence),
				SF5_CardPresence:                              fmt.Sprintf("%d", r.MastercardSchemeData.Request.PointOfServiceData.CardPresence),
				SF6_CardCaptureCapabilities:                   fmt.Sprintf("%d", r.MastercardSchemeData.Request.PointOfServiceData.CardCaptureCapabilities),
				SF7_TransactionStatus:                         fmt.Sprintf("%d", r.MastercardSchemeData.Request.PointOfServiceData.TransactionStatus),
				SF8_TransactionSecurity:                       fmt.Sprintf("%d", r.MastercardSchemeData.Request.PointOfServiceData.TransactionSecurity),
				SF10_CardholderActivatedTerminalLevel:         fmt.Sprintf("%d", r.MastercardSchemeData.Request.PointOfServiceData.CardHolderActivatedTerminalLevel),
				SF11_CardDataTerminalInputCapabilityIndicator: fmt.Sprintf("%d", r.MastercardSchemeData.Request.PointOfServiceData.CardDataTerminalInputCapabilityIndicator),
				SF12_AuthorizationLifeCycle:                   r.MastercardSchemeData.Request.PointOfServiceData.AuthorizationLifeCycle,
				SF13_CountryCode:                              countrycode.Must(r.MastercardSchemeData.Request.PointOfServiceData.CountryCode).Mastercard.Numeric(),
				SF14_PostalCode:                               r.MastercardSchemeData.Request.PointOfServiceData.PostalCode,
			},
		},
	}
}
