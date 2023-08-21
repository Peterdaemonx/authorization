package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/pos"
	mapping "gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"

	sql "gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/mastercard"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"

	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

type AuthorizationRepository struct {
	client       *spanner.Client
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewAuthorizationRepository(
	client *spanner.Client,
	readTimeout time.Duration,
	writeTimeout time.Duration) *AuthorizationRepository {
	return &AuthorizationRepository{
		client:       client,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (ar AuthorizationRepository) CreateAuthorization(ctx context.Context, a entity.Authorization) error {
	stmt := spanner.Statement{
		SQL: `INSERT INTO authorizations (
					authorization_id, masked_pan, pan_token_id,
                    amount, currency, localdatetime,
                    source, customer_reference, psp_id,
                    card_acceptor_id, card_acceptor_name, card_acceptor_postal_code,
                    card_acceptor_city, card_acceptor_country, card_acceptor_category_code,
                    is_initial_recurring, initial_trace_id, exemption,
                    threeds_version, threeds_directory_server_transaction_id, threeds_original_ecommerce_indicator,
                    created_at, status, card_scheme,
                    card_issuer_id, card_issuer_name, card_issuer_countrycode,
                    accountholder_authentication_value, transaction_initiated_by, transaction_subcategory,
                    terminal_id, card_holder_activated_terminal_level, terminal_capability, card_sequence,
                    card_holder_verification_method
                ) VALUES (
                    @authorization_id, @masked_pan, @pan_token_id,
                    @amount, @currency, @localdatetime,
                    @source, @customer_reference, @psp_id,
                    @card_acceptor_id, @card_acceptor_name, @card_acceptor_postal_code,
                    @card_acceptor_city, @card_acceptor_country, @card_acceptor_category_code,
                    @is_initial_recurring, @initial_trace_id, @exemption,
                    @threeds_version, @threeds_directory_server_transaction_id, @threeds_original_ecommerce_indicator,
                    @created_at, 'new', @card_scheme,
                    @card_issuer_id, @card_issuer_name, @card_issuer_countrycode,
                    @accountholder_authentication_value, @transaction_initiated_by, @transaction_subcategory,
                    @terminal_id, @card_holder_activated_terminal_level, @terminal_capability, @card_sequence,
                    @card_holder_verification_method
                )`,
		Params: mapCreateAuthParams(a),
	}

	_, err := ar.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, ar.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, stmt)

		return err
	})

	return err
}

func mapCreateAuthParams(a entity.Authorization) map[string]interface{} {
	return map[string]interface{}{
		"authorization_id":     a.ID.String(),
		"masked_pan":           a.Card.MaskedPan,
		"pan_token_id":         a.Card.PanTokenID,
		"card_scheme":          a.Card.Info.Scheme,
		"amount":               a.Amount,
		"currency":             a.Currency.Alpha3(),
		"localdatetime":        a.LocalTransactionDateTime,
		"source":               string(a.Source),
		"is_initial_recurring": a.Recurring.Initial,
		"threeds_version":      sql.NewNullString(a.ThreeDSecure.Version),
		"created_at":           time.Now(),
		"customer_reference":   sql.NewNullString(a.CustomerReference),
		"initial_trace_id":     sql.NewNullString(a.Recurring.TraceID),
		"threeds_directory_server_transaction_id": sql.NewNullString(a.ThreeDSecure.DirectoryServerID),
		"psp_id":                               a.Psp.ID.String(),
		"card_issuer_id":                       sql.NewNullString(a.Card.Info.IssuerID),
		"card_issuer_name":                     sql.NewNullString(a.Card.Info.IssuerName),
		"card_issuer_countrycode":              sql.NewNullString(a.Card.Info.IssuerCountryCode),
		"card_acceptor_name":                   a.CardAcceptor.Name,
		"card_acceptor_city":                   a.CardAcceptor.Address.City,
		"card_acceptor_country":                a.CardAcceptor.Address.CountryCode,
		"card_acceptor_id":                     a.CardAcceptor.ID,
		"card_acceptor_postal_code":            a.CardAcceptor.Address.PostalCode,
		"card_acceptor_category_code":          a.CardAcceptor.CategoryCode,
		"threeds_original_ecommerce_indicator": a.ThreeDSecure.EcommerceIndicator,
		"exemption":                            sql.NewNullString(string(a.Exemption)),
		"transaction_initiated_by":             sql.NewNullString(string(a.CitMitIndicator.InitiatedBy)),
		"transaction_subcategory":              sql.NewNullString(string(a.CitMitIndicator.SubCategory)),
		"accountholder_authentication_value":   a.ThreeDSecure.AuthenticationVerificationValue,
		"terminal_id":                          a.Terminal.TerminalId,
		"card_holder_activated_terminal_level": a.Terminal.TerminalLevel,
		"terminal_capability":                  a.Terminal.TerminalCapability,
		"card_sequence":                        a.Card.SequenceNumber,
		"card_holder_verification_method":      a.CardSchemeData.Request.CardHolderVerificationMethod,
	}
}

func (ar AuthorizationRepository) CreateMastercardAuthorization(ctx context.Context, a entity.Authorization) error {
	_, err := ar.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		statement := spanner.Statement{
			SQL: `
				INSERT INTO mastercard_authorizations (
					authorization_id, network_reporting_date, financial_network_code,
					reference,
					pos_pin_capture_code, pin_service_code, security_protocol,
				    cardholder_authentication, ucaf_collection_indicator, banknet_reference_number,
					created_at, authorization_type, 
					terminal_attendance, terminal_location, card_holder_presence,
					card_presence, card_capture_capabilities, transaction_status,
					transaction_security, card_holder_activated_terminal_level, card_data_terminal_input_capability_indicator,
					authorization_life_cycle, country_code, postal_code,
				    reason_ucaf_downgrade, card_program_id, card_product_id
				) VALUES (
					@authorization_id, @network_reporting_date, @financial_network_code,
					@reference,
					@pos_pin_capture_code, @pin_service_code, @security_protocol,
				    @cardholder_authentication, @ucaf_collection_indicator, @banknet_reference_number,
				    @created_at, @authorization_type, 
				    @terminal_attendance, @terminal_location, @card_holder_presence,
					@card_presence, @card_capture_capabilities, @transaction_status,
					@transaction_security, @card_holder_activated_terminal_level, @card_data_terminal_input_capability_indicator,
					@authorization_life_cycle, @country_code, @postal_code,
				    @reason_ucaf_downgrade, @card_program_id, @card_product_id
				)`,
			Params: mapMastercardAuthorizationParams(a),
		}

		ctx, cancel := context.WithTimeout(ctx, ar.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, statement)

		return err
	})

	return err
}

