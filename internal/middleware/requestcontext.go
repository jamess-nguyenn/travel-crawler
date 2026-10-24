package middleware

import (
	"net/http"

	"github.com/gorilla/mux"

	"travel-crawler/internal/helper"
	"travel-crawler/internal/requestcontext"
)

func RequestContext() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			traceID := getTraceID(r)
			requestID := helper.NewUUID()

			ctx := requestcontext.WithRequestContext(
				r.Context(),
				traceID,
				requestID,
			)

			w.Header().Set(
				requestcontext.TraceIDHeader,
				traceID,
			)

			w.Header().Set(
				requestcontext.RequestIDHeader,
				requestID,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}

func getTraceID(r *http.Request) string {
	return getHeader(
		r,
		requestcontext.TraceIDHeader,
	)
}

func getHeader(
	r *http.Request,
	header string,
) string {
	value := r.Header.Get(header)

	return helper.ResolveTraceID(value)
}
