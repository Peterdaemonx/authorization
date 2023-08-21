//go:build integration
// +build integration

package spanner

import (
	"context"
	"testing"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"github.com/google/uuid"
)

func TestPaymentServiceProviderRepository_GetPspByID(t *testing.T) {
	repository := connectToPaymentServiceProviderRepository()
	tests := []struct {
		name string
		arg  string
		psp  entity.PSP
	}{
		{
			name: "get_correct_psp_by_id",
			arg:  "0cd8d732-66c2-4dae-bb99-16494dea7796",
			psp:  entity.PSP{ID: uuid.MustParse("0cd8d732-66c2-4dae-bb99-16494dea7796"), Name: "test_psp", Prefix: "tps"},
		},
		{
			name: "get_psp_without_prefix",
			arg:  "0cd8d732-66c2-4dae-bb99-16494dea7796",
			psp:  entity.PSP{ID: uuid.MustParse("0cd8d732-66c2-4dae-bb99-16494dea7796"), Name: "test_psp"},
		},
		{
			name: "incorrect_id_no_psp",
			arg:  "",
			psp:  entity.PSP{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.arg != "" {
				id, err := uuid.Parse(tt.arg)
				if err != nil {
					t.Errorf("failed parsing string to uuid: %v", err)
				}
				psp, err := repository.GetPspByID(context.Background(), id)
				if err != nil {
					t.Errorf("got: %v", err)
				}
				if tt.arg != "" && psp == (entity.PSP{}) {
					t.Errorf("got: %v, wanted: %v", psp.ID, tt.psp.ID)
				}
			}
		})
	}
}

func TestPaymentServiceProviderRepository_GetPspByAPIKey(t *testing.T) {
	repository := connectToPaymentServiceProviderRepository()
	tests := []struct {
		name string
		arg  string
		psp  entity.PSP
	}{
		{
			name: "get correct psp for api key",
			arg:  "6247c10c-84a0-4fa1-b330-77eea1e944d3",
			psp:  entity.PSP{ID: uuid.MustParse("0cd8d732-66c2-4dae-bb99-16494dea7796"), Name: "test_psp", Prefix: "tps"},
		},
		{
			name: "incorrect api key results in no psp.",
			arg:  "",
			psp:  entity.PSP{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			psp, err := repository.GetPspByAPIKey(context.Background(), tt.arg)
			if err != nil {
				t.Errorf("got: %v", err)
			}
			if tt.arg != "" && psp == (entity.PSP{}) {
				t.Errorf("got: %v, wanted: %v", psp.ID.String(), tt.psp.ID.String())
			}
		})
	}
}

func connectToPaymentServiceProviderRepository() *PaymentServiceProviderRepository {
	wtimeout, _ := time.ParseDuration("10s")
	rtimeout, _ := time.ParseDuration("10s")
	return NewPaymentServiceProviderRepository(Client, rtimeout, wtimeout)
}
