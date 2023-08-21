package adapters

import (
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/pos"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/events"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

// TODO The concept of "event" isn't tied to any specific adapter, so this should be application/domain code
func CreatePublishAuthorizationRequest(a entity.Authorization, c entity.Capture) interface{} {
	message := events.AuthorizationCaptureMessageV1{
		AuthorizationID: a.ID.String(),
		CaptureID:       c.ID.String(),
		PSP: events.PSP{
			ID:   a.Psp.ID.String(),
			Name: a.Psp.Name,
		},
		PanTokenID:                             a.Card.PanTokenID,
		MaskedPan:                              a.Card.MaskedPan,
		CardScheme:                             a.Card.Info.Scheme,
		CardIssuerCountry:                      a.Card.Info.IssuerCountryCode,
		CardProductID:                          a.Card.Info.ProductID,
		CardProgramID:                          a.Card.Info.ProgramID,
		Amount:                                 a.Amount,
		Currency:                               a.Currency.Alpha3(),
		LocalTransactionDateTime:               time.Time(a.LocalTransactionDateTime),
		PartialCapture:                         !c.IsFinal,
		CapturedAt:                             time.Now(),
		ProcessingDate:                         a.ProcessingDate,
		ThreeDSVersion:                         a.ThreeDSecure.Version,
		ThreeDSAuthenticationVerificationValue: a.ThreeDSecure.AuthenticationVerificationValue,
		ThreeDSDirectoryServerTransactionID:    a.ThreeDSecure.DirectoryServerID,
		CustomerReference:                      c.Reference,
		ResponseCode:                           a.CardSchemeData.Response.ResponseCode.Value,
		Source:                                 string(a.Source),
		CardAcceptor: events.CardAcceptor{
			Name:           a.CardAcceptor.Name,
			City:           a.CardAcceptor.Address.City,
			Country:        a.CardAcceptor.Address.CountryCode,
			CategoryCode:   a.CardAcceptor.CategoryCode,
			CardAcceptorID: a.CardAcceptor.ID,
		},
		ProcessingCode: events.ProcessingCode{
			TransactionTypeCode: a.CardSchemeData.Request.ProcessingCode.TransactionTypeCode,
			FromAccountTypeCode: a.CardSchemeData.Request.ProcessingCode.FromAccountTypeCode,
			ToAccountTypeCode:   a.CardSchemeData.Request.ProcessingCode.ToAccountTypeCode,
		},
	}

	switch a.Card.Info.Scheme {
	case entity.Mastercard:
		message.MastercardSchemeData = mapMastercardSchemeData(a.Stan, a.MastercardSchemeData, a.CardSchemeData)
	}

	return message
}

func CreatePublishRefundRequest(r entity.Refund, c entity.RefundCapture) interface{} {
	message := events.RefundCaptureMessageV1{
		RefundID:  r.ID.String(),
		CaptureID: c.ID.String(),
		PSP: events.PSP{
			ID:   r.Psp.ID.String(),
			Name: r.Psp.Name,
		},
		PanTokenID:               r.Card.PanTokenID,
		MaskedPan:                r.Card.MaskedPan,
		CardScheme:               r.Card.Info.Scheme,
		CardIssuerCountry:        r.Card.Info.IssuerCountryCode,
		CardProductID:            r.Card.Info.ProductID,
		CardProgramID:            r.Card.Info.ProgramID,
		Amount:                   r.Amount,
		Currency:                 r.Currency.Alpha3(),
		LocalTransactionDateTime: time.Time(r.LocalTransactionDateTime),
		PartialCapture:           !c.IsFinal,
		CapturedAt:               time.Now(),
		ProcessingDate:           r.ProcessingDate,
		ResponseCode:             r.CardSchemeData.Response.ResponseCode.Value,
		Source:                   string(r.Source),
		CustomerReference:        r.CustomerReference,
		CardAcceptor: events.CardAcceptor{
			Name:           r.CardAcceptor.Name,
			City:           r.CardAcceptor.Address.City,
			Country:        r.CardAcceptor.Address.CountryCode,
			CategoryCode:   r.CardAcceptor.CategoryCode,
			CardAcceptorID: r.CardAcceptor.ID,
		},
		ProcessingCode: events.ProcessingCode{
			TransactionTypeCode: r.CardSchemeData.Request.ProcessingCode.TransactionTypeCode,
			FromAccountTypeCode: r.CardSchemeData.Request.ProcessingCode.FromAccountTypeCode,
			ToAccountTypeCode:   r.CardSchemeData.Request.ProcessingCode.ToAccountTypeCode,
		},
	}

	switch r.Card.Info.Scheme {
	case entity.Mastercard:
		message.MastercardSchemeData = mapMastercardSchemeData(r.Stan, r.MastercardSchemeData, r.CardSchemeData)
	}

	return message
}

func mapMastercardSchemeData(stan int, mcData entity.MastercardSchemeData, csData entity.CardSchemeData) *events.MastercardSchemeData {
	schemeData := &events.MastercardSchemeData{
		AuthorizationType:      string(mcData.Request.AuthorizationType),
		SystemTraceAuditNumber: stan,
		TraceID:                mcData.Response.TraceID.String(),
		FinancialNetworkCode:   mcData.Response.TraceID.FinancialNetworkCode,
		BanknetReferenceNumber: mcData.Response.TraceID.BanknetReferenceNumber,
		NetworkReportingDate:   mcData.Response.TraceID.NetworkReportingDate,
		PosPinCaptureCode:      mcData.Request.PosPinCaptureCode,
		POSEntryMode: events.POSEntryMode{
			PanEntryMode: pos.PanEntryCode(csData.Request.POSEntryMode.PanEntryMode),
			PinEntryMode: pos.PinEntryCode(csData.Request.POSEntryMode.PinEntryMode),
		},
		AdditionalData: events.AdditionalData{
			PinServiceCode: mcData.Request.AdditionalData.PinServiceCode,
		},
		PointOfServiceData: events.PointOfServiceData{
			TerminalAttendance:                       mcData.Request.PointOfServiceData.TerminalAttendance,
			TerminalLocation:                         mcData.Request.PointOfServiceData.TerminalLocation,
			CardHolderPresence:                       mcData.Request.PointOfServiceData.CardHolderPresence,
			CardPresence:                             mcData.Request.PointOfServiceData.CardPresence,
			CardCaptureCapabilities:                  mcData.Request.PointOfServiceData.CardCaptureCapabilities,
			CardHolderActivatedTerminalLevel:         mcData.Request.PointOfServiceData.CardHolderActivatedTerminalLevel,
			CardDataTerminalInputCapabilityIndicator: mcData.Request.PointOfServiceData.CardDataTerminalInputCapabilityIndicator,
		},
		ResponseReference: csData.Response.AuthorizationIDResponse,
	}

	if mcData.Response.AdditionalData.AppliedEcommerceIndicator != nil {
		schemeData.AdditionalData.EcommerceIndicators = events.EcommerceIndicators{
			SecurityProtocol:         mcData.Response.AdditionalData.AppliedEcommerceIndicator.SecurityProtocol,
			CardholderAuthentication: mcData.Response.AdditionalData.AppliedEcommerceIndicator.CardholderAuthentication,
			UCAFCollectionIndicator:  mcData.Response.AdditionalData.AppliedEcommerceIndicator.UCAFCollectionIndicator,
		}
	}

	return schemeData
}