func (ar AuthorizationRepository) UpdateAuthorizationResponse(ctx context.Context, a entity.Authorization) error {
	stmt := spanner.Statement{
		SQL: `UPDATE authorizations
				SET status = @status,
				    system_trace_audit_number = @system_trace_audit_number,
					updated_at = @updated_at,
					authorization_id_response = @authorization_id_response,
					retrieval_reference_number = @retrieval_reference_number,
					response_code = @response_code,
					transmitted_at = @transmitted_at,
					cardholder_transaction_type_code = @card_holder_transaction_type_code,
					cardholder_from_account_type_code = @cardholder_from_account_type_code,
					cardholder_to_account_type_code = @cardholder_to_account_type_code,
					point_of_service_pan_entry_mode = @point_of_service_pan_entry_mode,
					point_of_service_pin_entry_mode = @point_of_service_pin_entry_mode
				WHERE authorization_id = @authorization_id`,
		Params: mapUpdateAuthorizationResponseParams(a),
	}

	_, err := ar.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, ar.writeTimeout)
		defer cancel()

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to update authorization: %w", err)
		}
		if rowCount != 1 {
			return fmt.Errorf("no record found with ID: %s", a.ID.String())
		}

		return err
	})

	return err
}

func mapUpdateAuthorizationResponseParams(a entity.Authorization) map[string]interface{} {
	return map[string]interface{}{
		"authorization_id":                  a.ID.String(),
		"status":                            a.CardSchemeData.Response.Status.String(),
		"system_trace_audit_number":         a.Stan,
		"updated_at":                        time.Now(),
		"response_code":                     sql.NewNullString(a.CardSchemeData.Response.ResponseCode.Value),
		"authorization_id_response":         sql.NewNullString(a.CardSchemeData.Response.AuthorizationIDResponse),
		"retrieval_reference_number":        sql.NewNullString(a.CardSchemeData.Request.RetrievalReferenceNumber),
		"transmitted_at":                    a.ProcessingDate,
		"card_holder_transaction_type_code": a.CardSchemeData.Request.ProcessingCode.TransactionTypeCode,
		"cardholder_from_account_type_code": a.CardSchemeData.Request.ProcessingCode.FromAccountTypeCode,
		"cardholder_to_account_type_code":   a.CardSchemeData.Request.ProcessingCode.ToAccountTypeCode,
		"point_of_service_pan_entry_mode":   pos.PanEntryCode(a.CardSchemeData.Request.POSEntryMode.PanEntryMode),
		"point_of_service_pin_entry_mode":   pos.PinEntryCode(a.CardSchemeData.Request.POSEntryMode.PinEntryMode),
	}
}

func (ar AuthorizationRepository) GetAllAuthorizations(ctx context.Context, pspID uuid.UUID, filters entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Authorization, error) {
	var sortColumn string
	var dateRangeStart, dateRangeEnd time.Time

	// we're using different terms and styling on the FO, API endpoints and in the DB so we need to transform the sort to an existing column.
	switch filters.SortColumn() {
	case "createdAt":
		sortColumn = "created_at"
	case "processingDate":
		sortColumn = "transmitted_at"
	case "exemption":
		sortColumn = "exemption"
	default:
		sortColumn = filters.SortColumn()
	}

	// these parameters are already validated on the ports layer. We're sure these are of type time.Time
	dateRangeStart, dateRangeEnd = determineDateRange(params["startDate"].(time.Time), params["endDate"].(time.Time))

	// SQL injection is caught by valdation of the possible sort input and in the data.filters.SortColumn function
	stmt := spanner.Statement{
		SQL: fmt.Sprintf(`
			SELECT 
				a.authorization_id, a.psp_id, a.card_scheme,
				a.masked_pan, a.amount, a.currency,
				a.localdatetime, a.status, a.source,
				a.is_initial_recurring, a.customer_reference,
				a.exemption, a.threeds_version, a.threeds_original_ecommerce_indicator,
				ma.ucaf_collection_indicator, a.threeds_directory_server_transaction_id, a.card_acceptor_name,
				a.card_acceptor_city, a.card_acceptor_country, a.card_acceptor_category_code,
				a.card_acceptor_id, a.card_acceptor_postal_code, a.transmitted_at,
				a.response_code, a.created_at,
				a.updated_at, a.initial_trace_id, a.transaction_initiated_by, a.transaction_subcategory,
				ma.authorization_type, ma.financial_network_code, ma.banknet_reference_number, ma.network_reporting_date,
				ma.reason_ucaf_downgrade, ma.card_program_id, ma.card_product_id
			FROM authorizations a
			JOIN mastercard_authorizations ma on a.authorization_id = ma.authorization_id
			WHERE a.psp_id = @psp_id
			AND (a.masked_pan LIKE @masked_pan OR @masked_pan = '')
			AND (a.status = @status OR @status = '')
			AND (a.amount = @amount OR @amount = -1)
			AND (a.exemption = @exemption OR @exemption = '')
			AND (a.response_code = @response_code OR @response_code = '')
			AND (a.customer_reference = @customer_reference OR @customer_reference = '')
			AND (a.transmitted_at = @transmitted_at OR @transmitted_at = '0001-01-01T00:00:00Z')
			AND (a.created_at BETWEEN @created_after AND @created_before)
			ORDER BY %s %s
			LIMIT %v OFFSET %v`,
			sortColumn,
			filters.SortDirection(),
			filters.Limit()+1, // We fetch one extra record to see if we have fetched the last data set
			filters.Offset(),
		),
		Params: map[string]interface{}{
			"psp_id":             pspID.String(),
			"masked_pan":         fmt.Sprint("%", params["pan"].(string)),
			"status":             params["status"].(string),
			"amount":             params["amount"].(int),
			"exemption":          params["exemption"].(string),
			"transmitted_at":     params["processingDate"].(time.Time),
			"response_code":      params["responseCode"].(string),
			"customer_reference": params["reference"].(string),
			"created_before":     dateRangeEnd.Format("2006-01-02 15:04:05"),
			"created_after":      dateRangeStart.Format("2006-01-02 15:04:05"),
		},
	}

	var (
		authorizations []entity.Authorization
		fetchedRecords = 0
	)

	ctx, cancel := context.WithTimeout(ctx, ar.readTimeout)
	defer cancel()

	iter := ar.client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return entity.Metadata{}, nil, fmt.Errorf("failed to iterate over rows: %w", err)
		}

		var mar MastercardAuthorizationRecord
		if err = row.ToStruct(&mar); err != nil {
			return entity.Metadata{}, nil, fmt.Errorf("failed to parse row into struct: %w", err)
		}

		authorizations = append(authorizations, mapRowToAuthorizationEntity(mar))

		fetchedRecords++
	}

	// Remove the last record if we fetch more records than set on the page size.
	if fetchedRecords > filters.Limit() {
		authorizations = authorizations[:filters.Limit()]
	}

	metadata := entity.CalculateMetadata(fetchedRecords, filters.Page, filters.PageSize)

	return metadata, authorizations, nil
}

