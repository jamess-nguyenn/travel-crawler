package applog

import (
	"go.uber.org/zap"
)

func WithGet(
	logger *zap.SugaredLogger,
	conditions map[string]any,
	columns []string,
	data any,
) *zap.SugaredLogger {
	return logger.With(
		"conditions", conditions,
		"columns", columns,
		"data", data,
	)
}
