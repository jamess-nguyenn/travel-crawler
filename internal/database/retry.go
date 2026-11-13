package database

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	errLockDeadlock    = 1213
	errLockWaitTimeout = 1205

	maxDeadlockRetry  = 4
	deadlockBaseDelay = 100 * time.Millisecond
	deadlockMaxDelay  = 2 * time.Second
)

// isRetryableError reports whether a database error
// is a transient MySQL deadlock or lock-wait timeout.
func isRetryableError(err error) bool {
	var mysqlErr *mysql.MySQLError

	if !errors.As(err, &mysqlErr) {
		return false
	}

	switch mysqlErr.Number {
	case errLockDeadlock, errLockWaitTimeout:
		return true
	default:
		return false
	}
}

func retryDelay(attempt int) time.Duration {
	delay := deadlockBaseDelay

	for i := 0; i < attempt; i++ {
		if delay >= deadlockMaxDelay/2 {
			delay = deadlockMaxDelay
			break
		}

		delay *= 2
	}

	// Add jitter to avoid synchronized retries.
	jitterLimit := delay / 2
	if jitterLimit > 0 {
		delay += time.Duration(rand.Int64N(int64(jitterLimit)))
	}

	if delay > deadlockMaxDelay {
		return deadlockMaxDelay
	}

	return delay
}

func retryOperation(ctx context.Context, operation func() error) error {
	for attempt := 0; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := operation()
		if err == nil {
			return nil
		}

		if !isRetryableError(err) || attempt >= maxDeadlockRetry {
			return err
		}

		delay := retryDelay(attempt)

		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return ctx.Err()

		case <-timer.C:
		}
	}
}