func (ar AuthorizationRepository) UpdateAuthorizationStatus(ctx context.Context, authorizationID uuid.UUID, status entity.Status) error {
	stmt := spanner.NewStatement(
		`UPDATE authorizations
				SET status = @status
				WHERE authorization_id = @authorization_id`)
	stmt.Params["status"] = status
	stmt.Params["authorization_id"] = authorizationID.String()

	_, err := ar.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, ar.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to update authorization: %w", err)
		}

		return err
	})
	return err
}

func (ar AuthorizationRepository) GetAuthorization(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error) {
	stmt := spanner.NewStatement(`
		SELECT authorization_id, status, masked_pan, pan_token_id, 
		       card_scheme, amount, currency, localdatetime, source, 
		       is_initial_recurring, threeds_version, created_at, 
		       customer_reference, updated_at, initial_trace_id, response_code, 
		       system_trace_audit_number, threeds_directory_server_transaction_id, 
		       psp_id, threeds_original_ecommerce_indicator, transmitted_at, 
		       card_issuer_id, card_issuer_name, card_issuer_countrycode, 
		       cardholder_transaction_type_code, cardholder_from_account_type_code, 
		       cardholder_to_account_type_code, card_acceptor_name, card_acceptor_city, 
		       card_acceptor_country, card_acceptor_id, card_acceptor_postal_code, 
		       card_acceptor_category_code, exemption, accountholder_authentication_value,
		       transaction_initiated_by, transaction_subcategory
		FROM authorizations AS a
		WHERE a.authorization_id = @authorizationID
    `)
	stmt.Params["authorizationID"] = authorizationID.String()

	iter := ar.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return entity.Authorization{}, entity.ErrRecordNotFound
		}
		return entity.Authorization{}, err
	}

	var a AuthorizationRecord
	if err = row.ToStruct(&a); err != nil {
		return entity.Authorization{}, err
	}

	return mapAuthorizationEntity(a), err
}

func mapAuthorizationEntity(a AuthorizationRecord) entity.Authorization {
	return entity.Authorization{
		ID:                       uuid.MustParse(a.ID),
		LogID:                    uuid.UUID{},
		Amount:                   int(a.Amount),
		Currency:                 currencycode.Must(a.Currency),
		CustomerReference:        a.CustomerReference.StringVal,
		Source:                   entity.Source(a.Source),
		LocalTransactionDateTime: data.LocalTransactionDateTime(a.LocalDateTime),
		Status:                   entity.Status(a.Status),
		Stan:                     int(a.Stan.Int64),
		ProcessingDate:           a.TransmissionDate.Time,
		CreatedAt:                a.CreatedAt,
		Recurring: entity.Recurring{
			Initial:    a.IsInitialRecurring.Bool,
			Subsequent: !a.IsInitialRecurring.Bool,
			TraceID:    a.InitialTraceID.StringVal,
		},
		Card: entity.Card{
			MaskedPan:  a.MaskedPan,
			PanTokenID: a.PanTokenID,
			Info: cardinfo.Range{
				Scheme:            a.CardScheme,
				IssuerID:          a.CardIssuerID.StringVal,
				IssuerName:        a.CardIssuerName.StringVal,
				IssuerCountryCode: a.CardIssuerCountryCode.StringVal,
			},
		},
		CardAcceptor: entity.CardAcceptor{
			CategoryCode: a.CategoryCode,
			ID:           a.AcceptorID,
			Name:         a.AcceptorName,
			Address: entity.CardAcceptorAddress{
				PostalCode:  a.AcceptorPostalCode.StringVal,
				City:        a.AcceptorCity,
				CountryCode: a.AcceptorCountry,
			},
		},
		Psp: entity.PSP{
			ID:     uuid.MustParse(a.PspID),
			Name:   a.PspName.StringVal,
			Prefix: a.PspPrefix.StringVal,
		},
		Exemption: entity.ExemptionType(a.Exemption.StringVal),
		ThreeDSecure: entity.ThreeDSecure{
			Version:                         a.ThreedsVersion.StringVal,
			AuthenticationVerificationValue: a.AccountholderAuthenticationValue.StringVal,
			DirectoryServerID:               a.DirectoryServerID.StringVal,
			EcommerceIndicator:              int(a.ThreedsOriginalEcommerceIndicator.Int64),
		},
		CardSchemeData: entity.CardSchemeData{
			Request: entity.CardSchemeRequest{
				ProcessingCode: entity.ProcessingCode{
					TransactionTypeCode: a.TransactionTypeCode.StringVal,
					FromAccountTypeCode: a.FromAccountTypeCode.StringVal,
					ToAccountTypeCode:   a.ToAccountTypeCode.StringVal,
				},
				POSEntryMode: entity.POSEntryMode{
					PanEntryMode: pos.PanEntryFromCode(a.PanEntryMode),
					PinEntryMode: pos.PinEntryFromCode(a.PinEntryMode),
				},
				RetrievalReferenceNumber: a.RetrievalReferenceNumber.StringVal,
			},
			Response: entity.CardSchemeResponse{
				ResponseCode: entity.ResponseCode{
					Value:       a.ResponseCode.StringVal,
					Description: entity.ResponseCodeFromString(a.ResponseCode.StringVal).Description,
				},
				AuthorizationIDResponse: a.AuthorizationIDResponse.StringVal,
			},
		},
		CitMitIndicator: entity.CitMitIndicator{
			InitiatedBy: entity.MapInitiatedByFromStr(a.TransactionInitiatedBy.StringVal),
			SubCategory: entity.MapSubCategoryFromStr(a.TransactionSubCategory.StringVal),
		},
	}
}

