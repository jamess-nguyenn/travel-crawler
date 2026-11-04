package httpclient

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"travel-crawler/internal/applog"
)

func RegisterLogger(
	client *resty.Client,
	logger *zap.SugaredLogger,
) {
	client.OnAfterResponse(func(
		_ *resty.Client,
		resp *resty.Response,
	) error {
		if resp == nil || resp.Request == nil {
			return nil
		}

		r := resp.Request
		size := resp.Size()

		requestLogger := applog.FromContext(logger, r.Context()).With(
			"method", r.Method,
			"url", r.URL,
			"query", r.QueryParam,
			"path", r.PathParams,
			"form", r.FormData,
			"body", r.Body,
			"header", safeHeader(r.Header),
			"status_code", resp.StatusCode(),
			"response_size_bytes", size,
			"duration", resp.Time().String(),
		)

		// Include the parsed response only when it is small enough.
		if shouldLogResponseBody(size) {
			requestLogger = requestLogger.With(
				"result", resp.Result(),
			)
		}

		if resp.IsError() {
			requestLogger.Errorw(
				"HTTP request completed with error status",
				"error", resp.Error(),
			)

			return nil
		}

		requestLogger.Info("HTTP request completed")

		return nil
	})
}

func safeHeader(headers http.Header) map[string][]string {
	result := make(map[string][]string, len(headers))

	for key, values := range headers {
		normalized := strings.ToLower(key)

		switch normalized {
		case "authorization",
			"x-api-key",
			"api-key",
			"cookie",
			"set-cookie",
			"x-auth-token":

			result[key] = []string{"[REDACTED]"}

		default:
			result[key] = values
		}
	}

	return result
}
