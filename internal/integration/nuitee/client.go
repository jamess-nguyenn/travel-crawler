package nuitee

import (
	"context"
	"go.uber.org/zap"
	"travel-crawler/internal/redis"

	"travel-crawler/internal/config"
	"travel-crawler/internal/httpclient"
)

type Client struct {
	http   *httpclient.Client
	config config.Nuitee
	cache  *redis.Cache
	logger *zap.SugaredLogger
}

func NewClient(
	http *httpclient.Client,
	cfg config.Nuitee,
	cache *redis.Cache,
	logger *zap.SugaredLogger,
) *Client {
	return &Client{
		http:   http,
		config: cfg,
		cache:  cache,
		logger: logger,
	}
}

func (c *Client) buildURL(path string) string {
	return c.config.BaseURL + path
}

func (c *Client) buildHeaders() map[string]string {
	headers := httpclient.JSONHeaders()
	headers[httpclient.HeaderAPIKey] = c.config.APIKey

	return headers
}

func (c *Client) Get(ctx context.Context, path string, queryParams map[string]string, result any) error {
	req := &httpclient.HTTPRequest{
		Headers:     c.buildHeaders(),
		QueryParams: queryParams,
	}

	return c.http.Get(ctx, c.buildURL(path), req, result)
}

func (c *Client) Post(ctx context.Context, path string, body any, result any) error {
	req := &httpclient.HTTPRequest{
		Headers: c.buildHeaders(),
		Body:    body,
	}

	return c.http.Post(ctx, c.buildURL(path), req, result)
}