func mapVisaAuthorizationEntity(v VisaAuthorizationRecord) entity.Authorization {
	// TODO requires a field in entity.Authorization for retrievel_refence_number and created_at
	auth := mapAuthorizationEntity(v.AuthorizationRecord)
	auth.VisaSchemeData = entity.VisaSchemeData{
		Request: entity.VisaSchemeRequest{
			PosConditionCode: v.PointOfServiceConditionCode.StringVal,
			AdditionalPOSInformation: entity.AdditionalPOSInformation{
				TerminalType:                               v.AdditionalPosInfoTerminalType.StringVal,
				TerminalEntryCapability:                    v.AdditionalPosInfoTerminalEntryCapability.StringVal,
				TypeOrLevelIndicator:                       v.AdditionalPosInfoTypeOrLevelIndicator.StringVal,
				ChipConditionCode:                          v.ChipConditionCode.String(),
				SpecialConditionIndicator:                  v.SpecialConditionIndicator.String(),
				ChipTransactionIndicator:                   v.ChipTransactionIndicator.String(),
				ChipCardAuthenticationReliabilityIndicator: v.ChipCardAuthenticationReliabilityIndicator.String(),
				CardholderIDMethodIndicator:                v.CardholderIDMethodIndicator.String(),
				AdditionalAuthorizationIndicators:          v.AdditionalAuthorizationIndicators.String(),
			},
			PrivateUseFields: entity.PrivateUseFields{},
		},
		Response: entity.VisaSchemeResponse{
			TransactionId: int(v.TransactionIdentifier.Int64),
		},
	}
	return auth
}

func (ar AuthorizationRepository) GetAuthorizationWithSchemeData(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error) {
	stmt := spanner.Statement{
		SQL: `
			SELECT
				a.authorization_id, a.status, a.masked_pan,
				a.pan_token_id, a.card_scheme, a.amount,
				a.currency, a.localdatetime, a.source,
				a.is_initial_recurring, a.threeds_version,
				a.customer_reference, a.initial_trace_id, a.response_code,
				a.system_trace_audit_number, a.threeds_directory_server_transaction_id, a.psp_id,
				a.transmitted_at, a.card_issuer_id,
				a.card_issuer_name, a.card_issuer_countrycode,
				m.card_product_id, m.card_program_id,
				a.cardholder_transaction_type_code, a.cardholder_from_account_type_code,
				a.cardholder_to_account_type_code, a.card_acceptor_name, a.card_acceptor_city,
				a.card_acceptor_country, a.card_acceptor_id, a.card_acceptor_postal_code,
				a.card_acceptor_category_code, a.exemption, a.threeds_original_ecommerce_indicator,
				a.transaction_initiated_by, a.transaction_subcategory,
				a.point_of_service_pan_entry_mode, a.point_of_service_pin_entry_mode,
				p.name, p.prefix,
				m.network_reporting_date, m.financial_network_code, m.banknet_reference_number, m.pos_pin_capture_code,
				m.pin_service_code, m.security_protocol, m.cardholder_authentication, m.ucaf_collection_indicator,
				m.authorization_type, a.authorization_id_response,
				m.terminal_attendance, m.terminal_location, m.card_holder_presence,
				m.card_presence, m.card_capture_capabilities, m.transaction_status,
				m.transaction_security, m.card_holder_activated_terminal_level, m.card_data_terminal_input_capability_indicator,
				m.authorization_life_cycle, m.country_code, m.postal_code,
				m.reason_ucaf_downgrade, a.accountholder_authentication_value,
				v.point_of_service_condition_code, v.additional_pos_info_terminal_entry_capability,
				v.additional_pos_info_type_or_level_indicator, 
				a.retrieval_reference_number, v.additional_pos_info_terminal_type, v.transaction_identifier,
				v.chip_condition_code, v.special_condition_indicator, v.chip_transaction_indicator,
				v.chip_card_authentication_reliability_indicator, v.cardholder_id_method_indicator,
				v.additional_authorization_indicators
			FROM authorizations AS a
				 LEFT JOIN mastercard_authorizations AS m ON a.authorization_id = m.authorization_id
			     LEFT JOIN visa_authorizations AS v ON a.authorization_id = v.authorization_id 
				 JOIN psp AS p ON a.psp_id = p.psp_id
			WHERE a.authorization_id = @authorization_id
			  AND a.psp_id = @psp_id
		`,
		Params: map[string]interface{}{
			"authorization_id": authorizationID.String(),
			"psp_id":           pspID.String(),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, ar.readTimeout)
	defer cancel()

	iter := ar.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return entity.Authorization{}, entity.ErrRecordNotFound
		}
		return entity.Authorization{}, err
	}
	var cardScheme string
	if err = row.ColumnByName("card_scheme", &cardScheme); err != nil {
		return entity.Authorization{}, err
	}
	switch cardScheme {
	case "mastercard":
		var m MastercardAuthorizationRecord
		if err = row.ToStructLenient(&m); err != nil {
			return entity.Authorization{}, err
		}
		return mapMastercardAuthorizationEntity(m), nil
	case "visa":
		var v VisaAuthorizationRecord
		if err = row.ToStructLenient(&v); err != nil {
			return entity.Authorization{}, err
		}
		return mapVisaAuthorizationEntity(v), nil
	default:
		var a AuthorizationRecord
		if err = row.ToStructLenient(&a); err != nil {
			return entity.Authorization{}, err
		}
		return mapAuthorizationEntity(a), nil
	}
}

