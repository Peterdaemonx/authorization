package mastercard

import (
	"fmt"
	"time"

	"gitlab.cmpayments.local/creditcard/platform/countrycode"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/pos"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/mastercard/cis"
)

const (
	AuthorizationRequestMTI  = `0100`
	AuthorizationResponseMTI = `0110`
)

type LowRiskMerchantIndicator int

func authorizationSchemeData(a *entity.Authorization) {
	a.CardSchemeData.Request.ProcessingCode = entity.ProcessingCode{
		TransactionTypeCode: "00", // Purchase
		FromAccountTypeCode: "00", // Default Account
		ToAccountTypeCode:   "00", // Default Account
	}
	a.CardSchemeData.Request.POSEntryMode = posEntryMode(a.Source, a.Recurring.Subsequent)
	a.MastercardSchemeData.Request.AdditionalData = additionalData(*a)
	a.MastercardSchemeData.Request.PointOfServiceData = pointOfServiceData(a.Source, a.Recurring.Subsequent, a.CardAcceptor.Address)
}

func mapCardHolderVerificationMethod(cvm entity.CardHolderVerificationMethod) string {
	switch cvm {
	case entity.CardHolderVerificationMethodSignature, entity.CardHolderVerificationMethodUnattendedNoPin:
		return "S"
	default:
		return "P"
	}
}

func posEntryMode(s entity.Source, subseqRecurring bool) entity.POSEntryMode {
	switch {
	case s == entity.Moto || subseqRecurring:
		return entity.POSEntryMode{
			PanEntryMode: entity.PANEntryCredentialOnFile, // Credential on File
			PinEntryMode: entity.PINEntryUnspecified,      // Unspecified or unknown
		}
	default:
		return entity.POSEntryMode{
			PanEntryMode: entity.PANEntryViaEcomWithOpId, // PAN/Token entry via electronic commerce with optional Identity Check-AAV or DSRP cryptogram in UCAF
			PinEntryMode: entity.PINEntryUnspecified,     // Unspecified or unknown
		}
	}
}

func additionalData(a entity.Authorization) entity.AdditionalRequestData {
	ad := entity.AdditionalRequestData{
		TransactionCategoryCode:    transactionCategoryCode(a.Source),
		OriginalEcommerceIndicator: originalEcommerceIndicator(a.ThreeDSecure, a.Recurring, a.Exemption),
		AuthenticationData:         authenticationData(a.ThreeDSecure),
		PinServiceCode:             "",
	}

	if a.Exemption != "" {
		ad.LowRiskIndicator = mapLowRiskIndicator(string(a.Exemption))
	}

	return ad
}

func pointOfServiceData(source entity.Source, subseqRecurring bool, caa entity.CardAcceptorAddress) entity.PointOfServiceData {
	switch source {
	case entity.Moto:
		return entity.PointOfServiceData{
			TerminalAttendance:                       1, // Unattended terminal
			TerminalLocation:                         0, // On premises of card acceptor facility
			CardHolderPresence:                       mapCardHolderPresence(source, subseqRecurring),
			CardPresence:                             1, // Card not present
			CardCaptureCapabilities:                  0, // Terminal/operator does not have card capture capability
			TransactionStatus:                        2, // Identity Check Phone Order
			TransactionSecurity:                      0, // No security concern
			CardHolderActivatedTerminalLevel:         6, // Authorized Level 6 CAT: Electronic commerce
			CardDataTerminalInputCapabilityIndicator: 0, // Input capability unknown or unspecified
			AuthorizationLifeCycle:                   "00",
			CountryCode:                              caa.CountryCode,
			PostalCode:                               caa.PostalCode,
		}
	default:
		return entity.PointOfServiceData{
			TerminalAttendance:                       1, // Unattended terminal
			TerminalLocation:                         4, // On premises of card acceptor facility
			CardHolderPresence:                       mapCardHolderPresence(source, subseqRecurring),
			CardPresence:                             1, // Card not present
			CardCaptureCapabilities:                  0, // Terminal/operator does not have card capture capability
			TransactionStatus:                        0, // Normal request (original presentment)
			TransactionSecurity:                      0, // No security concern
			CardHolderActivatedTerminalLevel:         6, // Authorized Level 6 CAT: Electronic commerce
			CardDataTerminalInputCapabilityIndicator: 0, // Input capability unknown or unspecified
			AuthorizationLifeCycle:                   "00",
			CountryCode:                              caa.CountryCode,
			PostalCode:                               caa.PostalCode,
		}
	}
}

