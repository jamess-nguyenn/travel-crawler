package airline

import (
	"context"
	"travel-crawler/internal/integration/nuitee"
)

const (
	URIDataAirlines = "/data/flights/airlines"
)

type Client interface {
	GetAirlines(ctx context.Context, queryParams map[string]string) (*AirlineResponse, error)
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

func (c *client) GetAirlines(ctx context.Context, queryParams map[string]string) (*AirlineResponse, error) {
	var result AirlineResponse

	if err := c.api.Get(ctx, URIDataAirlines, queryParams, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
