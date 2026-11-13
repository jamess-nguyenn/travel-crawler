package nuitee

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/provider"
)

const defaultCooldown = 2 * time.Second

type RequestGuard struct {
	limiter         provider.RateLimiter
	defaultCooldown time.Duration
}

func NewRequestGuard(
	limiter provider.RateLimiter,
) *RequestGuard {
	return &RequestGuard{
		limiter:         limiter,
		defaultCooldown: defaultCooldown,
	}
}

// Wait acquires a shared request slot before calling Nuitee.
func (g *RequestGuard) Wait(ctx context.Context) error {
	if g.limiter == nil {
		return nil
	}

	if err := g.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("wait for Nuitee rate limiter: %w", err)
	}

	return nil
}

// HandleError handles Nuitee-specific HTTP errors.
func (g *RequestGuard) HandleError(
	ctx context.Context,
	err error,
) error {
	if err == nil {
		return nil
	}

	var httpErr *httpclient.HTTPError

	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTooManyRequests {
		return err
	}

	if g.limiter == nil {
		return err
	}

	cooldown := httpErr.RetryAfter()
	if cooldown <= 0 {
		cooldown = g.defaultCooldown
	}

	if cooldownErr := g.limiter.Cooldown(ctx, cooldown); cooldownErr != nil {
		return errors.Join(
			err,
			fmt.Errorf("activate Nuitee shared cooldown: %w", cooldownErr),
		)
	}

	return err
}