func mapCardHolderPresence(s entity.Source, subseqRecurring bool) int {
	switch {
	case s == entity.Moto:
		return 2 // Cardholder not present (mail/facsimile order)
	case subseqRecurring:
		return 4 // Standing order/recurring transactions
	case s == entity.Ecommerce:
		return 5 // Cardholder not present (Electronic order [home PC, Internet, mobile phone, PDA])
	}

	return 1 // Cardholder not present, unspecified
}

func transactionCategoryCode(s entity.Source) string {
	switch s {
	case entity.Ecommerce:
		return "T"
	case entity.Moto:
		return "T"
	default:
		return " " // We're returning a space since this is what MC wants
	}
}

func mapLowRiskIndicator(e string) string {
	switch e {
	case "merchantInitiated":
		return "1"
	case "transactionRiskAnalysis":
		return "2"
	case "recurring":
		return "3"
	case "lowValue":
		return "4"
	case "scaDelegation":
		return "5"
	case "secureCorporate":
		return "6"
	default:
		return ""
	}
}

// SE42 Electronic Commerce Indicators SF 1 = Security Protocol.
// Available values are "Channel" [2] and "None" [9]
func mapSecurityProtocol(t entity.ThreeDSecure, r entity.Recurring) int {
	switch {
	case t == entity.ThreeDSecure{} || t.EcommerceIndicator == 0:
		return 9
	default:
		return 2
	}
}

// SE42 Electronic Commerce Indicators SF 3 UCAF Collection Indicator
func mapUCAFCollectionIndicator(t entity.ThreeDSecure, r entity.Recurring) int {
	switch {
	case r.Subsequent:
		return 7 // authentication from previous payment
	default: // default we get the right ecommerce indicator in th 3DS object
		return t.EcommerceIndicator
	}
}

func originalEcommerceIndicator(t entity.ThreeDSecure, r entity.Recurring, exemption entity.ExemptionType) entity.SLI {
	switch {
	case exemption == entity.MerchantInitiatedExemption:
		return entity.SLI{
			SecurityProtocol:         2,
			CardholderAuthentication: 1,
			UCAFCollectionIndicator:  7,
		}
	case exemption != "":
		return entity.SLI{
			SecurityProtocol:         2,
			CardholderAuthentication: 1,
			UCAFCollectionIndicator:  6,
		}
	default:
		return entity.SLI{
			SecurityProtocol:         mapSecurityProtocol(t, r),
			CardholderAuthentication: 1,
			UCAFCollectionIndicator:  mapUCAFCollectionIndicator(t, r),
		}
	}
}

func mapThreeDVersion(a entity.ThreeDSecure) string {
	if len(a.Version) > 1 {
		return a.Version[0:1]
	}

	return a.Version
}

// SE66 Authentication Data
func authenticationData(t entity.ThreeDSecure) entity.AuthenticationData {
	return entity.AuthenticationData{
		ProgramProtocol:              mapThreeDVersion(t),
		DirectoryServerTransactionID: t.DirectoryServerID,
	}
}

