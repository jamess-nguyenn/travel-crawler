package database

import "travel-crawler/internal/database/scope"

type ListQuery struct {
	Conditions map[string]any
	Filters    []scope.Scope
	Scopes     []scope.Scope
	Columns    []string
}
