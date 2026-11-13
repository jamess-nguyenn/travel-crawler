package accommodation

import (
	"context"
	"travel-crawler/internal/integration/nuitee"
)

const (
	URIDataHotels = "/data/hotels"
)

type Client interface {
	GetHotels(ctx context.Context, queryParams map[string]string) (*HotelResponse, error)
}

type client struct {
	api *nuitee.Client
}

func NewClient(
	api *nuitee.Client,
) Client {
	return &client{
		api: api,
	}
}

func (c *client) GetHotels(ctx context.Context, queryParams map[string]string) (*HotelResponse, error) {
	var result HotelResponse

	if err := c.api.Get(ctx, URIDataHotels, queryParams, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
