package mock

import (
	"context"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
)

type CaptureRepo struct{}

func (r CaptureRepo) CreateRefundCapture(ctx context.Context, capture entity.RefundCapture) error {
	//TODO implement me
	panic("implement me")
}

func (r CaptureRepo) GetCaptureRefundSummary(ctx context.Context, refund entity.Refund) (entity.CaptureRefundSummary, error) {
	//TODO implement me
	panic("implement me")
}

func (r CaptureRepo) GetCaptureSummary(ctx context.Context, authorization entity.Authorization) (entity.CaptureSummary, error) {
	amount, final, err := r.CapturedAmountByAuthorizationID(ctx, authorization.ID)
	return entity.CaptureSummary{
		Authorization:       authorization,
		TotalCapturedAmount: amount,
		HasFinalCapture:     final,
	}, err
}

func (r CaptureRepo) CreateCapture(_ context.Context, _ entity.Capture) error {
	return nil
}

func (r CaptureRepo) GetCapturesByAuthorizationIDs(ctx context.Context, ids []string) ([]entity.Capture, error) {
	return nil, nil
}

func (r CaptureRepo) CapturedAmountByAuthorizationID(_ context.Context, _ uuid.UUID) (int, bool, error) {
	return 0, false, nil
}

func (r CaptureRepo) UpdateCapture(_ context.Context, _ entity.Capture) error {
	return nil
}

func (r CaptureRepo) FinalCaptureExists(ctx context.Context, authorizationID uuid.UUID) (bool, error) {
	return false, nil
}
