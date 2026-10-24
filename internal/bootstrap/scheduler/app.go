package scheduler

import (
	"go.uber.org/zap"

	"travel-crawler/internal/job/accommodation"
	"travel-crawler/internal/kafka"
)

type App struct {
	AccommodationPublisher *accommodation.Publisher
}

func NewApp(
	kafka *kafka.Client,
	logger *zap.SugaredLogger,
) (*App, error) {
	app := &App{}

	ctr := NewContainer(kafka, logger)

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
