package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/pos"
	mapping "gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/mastercard"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
	"google.golang.org/api/iterator"
)

type RefundRepository struct {
	client       *spanner.Client
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewRefundRepository(
	client *spanner.Client,
	readTimeout time.Duration,
	writeTimeout time.Duration) *RefundRepository {
	return &RefundRepository{
		client:       client,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (rr RefundRepository) CreateRefund(ctx context.Context, r entity.Refund) error {
	_, err := rr.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		statement := spanner.Statement{
			SQL: `
				INSERT INTO refunds (
					refund_id, masked_pan, pan_token_id,
                    amount, currency, localdatetime,
                    source, customer_reference, psp_id,
                    card_acceptor_id, card_acceptor_name, card_acceptor_postal_code,
                    card_acceptor_city, card_acceptor_country, card_acceptor_category_code,
                    created_at, status, card_scheme,
                    card_issuer_id, card_issuer_name, card_issuer_countrycode
				) VALUES (
					@refund_id, @masked_pan, @pan_token_id,
                    @amount, @currency, @localdatetime,
                    @source, @customer_reference, @psp_id,
                    @card_acceptor_id, @card_acceptor_name, @card_acceptor_postal_code,
                    @card_acceptor_city, @card_acceptor_country, @card_acceptor_category_code,
                    @created_at, 'new', @card_scheme,
                    @card_issuer_id, @card_issuer_name, @card_issuer_countrycode
				)`,
			Params: mapCreateRefundParams(r),
		}

		ctx, cancel := context.WithTimeout(ctx, rr.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, statement)

		return err
	})

	return err
}

func mapCreateRefundParams(r entity.Refund) map[string]interface{} {
	return map[string]interface{}{
		"refund_id":                   r.ID.String(),
		"masked_pan":                  r.Card.MaskedPan,
		"pan_token_id":                r.Card.PanTokenID,
		"card_scheme":                 r.Card.Info.Scheme,
		"amount":                      r.Amount,
		"currency":                    r.Currency.Alpha3(),
		"localdatetime":               r.LocalTransactionDateTime,
		"source":                      string(r.Source),
		"created_at":                  time.Now(),
		"customer_reference":          mapping.NewNullString(r.CustomerReference),
		"psp_id":                      r.Psp.ID.String(),
		"card_issuer_id":              mapping.NewNullString(r.Card.Info.IssuerID),
		"card_issuer_name":            mapping.NewNullString(r.Card.Info.IssuerName),
		"card_issuer_countrycode":     mapping.NewNullString(r.Card.Info.IssuerCountryCode),
		"card_product_id":             mapping.NewNullString(r.Card.Info.ProductID),
		"card_program_id":             mapping.NewNullString(r.Card.Info.ProgramID),
		"card_acceptor_name":          r.CardAcceptor.Name,
		"card_acceptor_city":          r.CardAcceptor.Address.City,
		"card_acceptor_country":       r.CardAcceptor.Address.CountryCode,
		"card_acceptor_id":            r.CardAcceptor.ID,
		"card_acceptor_postal_code":   r.CardAcceptor.Address.PostalCode,
		"card_acceptor_category_code": r.CardAcceptor.CategoryCode,
	}
}

func (rr RefundRepository) CreateMastercardRefund(ctx context.Context, r entity.Refund) error {
	_, err := rr.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		statement := spanner.Statement{
			SQL: `
				INSERT INTO mastercard_refunds (
					refund_id, network_reporting_date, financial_network_code,
					reference,
					banknet_reference_number, created_at, authorization_type, 
				    authorization_id_response, terminal_attendance, terminal_location,
				    card_holder_presence, card_presence, card_capture_capabilities, 
				    transaction_status, transaction_security, card_holder_activated_terminal_level,
				    card_data_terminal_input_capability_indicator, authorization_life_cycle, country_code,
				    postal_code, card_product_id, card_program_id
				) VALUES (
					@refund_id, @network_reporting_date, @financial_network_code,
					@reference,
					@banknet_reference_number, @created_at, @authorization_type,
				    @authorization_id_response, @terminal_attendance, @terminal_location, 
				    @card_holder_presence, @card_presence, @card_capture_capabilities,
				    @transaction_status, @transaction_security, @card_holder_activated_terminal_level,
				    @card_data_terminal_input_capability_indicator, @authorization_life_cycle, @country_code,
				    @postal_code, @card_product_id, @card_program_id
				)`,
			Params: mapMastercardRefundParams(r),
		}

		ctx, cancel := context.WithTimeout(ctx, rr.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, statement)

		return err
	})

	return err
}

