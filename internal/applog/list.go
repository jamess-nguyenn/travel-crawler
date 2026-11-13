package applog

import (
	"go.uber.org/zap"
	"travel-crawler/internal/database"
)

type ListQueryLog struct {
	Conditions map[string]any `json:"conditions,omitempty"`
	Filters    int            `json:"filters"`
	Scopes     int            `json:"scopes"`
	Columns    []string       `json:"columns,omitempty"`
}

func WithList(
	logger *zap.SugaredLogger,
	listQuery database.ListQuery,
	total int64,
	data any,
) *zap.SugaredLogger {
	return logger.With(
		"list_query", ToListQueryLog(listQuery),
		"total", total,
		"data", data,
	)
}

func ToListQueryLog(query database.ListQuery) ListQueryLog {
	return ListQueryLog{
		Conditions: query.Conditions,
		Filters:    len(query.Filters),
		Scopes:     len(query.Scopes),
		Columns:    query.Columns,
	}
}
