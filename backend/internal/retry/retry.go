// Package retry retries startup operations with exponential backoff.
package retry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// Policy describes how an operation is retried.
type Policy struct {
	Initial        time.Duration // first delay
	Max            time.Duration // delay cap
	MaxElapsed     time.Duration // give up once this much time has passed
	AttemptTimeout time.Duration // per-attempt context timeout (0: none)
}

// Startup is the policy used for connecting to dependencies at startup.
var Startup = Policy{
	Initial:        500 * time.Millisecond,
	Max:            10 * time.Second,
	MaxElapsed:     90 * time.Second,
	AttemptTimeout: 15 * time.Second,
}

type permanentError struct{ err error }

func (p permanentError) Error() string { return p.err.Error() }
func (p permanentError) Unwrap() error { return p.err }

// Permanent marks err as not worth retrying.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return permanentError{err}
}

// Do calls fn until it succeeds, returns a Permanent error, ctx is done or
// MaxElapsed has passed. Failed attempts are logged with what as context.
func (p Policy) Do(ctx context.Context, logger *slog.Logger, what string, fn func(context.Context) error) error {
	start := time.Now()
	delay := p.Initial
	for attempt := 1; ; attempt++ {
		actx, cancel := ctx, context.CancelFunc(func() {})
		if p.AttemptTimeout > 0 {
			actx, cancel = context.WithTimeout(ctx, p.AttemptTimeout)
		}
		err := fn(actx)
		cancel()
		if err == nil {
			if attempt > 1 {
				logger.Info(what+": succeeded", "attempt", attempt)
			}
			return nil
		}
		if pe := (permanentError{}); errors.As(err, &pe) {
			return fmt.Errorf("%s: %w", what, pe.err)
		}
		if time.Since(start)+delay > p.MaxElapsed {
			return fmt.Errorf("%s: giving up after %d attempts: %w", what, attempt, err)
		}
		logger.Warn(what+": failed, retrying", "attempt", attempt, "retry_in", delay, "error", err)
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: %w", what, ctx.Err())
		case <-time.After(delay):
		}
		delay = min(delay*2, p.Max)
	}
}
