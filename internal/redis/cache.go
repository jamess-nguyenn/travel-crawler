package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"travel-crawler/internal/config"
	"travel-crawler/internal/helper"

	goredis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client *goredis.Client
	prefix string
}

func NewCache(cfg config.Redis, prefix string) (*Cache, error) {
	client, err := Connect(cfg)
	if err != nil {
		return nil, err
	}

	return &Cache{
		client: client,
		prefix: prefix,
	}, nil
}

func (c *Cache) key(key string) string {
	if c.prefix == "" {
		return key
	}

	return helper.JoinKey(c.prefix, key)
}

func (c *Cache) Set(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	redisKey := c.key(key)

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf(
			"failed to marshal redis value for key %q: %w",
			redisKey,
			err,
		)
	}

	if err := c.client.Set(ctx, redisKey, data, ttl).Err(); err != nil {
		return fmt.Errorf(
			"failed to set redis key %q: %w",
			redisKey,
			err,
		)
	}

	return nil
}

func (c *Cache) Get(
	ctx context.Context,
	key string,
	dest any,
) (bool, error) {
	redisKey := c.key(key)

	data, err := c.client.Get(ctx, redisKey).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf(
			"failed to get redis key %q: %w",
			redisKey,
			err,
		)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf(
			"failed to unmarshal redis value for key %q: %w",
			redisKey,
			err,
		)
	}

	return true, nil
}

func (c *Cache) Delete(
	ctx context.Context,
	key string,
) error {
	redisKey := c.key(key)

	if err := c.client.Del(ctx, redisKey).Err(); err != nil {
		return fmt.Errorf(
			"failed to delete redis key %q: %w",
			redisKey,
			err,
		)
	}

	return nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}
