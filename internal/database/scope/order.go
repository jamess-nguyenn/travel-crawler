package scope

import (
	"gorm.io/gorm"
	"travel-crawler/internal/request"
)

func Order(column string, sortOrder string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(column + " " + sortOrder)
	}
}

func OrderBy(column string) Scope {
	return Order(column, request.SortOrderAsc)
}

func OrderByDesc(column string) Scope {
	return Order(column, request.SortOrderDesc)
}
