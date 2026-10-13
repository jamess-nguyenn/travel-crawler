package httpclient

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type HTTPRequest struct {
	QueryParams map[string]string
	FormData    map[string]string
	PathParams  map[string]string
	Headers     map[string]string
	Body        any
}

func (c *Client) newRequest(
	ctx context.Context,
	result any,
) *resty.Request {
	return c.client.
		R().
		SetContext(ctx).
		SetResult(result)
}

func (c *Client) applyRequestData(
	r *resty.Request,
	req *HTTPRequest,
) {
	if req == nil {
		return
	}

	if req.QueryParams != nil {
		r.SetQueryParams(req.QueryParams)
	}

	if req.FormData != nil {
		r.SetFormData(req.FormData)
	}

	if req.PathParams != nil {
		r.SetPathParams(req.PathParams)
	}

	if req.Headers != nil {
		r.SetHeaders(req.Headers)
	}

	if req.Body != nil {
		r.SetBody(req.Body)
	}
}
