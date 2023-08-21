package cis

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
)

//nolint:lll
type DE48_AdditionalData struct {
	TransactionCategoryCode                         string                                  //	This is not a proper subfield / subelement, but the prefix
	SE15_AuthorizationSystemAdviceDateTime          string                                  `iso8583:"15=n-10"`
	SE20_CardholderVerificationMethod               string                                  `iso8583:"20=a-1"`
	SE22_MultiPurposeMerchantIndicator              *DE48_SE22MultiPurposeMerchantIndicator `iso8583:"22=ans...21"`
	SE23_PaymentInitiationChannel                   string                                  `iso8583:"23=an-2"`
	SE26_WalletProgramData                          string                                  `iso8583:"26=an-3"`
	SE32_MastercardAssignedId                       string                                  `iso8583:"32=an-6"`
	SE33_PanMappingFileInformation                  *DE48_SE33_PANMappingFileInformation    `iso8583:"33=ans...93"`
	SE36_VisaMvv                                    string                                  `iso8583:"36=an-14"`
	SE37_AdditionalMerchantData                     *DE48_SE37_AdditionalMerchantData       `iso8583:"37=ans...49"`
	SE42_ElectronicCommerceIndicators               *DE48_SE42_ElectronicCommerceIndicators `iso8583:"42=n..19"`
	SE43_UniversalCardholderAuthenticationField     string                                  `iso8583:"43=ans...32"`
	SE44_TransactionIdentifier                      string                                  `iso8583:"44=b-20"`
	SE51_MerchantOnBehalfServices                   string                                  `iso8583:"51=ans..99"`
	SE52_TransactionIntegretyClass                  string                                  `iso8583:"52=an..2"`
	SE55_MerchantFraudScoringData                   string                                  `iso8583:"55=an..32"`
	SE57_SecurityServicesAdditionalDataForAcquirers string                                  `iso8583:"57=an-99"`
	SE61_ExtendedConditionCodes                     *DE48_SE61_ExtendedConditionCodes       `iso8583:"61=n-5, autofill"`
	SE63_TraceId                                    *DE48_SE63_TraceId                      `iso8583:"63=ans-15"`
	SE64_TransitProgram                             string                                  `iso8583:"64=n-4"`
	SE66_AuthenticationData                         *DE48_SE66_AuthenticationData           `iso8583:"66=ans..45, omitempty"`
	SE67_MoneySendInformation                       string                                  `iso8583:"67=ans..99"`
	SE76_ElectronicAcceptanceIndicator              string                                  `iso8583:"76=a-1"`
	SE77_TransactionTypeIdentifier                  string                                  `iso8583:"77=an-3"`
	SE80_PinServiceCode                             string                                  `iso8583:"80=a-2"`
	SE82_AddressVerificationServiceRequest          string                                  `iso8583:"82=n-2"`
	SE83_AddressVerificationServiceResponse         string                                  `iso8583:"83=a-1"`
	SE84_MerchantAdviceCode                         string                                  `iso8583:"84=an-2"`
	SE86_RelationshipParticipantIndicator           string                                  `iso8583:"86=a-1"`
	SE87_Cvv2Response                               string                                  `iso8583:"87=a-1"`
	SE90_LodgingAndAutoRentalIndicator              string                                  `iso8583:"90=a-1"`
	SE92_CardholderVerificationCode                 string                                  `iso8583:"92=n-3"`
	SE95_PromotionCode                              string                                  `iso8583:"95=an-6"`
}

