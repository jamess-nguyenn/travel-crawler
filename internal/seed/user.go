package seed

import (
	"time"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/model"
	"travel-crawler/internal/user"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedUser(db *gorm.DB) error {
	now := time.Now().UTC()

	countries := []user.User{
		{
			Model: model.Model{
				ID:        constant.SystemUserID,
				CreatedAt: now,
				CreatedBy: constant.SystemUserID,
				UpdatedAt: now,
				UpdatedBy: constant.SystemUserID,
			},
			Email:           "system@travel-crawler.xyz",
			Auth0ID:         "system|00000000000000000001001",
			EmailVerifiedAt: &now,
			FirstName:       "System",
			LastName:        "User",
			FullName:        "System User",
			Status:          user.StatusActive,
		},
	}

	return db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "auth0_id"},
			},
			DoNothing: true,
		}).
		Create(&countries).
		Error
}
