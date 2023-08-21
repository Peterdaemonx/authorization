package risk

import (
	"context"
	"testing"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
)

func Test_lowValueExemptionRule_Support(t *testing.T) {
	type args struct {
		trx interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "rule should be supported",
			args: args{trx: entity.Authorization{
				Exemption: entity.LowValueExemption,
			},
			},
			want: true,
		},
		{
			name: "rule should not be supported when transaction type is different of authorization",
			args: args{trx: processing.ThreeDSecure{}},
			want: false,
		},
		{
			name: "rule should not be supported when exemption type is different from low value",
			args: args{trx: entity.Authorization{
				Exemption: entity.MerchantInitiatedExemption,
			},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := lowValueExemption{}
			if got := v.Support(tt.args.trx); got != tt.want {
				t.Errorf("Support() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lowValueExemptionRule_Eval(t *testing.T) {
	type args struct {
		in0 context.Context
		trx interface{}
		in2 *Assessment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should no return error when amount is less than 30",
			args: args{trx: entity.Authorization{
				Amount:    30,
				Exemption: entity.LowValueExemption,
			},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := lowValueExemption{}
			if err := v.Eval(tt.args.in0, tt.args.trx, tt.args.in2); (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