func authorizationResultFromMessage(msg Message) (entity.CardSchemeResponse, entity.MastercardSchemeResponse) {
	csr := entity.CardSchemeResponse{
		Status: entity.AuthorizationStatusFromCardSchemeResponseCode(msg.DataElements.DE39_ResponseCode),
		ResponseCode: entity.ResponseCode{
			Value:       msg.DataElements.DE39_ResponseCode,
			Description: entity.ResponseDescriptionFromCode(msg.DataElements.DE39_ResponseCode),
		},
		AuthorizationIDResponse: msg.DataElements.DE38_AuthorizationIdResponse,
		TraceId:                 entity.TraceIDFromString(fmt.Sprintf("%s%s", msg.DataElements.DE63_NetworkData, msg.DataElements.DE15_SettlementDate)).String(),
	}

	msr := entity.MastercardSchemeResponse{
		AdditionalResponseData: parseAdditionalResponseData(msg.DataElements.DE44_AdditionalResponseData, msg.DataElements.DE39_ResponseCode),
		TraceID:                entity.TraceIDFromString(fmt.Sprintf("%s%s", msg.DataElements.DE63_NetworkData, msg.DataElements.DE15_SettlementDate)),
		AdditionalData: entity.AdditionalResponseData{
			AppliedEcommerceIndicator: &entity.SLI{},
			ReasonForUCAFDowngrade:    nil,
		},
	}

	if msg.DataElements.DE48_AdditionalData != nil {
		if msg.DataElements.DE48_AdditionalData.SE42_ElectronicCommerceIndicators != nil {
			csr.EcommerceIndicator = entity.NewAppliedEcommerceIndicator(entity.SLIFromString(msg.DataElements.DE48_AdditionalData.SE42_ElectronicCommerceIndicators.SF1_SecurityLevelIndicatorAndUCAFCollectionIndicator)).UCAFCollectionIndicator
			msr.AdditionalData = entity.AdditionalResponseData{
				AppliedEcommerceIndicator: entity.NewAppliedEcommerceIndicator(entity.SLIFromString(msg.DataElements.DE48_AdditionalData.SE42_ElectronicCommerceIndicators.SF1_SecurityLevelIndicatorAndUCAFCollectionIndicator)),
				ReasonForUCAFDowngrade:    entity.ReasonForUCAFDowngradeFromString(msg.DataElements.DE48_AdditionalData.SE42_ElectronicCommerceIndicators.SF3_ReasonForUCAFDowngrade),
			}
		}
	}

	return csr, msr
}

func MapResponseCode(code string) (string, string) {
	// [Customer Interface Specifications.pdf] p314-316
	switch code {
	case "00", "10":
		return "approved", "Approved"
	case "01":
		return "issuer_declined", "Issuer Declined"
	case "03":
		return "invalid_merchant", "Invalid merchant"
	case "04":
		return "capture_card", "Capture card"
	case "05":
		return "do_not_honor", "Do not honor"
	case "08":
		return "honor_with_id", "Honor with ID"
	case "12":
		return "invalid_transaction", "Invalid transaction"
	case "13":
		return "Invalid amount", "Invalid amount"
	case "14":
		return "invalid_cardnumber", "Invalid card number"
	case "15":
		return "invalid_issuer", "Invalid issuer"
	case "30":
		return "format_error", "Format error"
	case "41":
		return "lost_card", "Lost card"
	case "43":
		return "stolen_card", "Stolen card"
	case "51":
		return "insufficient_funds", "Insufficient funds"
	case "54":
		return "expired_card", "Expired card"
	case "55":
		return "invalid_pin", "Invalid PIN"
	case "57":
		return "transaction_not_permitted_to_issuer", "Transaction not permitted to issuer/cardholder"
	case "58":
		return "transaction_not_permitted_to_terminal", "Transaction not permitted to acquirer/terminal"
	case "61":
		return "exceeds_limit", "Exceeds withdrawal amount limit"
	case "62":
		return "restricted_card", "Restricted card"
	case "63":
		return "security_violation", "Security violation"
	case "65":
		return "soft_decline", "Soft Decline"
	case "70":
		return "issuer_declined", "Contact Card Issuer"
	case "71":
		return "pin_not_changed", "PIN Not Changed"
	case "75":
		return "pin_tries_exceeded", "Allowable number of PIN tries exceeded"
	case "76":
		return "invalid_transaction", "Invalid/nonexistent “To Account” specified"
	case "77":
		return "invalid_transaction", "Invalid/nonexistent “From Account” specified"
	case "78":
		return "invalid_transaction", "Invalid/nonexistent account specified (general)"
	case "81":
		return "invalid_transaction", "Domestic Debit Transaction Not Allowed (Regional use only)"
	case "84":
		return "invalid_transaction", "Invalid Authorization Life Cycle"
	case "85":
		return "valid ", "Not declined Valid for all zero amount transactions."
	case "86":
		return "pin_not_possible", "PIN Validation not possible"
	case "87":
		return "no_cash_back_allowed ", "Purchase Amount Only, No Cash Back Allowed"
	case "88":
		return "cryptographic_failure", "Cryptographic failure"
	case "89":
		return "invalid_pin", "Unacceptable PIN-Transaction Declined-Retry"
	case "91":
		return "issuer_declined", "Authorization System or issuer system inoperative"
	case "92":
		return "invalid_issuer", "Unable to route transaction"
	case "94":
		return "duplicate", "Duplicate transmission detected"
	case "96":
		return "system_error", "System error"
	default:
		return "unknown_code", "received unknown code " + code
	}
}

