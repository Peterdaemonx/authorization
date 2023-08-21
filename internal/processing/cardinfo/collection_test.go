package cardinfo

import (
	"testing"
)

func TestCollection_isBlocked(t *testing.T) {
	tests := []struct {
		name        string
		blockedbins []string
		low         string
		want        bool
	}{
		{
			name:        "low_has_correct_prefix",
			blockedbins: []string{"1234", "4321"},
			low:         "123456789",
			want:        true,
		},
		{
			name:        "low_has_incorrect_prefix",
			blockedbins: []string{"1234", "4321"},
			low:         "678912345",
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Collection{
				blockedBins: tt.blockedbins,
			}
			if got := c.isBlocked(tt.low); got != tt.want {
				t.Errorf("isBlocked() = %v, want %v", got, tt.want)
			}
		})
	}
}
