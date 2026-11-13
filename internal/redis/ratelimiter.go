package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// waitRateLimitScript atomically:
// 1. Reads Redis server time.
// 2. Checks the shared cooldown.
// 3. Checks the next available request slot.
// 4. Reserves the next slot before returning success.
//
// All workers share the same reservation state.
//
// Requests are evenly distributed rather than burst-based.
var waitRateLimitScript = goredis.NewScript(`
local time = redis.call("TIME")

local now_us = (tonumber(time[1]) * 1000000) + tonumber(time[2])
local now_ms = math.floor(now_us / 1000)

-- Check shared cooldown.
local cooldown_until = tonumber(redis.call("GET", KEYS[2]) or "0")

if cooldown_until > now_ms then
    return (cooldown_until - now_ms) * 1000
end

-- Check next available request slot.
local next_allowed = tonumber(redis.call("GET", KEYS[1]) or "0")

if next_allowed > now_us then
    return next_allowed - now_us
end

local request_per_second = tonumber(ARGV[1])

-- Calculate interval in microseconds.
-- Ceiling prevents exceeding the configured rate.
local interval_us = math.ceil(1000000 / request_per_second)

-- Reserve the slot atomically BEFORE returning success.
local next_at = now_us + interval_us

-- Keep the rate key alive for at least twice the interval.
local ttl_ms = math.max(1, math.ceil(interval_us * 2 / 1000))

redis.call("SET", KEYS[1], next_at, "PX", ttl_ms)

return 0
`)

// setCooldownScript extends the shared cooldown without shortening
// an existing cooldown.
var setCooldownScript = goredis.NewScript(`
local time = redis.call("TIME")

local now = (tonumber(time[1]) * 1000) + math.floor(tonumber(time[2]) / 1000)

local duration = tonumber(ARGV[1])
local requested_until = now + duration

local existing_until = tonumber(redis.call("GET", KEYS[1]) or "0")

if existing_until > requested_until then
    requested_until = existing_until
end

local ttl = requested_until - now

redis.call("SET", KEYS[1], requested_until, "PX", ttl)

return ttl
`)

// WaitRateLimit waits until a request slot is available.
//
// requestPerSecond is the maximum configured request frequency.
//
// rateKey and cooldownKey must be shared across all workers
// using the same API credential and quota.
func (c *Cache) WaitRateLimit(
	ctx context.Context,
	rateKey string,
	cooldownKey string,
	requestPerSecond int,
) error {
	if requestPerSecond <= 0 || requestPerSecond > 1000000 {
		return fmt.Errorf("request per second must be between 1 and 1000000")
	}

	redisRateKey := c.key(rateKey)
	redisCooldownKey := c.key(cooldownKey)

	for {
		waitUS, err := waitRateLimitScript.Run(
			ctx,
			c.client,
			[]string{redisRateKey, redisCooldownKey},
			requestPerSecond,
		).Int64()

		if err != nil {
			return fmt.Errorf(
				"failed to acquire distributed rate limit: %w",
				err,
			)
		}

		if waitUS <= 0 {
			// Slot has already been reserved atomically by Redis.
			return nil
		}

		timer := time.NewTimer(time.Duration(waitUS) * time.Microsecond)

		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}

			return ctx.Err()

		case <-timer.C:
			// Recheck Redis because another worker may have:
			// - Reserved the next available slot.
			// - Activated a shared cooldown.
		}
	}
}

// SetCooldown extends the shared cooldown.
//
// If another worker has already set a longer cooldown,
// this method preserves the longer expiry.
func (c *Cache) SetCooldown(
	ctx context.Context,
	cooldownKey string,
	duration time.Duration,
) error {
	if duration <= 0 {
		return nil
	}

	durationMS := duration.Milliseconds()
	if durationMS < 1 {
		durationMS = 1
	}

	redisKey := c.key(cooldownKey)

	if err := setCooldownScript.Run(
		ctx,
		c.client,
		[]string{redisKey},
		durationMS,
	).Err(); err != nil {
		return fmt.Errorf("failed to set shared cooldown for key %q: %w", redisKey, err)
	}

	return nil
}
