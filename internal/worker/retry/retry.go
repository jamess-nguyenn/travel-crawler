package retry

import (
	"context"
	"errors"
	"math/rand/v2"
	"net"
	"net/http"
	"time"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/worker"
)

const (
	MaxAttempt        = 3
	DefaultRetryDelay = 2 * time.Second
	MaxRetryDelay     = 30 * time.Second
)

type RetryAfterError interface {
	RetryAfter() time.Duration
}

// IsRetryable reports whether an error should be retried.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var invalidTaskErr *worker.InvalidTaskError
	if errors.As(err, &invalidTaskErr) {
		return false
	}

	if errors.Is(err, context.Canceled) {
		return false
	}

	var httpErr *httpclient.HTTPError
	if errors.As(err, &httpErr) {
		switch {
		case httpErr.StatusCode == http.StatusRequestTimeout:
			return true
		case httpErr.StatusCode == http.StatusTooManyRequests:
			return true
		case httpErr.StatusCode >= 500 && httpErr.StatusCode <= 599:
			return true
		default:
			return false
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	return false
}

// Delay returns the delay before the next attempt.
// failedAttempt is 1 for the delay after the first failure.
func Delay(err error, failedAttempt int) time.Duration {
	var retryErr RetryAfterError
	if errors.As(err, &retryErr) {
		if delay := retryErr.RetryAfter(); delay > 0 {
			if delay > MaxRetryDelay {
				return MaxRetryDelay
			}

			return delay
		}
	}

	if failedAttempt < 1 {
		failedAttempt = 1
	}

	delay := DefaultRetryDelay

	for i := 1; i < failedAttempt; i++ {
		if delay >= MaxRetryDelay/2 {
			delay = MaxRetryDelay
			break
		}

		delay *= 2
	}

	if delay > MaxRetryDelay {
		delay = MaxRetryDelay
	}

	// Full jitter: random delay between zero and the cap.
	return time.Duration(rand.Int64N(int64(delay) + 1))
}

// Wait waits for the specified delay or until the context is canceled.
func Wait(ctx context.Context, delay time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