// MarshalIso8583 marshals the DE48_AdditionalData subelements
// The first 3 subfields are delimited by a backslash
//
//nolint:wrapcheck
func (de *DE48_AdditionalData) MarshalIso8583() ([]byte, error) {
	// Read field definitions from struct
	definitions, err := iso8583.StructDefinitions(de)
	if err != nil {
		return nil, err
	}

	values := reflect.Indirect(reflect.ValueOf(de))

	// Marshal each field
	var buf bytes.Buffer

	if _, err := buf.WriteString(de.TransactionCategoryCode); err != nil {
		return nil, errors.New(`cannot write TCC for de48: ` + err.Error())
	}

	for _, se := range definitions {
		field := values.Field(se.Field)

		// Skip empty values
		if field.IsZero() {
			continue
		}

		// Each subelement consists of three components
		// - Tag Field
		// - Length Field
		// - Data Field

		// The first two bytes contain the subelement tag ID, a numeric value in the
		// range 01–99, to identify uniquely the subelement.
		if _, err := fmt.Fprintf(&buf, "%02d", se.Number); err != nil {
			return nil, fmt.Errorf(
				`cannot write subelement id and length for de48/se%s: %w`,
				strconv.Itoa(se.Number), err,
			)
		}

		// The length indicator always 2 bytes for these subelements, regardless of their type
		se.LengthIndicator = 2

		// The remainder of the subelement consists of the actual data
		if err := iso8583.Encode(field.Interface(), se, &buf); err != nil {
			return nil, fmt.Errorf(`cannot encode value for de48/se%s: %w`, strconv.Itoa(se.Number), err)
		}
	}

	return buf.Bytes(), nil
}

//nolint:wrapcheck
func (de *DE48_AdditionalData) UnmarshalIso8583(d []byte) error {
	r := bytes.NewBuffer(d)

	tcc, err := r.ReadByte()
	if err != nil {
		return err
	}

	de.TransactionCategoryCode = string(tcc)

	// Read field definitions from struct
	definitions, err := iso8583.StructDefinitions(de)
	if err != nil {
		return err
	}

	values := reflect.Indirect(reflect.ValueOf(de))

	for {
		// The first two bytes of each SE must contain the number ID, a numeric value in the range 0001–9999,
		// to uniquely identify the subelement.
		buf := make([]byte, 2)
		if _, err := io.ReadFull(r, buf); err != nil {
			if errors.Is(err, io.EOF) {
				// Continue with next buffer
				break
			}
			// Wrap the error if its not EOF
			return fmt.Errorf("could not read SE number; %w", err)
		}

		number, err := strconv.Atoi(string(buf))
		if err != nil {
			return fmt.Errorf("invalid SE number; %w", err)
		}

		element, ok := definitions[number]
		if !ok {
			return fmt.Errorf("missing definition for SE %d", number)
		}

		// The next two bytes of each SE must contain the subelement length, a numeric value in the
		// range 01–99, to specify the total length (in bytes) of the SE
		element.LengthIndicator = 2

		// Get the variable we must decode into
		field := iso8583.StructField(values, element.Field)

		// Pass a pointer to the variable where the data must be decoded into
		if err := iso8583.Decode(field.Addr().Interface(), element, r, iso8583.LengthEncodingAscii); err != nil {
			return fmt.Errorf(`cannot decode value for de48/se%s: %w`, strconv.Itoa(element.Number), err)
		}
	}

	return nil
}

type DE48_SE22MultiPurposeMerchantIndicator struct {
	SF1_LowRisk                       string `iso8583:"1=an..2, minlength=2,omitempty,justify=right"`
	_                                 string `iso8583:"2=an..1, minlength=0,omitempty"`
	_                                 string `iso8583:"3=an..1, minlength=0,omitempty"`
	_                                 string `iso8583:"4=an..1, minlength=0,omitempty"`
	SF5_InitiatedTransactionIndicator string `iso8583:"5=an..4, minlength=4,omitempty"`
}

func NewDE48_SE22MultiPurposeMerchantIndicator(lowRisk string, indicator string) *DE48_SE22MultiPurposeMerchantIndicator {
	if lowRisk == "" && indicator == "" {
		return nil
	}

	return &DE48_SE22MultiPurposeMerchantIndicator{
		SF1_LowRisk:                       lowRisk,
		SF5_InitiatedTransactionIndicator: indicator,
	}
}

func (se *DE48_SE22MultiPurposeMerchantIndicator) MarshalIso8583() ([]byte, error) {
	return MarshallDE48SEWithPositionIndicators(se)
}

func (se *DE48_SE22MultiPurposeMerchantIndicator) UnmarshalIso8583(d []byte) error {
	return UnmarshallDE48SEWithPositionIndicators(se, d)
}

