package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/config"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/database"
	"travel-crawler/internal/requestcontext"
	"travel-crawler/internal/seed"
)

func main() {
	if err := run(); err != nil {
		log.Printf("seeding command failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := requestcontext.WithNewRequestContext(
		context.Background(),
	)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := applog.FromContext(
		applog.New(
			cfg.Log,
			cfg.App.Name,
		),
		ctx,
	).Named(constant.LoggerSeed)

	db, err := database.Connect(cfg.MySQL)
	if err != nil {
		logger.Errorw(
			"failed to connect database",
			"error", err,
		)

		return fmt.Errorf("failed to connect database: %w", err)
	}

	defer func() {
		if err := database.Close(db); err != nil {
			logger.Warnw(
				"failed to close database",
				"error", err,
			)
		}
	}()

	logger.Info("starting database seeding")

	if err := seed.Run(db); err != nil {
		logger.Errorw(
			"database seeding failed",
			"error", err,
		)

		return err
	}

	logger.Info("database seeding completed successfully")

	return nil
}
