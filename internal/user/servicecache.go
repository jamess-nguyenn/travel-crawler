package user

import (
	"context"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/cache"
)

func (s *service) invalidateUserCache(
	ctx context.Context,
	ID uint64,
	auth0ID string,
) {
	logger := applog.FromContext(s.logger, ctx)

	cache.Invalidate(
		ctx,
		s.cache,
		CacheKeyByID(ID),
		logger,
	)

	cache.Invalidate(
		ctx,
		s.cache,
		CacheKeyByAuth0ID(auth0ID),
		logger,
	)
}
