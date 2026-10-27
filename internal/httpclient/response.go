package httpclient

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
	"travel-crawler/internal/helper"
)

const (
	maxResponseLogSize int64 = 10 * 1024
)

func (c *Client) handleResponse(resp *resty.Response) error {
	if resp == nil || resp.IsSuccess() {
		return nil
	}

	var method, url, body string

	if r := resp.Request; r != nil {
		method = r.Method
		url = r.URL
	}

	if shouldLogResponseBody(resp.Size()) {
		body = resp.String()
	}

	return &HTTPError{
		Method:          method,
		URL:             url,
		StatusCode:      resp.StatusCode(),
		ResponseBody:    body,
		RetryAfterDelay: parseRetryAfter(resp.Header().Get(HeaderRetryAfter), helper.NowUTC()),
	}
}

func shouldLogResponseBody(size int64) bool {
	return size >= 0 && size <= maxResponseLogSize
}

func parseRetryAfter(value string, now time.Time) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}

	// Retry-After can be a delay in seconds.
	if second, err := helper.ParseInt64(value); err == nil {
		if second <= 0 {
			return 0
		}

		maxSecond := (1<<63 - 1) / int64(time.Second)
		if second > maxSecond {
			return time.Duration(1<<63 - 1)
		}

		return time.Duration(second) * time.Second
	}

	// Or it can be an HTTP date.
	retryAt, err := http.ParseTime(value)
	if err != nil {
		return 0
	}

	delay := retryAt.Sub(now)
	if delay <= 0 {
		return 0
	}

	return delay
}
