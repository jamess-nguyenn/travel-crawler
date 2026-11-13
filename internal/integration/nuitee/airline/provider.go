package airline

import (
	"context"
	"fmt"
	"travel-crawler/internal/airline"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/integration/nuitee"
)

type Provider struct {
	client Client
	guard  *nuitee.RequestGuard
}

func New(
	client Client,
	guard *nuitee.RequestGuard,
) *Provider {
	return &Provider{
		client: client,
		guard:  guard,
	}
}

func (p *Provider) Code() string {
	return constant.ProviderCodeNuitee
}

func (p *Provider) Search(
	ctx context.Context,
	req airline.SearchRequest,
) (*airline.SearchResponse, error) {
	if err := p.guard.Wait(ctx); err != nil {
		return nil, err
	}

	queryParams := map[string]string{
		"limit": helper.ParseString(req.PageSize),
	}

	resp, err := p.client.GetAirlines(ctx, queryParams)
	if err != nil {
		return nil, fmt.Errorf(
			"get Nuitee airlines: %w",
			p.guard.HandleError(ctx, err),
		)
	}

	return &airline.SearchResponse{
		Airlines: mapAirlines(resp.Data[0].Airlines),
	}, nil
}

var _ airline.Provider = (*Provider)(nil)
