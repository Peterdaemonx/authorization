package mock

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/sequences"
)

// SequenceStore is a IMPERFECT implementation of sequences.SequenceStore
// It does not work well in some non-typical situations, like when the rollover point isn't a multiple of the
// minValue
// So... don't use it like that (you shouldn't use it outside of development anyway).
type SequenceStore struct {
	cur int
	rv  string
}

func (s *SequenceStore) NextBatch(ctx context.Context, batchSize int, rolloverCond sequences.RolloverCondition, rolloverVal sequences.Block) (int, error) {
	if rolloverCond(sequences.Block{NextValue: s.cur, RolloverValue: s.rv}) {
		s.cur = rolloverVal.NextValue
		s.rv = rolloverVal.RolloverValue
	}

	n := s.cur
	s.cur += batchSize

	return n, nil
}
