package cache

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type Store interface {
	Set(
		ctx context.Context,
		key string,
		value any,
		ttl time.Duration,
	) error

	Get(
		ctx context.Context,
		key string,
		dest any,
	) (bool, error)

	Delete(
		ctx context.Context,
		key string,
	) error
}

func GetOrLoad[T any](
	ctx context.Context,
	store Store,
	key string,
	ttl time.Duration,
	logger *zap.SugaredLogger,
	loader func() (T, error),
) (T, bool, error) {
	var zero T

	// Cache disabled: load directly from the database.
	if store == nil {
		value, err := loader()
		return value, false, err
	}

	var cached T

	// Try to retrieve the cached value.
	found, err := store.Get(ctx, key, &cached)
	if err != nil {
		logger.Warnw(
			"Fail to get data from cache",
			"key", key,
			"error", err,
		)
	}

	if found {
		return cached, true, nil
	}

	// Cache miss: load from the database.
	value, err := loader()
	if err != nil {
		return zero, false, err
	}

	// Cache the loaded value.
	if err := store.Set(ctx, key, value, ttl); err != nil {
		logger.Warnw(
			"Fail to store data in cache",
			"key", key,
			"error", err,
		)
	}

	return value, false, nil
}

func Invalidate(
	ctx context.Context,
	store Store,
	key string,
	logger *zap.SugaredLogger,
) {
	if store == nil {
		return
	}

	if err := store.Delete(ctx, key); err != nil {
		logger.Warnw(
			"Fail to invalidate cache",
			"key", key,
			"error", err,
		)
	}
}
