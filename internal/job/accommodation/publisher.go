package accommodation

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"travel-crawler/internal/applog"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/queue/kafka"
	"travel-crawler/internal/requestcontext"
)

type Publisher struct {
	publisher kafka.Publisher
	logger    *zap.SugaredLogger
}

func NewPublisher(
	publisher kafka.Publisher,
	logger *zap.SugaredLogger,
) *Publisher {
	return &Publisher{
		publisher: publisher,
		logger:    logger,
	}
}

// Publish creates and publishes an accommodation crawl task.
func (p *Publisher) Publish(ctx context.Context) {
	traceID := helper.NewUUID()
	requestID := helper.NewUUID()

	publishCtx := requestcontext.WithRequestContext(
		ctx,
		traceID,
		requestID,
	)

	logger := applog.FromContext(p.logger, publishCtx)

	taskID := helper.NewUUID()

	task := CrawlTask{
		ID:           taskID,
		Type:         CrawlTaskType,
		Version:      CrawlTaskVersion,
		ProviderCode: constant.ProviderCodeNuitee,
		Country:      "VN",
		RequestedAt:  helper.NowUTC(),
		TraceID:      traceID,
	}

	payload, err := json.Marshal(task)
	if err != nil {
		logger.Errorw(
			"Fail to marshal accommodation crawl task",
			"task", task,
			"error", err,
		)
		return
	}

	if err := p.publisher.Publish(publishCtx, taskID, payload); err != nil {
		logger.Errorw(
			"Fail to publish accommodation crawl task",
			"task", task,
			"error", err,
		)
		return
	}

	logger.Infow(
		"Publish accommodation crawl task",
		"task", task,
	)
}