// nolint:unused
func parseAdditionalResponseData(ard, rc string) string {
	switch rc {
	case "01":
		return fmt.Sprintf("call issuer: %s", ard)
	case "30":
		switch len(ard) {
		case 3:
			return fmt.Sprintf("format error in element %s", ard)
		case 5:
			return fmt.Sprintf("format error in element %s, subfield %s", ard[0:3], ard[3:5])
		case 6:
			return fmt.Sprintf("format error in element %s, subfield %s", ard[0:3], ard[3:6])
		default:
			return ard
		}
	default:
		return ard
	}
}

func mapInitiator(file entity.CitMitIndicator) string {
	var initiatedBy, subCategory string

	switch file.InitiatedBy {
	case entity.Cardholder:
		initiatedBy = "C1"
	case entity.MITRecurring:
		initiatedBy = "M1"
	case entity.MITIndustryPractice:
		initiatedBy = "M2"
	}

	switch file.SubCategory {
	case entity.CredentialOnFile, entity.UnscheduledCredentialOnFile:
		subCategory = "01"
	case entity.StandingOrder:
		subCategory = "02"
	case entity.Subscription:
		subCategory = "03"
	case entity.Installment:
		subCategory = "04"
	case entity.PartialShipment:
		subCategory = "05"
	case entity.DelayedCharge:
		subCategory = "06"
	case entity.NoShow:
		subCategory = "07"
	case entity.Resubmission:
		subCategory = "08"
	}

	return initiatedBy + subCategory
}