func (ar AuthorizationRepository) CreateVisaAuthorization(ctx context.Context, a entity.Authorization) error {
	_, err := ar.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		statement := spanner.Statement{
			SQL: `
				INSERT INTO visa_authorizations (
					authorization_id,
				    point_of_service_condition_code,
				    additional_pos_info_terminal_type, additional_pos_info_terminal_entry_capability,
				    additional_pos_info_type_or_level_indicator, transaction_identifier, created_at,
				    chip_condition_code, special_condition_indicator, chip_transaction_indicator,
					chip_card_authentication_reliability_indicator, cardholder_id_method_indicator,
					additional_authorization_indicators
				) VALUES (
					@authorization_id,
				    @point_of_service_condition_code,
				    @additional_pos_info_terminal_type, @additional_pos_info_terminal_entry_capability,
				    @additional_pos_info_type_or_level_indicator, @transaction_identifier, @created_at,
				    @chip_condition_code, @special_condition_indicator, @chip_transaction_indicator,
					@chip_card_authentication_reliability_indicator, @cardholder_id_method_indicator,
					@additional_authorization_indicators
				)`,
			Params: mapVisaAuthorizationParams(a),
		}

		ctx, cancel := context.WithTimeout(ctx, ar.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, statement)

		return err
	})

	return err
}

// AuthorizationRecord is exported to allow spanner.ToStructLenient() to fill it - do not use outside this file.
type AuthorizationRecord struct {
	ID                                string             `spanner:"authorization_id"`
	Status                            string             `spanner:"status"`
	MaskedPan                         string             `spanner:"masked_pan"`
	PanTokenID                        string             `spanner:"pan_token_id"`
	CardScheme                        string             `spanner:"card_scheme"`
	Amount                            int64              `spanner:"amount"`
	Currency                          string             `spanner:"currency"`
	LocalDateTime                     time.Time          `spanner:"localdatetime"`
	Source                            string             `spanner:"source"`
	IsInitialRecurring                spanner.NullBool   `spanner:"is_initial_recurring"`
	ThreedsVersion                    spanner.NullString `spanner:"threeds_version"`
	CreatedAt                         time.Time          `spanner:"created_at"`
	CustomerReference                 spanner.NullString `spanner:"customer_reference"`
	UpdatedAt                         spanner.NullTime   `spanner:"updated_at"`
	InitialTraceID                    spanner.NullString `spanner:"initial_trace_id"`
	ResponseCode                      spanner.NullString `spanner:"response_code"`
	Stan                              spanner.NullInt64  `spanner:"system_trace_audit_number"`
	DirectoryServerID                 spanner.NullString `spanner:"threeds_directory_server_transaction_id"`
	PspID                             string             `spanner:"psp_id"`
	TransmissionDate                  spanner.NullTime   `spanner:"transmitted_at"`
	CardIssuerID                      spanner.NullString `spanner:"card_issuer_id"`
	CardIssuerName                    spanner.NullString `spanner:"card_issuer_name"`
	CardIssuerCountryCode             spanner.NullString `spanner:"card_issuer_countrycode"`
	TransactionTypeCode               spanner.NullString `spanner:"cardholder_transaction_type_code"`
	FromAccountTypeCode               spanner.NullString `spanner:"cardholder_from_account_type_code"`
	ToAccountTypeCode                 spanner.NullString `spanner:"cardholder_to_account_type_code"`
	AcceptorName                      string             `spanner:"card_acceptor_name"`
	AcceptorCity                      string             `spanner:"card_acceptor_city"`
	AcceptorCountry                   string             `spanner:"card_acceptor_country"`
	AcceptorID                        string             `spanner:"card_acceptor_id"`
	AcceptorPostalCode                spanner.NullString `spanner:"card_acceptor_postal_code"`
	CategoryCode                      string             `spanner:"card_acceptor_category_code"`
	ThreedsOriginalEcommerceIndicator spanner.NullInt64  `spanner:"threeds_original_ecommerce_indicator"`
	PspName                           spanner.NullString `spanner:"name"`
	PspPrefix                         spanner.NullString `spanner:"prefix"`
	Exemption                         spanner.NullString `spanner:"exemption"`
	TransactionInitiatedBy            spanner.NullString `spanner:"transaction_initiated_by"`
	TransactionSubCategory            spanner.NullString `spanner:"transaction_subcategory"`
	AccountholderAuthenticationValue  spanner.NullString `spanner:"accountholder_authentication_value"`
	PanEntryMode                      string             `spanner:"point_of_service_pan_entry_mode"`
	PinEntryMode                      string             `spanner:"point_of_service_pin_entry_mode"`
	AuthorizationIDResponse           spanner.NullString `spanner:"authorization_id_response"`
	RetrievalReferenceNumber          spanner.NullString `spanner:"retrieval_reference_number"`
}

// MastercardAuthorizationRecord is exported to allow spanner.ToStructLenient() to fill it - do not use outside this file.
type MastercardAuthorizationRecord struct {
	AuthorizationRecord
	PosPinCaptureCode                        spanner.NullString `spanner:"pos_pin_capture_code"`
	PinServiceCode                           spanner.NullString `spanner:"pin_service_code"`
	SecurityProtocol                         int64              `spanner:"security_protocol"`
	CardholderAuthentication                 int64              `spanner:"cardholder_authentication"`
	UCAFCollectionIndicator                  int64              `spanner:"ucaf_collection_indicator"`
	NetworkReportingDate                     spanner.NullString `spanner:"network_reporting_date"`
	FinancialNetworkCode                     spanner.NullString `spanner:"financial_network_code"`
	BanknetReferenceNumber                   spanner.NullString `spanner:"banknet_reference_number"`
	Reference                                spanner.NullString `spanner:"reference"`
	AuthorizationType                        string             `spanner:"authorization_type"`
	AuthorizationIDResponse                  spanner.NullString `spanner:"authorization_id_response"`
	TerminalAttendance                       spanner.NullInt64  `spanner:"terminal_attendance"`
	TerminalLocation                         spanner.NullInt64  `spanner:"terminal_location"`
	CardHolderPresence                       spanner.NullInt64  `spanner:"card_holder_presence"`
	CardPresence                             spanner.NullInt64  `spanner:"card_presence"`
	CardCaptureCapabilities                  spanner.NullInt64  `spanner:"card_capture_capabilities"`
	TransactionStatus                        spanner.NullInt64  `spanner:"transaction_status"`
	TransactionSecurity                      spanner.NullInt64  `spanner:"transaction_security"`
	CardHolderActivatedTerminalLevel         spanner.NullInt64  `spanner:"card_holder_activated_terminal_level"`
	CardDataTerminalInputCapabilityIndicator spanner.NullInt64  `spanner:"card_data_terminal_input_capability_indicator"`
	AuthorizationLifeCycle                   spanner.NullString `spanner:"authorization_life_cycle"`
	CountryCode                              spanner.NullString `spanner:"country_code"`
	PostalCode                               spanner.NullString `spanner:"postal_code"`
	ReasonForUCAFDowngrade                   spanner.NullInt64  `spanner:"reason_ucaf_downgrade"`
	CardProductId                            spanner.NullString `spanner:"card_product_id"`
	CardProgramId                            spanner.NullString `spanner:"card_program_id"`
}

