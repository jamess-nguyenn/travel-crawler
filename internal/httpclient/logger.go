package httpclient

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
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
		r := resp.Request

		requestLogger := applog.FromContext(logger, r.Context()).With(
			"method", r.Method,
			"url", r.URL,
			"query", r.QueryParam,
			"path", r.PathParams,
			"form", r.FormData,
			"body", r.Body,
			"status_code", resp.StatusCode(),
			"duration", resp.Time().String(),
		)

		if resp.IsError() {
			requestLogger.Errorw(
				"HTTP request completed with error status",
				"error", resp.Error(),
			)

			return nil
		}

		requestLogger.Infow(
			"HTTP request completed",
			"result", resp.Result(),
		)

		return nil
	})
}
