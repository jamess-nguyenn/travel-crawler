package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"travel-crawler/internal/config"

	kafkago "github.com/segmentio/kafka-go"
)

const (
	defaultConnectionTimeout = 5 * time.Second
	defaultReadTimeout       = 10 * time.Second
	defaultWriteTimeout      = 10 * time.Second
)

type Client struct {
	brokers []string

	mu      sync.Mutex
	writers []*kafkago.Writer
	readers []*kafkago.Reader
	closed  bool
}

// Connect validates the Kafka configuration and checks connectivity.
func Connect(ctx context.Context, cfg config.Kafka) (*Client, error) {
	brokers := parseBrokers(cfg.Brokers)
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must not be empty")
	}

	dialer := &kafkago.Dialer{
		Timeout: defaultConnectionTimeout,
	}

	// Verify connectivity against the first configured broker.
	conn, err := dialer.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return nil, fmt.Errorf(
			"connect to kafka broker %q: %w",
			brokers[0],
			err,
		)
	}

	if err := conn.Close(); err != nil {
		return nil, fmt.Errorf("close kafka connection: %w", err)
	}

	return &Client{
		brokers: brokers,
	}, nil
}

// NewWriter creates a synchronous Kafka producer for a topic.
func (c *Client) NewWriter(topic string) (*kafkago.Writer, error) {
	if strings.TrimSpace(topic) == "" {
		return nil, errors.New("kafka topic must not be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil, errors.New("kafka client is closed")
	}

	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(c.brokers...),
		Topic:        topic,
		Balancer:     &kafkago.Hash{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
		WriteTimeout: defaultWriteTimeout,
	}

	c.writers = append(c.writers, writer)

	return writer, nil
}

// NewReader creates a Kafka consumer using a consumer group.
func (c *Client) NewReader(
	topic string,
	groupID string,
) (*kafkago.Reader, error) {
	if strings.TrimSpace(topic) == "" {
		return nil, errors.New("kafka topic must not be empty")
	}

	if strings.TrimSpace(groupID) == "" {
		return nil, errors.New("kafka consumer group ID must not be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil, errors.New("kafka client is closed")
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        c.brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1,
		MaxBytes:       10e6,
		MaxWait:        defaultReadTimeout,
		CommitInterval: 0,
	})

	c.readers = append(c.readers, reader)

	return reader, nil
}

// Brokers returns a copy of the configured broker addresses.
func (c *Client) Brokers() []string {
	return append([]string(nil), c.brokers...)
}

// Close releases all readers and writers created by this client.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	var errs []error

	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close kafka reader: %w", err))
		}
	}

	for _, writer := range c.writers {
		if err := writer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close kafka writer: %w", err))
		}
	}

	c.readers = nil
	c.writers = nil

	return errors.Join(errs...)
}

func parseBrokers(value string) []string {
	parts := strings.Split(value, ",")
	brokers := make([]string, 0, len(parts))

	for _, part := range parts {
		if broker := strings.TrimSpace(part); broker != "" {
			brokers = append(brokers, broker)
		}
	}

	return brokers
}
