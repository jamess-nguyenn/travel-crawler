package applog

import "go.uber.org/zap"

func WithCreate(
	logger *zap.SugaredLogger,
	data any,
) *zap.SugaredLogger {
	return logger.With(
		"data", data,
	)
}
