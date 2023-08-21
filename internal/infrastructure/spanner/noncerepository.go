package spanner

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/web"
	"google.golang.org/grpc/codes"

	"cloud.google.com/go/spanner"
)

type nonceRepository struct {
	client       *spanner.Client
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewNonceRepository(
	client *spanner.Client,
	readTimeout time.Duration,
	writeTimeout time.Duration) *nonceRepository {
	return &nonceRepository{
		client:       client,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (r nonceRepository) StoreNonce(ctx context.Context, pspID uuid.UUID, nonce string) error {

	statement := spanner.NewStatement("INSERT INTO nonces (psp_id, nonce, used_at) VALUES (@psp_id, @nonce, CURRENT_TIMESTAMP())")
	statement.Params["psp_id"] = pspID.String()
	statement.Params["nonce"] = nonce

	_, err := r.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		ctx, cancel := context.WithTimeout(ctx, r.writeTimeout)
		defer cancel()

		_, err := txn.Update(ctx, statement)

		if spanner.ErrCode(err) == codes.AlreadyExists {
			return web.ErrNonceAlreadyUsed
		}

		return err
	})

	return err
}
