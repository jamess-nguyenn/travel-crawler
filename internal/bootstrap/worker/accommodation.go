package worker

import (
	"fmt"
	"travel-crawler/internal/accommodation"
	accomprovider "travel-crawler/internal/provider/accommodation"
	"travel-crawler/internal/queue/kafka"
	accomworker "travel-crawler/internal/worker/accommodation"
)

type AccommodationModule struct{}

func (AccommodationModule) Init(app *App, ctr *Container) error {
	consumer, err := kafka.NewConsumer(
		ctr.kafka,
		kafka.TopicAccommodationCrawl,
		kafka.GroupAccommodationWorker,
	)
	if err != nil {
		return fmt.Errorf("create accommodation consumer: %w", err)
	}

	retryProducer, err := kafka.NewProducer(
		ctr.kafka,
		kafka.TopicAccommodationCrawlRetry,
	)
	if err != nil {
		_ = consumer.Close()
		return fmt.Errorf("create accommodation retry producer: %w", err)
	}

	dlqProducer, err := kafka.NewProducer(
		ctr.kafka,
		kafka.TopicAccommodationCrawlDLQ,
	)
	if err != nil {
		_ = consumer.Close()
		return fmt.Errorf("create accommodation DLQ producer: %w", err)
	}

	providers, err := accomprovider.New(ctr.http, *ctr.cfg, ctr.cache, ctr.logger)
	if err != nil {
		return fmt.Errorf("initialize accommodation providers: %w", err)
	}

	service := accommodation.NewSyncService(
		ctr.txManager,
		ctr.accommodationRepo,
		providers,
		ctr.logger,
	)

	processor := accomworker.NewAccommodationProcessor(
		service,
		ctr.logger,
	)

	handler := accomworker.NewRetryHandler(
		processor,
		retryProducer,
		dlqProducer,
	)

	app.AccommodationWorker = accomworker.NewWorker(
		consumer,
		handler,
	)

	return nil
}
