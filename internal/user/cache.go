package user

import (
	"time"
	"travel-crawler/internal/cache"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/redis"
)

const (
	cachePrefix = "user"
	cacheTTL    = 15 // days
)

func CacheKeyByID(ID uint64) string {
	return helper.JoinKey(
		cachePrefix,
		"id",
		ID,
	)
}

func CacheKeyByAuth0ID(auth0ID string) string {
	return helper.JoinKey(
		cachePrefix,
		"auth0_id",
		auth0ID,
	)
}

func CacheDuration() time.Duration {
	return helper.DayTTL(cacheTTL)
}

var _ cache.Store = (*redis.Cache)(nil)
