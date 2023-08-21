package adapters

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"cloud.google.com/go/spanner"

	infraSpanner "gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"
)

type reversalRepository struct {
	client       *spanner.Client
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewReversalRepository(
	client *spanner.Client,
	readTimeout time.Duration,
	writeTimeout time.Duration) *reversalRepository {
	return &reversalRepository{
		client:       client,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (rr *reversalRepository) CreateReversal(ctx context.Context, r entity.Reversal) error {
	_, err := rr.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		statement := spanner.Statement{
			SQL: `
				INSERT INTO authorization_reversals (
					authorization_id, reversal_id, status,
					reason, created_at, amount
				) VALUES (
					@authorization_id, @reversal_id, @status,
					@reason, @created_at, @amount
				)`,
			Params: mapCreateReversalParams(r),
		}

		ctx, cancel := context.WithTimeout(ctx, rr.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, statement)
		if spanner.ErrCode(err) == codes.AlreadyExists {
			return entity.ErrDupValOnIndex
		}

		return err
	})

	return err
}

func mapCreateReversalParams(r entity.Reversal) map[string]interface{} {
	return map[string]interface{}{
		"authorization_id": r.AuthorizationID.String(),
		"reversal_id":      r.ID.String(),
		"created_at":       time.Now(),
		"status":           r.Status,
		"amount":           r.Amount,
		"reason":           infraSpanner.NewNullError(r.Reason),
	}
}

func (rr *reversalRepository) UpdateReversalResponse(ctx context.Context, r entity.Reversal) error {
	stmt := spanner.Statement{
		SQL: `UPDATE authorization_reversals
				SET status = @status,
				    stan = @stan,
				    response_code = @response_code,
				    processing_date = @processing_date
				WHERE reversal_id = @reversal_id`,
		Params: mapUpdateReversalResponseParams(r),
	}

	_, err := rr.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, rr.writeTimeout)
		defer cancel()

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to update reversal: %w", err)
		}
		if rowCount != 1 {
			return fmt.Errorf("no record found with ID: %s", r.ID.String())
		}

		return err
	})

	return err
}

func mapUpdateReversalResponseParams(r entity.Reversal) map[string]interface{} {
	return map[string]interface{}{
		"reversal_id":     r.ID.String(),
		"status":          r.Status,
		"stan":            infraSpanner.NewNullInt64(r.Authorization.Stan),
		"response_code":   infraSpanner.NewNullString(r.CardSchemeData.Response.ResponseCode.Value),
		"processing_date": infraSpanner.NewNullTime(r.ProcessingDate),
	}
}
