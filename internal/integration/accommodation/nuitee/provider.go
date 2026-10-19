package nuitee

import (
	"context"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/provider"
)

type Provider struct {
	client Client
}

func New(client Client) *Provider {
	return &Provider{
		client: client,
	}
}

func (p *Provider) Code() string {
	return constant.ProviderCodeNuitee
}

func (p *Provider) Search(
	ctx context.Context,
	req accommodation.SearchRequest,
) (*accommodation.SearchResponse, error) {
	offset := helper.CalculateOffset(req.Page, req.PageSize)

	queryParams := map[string]string{
		"countryCode": req.Country,
		"offset":      helper.ParseString(offset),
		"limit":       helper.ParseString(req.PageSize),
	}

	resp, err := p.client.GetHotels(ctx, queryParams)
	if err != nil {
		return nil, err
	}

	return &accommodation.SearchResponse{
		Accommodations: mapHotels(resp.Data),
		Pagination:     provider.Pagination{Total: &resp.Total},
	}, nil
}

var _ accommodation.Provider = (*Provider)(nil)
