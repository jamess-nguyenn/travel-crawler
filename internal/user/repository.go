package user

import (
	"context"
	"gorm.io/gorm"
	"travel-crawler/internal/database"
)

type UserRepository interface {
	database.BaseRepository[User]

	GetByEmailOrAuth0ID(
		ctx context.Context,
		email string,
		auth0ID string,
	) ([]User, error)
}

type repository struct {
	*database.Repository[User]
}

func NewRepository(db *gorm.DB) UserRepository {
	return &repository{
		Repository: database.NewRepository[User](db),
	}
}

func (r *repository) GetByEmailOrAuth0ID(
	ctx context.Context,
	email string,
	auth0ID string,
) ([]User, error) {
	var users []User

	if err := r.Query(ctx).
		Where(
			"email = ? OR auth0_id = ?",
			email,
			auth0ID,
		).
		Limit(2).
		Find(&users).
		Error; err != nil {
		return nil, err
	}

	return users, nil
}
