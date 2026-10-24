package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"

	kafkago "github.com/segmentio/kafka-go"
	"travel-crawler/internal/kafka"
)

type Publisher interface {
	Publish(ctx context.Context, key string, payload []byte) error
}

type Producer struct {
	writer *kafkago.Writer
	topic  string
}

func NewProducer(
	client *kafka.Client,
	topic string,
) (*Producer, error) {
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return nil, errors.New("Kafka topic is required")
	}

	writer, err := client.NewWriter(topic)
	if err != nil {
		return nil, fmt.Errorf("create Kafka writer for topic %q: %w", topic, err)
	}

	return &Producer{
		writer: writer,
		topic:  topic,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, key string, payload []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	message := kafkago.Message{
		Key:   []byte(key),
		Value: payload,
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("write Kafka message to topic %s: %w", p.topic, err)
	}

	return nil
}

func (p *Producer) Topic() string {
	return p.topic
}
