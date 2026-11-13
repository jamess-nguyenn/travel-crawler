package airport

import (
	"context"

	"go.uber.org/zap"

	"travel-crawler/internal/airport"
)

type SyncJob struct {
	service airport.SyncService
	logger  *zap.SugaredLogger
}

func NewSyncJob(
	service airport.SyncService,
	logger *zap.SugaredLogger,
) *SyncJob {
	return &SyncJob{
		service: service,
		logger:  logger,
	}
}

func (j *SyncJob) Run(ctx context.Context) {
	if err := j.service.Sync(ctx); err != nil {
		j.logger.Errorw(
			"Failed to sync airport service",
			"error", err,
		)
	}
}
