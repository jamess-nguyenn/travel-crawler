package airline

import (
	"context"
)

type SearchRequest struct {
	Page     int
	PageSize int
}

type SearchResponse struct {
	Airlines []Airline
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
