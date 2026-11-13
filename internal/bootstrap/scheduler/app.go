package scheduler

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"travel-crawler/internal/config"
	airlinejob "travel-crawler/internal/job/airline"
	airportjob "travel-crawler/internal/job/airport"
	"travel-crawler/internal/redis"

	"travel-crawler/internal/job/accommodation"
	"travel-crawler/internal/kafka"
)

type App struct {
	AccommodationPublisher *accommodation.Publisher
	AirlineSyncJob         *airlinejob.SyncJob
	AirportSyncJob         *airportjob.SyncJob
}

func NewApp(
	db *gorm.DB,
	cfg *config.Config,
	cache *redis.Cache,
	kafka *kafka.Client,
	logger *zap.SugaredLogger,
) (*App, error) {
	app := &App{}

	ctr := NewContainer(db, cfg, cache, kafka, logger)

	modules := []Module{
		AccommodationModule{},
		AirlineModule{},
		AirportModule{},
	}

	for _, module := range modules {
		if err := module.Init(app, ctr); err != nil {
			return nil, err
		}
	}

	return app, nil
}
