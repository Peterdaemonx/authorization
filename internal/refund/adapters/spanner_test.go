package adapters

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/spanner"
)

func TestCaptureRepository_GetRefund(t *testing.T) {
	if testing.Short() {
		t.Skip("spanner: skipping integration test")
	}

	client := spanner.NewTestDB(t, "./testdata/insert_refund")

	repo := NewRefundRepository(
		client,
		3*time.Second,
		5*time.Second,
	)

	tests := []struct {
		name      string
		pspID     uuid.UUID
		refundID  uuid.UUID
		wantedErr error
	}{
		{
			name:      "get_refund",
			pspID:     uuid.MustParse("1779edcd-4f14-4c97-a61e-29a827e7ed89"),
			refundID:  uuid.MustParse("7a798829-4efc-4ad6-82ea-715307f9a8dd"),
			wantedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err := repo.GetRefund(ctx, tt.pspID, tt.refundID)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
