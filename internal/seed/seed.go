package seed

import (
	"fmt"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	if err := seedUser(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	return nil
}
