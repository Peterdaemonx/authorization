//go:build integration
// +build integration

package spanner

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_nonceRepository_StoreNonce(t *testing.T) {
	ctx := context.Background()
	pspID := uuid.MustParse("0cd8d732-66c2-4dae-bb99-16494dea7796")
	nonce, _ := generateNonce()
	otherNonce, _ := generateNonce()
	repository := connectToNonceRepository()

	err := repository.StoreNonce(context.Background(), pspID, otherNonce)
	if err != nil {
		t.Error(t, err, fmt.Sprintf("StoreNonce(%v, %v, %v)", ctx, pspID, otherNonce))
	}

	type args struct {
		pspID uuid.UUID
		nonce string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "store_valid_nonce",
			args: args{
				pspID: pspID,
				nonce: nonce,
			},
			wantErr: false,
		},
		{
			name: "store_duplicate_nonce",
			args: args{
				pspID: pspID,
				nonce: otherNonce,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.StoreNonce(context.Background(), tt.args.pspID, tt.args.nonce)
			if (err != nil) != tt.wantErr {
				t.Error(t, err, fmt.Sprintf("StoreNonce(%v, %v, %v)", ctx, tt.args.pspID, tt.args.nonce))
			}
		})
	}
}

func connectToNonceRepository() *nonceRepository {
	timeout, _ := time.ParseDuration("10s")
	return NewNonceRepository(Client, timeout, timeout)
}

func generateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}
