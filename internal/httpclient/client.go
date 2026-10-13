package httpclient

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Client struct {
	client *resty.Client
	logger *zap.SugaredLogger
}

func New(
	logger *zap.SugaredLogger,
) *Client {
	c := resty.New()

	client := &Client{
		client: c,
		logger: logger,
	}

	RegisterLogger(c, logger)

	return client
}

func (c *Client) execute(ctx context.Context, method string, url string, req *HTTPRequest, result any) error {
	r := c.newRequest(ctx, result)

	c.applyRequestData(r, req)

	resp, err := r.Execute(method, url)
	if err != nil {
		return err
	}

	return c.handleResponse(resp)
}

func (c *Client) Get(ctx context.Context, url string, req *HTTPRequest, result any) error {
	return c.execute(ctx, http.MethodGet, url, req, result)
}

func (c *Client) Post(ctx context.Context, url string, req *HTTPRequest, result any) error {
	return c.execute(ctx, http.MethodPost, url, req, result)
}

func (c *Client) Put(ctx context.Context, url string, req *HTTPRequest, result any) error {
	return c.execute(ctx, http.MethodPut, url, req, result)
}
