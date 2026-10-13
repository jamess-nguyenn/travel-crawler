package requestcontext

import (
	"context"
	"travel-crawler/internal/helper"
)

type contextKey string

const (
	traceIDKey   contextKey = "trace_id"
	requestIDKey contextKey = "request_id"
)

const (
	TraceIDHeader   = "X-Trace-ID"
	RequestIDHeader = "X-Request-ID"
)

func WithTraceID(
	ctx context.Context,
	traceID string,
) context.Context {
	return context.WithValue(
		ctx,
		traceIDKey,
		traceID,
	)
}

func WithRequestID(
	ctx context.Context,
	requestID string,
) context.Context {
	return context.WithValue(
		ctx,
		requestIDKey,
		requestID,
	)
}

func WithRequestContext(
	ctx context.Context,
	traceID string,
	requestID string,
) context.Context {
	ctx = WithTraceID(ctx, traceID)
	return WithRequestID(ctx, requestID)
}

func WithNewRequestContext(
	ctx context.Context,
) context.Context {
	return WithRequestContext(
		ctx,
		helper.NewUUID(),
		helper.NewUUID(),
	)
}

func TraceID(ctx context.Context) string {
	traceID, _ := ctx.Value(traceIDKey).(string)
	return traceID
}

func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey).(string)
	return requestID
}
