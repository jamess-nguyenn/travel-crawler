package accommodation

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/data"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/job"
	"travel-crawler/internal/queue/kafka"
	"travel-crawler/internal/requestcontext"
)

type Publisher struct {
	publisher kafka.Publisher
	countries []data.Country
	logger    *zap.SugaredLogger
}

func NewPublisher(
	publisher kafka.Publisher,
	countries []data.Country,
	logger *zap.SugaredLogger,
) *Publisher {
	return &Publisher{
		publisher: publisher,
		countries: countries,
		logger:    logger,
	}
}

// Publish creates and publishes an accommodation crawl task.
// When countryCodes is empty, all configured countries are included.
// Otherwise, only the requested countries are included.
func (p *Publisher) Publish(ctx context.Context, dateFrom *time.Time, countryCodes []string) {
	traceID := helper.NewUUID()
	requestID := helper.NewUUID()

	publishCtx := requestcontext.WithRequestContext(
		ctx,
		traceID,
		requestID,
	)

	logger := applog.FromContext(p.logger, publishCtx)

	countries := p.filterCountries(countryCodes)

	total := len(countries)

	if total == 0 {
		logger.Warn("No accommodation countries to publish")
		return
	}

	messages := make([]kafka.Message, 0, total)
	tasks := make([]CrawlTask, 0, total)

	for _, country := range countries {
		if err := publishCtx.Err(); err != nil {
			logger.Warnw(
				"Stopped creating accommodation crawl tasks",
				"country", country.Code,
				"error", err,
			)
			return
		}

		taskID := helper.NewUUID()

		task := CrawlTask{
			TaskMetadata: job.TaskMetadata{
				ID:          taskID,
				Type:        CrawlTaskType,
				Version:     CrawlTaskVersion,
				RequestedAt: helper.NowUTC(),
				TraceID:     traceID,
			},
			ProviderCode: constant.ProviderCodeNuitee,
			Country:      country.Code,
			UpdatedAt:    dateFrom,
		}

		payload, err := json.Marshal(task)

		if err != nil {
			logger.Errorw(
				"Fail to marshal accommodation crawl task",
				"task", task,
				"error", err,
			)
			continue
		}

		messages = append(messages, kafka.Message{
			Key:     taskID,
			Payload: payload,
		})

		tasks = append(tasks, task)
	}

	messageTotal := len(messages)

	if messageTotal == 0 {
		logger.Warn("No accommodation crawl tasks to publish")
		return
	}

	publishStartedAt := time.Now()

	err := p.publisher.PublishBatch(publishCtx, messages)

	publishDuration := time.Since(publishStartedAt)

	publishLogger := logger.With(
		"total", messageTotal,
		"publish_duration", publishDuration,
	)

	if err != nil {
		publishLogger.Errorw(
			"Fail to publish accommodation crawl task batch",
			"error", err,
		)
		return
	}

	publishLogger.Info(
		"Publish accommodation crawl task batch",
	)

	batches := helper.Chunk(tasks, constant.KafkaBatchSize)
	totalBatch := len(batches)

	for i, batch := range batches {
		batchNumber := i + 1
		batchProgress := fmt.Sprintf("%d/%d", batchNumber, totalBatch)

		applog.WithBatch(
			logger,
			messageTotal,
			batchProgress,
			len(batch),
			batch,
		).Info("Accommodation crawl task chunk summary")
	}
}

// filterCountries returns all configured countries when countryCodes is empty.
// Otherwise, it returns only countries matching the requested country codes.
func (p *Publisher) filterCountries(countryCodes []string) []data.Country {
	if len(countryCodes) == 0 {
		return p.countries
	}

	countrySet := make(map[string]struct{}, len(countryCodes))

	for _, code := range countryCodes {
		code = strings.ToUpper(strings.TrimSpace(code))
		if code == "" {
			continue
		}

		countrySet[code] = struct{}{}
	}

	if len(countrySet) == 0 {
		return p.countries
	}

	countries := make([]data.Country, 0, len(countrySet))

	for _, country := range p.countries {
		if _, exists := countrySet[country.Code]; !exists {
			continue
		}

		countries = append(countries, country)
		delete(countrySet, country.Code)
	}

	if len(countrySet) > 0 {
		p.logger.Warnw(
			"Requested countries are not configured",
			"country_codes", helper.MapKeys(countrySet),
		)
	}

	return countries
}
