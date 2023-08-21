package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrAuthorizationNotApproved = errors.New("failed to reverse authorization, authorization was not approved")
	ErrDupValOnIndex            = errors.New("already exists")
)

type ReversalStatus string

const (
	ReversalNew       ReversalStatus = "new"
	ReversalFailed    ReversalStatus = "failed"
	ReversalSucceeded ReversalStatus = "succeeded"
)

type Reversal struct {
	ID              uuid.UUID
	LogID           uuid.UUID
	Status          ReversalStatus
	AuthorizationID uuid.UUID
	Authorization   Authorization
	CardSchemeData  CardSchemeData
	ProcessingDate  time.Time
	Reason          error
	Amount          int
}