func messageFromAuthorization(a entity.Authorization) *Message {
	return &Message{
		Mti: iso8583.NewMti(AuthorizationRequestMTI),
		DataElements: cis.DataElements{
			DE2_PrimaryAccountNumber:             a.Card.Number,
			DE3_ProcessingCode:                   a.CardSchemeData.Request.ProcessingCode.String(),
			DE4_TransactionAmount:                int64(a.Amount),
			DE7_TransmissionDateTime:             cis.DE7FromTime(time.Now()),
			DE11_SystemTraceAuditNumber:          fmt.Sprintf("%06d", a.Stan),
			DE12_LocalTransactionTime:            a.LocalTransactionDateTime.Format("150405"),
			DE13_LocalTransactionDate:            a.LocalTransactionDateTime.Format("0102"),
			DE14_ExpirationDate:                  a.Card.Expiry.String(),
			DE18_MerchantType:                    a.CardAcceptor.CategoryCode,
			DE20_PrimaryAccountNumberCountryCode: a.Card.Info.IssuerCountryCode,
			DE22_PointOfServiceEntryMode:         fmt.Sprintf("%s%s", pos.PanEntryCode(a.CardSchemeData.Request.POSEntryMode.PanEntryMode), pos.PinEntryCode(a.CardSchemeData.Request.POSEntryMode.PinEntryMode)),
			DE32_AcquringInstitutionCode:         entity.MastercardInstitutionID,
			DE42_CardAcceptorCodeId:              a.CardAcceptor.ID,
			DE43_CardAcceptorNameAndLocation: &cis.DE43_CardAcceptorNameAndLocation{
				SF1_Name:               a.CardAcceptor.Name,
				SF3_City:               a.CardAcceptor.Address.City,
				SF5_StateOrCountryCode: a.CardAcceptor.Address.CountryCode,
			},
			DE48_AdditionalData: &cis.DE48_AdditionalData{
				TransactionCategoryCode:                     a.MastercardSchemeData.Request.AdditionalData.TransactionCategoryCode,
				SE22_MultiPurposeMerchantIndicator:          cis.NewDE48_SE22MultiPurposeMerchantIndicator(a.MastercardSchemeData.Request.AdditionalData.LowRiskIndicator, mapInitiator(a.CitMitIndicator)),
				SE42_ElectronicCommerceIndicators:           cis.NewDE48_SE42_ElectronicCommerceIndicators(a.MastercardSchemeData.Request.AdditionalData.OriginalEcommerceIndicator),
				SE43_UniversalCardholderAuthenticationField: a.ThreeDSecure.AuthenticationVerificationValue,
				SE63_TraceId:                    cis.NewDE48_SE63_TraceId(entity.TraceIDFromString(a.Recurring.TraceID)),
				SE66_AuthenticationData:         cis.NewDE48_SE66_AuthenticationData(a.MastercardSchemeData.Request.AdditionalData.AuthenticationData),
				SE92_CardholderVerificationCode: a.Card.Cvv,
			},
			DE49_TransactionCurrencyCode: a.Currency.Numeric(),
			DE61_PointOfServiceData: &cis.DE61_PointOfServiceData{
				SF1_TerminalAttendance:                        fmt.Sprintf("%d", a.MastercardSchemeData.Request.PointOfServiceData.TerminalAttendance),
				SF3_TerminalLocation:                          fmt.Sprintf("%d", a.MastercardSchemeData.Request.PointOfServiceData.TerminalLocation),
				SF4_CardholderPresence:                        fmt.Sprintf("%d", a.MastercardSchemeData.Request.PointOfServiceData.CardHolderPresence),
				SF5_CardPresence:                              fmt.Sprintf("%d", a.MastercardSchemeData.Request.PointOfServiceData.CardPresence),
				SF6_CardCaptureCapabilities:                   fmt.Sprintf("%d", a.MastercardSchemeData.Request.PointOfServiceData.CardCaptureCapabilities),
				SF7_TransactionStatus:                         fmt.Sprintf("%d", a.MastercardSchemeData.Request.PointOfServiceData.TransactionStatus),
				SF8_TransactionSecurity:                       fmt.Sprintf("%d", a.MastercardSchemeData.Request.PointOfServiceData.TransactionSecurity),
				SF10_CardholderActivatedTerminalLevel:         fmt.Sprintf("%d", a.MastercardSchemeData.Request.PointOfServiceData.CardHolderActivatedTerminalLevel),
				SF11_CardDataTerminalInputCapabilityIndicator: fmt.Sprintf("%d", a.MastercardSchemeData.Request.PointOfServiceData.CardDataTerminalInputCapabilityIndicator),
				SF12_AuthorizationLifeCycle:                   a.MastercardSchemeData.Request.PointOfServiceData.AuthorizationLifeCycle,
				SF13_CountryCode:                              countrycode.Must(a.MastercardSchemeData.Request.PointOfServiceData.CountryCode).Mastercard.Numeric(),
				SF14_PostalCode:                               a.MastercardSchemeData.Request.PointOfServiceData.PostalCode,
			},
		},
	}
}
