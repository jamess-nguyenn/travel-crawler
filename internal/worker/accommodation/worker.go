package accommodation

import (
	"context"
	accomjob "travel-crawler/internal/job/accommodation"
	"travel-crawler/internal/queue/kafka"
)

type Handler interface {
	Handle(ctx context.Context, payload []byte) error
}

type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

type Processor interface {
	Process(ctx context.Context, task accomjob.CrawlTask) error
}

type Worker struct {
	consumer Consumer
	handler  Handler
}

func NewWorker(
	consumer Consumer,
	handler Handler,
) *Worker {
	return &Worker{
		consumer: consumer,
		handler:  handler,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	return w.consumer.Consume(ctx, w.handler.Handle)
}
