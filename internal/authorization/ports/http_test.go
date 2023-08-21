package ports

import "testing"

func TestMapRecurring(t *testing.T) {

	type args struct {
		initialRecurring bool
		traceID          string
	}

	tests := []struct {
		name string
		args
		expectedErr error
	}{
		{
			name: "initial_recurring_without_traceId",
			args: args{initialRecurring: true},
		},
		{
			name: "subsequent_rucurring_with_traceId",
			args: args{initialRecurring: false, traceID: "123456"},
		},
		{
			name: "subsequent_recurring_without_traceId",
			args: args{
				initialRecurring: false,
			},
		},
		{
			name: "initial_recurring_with_traceId",
			args: args{
				initialRecurring: true,
				traceID:          "123456",
			},
			expectedErr: ErrSubseqSpecifiedWithEmptyTraceId,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recurring, err := mapRecurring(tt.args.traceID, tt.args.initialRecurring)
			if err == nil && tt.expectedErr != nil {
				t.Errorf("got: %v", err)
			}
			if tt.expectedErr == nil {
				if recurring.Initial != tt.initialRecurring {
					t.Errorf("expected: %v, got: %v", tt.initialRecurring, recurring.Initial)
				}
				if recurring.TraceID != tt.traceID {
					t.Errorf("expected: %v, got: %v", tt.traceID, recurring.TraceID)
				}
			}
		})
	}
}
