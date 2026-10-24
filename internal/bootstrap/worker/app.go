package worker

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"travel-crawler/internal/config"
	"travel-crawler/internal/kafka"
	"travel-crawler/internal/redis"
	"travel-crawler/internal/worker/accommodation"
)

type App struct {
	AccommodationWorker *accommodation.Worker
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
	}

	for _, module := range modules {
		if err := module.Init(app, ctr); err != nil {
			return nil, err
		}
	}

	return app, nil
}
