package airport

import (
	"gorm.io/gorm"
	"travel-crawler/internal/database"
)

type AirportRepository interface {
	database.BaseRepository[Airport]
}

type repository struct {
	*database.Repository[Airport]
}

func NewRepository(db *gorm.DB) AirportRepository {
	return &repository{
		Repository: database.NewRepository[Airport](db),
	}
}

var _ AirportRepository = (*repository)(nil)
