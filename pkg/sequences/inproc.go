package sequences

import "sync"

func InProc(min, max int) *inProc {
	return &inProc{
		min: min,
		max: max,
		cur: min,
	}
}

type inProc struct {
	min, max, cur int
	sync.Mutex
}

func (ip *inProc) Next() int {
	ip.Lock()
	defer ip.Unlock()

	r := ip.cur
	ip.cur++
	if ip.cur > ip.max {
		ip.cur = ip.min
	}

	return r
}
