package core_test

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
)

func TestMiddlewareChain(t *testing.T) {
	calls := []string{}
	base := func(_ context.Context, _ core.Message) error {
		calls = append(calls, "base")
		return nil
	}
	mw1 := func(next core.Handler) core.Handler {
		return func(ctx context.Context, msg core.Message) error {
			calls = append(calls, "mw1")
			return next(ctx, msg)
		}
	}
	mw2 := func(next core.Handler) core.Handler {
		return func(ctx context.Context, msg core.Message) error {
			calls = append(calls, "mw2")
			return next(ctx, msg)
		}
	}
	chain := mw1(mw2(base))
	_ = chain(context.Background(), nil)
	if len(calls) != 3 || calls[0] != "mw1" || calls[1] != "mw2" || calls[2] != "base" {
		t.Errorf("middleware chain order wrong: %v", calls)
	}
}