type DE48_SE33_PANMappingFileInformation struct {
	_                    string `iso8583:"1=an..1, minlength=0, omitempty"`
	_                    string `iso8583:"2=n..19, minlength=0, omitempty"`
	_                    string `iso8583:"3=n..4, minlength=0, omitempty"`
	_                    string `iso8583:"4=an..3, minlength=0, omitempty"`
	_                    string `iso8583:"5=n..2, minlength=0, omitempty"`
	SF6_TokenRequestorID string `iso8583:"6=n..11, minlength=0, omitempty"`
	_                    string `iso8583:"7=n..19, minlength=0, omitempty"`
	_                    string `iso8583:"8=an..2, minlength=0, omitempty"`
}

func (se *DE48_SE33_PANMappingFileInformation) MarshalIso8583() ([]byte, error) {
	return MarshallDE48SEWithPositionIndicators(se)
}

func (se *DE48_SE33_PANMappingFileInformation) UnmarshalIso8583(d []byte) error {
	return UnmarshallDE48SEWithPositionIndicators(se, d)
}

type DE48_SE37_AdditionalMerchantData struct {
	SF1_PaymentFacilitatorId           string `iso8583:"1=n-11, minlength=0, omitempty"`
	SF2_IndependantSalesOrganisationId string `iso8583:"2=n-11, minlength=0, omitempty"`
	SF3_SubMerchantId                  string `iso8583:"3=ans-15, minlength=0, omitempty"`
	SF4_MerchantCountryOfOrigin        string `iso8583:"4=n-3, minlength=0, omitempty"`
}

type DE48_SE42_ElectronicCommerceIndicators struct {
	SF1_SecurityLevelIndicatorAndUCAFCollectionIndicator         string `iso8583:"1=n..3"`
	SF2_OriginalSecurityLevelIndicatorAndUCAFCollectionIndicator string `iso8583:"2=n..3, omitempty"`
	SF3_ReasonForUCAFDowngrade                                   string `iso8583:"3=n..1, omitempty"`
}

func NewDE48_SE42_ElectronicCommerceIndicators(i entity.SLI) *DE48_SE42_ElectronicCommerceIndicators {
	return &DE48_SE42_ElectronicCommerceIndicators{
		SF1_SecurityLevelIndicatorAndUCAFCollectionIndicator: fmt.Sprintf("%d%d%d", i.SecurityProtocol, i.CardholderAuthentication, i.UCAFCollectionIndicator),
	}
}

func (se *DE48_SE42_ElectronicCommerceIndicators) MarshalIso8583() ([]byte, error) {
	return MarshallDE48SEWithPositionIndicators(se)
}

func (se *DE48_SE42_ElectronicCommerceIndicators) UnmarshalIso8583(d []byte) error {
	r := bytes.NewBuffer(d)

	// Read field definitions from struct
	definitions, err := iso8583.StructDefinitions(se)
	if err != nil {
		return fmt.Errorf("struct definition not found: %w", err)
	}

	values := reflect.Indirect(reflect.ValueOf(se))

	// we have 3 subfields
	for number := 1; number <= 3; number++ {
		// The first two bytes of each SE must contain the number ID, a numeric value in the range 0001–9999,
		// to uniquely identify the subelement.
		buf := make([]byte, 2)
		if _, err := io.ReadFull(r, buf); err != nil {
			if errors.Is(err, io.EOF) {
				// Continue with next buffer
				break
			}
			// Wrap the error if its not EOF
			return fmt.Errorf("could not read SE number; %w", err)
		}

		element, ok := definitions[number]
		if !ok {
			return iso8583.NewElementError(number, fmt.Errorf("missing definition")) //nolint:wrapcheck
		}

		// Get the variable we must decode into
		field := iso8583.StructField(values, element.Field)

		// Pass a pointer to the variable where the data must be decoded into
		if err := iso8583.Decode(field.Addr().Interface(), element, r, iso8583.LengthEncodingAscii); err != nil {
			if errors.Is(err, io.EOF) {
				// Not all subfields have to be present (ex PDS 158)
				break
			}

			return iso8583.NewElementError(number, err) //nolint:wrapcheck
		}
	}

	return nil
}

