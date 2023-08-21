package timing

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

func newRing(interval, duration time.Duration) *ring {
	numRings := int(duration / interval)
	if numRings == 0 {
		return &ring{}
	}

	first := &frame{nr: 0, timings: make([]timer, 0)}
	cur := first

	for i := 1; i < numRings; i++ {
		next := &frame{nr: i}
		cur.next = next
		cur = next
	}
	cur.next = first

	return &ring{
		current:  first,
		interval: interval,
		duration: duration,
	}
}

type frame struct {
	next    *frame
	timings []timer
	nr      int
}

type ring struct {
	current  *frame
	mu       sync.Mutex
	interval time.Duration
	duration time.Duration
}

func (r *ring) rotate(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.mu.Lock()
			last := r.current
			r.current = last.next
			r.current.timings = make([]timer, 0)
			r.mu.Unlock()
			ticker.Reset(r.interval)
		}
	}

}

func (r *ring) capture(t timer) {
	r.mu.Lock()
	r.current.timings = append(r.current.timings, t)
	r.mu.Unlock()
}

func (r *ring) ShortLabel() string {
	return fmt.Sprintf("last_%s", r.duration.String())
}

func (r *ring) Label() string {
	return fmt.Sprintf("last_%s_per_%s", r.duration.String(), r.interval.String())
}

func (r *ring) Summarize(excludeSegments ...string) (*Timing, map[string]*Timing) {
	r.mu.Lock()
	defer r.mu.Unlock()

	first := r.current
	cur := first

	totalRes := NewTiming()
	tagsRes := map[string]*Timing{}

	for {
		for _, t := range cur.timings {
			netTime := t.NetTime()
			brutTime := t.BrutTime(excludeSegments...)

			tagKeys := []string{}
			for k := range t.tags {
				tagKeys = append(tagKeys, k)
			}
			sort.Strings(tagKeys)
			tag := ""
			for _, k := range tagKeys {
				tag += k + "=" + t.tags[k] + "|"
			}

			tagRes, ok := tagsRes[tag]
			if !ok {
				tagRes = NewTiming()
			}

			totalRes.Add(netTime, brutTime)
			tagRes.Add(netTime, brutTime)

			for sk, sv := range t.finishedTimers {
				totalRes.AddSegment(sk, sv)
				tagRes.AddSegment(sk, sv)
			}

			tagsRes[tag] = tagRes
		}

		// After a full rotation, stop
		if cur.next == first {
			return totalRes, tagsRes
		}

		cur = cur.next
	}
}

func NewTiming() *Timing {
	return &Timing{
		Segments: map[string]Segment{},
	}
}

type Timing struct {
	Amount    int
	NetMax    time.Duration
	NetTotal  time.Duration
	BrutMax   time.Duration
	BrutTotal time.Duration
	Segments  map[string]Segment
}

func (t *Timing) Add(net, brut time.Duration) {
	t.Amount++

	t.NetTotal += net
	if t.NetMax < net {
		t.NetMax = net
	}

	t.BrutTotal += brut
	if t.BrutMax < brut {
		t.BrutMax = brut
	}
}

func (t *Timing) AddSegment(l string, d time.Duration) {
	seg, ok := t.Segments[l]
	if !ok {
		seg = Segment{}
	}
	seg.Add(d)
	t.Segments[l] = seg
}

func (t *Timing) NetAvg() time.Duration {
	if t.NetTotal == 0 {
		return 0
	}

	return time.Duration(int64(t.NetTotal) / int64(t.Amount))
}

func (t *Timing) BrutAvg() time.Duration {
	if t.BrutTotal == 0 {
		return 0
	}

	return time.Duration(int64(t.BrutTotal) / int64(t.Amount))
}

type Segment struct {
	Amount int
	Max    time.Duration
	Total  time.Duration
}

func (s *Segment) Add(d time.Duration) {
	s.Amount++

	s.Total += d
	if s.Max < d {
		s.Max = d
	}
}

func (s *Segment) Avg() time.Duration {
	if s.Total == 0 {
		return 0
	}

	return time.Duration(int64(s.Total) / int64(s.Amount))
}
