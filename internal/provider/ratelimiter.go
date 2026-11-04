package provider

import (
	"context"
	"time"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
	Cooldown(ctx context.Context, duration time.Duration) error
}
