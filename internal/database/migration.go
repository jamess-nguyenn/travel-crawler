package database

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

const (
	dialectDB    = "mysql"
	migrationDir = "migrations"
)

func setMigrationDialect() error {
	if err := goose.SetDialect(dialectDB); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	return nil
}

func prepareMigration(db *gorm.DB) (*sql.DB, error) {
	if err := setMigrationDialect(); err != nil {
		return nil, err
	}

	sqlDB, err := getSQLDB(db)
	if err != nil {
		return nil, err
	}

	return sqlDB, nil
}

func Up(db *gorm.DB) error {
	sqlDB, err := prepareMigration(db)
	if err != nil {
		return err
	}

	if err := goose.Up(sqlDB, migrationDir); err != nil {
		return fmt.Errorf("failed to run migrations up: %w", err)
	}

	return nil
}

func Down(db *gorm.DB, step int) error {
	if step < 1 {
		return fmt.Errorf("migration down step must be greater than 0")
	}

	sqlDB, err := prepareMigration(db)
	if err != nil {
		return err
	}

	for i := 1; i <= step; i++ {
		version, err := goose.GetDBVersion(sqlDB)
		if err != nil {
			return fmt.Errorf(
				"failed to get current migration version at step %d: %w",
				i,
				err,
			)
		}

		// Nothing to roll back.
		if version == 0 {
			return nil
		}

		if err := goose.Down(sqlDB, migrationDir); err != nil {
			return fmt.Errorf(
				"failed to roll back migration at step %d: %w",
				i,
				err,
			)
		}
	}

	return nil
}

func Reset(db *gorm.DB) error {
	sqlDB, err := prepareMigration(db)
	if err != nil {
		return err
	}

	if err := goose.Reset(sqlDB, migrationDir); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	return nil
}

func Status(db *gorm.DB) error {
	sqlDB, err := prepareMigration(db)
	if err != nil {
		return err
	}

	if err := goose.Status(sqlDB, migrationDir); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}
