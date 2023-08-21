package sequences

import (
	"context"
	"fmt"
)

// NextBatch should return the starting point of the next batch. Sequence generators will start from that number,
// increment by 1 until next+batchSize is reached and then call NextBatch again.
type Store interface {
	NextBatch(ctx context.Context, batchSize int, rolloverCond RolloverCondition, rolloverTo Block) (int, error)
}

type Block struct {
	NextValue     int
	RolloverValue string
}

type RolloverCondition func(Block) bool

type UnknownSequence struct {
	Name string
}

func (us UnknownSequence) Error() string {
	return fmt.Sprintf("unknown sequence %s", us.Name)
}
