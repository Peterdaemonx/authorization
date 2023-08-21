package mock

import (
	"context"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
)

const (
	tokenID = "cd89ecb2-f50b-466f-a1a5-7b7f4d9a58d1"
)

type Tokenization struct {
}

func (t Tokenization) Tokenize(_ context.Context, _ string, _ entity.Card) (string, error) {
	return tokenID, nil
}

func (t Tokenization) Detokenize(ctx context.Context, merchantID string, card entity.Card) (entity.Card, error) {
	if card.PanTokenID == tokenID {
		return entity.Card{
			Number:     "5413330002001411",
			Holder:     "a holder",
			Expiry:     entity.Expiry{Year: time.Now().Add(time.Hour * 24 * 730).Format("06"), Month: "03"},
			PanTokenID: card.PanTokenID,
			Info:       card.Info,
		}, nil
	}
	return entity.Card{}, processing.ErrDetokenization
}
