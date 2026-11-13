package accommodation

import (
	"context"
	"fmt"
	"time"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/integration/nuitee"
	"travel-crawler/internal/provider"
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
	req accommodation.SearchRequest,
) (*accommodation.SearchResponse, error) {
	if err := p.guard.Wait(ctx); err != nil {
		return nil, err
	}

	offset := helper.CalculateOffset(req.Page, req.PageSize)

	queryParams := map[string]string{
		"countryCode": req.Country,
		"offset":      helper.ParseString(offset),
		"limit":       helper.ParseString(req.PageSize),
	}

	if req.UpdatedAt != nil {
		queryParams["lastUpdatedAt"] = req.UpdatedAt.Format(time.RFC3339)
	}

	resp, err := p.client.GetHotels(ctx, queryParams)
	if err != nil {
		return nil, fmt.Errorf(
			"get Nuitee hotels: %w",
			p.guard.HandleError(ctx, err),
		)
	}

	return &accommodation.SearchResponse{
		Accommodations: mapHotels(resp.Data),
		Pagination:     provider.Pagination{Total: &resp.Total},
	}, nil
}

var _ accommodation.Provider = (*Provider)(nil)
