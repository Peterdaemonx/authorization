package adapters

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
)

type captureRepository struct {
	client       *spanner.Client
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewCaptureRepository(
	client *spanner.Client,
	readTimeout time.Duration,
	writeTimeout time.Duration) *captureRepository {
	return &captureRepository{
		client:       client,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (c captureRepository) UpdateCapture(ctx context.Context, capture entity.Capture) error {
	stmt := spanner.NewStatement(`
		UPDATE authorization_captures set status=@status, updated_at=@updated_at, ird=@ird
		WHERE authorization_id=@authorization_id AND capture_id=@capture_id
	`)
	stmt.Params["status"] = capture.Status.String()
	stmt.Params["updated_at"] = time.Now()
	stmt.Params["ird"] = capture.IRD
	stmt.Params["authorization_id"] = capture.AuthorizationID.String()
	stmt.Params["capture_id"] = capture.ID.String()

	_, err := c.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, c.writeTimeout)
		defer cancel()

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		if rowCount != 1 {
			return fmt.Errorf("no capture was updated by Capture ID %v and Authorization ID %s", capture.ID, capture.AuthorizationID)
		}

		return err
	})

	return err

}

func (c captureRepository) CreateCapture(ctx context.Context, capture entity.Capture) error {
	stmt := spanner.NewStatement(`
				INSERT INTO authorization_captures (
					authorization_id, capture_id, amount,
				    currency, is_final, reference,
				    status, created_at
				) VALUES (
					@authorization_id, @capture_id, @amount,
					@currency, @is_final, @reference,
				    @status, @created_at
				)`)
	stmt.Params["authorization_id"] = capture.AuthorizationID.String()
	stmt.Params["capture_id"] = capture.ID.String()
	stmt.Params["amount"] = capture.Amount
	stmt.Params["currency"] = capture.Currency.Alpha3()
	stmt.Params["is_final"] = capture.IsFinal
	stmt.Params["reference"] = capture.Reference
	stmt.Params["status"] = capture.Status.String()
	stmt.Params["created_at"] = time.Now()

	_, err := c.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, c.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, stmt)

		return err
	})

	return err
}

func (c captureRepository) GetCaptureSummary(ctx context.Context, authorization entity.Authorization) (entity.CaptureSummary, error) {
	stmt := spanner.NewStatement(
		`
		SELECT SUM(amount) as total_amount, MAX(is_final) as is_final
		FROM authorization_captures
		WHERE authorization_id = @authorization_id AND status = @status
		GROUP BY authorization_id
		`)
	stmt.Params["authorization_id"] = authorization.ID.String()
	stmt.Params["status"] = entity.CaptureCreated.String()

	ctx, cancel := context.WithTimeout(ctx, c.readTimeout)
	defer cancel()

	var totalAmount int
	var isFinal bool
	iter := c.client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return entity.CaptureSummary{}, err
		}

		var amount spanner.NullInt64
		err = row.ColumnByName("total_amount", &amount)
		if err != nil {
			return entity.CaptureSummary{}, err
		}
		if amount.Valid {
			totalAmount = totalAmount + int(amount.Int64)
		}

		var final bool
		err = row.ColumnByName("is_final", &final)
		if err != nil {
			return entity.CaptureSummary{}, err
		}
		if amount.Valid {
			isFinal = final
		}
	}

	return entity.CaptureSummary{
		Authorization:       authorization,
		TotalCapturedAmount: totalAmount,
		HasFinalCapture:     isFinal,
	}, nil
}

func (c captureRepository) CreateRefundCapture(ctx context.Context, capture entity.RefundCapture) error {
	stmt := spanner.NewStatement(`
				INSERT INTO refund_captures (
					refund_id, capture_id, amount,
				    currency, is_final, reference,
				    status, created_at
				) VALUES (
					@refund_id, @capture_id, @amount,
					@currency, @is_final, @reference,
				    @status, @created_at
				)`)
	stmt.Params["refund_id"] = capture.RefundID.String()
	stmt.Params["capture_id"] = capture.ID.String()
	stmt.Params["amount"] = capture.Amount
	stmt.Params["currency"] = capture.Currency.Alpha3()
	stmt.Params["is_final"] = capture.IsFinal
	stmt.Params["reference"] = capture.Reference
	stmt.Params["status"] = capture.Status.String()
	stmt.Params["created_at"] = time.Now()

	_, err := c.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, c.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, stmt)

		return err
	})

	return err
}

