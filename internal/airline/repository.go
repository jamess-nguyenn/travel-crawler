package airline

import (
	"gorm.io/gorm"
	"travel-crawler/internal/database"
)

type AirlineRepository interface {
	database.BaseRepository[Airline]
}

type repository struct {
	*database.Repository[Airline]
}

func NewRepository(db *gorm.DB) AirlineRepository {
	return &repository{
		Repository: database.NewRepository[Airline](db),
	}
}

var _ AirlineRepository = (*repository)(nil)
