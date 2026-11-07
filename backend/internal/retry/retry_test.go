package retry

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
)

var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

var fast = Policy{Initial: time.Millisecond, Max: 2 * time.Millisecond, MaxElapsed: time.Second, AttemptTimeout: time.Second}

func TestSucceedsAfterFailures(t *testing.T) {
	calls := 0
	err := fast.Do(context.Background(), discard, "op", func(ctx context.Context) error {
		calls++
		if _, ok := ctx.Deadline(); !ok {
			t.Error("attempt context has no deadline")
		}
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil || calls != 3 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}

func TestGivesUp(t *testing.T) {
	p := fast
	p.MaxElapsed = 20 * time.Millisecond
	sentinel := errors.New("down")
	calls := 0
	err := p.Do(context.Background(), discard, "op", func(context.Context) error { calls++; return sentinel })
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want wrapped sentinel", err)
	}
	if calls < 2 {
		t.Errorf("calls = %d, expected several attempts", calls)
	}
}

func TestPermanentStopsImmediately(t *testing.T) {
	sentinel := errors.New("config error")
	calls := 0
	err := fast.Do(context.Background(), discard, "op", func(context.Context) error { calls++; return Permanent(sentinel) })
	if !errors.Is(err, sentinel) || calls != 1 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
	if Permanent(nil) != nil {
		t.Error("Permanent(nil) != nil")
	}
}

func TestContextCancel(t *testing.T) {
	p := Policy{Initial: time.Hour, Max: time.Hour, MaxElapsed: 10 * time.Hour}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(10 * time.Millisecond); cancel() }()
	err := p.Do(ctx, discard, "op", func(context.Context) error { return errors.New("down") })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}
