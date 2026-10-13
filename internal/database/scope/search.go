package scope

import (
	"gorm.io/gorm"
	"strings"
	"travel-crawler/internal/helper"
)

func Search(column string, keyword *string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		search := helper.StringValue(keyword)
		if search == "" {
			return db
		}

		return db.Where(column+" LIKE ?", "%"+search+"%")
	}
}

func SearchColumns(columns []string, keyword *string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		search := helper.StringValue(keyword)
		if search == "" {
			return db
		}

		search = "%" + search + "%"

		conditions := make([]string, 0, len(columns))
		args := make([]any, 0, len(columns))

		for _, column := range columns {
			conditions = append(conditions, column+" LIKE ?")
			args = append(args, search)
		}

		return db.Where(
			strings.Join(conditions, " OR "),
			args...,
		)
	}
}
