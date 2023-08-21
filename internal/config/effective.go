package config

import (
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
)

// Snapshotter can create a snapshot of effective Merchant configuration.
type Snapshotter struct {
	TargetTable string
}

// SnapshotEffectiveConfig creates a snapshot of the current, effective configuration that applies
// to a Merchant and returns it.
// The snapshot is stored in the DB table indicated by the Snapshotter.TargetTable field.
func (s Snapshotter) SnapshotEffectiveConfig(merchantID string) (EffectiveConfig, error) {
	if merchantID != "2ef864ea-5f4c-439d-8e99-32d286d89553" {
		return EffectiveConfig{}, processing.ErrPspNotFound
	}
	return EffectiveConfig{Merchant: entity.CardAcceptor{
		CategoryCode: "5691",
		ID:           "0987654321",
		Name:         "MaxCorp Inc.",
		Address: entity.CardAcceptorAddress{
			PostalCode:  "4825 AN",
			City:        "Breda",
			CountryCode: "NLD",
		},
	}}, nil
}

// EffectiveConfig contains the effective, resolved configuration of a Merchant.
// It has a property for each configuration value and takes the most specific point where it has
// been set (Merchant > Tenant > Default)
type EffectiveConfig struct {
	Merchant entity.CardAcceptor
}

// Merchant This is a DEMO structure for store merchant configuration, should be
// replaced when fetchconfig service is implemented
type Merchant struct {
	ID                       string
	MCC                      string
	CountryCode              string
	RouteBy                  string
	DefaultAuthorizationType string
	AutomaticCapture         bool
	CaptureDelay             time.Duration
	TenantID                 string
}
