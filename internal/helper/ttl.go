package helper

import (
	"math/rand/v2"
	"time"
)

const (
	day  = 24 * time.Hour
	week = 7 * day
)

const (
	shortTTLLimit  = time.Hour
	mediumTTLLimit = day
	longTTLLimit   = week
)

const (
	shortTTLMinPercent    = 90 // <= 1 hour
	mediumTTLMinPercent   = 95 // > 1 hour && <= 1 day
	longTTLMinPercent     = 97 // > 1 day && <= 1 week
	veryLongTTLMinPercent = 98 // > 1 week

	maxJitterPercent = 99
)

func jitterPercent(ttl time.Duration) int {
	var minPercent int

	switch {
	case ttl <= shortTTLLimit:
		minPercent = shortTTLMinPercent
	case ttl <= mediumTTLLimit:
		minPercent = mediumTTLMinPercent
	case ttl <= longTTLLimit:
		minPercent = longTTLMinPercent
	default:
		minPercent = veryLongTTLMinPercent
	}

	return rand.IntN(maxJitterPercent-minPercent+1) + minPercent
}

func Jitter(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return ttl
	}

	percent := time.Duration(jitterPercent(ttl))

	// Use integer arithmetic without overflowing ttl * percent.
	whole := ttl / 100 * percent
	remainder := (ttl%100*percent + 99) / 100

	return whole + remainder
}

func MinuteTTL(count int) time.Duration {
	return time.Duration(count) * time.Minute
}

func HourTTL(count int) time.Duration {
	return Jitter(time.Duration(count) * time.Hour)
}

func DayTTL(count int) time.Duration {
	return Jitter(time.Duration(count) * day)
}

func WeekTTL(count int) time.Duration {
	return Jitter(time.Duration(count) * week)
}
