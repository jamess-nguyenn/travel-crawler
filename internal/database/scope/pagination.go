package scope

import (
	"gorm.io/gorm"
	"travel-crawler/internal/helper"
)

func Limit(limit int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func Offset(offset int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func Paginate(page int, pageSize int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Offset(helper.CalculateOffset(page, pageSize)).
			Limit(pageSize)
	}
}
