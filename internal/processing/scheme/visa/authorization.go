package visa

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/pos"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/visa/base1"
	"gitlab.cmpayments.local/creditcard/platform/countrycode"
)

const (
	authorizationRequestMTI  = `0100`
	authorizationResponseMTI = `0110`
)

func AuthorizationSchemeData(a *entity.Authorization) {
	a.CardSchemeData.Request.ProcessingCode = entity.ProcessingCode{
		TransactionTypeCode: "00", // Purchase
		FromAccountTypeCode: "00", // Default Account
		ToAccountTypeCode:   "00", // Default Account
	}
	a.CardSchemeData.Request.POSEntryMode = posEntryMode(a.Recurring.Subsequent, a.CitMitIndicator)
	a.VisaSchemeData.Request.PosConditionCode = posConditionCode(a.Source, a.Recurring.Subsequent)
	a.VisaSchemeData.Request.AdditionalPOSInformation = additionalPOSInformation(a.Source, a.ThreeDSecure)
	a.VisaSchemeData.Request.PrivateUseFields = privateUseFields(a.Recurring)
}

func posEntryMode(subseqRecurring bool, indicator entity.CitMitIndicator) entity.POSEntryMode {
	if !subseqRecurring &&
		indicator.InitiatedBy == entity.Cardholder &&
		indicator.SubCategory == entity.CredentialOnFile {
		return entity.POSEntryMode{
			PanEntryMode: entity.PANEntryManual,      // Manual (key entry)
			PinEntryMode: entity.PINEntryUnspecified, // Unspecified or unknown
		}
	}

	if subseqRecurring &&
		indicator.InitiatedBy == entity.MITIndustryPractice &&
		indicator.SubCategory == entity.UnscheduledCredentialOnFile {
		return entity.POSEntryMode{
			PanEntryMode: entity.PANEntryCredentialOnFile, // Credential on file
			PinEntryMode: entity.PINEntryUnspecified,
		}
	}

	return entity.POSEntryMode{
		PanEntryMode: entity.PANEntryUnknown, // Unknown or terminal not used.
		PinEntryMode: entity.PINEntryUnspecified,
	}
}

func posConditionCode(s entity.Source, subseqRecurring bool) string {
	switch {
	case s == entity.Ecommerce:
		return "59"
	case s == entity.Moto:
		return "08"
	case subseqRecurring:
		return "08"
	default:
		return ""
	}
}

func additionalPOSInformation(s entity.Source, tds entity.ThreeDSecure) entity.AdditionalPOSInformation {
	return entity.AdditionalPOSInformation{
		TerminalType:                               "0", // unspecified
		TerminalEntryCapability:                    "1", // Terminal not used
		ChipConditionCode:                          "0",
		SpecialConditionIndicator:                  "0",                                  // Default value
		ChipTransactionIndicator:                   "0",                                  // Not applicable
		ChipCardAuthenticationReliabilityIndicator: "0",                                  // Fill for field 60.7 present, or subsequent subfields that are present.
		TypeOrLevelIndicator:                       mapElectricCommerceIndicator(s, tds), // Single transaction of a mail/phone order
		CardholderIDMethodIndicator:                "4",                                  // Mail/Telephone/Electronic Commerce
		AdditionalAuthorizationIndicators:          "0",                                  // Not applicable
	}
}

func mapElectricCommerceIndicator(s entity.Source, tds entity.ThreeDSecure) string {
	switch {
	case s == entity.Moto:
		return "01"
	case s == entity.Ecommerce && tds.EcommerceIndicator != 0:
		return fmt.Sprintf("%02d", tds.EcommerceIndicator)
	default:
		return "07"
	}
}

func privateUseFields(r entity.Recurring) entity.PrivateUseFields {
	switch {
	case r.Initial:
		return entity.PrivateUseFields{POSEnvironment: "C"}
	case r.Subsequent:
		return entity.PrivateUseFields{POSEnvironment: "R"}
	default:
		return entity.PrivateUseFields{}
	}
}

func mapInitiator(initiator entity.CitMitIndicator) string {
	var initiatedBy, subCategory string

	switch initiator.InitiatedBy {
	case entity.MITIndustryPractice:
		initiatedBy = "39"
	default:
		return ""
	}

	switch initiator.SubCategory {
	case entity.Resubmission:
		subCategory = "01"
	case entity.DelayedCharge:
		subCategory = "02"
	case entity.NoShow:
		subCategory = "04"
	default:
		return ""
	}

	return initiatedBy + subCategory
}

