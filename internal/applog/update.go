package applog

import "go.uber.org/zap"

func WithUpdate(
	logger *zap.SugaredLogger,
	conditions any,
	updates map[string]any,
) *zap.SugaredLogger {
	return logger.With(
		"conditions", conditions,
		"updates", updates,
	)
}
