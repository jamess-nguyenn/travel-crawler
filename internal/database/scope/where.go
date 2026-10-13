package scope

import "gorm.io/gorm"

func WhereIn(column string, values any) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column+" IN ?", values)
	}
}

func WhereNotIn(column string, values any) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column+" NOT IN ?", values)
	}
}