func (c captureRepository) GetCaptureRefundSummary(ctx context.Context, refund entity.Refund) (entity.CaptureRefundSummary, error) {
	stmt := spanner.NewStatement(
		`
		SELECT SUM(amount) as total_amount, MAX(is_final) as is_final
		FROM refund_captures
		WHERE refund_id = @refund_id AND status = @status
		GROUP BY refund_id
		`)
	stmt.Params["refund_id"] = refund.ID.String()
	stmt.Params["status"] = entity.CaptureCreated.String()

	ctx, cancel := context.WithTimeout(ctx, c.readTimeout)
	defer cancel()

	var totalAmount int
	var isFinal bool
	iter := c.client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done || err != nil {
			break
		}

		var amount spanner.NullInt64
		err = row.ColumnByName("total_amount", &amount)
		if err != nil {
			return entity.CaptureRefundSummary{}, err
		}
		if amount.Valid {
			totalAmount = totalAmount + int(amount.Int64)
		}

		var final bool
		err = row.ColumnByName("is_final", &final)
		if err != nil {
			return entity.CaptureRefundSummary{}, err
		}
		if amount.Valid {
			isFinal = final
		}

		if err != nil {
			return entity.CaptureRefundSummary{}, err
		}
	}

	return entity.CaptureRefundSummary{
		Refund:              refund,
		TotalCapturedAmount: totalAmount,
		HasFinalCapture:     isFinal,
	}, nil
}

func (c captureRepository) GetCapturesByAuthorizationIDs(ctx context.Context, ids []string) ([]entity.Capture, error) {
	var captures []entity.Capture

	stmt := spanner.Statement{
		SQL: `
			SELECT ac.authorization_id, ac.capture_id, ac.amount, ac.currency, ac.is_final, ac.status
			FROM authorizations AS a
			INNER JOIN authorization_captures AS ac
			ON a.authorization_id = ac.authorization_id
			AND a.authorization_id IN UNNEST (@ids)
			ORDER BY ac.created_at;
		`,
		Params: map[string]interface{}{"ids": ids},
	}

	ctx, cancel := context.WithTimeout(ctx, c.readTimeout)
	defer cancel()

	iter := c.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return captures, fmt.Errorf("failed to iterate over rows: %w", err)
		}

		var capture captureRecord
		if err = row.ToStruct(&capture); err != nil {
			return captures, err
		}

		captures = append(captures, mapCaptureRecordToCaptureEntity(capture))
	}

	if len(captures) == 0 {
		captures = []entity.Capture{}
	}

	return captures, nil
}

func mapCaptureRecordToCaptureEntity(c captureRecord) entity.Capture {
	capture := entity.Capture{
		AuthorizationID: uuid.MustParse(c.AuthorizationID),
		ID:              uuid.MustParse(c.CaptureID),
		Amount:          int(c.Amount),
		Currency:        currencycode.Must(c.Currency),
		IsFinal:         c.IsFinal,
		Reference:       c.Reference,
		Status:          entity.CaptureStatusFromString(c.Status),
	}

	if c.UpdatedAt.Valid {
		capture.UpdatedAt = c.UpdatedAt.Time
	}

	if c.IRD.Valid {
		capture.IRD = c.IRD.StringVal
	}

	return capture
}

type captureRecord struct {
	AuthorizationID string             `spanner:"authorization_id"`
	CaptureID       string             `spanner:"capture_id"`
	Amount          int64              `spanner:"amount"`
	Currency        string             `spanner:"currency"`
	IsFinal         bool               `spanner:"is_final"`
	Reference       string             `spanner:"reference"`
	Status          string             `spanner:"status"`
	IRD             spanner.NullString `spanner:"ird"`
	UpdatedAt       spanner.NullTime   `spanner:"updated_at"`
}

func (c captureRepository) FinalCaptureExists(ctx context.Context, authorizationID uuid.UUID) (bool, error) {
	stmt := spanner.Statement{
		SQL: `
			    SELECT 1
				FROM authorization_captures
				WHERE authorization_id = @authorization_id
				AND is_final = TRUE
		`,
		Params: map[string]interface{}{"authorization_id": authorizationID.String()},
	}

	ctx, cancel := context.WithTimeout(ctx, c.readTimeout)
	defer cancel()

	iter := c.client.Single().Query(ctx, stmt)
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
