package mastercard

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	type args struct {
		t      time.Time
		layout string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "julian date OK",
			args: args{
				layout: "j",
				t:      time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
			},
			want: "1335",
		},
		{
			name: "normal date OK",
			args: args{
				layout: "YYMMDD",
				t:      time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
			},
			want: "211201",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeFormat(tt.args.t, tt.args.layout); got != tt.want {
				t.Errorf("TimeFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeParse(t *testing.T) {
	type args struct {
		layout string
		value  string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "julian date OK",
			args: args{
				layout: "j",
				value:  "9364",
			},
			want:    time.Date(2019, time.Month(12), 30, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "normal date OK",
			args: args{
				layout: "YYMMDD",
				value:  "211211",
			},
			want:    time.Date(2021, time.Month(12), 11, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeParse(tt.args.layout, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeParse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