func (rr RefundRepository) CreateVisaRefund(ctx context.Context, r entity.Refund) error {
	stmt := spanner.NewStatement(`
		INSERT INTO visa_refunds (refund_id, created_at, transaction_identifier)
		VALUES
		    (@refund_id, @created_at, @transaction_identifier)
	`)

	stmt.Params["refund_id"] = r.ID.String()
	stmt.Params["created_at"] = time.Now()
	stmt.Params["transaction_identifier"] = r.VisaSchemeData.Response.TransactionId

	_, err := rr.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, rr.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, stmt)

		return err
	})

	return err
}

func mapMastercardRefundParams(r entity.Refund) map[string]interface{} {
	return map[string]interface{}{
		"refund_id":                                     r.ID.String(),
		"banknet_reference_number":                      mapping.NewNullString(r.MastercardSchemeData.Response.TraceID.BanknetReferenceNumber),
		"network_reporting_date":                        mapping.NewNullString(r.MastercardSchemeData.Response.TraceID.NetworkReportingDate),
		"financial_network_code":                        mapping.NewNullString(r.MastercardSchemeData.Response.TraceID.FinancialNetworkCode),
		"reference":                                     mapping.NewNullString(r.CardSchemeData.Response.AuthorizationIDResponse),
		"created_at":                                    time.Now(),
		"authorization_type":                            r.MastercardSchemeData.Request.AuthorizationType,
		"authorization_id_response":                     mapping.NewNullString(r.CardSchemeData.Response.AuthorizationIDResponse),
		"terminal_attendance":                           r.MastercardSchemeData.Request.PointOfServiceData.TerminalAttendance,
		"terminal_location":                             r.MastercardSchemeData.Request.PointOfServiceData.TerminalLocation,
		"card_holder_presence":                          r.MastercardSchemeData.Request.PointOfServiceData.CardHolderPresence,
		"card_presence":                                 r.MastercardSchemeData.Request.PointOfServiceData.CardPresence,
		"card_capture_capabilities":                     r.MastercardSchemeData.Request.PointOfServiceData.CardCaptureCapabilities,
		"transaction_status":                            r.MastercardSchemeData.Request.PointOfServiceData.TransactionStatus,
		"transaction_security":                          r.MastercardSchemeData.Request.PointOfServiceData.TransactionSecurity,
		"card_holder_activated_terminal_level":          r.MastercardSchemeData.Request.PointOfServiceData.CardHolderActivatedTerminalLevel,
		"card_data_terminal_input_capability_indicator": r.MastercardSchemeData.Request.PointOfServiceData.CardDataTerminalInputCapabilityIndicator,
		"authorization_life_cycle":                      r.MastercardSchemeData.Request.PointOfServiceData.AuthorizationLifeCycle,
		"country_code":                                  r.MastercardSchemeData.Request.PointOfServiceData.CountryCode,
		"postal_code":                                   r.MastercardSchemeData.Request.PointOfServiceData.PostalCode,
		"card_program_id":                               r.Card.Info.ProgramID,
		"card_product_id":                               r.Card.Info.ProductID,
	}
}

