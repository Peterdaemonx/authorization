package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/platform"

	captureApp "gitlab.cmpayments.local/creditcard/authorization/internal/capture/app"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization"
)

const (
	mastercard = "mastercard"
	visa       = "visa"
)

//go:generate mockgen -package=mock -source=./service.go -destination=./mock/service.go
type Tokenizer interface {
	Detokenize(ctx context.Context, merchantID string, card entity.Card) (entity.Card, error)
}

type ReversalRepository interface {
	CreateReversal(ctx context.Context, reversal entity.Reversal) error
	UpdateReversalResponse(ctx context.Context, reversal entity.Reversal) error
}

type AuthorizationRepository interface {
	GetAuthorizationWithSchemeData(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error)
	AuthorizationAlreadyReversed(ctx context.Context, id uuid.UUID) (bool, error)
}

type ReversalService struct {
	log          platform.Logger
	authRepo     AuthorizationRepository
	captureRepo  captureApp.CaptureRepository
	reversalRepo ReversalRepository
	tokenizer    Tokenizer
	mapper       *authorization.Mapper
}

func NewReversalService(
	logger platform.Logger,
	authRepo AuthorizationRepository,
	captureRepo captureApp.CaptureRepository,
	reversalRepo ReversalRepository,
	tokenizer Tokenizer,
	mapper *authorization.Mapper,
) ReversalService {
	return ReversalService{
		log:          logger,
		authRepo:     authRepo,
		captureRepo:  captureRepo,
		reversalRepo: reversalRepo,
		tokenizer:    tokenizer,
		mapper:       mapper,
	}
}

func (rs ReversalService) Reverse(ctx context.Context, pspID uuid.UUID, r *entity.Reversal) error {
	var err error
	if r.Authorization == (entity.Authorization{}) {
		r.Authorization, err = rs.authRepo.GetAuthorizationWithSchemeData(ctx, pspID, r.AuthorizationID)
		if err != nil {
			switch {
			case errors.Is(err, entity.ErrRecordNotFound):
				return entity.ErrRecordNotFound
			default:
				return fmt.Errorf("failed to fetch authorization %s for reversal %s: %w", r.AuthorizationID, r.ID, err)
			}
		}

		if r.Authorization.Status != entity.Approved {
			return entity.ErrAuthorizationNotApproved
		}
	}

	r.Authorization.Card, err = rs.tokenizer.Detokenize(ctx, pspID.String(), r.Authorization.Card)
	if err != nil {
		return fmt.Errorf("failed to detokenize pan token ID %s for reversal %s: %w", r.Authorization.Card.PanTokenID, r.ID, err)
	}

	alreadyReversed, err := rs.authRepo.AuthorizationAlreadyReversed(ctx, r.AuthorizationID)
	if err != nil {
		return fmt.Errorf("failed to fetch reversals for authorization %s: %w", r.AuthorizationID, err)
	}
	if alreadyReversed {
		return entity.ErrAuthAlreadyReversed
	}

	finalCaptureExists, err := rs.captureRepo.FinalCaptureExists(ctx, r.AuthorizationID)
	if err != nil {
		return fmt.Errorf("failed to fetch final capture for authorization %s: %w", r.AuthorizationID, err)
	}
	if finalCaptureExists {
		return entity.ErrFinalCaptureExists
	}

	summary, err := rs.captureRepo.GetCaptureSummary(ctx, r.Authorization)
	if err != nil {
		return fmt.Errorf("failed to GetCaptureSummary: %w", err)
	}
	r.Amount = r.Authorization.Amount - summary.TotalCapturedAmount

	err = rs.reversalRepo.CreateReversal(ctx, *r)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDupValOnIndex):
			return fmt.Errorf("reversal for authorization %s %w", r.AuthorizationID, err)
		default:
			return fmt.Errorf("failed to store reversal: %w", err)
		}
	}

	err = rs.mapper.SendReversal(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to send reversal to card scheme: %w", err)
	}

	err = rs.reversalRepo.UpdateReversalResponse(ctx, *r)
	if err != nil {
		return fmt.Errorf("failed to update reversal with response: %w", err)
	}

	return err
}
