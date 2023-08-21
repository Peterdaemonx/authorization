package fetcher

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"github.com/golang/mock/gomock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/config"
	storage_mocks "gitlab.cmpayments.local/creditcard/authorization/internal/config/mocks"
)

func Test_configFetcher_FetchConfig(t *testing.T) {
	var (
		ctx        = context.Background()
		merchantID = "merchantID"
		merchant   = entity.CardAcceptor{ID: merchantID}
	)

	tests := []struct {
		name    string
		cf      configFetcher
		want    entity.CardAcceptor
		wantErr bool
	}{
		{
			name:    "should return error",
			cf:      buildConfigFetcher(t, merchant, false),
			want:    merchant,
			wantErr: false,
		},
		{
			name:    "should return error",
			cf:      buildConfigFetcher(t, merchant, true),
			want:    entity.CardAcceptor{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cf.FetchConfig(ctx, merchantID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func buildConfigFetcher(t *testing.T, merchant entity.CardAcceptor, wantErr bool) configFetcher {
	snapshotter := storage_mocks.NewMockConfigSnapshotter(gomock.NewController(t))
	if wantErr {
		snapshotter.EXPECT().SnapshotEffectiveConfig(gomock.Any()).
			Return(config.EffectiveConfig{}, errors.New("dummy error"))
	} else {
		snapshotter.EXPECT().SnapshotEffectiveConfig(gomock.Any()).
			Return(config.EffectiveConfig{Merchant: merchant}, nil)
	}

	return configFetcher{snapshotter: snapshotter}
}
