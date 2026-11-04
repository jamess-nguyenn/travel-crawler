package nuitee

import (
	"context"
	"fmt"
	"time"
	"travel-crawler/internal/provider"

	"travel-crawler/internal/redis"
)

const (
	nuiteeHotelsRateLimitKey = "ratelimit:nuitee:hotels"
	nuiteeCooldownKey        = "ratelimit:nuitee:cooldown"
)

type redisRequestLimiter struct {
	cache    *redis.Cache
	interval time.Duration
}

func NewRedisRequestLimiter(
	cache *redis.Cache,
	interval time.Duration,
) (provider.RateLimiter, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("rate limit interval must be positive")
	}

	return &redisRequestLimiter{
		cache:    cache,
		interval: interval,
	}, nil
}

func (l *redisRequestLimiter) Wait(ctx context.Context) error {
	return l.cache.WaitRateLimit(
		ctx,
		nuiteeHotelsRateLimitKey,
		nuiteeCooldownKey,
		l.interval,
	)
}

func (l *redisRequestLimiter) Cooldown(
	ctx context.Context,
	duration time.Duration,
) error {
	return l.cache.SetCooldown(ctx, nuiteeCooldownKey, duration)
}

var _ provider.RateLimiter = (*redisRequestLimiter)(nil)
