package visa

import (
	"testing"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/visa/base1"
)

func TestMessageFromReversal(t *testing.T) {
	const (
	//acceptorName = "MaxCorp Inc."
	//city         = "Breda"
	)
	var ()

	type args struct {
		refund entity.Refund
	}

	tests := []struct {
		name string
		args args
		want *Message
	}{
		{
			name: "test_valid_message_from_reversal",
			args: args{
				entity.Refund{},
			},
			want: &Message{
				Mti:    iso8583.MTI{},
				Fields: base1.Fields{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
