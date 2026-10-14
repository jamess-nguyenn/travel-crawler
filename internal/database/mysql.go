package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"travel-crawler/internal/helper"

	"travel-crawler/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func buildDsn(cfg config.MySQL) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		helper.EscapeQuery(cfg.Timezone),
	)
}

func parseLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Error
	}
}

func open(dsn string, logLevel logger.LogLevel) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             30 * time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect mysql: %w", err)
	}

	return db, nil
}

func getSQLDB(db *gorm.DB) (*sql.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get mysql sql.DB: %w", err)
	}

	return sqlDB, nil
}

func configurePool(sqlDB *sql.DB) {
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
}

func ping(sqlDB *sql.DB) error {
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping mysql: %w", err)
	}

	return nil
}

func Connect(cfg config.MySQL) (*gorm.DB, error) {
	dsn := buildDsn(cfg)
	logLevel := parseLogLevel(cfg.LogLevel)

	db, err := open(dsn, logLevel)
	if err != nil {
		return nil, err
	}

	sqlDB, err := getSQLDB(db)
	if err != nil {
		return nil, err
	}

	configurePool(sqlDB)

	if err := ping(sqlDB); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := getSQLDB(db)
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close mysql: %w", err)
	}

	return nil
}
