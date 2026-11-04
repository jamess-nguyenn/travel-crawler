package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// waitRateLimitScript atomically checks the shared cooldown and reserves a request slot.
//
// A request is allowed only when:
//   - the shared cooldown has expired; and
//   - the next request slot is available.
//
// Redis TIME is used so all workers use the same clock.
var waitRateLimitScript = goredis.NewScript(`
local time = redis.call("TIME")
local now = (tonumber(time[1]) * 1000) + math.floor(tonumber(time[2]) / 1000)

local cooldown_until = tonumber(redis.call("GET", KEYS[2]) or "0")
if cooldown_until > now then
    return cooldown_until - now
end

local next_allowed = tonumber(redis.call("GET", KEYS[1]) or "0")
if next_allowed > now then
    return next_allowed - now
end

local interval = tonumber(ARGV[1])
local next_at = now + interval

redis.call("SET", KEYS[1], next_at, "PX", interval * 2)

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
// rateKey and cooldownKey are shared across all workers using the same
// API credential and endpoint class. Do not include country or task IDs.
func (c *Cache) WaitRateLimit(
	ctx context.Context,
	rateKey string,
	cooldownKey string,
	interval time.Duration,
) error {
	if interval <= 0 {
		return fmt.Errorf("rate limit interval must be positive")
	}

	intervalMS := interval.Milliseconds()
	if intervalMS < 1 {
		intervalMS = 1
	}

	redisRateKey := c.key(rateKey)
	redisCooldownKey := c.key(cooldownKey)

	for {
		waitMS, err := waitRateLimitScript.Run(
			ctx,
			c.client,
			[]string{redisRateKey, redisCooldownKey},
			intervalMS,
		).Int64()
		if err != nil {
			return fmt.Errorf("failed to acquire distributed rate limit: %w", err)
		}

		if waitMS <= 0 {
			return nil
		}

		timer := time.NewTimer(time.Duration(waitMS) * time.Millisecond)

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
			// Recheck Redis. Another worker may have acquired the slot
			// or a new shared cooldown may have been applied.
		}
	}
}

// SetCooldown extends the shared cooldown.
//
// If another worker has already set a longer cooldown, this method preserves the longer expiry.
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
