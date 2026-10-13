package helper

import "time"

func NowUTC() time.Time {
	return time.Now().UTC()
}

func StartOfDay(t time.Time) time.Time {
	t = t.UTC()
	y, m, d := t.Date()

	return time.Date(
		y,
		m,
		d,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func EndOfDay(t time.Time) time.Time {
	t = t.UTC()
	y, m, d := t.Date()

	return time.Date(
		y,
		m,
		d,
		23,
		59,
		59,
		999999999,
		time.UTC,
	)
}
