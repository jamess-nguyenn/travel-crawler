package nuitee

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/provider"
)

const defaultNuiteeCooldown = 2 * time.Second

type Provider struct {
	client  Client
	limiter provider.RateLimiter
}

func New(
	client Client,
	limiter provider.RateLimiter,
) *Provider {
	return &Provider{
		client:  client,
		limiter: limiter,
	}
}

func (p *Provider) Code() string {
	return constant.ProviderCodeNuitee
}

func (p *Provider) Search(
	ctx context.Context,
	req accommodation.SearchRequest,
) (*accommodation.SearchResponse, error) {
	// Acquire a shared request slot before calling Nuitee.
	if p.limiter != nil {
		if err := p.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf(
				"wait for Nuitee rate limiter: %w",
				err,
			)
		}
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
		return nil, p.handleRequestError(ctx, err)
	}

	return &accommodation.SearchResponse{
		Accommodations: mapHotels(resp.Data),
		Pagination:     provider.Pagination{Total: &resp.Total},
	}, nil
}

func (p *Provider) handleRequestError(
	ctx context.Context,
	err error,
) error {
	var httpErr *httpclient.HTTPError

	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTooManyRequests {
		return err
	}

	// Use the provider's Retry-After value when available.
	// Fall back to a short shared cooldown if the header is absent.
	cooldown := httpErr.RetryAfter()
	if cooldown <= 0 {
		cooldown = defaultNuiteeCooldown
	}

	if cooldownErr := p.limiter.Cooldown(ctx, cooldown); cooldownErr != nil {
		// Preserve the original HTTP error so the retry handler can
		// still identify this as an HTTP 429.
		return errors.Join(
			err,
			fmt.Errorf("activate Nuitee shared cooldown: %w", cooldownErr),
		)
	}

	return err
}

var _ accommodation.Provider = (*Provider)(nil)
