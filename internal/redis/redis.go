package redis

import (
	"context"
	"fmt"
	"time"
	"travel-crawler/internal/config"

	goredis "github.com/redis/go-redis/v9"
)

const connectionTimeout = 5 * time.Second

func open(cfg config.Redis) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:     cfg.Address(),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func ping(client *goredis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

func Connect(cfg config.Redis) (*goredis.Client, error) {
	client := open(cfg)

	if err := ping(client); err != nil {
		_ = client.Close()

		return nil, err
	}

	return client, nil
}