type DE48_SE66_AuthenticationData struct {
	SF1_ProgramProtocol              string `iso8583:"1=an..1, minlength=0, omitempty"`
	SF2_DirectoryServerTransactionId string `iso8583:"2=ans..36, minlength=0, omitempty"`
}

func NewDE48_SE66_AuthenticationData(data entity.AuthenticationData) *DE48_SE66_AuthenticationData {
	if (data == entity.AuthenticationData{}) {
		return nil
	}

	return &DE48_SE66_AuthenticationData{
		SF1_ProgramProtocol:              data.ProgramProtocol,
		SF2_DirectoryServerTransactionId: data.DirectoryServerTransactionID,
	}
}

func (se *DE48_SE66_AuthenticationData) MarshalIso8583() ([]byte, error) {
	return MarshallDE48SEWithPositionIndicators(se)
}

func (se *DE48_SE66_AuthenticationData) UnmarshalIso8583(d []byte) error {
	return UnmarshallDE48SEWithPositionIndicators(se, d)
}

type DE48_SE61_ExtendedConditionCodes struct {
	SF1_PartialApprovalTerminalSupportIndicator    int `iso8583:"1=n-1, autofill"`
	SF2_PurchaseAmountOnlyTerminalSupportIndicator int `iso8583:"2=n-1, autofill"`
	SF3_RealTimeSubstantiationIndicattor           int `iso8583:"3=n-1, autofill"`
	SF4_MerchantTransactionFroudScoringIndicator   int `iso8583:"4=n-1, autofill"`
	SF5_FinalAuthorizationIndicator                int `iso8583:"5=n-1, autofill"`
}

type DE48_SE63_TraceId struct {
	SF1_NetworkData    string `iso8583:"1=ans-9,omitempty"`
	SF2_DateSettlement string `iso8583:"2=ans-6,omitempty,justify=right"`
}

func NewDE48_SE63_TraceId(traceID entity.MTraceID) *DE48_SE63_TraceId {
	if (traceID == entity.MTraceID{}) {
		return nil
	}

	return &DE48_SE63_TraceId{
		SF1_NetworkData:    traceID.NetworkData(),
		SF2_DateSettlement: traceID.NetworkReportingDate,
	}
}

//nolint:wrapcheck
func MarshallDE48SEWithPositionIndicators(se interface{}) ([]byte, error) {
	// Read field definitions from struct
	definitions, err := iso8583.StructDefinitions(se)
	if err != nil {
		return nil, err
	}

	values := reflect.Indirect(reflect.ValueOf(se))

	// Marshal each field
	var buf bytes.Buffer

	for i, elem := range definitions {
		field := values.Field(elem.Field)

		if !field.CanInterface() {
			continue
		}

		if elem.OmitEmpty && field.IsZero() {
			continue
		}

		var ebuf bytes.Buffer

		if err := iso8583.Encode(field.Interface(), elem, &ebuf); err != nil {
			return nil, iso8583.NewElementError(i, err)
		}

		if ebuf.Len() > 0 {
			buf.WriteString(fmt.Sprintf("%0*d", 2, i))
			buf.Write(ebuf.Bytes())
		}
	}

	return buf.Bytes(), nil
}

//nolint:wrapcheck
func UnmarshallDE48SEWithPositionIndicators(se interface{}, d []byte) error {
	buf := bytes.NewBuffer(d)

	// Read field definitions from struct
	definitions, err := iso8583.StructDefinitions(se)
	if err != nil {
		return err
	}

	values := reflect.Indirect(reflect.ValueOf(se))
	// we know we have 6 subfields
	for number, element := range definitions {
		// Get the variable we must decode into
		field := iso8583.StructField(values, element.Field)

		if field.IsZero() {
			continue
		}

		// Pass a pointer to the variable where the data must be decoded into
		if err := iso8583.Decode(field.Addr().Interface(), element, buf, iso8583.LengthEncodingAscii); err != nil {
			return iso8583.NewElementError(number, err)
		}
	}

	return nil
}