// VisaAuthorizationRecord is exported to allow spanner.ToStructLenient() to fill it - do not use outside this file.
type VisaAuthorizationRecord struct {
	AuthorizationRecord
	PointOfServiceConditionCode                spanner.NullString `spanner:"point_of_service_condition_code"`
	AdditionalPosInfoTerminalType              spanner.NullString `spanner:"additional_pos_info_terminal_type"`
	AdditionalPosInfoTerminalEntryCapability   spanner.NullString `spanner:"additional_pos_info_terminal_entry_capability"`
	AdditionalPosInfoTypeOrLevelIndicator      spanner.NullString `spanner:"additional_pos_info_type_or_level_indicator"`
	TransactionIdentifier                      spanner.NullInt64  `spanner:"transaction_identifier"`
	CreatedAt                                  spanner.NullTime   `spanner:"created_at"`
	ChipConditionCode                          spanner.NullString `spanner:"chip_condition_code"`
	SpecialConditionIndicator                  spanner.NullString `spanner:"special_condition_indicator"`
	ChipTransactionIndicator                   spanner.NullString `spanner:"chip_transaction_indicator"`
	ChipCardAuthenticationReliabilityIndicator spanner.NullString `spanner:"chip_card_authentication_reliability_indicator"`
	CardholderIDMethodIndicator                spanner.NullString `spanner:"cardholder_id_method_indicator"`
	AdditionalAuthorizationIndicators          spanner.NullString `spanner:"additional_authorization_indicators"`
}

func mapRowToAuthorizationEntity(mar MastercardAuthorizationRecord) entity.Authorization {
	_, schemeResponseMessage := mastercard.MapResponseCode(mar.ResponseCode.StringVal)
	// TODO check if this needs to be added?
	//if r.Type.Valid {
	//	authorizationType = processing.AuthorizationType(r.Type.StringVal)
	//}

	return entity.Authorization{
		ID:                       uuid.MustParse(mar.ID),
		Amount:                   int(mar.Amount),
		Currency:                 currencycode.Must(mar.Currency),
		CustomerReference:        mar.CustomerReference.StringVal,
		Source:                   entity.Source(mar.Source),
		LocalTransactionDateTime: data.LocalTransactionDateTime(mar.LocalDateTime),
		Status:                   entity.Status(mar.Status),
		Stan:                     int(mar.Stan.Int64),
		ProcessingDate:           mar.TransmissionDate.Time,
		CreatedAt:                mar.CreatedAt,
		Recurring: entity.Recurring{
			Initial:    mar.IsInitialRecurring.Bool,
			Subsequent: mar.InitialTraceID.Valid,
			TraceID:    mar.InitialTraceID.StringVal,
		},
		Card: entity.Card{
			MaskedPan:  mar.MaskedPan,
			PanTokenID: mar.PanTokenID,
			Info: cardinfo.Range{
				Scheme:            mar.CardScheme,
				ProductID:         mar.CardProductId.StringVal,
				ProgramID:         mar.CardProgramId.StringVal,
				IssuerID:          mar.CardIssuerID.StringVal,
				IssuerName:        mar.CardIssuerName.StringVal,
				IssuerCountryCode: mar.CardIssuerCountryCode.StringVal,
			},
		},
		CardAcceptor: entity.CardAcceptor{
			CategoryCode: mar.CategoryCode,
			ID:           mar.AcceptorID,
			Name:         mar.AcceptorName,
			Address: entity.CardAcceptorAddress{
				PostalCode:  mar.AcceptorPostalCode.StringVal,
				City:        mar.AcceptorCity,
				CountryCode: mar.AcceptorCountry,
			},
		},
		Psp: entity.PSP{
			ID: uuid.MustParse(mar.PspID),
		},
		Exemption: entity.ExemptionType(mar.Exemption.StringVal),
		ThreeDSecure: entity.ThreeDSecure{
			Version:            mar.ThreedsVersion.StringVal,
			DirectoryServerID:  mar.DirectoryServerID.StringVal,
			EcommerceIndicator: int(mar.ThreedsOriginalEcommerceIndicator.Int64),
		},
		CardSchemeData: entity.CardSchemeData{Response: entity.CardSchemeResponse{
			Status: entity.AuthorizationStatusFromString(mar.Status),
			ResponseCode: entity.ResponseCode{
				Value:       mar.ResponseCode.StringVal,
				Description: schemeResponseMessage,
			},
			AuthorizationIDResponse: mar.Reference.StringVal,
		}},
		CitMitIndicator: entity.CitMitIndicator{
			InitiatedBy: entity.MapInitiatedByFromStr(mar.TransactionInitiatedBy.StringVal),
			SubCategory: entity.MapSubCategoryFromStr(mar.TransactionSubCategory.StringVal),
		},
		MastercardSchemeData: entity.MastercardSchemeData{
			Request: entity.MastercardSchemeRequest{
				AuthorizationType: entity.AuthorizationType(mar.AuthorizationType),
			},
			Response: entity.MastercardSchemeResponse{
				AdditionalData: entity.AdditionalResponseData{
					AppliedEcommerceIndicator: &entity.SLI{
						SecurityProtocol:         int(mar.SecurityProtocol),
						CardholderAuthentication: int(mar.CardholderAuthentication),
						UCAFCollectionIndicator:  int(mar.UCAFCollectionIndicator),
					},
					ReasonForUCAFDowngrade: mapping.Ptr(int(mar.ReasonForUCAFDowngrade.Int64)),
				},
				TraceID: entity.MTraceID{
					NetworkReportingDate:   mar.NetworkReportingDate.StringVal,
					BanknetReferenceNumber: mar.BanknetReferenceNumber.StringVal,
					FinancialNetworkCode:   mar.FinancialNetworkCode.StringVal,
				},
			},
		},
	}
}

func Ptr[T any](v T) *T {
	return &v
}

