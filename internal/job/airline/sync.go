package airline

import (
	"context"

	"go.uber.org/zap"

	"travel-crawler/internal/airline"
)

type SyncJob struct {
	service       airline.SyncService
	providerCodes []string
	logger        *zap.SugaredLogger
}

func NewSyncJob(
	service airline.SyncService,
	providerCodes []string,
	logger *zap.SugaredLogger,
) *SyncJob {
	return &SyncJob{
		service:       service,
		providerCodes: providerCodes,
		logger:        logger,
	}
}

func (j *SyncJob) Run(ctx context.Context) {
	for _, code := range j.providerCodes {
		if err := ctx.Err(); err != nil {
			j.logger.Warnw(
				"Stopped airline synchronization",
				"error", err,
			)
			return
		}

		if err := j.service.Sync(ctx, code, airline.SearchRequest{}); err != nil {
			j.logger.Errorw(
				"Failed to sync airline service",
				"provider", code,
				"error", err,
			)
		}
	}
}
