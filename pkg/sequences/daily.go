package sequences

import (
	"context"
	"fmt"
	"math"
	"time"
)

func NewDaily(batchSize, minValue, maxValue int, store Store) daily {
	prefetch := int(math.Max(float64(batchSize/10), float64(5)))
	d := daily{
		next:      make(chan int, 1),
		reserved:  make(chan int, batchSize+prefetch),
		batchSize: batchSize,
		prefetch:  prefetch,
		minVal:    minValue,
		maxVal:    maxValue,
		store:     store,
	}

	return d
}

type daily struct {
	next      chan int
	reserved  chan int
	batchSize int
	prefetch  int
	minVal    int
	maxVal    int
	store     Store
}

func (d *daily) Next() int {
	return <-d.next
}

func (d *daily) Fill(ctx context.Context) error {

	err := d.grow(ctx)
	if err != nil {
		return fmt.Errorf("cannot init batch: %w", err)
	}
	for {
		select {
		// d.next<- means "write to d.next"
		// <-d.reserved means "read from d.reserved"
		// So case d.next <- <-d.reserved means "when we write something from d.next that we read from d.reserved"
		case d.next <- <-d.reserved:
			if len(d.reserved) < d.prefetch {
				err = d.grow(ctx)
				if err != nil {
					return fmt.Errorf("cannot grow sequence batch: %w", err)
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (d *daily) grow(ctx context.Context) error {
	rc := anyRollover(
		maxValRollover(d.maxVal),
		diffRolloverRollover(func() string {
			return time.Now().Format("0102")
		}),
	)
	startVal, err := d.store.NextBatch(ctx, d.batchSize, rc, Block{
		NextValue:     d.minVal,
		RolloverValue: time.Now().Format("0102"),
	})

	if err != nil {
		return fmt.Errorf("cannot determine batch start: %w", err)
	}

	for i := 0; i < d.batchSize; i++ {
		d.reserved <- startVal + i
	}

	return nil
}
