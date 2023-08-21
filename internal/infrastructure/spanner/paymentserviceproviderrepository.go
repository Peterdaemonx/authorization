package spanner

import (
	"context"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type PaymentServiceProviderRepository struct {
	client       *spanner.Client
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type PspRecord struct {
	ID     string `spanner:"psp_id"`
	Name   string `spanner:"name"`
	Prefix string `spanner:"prefix"`
}

func NewPaymentServiceProviderRepository(
	client *spanner.Client,
	readTimeout time.Duration,
	writeTimeout time.Duration) *PaymentServiceProviderRepository {
	return &PaymentServiceProviderRepository{
		client:       client,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (p PaymentServiceProviderRepository) GetPspByID(ctx context.Context, id uuid.UUID) (entity.PSP, error) {
	const sqlQuery = `
    	SELECT psp.psp_id, psp.name, psp.prefix 
		FROM psp
		WHERE psp.psp_id = @id
    `
	stmt := spanner.Statement{
		SQL:    sqlQuery,
		Params: map[string]interface{}{"id": id.String()},
	}
	ctx, cancel := context.WithTimeout(ctx, p.readTimeout)
	defer cancel()

	iter := p.client.Single().Query(ctx, stmt)
	var psp Psp
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return mapRecordToDto(psp), nil
		}
		if err != nil {
			return entity.PSP{}, err
		}
		if err := row.ToStruct(&psp); err != nil {
			return entity.PSP{}, err
		}
	}
}

func (p PaymentServiceProviderRepository) GetPspByAPIKey(ctx context.Context, apiKey string) (entity.PSP, error) {
	const sqlQuery = `
    	SELECT psp.psp_id, psp.name, psp.prefix 
		FROM psp
		JOIN api_consumers ac on psp.psp_id = ac.psp_id
		WHERE ac.api_key = @api_key
    `
	stmt := spanner.Statement{
		SQL:    sqlQuery,
		Params: map[string]interface{}{"api_key": apiKey},
	}
	ctx, cancel := context.WithTimeout(ctx, p.readTimeout)
	defer cancel()

	iter := p.client.Single().Query(ctx, stmt)
	var psp Psp
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return mapRecordToDto(psp), nil
		}
		if err != nil {
			return entity.PSP{}, err
		}
		if err := row.ToStruct(&psp); err != nil {
			return entity.PSP{}, err
		}
	}
}

type Psp struct {
	Id     string `spanner:"psp_id"`
	Name   string `spanner:"name"`
	Prefix string `spanner:"prefix"`
}

func mapRecordToDto(psp Psp) entity.PSP {
	id, err := uuid.Parse(psp.Id)
	if err != nil {
		id = uuid.Nil
	}
	return entity.PSP{
		ID:     id,
		Name:   psp.Name,
		Prefix: psp.Prefix,
	}
}