func mapMastercardAuthorizationEntity(m MastercardAuthorizationRecord) entity.Authorization {
	return entity.Authorization{
		ID:                       uuid.MustParse(m.ID),
		Amount:                   int(m.Amount),
		Currency:                 currencycode.Must(m.Currency),
		CustomerReference:        m.CustomerReference.StringVal,
		Source:                   entity.Source(m.Source),
		LocalTransactionDateTime: data.LocalTransactionDateTime(m.LocalDateTime),
		Status:                   entity.Status(m.Status),
		Stan:                     int(m.Stan.Int64),
		ProcessingDate:           m.TransmissionDate.Time,
		CreatedAt:                m.CreatedAt,
		Recurring: entity.Recurring{
			Initial:    m.IsInitialRecurring.Bool,
			Subsequent: m.InitialTraceID.Valid,
			TraceID:    m.InitialTraceID.StringVal,
		},
		Card: entity.Card{
			MaskedPan:  m.MaskedPan,
			PanTokenID: m.PanTokenID,
			Info: cardinfo.Range{
				IssuerID:          m.CardIssuerID.String(),
				IssuerName:        m.CardIssuerName.String(),
				IssuerCountryCode: m.CardIssuerCountryCode.String(),
				ProductID:         m.CardProductId.String(),
				ProgramID:         m.CardProgramId.String(),
				Scheme:            m.CardScheme,
			},
		},
		CardAcceptor: entity.CardAcceptor{
			CategoryCode: m.CategoryCode,
			ID:           m.AcceptorID,
			Name:         m.AcceptorName,
			Address: entity.CardAcceptorAddress{
				PostalCode:  m.AcceptorPostalCode.StringVal,
				City:        m.AcceptorCity,
				CountryCode: m.AcceptorCountry,
			},
		},
		// TODO: check if this needs to be added.
		//AuthorizationType:        authorizationType,
		Psp: entity.PSP{
			ID:     uuid.MustParse(m.PspID),
			Name:   m.PspName.StringVal,
			Prefix: m.PspPrefix.StringVal,
		},
		Exemption: entity.ExemptionType(m.Exemption.String()),
		ThreeDSecure: entity.ThreeDSecure{
			Version:                         m.ThreedsVersion.StringVal,
			AuthenticationVerificationValue: m.AccountholderAuthenticationValue.StringVal,
			EcommerceIndicator:              int(m.ThreedsOriginalEcommerceIndicator.Int64),
			DirectoryServerID:               m.DirectoryServerID.StringVal,
		},
		CardSchemeData: entity.CardSchemeData{
			Request: entity.CardSchemeRequest{
				ProcessingCode: entity.ProcessingCode{
					TransactionTypeCode: m.TransactionTypeCode.StringVal,
					FromAccountTypeCode: m.FromAccountTypeCode.StringVal,
					ToAccountTypeCode:   m.ToAccountTypeCode.StringVal,
				},
				POSEntryMode: entity.POSEntryMode{
					PanEntryMode: pos.PanEntryFromCode(m.PanEntryMode),
					PinEntryMode: pos.PinEntryFromCode(m.PinEntryMode),
				},
			},
			Response: entity.CardSchemeResponse{
				ResponseCode: entity.ResponseCode{
					Value:       m.ResponseCode.StringVal,
					Description: entity.ResponseCodeFromString(m.ResponseCode.StringVal).Description,
				},
				AuthorizationIDResponse: m.AuthorizationIDResponse.StringVal,
			}},
		CitMitIndicator: entity.CitMitIndicator{
			InitiatedBy: entity.MapInitiatedByFromStr(m.TransactionInitiatedBy.StringVal),
			SubCategory: entity.MapSubCategoryFromStr(m.TransactionSubCategory.StringVal),
		},
		MastercardSchemeData: entity.MastercardSchemeData{
			Request: entity.MastercardSchemeRequest{
				AuthorizationType: entity.AuthorizationType(m.AuthorizationType),
				PosPinCaptureCode: m.PosPinCaptureCode.StringVal,
				AdditionalData: entity.AdditionalRequestData{
					PinServiceCode: m.PinServiceCode.StringVal,
				},
				PointOfServiceData: entity.PointOfServiceData{
					TerminalAttendance:                       int(m.TerminalAttendance.Int64),
					TerminalLocation:                         int(m.TerminalLocation.Int64),
					CardHolderPresence:                       int(m.CardHolderPresence.Int64),
					CardPresence:                             int(m.CardPresence.Int64),
					CardCaptureCapabilities:                  int(m.CardCaptureCapabilities.Int64),
					TransactionStatus:                        int(m.TransactionStatus.Int64),
					TransactionSecurity:                      int(m.TransactionSecurity.Int64),
					CardHolderActivatedTerminalLevel:         int(m.CardHolderActivatedTerminalLevel.Int64),
					CardDataTerminalInputCapabilityIndicator: int(m.CardDataTerminalInputCapabilityIndicator.Int64),
					AuthorizationLifeCycle:                   m.AuthorizationLifeCycle.StringVal,
					CountryCode:                              m.CountryCode.StringVal,
					PostalCode:                               m.PostalCode.StringVal,
				},
			},
			Response: entity.MastercardSchemeResponse{
				AdditionalData: entity.AdditionalResponseData{
					AppliedEcommerceIndicator: &entity.SLI{
						SecurityProtocol:         int(m.SecurityProtocol),
						CardholderAuthentication: int(m.CardholderAuthentication),
						UCAFCollectionIndicator:  int(m.UCAFCollectionIndicator),
					},
					ReasonForUCAFDowngrade: mapping.Ptr(int(m.ReasonForUCAFDowngrade.Int64)),
				},
				TraceID: entity.MTraceID{
					FinancialNetworkCode:   m.FinancialNetworkCode.StringVal,
					BanknetReferenceNumber: m.BanknetReferenceNumber.StringVal,
					NetworkReportingDate:   m.NetworkReportingDate.StringVal,
				},
			},
		},
	}
}

