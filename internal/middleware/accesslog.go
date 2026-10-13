package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/helper"
)

const maxResponseLogBytes = 64 * 1024

type responseWriter struct {
	http.ResponseWriter
	status      int
	written     int
	wroteHeader bool
	body        bytes.Buffer
	truncated   bool
}

func (w *responseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}

	w.status = status
	w.wroteHeader = true

	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(body []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	n, err := w.ResponseWriter.Write(body)
	w.written += n

	if n > 0 {
		remaining := maxResponseLogBytes - w.body.Len()

		if remaining > 0 {
			capture := n
			if capture > remaining {
				capture = remaining
			}

			_, _ = w.body.Write(body[:capture])
		}

		if n > remaining {
			w.truncated = true
		}
	}

	return n, err
}

func AccessLog(
	logger *zap.SugaredLogger,
) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			start := helper.NowUTC()

			writer := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(writer, r)

			fields := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"status", writer.status,
				"bytes", writer.written,
				"duration", time.Since(start).String(),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			}

			if !writer.truncated && writer.body.Len() > 0 {
				var responseBody any

				if err := json.Unmarshal(
					writer.body.Bytes(),
					&responseBody,
				); err == nil {
					fields = append(
						fields,
						"response", responseBody,
					)
				}
			}

			applog.FromContext(logger, r.Context()).Infow("HTTP request", fields...)
		})
	}
}
