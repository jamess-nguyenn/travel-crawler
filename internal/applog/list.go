package applog

import (
	"go.uber.org/zap"
	"travel-crawler/internal/database"
)

func WithList(
	logger *zap.SugaredLogger,
	listQuery database.ListQuery,
	total int64,
	data any,
) *zap.SugaredLogger {
	return logger.With(
		"list_query", listQuery,
		"total", total,
		"data", data,
	)
}
