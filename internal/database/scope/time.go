package scope

import (
	"gorm.io/gorm"
	"time"
)

func TimeFrom(column string, value time.Time) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column+" >= ?", value)
	}
}

func TimeTo(column string, value time.Time) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column+" <= ?", value)
	}
}

func TimeBetween(column string, from, to time.Time) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(
			column+" BETWEEN ? AND ?",
			from,
			to,
		)
	}
}

func TimeRange(column string, from, to *time.Time) Scope {
	return func(db *gorm.DB) *gorm.DB {
		if from != nil {
			db = db.Where(column+" >= ?", *from)
		}

		if to != nil {
			db = db.Where(column+" <= ?", *to)
		}

		return db
	}
}
