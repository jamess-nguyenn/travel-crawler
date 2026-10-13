package database

import "travel-crawler/internal/database/scope"

type ListQuery struct {
	Conditions map[string]any
	Filters    []scope.Scope
	Scopes     []scope.Scope
	Columns    []string
}

func (q ListQuery) NormalizeColumns() []string {
	if len(q.Columns) > 0 {
		return q.Columns
	}

	return []string{"*"}
}
