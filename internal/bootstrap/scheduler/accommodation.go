package scheduler

import (
	"fmt"
	"travel-crawler/internal/data"
	"travel-crawler/internal/job/accommodation"
	"travel-crawler/internal/queue/kafka"
)

type AccommodationModule struct{}

func (AccommodationModule) Init(app *App, ctr *Container) error {
	producer, err := kafka.NewProducer(
		ctr.kafka,
		kafka.TopicAccommodationCrawl,
	)
	if err != nil {
		return fmt.Errorf("create accommodation producer: %w", err)
	}

	countries, err := data.Countries()
	if err != nil {
		return fmt.Errorf("get countries: %w", err)
	}

	app.AccommodationPublisher = accommodation.NewPublisher(
		producer,
		countries,
		ctr.logger,
	)

	return nil
}
