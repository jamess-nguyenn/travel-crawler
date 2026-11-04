package accommodation

import (
	"context"
	"time"

	"travel-crawler/internal/provider"
)

type SearchRequest struct {
	Country   string
	UpdatedAt *time.Time
	Page      int
	PageSize  int
}

type SearchResponse struct {
	Accommodations []Accommodation
	Pagination     provider.Pagination
}

type Provider interface {
	Code() string

	Search(
		ctx context.Context,
		req SearchRequest,
	) (*SearchResponse, error)
}

type ProviderManager interface {
	Get(code string) (Provider, error)
}
