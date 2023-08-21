//go:build integration
// +build integration

package spanner

import (
	"context"
	"testing"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/sequences"
)

func TestSequenceRepository_NextBatch(t *testing.T) {
	sr := connectToSequenceRepository()

	noRollover := func(block sequences.Block) bool {
		return false
	}
	doRollover := func(block sequences.Block) bool {
		return true
	}
	rb := sequences.Block{NextValue: 0}

	val, err := sr.NextBatch(context.Background(), 5, noRollover, rb)
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}
	if val != 1000000 {
		t.Errorf("Expected 0 on first NextBatch() call, got %d", val)
	}

	val, err = sr.NextBatch(context.Background(), 5, noRollover, rb)
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}
	if val != 1000005 {
		t.Errorf("Expected 5 on second NextBatch() call, got %d", val)
	}

	val, err = sr.NextBatch(context.Background(), 5, doRollover, rb)
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}
	if val != 0 {
		t.Errorf("Expected 0 on third NextBatch() call, got %d", val)
	}

	val, err = sr.NextBatch(context.Background(), 5, noRollover, rb)
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}
	if val != 5 {
		t.Errorf("Expected 5 on fourth NextBatch() call, got %d", val)
	}
}

func connectToSequenceRepository() *SequenceRepository {
	wtimeout, _ := time.ParseDuration("10s")
	rtimeout, _ := time.ParseDuration("10s")
	return NewSequenceRepository("mastercard_stan", Client, rtimeout, wtimeout)
}