func (rr RefundRepository) GetRefund(ctx context.Context, pspID, refundID uuid.UUID) (entity.Refund, error) {
	stmt := spanner.Statement{
		SQL: `
			SELECT
				r.refund_id, r.status, r.masked_pan,
				r.pan_token_id, r.card_scheme, r.amount,
				r.currency, r.localdatetime, r.source,
				r.customer_reference, r.response_code, r.system_trace_audit_number,
				r.psp_id, r.transmitted_at, r.card_issuer_id,
				r.card_issuer_name, r.card_issuer_countrycode, r.cardholder_transaction_type_code,
				r.cardholder_from_account_type_code, r.cardholder_to_account_type_code, r.card_acceptor_name,
				r.card_acceptor_city, r.card_acceptor_country, r.card_acceptor_id,
				r.card_acceptor_postal_code, r.card_acceptor_category_code, p.name,
				r.point_of_service_pan_entry_mode, r.point_of_service_pin_entry_mode,
				p.prefix, m.network_reporting_date, m.financial_network_code,
				m.banknet_reference_number,
				m.authorization_type, m.authorization_id_response, m.terminal_attendance,
				m.terminal_location, m.card_holder_presence, m.card_presence,
				m.card_capture_capabilities, m.transaction_status, m.transaction_security,
				m.card_holder_activated_terminal_level, m.card_data_terminal_input_capability_indicator, m.authorization_life_cycle,
				m.country_code, m.postal_code, m.card_program_id, m.card_product_id
			FROM refunds AS r
					 LEFT OUTER JOIN mastercard_refunds AS m ON r.refund_id = m.refund_id
					 JOIN psp AS p ON r.psp_id = p.psp_id
			WHERE r.refund_id = @refund_id
			  AND r.psp_id = @psp_id
		`,
		Params: map[string]interface{}{
			"refund_id": refundID.String(),
			"psp_id":    pspID.String(),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, rr.readTimeout)
	defer cancel()

	iter := rr.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return entity.Refund{}, entity.ErrRecordNotFound
		}
		return entity.Refund{}, err
	}

	var cardScheme string
	if err = row.ColumnByName("card_scheme", &cardScheme); err != nil {
		return entity.Refund{}, err
	}

	switch cardScheme {
	case "mastercard":
		var m mastercardRefundRecord
		if err = row.ToStruct(&m); err != nil {
			return entity.Refund{}, err
		}
		return mapMastercardRowToRefundEntity(m), nil
	default:
		var r refundRecord
		if err = row.ToStructLenient(&r); err != nil {
			return entity.Refund{}, err
		}
		return mapRowToRefundEntity(r), nil
	}
}

func (rr RefundRepository) UpdateRefundResponse(ctx context.Context, r entity.Refund) error {
	stmt := spanner.Statement{
		SQL: `UPDATE refunds
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
				WHERE refund_id = @refund_id`,
		Params: mapUpdateRefundResponseParams(r),
	}

	_, err := rr.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, rr.writeTimeout)
		defer cancel()

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to update refund: %w", err)
		}
		if rowCount != 1 {
			return fmt.Errorf("no record found with ID: %s", r.ID.String())
		}

		return err
	})

	return err
}

func mapUpdateRefundResponseParams(r entity.Refund) map[string]interface{} {
	return map[string]interface{}{
		"refund_id":                         r.ID.String(),
		"status":                            r.CardSchemeData.Response.Status.String(),
		"system_trace_audit_number":         r.Stan,
		"updated_at":                        time.Now(),
		"response_code":                     mapping.NewNullString(r.CardSchemeData.Response.ResponseCode.Value),
		"authorization_id_response":         mapping.NewNullString(r.CardSchemeData.Response.AuthorizationIDResponse),
		"retrieval_reference_number":        mapping.NewNullString(r.CardSchemeData.Request.RetrievalReferenceNumber),
		"transmitted_at":                    r.ProcessingDate,
		"card_holder_transaction_type_code": r.CardSchemeData.Request.ProcessingCode.TransactionTypeCode,
		"cardholder_from_account_type_code": r.CardSchemeData.Request.ProcessingCode.FromAccountTypeCode,
		"cardholder_to_account_type_code":   r.CardSchemeData.Request.ProcessingCode.ToAccountTypeCode,
		"point_of_service_pan_entry_mode":   pos.PanEntryCode(r.CardSchemeData.Request.POSEntryMode.PanEntryMode),
		"point_of_service_pin_entry_mode":   pos.PinEntryCode(r.CardSchemeData.Request.POSEntryMode.PinEntryMode),
	}
}

