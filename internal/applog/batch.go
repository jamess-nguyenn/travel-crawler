package applog

import "go.uber.org/zap"

func WithBatch(
	logger *zap.SugaredLogger,
	totalRow int,
	batchProgress string,
	batchSize int,
	batchData any,
) *zap.SugaredLogger {
	return logger.With(
		"total_row", totalRow,
		"batch_progress", batchProgress,
		"batch_size", batchSize,
		"batch_data", batchData,
	)
}
