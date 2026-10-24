package accommodation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	accomjob "travel-crawler/internal/job/accommodation"
	"travel-crawler/internal/queue/kafka"
)

type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

type Processor interface {
	Process(ctx context.Context, task accomjob.CrawlTask) error
}

type Worker struct {
	consumer  Consumer
	processor Processor
}

func NewWorker(
	consumer Consumer,
	processor Processor,
) *Worker {
	return &Worker{
		consumer:  consumer,
		processor: processor,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	return w.consumer.Consume(ctx, w.Handle)
}

func (w *Worker) Handle(ctx context.Context, payload []byte) error {
	var task accomjob.CrawlTask

	if err := json.Unmarshal(payload, &task); err != nil {
		return fmt.Errorf("decode accommodation crawl task: %w", err)
	}

	if err := w.validateTask(task); err != nil {
		return err
	}

	if err := w.processor.Process(ctx, task); err != nil {
		return fmt.Errorf("process crawl task %s: %w", task.ID, err)
	}

	return nil
}

func (w *Worker) validateTask(task accomjob.CrawlTask) error {
	if strings.TrimSpace(task.ID) == "" {
		return fmt.Errorf("crawl task ID is required")
	}

	if task.Type != accomjob.CrawlTaskType {
		return fmt.Errorf("unsupported crawl task type %q", task.Type)
	}

	if task.Version != accomjob.CrawlTaskVersion {
		return fmt.Errorf("unsupported crawl task version %d", task.Version)
	}

	if strings.TrimSpace(task.ProviderCode) == "" {
		return fmt.Errorf("provider code is required")
	}

	return nil
}
