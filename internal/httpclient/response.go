package httpclient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func (c *Client) handleResponse(
	resp *resty.Response,
) error {
	if resp.IsSuccess() {
		return nil
	}

	return fmt.Errorf(
		"HTTP request failed: method=%s url=%s status=%s",
		resp.Request.Method,
		resp.Request.URL,
		resp.Status(),
	)
}
