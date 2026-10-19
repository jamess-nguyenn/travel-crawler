package accommodation

import (
	"gorm.io/gorm"
	"travel-crawler/internal/database"
)

type AccommodationRepository interface {
	database.BaseRepository[Accommodation]
}

type repository struct {
	*database.Repository[Accommodation]
}

func NewRepository(db *gorm.DB) AccommodationRepository {
	return &repository{
		Repository: database.NewRepository[Accommodation](db),
	}
}

var _ AccommodationRepository = (*repository)(nil)
