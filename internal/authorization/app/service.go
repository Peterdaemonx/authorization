package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/platform"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/risk"
	"gitlab.cmpayments.local/creditcard/authorization/internal/reversal/app"
)

//go:generate mockgen -package=mock -source=./service.go -destination=./mock/service.go
type Tokenizer interface {
	Tokenize(ctx context.Context, merchantID string, card entity.Card) (string, error)
	Detokenize(ctx context.Context, merchantID string, card entity.Card) (entity.Card, error)
}

type Repository interface {
	CreateAuthorization(ctx context.Context, a entity.Authorization) error
	CreateMastercardAuthorization(ctx context.Context, a entity.Authorization) error
	CreateVisaAuthorization(ctx context.Context, a entity.Authorization) error
	UpdateAuthorizationResponse(ctx context.Context, a entity.Authorization) error
	GetAllAuthorizations(ctx context.Context, pspID uuid.UUID, filters entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Authorization, error)
	GetAuthorizationWithSchemeData(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error)
	GetAuthorization(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error)
	AuthorizationAlreadyReversed(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateAuthorizationStatus(ctx context.Context, authorizationID uuid.UUID, status entity.Status) error
}

type AuthorizationService struct {
	log          platform.Logger
	repo         Repository
	tokenizer    Tokenizer
	revSer       app.ReversalService
	riskAssessor *risk.Assessor
	mapper       *authorization.Mapper
}

func NewAuthorizationService(
	logger platform.Logger,
	repo Repository,
	tokenizer Tokenizer,
	revSer app.ReversalService,
	mapper *authorization.Mapper,
) AuthorizationService {
	return AuthorizationService{
		log:          logger,
		repo:         repo,
		tokenizer:    tokenizer,
		revSer:       revSer,
		riskAssessor: risk.NewAssessor(risk.Rules()),
		mapper:       mapper,
	}
}

func (as AuthorizationService) Authorize(ctx context.Context, a *entity.Authorization) error {
	var err error

	a.Card.PanTokenID, err = as.tokenizer.Tokenize(ctx, a.Psp.ID.String(), a.Card)
	if err != nil {
		return fmt.Errorf("failed to tokenize consumer token: %w", err)
	}

	as.log.Info(ctx, "tokenized")

	err = as.repo.CreateAuthorization(ctx, *a)
	if err != nil {
		return fmt.Errorf("failed to store authorization: %w", err)
	}

	as.log.Info(ctx, "record inserted")

	err = as.riskAssessor.Process(ctx, a)
	if err != nil {
		return fmt.Errorf("failed analyzing risk information: %w", err)
	}
	as.log.Info(ctx, "risk assessed")

	err = as.mapper.SendAuthorization(ctx, a)
	if err != nil {
		return fmt.Errorf("failed to send authorization to card scheme: %w", err)
	}

	as.log.Info(ctx, "auth send")

	switch a.Card.Info.Scheme {
	case entity.Mastercard:
		err = as.repo.CreateMastercardAuthorization(ctx, *a)
		if err != nil {
			return fmt.Errorf("failed to store mastercard authorization: %w", err)
		}
	case entity.Visa:
		err = as.repo.CreateVisaAuthorization(ctx, *a)
		if err != nil {
			return fmt.Errorf("failed to store visa authorization: %w", err)
		}
	}

	err = as.repo.UpdateAuthorizationResponse(ctx, *a)
	if err != nil {
		return fmt.Errorf("failed to update authorization with response: %w", err)
	}

	return nil
}

//func (as AuthorizationService) updateAuthorizationAndReverse(ctx context.Context, a *entity.Authorization, err error) {
//	a.Status = entity.Failed
//	updateErr := as.repo.UpdateAuthorizationStatus(ctx, a.ID, a.Status)
//	if updateErr != nil {
//		as.log.Critical(ctx, fmt.Sprintf("failed updating authorization status to: %s, %s", a.Status, updateErr.Error()))
//	}
//
//	rev := entity.Reversal{
//		ID:              uuid.New(),
//		LogID:           uuid.MustParse(ctx.Value(logging.LogIDKey).(string)),
//		AuthorizationID: a.ID,
//		Authorization:   *a,
//		Status:          entity.ReversalNew,
//		Reason:          err,
//	}
//
//	revErr := as.revSer.Reverse(ctx, a.Psp.ID, &rev)
//	if revErr != nil {
//		as.log.Critical(ctx, fmt.Sprintf("failed reversing authorization: %s", revErr.Error()))
//	}
//}

func (as AuthorizationService) GetAuthorizations(ctx context.Context, pspID uuid.UUID, f entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Authorization, error) {
	return as.repo.GetAllAuthorizations(ctx, pspID, f, params)
}

func (as AuthorizationService) GetAuthorization(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error) {
	auth, err := as.repo.GetAuthorizationWithSchemeData(ctx, pspID, authorizationID)
	if err != nil {
		as.log.Error(ctx, fmt.Sprintf("failed to get authorization: %s", err.Error()))
		return entity.Authorization{}, err
	}
	return auth, nil
}