func mapMastercardAuthorizationParams(a entity.Authorization) map[string]interface{} {
	return map[string]interface{}{
		"authorization_id":                              a.ID.String(),
		"pos_pin_capture_code":                          sql.NewNullString(a.MastercardSchemeData.Request.PosPinCaptureCode),
		"pin_service_code":                              sql.NewNullString(a.MastercardSchemeData.Request.AdditionalData.PinServiceCode),
		"security_protocol":                             a.MastercardSchemeData.Response.AdditionalData.AppliedEcommerceIndicator.SecurityProtocol,
		"cardholder_authentication":                     a.MastercardSchemeData.Response.AdditionalData.AppliedEcommerceIndicator.CardholderAuthentication,
		"ucaf_collection_indicator":                     a.MastercardSchemeData.Response.AdditionalData.AppliedEcommerceIndicator.UCAFCollectionIndicator,
		"banknet_reference_number":                      sql.NewNullString(a.MastercardSchemeData.Response.TraceID.BanknetReferenceNumber),
		"network_reporting_date":                        sql.NewNullString(a.MastercardSchemeData.Response.TraceID.NetworkReportingDate),
		"financial_network_code":                        sql.NewNullString(a.MastercardSchemeData.Response.TraceID.FinancialNetworkCode),
		"reference":                                     sql.NewNullString(a.CardSchemeData.Response.AuthorizationIDResponse),
		"created_at":                                    time.Now(),
		"authorization_type":                            a.MastercardSchemeData.Request.AuthorizationType,
		"terminal_attendance":                           a.MastercardSchemeData.Request.PointOfServiceData.TerminalAttendance,
		"terminal_location":                             a.MastercardSchemeData.Request.PointOfServiceData.TerminalLocation,
		"card_holder_presence":                          a.MastercardSchemeData.Request.PointOfServiceData.CardHolderPresence,
		"card_presence":                                 a.MastercardSchemeData.Request.PointOfServiceData.CardPresence,
		"card_capture_capabilities":                     a.MastercardSchemeData.Request.PointOfServiceData.CardCaptureCapabilities,
		"transaction_status":                            a.MastercardSchemeData.Request.PointOfServiceData.TransactionStatus,
		"transaction_security":                          a.MastercardSchemeData.Request.PointOfServiceData.TransactionSecurity,
		"card_holder_activated_terminal_level":          a.MastercardSchemeData.Request.PointOfServiceData.CardHolderActivatedTerminalLevel,
		"card_data_terminal_input_capability_indicator": a.MastercardSchemeData.Request.PointOfServiceData.CardDataTerminalInputCapabilityIndicator,
		"authorization_life_cycle":                      a.MastercardSchemeData.Request.PointOfServiceData.AuthorizationLifeCycle,
		"country_code":                                  a.MastercardSchemeData.Request.PointOfServiceData.CountryCode,
		"postal_code":                                   a.MastercardSchemeData.Request.PointOfServiceData.PostalCode,
		"reason_ucaf_downgrade":                         sql.NewNullInt64Pointer(a.MastercardSchemeData.Response.AdditionalData.ReasonForUCAFDowngrade),
		"exemption":                                     sql.NewNullString(string(a.Exemption)),
		"card_product_id":                               sql.NewNullString(a.Card.Info.ProductID), // TODO: should we move this field to authorizations table?
		"card_program_id":                               sql.NewNullString(a.Card.Info.ProgramID), // TODO: should we move this field to authorizations table?
	}
}

func mapVisaAuthorizationParams(a entity.Authorization) map[string]interface{} {
	return map[string]interface{}{
		"authorization_id":                               a.ID.String(),
		"point_of_service_condition_code":                a.VisaSchemeData.Request.PosConditionCode,
		"additional_pos_info_terminal_type":              a.VisaSchemeData.Request.AdditionalPOSInformation.TerminalType,
		"additional_pos_info_terminal_entry_capability":  a.VisaSchemeData.Request.AdditionalPOSInformation.TerminalEntryCapability,
		"additional_pos_info_type_or_level_indicator":    a.VisaSchemeData.Request.AdditionalPOSInformation.TypeOrLevelIndicator,
		"transaction_identifier":                         sql.NewNullInt64(a.VisaSchemeData.Response.TransactionId),
		"created_at":                                     time.Now().Format(time.RFC3339),
		"chip_condition_code":                            a.VisaSchemeData.Request.AdditionalPOSInformation.ChipConditionCode,
		"special_condition_indicator":                    a.VisaSchemeData.Request.AdditionalPOSInformation.SpecialConditionIndicator,
		"chip_transaction_indicator":                     a.VisaSchemeData.Request.AdditionalPOSInformation.ChipTransactionIndicator,
		"chip_card_authentication_reliability_indicator": a.VisaSchemeData.Request.AdditionalPOSInformation.ChipCardAuthenticationReliabilityIndicator,
		"cardholder_id_method_indicator":                 a.VisaSchemeData.Request.AdditionalPOSInformation.CardholderIDMethodIndicator,
		"additional_authorization_indicators":            a.VisaSchemeData.Request.AdditionalPOSInformation.AdditionalAuthorizationIndicators,
	}
}

func determineDateRange(startDate time.Time, endDate time.Time) (time.Time, time.Time) {
	now := time.Now()
	var dateRangeStart, dateRangeEnd time.Time
	// the default date range will be between one week ago and today
	// The dateRange is submitted as yyyy-mm-dd and will have 00:00:00 as the timestamp
	// DateRangeStart needs to start at 00:00:00
	// DateRangeEnd needs to end at 23:59:59

	// Milliseconds for a week so we can count back if we have no dates
	durationWeek := time.Hour * 24 * 7

	//I think this isn't working when either start or enddate is empty and I dislike the fact that we have a dependency on the FO for this
	if startDate.IsZero() && endDate.IsZero() {
		dateRangeStart = now.Add(-durationWeek)
		dateRangeEnd = now
	} else {
		dateRangeStart = startDate
		dateRangeEnd = endDate
	}
	dateRangeStart = time.Date(dateRangeStart.Year(), dateRangeStart.Month(), dateRangeStart.Day(), 0, 0, 0, 0, dateRangeStart.Location())
	dateRangeEnd = time.Date(dateRangeEnd.Year(), dateRangeEnd.Month(), dateRangeEnd.Day(), 23, 59, 59, 0, dateRangeEnd.Location())
	return dateRangeStart, dateRangeEnd
}

func (ar AuthorizationRepository) AuthorizationAlreadyReversed(ctx context.Context, id uuid.UUID) (bool, error) {
	stmt := spanner.Statement{
		SQL: `
			    SELECT 1
				FROM authorization_reversals
				WHERE authorization_id = @authorization_id
				AND status = 'succeeded'
		`,
		Params: map[string]interface{}{"authorization_id": id.String()},
	}

	ctx, cancel := context.WithTimeout(ctx, ar.readTimeout)
	defer cancel()

	iter := ar.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	_, err := iter.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