func (rr RefundRepository) GetAllRefunds(ctx context.Context, pspID uuid.UUID, filters entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Refund, error) {
	stmt := spanner.Statement{
		SQL: fmt.Sprintf(`
			SELECT 
			    r.refund_id, r.psp_id, r.card_scheme,
				r.masked_pan, r.amount, r.currency,
				r.localdatetime, r.status, r.source,
				r.customer_reference, r.card_acceptor_city, r.card_acceptor_country,
				r.card_acceptor_category_code, r.card_acceptor_id, r.card_acceptor_postal_code,
				r.transmitted_at, r.response_code, r.created_at,r.updated_at,
				r.point_of_service_pin_entry_mode, r.point_of_service_pan_entry_mode,
				mr.authorization_type, mr.financial_network_code,
				mr.banknet_reference_number, mr.network_reporting_date
			FROM refunds r
			JOIN mastercard_refunds mr ON r.refund_id = mr.refund_id
			WHERE r.psp_id = @psp_id
			AND (r.transmitted_at = @transmitted_at OR @transmitted_at = '0001-01-01T00:00:00Z')
			AND (r.amount = @amount OR @amount = -1)
			AND (r.response_code = @response_code OR @response_code = '')
			ORDER BY r.transmitted_at %s
			LIMIT %v OFFSET %v`,
			filters.SortDirection(),
			filters.Limit()+1, // We fetch one extra record to see if we have fetched the last data set
			filters.Offset(),
		),
		Params: map[string]interface{}{
			"psp_id":         pspID.String(),
			"transmitted_at": mapping.MapNullTimeStr(params["processingDate"]),
			"amount":         params["amount"],
			"response_code":  params["responseCode"],
		},
	}

	var (
		refunds        []entity.Refund
		fetchedRecords = 0
	)

	ctx, cancel := context.WithTimeout(ctx, rr.readTimeout)
	defer cancel()

	iter := rr.client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return entity.Metadata{}, nil, fmt.Errorf("cannot iterate over rows: %w", err)
		}

		var mrr mastercardRefundRecord
		if err := row.ToStruct(&mrr); err != nil {
			return entity.Metadata{}, nil, fmt.Errorf("cannot parse row into struct: %w", err)
		}

		refunds = append(refunds, mapMastercardRowToRefundEntity(mrr))

		fetchedRecords++
	}

	if int(iter.RowCount) > filters.Limit() {
		refunds = refunds[:len(refunds)-1]
	}

	metadata := entity.CalculateMetadata(fetchedRecords, filters.Page, filters.PageSize)

	return metadata, refunds, nil
}

type refundRecord struct {
	ID                    string             `spanner:"refund_id"`
	Status                string             `spanner:"status"`
	MaskedPan             string             `spanner:"masked_pan"`
	PanTokenID            string             `spanner:"pan_token_id"`
	CardScheme            string             `spanner:"card_scheme"`
	Amount                int64              `spanner:"amount"`
	Currency              string             `spanner:"currency"`
	LocalDateTime         time.Time          `spanner:"localdatetime"`
	Source                string             `spanner:"source"`
	CreatedAt             time.Time          `spanner:"created_at"`
	CustomerReference     spanner.NullString `spanner:"customer_reference"`
	UpdatedAt             spanner.NullTime   `spanner:"updated_at"`
	ResponseCode          spanner.NullString `spanner:"response_code"`
	Stan                  spanner.NullInt64  `spanner:"system_trace_audit_number"`
	PspID                 string             `spanner:"psp_id"`
	TransmissionDate      spanner.NullTime   `spanner:"transmitted_at"`
	CardIssuerID          spanner.NullString `spanner:"card_issuer_id"`
	CardIssuerName        spanner.NullString `spanner:"card_issuer_name"`
	CardIssuerCountryCode spanner.NullString `spanner:"card_issuer_countrycode"`
	TransactionTypeCode   spanner.NullString `spanner:"cardholder_transaction_type_code"`
	FromAccountTypeCode   spanner.NullString `spanner:"cardholder_from_account_type_code"`
	ToAccountTypeCode     spanner.NullString `spanner:"cardholder_to_account_type_code"`
	AcceptorName          string             `spanner:"card_acceptor_name"`
	AcceptorCity          string             `spanner:"card_acceptor_city"`
	AcceptorCountry       string             `spanner:"card_acceptor_country"`
	AcceptorID            string             `spanner:"card_acceptor_id"`
	AcceptorPostalCode    spanner.NullString `spanner:"card_acceptor_postal_code"`
	CategoryCode          string             `spanner:"card_acceptor_category_code"`
	PanEntryMode          string             `spanner:"point_of_service_pan_entry_mode"`
	PinEntryMode          string             `spanner:"point_of_service_pin_entry_mode"`
	PspName               spanner.NullString `spanner:"name"`
	PspPrefix             spanner.NullString `spanner:"prefix"`
}

