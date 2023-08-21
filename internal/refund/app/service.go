package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/platform"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/authorization"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/risk"
)

type Tokenizer interface {
	Tokenize(ctx context.Context, merchantID string, card entity.Card) (string, error)
	Detokenize(ctx context.Context, merchantID string, card entity.Card) (entity.Card, error)
}

type Repository interface {
	CreateRefund(ctx context.Context, r entity.Refund) error
	CreateMastercardRefund(ctx context.Context, r entity.Refund) error
	CreateVisaRefund(ctx context.Context, r entity.Refund) error
	GetRefund(ctx context.Context, pspID, refundID uuid.UUID) (entity.Refund, error)
	UpdateRefundResponse(ctx context.Context, r entity.Refund) error
	GetAllRefunds(ctx context.Context, pspID uuid.UUID, filters entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Refund, error)
}

type RefundService struct {
	log          platform.Logger
	repo         Repository
	tokenizer    Tokenizer
	riskAssessor *risk.Assessor
	mapper       *authorization.Mapper
}

func NewRefundService(
	logger platform.Logger,
	repo Repository,
	tokenizer Tokenizer,
	mapper *authorization.Mapper,
) RefundService {
	return RefundService{
		log:          logger,
		repo:         repo,
		tokenizer:    tokenizer,
		riskAssessor: risk.NewAssessor(risk.Rules()),
		mapper:       mapper,
	}
}

func (rs RefundService) Authorize(ctx context.Context, r *entity.Refund) error {
	var err error

	r.Card.PanTokenID, err = rs.tokenizer.Tokenize(ctx, r.Psp.ID.String(), r.Card)
	if err != nil {
		return fmt.Errorf("failed to tokenize consumer token: %w", err)
	}

	rs.log.Info(ctx, "tokenized")

	err = rs.repo.CreateRefund(ctx, *r)
	if err != nil {
		return fmt.Errorf("failed to store refund: %w", err)
	}

	rs.log.Info(ctx, "record inserted")

	err = rs.mapper.SendRefund(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to send refund to card scheme: %w", err)
	}

	rs.log.Info(ctx, "refund send")

	switch r.Card.Info.Scheme {
	case entity.Mastercard:
		err = rs.repo.CreateMastercardRefund(ctx, *r)
		if err != nil {
			return fmt.Errorf("failed to store mastercard refund: %w", err)
		}
	case entity.Visa:
		err = rs.repo.CreateVisaRefund(ctx, *r)
		if err != nil {
			return fmt.Errorf("failed to store mastercard refund: %w", err)
		}
	}

	err = rs.repo.UpdateRefundResponse(ctx, *r)
	if err != nil {
		return fmt.Errorf("failed to update refund with response: %w", err)
	}

	return nil
}

func (rs RefundService) GetRefunds(ctx context.Context, pspID uuid.UUID, f entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Refund, error) {
	return rs.repo.GetAllRefunds(ctx, pspID, f, params)
}
