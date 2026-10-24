package accommodation

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/helper"
	accomjob "travel-crawler/internal/job/accommodation"
	"travel-crawler/internal/requestcontext"
)

type AccommodationProcessor struct {
	syncService accommodation.SyncService
	logger      *zap.SugaredLogger
}

func NewAccommodationProcessor(
	syncService accommodation.SyncService,
	logger *zap.SugaredLogger,
) *AccommodationProcessor {
	return &AccommodationProcessor{
		syncService: syncService,
		logger:      logger,
	}
}

func (p *AccommodationProcessor) Process(ctx context.Context, task accomjob.CrawlTask) error {
	wrkCtx := requestcontext.WithRequestContext(
		ctx,
		helper.ResolveTraceID(task.TraceID),
		helper.NewUUID(),
	)

	logger := applog.FromContext(p.logger, wrkCtx)

	logger.Infow(
		"Starting accommodation crawl",
		"task", task,
	)

	if err := p.syncService.Sync(
		wrkCtx,
		task.ProviderCode,
		accommodation.SearchRequest{
			Country: task.Country,
			City:    task.City,
		},
	); err != nil {
		logger.Errorw(
			"Fail to crawl accommodation",
			"task", task,
			"error", err,
		)

		return fmt.Errorf("sync accommodations: %w", err)
	}

	logger.Infow(
		"Crawl accommodation successfully",
		"task", task,
	)

	return nil
}

var _ Processor = (*AccommodationProcessor)(nil)
