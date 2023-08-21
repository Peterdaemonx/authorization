package timing

import "context"

type ctxKey int

const (
	_ ctxKey = iota
	timerCtxKey
)

// contextWithTimer amends a context.Context with the data needed to start tracking stuff
func contextWithTimer(ctx context.Context) context.Context {
	t, _ := ctx.Value(timerCtxKey).(*timer)
	if t == nil {
		ctx = context.WithValue(ctx, timerCtxKey, newTimer())
	}
	return ctx
}

func timerFromContext(ctx context.Context) *timer {
	t, _ := ctx.Value(timerCtxKey).(*timer)
	if t == nil {
		return newTimer()
	}
	return t
}

func StartRoot(ctx context.Context) {
	timerFromContext(ctx).startRoot()
}

func StopRoot(ctx context.Context) {
	timerFromContext(ctx).stopRoot()
}

func Start(ctx context.Context, label string) {
	timerFromContext(ctx).start(label)
}

func Stop(ctx context.Context, label string) {
	timerFromContext(ctx).stop(label)
}

func Tag(ctx context.Context, label, value string) {
	timerFromContext(ctx).tag(label, value)
}