type mastercardRefundRecord struct {
	refundRecord
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
	CardProductID                            string             `spanner:"card_product_id"`
	CardProgramID                            string             `spanner:"card_program_id"`
}

func mapMastercardRowToRefundEntity(mrr mastercardRefundRecord) entity.Refund {
	_, schemeResponseMessage := mastercard.MapResponseCode(mrr.ResponseCode.StringVal)

	return entity.Refund{
		ID:                       uuid.MustParse(mrr.ID),
		Amount:                   int(mrr.Amount),
		Currency:                 currencycode.Must(mrr.Currency),
		CustomerReference:        mrr.CustomerReference.StringVal,
		Source:                   entity.Source(mrr.Source),
		LocalTransactionDateTime: data.LocalTransactionDateTime(mrr.LocalDateTime),
		Status:                   entity.Status(mrr.Status),
		Stan:                     int(mrr.Stan.Int64),
		ProcessingDate:           mrr.TransmissionDate.Time,
		CreatedAt:                mrr.CreatedAt,
		Card: entity.Card{
			MaskedPan:  mrr.MaskedPan,
			PanTokenID: mrr.PanTokenID,
			Info: cardinfo.Range{
				Scheme:            mrr.CardScheme,
				ProductID:         mrr.CardProductID,
				ProgramID:         mrr.CardProgramID,
				IssuerID:          mrr.CardIssuerID.StringVal,
				IssuerName:        mrr.CardIssuerName.StringVal,
				IssuerCountryCode: mrr.CardIssuerCountryCode.StringVal,
			},
		},
		CardAcceptor: entity.CardAcceptor{
			CategoryCode: mrr.CategoryCode,
			ID:           mrr.AcceptorID,
			Name:         mrr.AcceptorName,
			Address: entity.CardAcceptorAddress{
				PostalCode:  mrr.AcceptorPostalCode.StringVal,
				City:        mrr.AcceptorCity,
				CountryCode: mrr.AcceptorCountry,
			},
		},
		Psp: entity.PSP{
			ID:   uuid.MustParse(mrr.PspID),
			Name: mrr.PspName.StringVal,
		},
		CardSchemeData: entity.CardSchemeData{
			Request: entity.CardSchemeRequest{
				ProcessingCode: entity.ProcessingCode{
					TransactionTypeCode: mrr.TransactionTypeCode.StringVal,
					FromAccountTypeCode: mrr.FromAccountTypeCode.StringVal,
					ToAccountTypeCode:   mrr.ToAccountTypeCode.StringVal,
				},
				POSEntryMode: entity.POSEntryMode{
					PanEntryMode: pos.PanEntryFromCode(mrr.PanEntryMode),
					PinEntryMode: pos.PinEntryFromCode(mrr.PinEntryMode),
				},
			},
			Response: entity.CardSchemeResponse{
				Status: entity.AuthorizationStatusFromString(mrr.Status),
				ResponseCode: entity.ResponseCode{
					Value:       mrr.ResponseCode.StringVal,
					Description: schemeResponseMessage,
				},
				AuthorizationIDResponse: mrr.AuthorizationIDResponse.StringVal,
			}},
		MastercardSchemeData: entity.MastercardSchemeData{
			Request: entity.MastercardSchemeRequest{
				AuthorizationType: entity.AuthorizationType(mrr.AuthorizationType),
				PointOfServiceData: entity.PointOfServiceData{
					TerminalAttendance:                       int(mrr.TerminalAttendance.Int64),
					TerminalLocation:                         int(mrr.TerminalLocation.Int64),
					CardHolderPresence:                       int(mrr.CardHolderPresence.Int64),
					CardPresence:                             int(mrr.CardPresence.Int64),
					CardCaptureCapabilities:                  int(mrr.CardCaptureCapabilities.Int64),
					TransactionStatus:                        int(mrr.TransactionStatus.Int64),
					TransactionSecurity:                      int(mrr.TransactionSecurity.Int64),
					CardHolderActivatedTerminalLevel:         int(mrr.CardHolderActivatedTerminalLevel.Int64),
					CardDataTerminalInputCapabilityIndicator: int(mrr.CardDataTerminalInputCapabilityIndicator.Int64),
					AuthorizationLifeCycle:                   mrr.AuthorizationLifeCycle.StringVal,
					CountryCode:                              mrr.CountryCode.StringVal,
					PostalCode:                               mrr.PostalCode.StringVal,
				},
			},
			Response: entity.MastercardSchemeResponse{
				TraceID: entity.MTraceID{
					NetworkReportingDate:   mrr.NetworkReportingDate.StringVal,
					BanknetReferenceNumber: mrr.BanknetReferenceNumber.StringVal,
					FinancialNetworkCode:   mrr.FinancialNetworkCode.StringVal,
				},
			},
		},
	}
}

