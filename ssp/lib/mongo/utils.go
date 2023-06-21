package mongo

import (
	"context"
	"time"
)

func SetContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, deadlinePresent := ctx.Deadline(); deadlinePresent {
		return context.WithCancel(ctx)
	}
	// context.WithTimeout will return a derived context with timeout <= parent context timeout
	return context.WithTimeout(ctx, timeout)
}
