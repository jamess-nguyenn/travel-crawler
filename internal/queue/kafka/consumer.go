package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"

	kafkago "github.com/segmentio/kafka-go"
	"travel-crawler/internal/kafka"
)

type MessageHandler = func(ctx context.Context, payload []byte) error

type Consumer struct {
	reader *kafkago.Reader
}

func NewConsumer(
	client *kafka.Client,
	topic string,
	groupID string,
) (*Consumer, error) {
	if strings.TrimSpace(topic) == "" {
		return nil, errors.New("Kafka topic is required")
	}

	if strings.TrimSpace(groupID) == "" {
		return nil, errors.New("Kafka consumer group is required")
	}

	reader, err := client.NewReader(topic, groupID)
	if err != nil {
		return nil, fmt.Errorf("create Kafka reader: %w", err)
	}

	return &Consumer{
		reader: reader,
	}, nil
}

// Consume fetches and processes messages sequentially.
// On handler failure, it returns without committing the failed message.
func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) error {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if contextCanceled(ctx) {
				return nil
			}

			return fmt.Errorf("fetch Kafka message: %w", err)
		}

		if err := handler(ctx, message.Value); err != nil {
			if contextCanceled(ctx) {
				return nil
			}

			return fmt.Errorf(
				"handle Kafka message topic=%s partition=%d offset=%d: %w",
				message.Topic,
				message.Partition,
				message.Offset,
				err,
			)
		}

		if err := c.reader.CommitMessages(ctx, message); err != nil {
			if contextCanceled(ctx) {
				return nil
			}

			return fmt.Errorf(
				"commit Kafka message topic=%s partition=%d offset=%d: %w",
				message.Topic,
				message.Partition,
				message.Offset,
				err,
			)
		}
	}
}

func contextCanceled(ctx context.Context) bool {
	return ctx.Err() != nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

var _ interface {
	Consume(context.Context, MessageHandler) error
} = (*Consumer)(nil)