func mapRowToRefundEntity(r refundRecord) entity.Refund {
	_, schemeResponseMessage := mastercard.MapResponseCode(r.ResponseCode.StringVal)

	return entity.Refund{
		ID:                       uuid.MustParse(r.ID),
		Amount:                   int(r.Amount),
		Currency:                 currencycode.Must(r.Currency),
		CustomerReference:        r.CustomerReference.StringVal,
		Source:                   entity.Source(r.Source),
		LocalTransactionDateTime: data.LocalTransactionDateTime(r.LocalDateTime),
		Status:                   entity.Status(r.Status),
		Stan:                     int(r.Stan.Int64),
		ProcessingDate:           r.TransmissionDate.Time,
		CreatedAt:                r.CreatedAt,
		Card: entity.Card{
			MaskedPan:  r.MaskedPan,
			PanTokenID: r.PanTokenID,
			Info: cardinfo.Range{
				Scheme:            r.CardScheme,
				IssuerID:          r.CardIssuerID.StringVal,
				IssuerName:        r.CardIssuerName.StringVal,
				IssuerCountryCode: r.CardIssuerCountryCode.StringVal,
			},
		},
		CardAcceptor: entity.CardAcceptor{
			CategoryCode: r.CategoryCode,
			ID:           r.AcceptorID,
			Name:         r.AcceptorName,
			Address: entity.CardAcceptorAddress{
				PostalCode:  r.AcceptorPostalCode.StringVal,
				City:        r.AcceptorCity,
				CountryCode: r.AcceptorCountry,
			},
		},
		Psp: entity.PSP{
			ID:   uuid.MustParse(r.PspID),
			Name: r.PspName.StringVal,
		},
		CardSchemeData: entity.CardSchemeData{
			Request: entity.CardSchemeRequest{
				ProcessingCode: entity.ProcessingCode{
					TransactionTypeCode: r.TransactionTypeCode.StringVal,
					FromAccountTypeCode: r.FromAccountTypeCode.StringVal,
					ToAccountTypeCode:   r.ToAccountTypeCode.StringVal,
				},
				POSEntryMode: entity.POSEntryMode{
					PanEntryMode: pos.PanEntryFromCode(r.PanEntryMode),
					PinEntryMode: pos.PinEntryFromCode(r.PinEntryMode),
				},
			},
			Response: entity.CardSchemeResponse{
				Status: entity.AuthorizationStatusFromString(r.Status),
				ResponseCode: entity.ResponseCode{
					Value:       r.ResponseCode.StringVal,
					Description: schemeResponseMessage,
				},
			},
		},
	}
}
