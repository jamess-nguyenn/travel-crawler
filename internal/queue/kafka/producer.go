package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"

	kafkago "github.com/segmentio/kafka-go"
	"travel-crawler/internal/kafka"
)

type Message struct {
	Key     string
	Payload []byte
}

type Publisher interface {
	Publish(ctx context.Context, key string, payload []byte) error
	PublishBatch(ctx context.Context, messages []Message) error
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
	return p.PublishBatch(ctx, []Message{
		{
			Key:     key,
			Payload: payload,
		},
	})
}

func (p *Producer) PublishBatch(ctx context.Context, messages []Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	total := len(messages)

	if total == 0 {
		return nil
	}

	kafkaMessages := make([]kafkago.Message, 0, total)

	for _, message := range messages {
		kafkaMessages = append(kafkaMessages, kafkago.Message{
			Key:   []byte(message.Key),
			Value: message.Payload,
		})
	}

	if err := p.writer.WriteMessages(ctx, kafkaMessages...); err != nil {
		return fmt.Errorf(
			"write %d Kafka messages to topic %s: %w",
			total,
			p.topic,
			err,
		)
	}

	return nil
}

func (p *Producer) Topic() string {
	return p.topic
}
