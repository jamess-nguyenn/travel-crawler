package nuitee

import (
	"context"
	"fmt"
	"time"
	"travel-crawler/internal/provider"

	"travel-crawler/internal/redis"
)

const (
	nuiteeRateLimitKey = "ratelimit:nuitee"
	nuiteeCooldownKey  = "ratelimit:nuitee:cooldown"
)

type redisRequestLimiter struct {
	cache            *redis.Cache
	requestPerSecond int
}

func NewRedisRequestLimiter(
	cache *redis.Cache,
	requestPerSecond int,
) (provider.RateLimiter, error) {
	if requestPerSecond <= 0 || requestPerSecond > 1000000 {
		return nil, fmt.Errorf("request per second must be between 1 and 1000000")
	}

	return &redisRequestLimiter{
		cache:            cache,
		requestPerSecond: requestPerSecond,
	}, nil
}

func (l *redisRequestLimiter) Wait(ctx context.Context) error {
	return l.cache.WaitRateLimit(
		ctx,
		nuiteeRateLimitKey,
		nuiteeCooldownKey,
		l.requestPerSecond,
	)
}

func (l *redisRequestLimiter) Cooldown(
	ctx context.Context,
	duration time.Duration,
) error {
	return l.cache.SetCooldown(ctx, nuiteeCooldownKey, duration)
}

var _ provider.RateLimiter = (*redisRequestLimiter)(nil)
