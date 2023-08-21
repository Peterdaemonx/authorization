// Code generated by timing/wrappers/generate, DO NOT EDIT.
// Generated on Mon May 15 15:07 2023
package timingwrappers

import (
	"context"
	uuid "github.com/google/uuid"
	app "gitlab.cmpayments.local/creditcard/authorization/internal/capture/app"
	entity "gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	timing "gitlab.cmpayments.local/creditcard/authorization/internal/timing"
)

type CaptureRepository struct {
	Base app.CaptureRepository
}

func (w CaptureRepository) CreateCapture(ctx context.Context, capture entity.Capture) error {
	timing.Start(ctx, "CaptureRepository.CreateCapture")
	defer timing.Stop(ctx, "CaptureRepository.CreateCapture")
	return w.Base.CreateCapture(ctx, capture)
}
func (w CaptureRepository) CreateRefundCapture(ctx context.Context, capture entity.RefundCapture) error {
	timing.Start(ctx, "CaptureRepository.CreateRefundCapture")
	defer timing.Stop(ctx, "CaptureRepository.CreateRefundCapture")
	return w.Base.CreateRefundCapture(ctx, capture)
}
func (w CaptureRepository) FinalCaptureExists(ctx context.Context, authorizationID uuid.UUID) (bool, error) {
	timing.Start(ctx, "CaptureRepository.FinalCaptureExists")
	defer timing.Stop(ctx, "CaptureRepository.FinalCaptureExists")
	return w.Base.FinalCaptureExists(ctx, authorizationID)
}
func (w CaptureRepository) GetCaptureRefundSummary(ctx context.Context, refund entity.Refund) (entity.CaptureRefundSummary, error) {
	timing.Start(ctx, "CaptureRepository.GetCaptureRefundSummary")
	defer timing.Stop(ctx, "CaptureRepository.GetCaptureRefundSummary")
	return w.Base.GetCaptureRefundSummary(ctx, refund)
}
func (w CaptureRepository) GetCaptureSummary(ctx context.Context, authorization entity.Authorization) (entity.CaptureSummary, error) {
	timing.Start(ctx, "CaptureRepository.GetCaptureSummary")
	defer timing.Stop(ctx, "CaptureRepository.GetCaptureSummary")
	return w.Base.GetCaptureSummary(ctx, authorization)
}
func (w CaptureRepository) GetCapturesByAuthorizationIDs(ctx context.Context, ids []string) ([]entity.Capture, error) {
	timing.Start(ctx, "CaptureRepository.GetCapturesByAuthorizationIDs")
	defer timing.Stop(ctx, "CaptureRepository.GetCapturesByAuthorizationIDs")
	return w.Base.GetCapturesByAuthorizationIDs(ctx, ids)
}
func (w CaptureRepository) UpdateCapture(ctx context.Context, capture entity.Capture) error {
	timing.Start(ctx, "CaptureRepository.UpdateCapture")
	defer timing.Stop(ctx, "CaptureRepository.UpdateCapture")
	return w.Base.UpdateCapture(ctx, capture)
}
