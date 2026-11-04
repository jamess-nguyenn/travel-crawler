package accommodation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/job"
	accomjob "travel-crawler/internal/job/accommodation"
	"travel-crawler/internal/queue/kafka"
	"travel-crawler/internal/worker"
	"travel-crawler/internal/worker/retry"
)

type FailurePublisher interface {
	Publish(ctx context.Context, key string, payload []byte) error
}

type RetryHandler struct {
	processor    Processor
	dlqPublisher FailurePublisher
	maxAttempt   int
	now          func() time.Time
}

func NewRetryHandler(
	processor Processor,
	dlqPublisher FailurePublisher,
) *RetryHandler {
	return &RetryHandler{
		processor:    processor,
		dlqPublisher: dlqPublisher,
		maxAttempt:   retry.MaxAttempt,
		now:          time.Now,
	}
}

func (h *RetryHandler) Handle(ctx context.Context, payload []byte) error {
	var task accomjob.CrawlTask

	if err := json.Unmarshal(payload, &task); err != nil {
		return h.publishDLQ(
			ctx,
			payload,
			nil,
			0,
			"invalid_json",
			fmt.Errorf("decode crawl task: %w", err),
			h.now(),
		)
	}

	if err := validateTask(task); err != nil {
		return h.publishDLQ(
			ctx,
			payload,
			task,
			0,
			"invalid_task",
			err,
			h.now(),
		)
	}

	firstFailedAt := time.Time{}

	for attempt := 1; attempt <= h.maxAttempt; attempt++ {
		if err := ctx.Err(); err != nil {
			// Cancellation is not a task failure. Return the error so
			// the consumer does not commit the original message.
			return err
		}

		err := h.processor.Process(ctx, task)
		if err == nil {
			return nil
		}

		now := h.now()
		if firstFailedAt.IsZero() {
			firstFailedAt = now
		}

		// Do not route a task to the DLQ merely because the worker
		// is shutting down or its parent context has been cancelled.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !retry.IsRetryable(err) {
			return h.publishDLQ(
				ctx,
				payload,
				task,
				attempt,
				"classified_permanent",
				err,
				firstFailedAt,
			)
		}

		if attempt == h.maxAttempt {
			return h.publishDLQ(
				ctx,
				payload,
				task,
				attempt,
				"retries_exhausted",
				err,
				firstFailedAt,
			)
		}

		delay := retry.Delay(err, attempt)

		if err := retry.Wait(ctx, delay); err != nil {
			return err
		}
	}

	return nil
}

func validateTask(task accomjob.CrawlTask) error {
	dimension := constant.DimensionAccommodation

	if task.GetID() == "" {
		return &worker.InvalidTaskError{Dimension: dimension, Reason: "id is required"}
	}

	if task.GetType() != accomjob.CrawlTaskType {
		return &worker.InvalidTaskError{Dimension: dimension, Reason: "unsupported task type"}
	}

	if task.GetVersion() != accomjob.CrawlTaskVersion {
		return &worker.InvalidTaskError{Dimension: dimension, Reason: "unsupported task version"}
	}

	if strings.TrimSpace(task.ProviderCode) == "" {
		return &worker.InvalidTaskError{Dimension: dimension, Reason: "provider code is required"}
	}

	return nil
}

func (h *RetryHandler) publishDLQ(
	ctx context.Context,
	payload []byte,
	metadata job.Metadata,
	attempt int,
	errorClass string,
	cause error,
	firstFailedAt time.Time,
) error {
	if h.dlqPublisher == nil {
		return errors.New("DLQ publisher is not configured")
	}

	now := h.now()

	taskID := ""
	taskType := ""
	traceID := helper.NewUUID()

	if metadata != nil {
		taskID = metadata.GetID()
		taskType = metadata.GetType()
		traceID = metadata.GetTraceID()
	}

	if firstFailedAt.IsZero() {
		firstFailedAt = now
	}

	message := kafka.FailureMessage{
		TaskID:          taskID,
		TaskType:        taskType,
		OriginalPayload: payload,
		Attempt:         attempt,
		MaxAttempt:      h.maxAttempt,
		ErrorClass:      errorClass,
		ErrorMessage:    cause.Error(),
		FirstFailedAt:   firstFailedAt,
		FailedAt:        now,
		TraceID:         traceID,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal DLQ message: %w", err)
	}

	key := taskID
	if key == "" {
		key = helper.NewUUID()
	}

	return h.dlqPublisher.Publish(ctx, key, data)
}
