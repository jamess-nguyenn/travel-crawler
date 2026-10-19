package nuitee

import (
	"context"
	"go.uber.org/zap"
	"travel-crawler/internal/config"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/redis"
)

const (
	URIDataHotels = "/data/hotels"
)

type Client interface {
	GetHotels(ctx context.Context, queryParams map[string]string) (*HotelResponse, error)
}

type client struct {
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
) Client {
	return &client{
		http:   http,
		config: cfg,
		cache:  cache,
		logger: logger,
	}
}

func (c *client) buildHeaders() map[string]string {
	headers := httpclient.JSONHeaders()
	headers[httpclient.HeaderAPIKey] = c.config.APIKey

	return headers
}

func (c *client) GetHotels(ctx context.Context, queryParams map[string]string) (*HotelResponse, error) {
	var result HotelResponse

	req := &httpclient.HTTPRequest{
		Headers:     c.buildHeaders(),
		QueryParams: queryParams,
	}

	url := c.config.BaseURL + URIDataHotels

	if err := c.http.Get(ctx, url, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
