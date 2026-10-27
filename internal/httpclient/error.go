package httpclient

import (
	"fmt"
	"time"
)

type HTTPError struct {
	Method          string
	URL             string
	StatusCode      int
	ResponseBody    string
	RetryAfterDelay time.Duration
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf(
		"HTTP request failed: method=%s url=%s status_code=%d",
		e.Method,
		e.URL,
		e.StatusCode,
	)
}

func (e *HTTPError) RetryAfter() time.Duration {
	return e.RetryAfterDelay
}
