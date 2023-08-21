package tokenization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gitlab.cmpayments.local/creditcard/platform"
	"gitlab.cmpayments.local/creditcard/tokenization-api/pkg/tokenization"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

var (
	ErrFailedTokenize   = errors.New("failed to tokenize")
	ErrFailedDetokenize = errors.New("failed to detokenize")
)

func NewService(tokenizationClient tokenization.Client, detokenizationClient tokenization.Client, logger platform.Logger) *service {
	return &service{tokenizationClient, detokenizationClient, logger}
}

type service struct {
	tokenizationClient   tokenization.Client
	detokenizationClient tokenization.Client
	logger               platform.Logger
}

func (s service) Tokenize(ctx context.Context, merchantID string, card entity.Card) (string, error) {
	req := tokenization.InputCcard{
		Number: card.Number,
	}

	if card.Holder != "" {
		req.Holder = &card.Holder
	}

	if (card.Expiry != entity.Expiry{}) {
		req.Expiry = &tokenization.Expiry{
			Year:  card.Expiry.MustYearToInt(),
			Month: card.Expiry.MustMonthToInt(),
		}
	}

	res, err := s.tokenizationClient.TokenizeCcard(ctx, merchantID, req)
	if err != nil {
		return "", handleError("tokenizing", err)
	}

	return res, nil
}

func (s service) Detokenize(ctx context.Context, merchantID string, card entity.Card) (entity.Card, error) {
	res, err := s.detokenizationClient.DetokenizeCcard(ctx, merchantID, card.PanTokenID)
	if err != nil {
		return entity.Card{}, handleError("detokenizing", err)
	}

	detokenizedCard := entity.Card{
		Number:     res.Pan,
		Holder:     *res.Holder,
		PanTokenID: card.PanTokenID,
		Info:       card.Info,
	}

	if res.Expiry != nil {
		detokenizedCard.Expiry = entity.Expiry{
			Year:  strconv.Itoa(res.Expiry.Year % 100),
			Month: strconv.Itoa(res.Expiry.Month),
		}
	}

	return detokenizedCard, nil
}

func handleError(service string, err error) error {
	var errJson *json.SyntaxError
	var errResponse error
	switch {
	case errors.As(err, &errJson):
		errResponse = fmt.Errorf("cannot execute HTTP request: %v", err)
	default:
		errResponse = fmt.Errorf("unknown error has occurred %s credit card: %v", service, err)
	}
	return fmt.Errorf("%w: %v", ErrFailedTokenize, errResponse)
}