func authorizationResultFromMessage(msg Message) (entity.CardSchemeResponse, entity.VisaSchemeResponse) {
	return entity.CardSchemeResponse{
		Status: entity.AuthorizationStatusFromCardSchemeResponseCode(msg.Fields.F039_ResponseCode),
		ResponseCode: entity.ResponseCode{
			Value:       msg.Fields.F039_ResponseCode,
			Description: entity.ResponseDescriptionFromCode(msg.Fields.F039_ResponseCode),
		},
		AuthorizationIDResponse: msg.Fields.F038_AuthorizationIdenticationResponse,
		TraceId:                 fromTraceId(msg.Fields.F062_CustomPaymentServiceFields.SF2_TransactionIdentifier),
	}, entity.VisaSchemeResponse{TransactionId: msg.Fields.F062_CustomPaymentServiceFields.SF2_TransactionIdentifier}
}

func fromTraceId(traceId int) string {
	if traceId == 0 {
		return ""
	}
	return strconv.Itoa(traceId)
}

func MessageFromAuthorization(a entity.Authorization) (*Message, error) {
	cavv, err := cavvRequestData(a.ThreeDSecure)
	if err != nil {
		return nil, err
	}
	return &Message{
		Mti: iso8583.NewMti(authorizationRequestMTI),
		Fields: base1.Fields{
			F002_PrimaryAccountNumber:                   a.Card.Number,
			F003_ProcessingCode:                         a.CardSchemeData.Request.ProcessingCode.String(),
			F004_TransactionAmount:                      int64(a.Amount),
			F007_TransmissionDateTime:                   base1.F007FromTime(time.Now()),
			F011_SystemTraceAuditNumber:                 fmt.Sprintf("%06d", a.Stan),
			F012_LocalTransactionTime:                   a.LocalTransactionDateTime.Format(`150405`),
			F013_LocalTransactionDate:                   a.LocalTransactionDateTime.Format(`0102`),
			F014_ExpirationDate:                         a.Card.Expiry.String(),
			F018_MerchantType:                           a.CardAcceptor.CategoryCode,
			F019_AcquiringInstituteCountryCode:          countrycode.Must("NLD").Numeric(),
			F022_PosEntryMode:                           fmt.Sprintf("%s%s0", pos.PanEntryCode(a.CardSchemeData.Request.POSEntryMode.PanEntryMode), pos.PinEntryCode(a.CardSchemeData.Request.POSEntryMode.PinEntryMode)),
			F025_PosCondition:                           a.VisaSchemeData.Request.PosConditionCode,
			F032_AcquiringInstitutionIdentificationCode: entity.VisaInstitutionID,
			F034_ElectronicCommerceData:                 mapElectronicCommerceData(a.Exemption),
			F037_RetrievalReferenceNumber:               retrievalReferenceNumber(a.Stan),
			F042_CardAcceptorIdentificationCode:         fmt.Sprintf("%*s", -15, fmt.Sprintf("%s%s", a.Psp.Prefix, a.CardAcceptor.ID)),
			F043_CardAcceptorNameLocation: base1.F043_CardAcceptorNameLocation{
				SF1_CarAcceptorName:  a.CardAcceptor.Name,
				SF2_CardAcceptorCity: a.CardAcceptor.Address.City,
				SF3_CountryCode:      countrycode.Must(a.CardAcceptor.Address.CountryCode).Alpha2(),
			},
			F049_TransactionCurrencyCode: a.Currency.Numeric(),
			F060_AdditionalPointOfServiceInformation: base1.F060_AdditionalPOSInformation{
				B1: fmt.Sprintf("%s%s", a.VisaSchemeData.Request.AdditionalPOSInformation.TerminalType, a.VisaSchemeData.Request.AdditionalPOSInformation.TerminalEntryCapability),
				B2: fmt.Sprintf("%s%s", a.VisaSchemeData.Request.AdditionalPOSInformation.ChipConditionCode, a.VisaSchemeData.Request.AdditionalPOSInformation.SpecialConditionIndicator),
				B3: "00",
				B4: fmt.Sprintf("%s%s", a.VisaSchemeData.Request.AdditionalPOSInformation.ChipTransactionIndicator, a.VisaSchemeData.Request.AdditionalPOSInformation.ChipCardAuthenticationReliabilityIndicator),
				B5: a.VisaSchemeData.Request.AdditionalPOSInformation.TypeOrLevelIndicator,
				B6: fmt.Sprintf("%s%s", a.VisaSchemeData.Request.AdditionalPOSInformation.CardholderIDMethodIndicator, a.VisaSchemeData.Request.AdditionalPOSInformation.AdditionalAuthorizationIndicators),
			},
			F062_CustomPaymentServiceFields: base1.F062_CustomPaymentService{
				SF1_AuthorizationCharacteristicsIndicator: "Y", //TODO Different for Recurring
				SF2_TransactionIdentifier:                 toTraceId(a.Recurring),
			},
			F063_NetworkData: base1.F063_NetworkData{
				SF1_NetworkID:         mapNetworkId(),
				SF3_MessageReasonCode: mapInitiator(a.CitMitIndicator),
			},
			F126_PrivateUseFields: base1.F126_PrivateUseFields{
				SF9_CAVVData:                      cavv,
				SF10_CVV2AuthorizationRequestData: cvvRequestData(a.Card.Cvv, a.Recurring.Subsequent),
				SF13_POSEnvironment:               a.VisaSchemeData.Request.PrivateUseFields.POSEnvironment,
			},
		},
	}, nil
}

