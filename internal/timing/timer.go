package timing

import "time"

func newTimer() *timer {
	return &timer{
		tags:           map[string]string{},
		startedTimers:  map[string]time.Time{},
		finishedTimers: map[string]time.Duration{},
	}
}

// timer is one specific timer collection of metrics for one specific
type timer struct {
	tags           map[string]string
	startedTimers  map[string]time.Time
	finishedTimers map[string]time.Duration
	rootStart      time.Time
	rootDuration   time.Duration
}

func (t *timer) tag(label, value string) {
	t.tags[label] = value
}

func (t *timer) startRoot() {
	t.rootStart = time.Now()
}

func (t *timer) stopRoot() {
	t.rootDuration = time.Since(t.rootStart)
}

func (t *timer) start(label string) {
	t.startedTimers[label] = time.Now()
}

func (t *timer) stop(label string) {
	started, ok := t.startedTimers[label]
	if !ok {
		t.finishedTimers[label] = time.Duration(0)
	}
	t.finishedTimers[label] = time.Since(started)
}

func (t *timer) NetTime() time.Duration {
	return t.rootDuration
}

func (t *timer) BrutTime(excludeSegments ...string) time.Duration {
	brut := t.rootDuration

	etMap := map[string]bool{}
	for _, et := range excludeSegments {
		etMap[et] = true
	}

	for tag, tagDur := range t.finishedTimers {
		if _, ok := etMap[tag]; ok {
			brut = brut - tagDur
		}
	}

	return brut
}
