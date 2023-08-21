package mock

import (
	"context"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

type RefundRepo struct{}

func (rs RefundRepo) CreateRefund(ctx context.Context, r entity.Refund) error {
	//TODO implement me
	panic("implement me")
}

func (rs RefundRepo) CreateMastercardRefund(ctx context.Context, r entity.Refund) error {
	//TODO implement me
	panic("implement me")
}

func (rs RefundRepo) CreateVisaRefund(ctx context.Context, r entity.Refund) error {
	//TODO implement me
	panic("implement me")
}

func (rs RefundRepo) UpdateRefundResponse(ctx context.Context, r entity.Refund) error {
	//TODO implement me
	panic("implement me")
}

func (rs RefundRepo) GetAllRefunds(ctx context.Context, pspID uuid.UUID, filters entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Refund, error) {
	//TODO implement me
	panic("implement me")
}

func (rs RefundRepo) GetRefund(ctx context.Context, pspID, refundID uuid.UUID) (entity.Refund, error) {
	//TODO implement me
	panic("implement me")
}