func toTraceId(r entity.Recurring) int {
	tId, err := strconv.Atoi(r.TraceID)
	if err != nil {
		return 0
	}
	return tId
}

func is3dSecure(source entity.Source, secure entity.ThreeDSecure) string {
	if source == entity.Moto {
		return "01"
	}

	if secure == (entity.ThreeDSecure{}) {
		// Non-authenticated security transaction: Use to identify an electronic commerce
		// transaction that uses data encryption for security however,
		// cardholder authentication is not performed using a Visa approved protocol,
		// such as 3-D Secure. Reference: visanet-authorization-only-online-messages-technical-specifications.pdf page 349
		return "07"
	}

	return fmt.Sprintf("%02d", secure.EcommerceIndicator)
}

func mapNetworkId() string {
	return "0002"
}

func cvvRequestData(cvv string, ssr bool) string {
	if ssr {
		return ""
	}

	if cvv != "" {
		// position 1 Presence Indicator (0 = not present,  1 = present)
		// position 2 Response Type (0 or 1, if 1 then we receive the CVV2 auth result in F44.10)
		// position 3 - 6 CVV, if CVV is 3 char's (MC and Visa) then pos3 must contain a space
		if len(cvv) == 4 {
			return "11" + cvv
		}
		return "11" + " " + cvv
	}

	return "01    "
}

var CavvErrorNumeric = errors.New("CAVV is not a numeric value")

func cavvRequestData(secure entity.ThreeDSecure) (string, error) {
	if secure.NotSet() {
		return "", nil
	}

	if secure.AuthenticationVerificationValue == "" {
		return "", nil
	}

	decoded, err := base64.StdEncoding.DecodeString(secure.AuthenticationVerificationValue)
	if err != nil {
		return "", err
	}

	if !regexp.MustCompile(`\d`).Match(decoded) {
		return "", CavvErrorNumeric
	}

	return string(decoded), nil
}

func retrievalReferenceNumber(stan int) string {
	return time.Now().Format(`0600215`)[1:] + strconv.Itoa(stan)
}

func mapElectronicCommerceData(exemption entity.ExemptionType) base1.F034_ElectronicCommerceData {
	var ecd base1.F034_ElectronicCommerceData

	if exemption == "" {
		return base1.F034_ElectronicCommerceData{}
	}

	// TODO: map 3D Secure Protocol (HEX01_AuthenticationData => T86_3DSecureProtocolVersionNumber)
	//ecd.HEX01_AuthenticationData.T86_3DSecureProtocolVersionNumber = "UNKNOWN"

	switch exemption {
	case entity.MerchantInitiatedExemption:
		ecd.HEX02_AcceptanceEnvironmentAdditionalData.T80_InitiatingPartyIndicator = "1"
	case entity.LowValueExemption:
		ecd.HEX4A_StrongConsumerAuthentication.T87_LowValueExemptionIndicator = "1"
	case entity.SecureCorporateExemption:
		ecd.HEX4A_StrongConsumerAuthentication.T88_SecureCorporatePaymentIndicator = "1"
	case entity.TransactionRiskAnalysisExemption:
		ecd.HEX4A_StrongConsumerAuthentication.T89_TransactionRiskAnalysisExemptionIndicator = "1"
	case entity.SCADelegationExemption:
		ecd.HEX4A_StrongConsumerAuthentication.T8A_DelegatedAuthenticationIndicator = "1"
	}

	return ecd
}
