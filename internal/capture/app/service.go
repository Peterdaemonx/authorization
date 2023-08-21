package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/platform/events/pubsub"

	"gitlab.cmpayments.local/creditcard/authorization/internal/capture/adapters"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

//go:generate mockgen -package=mock -source=./service.go -destination=./mock/service.go
type CaptureRepository interface {
	CreateCapture(ctx context.Context, capture entity.Capture) error
	CreateRefundCapture(ctx context.Context, capture entity.RefundCapture) error
	GetCapturesByAuthorizationIDs(ctx context.Context, ids []string) ([]entity.Capture, error)
	UpdateCapture(ctx context.Context, capture entity.Capture) error
	GetCaptureSummary(ctx context.Context, authorization entity.Authorization) (entity.CaptureSummary, error)
	GetCaptureRefundSummary(ctx context.Context, refund entity.Refund) (entity.CaptureRefundSummary, error)
	FinalCaptureExists(ctx context.Context, authorizationID uuid.UUID) (bool, error)
}

type AuthorizationRepository interface {
	GetAuthorizationWithSchemeData(ctx context.Context, pspID, authorizationID uuid.UUID) (entity.Authorization, error)
	AuthorizationAlreadyReversed(ctx context.Context, id uuid.UUID) (bool, error)
}

type RefundRepository interface {
	GetRefund(ctx context.Context, pspID, refundID uuid.UUID) (entity.Refund, error)
}
type CaptureService struct {
	authRepo    AuthorizationRepository
	refundRepo  RefundRepository
	captureRepo CaptureRepository
	publisher   pubsub.Publisher
	authTopic   string
	refundTopic string
}

func NewCaptureService(authRepo AuthorizationRepository, refundRepo RefundRepository, captureRepo CaptureRepository, publisher pubsub.Publisher, authTopic, refundTopic string) CaptureService {
	return CaptureService{
		authRepo:    authRepo,
		refundRepo:  refundRepo,
		captureRepo: captureRepo,
		publisher:   publisher,
		authTopic:   authTopic,
		refundTopic: refundTopic,
	}
}

func (cs CaptureService) Capture(ctx context.Context, pspID uuid.UUID, c entity.Capture) error {
	a, err := cs.authRepo.GetAuthorizationWithSchemeData(ctx, pspID, c.AuthorizationID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			return entity.ErrRecordNotFound
		default:
			return fmt.Errorf("failed to fetch authorization %s for capture %s: %w", c.AuthorizationID, c.ID, err)
		}
	}

	summary, err := cs.captureRepo.GetCaptureSummary(ctx, a)
	if err != nil {
		return fmt.Errorf("failed to GetCaptureSummary: %w", err)
	}

	validationErr := summary.ValidateExtraCapture(c.Amount)
	if validationErr != nil {
		return validationErr
	}

	// If the PSP doesn't request a Final Capture, but based on the amount they should,
	// treat it as Final anyway.
	if !c.IsFinal {
		c.IsFinal = summary.IsFinalizedWith(c.Amount)
	}

	alreadyReversed, err := cs.authRepo.AuthorizationAlreadyReversed(ctx, a.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch reversals for authorization %s: %w", a.ID, err)
	}
	if alreadyReversed {
		return entity.ErrAuthAlreadyReversed
	}

	err = cs.captureRepo.CreateCapture(ctx, c)
	if err != nil {
		return fmt.Errorf("failed to store capture: %w", err)
	}

	err = cs.publisher.Publish(ctx, cs.authTopic, adapters.CreatePublishAuthorizationRequest(a, c))
	if err != nil {
		return fmt.Errorf("failed to publish capture: %w", err)
	}

	return nil
}

func (cs CaptureService) CaptureRefund(ctx context.Context, pspID uuid.UUID, c entity.RefundCapture) error {
	r, err := cs.refundRepo.GetRefund(ctx, pspID, c.RefundID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			return entity.ErrRecordNotFound
		default:
			return fmt.Errorf("failed to fetch refund %s for capture %s: %w", c.RefundID, c.ID, err)
		}
	}

	summary, err := cs.captureRepo.GetCaptureRefundSummary(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to GetCaptureSummary: %w", err)
	}

	validationErr := summary.ValidateExtraCapture(c.Amount)
	if validationErr != nil {
		return validationErr
	}

	// If the PSP doesn't request a Final Capture, but based on the amount they should,
	// treat it as Final anyway.
	if !c.IsFinal {
		c.IsFinal = summary.IsFinalizedWith(c.Amount)
	}

	err = cs.captureRepo.CreateRefundCapture(ctx, c)
	if err != nil {
		return fmt.Errorf("failed to store refund capture: %w", err)
	}

	//TODO: This change needs some changes in clearing. Validations block the capture of refunds.
	err = cs.publisher.Publish(ctx, cs.refundTopic, adapters.CreatePublishRefundRequest(r, c))
	if err != nil {
		return fmt.Errorf("failed to publish refund capture: %w", err)
	}

	return nil
}

func (cs CaptureService) GetCapturesByAuthorizationIDs(ctx context.Context, ids []string) ([]entity.Capture, error) {
	return cs.captureRepo.GetCapturesByAuthorizationIDs(ctx, ids)
}
