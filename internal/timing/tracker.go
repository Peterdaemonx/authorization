package timing

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Slice struct {
	Interval time.Duration
	Duration time.Duration
}

func (s Slice) Label() string {
	return fmt.Sprintf("last_%s_per_%s", s.Duration.String(), s.Interval.String())
}

func NewTracker(slices ...Slice) *Tracker {
	t := Tracker{
		rings:  make([]*ring, 0),
		timers: make(chan timer),
	}

	for _, slice := range slices {
		t.rings = append(t.rings, newRing(slice.Interval, slice.Duration))
	}

	return &t
}

// Tracker manages the timing of things, collecting all data and exposing the results
type Tracker struct {
	rings  []*ring
	timers chan timer
}

func (t *Tracker) Process(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case timed := <-t.timers:
				for _, r := range t.rings {
					r.capture(timed)
				}
			}
		}
	}()

	for _, r := range t.rings {
		go r.rotate(ctx)
	}
}

func (t *Tracker) Http(label string, handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(contextWithTimer(request.Context()))
		StartRoot(request.Context())
		handler(writer, request)
		StopRoot(request.Context())
		timed := timerFromContext(request.Context())
		t.timers <- *timed
	}
}
