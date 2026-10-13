package applog

import (
	"context"
	"go.uber.org/zap"
	"travel-crawler/internal/requestcontext"
)

func FromContext(logger *zap.SugaredLogger, ctx context.Context) *zap.SugaredLogger {
	return logger.With(
		"trace_id", requestcontext.TraceID(ctx),
		"request_id", requestcontext.RequestID(ctx),
	)
}
