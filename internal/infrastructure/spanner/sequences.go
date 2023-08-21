package spanner

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/sequences"
	"google.golang.org/api/iterator"
)

func NewSequenceRepository(
	name string,
	client *spanner.Client,
	readTimeout, writeTimeout time.Duration) *SequenceRepository {
	return &SequenceRepository{
		name:         name,
		client:       client,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

type SequenceRepository struct {
	name         string
	client       *spanner.Client
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (sr SequenceRepository) NextBatch(ctx context.Context, batchSize int, rolloverCond sequences.RolloverCondition, rolloverVal sequences.Block) (int, error) {

	var nextVal int
	_, err := sr.client.ReadWriteTransaction(ctx, func(ctx context.Context, trx *spanner.ReadWriteTransaction) error {
		ctx, cancelRead := context.WithTimeout(ctx, sr.readTimeout)
		defer cancelRead()

		curStmt := spanner.NewStatement("SELECT next_value, rollover_value FROM sequences WHERE name = @name")
		curStmt.Params["name"] = sr.name

		res := trx.Query(ctx, curStmt)

		defer res.Stop()
		row, err := res.Next()
		if err != nil {
			if err == iterator.Done {
				return sequences.UnknownSequence{Name: sr.name}
			}
			return fmt.Errorf("cannot select sequence %s: %w", sr.name, err)
		}

		var nextSeq sequenceRow
		err = row.ToStruct(&nextSeq)
		if err != nil {
			return fmt.Errorf("cannot parse row for sequence %s: %w", sr.name, err)
		}

		nextVal = int(nextSeq.NextValue)
		nextRolVal := nextSeq.RolloverValue
		rb := sequences.Block{
			NextValue:     nextVal + batchSize,
			RolloverValue: nextRolVal,
		}
		if rolloverCond(rb) {
			nextVal = rolloverVal.NextValue
			nextRolVal = rolloverVal.RolloverValue
		}

		ctx, cancelWrite := context.WithTimeout(ctx, sr.writeTimeout)
		defer cancelWrite()

		upStmt := spanner.NewStatement("UPDATE sequences SET next_value = @nextVal, rollover_value = @rolVal WHERE name = @name")
		upStmt.Params["name"] = sr.name
		upStmt.Params["nextVal"] = nextVal + batchSize
		upStmt.Params["rolVal"] = nextRolVal

		_, err = trx.Update(ctx, upStmt)
		if err != nil {
			return fmt.Errorf("cannot update sequence %s: %w", sr.name, err)
		}

		return nil
	})

	return nextVal, err
}

type sequenceRow struct {
	//Name          string `spanner:"name"`
	NextValue     int64  `spanner:"next_value"`
	RolloverValue string `spanner:"rollover_value"`
}
